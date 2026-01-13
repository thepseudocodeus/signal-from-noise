package datalake

import (
	"fmt"
	"os"
	"path/filepath"
)

// DataLakeService handles data lake operations
type DataLakeService struct {
	rootPath string
}

// NewDataLakeService creates a new data lake service
func NewDataLakeService(rootPath string) *DataLakeService {
	return &DataLakeService{
		rootPath: rootPath,
	}
}

// ValidateDataLake checks if the data lake path is accessible
func (d *DataLakeService) ValidateDataLake() error {
	if d.rootPath == "" {
		return fmt.Errorf("data lake path is empty")
	}

	// Check if path exists
	info, err := os.Stat(d.rootPath)
	if os.IsNotExist(err) {
		return fmt.Errorf("data lake path does not exist: %s", d.rootPath)
	}
	if err != nil {
		return fmt.Errorf("error checking data lake path: %w", err)
	}

	// Check if it's a directory
	if !info.IsDir() {
		return fmt.Errorf("data lake path is not a directory: %s", d.rootPath)
	}

	// Check if readable
	file, err := os.Open(d.rootPath)
	if err != nil {
		return fmt.Errorf("data lake path is not readable: %w", err)
	}
	file.Close()

	return nil
}

// DiscoverEmailFiles discovers all email files in the email final directory
// For demo: returns file list (first N files for performance)
func (d *DataLakeService) DiscoverEmailFiles(emailPath string, maxFiles int) ([]string, error) {
	if emailPath == "" {
		return nil, fmt.Errorf("email path is empty")
	}

	// Check if email directory exists
	if _, err := os.Stat(emailPath); os.IsNotExist(err) {
		return []string{}, nil // Empty directory is OK
	}

	var files []string
	count := 0

	err := filepath.Walk(emailPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil // Skip files with errors
		}

		// Only process files (not directories)
		if !info.IsDir() {
			// For demo: common email file extensions
			ext := filepath.Ext(path)
			if isEmailFile(ext) {
				files = append(files, path)
				count++

				// Limit files for demo performance
				if maxFiles > 0 && count >= maxFiles {
					return filepath.SkipAll
				}
			}
		}

		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("error walking email directory: %w", err)
	}

	return files, nil
}

// isEmailFile checks if file extension indicates an email file
func isEmailFile(ext string) bool {
	emailExts := []string{".parquet", ".parq", ".csv", ".json", ".pst"}
	for _, e := range emailExts {
		if ext == e {
			return true
		}
	}
	return false
}

// GetEmailFileCount returns the count of email files (for status)
func (d *DataLakeService) GetEmailFileCount(emailPath string) (int, error) {
	files, err := d.DiscoverEmailFiles(emailPath, 0) // 0 = no limit for counting
	if err != nil {
		return 0, err
	}
	return len(files), nil
}
