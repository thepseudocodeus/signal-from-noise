package main

import (
	"context"
	"fmt"
	"os"

	"signal-from-noise/config"
	"signal-from-noise/datalake"
)

// App struct
type App struct {
	ctx            context.Context
	config         *config.Config
	dataLakeService *datalake.DataLakeService
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
