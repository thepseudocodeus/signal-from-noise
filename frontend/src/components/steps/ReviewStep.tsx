import { Button, Card, Badge } from 'flowbite-react';
import { UserQuery, Category } from '../../types/query';
import { ProgressIndicator } from '../shared/ProgressIndicator';
import { format } from 'date-fns';

interface ReviewStepProps {
  query: UserQuery;
  onConfirm: () => void;
  onEdit: (step: number) => void;
}

export function ReviewStep({ query, onConfirm, onEdit }: ReviewStepProps) {
  const formatDate = (date: Date | null) => {
    if (!date) return 'Not set';
    return format(date, 'MMM dd, yyyy');
  };

  return (
    <div className="max-w-3xl mx-auto px-4 py-8">
      <ProgressIndicator current={5} total={5} />

      <div className="mb-8">
        <h2 className="text-3xl font-bold text-gray-900 mb-3">
          Review your selections
        </h2>
        <p className="text-gray-600">
          Please review your query before submitting
        </p>
      </div>

      <div className="space-y-4 mb-8">
        <Card>
          <div className="flex justify-between items-start">
            <div>
              <h3 className="font-semibold text-gray-900 mb-1">Query Type</h3>
              <p className="text-gray-600">
                {query.informationRequest?.title || 'Not selected'}
              </p>
              {query.informationRequest?.description && (
                <p className="text-sm text-gray-500 mt-1">
                  {query.informationRequest.description}
                </p>
              )}
            </div>
            <Button
              color="light"
              size="sm"
              onClick={() => onEdit(1)}
            >
              Edit
            </Button>
          </div>
        </Card>

        <Card>
          <div className="flex justify-between items-start">
            <div>
              <h3 className="font-semibold text-gray-900 mb-1">Date Range</h3>
              <p className="text-gray-600">
                {formatDate(query.dateRange.start)} - {formatDate(query.dateRange.end)}
              </p>
            </div>
            <Button
              color="light"
              size="sm"
              onClick={() => onEdit(2)}
            >
              Edit
            </Button>
          </div>
        </Card>

        <Card>
          <div className="flex justify-between items-start">
            <div>
              <h3 className="font-semibold text-gray-900 mb-2">Categories</h3>
              <div className="flex flex-wrap gap-2">
                {query.categories.length > 0 ? (
                  query.categories.map((cat: Category) => (
                    <Badge key={cat} color="info">
                      {cat}
                    </Badge>
                  ))
                ) : (
                  <span className="text-gray-500">Not selected</span>
                )}
              </div>
            </div>
            <Button
              color="light"
              size="sm"
              onClick={() => onEdit(3)}
            >
              Edit
            </Button>
          </div>
        </Card>

        <Card>
          <div className="flex justify-between items-start">
            <div>
              <h3 className="font-semibold text-gray-900 mb-2">Keywords</h3>
              <div className="flex flex-wrap gap-2">
                {query.keywords.length > 0 ? (
                  query.keywords.map((keyword: string) => (
                    <Badge key={keyword} color="gray">
                      {keyword}
                    </Badge>
                  ))
                ) : (
                  <span className="text-gray-500">None</span>
                )}
              </div>
            </div>
            <Button
              color="light"
              size="sm"
              onClick={() => onEdit(4)}
            >
              Edit
            </Button>
          </div>
        </Card>
      </div>

      <div className="flex justify-between">
        <Button
          color="gray"
          onClick={() => onEdit(4)}
        >
          Back
        </Button>
        <Button
          onClick={onConfirm}
          size="xl"
          className="px-8"
        >
          Confirm & Submit
        </Button>
      </div>
    </div>
  );
}
