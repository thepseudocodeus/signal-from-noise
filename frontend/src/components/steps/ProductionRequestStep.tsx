import { Button, Card, Label, Select, Spinner } from "flowbite-react";
import { ProgressIndicator } from "../shared/ProgressIndicator";
import { ProductionRequest } from "../../types/query";

interface ProductionRequestStepProps {
  requests: ProductionRequest[];
  selected: ProductionRequest | null;
  onSelect: (request: ProductionRequest) => void;
  onNext: () => void;
  onBack?: () => void;
  loading?: boolean;
}

export function ProductionRequestStep({
  requests,
  selected,
  onSelect,
  onNext,
  onBack,
  loading = false,
}: ProductionRequestStepProps) {
  const isValid = selected !== null;

  return (
    <div className="max-w-2xl mx-auto px-4 py-8 flex flex-col">
      <ProgressIndicator current={1} total={3} />

      <div className="mb-8">
        <h2 className="text-3xl font-bold text-gray-900 mb-3">
          Select Production Request
        </h2>
        <p className="text-gray-600">
          Choose one production request from the dropdown below
        </p>
      </div>

      <div className="flex-grow">
        {loading ? (
          <div className="flex justify-center items-center py-12 space-x-4">
            <Spinner size="xl" />
            <span className="text-gray-600">Loading requests...</span>
          </div>
        ) : (
          <Card className="mb-8">
            <div className="mb-4">
              <Label
                htmlFor="request-select"
                value="Select a production request"
              />
            </div>
            <Select
              id="request-select"
              value={selected?.id || ""}
              onChange={(e) => {
                const id = parseInt(e.target.value, 10);
                const req = requests.find((r) => r.id === id) || null;
                onSelect(req);
              }}
            >
              <option value="">Choose a request...</option>
              {requests.map((req) => (
                <option key={req.id} value={req.id}>
                  {req.title || `REQUEST FOR PRODUCTION NO: ${req.id}`}
                </option>
              ))}
            </Select>
            {selected && (
              <div className="mt-4 p-4 bg-blue-50 rounded-lg">
                <p className="text-sm font-semibold text-gray-900 mb-2">
                  {selected.title}
                </p>
                <p className="text-sm text-gray-700">{selected.description}</p>
              </div>
            )}
          </Card>
        )}
      </div>

      {/* Always reserve space for buttons */}
      <div className="flex justify-between mt-4">
        {onBack ? (
          <Button color="gray" onClick={onBack}>
            Back
          </Button>
        ) : (
          <div />
        )}
        <Button onClick={onNext} className="px-6" disabled={!isValid}>
          Next
        </Button>
      </div>
    </div>
  );
}
