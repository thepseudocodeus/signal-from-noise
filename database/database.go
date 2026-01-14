package database

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

// File represents a file in the database
type File struct {
	ID            int64     `json:"id"`
	Path          string    `json:"path"`
	Directory     string    `json:"directory"`
	Category      string    `json:"category"` // "email", "claim", "other"
	Date          time.Time `json:"date"`
	Size          int64     `json:"size"`
	Privileged    bool      `json:"privileged"`
	DuplicateHash string    `json:"duplicate_hash"`
	FileName      string    `json:"file_name"`
	// Email-specific fields (NULL for non-email files)
	Subject     string `json:"subject"`      // Email subject line
	FromEmail   string `json:"from_email"`   // Sender email address
	ToEmail     string `json:"to_email"`     // Recipient email address (first/main)
	Sentiment   string `json:"sentiment"`    // "positive", "negative", "neutral", "unknown"
	IsInternal bool  `json:"is_internal"`   // true if internal email, false if external
	Topic       string `json:"topic"`        // Extracted topic from subject
}

// FileFilters represents filters for querying files
// Each filter incrementally reduces the problem space
type FileFilters struct {
	ProductionRequestID string
	DateStart          *time.Time
	DateEnd            *time.Time
	Categories         []string // "email", "claim", "other"
	ExcludePrivileged  bool
	// New incremental filters
	Topics    []string // Topics extracted from email subjects
	People    []string // Email addresses (from FROM or TO fields)
	Sentiment string   // "positive", "negative", "neutral", "unknown", "all"
	// People filter options
	PeopleFilterType string // "internal", "external", "specific", "all"
	Page             int
	PageSize         int
}

// FileResult represents paginated file results
type FileResult struct {
	Files      []File `json:"files"`
	TotalCount int    `json:"total_count"`
	Page       int    `json:"page"`
	PageSize   int    `json:"page_size"`
	TotalPages int    `json:"total_pages"`
}

// DB wraps the database connection
type DB struct {
	db *sql.DB
}

// NewDB creates a new database connection
func NewDB(dbPath string) (*DB, error) {
	dir := filepath.Dir(dbPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create database directory: %w", err)
	}

	db, err := sql.Open("sqlite3", dbPath+"?_journal_mode=WAL")
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	database := &DB{db: db}
	if err := database.initSchema(); err != nil {
		return nil, fmt.Errorf("failed to initialize schema: %w", err)
	}

	// Check if we need to seed data
	count, err := database.getFileCount()
	if err != nil {
		return nil, fmt.Errorf("failed to check file count: %w", err)
	}

	if count == 0 {
		if err := database.SeedMockData(); err != nil {
			return nil, fmt.Errorf("failed to seed mock data: %w", err)
		}
	}

	return database, nil
}

// Close closes the database connection
func (d *DB) Close() error {
	return d.db.Close()
}

// initSchema creates the database tables
func (d *DB) initSchema() error {
	schema := `
	CREATE TABLE IF NOT EXISTS files (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		path TEXT NOT NULL,
		directory TEXT NOT NULL,
		category TEXT NOT NULL,
		date TEXT NOT NULL,
		size INTEGER NOT NULL,
		privileged INTEGER NOT NULL DEFAULT 0,
		duplicate_hash TEXT,
		file_name TEXT NOT NULL,
		-- Email-specific fields (NULL for non-email files)
		subject TEXT,
		from_email TEXT,
		to_email TEXT,
		sentiment TEXT,
		is_internal INTEGER DEFAULT 0,
		topic TEXT,
		created_at TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP
	);

	CREATE INDEX IF NOT EXISTS idx_files_category ON files(category);
	CREATE INDEX IF NOT EXISTS idx_files_date ON files(date);
	CREATE INDEX IF NOT EXISTS idx_files_directory ON files(directory);
	CREATE INDEX IF NOT EXISTS idx_files_privileged ON files(privileged);
	CREATE INDEX IF NOT EXISTS idx_files_duplicate_hash ON files(duplicate_hash);
	-- New indexes for filtering
	CREATE INDEX IF NOT EXISTS idx_files_topic ON files(topic);
	CREATE INDEX IF NOT EXISTS idx_files_sentiment ON files(sentiment);
	CREATE INDEX IF NOT EXISTS idx_files_is_internal ON files(is_internal);
	CREATE INDEX IF NOT EXISTS idx_files_from_email ON files(from_email);
	CREATE INDEX IF NOT EXISTS idx_files_to_email ON files(to_email);

	CREATE TABLE IF NOT EXISTS production_requests (
		id TEXT PRIMARY KEY,
		title TEXT NOT NULL,
		description TEXT NOT NULL,
		created_at TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP
	);
	`

	_, err := d.db.Exec(schema)
	return err
}

