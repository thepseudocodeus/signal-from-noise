package database

import (
	"fmt"
	"strings"
	"time"

	"signal-from-noise/assert"
	"signal-from-noise/logging"
)

// QueryBuilder builds SQL queries incrementally
// This makes the "physics" explicit: each filter reduces the problem space
// ASSUMPTION: SQL queries can be built incrementally by adding WHERE clauses
type QueryBuilder struct {
	whereClauses []string
	args         []interface{}
	baseQuery    string
}

// NewQueryBuilder creates a new query builder
// ASSUMPTION: Base query structure is valid
func NewQueryBuilder() *QueryBuilder {
	return &QueryBuilder{
		whereClauses: []string{"1=1"}, // Start with always-true to simplify AND logic
		args:         []interface{}{},
		baseQuery:    "SELECT id, path, directory, category, date, size, privileged, duplicate_hash, file_name, subject, from_email, to_email, sentiment, is_internal, topic FROM files",
	}
}

// AddDateRange adds date range filter
// ASSUMPTION: Date range reduces problem space by temporal filtering
// If dates are invalid, this assumption is falsified
func (qb *QueryBuilder) AddDateRange(start, end *time.Time) *QueryBuilder {
	op := logging.StartOperation("AddDateRange", map[string]interface{}{
		"has_start": start != nil,
		"has_end":   end != nil,
	})
	defer op.EndOperation()

	// ASSUMPTION: If both dates provided, start <= end
	// If this fails, the date range is invalid
	if start != nil && end != nil {
		assert.That(!end.Before(*start), "date range must be valid: start <= end (temporal constraint)")
	}

	if start != nil {
		qb.whereClauses = append(qb.whereClauses, "date >= ?")
		qb.args = append(qb.args, start.Format(time.RFC3339))
		logging.LogCheckpoint("AddDateRange", map[string]interface{}{
			"filter": "date_start",
			"value":  start.Format(time.RFC3339),
		})
	}

	if end != nil {
		qb.whereClauses = append(qb.whereClauses, "date <= ?")
		qb.args = append(qb.args, end.Format(time.RFC3339))
		logging.LogCheckpoint("AddDateRange", map[string]interface{}{
			"filter": "date_end",
			"value":  end.Format(time.RFC3339),
		})
	}

	op.EndOperationWithResult(map[string]interface{}{
		"clauses_added": len(qb.whereClauses) - 1, // Subtract initial "1=1"
	})
	return qb
}

// AddCategories adds category filter
// ASSUMPTION: Category filter reduces problem space by file type
// If no categories selected, this doesn't reduce space (assumption falsified)
func (qb *QueryBuilder) AddCategories(categories []string) *QueryBuilder {
	op := logging.StartOperation("AddCategories", map[string]interface{}{
		"category_count": len(categories),
	})
	defer op.EndOperation()

	if len(categories) == 0 {
		logging.LogCheckpoint("AddCategories", map[string]interface{}{
			"note": "no_categories_selected_no_reduction",
		})
		return qb // No reduction if no categories
	}

	// ASSUMPTION: At least one category must be selected for reduction
	// If empty, we're not filtering (no reduction)
	assert.That(len(categories) > 0, "at least one category must be selected for category filter to reduce problem space")

	placeholders := strings.Repeat("?,", len(categories))
	placeholders = placeholders[:len(placeholders)-1] // Remove trailing comma

	qb.whereClauses = append(qb.whereClauses, fmt.Sprintf("category IN (%s)", placeholders))
	for _, cat := range categories {
		qb.args = append(qb.args, cat)
	}

	logging.LogCheckpoint("AddCategories", map[string]interface{}{
		"categories": categories,
		"reduction_expected": fmt.Sprintf("~%d/%d of space", len(categories), 3), // Assuming 3 categories
	})

	op.EndOperationWithResult(map[string]interface{}{
		"categories": categories,
	})
	return qb
}

// AddTopics adds topic filter
// ASSUMPTION: Topic filter reduces problem space by email subject topics
// Only applies to email files (category='email')
func (qb *QueryBuilder) AddTopics(topics []string) *QueryBuilder {
	op := logging.StartOperation("AddTopics", map[string]interface{}{
		"topic_count": len(topics),
	})
	defer op.EndOperation()

	if len(topics) == 0 {
		return qb // No reduction if no topics
	}

	// ASSUMPTION: Topics only apply to email files
	// If topics selected but no email category, this filter does nothing
	// (This is a constraint we should verify)
	placeholders := strings.Repeat("?,", len(topics))
	placeholders = placeholders[:len(placeholders)-1]

	qb.whereClauses = append(qb.whereClauses, fmt.Sprintf("(category != 'email' OR topic IN (%s))", placeholders))
	for _, topic := range topics {
		qb.args = append(qb.args, topic)
	}

	logging.LogCheckpoint("AddTopics", map[string]interface{}{
		"topics": topics,
		"note":   "applies_to_emails_only",
	})

	op.EndOperationWithResult(map[string]interface{}{
		"topics": topics,
	})
	return qb
}

