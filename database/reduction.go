package database

import (
	"database/sql"
	"fmt"

	"signal-from-noise/assert"
	"signal-from-noise/logging"
)

// ReductionMetrics tracks problem space reduction at each step
// This makes the "physics" explicit: how each filter reduces complexity
type ReductionMetrics struct {
	InitialSpace    int                    `json:"initial_space"`    // |U| = total files
	StepReductions  []StepReduction        `json:"step_reductions"`   // Reduction at each step
	FinalSpace      int                    `json:"final_space"`      // |F| = final filtered files
	TotalReduction  float64                `json:"total_reduction"`  // (|U| - |F|) / |U|
	Confidence      float64                `json:"confidence"`       // Confidence that result satisfies request
}

// StepReduction tracks reduction at a single filtering step
type StepReduction struct {
	StepName      string  `json:"step_name"`       // e.g., "category", "date_range"
	SpaceBefore   int     `json:"space_before"`    // |Before|
	SpaceAfter    int     `json:"space_after"`     // |After|
	Reduction     float64 `json:"reduction"`       // (|Before| - |After|) / |Before|
	ExpectedReduction float64 `json:"expected_reduction"` // Expected reduction factor (0.0 to 1.0)
	AssumptionMet bool    `json:"assumption_met"`  // Did actual match expected?
}