// getFileCount returns the total number of files
func (d *DB) getFileCount() (int, error) {
	var count int
	err := d.db.QueryRow("SELECT COUNT(*) FROM files").Scan(&count)
	return count, err
}

// SearchFiles searches for files based on filters
func (d *DB) SearchFiles(filters FileFilters) (*FileResult, error) {
	if d.db == nil {
		return nil, fmt.Errorf("database connection required")
	}

	// ASSUMPTION: Page number must be positive for pagination to work correctly
	// Negative or zero pages would cause incorrect offset calculations
	if filters.Page < 1 {
		filters.Page = 1
	}
	// ASSUMPTION: Page size must be positive to return results
	// Zero or negative page size would return no results or cause errors
	if filters.PageSize < 1 {
		filters.PageSize = 50
	}

	// Build WHERE clause
	whereClause := "1=1"
	args := []interface{}{}

	// Date range filter
	if filters.DateStart != nil {
		whereClause += " AND date >= ?"
		args = append(args, filters.DateStart.Format(time.RFC3339))
	}
	if filters.DateEnd != nil {
		whereClause += " AND date <= ?"
		args = append(args, filters.DateEnd.Format(time.RFC3339))
	}

	// Category filter
	if len(filters.Categories) > 0 {
		placeholders := ""
		for i, cat := range filters.Categories {
			if i > 0 {
				placeholders += ","
			}
			placeholders += "?"
			args = append(args, cat)
		}
		whereClause += fmt.Sprintf(" AND category IN (%s)", placeholders)
	}

	// Topic filter (incremental complexity reduction)
	// ASSUMPTION: Topics exist in database and can filter email files
	// If topics selected, only files with matching topics are returned
	if len(filters.Topics) > 0 {
		placeholders := ""
		for i, topic := range filters.Topics {
			if i > 0 {
				placeholders += ","
			}
			placeholders += "?"
			args = append(args, topic)
		}
		whereClause += fmt.Sprintf(" AND topic IN (%s)", placeholders)
	}

	// People filter (incremental complexity reduction)
	// ASSUMPTION: People filter reduces result set by email addresses
	// Internal/external classification precomputed for performance
	if filters.PeopleFilterType != "" && filters.PeopleFilterType != "all" {
		if filters.PeopleFilterType == "internal" {
			whereClause += " AND is_internal = 1"
		} else if filters.PeopleFilterType == "external" {
			whereClause += " AND is_internal = 0"
		} else if filters.PeopleFilterType == "specific" && len(filters.People) > 0 {
			placeholders := ""
			for i, email := range filters.People {
				if i > 0 {
					placeholders += " OR "
				}
				placeholders += "(from_email = ? OR to_email = ?)"
				args = append(args, email, email)
			}
			whereClause += fmt.Sprintf(" AND (%s)", placeholders)
		}
	}

	// Sentiment filter (incremental complexity reduction)
	// ASSUMPTION: Sentiment values are valid and filter email files
	// "all" means no sentiment filter applied
	if filters.Sentiment != "" && filters.Sentiment != "all" {
		whereClause += " AND sentiment = ?"
		args = append(args, filters.Sentiment)
	}

	// Exclude privileged
	if filters.ExcludePrivileged {
		whereClause += " AND privileged = 0"
	}

	// Get total count
	// ASSUMPTION: SQL query will execute successfully and return a count
	// If this fails, the database schema or query structure is invalid
	countQuery := "SELECT COUNT(*) FROM files WHERE " + whereClause
	var totalCount int
	err := d.db.QueryRow(countQuery, args...).Scan(&totalCount)
	if err != nil {
		return nil, fmt.Errorf("failed to get file count: %w", err)
	}

	offset := (filters.Page - 1) * filters.PageSize
	totalPages := (totalCount + filters.PageSize - 1) / filters.PageSize

	// Get files
	// ASSUMPTION: SQL query structure matches database schema
	// Column names must exist in the files table
	// Include all fields for frontend display
	query := fmt.Sprintf(`
		SELECT id, path, directory, category, date, size, privileged, duplicate_hash, file_name,
		       subject, from_email, to_email, sentiment, is_internal, topic
		FROM files
		WHERE %s
		ORDER BY date DESC
		LIMIT ? OFFSET ?
	`, whereClause)

	args = append(args, filters.PageSize, offset)
	rows, err := d.db.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to query files: %w", err)
	}
	defer rows.Close()

	var files []File
	for rows.Next() {
		var f File
		var dateStr string
		var subject, fromEmail, toEmail, sentiment, topic sql.NullString
		var isInternal sql.NullBool

		// ASSUMPTION: Row structure matches SELECT statement
		// All columns must be scannable into the File struct
		err := rows.Scan(
			&f.ID,
			&f.Path,
			&f.Directory,
			&f.Category,
			&dateStr,
			&f.Size,
			&f.Privileged,
			&f.DuplicateHash,
			&f.FileName,
			&subject,
			&fromEmail,
			&toEmail,
			&sentiment,
			&isInternal,
			&topic,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan file: %w", err)
		}

		f.Date, err = time.Parse(time.RFC3339, dateStr)
		if err != nil {
			return nil, fmt.Errorf("failed to parse date: %w", err)
		}

		// Handle nullable email fields
		if subject.Valid {
			f.Subject = subject.String
		}
		if fromEmail.Valid {
			f.FromEmail = fromEmail.String
		}
		if toEmail.Valid {
			f.ToEmail = toEmail.String
		}
		if sentiment.Valid {
			f.Sentiment = sentiment.String
		}
		if topic.Valid {
			f.Topic = topic.String
		}
		if isInternal.Valid {
			f.IsInternal = isInternal.Bool
		}

		files = append(files, f)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating rows: %w", err)
	}

	return &FileResult{
		Files:      files,
		TotalCount: totalCount,
		Page:       filters.Page,
		PageSize:   filters.PageSize,
		TotalPages: totalPages,
	}, nil
}

