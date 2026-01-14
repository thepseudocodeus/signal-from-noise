package database

import (
	"database/sql"
	"fmt"
	"log"
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
}

// FileFilters represents filters for querying files
type FileFilters struct {
	ProductionRequestID string
	DateStart          *time.Time
	DateEnd            *time.Time
	Categories         []string // "email", "claim", "other"
	ExcludePrivileged  bool
	Page               int
	PageSize           int
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
	// Ensure directory exists
	dir := filepath.Dir(dbPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create database directory: %w", err)
	}

	db, err := sql.Open("sqlite3", dbPath+"?_journal_mode=WAL")
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	// Test connection
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	database := &DB{db: db}

	// Initialize schema
	if err := database.initSchema(); err != nil {
		return nil, fmt.Errorf("failed to initialize schema: %w", err)
	}

	// Check if we need to seed data
	count, err := database.getFileCount()
	if err != nil {
		return nil, fmt.Errorf("failed to check file count: %w", err)
	}

	if count == 0 {
		log.Println("Database is empty, seeding mock data...")
		if err := database.SeedMockData(); err != nil {
			return nil, fmt.Errorf("failed to seed mock data: %w", err)
		}
		log.Println("Mock data seeded successfully")
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
		created_at TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP
	);

	CREATE INDEX IF NOT EXISTS idx_files_category ON files(category);
	CREATE INDEX IF NOT EXISTS idx_files_date ON files(date);
	CREATE INDEX IF NOT EXISTS idx_files_directory ON files(directory);
	CREATE INDEX IF NOT EXISTS idx_files_privileged ON files(privileged);
	CREATE INDEX IF NOT EXISTS idx_files_duplicate_hash ON files(duplicate_hash);

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
	log.Printf("SearchFiles called with filters: %+v", filters)

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

	// Exclude privileged
	if filters.ExcludePrivileged {
		whereClause += " AND privileged = 0"
	}

	// Get total count
	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM files WHERE %s", whereClause)
	var totalCount int
	err := d.db.QueryRow(countQuery, args...).Scan(&totalCount)
	if err != nil {
		return nil, fmt.Errorf("failed to get file count: %w", err)
	}

	log.Printf("Total files matching filters: %d", totalCount)

	// Calculate pagination
	if filters.Page < 1 {
		filters.Page = 1
	}
	if filters.PageSize < 1 {
		filters.PageSize = 50
	}
	offset := (filters.Page - 1) * filters.PageSize
	totalPages := (totalCount + filters.PageSize - 1) / filters.PageSize

	// Get files
	query := fmt.Sprintf(`
		SELECT id, path, directory, category, date, size, privileged, duplicate_hash, file_name
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
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan file: %w", err)
		}

		f.Date, err = time.Parse(time.RFC3339, dateStr)
		if err != nil {
			return nil, fmt.Errorf("failed to parse date: %w", err)
		}

		files = append(files, f)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating rows: %w", err)
	}

	log.Printf("Returning %d files (page %d of %d)", len(files), filters.Page, totalPages)

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
		SELECT id, path, directory, category, date, size, privileged, duplicate_hash, file_name
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
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan file: %w", err)
		}

		f.Date, err = time.Parse(time.RFC3339, dateStr)
		if err != nil {
			return nil, fmt.Errorf("failed to parse date: %w", err)
		}

		files = append(files, f)
	}

	return files, nil
}
