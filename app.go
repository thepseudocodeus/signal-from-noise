package main

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"signal-from-noise/config"
	"signal-from-noise/database"
	"signal-from-noise/datalake"
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
func (a *App) SearchFiles(req FileSearchRequest) (*FileSearchResult, error) {
	fmt.Printf("SearchFiles called: %+v\n", req)

	if a.db == nil {
		return nil, fmt.Errorf("database not initialized")
	}

	// Parse dates
	var dateStart, dateEnd *time.Time
	if req.DateStart != "" {
		t, err := time.Parse(time.RFC3339, req.DateStart)
		if err != nil {
			return nil, fmt.Errorf("invalid date_start format: %w", err)
		}
		dateStart = &t
	}
	if req.DateEnd != "" {
		t, err := time.Parse(time.RFC3339, req.DateEnd)
		if err != nil {
			return nil, fmt.Errorf("invalid date_end format: %w", err)
		}
		dateEnd = &t
	}

	// Normalize categories (frontend uses "Email", "Claims", "Other", backend uses lowercase)
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
	if req.Page < 1 {
		req.Page = 1
	}
	if req.PageSize < 1 {
		req.PageSize = 50
	}

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
		return nil, fmt.Errorf("database search failed: %w", err)
	}

	// Convert to frontend format
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

	fmt.Printf("SearchFiles returning %d files (page %d of %d)\n", len(files), result.Page, result.TotalPages)

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
func (a *App) CreateZip(req CreateZipRequest) (*CreateZipResponse, error) {
	fmt.Printf("CreateZip called: production_request=%s, file_count=%d\n", req.ProductionRequestID, len(req.FileIDs))

	if a.db == nil {
		return &CreateZipResponse{
			Success: false,
			Message: "Database not initialized",
		}, fmt.Errorf("database not initialized")
	}

	if len(req.FileIDs) == 0 {
		return &CreateZipResponse{
			Success: false,
			Message: "No files selected",
		}, fmt.Errorf("no files selected")
	}

	// Create output directory in temp
	outputDir := filepath.Join(os.TempDir(), "signal-from-noise-zips")
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return &CreateZipResponse{
			Success: false,
			Message: fmt.Sprintf("Failed to create output directory: %v", err),
		}, fmt.Errorf("failed to create output directory: %w", err)
	}

	// Create zip file
	zipPath, err := a.db.CreateZipFile(req.ProductionRequestID, req.FileIDs, outputDir)
	if err != nil {
		return &CreateZipResponse{
			Success: false,
			Message: fmt.Sprintf("Failed to create zip: %v", err),
		}, fmt.Errorf("failed to create zip: %w", err)
	}

	fmt.Printf("Zip file created successfully: %s\n", zipPath)

	return &CreateZipResponse{
		ZipPath: zipPath,
		Success: true,
		Message: fmt.Sprintf("Zip file created: %s", zipPath),
	}, nil
}
