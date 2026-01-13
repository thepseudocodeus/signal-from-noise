import { Button, Card, Badge } from 'flowbite-react';
import { QueryState } from '../../types/query';
import { ProgressIndicator } from '../shared/ProgressIndicator';

interface ResultsStepProps {
  query: QueryState;
  fileCount: number;
  fileSize: string;
  onNewSearch: () => void;
}

export function ResultsStep({ query, fileCount, fileSize, onNewSearch }: ResultsStepProps) {
  return (
    <div className="max-w-3xl mx-auto px-4 py-8">
      <ProgressIndicator current={3} total={3} />

      <div className="mb-8">
        <h2 className="text-3xl font-bold text-gray-900 mb-3">
          Query Results
        </h2>
        <p className="text-gray-600">
          Your search has been completed
        </p>
      </div>

      {/* Results Summary */}
      <Card className="mb-6">
        <div className="space-y-4">
          <div>
            <h3 className="font-semibold text-gray-900 mb-2">Files Found</h3>
            <p className="text-2xl font-bold text-blue-600">{fileCount.toLocaleString()}</p>
          </div>
          <div>
            <h3 className="font-semibold text-gray-900 mb-2">Total Size</h3>
            <p className="text-2xl font-bold text-blue-600">{fileSize}</p>
          </div>
        </div>
      </Card>

      {/* Query Summary */}
      <Card className="mb-6">
        <h3 className="font-semibold text-gray-900 mb-4">Query Summary</h3>
        <div className="space-y-2 text-sm">
          <div className="flex justify-between">
            <span className="text-gray-600">Production Request:</span>
            <span className="font-medium">{query.productionRequest?.title || 'N/A'}</span>
          </div>
          <div className="flex justify-between">
            <span className="text-gray-600">Year Range:</span>
            <span className="font-medium">
              {query.yearRange.startYear} - {query.yearRange.endYear}
            </span>
          </div>
        </div>
      </Card>

      <div className="flex justify-center">
        <Button
          onClick={onNewSearch}
          size="xl"
          className="px-8"
        >
          New Search
        </Button>
      </div>
    </div>
  );
}
