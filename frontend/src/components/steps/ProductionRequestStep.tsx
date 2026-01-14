import { Button, Card } from "flowbite-react";
import { ProductionRequest } from "../../types/query";
import { ProgressIndicator } from "../shared/ProgressIndicator";

interface ProductionRequestStepProps {
  requests: ProductionRequest[];
  selected: ProductionRequest | null;
  onSelect: (request: ProductionRequest) => void;
  onNext: () => void;
}

export function ProductionRequestStep({
  requests,
  selected,
  onSelect,
  onNext,
}: ProductionRequestStepProps) {
  const isValid = selected !== null;

  return (
    <div className="max-w-3xl mx-auto px-4 py-8">
      <ProgressIndicator current={1} total={3} />

      <div className="mb-8">
        <h2 className="text-3xl font-bold text-gray-900 mb-3">
          Select Production Request
        </h2>
        <p className="text-gray-600">Choose one production request</p>
      </div>

      <div className="grid grid-cols-1 md:grid-cols-2 gap-4 mb-8">
        {requests.map((request) => {
          const isSelected = selected?.id === request.id;

          return (
            <Card
              key={request.id}
              role="button"
              tabIndex={0}
              onClick={() => onSelect(request)}
              onKeyDown={(e) => {
                if (e.key === "Enter" || e.key === " ") {
                  e.preventDefault();
                  onSelect(request);
                }
              }}
              className={`cursor-pointer transition-all hover:shadow-lg ${
                isSelected
                  ? "ring-2 ring-blue-500 border-blue-500"
                  : "hover:border-gray-300"
              }`}
            >
              <h3 className="text-lg font-semibold text-gray-900 mb-2">
                {request.title}
              </h3>
              <p className="text-sm text-gray-600">{request.description}</p>
            </Card>
          );
        })}
      </div>

      <div className="flex justify-end">
        <Button onClick={onNext} className="px-6">
          Next
        </Button>
        {/* <Button onClick={onNext} disabled={!isValid} className="px-6">
          Next
        </Button> */}
      </div>
    </div>
  );
}