// GetFileByID retrieves a file by its ID
func (d *DB) GetFileByID(id int64) (*File, error) {
	query := `
		SELECT id, path, directory, category, date, size, privileged, duplicate_hash, file_name
		FROM files
		WHERE id = ?
	`

	var f File
	var dateStr string
	err := d.db.QueryRow(query, id).Scan(
		&f.ID,
		&f.Path,
		&f.Directory,
		&f.Category,
		&dateStr,
		&f.Size,
		&f.Privileged,
		&f.DuplicateHash,
		&f.FileName,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get file: %w", err)
	}

	f.Date, err = time.Parse(time.RFC3339, dateStr)
	if err != nil {
		return nil, fmt.Errorf("failed to parse date: %w", err)
	}

	return &f, nil
}

// GetFilesByIDs retrieves multiple files by their IDs
func (d *DB) GetFilesByIDs(ids []int64) ([]File, error) {
	if len(ids) == 0 {
		return []File{}, nil
	}

	placeholders := ""
	args := make([]interface{}, len(ids))
	for i, id := range ids {
		if i > 0 {
			placeholders += ","
		}
		placeholders += "?"
		args[i] = id
	}

	query := fmt.Sprintf(`
		SELECT id, path, directory, category, date, size, privileged, duplicate_hash, file_name,
		       subject, from_email, to_email, sentiment, is_internal, topic
		FROM files
		WHERE id IN (%s)
	`, placeholders)

	rows, err := d.db.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to query files: %w", err)
	}
	defer rows.Close()

	var files []File
	for rows.Next() {
		var f File
		var dateStr string
		var subject, fromEmail, toEmail, sentiment, topic sql.NullString
		var isInternal sql.NullBool
		err := rows.Scan(
			&f.ID,
			&f.Path,
			&f.Directory,
			&f.Category,
			&dateStr,
			&f.Size,
			&f.Privileged,
			&f.DuplicateHash,
			&f.FileName,
			&subject,
			&fromEmail,
			&toEmail,
			&sentiment,
			&isInternal,
			&topic,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan file: %w", err)
		}

		f.Date, err = time.Parse(time.RFC3339, dateStr)
		if err != nil {
			return nil, fmt.Errorf("failed to parse date: %w", err)
		}

		// Handle nullable fields
		if subject.Valid {
			f.Subject = subject.String
		}
		if fromEmail.Valid {
			f.FromEmail = fromEmail.String
		}
		if toEmail.Valid {
			f.ToEmail = toEmail.String
		}
		if sentiment.Valid {
			f.Sentiment = sentiment.String
		}
		if topic.Valid {
			f.Topic = topic.String
		}
		if isInternal.Valid {
			f.IsInternal = isInternal.Bool
		}

		files = append(files, f)
	}

	return files, nil
}
