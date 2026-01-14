import { Button, Datepicker, Label } from "flowbite-react";
import { useEffect, useState } from "react";
import { DateRange } from "../../types/query";
import { ProgressIndicator } from "../shared/ProgressIndicator";

interface DateRangeStepProps {
  dateRange: DateRange;
  onDateRangeChange: (range: DateRange) => void;
  onNext: () => void;
  onBack: () => void;
}

export function DateRangeStep({
  dateRange,
  onDateRangeChange,
  onNext,
  onBack,
}: DateRangeStepProps) {
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    // Validate date range
    if (dateRange.start && dateRange.end) {
      if (dateRange.end < dateRange.start) {
        setError("End date must be after start date");
      } else {
        setError(null);
      }
    } else {
      setError(null);
    }
  }, [dateRange]);

  const isValid =
    dateRange.start !== null && dateRange.end !== null && error === null;

  const handleStartDateChange = (date: Date) => {
    onDateRangeChange({
      ...dateRange,
      start: date,
    });
  };

  const handleEndDateChange = (date: Date) => {
    onDateRangeChange({
      ...dateRange,
      end: date,
    });
  };

  return (
    <div className="max-w-2xl mx-auto px-4 py-8">
      <ProgressIndicator current={2} total={5} />

      <div className="mb-8">
        <h2 className="text-3xl font-bold text-gray-900 mb-3">
          Select date range
        </h2>
        <p className="text-gray-600">Choose the time period for your query</p>
      </div>

      <div className="grid grid-cols-1 md:grid-cols-2 gap-6 mb-6">
        <div>
          <Label htmlFor="start-date" className="mb-2 block">
            Start Date
          </Label>
          <Datepicker
            id="start-date"
            value={dateRange.start?.toLocaleDateString() || ""}
            onSelectedDateChanged={handleStartDateChange}
            maxDate={dateRange.end || new Date()}
          />
        </div>

        <div>
          <Label htmlFor="end-date" className="mb-2 block">
            End Date
          </Label>
          <Datepicker
            id="end-date"
            value={dateRange.end?.toLocaleDateString() || ""}
            onSelectedDateChanged={handleEndDateChange}
            minDate={dateRange.start || undefined}
            maxDate={new Date()}
          />
        </div>
      </div>

      {error && (
        <div className="mb-6 p-4 bg-red-50 border border-red-200 rounded-lg">
          <p className="text-red-600 text-sm">{error}</p>
        </div>
      )}

      <div className="flex justify-between">
        <Button color="gray" onClick={onBack}>
          Back
        </Button>
        <Button onClick={onNext} disabled={!isValid} className="px-6">
          Next
        </Button>
      </div>
    </div>
  );
}
