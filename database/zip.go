package database

import (
	"archive/zip"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"
)

// CreateZipFile creates a zip file containing the specified files
// Returns the path to the created zip file
func (d *DB) CreateZipFile(productionRequestID string, fileIDs []int64, outputDir string) (string, error) {
	if len(fileIDs) == 0 {
		return "", fmt.Errorf("no files selected")
	}

	// Get files from database
	files, err := d.GetFilesByIDs(fileIDs)
	if err != nil {
		return "", fmt.Errorf("failed to get files: %w", err)
	}

	if len(files) == 0 {
		return "", fmt.Errorf("no files found for selected IDs")
	}

	// Ensure output directory exists
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return "", fmt.Errorf("failed to create output directory: %w", err)
	}

	// Create zip file name: {productionRequestID}_{timestamp}.zip
	timestamp := time.Now().Format("20060102_150405")
	zipFileName := fmt.Sprintf("%s_%s.zip", productionRequestID, timestamp)
	zipPath := filepath.Join(outputDir, zipFileName)

	// Create zip file
	zipFile, err := os.Create(zipPath)
	if err != nil {
		return "", fmt.Errorf("failed to create zip file: %w", err)
	}
	defer zipFile.Close()

	zipWriter := zip.NewWriter(zipFile)
	defer zipWriter.Close()

	// Add files to zip
	// Note: In a real implementation, we would read actual file contents
	// For now, we'll create placeholder files with metadata
	for _, file := range files {
		// Create a file entry in the zip
		fileWriter, err := zipWriter.Create(file.Path)
		if err != nil {
			return "", fmt.Errorf("failed to create file in zip: %w", err)
		}

		// Write file metadata as content (in real app, this would be actual file content)
		metadata := fmt.Sprintf(`File: %s
Directory: %s
Category: %s
Date: %s
Size: %d bytes
Privileged: %v
Duplicate Hash: %s
`, file.FileName, file.Directory, file.Category, file.Date.Format(time.RFC3339), file.Size, file.Privileged, file.DuplicateHash)

		if _, err := fileWriter.Write([]byte(metadata)); err != nil {
			return "", fmt.Errorf("failed to write file to zip: %w", err)
		}
	}

	// Create manifest file
	manifestWriter, err := zipWriter.Create("manifest.json")
	if err != nil {
		return "", fmt.Errorf("failed to create manifest: %w", err)
	}

	manifest := fmt.Sprintf(`{
  "production_request_id": "%s",
  "created_at": "%s",
  "file_count": %d,
  "total_size": %d
}`, productionRequestID, time.Now().Format(time.RFC3339), len(files), d.calculateTotalSize(files))

	if _, err := manifestWriter.Write([]byte(manifest)); err != nil {
		return "", fmt.Errorf("failed to write manifest: %w", err)
	}

	return zipPath, nil
}

// calculateTotalSize calculates the total size of files
func (d *DB) calculateTotalSize(files []File) int64 {
	var total int64
	for _, file := range files {
		total += file.Size
	}
	return total
}

// CopyFile copies a file from source to destination
func CopyFile(src, dst string) error {
	sourceFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer sourceFile.Close()

	destFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer destFile.Close()

	_, err = io.Copy(destFile, sourceFile)
	return err
}
