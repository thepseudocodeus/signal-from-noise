import { Alert, Datepicker, Label } from "flowbite-react";
import { useEffect, useState } from "react";

/**
 * DateRangeSelector - A reusable date range selection component
 *
 * Single Responsibility: Date range selection only
 * Explicit Contracts: Clear TypeScript interface
 * Fail Fast: Validation on every change
 * Type Safety: Full TypeScript types
 *
 * Uses Flowbite React Datepicker component (default, no custom replacements)
 */
export interface DateRangeSelectorProps {
  /** Minimum date allowed (from data) */
  minDate: Date;
  /** Maximum date allowed (from data) */
  maxDate: Date;
  /** Currently selected start date */
  startDate: Date | null;
  /** Currently selected end date */
  endDate: Date | null;
  /** Callback when start date changes */
  onStartChange: (date: Date | null) => void;
  /** Callback when end date changes */
  onEndChange: (date: Date | null) => void;
  /** Optional label for start date picker */
  startLabel?: string;
  /** Optional label for end date picker */
  endLabel?: string;
  /** Optional CSS class name */
  className?: string;
}

/**
 * DateRangeSelector Component
 *
 * ASSUMPTION: Flowbite Datepicker component works correctly
 * ASSUMPTION: minDate <= maxDate (data constraint)
 * ASSUMPTION: If startDate and endDate are set, startDate <= endDate (user constraint)
 */
export function DateRangeSelector({
  minDate,
  maxDate,
  startDate,
  endDate,
  onStartChange,
  onEndChange,
  startLabel = "Start Date",
  endLabel = "End Date",
  className = "",
}: DateRangeSelectorProps) {
  const [error, setError] = useState<string | null>(null);

  // ASSUMPTION: minDate must be <= maxDate for component to work correctly
  // If this is false, the date range constraints are invalid
  useEffect(() => {
    if (minDate > maxDate) {
      console.error(
        "DateRangeSelector: minDate > maxDate, constraints are invalid"
      );
    }
  }, [minDate, maxDate]);

  // Validation: End date must be >= start date
  // ASSUMPTION: If both dates are set, startDate <= endDate must be true
  // If this fails, show error and prevent invalid state
  useEffect(() => {
    if (startDate && endDate) {
      if (endDate < startDate) {
        setError("End date must be on or after start date");
      } else {
        setError(null);
      }
    } else {
      setError(null);
    }
  }, [startDate, endDate]);

  /**
   * Handle start date change
   * ASSUMPTION: Datepicker onChange provides Date | null
   * If new start date is after end date, end date will be constrained by minDate
   */
  const handleStartChange = (date: Date | null) => {
    onStartChange(date);

    // If new start date is after current end date, clear end date
    // This prevents invalid state (end < start)
    if (date && endDate && endDate < date) {
      onEndChange(null);
      setError("End date must be on or after start date. End date cleared.");
    }
  };

  /**
   * Handle end date change
   * ASSUMPTION: Datepicker onChange provides Date | null
   * Validation ensures endDate >= startDate
   */
  const handleEndChange = (date: Date | null) => {
    // Fail fast: If date is before start date, reject it
    if (date && startDate && date < startDate) {
      setError("End date must be on or after start date");
      return; // Prevent invalid selection
    }

    setError(null);
    onEndChange(date);
  };

  // Calculate dynamic constraints
  // ASSUMPTION: End date's minDate should be startDate (if set) or data minDate
  // This prevents selecting end date before start date
  const endDateMinDate = startDate || minDate;

  // ASSUMPTION: Start date's maxDate should be endDate (if set) or data maxDate
  // This prevents selecting start date after end date
  const startDateMaxDate = endDate || maxDate;

  return (
    <div className={className}>
      <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
        {/* Start Date Picker */}
        <div>
          <Label htmlFor="start-date" className="mb-2 block">
            {startLabel}
          </Label>
          <Datepicker
            id="start-date"
            value={startDate || null}
            onChange={handleStartChange}
            minDate={minDate}
            maxDate={startDateMaxDate}
            labelClearButton="Clear"
            labelTodayButton="Today"
          />
        </div>

        {/* End Date Picker */}
        <div>
          <Label htmlFor="end-date" className="mb-2 block">
            {endLabel}
          </Label>
          <Datepicker
            id="end-date"
            value={endDate || null}
            onChange={handleEndChange}
            minDate={endDateMinDate}
            maxDate={maxDate}
            labelClearButton="Clear"
            labelTodayButton="Today"
          />
        </div>
      </div>

      {/* Error Display */}
      {error && (
        <Alert color="failure" className="mt-4">
          <span className="font-medium">Validation Error:</span> {error}
        </Alert>
      )}

      {/* Success Display (both dates valid) */}
      {startDate && endDate && !error && (
        <Alert color="success" className="mt-4">
          <span className="font-medium">Date Range Selected:</span>{" "}
          {startDate.toLocaleDateString()} to {endDate.toLocaleDateString()}
        </Alert>
      )}
    </div>
  );
}
