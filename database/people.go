package database

import (
	"database/sql"
	"fmt"

	"signal-from-noise/assert"
	"signal-from-noise/logging"
)

// PeopleList represents internal and external email addresses
// ASSUMPTION: People can be classified as internal or external
// ASSUMPTION: Internal/external sets are precomputed (is_internal field)
type PeopleList struct {
	Internal []string `json:"internal"` // Precomputed internal email addresses
	External []string `json:"external"` // Precomputed external email addresses
	All      []string `json:"all"`      // Union of all email addresses
}

// GetPeople returns lists of internal and external email addresses
// ASSUMPTION: People appear in FROM or TO fields of emails
// ASSUMPTION: Internal/external classification is precomputed (is_internal field)
// ASSUMPTION: Union of FROM and TO gives complete people list
func (d *DB) GetPeople() (*PeopleList, error) {
	op := logging.StartOperation("GetPeople", map[string]interface{}{})
	defer op.EndOperation()

	// ASSUMPTION: Database connection exists
	assert.ThatNotNil(d.db, "database connection must exist to get people")

	// Get internal emails (precomputed set)
	// ASSUMPTION: is_internal=1 means internal email
	internalQuery := `
		SELECT DISTINCT from_email
		FROM files
		WHERE category = 'email' AND is_internal = 1 AND from_email IS NOT NULL AND from_email != ''
		UNION
		SELECT DISTINCT to_email
		FROM files
		WHERE category = 'email' AND is_internal = 1 AND to_email IS NOT NULL AND to_email != ''
		ORDER BY from_email
	`

	logging.LogQuery(internalQuery, map[string]interface{}{
		"note": "precomputed_internal_emails_union_from_and_to",
	})

	internalRows, err := d.db.Query(internalQuery)
	if err != nil {
		logging.LogError("GetPeople", err, map[string]interface{}{
			"operation": "query_internal",
		})
		return nil, fmt.Errorf("failed to query internal emails: %w", err)
	}
	defer internalRows.Close()

	internalSet := make(map[string]bool)
	var internal []string
	for internalRows.Next() {
		var email string
		if err := internalRows.Scan(&email); err != nil {
			return nil, fmt.Errorf("failed to scan internal email: %w", err)
		}
		if !internalSet[email] {
			internalSet[email] = true
			internal = append(internal, email)
		}
	}

	// Get external emails (precomputed set)
	// ASSUMPTION: is_internal=0 means external email
	externalQuery := `
		SELECT DISTINCT from_email
		FROM files
		WHERE category = 'email' AND is_internal = 0 AND from_email IS NOT NULL AND from_email != ''
		UNION
		SELECT DISTINCT to_email
		FROM files
		WHERE category = 'email' AND is_internal = 0 AND to_email IS NOT NULL AND to_email != ''
		ORDER BY from_email
	`

	logging.LogQuery(externalQuery, map[string]interface{}{
		"note": "precomputed_external_emails_union_from_and_to",
	})

	externalRows, err := d.db.Query(externalQuery)
	if err != nil {
		logging.LogError("GetPeople", err, map[string]interface{}{
			"operation": "query_external",
		})
		return nil, fmt.Errorf("failed to query external emails: %w", err)
	}
	defer externalRows.Close()

	externalSet := make(map[string]bool)
	var external []string
	for externalRows.Next() {
		var email string
		if err := externalRows.Scan(&email); err != nil {
			return nil, fmt.Errorf("failed to scan external email: %w", err)
		}
		if !externalSet[email] {
			externalSet[email] = true
			external = append(external, email)
		}
	}

	// Create union of all emails
	// ASSUMPTION: Union of internal and external gives complete people list
	allSet := make(map[string]bool)
	var all []string
	for _, email := range internal {
		if !allSet[email] {
			allSet[email] = true
			all = append(all, email)
		}
	}
	for _, email := range external {
		if !allSet[email] {
			allSet[email] = true
			all = append(all, email)
		}
	}

	// ASSUMPTION: At least some people exist
	// If no people found, either no emails in database or email fields are NULL
	if len(all) == 0 {
		logging.LogCheckpoint("GetPeople", map[string]interface{}{
			"warning": "no_people_found",
			"note":    "may_indicate_no_emails_or_null_email_fields",
		})
	}

	op.EndOperationWithResult(map[string]interface{}{
		"internal_count": len(internal),
		"external_count": len(external),
		"total_count":    len(all),
	})

	return &PeopleList{
		Internal: internal,
		External: external,
		All:      all,
	}, nil
}
