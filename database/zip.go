package database

import (
	"archive/zip"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"

	"signal-from-noise/assert"
	"signal-from-noise/logging"
)

// CreateZipFile creates a zip file containing the specified files
// Returns the path to the created zip file
// Assumption: At least one file ID is provided
// Assumption: All file IDs exist in database
// Assumption: Output directory is writable
func (d *DB) CreateZipFile(productionRequestID string, fileIDs []int64, outputDir string) (string, error) {
	op := logging.StartOperation("CreateZipFile", map[string]interface{}{
		"production_request_id": productionRequestID,
		"file_count":           len(fileIDs),
		"output_dir":           outputDir,
	})
	defer op.EndOperation()

	// ASSUMPTION: At least one file must be selected for zip creation
	// Empty file list would create an invalid zip file
	assert.That(len(fileIDs) > 0, "at least one file ID must be provided for zip creation")

	// ASSUMPTION: Production request ID must be non-empty for file naming
	// Empty ID would create incorrectly named zip files
	assert.That(productionRequestID != "", "production request ID must be non-empty for zip file naming")

	// Get files from database
	// ASSUMPTION: All file IDs exist in database and can be retrieved
	// If files don't exist, zip creation cannot proceed
	files, err := d.GetFilesByIDs(fileIDs)
	if err != nil {
		logging.LogError("CreateZipFile", err, map[string]interface{}{
			"operation": "get_files_by_ids",
			"file_count": len(fileIDs),
		})
		return "", fmt.Errorf("failed to get files: %w", err)
	}

	// ASSUMPTION: All requested files must exist in database
	// If some files are missing, the user selected invalid file IDs
	assert.That(len(files) > 0, "at least one file must be found for selected IDs (database integrity check)")
	assert.That(len(files) == len(fileIDs), "all requested file IDs must exist in database (data consistency check)")

	logging.LogResult("GetFilesByIDs", len(files), map[string]interface{}{
		"requested_count": len(fileIDs),
		"found_count":     len(files),
	})

	// Ensure output directory exists
	// ASSUMPTION: Output directory path is valid and can be created
	// If this fails, the file system is not accessible or permissions are wrong
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		logging.LogError("CreateZipFile", err, map[string]interface{}{
			"operation":  "create_output_directory",
			"output_dir": outputDir,
		})
		return "", fmt.Errorf("failed to create output directory: %w", err)
	}

	// Create zip file name: {productionRequestID}_{timestamp}.zip
	// ASSUMPTION: Timestamp format is valid and creates unique filenames
	// Format ensures no collisions and sortable by creation time
	timestamp := time.Now().Format("20060102_150405")
	zipFileName := fmt.Sprintf("%s_%s.zip", productionRequestID, timestamp)
	zipPath := filepath.Join(outputDir, zipFileName)

	logging.LogCheckpoint("CreateZipFile", map[string]interface{}{
		"zip_file_name": zipFileName,
		"zip_path":      zipPath,
		"files_to_add":  len(files),
	})

	// Create zip file
	// ASSUMPTION: Output directory is writable and file can be created
	// If this fails, disk is full or permissions are incorrect
	zipFile, err := os.Create(zipPath)
	if err != nil {
		logging.LogError("CreateZipFile", err, map[string]interface{}{
			"operation": "create_zip_file",
			"zip_path":   zipPath,
		})
		return "", fmt.Errorf("failed to create zip file: %w", err)
	}
	defer zipFile.Close()

	zipWriter := zip.NewWriter(zipFile)
	defer zipWriter.Close()

	// Add files to zip
	// Note: In a real implementation, we would read actual file contents
	// For now, we'll create placeholder files with metadata
	// ASSUMPTION: All files can be added to zip without errors
	// If any file fails, the zip is corrupted
	filesAdded := 0
	for _, file := range files {
		// ASSUMPTION: File path is valid for zip entry creation
		// Invalid paths would cause zip creation to fail
		assert.That(file.Path != "", "file path must be non-empty for zip entry creation")

		// Create a file entry in the zip
		fileWriter, err := zipWriter.Create(file.Path)
		if err != nil {
			logging.LogError("CreateZipFile", err, map[string]interface{}{
				"operation": "create_zip_entry",
				"file_path": file.Path,
			})
			return "", fmt.Errorf("failed to create file in zip: %w", err)
		}

		// Write file metadata as content (in real app, this would be actual file content)
		// ASSUMPTION: File date can be formatted as RFC3339
		// If formatting fails, the date structure is invalid
		metadata := fmt.Sprintf(`File: %s
Directory: %s
Category: %s
Date: %s
Size: %d bytes
Privileged: %v
Duplicate Hash: %s
`, file.FileName, file.Directory, file.Category, file.Date.Format(time.RFC3339), file.Size, file.Privileged, file.DuplicateHash)

		if _, err := fileWriter.Write([]byte(metadata)); err != nil {
			logging.LogError("CreateZipFile", err, map[string]interface{}{
				"operation": "write_file_to_zip",
				"file_path": file.Path,
			})
			return "", fmt.Errorf("failed to write file to zip: %w", err)
		}
		filesAdded++
	}

	// ASSUMPTION: All files were successfully added to zip
	// If count doesn't match, some files failed silently
	assert.That(filesAdded == len(files), "all files must be added to zip (zip integrity check)")

	logging.LogCheckpoint("CreateZipFile", map[string]interface{}{
		"files_added": filesAdded,
		"total_files": len(files),
	})

	// Create manifest file
	// ASSUMPTION: Manifest entry can be created in zip
	// If this fails, zip structure is corrupted
	manifestWriter, err := zipWriter.Create("manifest.json")
	if err != nil {
		logging.LogError("CreateZipFile", err, map[string]interface{}{
			"operation": "create_manifest_entry",
		})
		return "", fmt.Errorf("failed to create manifest: %w", err)
	}

	totalSize := d.calculateTotalSize(files)
	manifest := fmt.Sprintf(`{
  "production_request_id": "%s",
  "created_at": "%s",
  "file_count": %d,
  "total_size": %d
}`, productionRequestID, time.Now().Format(time.RFC3339), len(files), totalSize)

	if _, err := manifestWriter.Write([]byte(manifest)); err != nil {
		logging.LogError("CreateZipFile", err, map[string]interface{}{
			"operation": "write_manifest",
		})
		return "", fmt.Errorf("failed to write manifest: %w", err)
	}

	// ASSUMPTION: Zip file was successfully created and is accessible
	// If file doesn't exist after creation, the operation silently failed
	// Note: We can't check file existence here because zipWriter.Close() hasn't been called yet
	// The defer will close it, so we verify existence would need to happen after this function

	op.EndOperationWithResult(map[string]interface{}{
		"zip_path":    zipPath,
		"files_added": filesAdded,
		"total_size":  totalSize,
		"success":     true,
	})

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