// AddPeople adds people filter
// ASSUMPTION: People filter reduces problem space by email participants
// ASSUMPTION: People appear in FROM or TO fields (union required)
func (qb *QueryBuilder) AddPeople(people []string, filterType string) *QueryBuilder {
	op := logging.StartOperation("AddPeople", map[string]interface{}{
		"people_count": len(people),
		"filter_type":  filterType,
	})
	defer op.EndOperation()

	// ASSUMPTION: Filter type is valid
	assert.That(filterType == "internal" || filterType == "external" || filterType == "specific" || filterType == "all",
		"people filter type must be: internal, external, specific, or all")

	if filterType == "all" {
		return qb // No reduction
	}

	if filterType == "internal" {
		// ASSUMPTION: Internal emails have is_internal=1
		qb.whereClauses = append(qb.whereClauses, "(category != 'email' OR is_internal = 1)")
		logging.LogCheckpoint("AddPeople", map[string]interface{}{
			"filter": "internal_only",
		})
	} else if filterType == "external" {
		// ASSUMPTION: External emails have is_internal=0
		qb.whereClauses = append(qb.whereClauses, "(category != 'email' OR is_internal = 0)")
		logging.LogCheckpoint("AddPeople", map[string]interface{}{
			"filter": "external_only",
		})
	} else if filterType == "specific" && len(people) > 0 {
		// ASSUMPTION: People appear in FROM or TO fields (union)
		placeholders := strings.Repeat("?,", len(people))
		placeholders = placeholders[:len(placeholders)-1]

		// Union of FROM and TO fields
		qb.whereClauses = append(qb.whereClauses,
			fmt.Sprintf("(category != 'email' OR from_email IN (%s) OR to_email IN (%s))", placeholders, placeholders))
		for _, person := range people {
			qb.args = append(qb.args, person)
			qb.args = append(qb.args, person) // Once for FROM, once for TO
		}

		logging.LogCheckpoint("AddPeople", map[string]interface{}{
			"filter": "specific_people",
			"people": people,
			"note":   "union_of_from_and_to",
		})
	}

	op.EndOperationWithResult(map[string]interface{}{
		"filter_type": filterType,
		"people":      people,
	})
	return qb
}

// AddSentiment adds sentiment filter
// ASSUMPTION: Sentiment filter reduces problem space by email sentiment
func (qb *QueryBuilder) AddSentiment(sentiment string) *QueryBuilder {
	op := logging.StartOperation("AddSentiment", map[string]interface{}{
		"sentiment": sentiment,
	})
	defer op.EndOperation()

	// ASSUMPTION: Sentiment value is valid
	assert.That(sentiment == "positive" || sentiment == "negative" || sentiment == "neutral" || sentiment == "unknown" || sentiment == "all",
		"sentiment must be: positive, negative, neutral, unknown, or all")

	if sentiment == "all" {
		return qb // No reduction
	}

	// ASSUMPTION: Sentiment only applies to email files
	qb.whereClauses = append(qb.whereClauses, fmt.Sprintf("(category != 'email' OR sentiment = ?)"))
	qb.args = append(qb.args, sentiment)

	logging.LogCheckpoint("AddSentiment", map[string]interface{}{
		"sentiment": sentiment,
		"note":      "applies_to_emails_only",
	})

	op.EndOperationWithResult(map[string]interface{}{
		"sentiment": sentiment,
	})
	return qb
}

// AddExcludePrivileged adds privileged exclusion filter
func (qb *QueryBuilder) AddExcludePrivileged(exclude bool) *QueryBuilder {
	if exclude {
		qb.whereClauses = append(qb.whereClauses, "privileged = 0")
	}
	return qb
}

// Build constructs the final SQL query
// ASSUMPTION: Query is valid SQL that will execute
// If this fails, query construction logic is wrong
func (qb *QueryBuilder) Build(page, pageSize int) (string, []interface{}, error) {
	op := logging.StartOperation("QueryBuilder.Build", map[string]interface{}{
		"where_clauses": len(qb.whereClauses),
		"args_count":    len(qb.args),
		"page":          page,
		"page_size":     pageSize,
	})
	defer op.EndOperation()

	// ASSUMPTION: Page and page size are positive
	assert.That(page > 0, "page must be positive for pagination")
	assert.That(pageSize > 0, "page size must be positive for pagination")

	whereClause := strings.Join(qb.whereClauses, " AND ")
	offset := (page - 1) * pageSize

	query := fmt.Sprintf(`
		%s
		WHERE %s
		ORDER BY date DESC
		LIMIT ? OFFSET ?
	`, qb.baseQuery, whereClause)

	args := append(qb.args, pageSize, offset)

	logging.LogQuery(query, map[string]interface{}{
		"where_clauses": len(qb.whereClauses),
		"args":          len(args),
		"page":          page,
		"page_size":     pageSize,
		"offset":        offset,
	})

	op.EndOperationWithResult(map[string]interface{}{
		"query":      query,
		"args_count": len(args),
	})

	return query, args, nil
}

// GetCountQuery builds a COUNT query with same filters
// ASSUMPTION: Count query uses same WHERE clause as main query
func (qb *QueryBuilder) GetCountQuery() (string, []interface{}) {
	whereClause := strings.Join(qb.whereClauses, " AND ")
	query := fmt.Sprintf("SELECT COUNT(*) FROM files WHERE %s", whereClause)
	return query, qb.args
}
