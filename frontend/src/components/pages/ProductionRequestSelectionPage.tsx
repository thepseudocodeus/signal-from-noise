import { Button } from "flowbite-react";
import { useState } from "react";

const REQUEST_TITLE = "Production Request";
const REQUEST_DESCRIPTIONS = [
  "Description 1",
  "Description 2",
  "Description 3",
  "Description 4",
  "Description 5",
  "Description 6",
  "Description 7",
  "Description 8",
  "Description 9",
  "Description 10",
  "Description 11",
  "Description 12",
  "Description 13",
  "Description 14",
  "Description 15",
  "Description 16",
  "Description 17",
  "Description 18",
  "Description 19",
  "Description 20",
];

const make_requests = () => {
  return REQUEST_DESCRIPTIONS.map((description, index) => ({
    id: index + 1,
    title: `${REQUEST_TITLE} ${index + 1}`,
    description: description,
  }));
};

const PRODUCTION_REQUESTS = make_requests();

interface ProductionRequestSelectionPageProps {
  onStart: (requestNumber: number) => void;
}

export function ProductionRequestSelectionPage({
  onStart,
}: ProductionRequestSelectionPageProps) {
  const [selectedRequest, setSelectedRequest] = useState<number | null>(null);
  const [isOpen, setIsOpen] = useState(false);

  const productionRequests = PRODUCTION_REQUESTS;

  const handleSelect = (requestNumber: number) => {
    setSelectedRequest(requestNumber);
    setIsOpen(false);
  };

  const handleStart = () => {
    if (selectedRequest) {
      onStart(selectedRequest);
    }
  };

  return (
    <main className="min-h-screen flex items-center justify-center px-6 bg-white">
      <div className="w-full max-w-xl text-center">
        {/* Eyebrow */}
        <span className="mb-4 inline-block text-xs tracking-widest uppercase text-blue-600/70">
          Production Portal
        </span>

        {/* Title */}
        <h1 className="text-4xl md:text-5xl font-medium mb-4 text-gray-900">
          Start a Production Request
        </h1>

        {/* Subtext */}
        <p className="text-gray-600 mb-10 max-w-md mx-auto">
          Select a production request number to begin your application.
        </p>

        {/* Dropdown */}
        <div className="mb-8 relative inline-block w-full max-w-sm">
          <button
            type="button"
            onClick={() => setIsOpen(!isOpen)}
            className="inline-flex items-center justify-between w-full rounded-md border border-gray-200 bg-white px-5 py-4 text-sm font-medium text-gray-700 shadow-sm hover:bg-gray-50 focus:outline-none focus:ring-2 focus:ring-blue-500"
          >
            <span>
              {selectedRequest
                ? `Production Request No #${selectedRequest}`
                : "Choose a Production Request"}
            </span>
            <svg
              className={`w-4 h-4 ml-2 transition-transform ${
                isOpen ? "rotate-180" : ""
              }`}
              fill="none"
              stroke="currentColor"
              viewBox="0 0 24 24"
            >
              <path
                strokeLinecap="round"
                strokeLinejoin="round"
                strokeWidth={2}
                d="M19 9l-7 7-7-7"
              />
            </svg>
          </button>

          {/* Dropdown Menu */}
          {isOpen && (
            <>
              <div
                className="fixed inset-0 z-10"
                onClick={() => setIsOpen(false)}
              ></div>
              <div className="absolute z-20 mt-2 w-full rounded-md border border-gray-200 bg-white shadow-lg max-h-72 overflow-y-auto">
                <div className="px-4 py-2 text-xs font-semibold text-gray-500 uppercase border-b border-gray-200">
                  PRODUCTION REQUESTS
                </div>
                <ul className="text-sm text-left">
                  {productionRequests.map((requestNumber) => (
                    <li key={requestNumber.id}>
                      <button
                        type="button"
                        onClick={() => handleSelect(requestNumber.id)}
                        className={`block w-full px-4 py-3 text-left hover:bg-gray-50 ${
                          selectedRequest === requestNumber.id
                            ? "bg-blue-50 text-blue-600"
                            : "text-gray-700"
                        }`}
                      >
                        Production Request No #{requestNumber.id}
                      </button>
                    </li>
                  ))}
                </ul>
              </div>
            </>
          )}
        </div>

        {/* Helper */}
        <p className="mb-8 text-xs text-gray-500">
          This will take approximately 3–5 minutes to complete.
        </p>

        {/* Start Button - Only shown when selection made */}
        {selectedRequest && (
          <Button size="xl" onClick={handleStart} className="px-8 py-3">
            Start
          </Button>
        )}
      </div>
    </main>
  );
}
