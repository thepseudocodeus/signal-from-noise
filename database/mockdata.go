package database

import (
	"fmt"
	"math/rand"
	"time"
)

// SeedMockData populates the database with mock data
func (d *DB) SeedMockData() error {
	// Seed random number generator
	rand.Seed(time.Now().UnixNano())

	// Generate dates across a 3-year range
	baseDate := time.Date(2022, 1, 1, 0, 0, 0, 0, time.UTC)
	endDate := time.Date(2024, 12, 31, 23, 59, 59, 0, time.UTC)

	// Create directories
	directories := []struct {
		name     string
		category string
	}{
		{"Email_Folder_2022", "email"},
		{"Email_Folder_2023", "email"},
		{"Email_Folder_2024", "email"},
		{"Claim_Documents_2022", "claim"},
		{"Claim_Documents_2023", "claim"},
		{"Claim_Documents_2024", "claim"},
		{"Claim_Evidence_2023", "claim"},
		{"Other_Documents", "other"},
		{"Misc_Files", "other"},
		{"Archive_2022", "other"},
	}

	// Generate files
	fileCount := 0
	duplicateHashes := make(map[string]int) // Track duplicates

	for _, dir := range directories {
		// Generate 50-150 files per directory
		filesInDir := 50 + rand.Intn(101)

		for i := 0; i < filesInDir; i++ {
			// Random date within range
			daysOffset := rand.Intn(int(endDate.Sub(baseDate).Hours() / 24))
			fileDate := baseDate.AddDate(0, 0, daysOffset)

			// Generate file name
			fileName := fmt.Sprintf("file_%d_%d.pdf", fileCount, rand.Intn(10000))

			// Determine if privileged (only for emails, 20% chance)
			privileged := false
			if dir.category == "email" && rand.Float32() < 0.2 {
				privileged = true
			}

			// Create some duplicates (10% chance of being a duplicate)
			duplicateHash := ""
			if rand.Float32() < 0.1 && len(duplicateHashes) > 0 {
				// Pick a random existing hash
				idx := rand.Intn(len(duplicateHashes))
				counter := 0
				for hash := range duplicateHashes {
					if counter == idx {
						duplicateHash = hash
						duplicateHashes[hash]++
						break
					}
					counter++
				}
			} else {
				// Generate new unique hash
				duplicateHash = fmt.Sprintf("hash_%d_%d", fileCount, rand.Intn(1000000))
				duplicateHashes[duplicateHash] = 1
			}

			// Random file size (1KB to 10MB)
			fileSize := int64(1024 + rand.Intn(10*1024*1024-1024))

			// Insert file
			query := `
				INSERT INTO files (path, directory, category, date, size, privileged, duplicate_hash, file_name)
				VALUES (?, ?, ?, ?, ?, ?, ?, ?)
			`

			path := fmt.Sprintf("%s/%s", dir.name, fileName)
			_, err := d.db.Exec(query,
				path,
				dir.name,
				dir.category,
				fileDate.Format(time.RFC3339),
				fileSize,
				privileged,
				duplicateHash,
				fileName,
			)
			if err != nil {
				return fmt.Errorf("failed to insert file: %w", err)
			}

			fileCount++
		}
	}

	// Insert production requests
	productionRequests := []struct {
		id          string
		title       string
		description string
	}{
		{"PR-001", "Production Request #1", "Email communication timeline"},
		{"PR-002", "Production Request #2", "Claims-related communications"},
		{"PR-003", "Production Request #3", "Document review"},
		{"PR-004", "Production Request #4", "Evidence compilation"},
		{"PR-005", "Production Request #5", "Correspondence analysis"},
	}

	for _, pr := range productionRequests {
		query := `
			INSERT OR REPLACE INTO production_requests (id, title, description)
			VALUES (?, ?, ?)
		`
		_, err := d.db.Exec(query, pr.id, pr.title, pr.description)
		if err != nil {
			return fmt.Errorf("failed to insert production request: %w", err)
		}
	}

	fmt.Printf("Seeded %d files and %d production requests\n", fileCount, len(productionRequests))
	return nil
}
