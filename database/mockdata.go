package database

import (
	"database/sql"
	"fmt"
	"math/rand"
	"time"

	"signal-from-noise/logging"
)

// Helper functions for nullable fields
func getNullString(ns sql.NullString) interface{} {
	if ns.Valid {
		return ns.String
	}
	return nil
}

func getNullInt64(ni sql.NullInt64) interface{} {
	if ni.Valid {
		return ni.Int64
	}
	return nil
}

// SeedMockData populates the database with mock data
// ASSUMPTION: Database schema is initialized before seeding
// ASSUMPTION: Mock data represents real-world patterns (topics, people, sentiment)
func (d *DB) SeedMockData() error {
	op := logging.StartOperation("SeedMockData", map[string]interface{}{})
	defer op.EndOperation()

	// Seed random number generator
	rand.Seed(time.Now().UnixNano())

	// Generate dates across a 3-year range
	baseDate := time.Date(2022, 1, 1, 0, 0, 0, 0, time.UTC)
	endDate := time.Date(2024, 12, 31, 23, 59, 59, 0, time.UTC)

	// Precomputed sets for realistic data
	// ASSUMPTION: These represent realistic email patterns
	topics := []string{
		"Project Update", "Meeting Request", "Contract Review", "Budget Approval",
		"Status Report", "Client Communication", "Legal Matter", "Invoice",
		"Proposal", "Follow-up", "Urgent Action", "Documentation",
	}

	internalEmails := []string{
		"john.doe@company.com", "jane.smith@company.com", "bob.jones@company.com",
		"alice.brown@company.com", "charlie.wilson@company.com", "diana.miller@company.com",
	}

	externalEmails := []string{
		"client1@external.com", "vendor@supplier.com", "partner@business.com",
		"customer@client.com", "consultant@firm.com", "lawyer@legal.com",
	}

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

			// Generate email-specific fields (only for email category)
			// ASSUMPTION: Email files have subject, from, to, sentiment, topic
			// Non-email files have NULL for these fields
			var subject, fromEmail, toEmail, sentiment, topic string
			var isInternal int

			if dir.category == "email" {
				// Generate subject with topic
				topic = topics[rand.Intn(len(topics))]
				subject = fmt.Sprintf("%s - %s", topic, fmt.Sprintf("Email %d", fileCount))

				// Determine if internal or external (70% internal, 30% external)
				isInternal = 0
				if rand.Float32() < 0.7 {
					isInternal = 1
					fromEmail = internalEmails[rand.Intn(len(internalEmails))]
					toEmail = internalEmails[rand.Intn(len(internalEmails))]
					// Ensure from and to are different
					for toEmail == fromEmail {
						toEmail = internalEmails[rand.Intn(len(internalEmails))]
					}
				} else {
					// External email - mix of internal and external addresses
					if rand.Float32() < 0.5 {
						fromEmail = internalEmails[rand.Intn(len(internalEmails))]
						toEmail = externalEmails[rand.Intn(len(externalEmails))]
					} else {
						fromEmail = externalEmails[rand.Intn(len(externalEmails))]
						toEmail = internalEmails[rand.Intn(len(internalEmails))]
					}
				}

				// Assign sentiment (weighted: 40% neutral, 30% positive, 20% negative, 10% unknown)
				sentimentRand := rand.Float32()
				switch {
				case sentimentRand < 0.4:
					sentiment = "neutral"
				case sentimentRand < 0.7:
					sentiment = "positive"
				case sentimentRand < 0.9:
					sentiment = "negative"
				default:
					sentiment = "unknown"
				}
			}

			// Insert file
			// ASSUMPTION: All fields match schema (including nullable email fields)
			query := `
				INSERT INTO files (path, directory, category, date, size, privileged, duplicate_hash, file_name,
				                   subject, from_email, to_email, sentiment, is_internal, topic)
				VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
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
				subject,  // NULL for non-email files
				fromEmail, // NULL for non-email files
				toEmail,   // NULL for non-email files
				sentiment, // NULL for non-email files
				isInternal,
				topic,     // NULL for non-email files
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

	logging.LogResult("SeedMockData", fileCount, map[string]interface{}{
		"production_requests": len(productionRequests),
		"topics_available":     len(topics),
		"internal_emails":      len(internalEmails),
		"external_emails":     len(externalEmails),
	})

	op.EndOperationWithResult(map[string]interface{}{
		"files_seeded":         fileCount,
		"production_requests": len(productionRequests),
	})

	logging.LogResult("SeedMockData", fileCount, map[string]interface{}{
		"production_requests": len(productionRequests),
		"topics_available":     len(topics),
		"internal_emails":     len(internalEmails),
		"external_emails":     len(externalEmails),
	})

	op.EndOperationWithResult(map[string]interface{}{
		"files_seeded":        fileCount,
		"production_requests": len(productionRequests),
	})

	fmt.Printf("Seeded %d files and %d production requests\n", fileCount, len(productionRequests))
	return nil
}
