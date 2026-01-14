package database

import (
	"fmt"

	"signal-from-noise/assert"
	"signal-from-noise/logging"
)

// GetTopics returns unique topics from email subjects
// ASSUMPTION: Topics can be extracted from email subjects
// ASSUMPTION: Topics exist in database (topic field is populated)
// If no topics found, assumption about topic extraction is falsified
func (d *DB) GetTopics() ([]string, error) {
	op := logging.StartOperation("GetTopics", map[string]interface{}{})
	defer op.EndOperation()

	// ASSUMPTION: Database connection exists
	assert.ThatNotNil(d.db, "database connection must exist to get topics")

	// ASSUMPTION: Topics are stored in topic field
	// If topic field is NULL for all emails, assumption about topic extraction is false
	query := `
		SELECT DISTINCT topic
		FROM files
		WHERE category = 'email' AND topic IS NOT NULL AND topic != ''
		ORDER BY topic
	`

	logging.LogQuery(query, map[string]interface{}{
		"note": "extracting_unique_topics_from_email_subjects",
	})

	rows, err := d.db.Query(query)
	if err != nil {
		logging.LogError("GetTopics", err, map[string]interface{}{
			"operation": "query_topics",
		})
		return nil, fmt.Errorf("failed to query topics: %w", err)
	}
	defer rows.Close()

	var topics []string
	for rows.Next() {
		var topic string
		if err := rows.Scan(&topic); err != nil {
			logging.LogError("GetTopics", err, map[string]interface{}{
				"operation": "scan_topic",
			})
			return nil, fmt.Errorf("failed to scan topic: %w", err)
		}
		topics = append(topics, topic)
	}

	// ASSUMPTION: At least some topics exist
	// If no topics found, either no emails in database or topic extraction failed
	if len(topics) == 0 {
		logging.LogCheckpoint("GetTopics", map[string]interface{}{
			"warning": "no_topics_found",
			"note":    "may_indicate_topic_extraction_failed_or_no_emails",
		})
	}

	op.EndOperationWithResult(map[string]interface{}{
		"topic_count": len(topics),
		"topics":      topics,
	})

	return topics, nil
}
