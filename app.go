package main

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"signal-from-noise/assert"
	"signal-from-noise/config"
	"signal-from-noise/database"
	"signal-from-noise/datalake"
	"signal-from-noise/logging"
)

// App struct
type App struct {
	ctx            context.Context
	config         *config.Config
	dataLakeService *datalake.DataLakeService
	db             *database.DB
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx

	// Load configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		// Log error but don't fail startup (for demo, we'll handle gracefully)
		fmt.Printf("Warning: Could not load config: %v\n", err)
	} else {
		a.config = cfg
		if cfg != nil {
			a.dataLakeService = datalake.NewDataLakeService(cfg.GetDataLakePath())
		}
	}

	// Initialize database
	dbPath := filepath.Join(os.TempDir(), "signal-from-noise.db")
	fmt.Printf("Initializing database at: %s\n", dbPath)
	
	db, err := database.NewDB(dbPath)
	if err != nil {
		fmt.Printf("Error initializing database: %v\n", err)
		// Continue without database for now
	} else {
		a.db = db
		fmt.Println("Database initialized successfully")
	}
}

// Greet returns a greeting for the given name
func (a *App) Greet(name string) string {
	return fmt.Sprintf("Hello %s, It's show time!", name)
}

// GetDataLakeStatus checks if the data lake is accessible
// Returns status string: "available", "unavailable", or error message
func (a *App) GetDataLakeStatus() (string, error) {
	if a.dataLakeService == nil {
		return "unavailable", fmt.Errorf("data lake service not initialized")
	}

	err := a.dataLakeService.ValidateDataLake()
	if err != nil {
		return "unavailable", err
	}

	return "available", nil
}

// GetEmailFileCount returns the count of email files in the email final directory
func (a *App) GetEmailFileCount() (int, error) {
	if a.config == nil {
		return 0, fmt.Errorf("configuration not loaded")
	}

	if a.dataLakeService == nil {
		return 0, fmt.Errorf("data lake service not initialized")
	}

	emailPath := a.config.GetEmailFinalPath()
	count, err := a.dataLakeService.GetEmailFileCount(emailPath)
	if err != nil {
		return 0, fmt.Errorf("error getting email file count: %w", err)
	}

	return count, nil
}

// ListEmailFiles returns a list of email files (limited for demo performance)
// maxFiles: maximum number of files to return (0 = unlimited, but not recommended for large datasets)
func (a *App) ListEmailFiles(maxFiles int) ([]string, error) {
	if a.config == nil {
		return nil, fmt.Errorf("configuration not loaded")
	}

	if a.dataLakeService == nil {
		return nil, fmt.Errorf("data lake service not initialized")
	}

	emailPath := a.config.GetEmailFinalPath()

	// For demo: default to 50 files if maxFiles is 0
	if maxFiles == 0 {
		maxFiles = 50
	}

	files, err := a.dataLakeService.DiscoverEmailFiles(emailPath, maxFiles)
	if err != nil {
		return nil, fmt.Errorf("error discovering email files: %w", err)
	}

	return files, nil
}

// GetEmailPath returns the path to the email final directory
func (a *App) GetEmailPath() (string, error) {
	if a.config == nil {
		return "", fmt.Errorf("configuration not loaded")
	}

	emailPath := a.config.GetEmailFinalPath()

	// Verify path exists
	if _, err := os.Stat(emailPath); os.IsNotExist(err) {
		return emailPath, nil // Return path even if doesn't exist (empty directory is OK)
	}

	return emailPath, nil
}

// GetDataLakePath returns the root data lake path
func (a *App) GetDataLakePath() (string, error) {
	if a.config == nil {
		return "", fmt.Errorf("configuration not loaded")
	}

	return a.config.GetDataLakePath(), nil
}

// FileSearchRequest represents a file search request from the frontend
type FileSearchRequest struct {
	ProductionRequestID string   `json:"production_request_id"`
	DateStart          string   `json:"date_start"` // ISO 8601 format
	DateEnd            string   `json:"date_end"`   // ISO 8601 format
	Categories         []string `json:"categories"` // "email", "claim", "other"
	ExcludePrivileged  bool     `json:"exclude_privileged"`
	Page               int      `json:"page"`
	PageSize           int      `json:"page_size"`
}

// FileSearchResult represents the result of a file search
type FileSearchResult struct {
	Files      []FileInfo `json:"files"`
	TotalCount int        `json:"total_count"`
	Page       int        `json:"page"`
	PageSize   int        `json:"page_size"`
	TotalPages int        `json:"total_pages"`
}

