package database

// GetSentimentOptions returns available sentiment values
// ASSUMPTION: Sentiment values are standardized in database
func (d *DB) GetSentimentOptions() ([]string, error) {
	// Sentiment options are fixed - no need to query database
	// ASSUMPTION: These are the only valid sentiment values
	// If database has different values, they should be normalized
	return []string{"all", "positive", "negative", "neutral", "unknown"}, nil
}
