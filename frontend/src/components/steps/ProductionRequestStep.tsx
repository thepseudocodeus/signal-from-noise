import { Button, Card } from 'flowbite-react';
import { ProductionRequest } from '../../types/query';
import { ProgressIndicator } from '../shared/ProgressIndicator';

interface ProductionRequestStepProps {
  requests: ProductionRequest[];
  selected: ProductionRequest | null;
  onSelect: (request: ProductionRequest) => void;
  onNext: () => void;
  onBack: () => void;
}

export function ProductionRequestStep({
  requests,
  selected,
  onSelect,
  onNext,
  onBack
}: ProductionRequestStepProps) {
  const isValid = selected !== null;

  return (
    <div className="max-w-3xl mx-auto px-4 py-8">
      <ProgressIndicator current={1} total={3} />

      <div className="mb-8">
        <h2 className="text-3xl font-bold text-gray-900 mb-3">
          Select Production Request
        </h2>
        <p className="text-gray-600">
          Choose one production request or explore with a custom search
        </p>
      </div>

      <div className="grid grid-cols-1 md:grid-cols-2 gap-4 mb-8">
        {requests.map((request) => (
          <Card
            key={request.id}
            className={`cursor-pointer transition-all hover:shadow-lg ${
              selected?.id === request.id
                ? 'ring-4 ring-blue-500 border-blue-500'
                : 'hover:border-gray-300'
            }`}
            onClick={() => onSelect(request)}
          >
            <h3 className="text-lg font-semibold text-gray-900 mb-2">
              {request.title}
            </h3>
            <p className="text-sm text-gray-600">
              {request.description}
            </p>
          </Card>
        ))}
      </div>

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