// FileInfo represents file information returned to frontend
type FileInfo struct {
	ID            int64  `json:"id"`
	Path          string `json:"path"`
	Directory     string `json:"directory"`
	Category      string `json:"category"`
	Date          string `json:"date"` // ISO 8601 format
	Size          int64  `json:"size"`
	Privileged    bool   `json:"privileged"`
	DuplicateHash string `json:"duplicate_hash"`
	FileName      string `json:"file_name"`
}

// SearchFiles searches for files based on the provided filters
// Assumption: Database is initialized and accessible
// Assumption: Request contains valid filter parameters
func (a *App) SearchFiles(req FileSearchRequest) (*FileSearchResult, error) {
	op := logging.StartOperation("SearchFiles", map[string]interface{}{
		"production_request_id": req.ProductionRequestID,
		"date_start":            req.DateStart,
		"date_end":              req.DateEnd,
		"categories":             req.Categories,
		"exclude_privileged":    req.ExcludePrivileged,
		"page":                  req.Page,
		"page_size":             req.PageSize,
	})
	defer op.EndOperation()

	// ASSUMPTION: Database must be initialized before queries can execute
	// If db is nil, the application startup failed or database initialization was skipped
	assert.ThatNotNil(a.db, "database must be initialized for file searches to work")

	logging.LogAssumption("Request contains valid filter parameters", map[string]interface{}{
		"has_date_start": req.DateStart != "",
		"has_date_end":   req.DateEnd != "",
		"category_count": len(req.Categories),
	})

	// Parse dates
	// ASSUMPTION: Date strings are in RFC3339 format if provided
	// If parsing fails, the frontend sent invalid date format
	var dateStart, dateEnd *time.Time
	if req.DateStart != "" {
		t, err := time.Parse(time.RFC3339, req.DateStart)
		if err != nil {
			logging.LogError("SearchFiles", err, map[string]interface{}{
				"operation":  "parse_date_start",
				"date_string": req.DateStart,
			})
			return nil, fmt.Errorf("invalid date_start format: %w", err)
		}
		dateStart = &t
	}
	if req.DateEnd != "" {
		t, err := time.Parse(time.RFC3339, req.DateEnd)
		if err != nil {
			logging.LogError("SearchFiles", err, map[string]interface{}{
				"operation":  "parse_date_end",
				"date_string": req.DateEnd,
			})
			return nil, fmt.Errorf("invalid date_end format: %w", err)
		}
		dateEnd = &t
	}

	// ASSUMPTION: If both dates provided, start must be before or equal to end
	// Invalid date ranges would return no results or cause logical errors
	if dateStart != nil && dateEnd != nil {
		assert.That(!dateEnd.Before(*dateStart), "date range must be valid: start date must be before or equal to end date")
	}

	// Normalize categories (frontend uses "Email", "Claims", "Other", backend uses lowercase)
	// ASSUMPTION: Frontend sends valid category names that can be mapped
	// Unknown categories are passed through, but may not match database values
	categories := make([]string, len(req.Categories))
	for i, cat := range req.Categories {
		switch cat {
		case "Email":
			categories[i] = "email"
		case "Claims":
			categories[i] = "claim"
		case "Other":
			categories[i] = "other"
		default:
			categories[i] = cat
		}
	}

	// Set defaults
	// ASSUMPTION: Page and page size must be positive for pagination
	// Zero or negative values would cause incorrect SQL queries
	if req.Page < 1 {
		req.Page = 1
	}
	if req.PageSize < 1 {
		req.PageSize = 50
	}

	logging.LogInvariant("pagination_parameters", map[string]interface{}{
		"page":      req.Page,
		"page_size": req.PageSize,
	})

	// Search database
	filters := database.FileFilters{
		ProductionRequestID: req.ProductionRequestID,
		DateStart:          dateStart,
		DateEnd:            dateEnd,
		Categories:         categories,
		ExcludePrivileged:  req.ExcludePrivileged,
		Page:               req.Page,
		PageSize:           req.PageSize,
	}

	result, err := a.db.SearchFiles(filters)
	if err != nil {
		logging.LogError("SearchFiles", err, map[string]interface{}{
			"operation": "database_search",
		})
		return nil, fmt.Errorf("database search failed: %w", err)
	}

	// ASSUMPTION: Database result is valid and contains files array
	// If result is nil, the database query failed unexpectedly
	assert.ThatNotNil(result, "database search must return a valid result")

	// Convert to frontend format
	// ASSUMPTION: All file dates can be formatted as RFC3339
	// If formatting fails, the date structure is invalid
	files := make([]FileInfo, len(result.Files))
	for i, f := range result.Files {
		files[i] = FileInfo{
			ID:            f.ID,
			Path:          f.Path,
			Directory:     f.Directory,
			Category:      f.Category,
			Date:          f.Date.Format(time.RFC3339),
			Size:          f.Size,
			Privileged:    f.Privileged,
			DuplicateHash: f.DuplicateHash,
			FileName:      f.FileName,
		}
	}

	// ASSUMPTION: Result structure is consistent (files count matches result structure)
	// Mismatch indicates data transformation error
	assert.That(len(files) == len(result.Files), "converted files count must match database result count")

	op.EndOperationWithResult(map[string]interface{}{
		"files_returned": len(files),
		"total_count":    result.TotalCount,
		"page":           result.Page,
		"total_pages":    result.TotalPages,
	})

	return &FileSearchResult{
		Files:      files,
		TotalCount: result.TotalCount,
		Page:       result.Page,
		PageSize:   result.PageSize,
		TotalPages: result.TotalPages,
	}, nil
}

