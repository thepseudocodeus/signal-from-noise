package database

import (
	"fmt"
	"strings"

	"signal-from-noise/assert"
	"signal-from-noise/logging"
)

// GetTopics returns unique topics extracted from email subjects
// ASSUMPTION: Topics exist in database and can be extracted
// This precomputes the set of available topics for user selection
func (d *DB) GetTopics() ([]string, error) {
	op := logging.StartOperation("GetTopics", map[string]interface{}{})
	defer op.EndOperation()

	// ASSUMPTION: Database connection exists
	assert.ThatNotNil(d.db, "database connection must exist to get topics")

	// ASSUMPTION: Topic column exists and contains topic values
	// If column doesn't exist, schema migration is needed
	query := `
		SELECT DISTINCT topic
		FROM files
		WHERE topic IS NOT NULL AND topic != ''
		ORDER BY topic
	`

	logging.LogQuery(query, map[string]interface{}{
		"operation": "get_topics",
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

	// ASSUMPTION: Topics list is valid (may be empty if no topics in data)
	// Empty list is valid - means no topics available yet
	logging.LogResult("GetTopics", len(topics), map[string]interface{}{
		"topics": topics,
	})

	op.EndOperationWithResult(map[string]interface{}{
		"topic_count": len(topics),
	})

	return topics, nil
}

// GetPeople returns lists of internal and external email addresses
// ASSUMPTION: Email addresses exist in FROM and TO fields
// Precomputes internal/external sets for efficient filtering
func (d *DB) GetPeople() (internal []string, external []string, err error) {
	op := logging.StartOperation("GetPeople", map[string]interface{}{})
	defer op.EndOperation()

	// ASSUMPTION: Database connection exists
	assert.ThatNotNil(d.db, "database connection must exist to get people")

	// Get internal emails (from FROM and TO fields where is_internal = 1)
	// ASSUMPTION: is_internal flag correctly classifies emails
	internalQuery := `
		SELECT DISTINCT email
		FROM (
			SELECT from_email AS email FROM files WHERE is_internal = 1 AND from_email IS NOT NULL AND from_email != ''
			UNION
			SELECT to_email AS email FROM files WHERE is_internal = 1 AND to_email IS NOT NULL AND to_email != ''
		)
		ORDER BY email
	`

	logging.LogQuery(internalQuery, map[string]interface{}{
		"operation": "get_internal_people",
	})

	rows, err := d.db.Query(internalQuery)
	if err != nil {
		logging.LogError("GetPeople", err, map[string]interface{}{
			"operation": "query_internal_people",
		})
		return nil, nil, fmt.Errorf("failed to query internal people: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var email string
		if err := rows.Scan(&email); err != nil {
			logging.LogError("GetPeople", err, map[string]interface{}{
				"operation": "scan_internal_email",
			})
			return nil, nil, fmt.Errorf("failed to scan internal email: %w", err)
		}
		internal = append(internal, email)
	}
	rows.Close()

	// Get external emails (from FROM and TO fields where is_internal = 0)
	externalQuery := `
		SELECT DISTINCT email
		FROM (
			SELECT from_email AS email FROM files WHERE is_internal = 0 AND from_email IS NOT NULL AND from_email != ''
			UNION
			SELECT to_email AS email FROM files WHERE is_internal = 0 AND to_email IS NOT NULL AND to_email != ''
		)
		ORDER BY email
	`

	logging.LogQuery(externalQuery, map[string]interface{}{
		"operation": "get_external_people",
	})

	rows, err = d.db.Query(externalQuery)
	if err != nil {
		logging.LogError("GetPeople", err, map[string]interface{}{
			"operation": "query_external_people",
		})
		return nil, nil, fmt.Errorf("failed to query external people: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var email string
		if err := rows.Scan(&email); err != nil {
			logging.LogError("GetPeople", err, map[string]interface{}{
				"operation": "scan_external_email",
			})
			return nil, nil, fmt.Errorf("failed to scan external email: %w", err)
		}
		external = append(external, email)
	}

	// ASSUMPTION: People lists are valid (may be empty)
	// Empty lists are valid - means no people in that category
	logging.LogResult("GetPeople", len(internal)+len(external), map[string]interface{}{
		"internal_count": len(internal),
		"external_count": len(external),
	})

	op.EndOperationWithResult(map[string]interface{}{
		"internal_count": len(internal),
		"external_count": len(external),
	})

	return internal, external, nil
}

// GetSentimentOptions returns available sentiment values
// ASSUMPTION: Sentiment values are standardized in database
func (d *DB) GetSentimentOptions() ([]string, error) {
	// Sentiment options are fixed - no need to query database
	// ASSUMPTION: These are the only valid sentiment values
	// If database has different values, they should be normalized
	return []string{"all", "positive", "negative", "neutral", "unknown"}, nil
}
