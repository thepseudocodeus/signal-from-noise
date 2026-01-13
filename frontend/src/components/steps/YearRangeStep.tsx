import { Button, Label, Select } from 'flowbite-react';
import { YearRange } from '../../types/query';
import { ProgressIndicator } from '../shared/ProgressIndicator';

interface YearRangeStepProps {
  yearRange: YearRange;
  onYearRangeChange: (range: YearRange) => void;
  onNext: () => void;
  onBack: () => void;
}

const CURRENT_YEAR = new Date().getFullYear();
const START_YEAR = 2000;
const YEARS = Array.from({ length: CURRENT_YEAR - START_YEAR + 1 }, (_, i) => START_YEAR + i).reverse();

export function YearRangeStep({
  yearRange,
  onYearRangeChange,
  onNext,
  onBack
}: YearRangeStepProps) {
  const isValid = yearRange.startYear !== null &&
                  yearRange.endYear !== null &&
                  (yearRange.endYear >= (yearRange.startYear || 0));

  const handleStartYearChange = (e: React.ChangeEvent<HTMLSelectElement>) => {
    const year = e.target.value ? parseInt(e.target.value) : null;
    onYearRangeChange({
      ...yearRange,
      startYear: year
    });
  };

  const handleEndYearChange = (e: React.ChangeEvent<HTMLSelectElement>) => {
    const year = e.target.value ? parseInt(e.target.value) : null;
    onYearRangeChange({
      ...yearRange,
      endYear: year
    });
  };

  return (
    <div className="max-w-2xl mx-auto px-4 py-8">
      <ProgressIndicator current={2} total={3} />

      <div className="mb-8">
        <h2 className="text-3xl font-bold text-gray-900 mb-3">
          Select Year Range
        </h2>
        <p className="text-gray-600">
          Choose the time period for your search
        </p>
      </div>

      <div className="grid grid-cols-1 md:grid-cols-2 gap-6 mb-6">
        <div>
          <Label htmlFor="start-year" className="mb-2 block">
            Start Year
          </Label>
          <Select
            id="start-year"
            value={yearRange.startYear || ''}
            onChange={handleStartYearChange}
          >
            <option value="">Select start year</option>
            {YEARS.map(year => (
              <option key={year} value={year}>{year}</option>
            ))}
          </Select>
        </div>

        <div>
          <Label htmlFor="end-year" className="mb-2 block">
            End Year
          </Label>
          <Select
            id="end-year"
            value={yearRange.endYear || ''}
            onChange={handleEndYearChange}
          >
            <option value="">Select end year</option>
            {YEARS.map(year => (
              <option key={year} value={year}>{year}</option>
            ))}
          </Select>
        </div>
      </div>

      {yearRange.startYear && yearRange.endYear && yearRange.endYear < yearRange.startYear && (
        <div className="mb-6 p-4 bg-red-50 border border-red-200 rounded-lg">
          <p className="text-red-600 text-sm">End year must be after start year</p>
        </div>
      )}

      <div className="flex justify-between">
        <Button
          color="gray"
          onClick={onBack}
        >
          Back
        </Button>
        <Button
          onClick={onNext}
          disabled={!isValid}
          className="px-6"
        >
          Next
        </Button>
      </div>
    </div>
  );
}
