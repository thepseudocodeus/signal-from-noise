import { useState, KeyboardEvent } from 'react';
import { Button, TextInput, Label, Badge } from 'flowbite-react';
import { ProgressIndicator } from '../shared/ProgressIndicator';

interface KeywordsStepProps {
  keywords: string[];
  onKeywordsChange: (keywords: string[]) => void;
  onNext: () => void;
  onBack: () => void;
  optional?: boolean;
}

export function KeywordsStep({
  keywords,
  onKeywordsChange,
  onNext,
  onBack,
  optional = true
}: KeywordsStepProps) {
  const [inputValue, setInputValue] = useState('');

  const handleAddKeyword = () => {
    const trimmed = inputValue.trim();
    if (trimmed && !keywords.includes(trimmed)) {
      onKeywordsChange([...keywords, trimmed]);
      setInputValue('');
    }
  };

  const handleRemoveKeyword = (keyword: string) => {
    onKeywordsChange(keywords.filter(k => k !== keyword));
  };

  const handleKeyPress = (e: KeyboardEvent<HTMLInputElement>) => {
    if (e.key === 'Enter') {
      e.preventDefault();
      handleAddKeyword();
    }
  };

  return (
    <div className="max-w-2xl mx-auto px-4 py-8">
      <ProgressIndicator current={4} total={5} />

      <div className="mb-8">
        <h2 className="text-3xl font-bold text-gray-900 mb-3">
          Keywords {optional && <span className="text-gray-500 text-xl">(Optional)</span>}
        </h2>
        <p className="text-gray-600">
          Add keywords to filter your search
        </p>
      </div>

      <div className="mb-6">
        <Label htmlFor="keyword-input" className="mb-2 block">
          Add Keyword
        </Label>
        <div className="flex gap-2">
          <TextInput
            id="keyword-input"
            type="text"
            value={inputValue}
            onChange={(e) => setInputValue(e.target.value)}
            onKeyPress={handleKeyPress}
            placeholder="Enter keyword and press Enter"
            className="flex-grow"
          />
          <Button onClick={handleAddKeyword}>
            Add
          </Button>
        </div>
      </div>

      {keywords.length > 0 && (
        <div className="mb-6">
          <Label className="mb-2 block">Selected Keywords</Label>
          <div className="flex flex-wrap gap-2">
            {keywords.map((keyword) => (
              <Badge
                key={keyword}
                color="info"
                className="text-sm py-2 px-3"
              >
                {keyword}
                <button
                  onClick={() => handleRemoveKeyword(keyword)}
                  className="ml-2 text-blue-600 hover:text-blue-800"
                >
                  ×
                </button>
              </Badge>
            ))}
          </div>
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
          className="px-6"
        >
          Next
        </Button>
      </div>
    </div>
  );
}
