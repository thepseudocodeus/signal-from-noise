package database

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"signal-from-noise/assert"
	"signal-from-noise/logging"

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
// Assumption: dbPath must be a valid path string (can be empty for in-memory)
func NewDB(dbPath string) (*DB, error) {
	op := logging.StartOperation("NewDB", map[string]interface{}{
		"db_path": dbPath,
	})
	defer op.EndOperation()

	// ASSUMPTION: Directory path must be valid for file operations
	// If this fails, we cannot create the database file
	dir := filepath.Dir(dbPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		logging.LogError("NewDB", err, map[string]interface{}{
			"operation": "create_directory",
			"path":      dir,
		})
		return nil, fmt.Errorf("failed to create database directory: %w", err)
	}

	// ASSUMPTION: SQLite driver is available and can open connections
	// If this fails, the database system is not properly configured
	db, err := sql.Open("sqlite3", dbPath+"?_journal_mode=WAL")
	if err != nil {
		logging.LogError("NewDB", err, map[string]interface{}{
			"operation": "open_database",
			"path":      dbPath,
		})
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	// ASSUMPTION: Database connection is alive and responsive
	// If ping fails, the database file is corrupted or inaccessible
	if err := db.Ping(); err != nil {
		logging.LogError("NewDB", err, map[string]interface{}{
			"operation": "ping_database",
		})
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	database := &DB{db: db}

	// ASSUMPTION: Schema can be initialized on this database
	// If this fails, the database structure is invalid
	if err := database.initSchema(); err != nil {
		logging.LogError("NewDB", err, map[string]interface{}{
			"operation": "init_schema",
		})
		return nil, fmt.Errorf("failed to initialize schema: %w", err)
	}

	// Check if we need to seed data
	// ASSUMPTION: We can query the database to check if it's empty
	count, err := database.getFileCount()
	if err != nil {
		logging.LogError("NewDB", err, map[string]interface{}{
			"operation": "check_file_count",
		})
		return nil, fmt.Errorf("failed to check file count: %w", err)
	}

	if count == 0 {
		logging.LogCheckpoint("NewDB", map[string]interface{}{
			"action": "seeding_mock_data",
			"reason": "database_empty",
		})
		if err := database.SeedMockData(); err != nil {
			logging.LogError("NewDB", err, map[string]interface{}{
				"operation": "seed_mock_data",
			})
			return nil, fmt.Errorf("failed to seed mock data: %w", err)
		}
		logging.LogResult("SeedMockData", 0, map[string]interface{}{
			"status": "success",
		})
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
// Assumption: Database connection is valid and queries can execute
// Assumption: Filters contain valid values (dates are valid, categories are valid strings)
func (d *DB) SearchFiles(filters FileFilters) (*FileResult, error) {
	op := logging.StartOperation("SearchFiles", map[string]interface{}{
		"date_start":         filters.DateStart,
		"date_end":           filters.DateEnd,
		"categories":         filters.Categories,
		"exclude_privileged": filters.ExcludePrivileged,
		"page":               filters.Page,
		"page_size":          filters.PageSize,
	})
	defer op.EndOperation()

	// ASSUMPTION: Database connection exists and is valid
	// If db is nil, the database was not properly initialized
	assert.ThatNotNil(d.db, "database connection must exist for queries to execute")

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

	// Exclude privileged
	if filters.ExcludePrivileged {
		whereClause += " AND privileged = 0"
	}

	// Get total count
	// ASSUMPTION: SQL query will execute successfully and return a count
	// If this fails, the database schema or query structure is invalid
	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM files WHERE %s", whereClause)
	logging.LogQuery(countQuery, map[string]interface{}{
		"args_count": len(args),
		"filters":    whereClause,
	})

	var totalCount int
	err := d.db.QueryRow(countQuery, args...).Scan(&totalCount)
	if err != nil {
		logging.LogError("SearchFiles", err, map[string]interface{}{
			"operation": "count_query",
			"query":     countQuery,
		})
		return nil, fmt.Errorf("failed to get file count: %w", err)
	}

	// ASSUMPTION: Count must be non-negative (database constraint)
	// Negative counts indicate data corruption or query error
	assert.That(totalCount >= 0, "file count must be non-negative (database integrity check)")

	logging.LogResult("SearchFiles", totalCount, map[string]interface{}{
		"filters_applied": len(filters.Categories),
	})

	// Calculate pagination
	// ASSUMPTION: Offset calculation is correct for pagination
	// Offset = (page - 1) * page_size ensures we skip the right number of records
	offset := (filters.Page - 1) * filters.PageSize
	totalPages := (totalCount + filters.PageSize - 1) / filters.PageSize

	// ASSUMPTION: Offset must be non-negative for SQL LIMIT/OFFSET to work
	// Negative offset would cause SQL error
	assert.That(offset >= 0, "offset must be non-negative for SQL pagination to work")

	logging.LogCheckpoint("SearchFiles", map[string]interface{}{
		"total_count":  totalCount,
		"total_pages":  totalPages,
		"current_page": filters.Page,
		"offset":       offset,
	})

	// Get files
	// ASSUMPTION: SQL query structure matches database schema
	// Column names must exist in the files table
	query := fmt.Sprintf(`
		SELECT id, path, directory, category, date, size, privileged, duplicate_hash, file_name
		FROM files
		WHERE %s
		ORDER BY date DESC
		LIMIT ? OFFSET ?
	`, whereClause)

	args = append(args, filters.PageSize, offset)
	logging.LogQuery(query, map[string]interface{}{
		"limit":  filters.PageSize,
		"offset": offset,
	})

	rows, err := d.db.Query(query, args...)
	if err != nil {
		logging.LogError("SearchFiles", err, map[string]interface{}{
			"operation": "query_files",
			"query":     query,
		})
		return nil, fmt.Errorf("failed to query files: %w", err)
	}
	defer rows.Close()

	var files []File
	for rows.Next() {
		var f File
		var dateStr string
		
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
		)
		if err != nil {
			logging.LogError("SearchFiles", err, map[string]interface{}{
				"operation": "scan_row",
			})
			return nil, fmt.Errorf("failed to scan file: %w", err)
		}

		// ASSUMPTION: Date string is in RFC3339 format as stored
		// If parsing fails, the database contains invalid date data
		f.Date, err = time.Parse(time.RFC3339, dateStr)
		if err != nil {
			logging.LogError("SearchFiles", err, map[string]interface{}{
				"operation": "parse_date",
				"date_str":  dateStr,
			})
			return nil, fmt.Errorf("failed to parse date: %w", err)
		}

		// ASSUMPTION: File ID must be positive (database constraint)
		// Zero or negative IDs indicate data corruption
		assert.That(f.ID > 0, "file ID must be positive (database integrity check)")

		files = append(files, f)
	}

	if err := rows.Err(); err != nil {
		logging.LogError("SearchFiles", err, map[string]interface{}{
			"operation": "iterate_rows",
		})
		return nil, fmt.Errorf("error iterating rows: %w", err)
	}

	// ASSUMPTION: Number of files returned should not exceed page size
	// More files than page size indicates pagination logic error
	assert.That(len(files) <= filters.PageSize, "returned files must not exceed page size (pagination logic check)")

	op.EndOperationWithResult(map[string]interface{}{
		"files_returned": len(files),
		"total_count":    totalCount,
		"page":           filters.Page,
		"total_pages":    totalPages,
	})

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