// CreateZipRequest represents a request to create a zip file
type CreateZipRequest struct {
	ProductionRequestID string  `json:"production_request_id"`
	FileIDs            []int64 `json:"file_ids"`
}

// CreateZipResponse represents the response from zip creation
type CreateZipResponse struct {
	ZipPath string `json:"zip_path"`
	Success bool   `json:"success"`
	Message string `json:"message"`
}

// CreateZip creates a zip file containing the selected files
// Assumption: Database is initialized and can retrieve file information
// Assumption: At least one file ID is provided
// Assumption: Production request ID is valid
func (a *App) CreateZip(req CreateZipRequest) (*CreateZipResponse, error) {
	op := logging.StartOperation("CreateZip", map[string]interface{}{
		"production_request_id": req.ProductionRequestID,
		"file_count":           len(req.FileIDs),
	})
	defer op.EndOperation()

	// ASSUMPTION: Database must be initialized to retrieve file information
	// If db is nil, we cannot look up files to include in zip
	assert.ThatNotNil(a.db, "database must be initialized to create zip files")

	// ASSUMPTION: At least one file must be selected for zip creation
	// Empty zip files serve no purpose and indicate user error
	assert.That(len(req.FileIDs) > 0, "at least one file must be selected for zip creation")

	// ASSUMPTION: Production request ID must be provided for file naming
	// Empty ID would create incorrectly named zip files
	assert.That(req.ProductionRequestID != "", "production request ID must be provided for zip file naming")

	// Create output directory in temp
	// ASSUMPTION: System temp directory is writable
	// If this fails, the application cannot create output files
	outputDir := filepath.Join(os.TempDir(), "signal-from-noise-zips")
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		logging.LogError("CreateZip", err, map[string]interface{}{
			"operation":   "create_output_directory",
			"output_dir":  outputDir,
		})
		return &CreateZipResponse{
			Success: false,
			Message: fmt.Sprintf("Failed to create output directory: %v", err),
		}, fmt.Errorf("failed to create output directory: %w", err)
	}

	// ASSUMPTION: Output directory exists and is writable
	// If directory creation succeeded, we should be able to write files
	logging.LogCheckpoint("CreateZip", map[string]interface{}{
		"output_directory": outputDir,
		"files_to_zip":     len(req.FileIDs),
	})

	// Create zip file
	// ASSUMPTION: All file IDs exist in database and can be retrieved
	// If files don't exist, zip creation will fail
	zipPath, err := a.db.CreateZipFile(req.ProductionRequestID, req.FileIDs, outputDir)
	if err != nil {
		logging.LogError("CreateZip", err, map[string]interface{}{
			"operation": "create_zip_file",
		})
		return &CreateZipResponse{
			Success: false,
			Message: fmt.Sprintf("Failed to create zip: %v", err),
		}, fmt.Errorf("failed to create zip: %w", err)
	}

	// ASSUMPTION: Zip file path is valid and file exists after creation
	// If file doesn't exist, creation silently failed
	assert.That(zipPath != "", "zip file path must be non-empty after creation")

	op.EndOperationWithResult(map[string]interface{}{
		"zip_path":    zipPath,
		"file_count":  len(req.FileIDs),
		"success":     true,
	})

	return &CreateZipResponse{
		ZipPath: zipPath,
		Success: true,
		Message: fmt.Sprintf("Zip file created: %s", zipPath),
	}, nil
}
