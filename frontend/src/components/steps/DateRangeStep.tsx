import { Button } from "flowbite-react";
import { DateRange } from "../../types/query";
import { DateRangeSelector } from "../shared/DateRangeSelector";
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
  // Determine min/max dates from data
  // ASSUMPTION: Data spans reasonable range (e.g., 2020-2025)
  // In production, these would come from actual data bounds
  const dataMinDate = new Date(2020, 0, 1); // January 1, 2020
  const dataMaxDate = new Date(); // Today

  const handleStartChange = (date: Date | null) => {
    onDateRangeChange({
      ...dateRange,
      start: date,
    });
  };

  const handleEndChange = (date: Date | null) => {
    onDateRangeChange({
      ...dateRange,
      end: date,
    });
  };

  // Validation: Both dates must be selected and valid
  // ASSUMPTION: DateRangeSelector handles validation internally
  // We just check that both dates are set
  const isValid = dateRange.start !== null && dateRange.end !== null;

  return (
    <div className="max-w-2xl mx-auto px-4 py-8">
      <ProgressIndicator current={2} total={5} />

      <div className="mb-8">
        <h2 className="text-3xl font-bold text-gray-900 mb-3">
          Select date range
        </h2>
        <p className="text-gray-600">Choose the time period for your query</p>
      </div>

      <div className="mb-6">
        <DateRangeSelector
          minDate={dataMinDate}
          maxDate={dataMaxDate}
          startDate={dateRange.start}
          endDate={dateRange.end}
          onStartChange={handleStartChange}
          onEndChange={handleEndChange}
          startLabel="Start Date"
          endLabel="End Date"
        />
      </div>

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