// CalculateReductionMetrics calculates problem space reduction through filtering steps
// ASSUMPTION: Each filter step reduces the problem space
// If reduction is less than expected, assumption is falsified
func (d *DB) CalculateReductionMetrics(filters FileFilters) (*ReductionMetrics, error) {
	op := logging.StartOperation("CalculateReductionMetrics", map[string]interface{}{
		"filters": fmt.Sprintf("%+v", filters),
	})
	defer op.EndOperation()

	// ASSUMPTION: Database connection exists
	assert.ThatNotNil(d.db, "database connection must exist for reduction calculations")

	metrics := &ReductionMetrics{
		StepReductions: []StepReduction{},
	}

	// Step 1: Get initial space |U|
	// ASSUMPTION: Initial space is all files in database
	var initialCount int
	err := d.db.QueryRow("SELECT COUNT(*) FROM files").Scan(&initialCount)
	if err != nil {
		logging.LogError("CalculateReductionMetrics", err, map[string]interface{}{
			"operation": "get_initial_count",
		})
		return nil, fmt.Errorf("failed to get initial count: %w", err)
	}
	metrics.InitialSpace = initialCount

	logging.LogCheckpoint("CalculateReductionMetrics", map[string]interface{}{
		"initial_space": initialCount,
		"formula":       "|U| = COUNT(*) FROM files",
	})

	// Step 2: Calculate reduction at each filtering step
	// This makes the incremental reduction explicit and verifiable
	currentCount := initialCount
	builder := NewQueryBuilder()

	// Category filter
	if len(filters.Categories) > 0 {
		builder.AddCategories(filters.Categories)
		countQuery, countArgs := builder.GetCountQuery()
		var afterCount int
		err := d.db.QueryRow(countQuery, countArgs...).Scan(&afterCount)
		if err != nil {
			return nil, fmt.Errorf("failed to count after category filter: %w", err)
		}

		reduction := calculateReduction(currentCount, afterCount)
		expectedReduction := 1.0 - (float64(len(filters.Categories)) / 3.0) // Assuming 3 categories
		
		stepReduction := StepReduction{
			StepName:         "category",
			SpaceBefore:      currentCount,
			SpaceAfter:       afterCount,
			Reduction:         reduction,
			ExpectedReduction: expectedReduction,
			AssumptionMet:    reduction >= expectedReduction*0.8, // Allow 20% tolerance
		}

		// ASSUMPTION: Category filter reduces space
		// If reduction is negative or zero, assumption is falsified
		assert.That(reduction > 0, fmt.Sprintf("category filter must reduce space (got reduction=%.2f)", reduction))

		metrics.StepReductions = append(metrics.StepReductions, stepReduction)
		currentCount = afterCount

		logging.LogCheckpoint("CalculateReductionMetrics", map[string]interface{}{
			"step":              "category",
			"before":            stepReduction.SpaceBefore,
			"after":             stepReduction.SpaceAfter,
			"reduction":         reduction,
			"expected":          expectedReduction,
			"assumption_met":   stepReduction.AssumptionMet,
		})
	}

	// Date range filter
	if filters.DateStart != nil || filters.DateEnd != nil {
		builder.AddDateRange(filters.DateStart, filters.DateEnd)
		countQuery, countArgs := builder.GetCountQuery()
		var afterCount int
		err := d.db.QueryRow(countQuery, countArgs...).Scan(&afterCount)
		if err != nil {
			return nil, fmt.Errorf("failed to count after date filter: %w", err)
		}

		reduction := calculateReduction(currentCount, afterCount)
		// Expected reduction depends on date range (simplified: assume 50% for now)
		expectedReduction := 0.5

		stepReduction := StepReduction{
			StepName:         "date_range",
			SpaceBefore:      currentCount,
			SpaceAfter:       afterCount,
			Reduction:         reduction,
			ExpectedReduction: expectedReduction,
			AssumptionMet:    reduction >= expectedReduction*0.7, // Allow 30% tolerance
		}

		assert.That(reduction > 0, fmt.Sprintf("date range filter must reduce space (got reduction=%.2f)", reduction))

		metrics.StepReductions = append(metrics.StepReductions, stepReduction)
		currentCount = afterCount

		logging.LogCheckpoint("CalculateReductionMetrics", map[string]interface{}{
			"step":            "date_range",
			"before":          stepReduction.SpaceBefore,
			"after":           stepReduction.SpaceAfter,
			"reduction":       reduction,
			"assumption_met": stepReduction.AssumptionMet,
		})
	}

	// Topics filter
	if len(filters.Topics) > 0 {
		builder.AddTopics(filters.Topics)
		countQuery, countArgs := builder.GetCountQuery()
		var afterCount int
		err := d.db.QueryRow(countQuery, countArgs...).Scan(&afterCount)
		if err != nil {
			return nil, fmt.Errorf("failed to count after topics filter: %w", err)
		}

		reduction := calculateReduction(currentCount, afterCount)
		expectedReduction := 0.2 // Assume topics reduce by ~80% (keep 20%)

		stepReduction := StepReduction{
			StepName:         "topics",
			SpaceBefore:      currentCount,
			SpaceAfter:       afterCount,
			Reduction:         reduction,
			ExpectedReduction: expectedReduction,
			AssumptionMet:    reduction >= expectedReduction*0.5, // Allow 50% tolerance (topics vary widely)
		}

		if reduction > 0 {
			metrics.StepReductions = append(metrics.StepReductions, stepReduction)
			currentCount = afterCount
		}

		logging.LogCheckpoint("CalculateReductionMetrics", map[string]interface{}{
			"step":            "topics",
			"before":          stepReduction.SpaceBefore,
			"after":           stepReduction.SpaceAfter,
			"reduction":       reduction,
			"assumption_met": stepReduction.AssumptionMet,
		})
	}

	// People filter
	if filters.PeopleFilterType != "" && filters.PeopleFilterType != "all" {
		builder.AddPeople(filters.People, filters.PeopleFilterType)
		countQuery, countArgs := builder.GetCountQuery()
		var afterCount int
		err := d.db.QueryRow(countQuery, countArgs...).Scan(&afterCount)
		if err != nil {
			return nil, fmt.Errorf("failed to count after people filter: %w", err)
		}

		reduction := calculateReduction(currentCount, afterCount)
		expectedReduction := 0.3 // Assume people filter reduces by ~70%

		stepReduction := StepReduction{
			StepName:         "people",
			SpaceBefore:      currentCount,
			SpaceAfter:       afterCount,
			Reduction:         reduction,
			ExpectedReduction: expectedReduction,
			AssumptionMet:    reduction >= expectedReduction*0.6,
		}

		if reduction > 0 {
			metrics.StepReductions = append(metrics.StepReductions, stepReduction)
			currentCount = afterCount
		}

		logging.LogCheckpoint("CalculateReductionMetrics", map[string]interface{}{
			"step":            "people",
			"before":          stepReduction.SpaceBefore,
			"after":           stepReduction.SpaceAfter,
			"reduction":       reduction,
			"assumption_met": stepReduction.AssumptionMet,
		})
	}

	// Sentiment filter
	if filters.Sentiment != "" && filters.Sentiment != "all" {
		builder.AddSentiment(filters.Sentiment)
		countQuery, countArgs := builder.GetCountQuery()
		var afterCount int
		err := d.db.QueryRow(countQuery, countArgs...).Scan(&afterCount)
		if err != nil {
			return nil, fmt.Errorf("failed to count after sentiment filter: %w", err)
		}

		reduction := calculateReduction(currentCount, afterCount)
		expectedReduction := 0.25 // Assume sentiment reduces by ~75% (keep 25%)

		stepReduction := StepReduction{
			StepName:         "sentiment",
			SpaceBefore:      currentCount,
			SpaceAfter:       afterCount,
			Reduction:         reduction,
			ExpectedReduction: expectedReduction,
			AssumptionMet:    reduction >= expectedReduction*0.5,
		}

		if reduction > 0 {
			metrics.StepReductions = append(metrics.StepReductions, stepReduction)
			currentCount = afterCount
		}

		logging.LogCheckpoint("CalculateReductionMetrics", map[string]interface{}{
			"step":            "sentiment",
			"before":          stepReduction.SpaceBefore,
			"after":           stepReduction.SpaceAfter,
			"reduction":       reduction,
			"assumption_met": stepReduction.AssumptionMet,
		})
	}

	// Calculate final metrics
	metrics.FinalSpace = currentCount
	if metrics.InitialSpace > 0 {
		metrics.TotalReduction = float64(metrics.InitialSpace-metrics.FinalSpace) / float64(metrics.InitialSpace)
	}

	// Calculate confidence
	// ASSUMPTION: Confidence = 1 - (false_positive_rate)
	// Simplified: confidence increases with reduction
	// More sophisticated: would measure actual relevance
	metrics.Confidence = calculateConfidence(metrics)

	// ASSUMPTION: Total reduction should be significant (>50%)
	// If reduction is too small, we haven't narrowed the problem enough
	assert.That(metrics.TotalReduction > 0.5, 
		fmt.Sprintf("total reduction must be >50%% to be manageable (got %.2f%%)", metrics.TotalReduction*100))

	// ASSUMPTION: Confidence should be >= 95% for production use
	// If confidence is lower, assumption about filter effectiveness is falsified
	if metrics.Confidence < 0.95 {
		logging.LogCheckpoint("CalculateReductionMetrics", map[string]interface{}{
			"warning":        "confidence_below_threshold",
			"confidence":     metrics.Confidence,
			"threshold":      0.95,
			"note":           "may_need_additional_filtering",
		})
	}

	op.EndOperationWithResult(map[string]interface{}{
		"initial_space":   metrics.InitialSpace,
		"final_space":     metrics.FinalSpace,
		"total_reduction": metrics.TotalReduction,
		"confidence":      metrics.Confidence,
		"steps":           len(metrics.StepReductions),
	})

	return metrics, nil
}

// calculateReduction calculates reduction factor: (before - after) / before
func calculateReduction(before, after int) float64 {
	if before == 0 {
		return 0.0
	}
	return float64(before-after) / float64(before)
}

// calculateConfidence calculates confidence that result satisfies request
// ASSUMPTION: Confidence formula is C = f(reduction, completeness, relevance)
// Simplified version: confidence based on reduction and step assumptions
func calculateConfidence(metrics *ReductionMetrics) float64 {
	if metrics.InitialSpace == 0 {
		return 0.0
	}

	// Base confidence from total reduction
	baseConfidence := metrics.TotalReduction

	// Adjust based on whether assumptions were met at each step
	assumptionBonus := 0.0
	metCount := 0
	for _, step := range metrics.StepReductions {
		if step.AssumptionMet {
			metCount++
		}
	}
	if len(metrics.StepReductions) > 0 {
		assumptionBonus = float64(metCount) / float64(len(metrics.StepReductions)) * 0.1 // Up to 10% bonus
	}

	confidence := baseConfidence + assumptionBonus
	if confidence > 1.0 {
		confidence = 1.0
	}

	return confidence
}
