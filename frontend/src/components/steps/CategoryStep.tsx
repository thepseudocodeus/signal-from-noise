import { Button, Checkbox, Label } from "flowbite-react";
import { DataCategory } from "../../types/query";
import { ProgressIndicator } from "../shared/ProgressIndicator";

interface CategoryStepProps {
  categories: DataCategory[];
  selected: DataCategory[];
  onToggle: (category: DataCategory) => void;
  onNext: () => void;
  onBack: () => void;
  minSelection?: number;
}

export function CategoryStep({
  categories,
  selected,
  onToggle,
  onNext,
  onBack,
  minSelection = 1,
}: CategoryStepProps) {
  const isValid = selected.length >= minSelection;

  const categoryLabels: Record<DataCategory, string> = {
    Claims: "Claims",
    Email: "Email",
    Other: "Other",
  };

  return (
    <div className="max-w-2xl mx-auto px-4 py-8">
      <ProgressIndicator current={3} total={5} />

      <div className="mb-8">
        <h2 className="text-3xl font-bold text-gray-900 mb-3">
          Select categories
        </h2>
        <p className="text-gray-600">
          Choose at least {minSelection} category{minSelection > 1 ? "s" : ""}
        </p>
      </div>

      <div className="space-y-4 mb-8">
        {categories.map((category) => (
          <div
            key={category}
            className="flex items-center p-4 border border-gray-200 rounded-lg hover:bg-gray-50 cursor-pointer"
            onClick={() => onToggle(category)}
          >
            <Checkbox
              id={category}
              checked={selected.includes(category)}
              onChange={() => onToggle(category)}
              className="mr-4"
            />
            <Label
              htmlFor={category}
              className="text-lg font-medium text-gray-900 cursor-pointer flex-grow"
            >
              {categoryLabels[category]}
            </Label>
          </div>
        ))}
      </div>

      {!isValid && (
        <div className="mb-6 p-4 bg-yellow-50 border border-yellow-200 rounded-lg">
          <p className="text-yellow-700 text-sm">
            Please select at least {minSelection} category
            {minSelection > 1 ? "s" : ""}
          </p>
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
