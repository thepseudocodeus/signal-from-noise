import { Button, Card, Label, Select, Spinner } from "flowbite-react";
import { useEffect, useState } from "react";
import { ProgressIndicator } from "./components/shared/ProgressIndicator";
import { CategoryStep } from "./components/steps/CategoryStep";
import { DataCategory } from "./types/query";

type Step = "request" | "categories" | "dashboard";

type ProductionRequest = {
  id: number;
  title: string;
  description: string;
};

function App() {
  const [step, setStep] = useState<Step>("request");
  const [selectedRequest, setSelectedRequest] = useState<string>("");
  const [selectedCategories, setSelectedCategories] = useState<DataCategory[]>(
    []
  );
  const [categories, setCategories] = useState<Record<string, number>>({});
  const [loading, setLoading] = useState(false);
  const [files, setFiles] = useState<any[]>([]);
  const [requests, setRequests] = useState<ProductionRequest[]>([]);

  // Load production requests on mount - wait for Wails runtime
  useEffect(() => {
    const waitForWailsAndLoad = async () => {
      // Wait for Wails runtime to initialize
      let attempts = 0;
      const maxAttempts = 50;

      while (attempts < maxAttempts) {
        if (typeof window !== "undefined" && (window as any).go?.main?.App) {
          break; // Runtime is ready
        }
        await new Promise((resolve) => setTimeout(resolve, 100));
        attempts++;
      }

      if (attempts >= maxAttempts) {
        console.warn("Wails runtime not available after waiting");
        setLoading(false);
        return;
      }

      // Runtime is ready, load requests
      setLoading(true);
      try {
        const { GetProductionRequests } = await import(
          "../wailsjs/go/main/App"
        );
        const loadedRequests = await GetProductionRequests();

        // Ensure all requests have required fields
        const validRequests = (loadedRequests || []).map((req: any) => ({
          id: req.id,
          title: req.title || `REQUEST FOR PRODUCTION NO: ${req.id}`,
          description: req.description || "",
        }));

        setRequests(validRequests);
      } catch (err: any) {
        console.error("Error loading production requests:", err);
        setRequests([]);
      } finally {
        setLoading(false);
      }
    };

    waitForWailsAndLoad();
  }, []);

  // Load categories when moving to categories step
  useEffect(() => {
    if (step === "categories") {
      loadCategories();
    }
  }, [step]);

  const loadCategories = async () => {
    setLoading(true);
    try {
      const { GetCategories } = await import("../wailsjs/go/main/App");
      const cats = await GetCategories();
      setCategories(cats || {});
    } catch (err) {
      console.error("Error loading categories:", err);
    } finally {
      setLoading(false);
    }
  };

  const handleCategoryToggle = (category: DataCategory) => {
    setSelectedCategories((prev) =>
      prev.includes(category)
        ? prev.filter((c) => c !== category)
        : [...prev, category]
    );
  };

  const handleRequestNext = () => {
    if (selectedRequest) {
      setStep("categories");
    }
  };

  const selectedRequestObj = requests.find(
    (r) => r.id.toString() === selectedRequest
  );

  const handleCategoriesNext = async () => {
    if (selectedCategories.length === 0) return;

    setLoading(true);
    try {
      const { GetFiles } = await import("../wailsjs/go/main/App");
      // Convert DataCategory to lowercase string for backend
      const categoryStrings = selectedCategories.map((cat) => {
        if (cat === "Email") return "email";
        if (cat === "Claims") return "claim";
        return "other";
      });
      const fileList = await GetFiles(categoryStrings);
      setFiles(fileList || []);
      setStep("dashboard");
    } catch (err) {
      console.error("Error loading files:", err);
      alert("Error loading files: " + err);
    } finally {
      setLoading(false);
    }
  };

  // Convert category counts to DataCategory array
  // Default to all categories if none loaded yet
  const availableCategories: DataCategory[] =
    Object.keys(categories).length > 0
      ? Object.keys(categories)
          .map((cat) => {
            if (cat === "email") return "Email";
            if (cat === "claim") return "Claims";
            return "Other";
          })
          .filter((cat): cat is DataCategory =>
            ["Email", "Claims", "Other"].includes(cat)
          )
      : ["Email", "Claims", "Other"]; // Default categories

  return (
    <div className="min-h-screen bg-gray-50 p-8">
      <div className="max-w-4xl mx-auto">
        <h1 className="text-4xl font-bold text-gray-900 mb-8">
          Signal from Noise
        </h1>

        {/* Step 1: Request Selection */}
        {step === "request" && (
          <div className="max-w-2xl mx-auto px-4 py-8">
            <ProgressIndicator current={1} total={5} />

            <div className="mb-8">
              <h2 className="text-3xl font-bold text-gray-900 mb-3">
                Select Production Request
              </h2>
              <p className="text-gray-600">
                Choose one production request from the dropdown below
              </p>
            </div>

            {loading && requests.length === 0 ? (
              <div className="flex justify-center items-center py-12">
                <Spinner size="xl" />
                <span className="ml-4 text-gray-600">Loading requests...</span>
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
                  value={selectedRequest}
                  onChange={(e) => setSelectedRequest(e.target.value)}
                >
                  <option value="">Choose a request...</option>
                  {requests.map((req) => (
                    <option key={req.id} value={req.id.toString()}>
                      {req.title || `REQUEST FOR PRODUCTION NO: ${req.id}`}
                    </option>
                  ))}
                </Select>
                {selectedRequestObj && (
                  <div className="mt-4 p-4 bg-blue-50 rounded-lg">
                    <p className="text-sm font-semibold text-gray-900 mb-2">
                      {selectedRequestObj.title ||
                        `REQUEST FOR PRODUCTION NO: ${selectedRequestObj.id}`}
                    </p>
                    <p className="text-sm text-gray-700">
                      {selectedRequestObj.description || ""}
                    </p>
                  </div>
                )}
              </Card>
            )}

            <div className="flex justify-between">
              <Button color="gray" onClick={() => {}} disabled={true}>
                Back
              </Button>
              <Button
                onClick={handleRequestNext}
                disabled={!selectedRequest || requests.length === 0}
                className="px-6"
              >
                Next
              </Button>
            </div>
          </div>
        )}

        {/* Step 2: Category Selection */}
        {step === "categories" && (
          <>
            {loading && availableCategories.length === 0 ? (
              <div className="flex justify-center items-center py-12">
                <Spinner size="xl" />
                <span className="ml-4 text-gray-600">
                  Loading categories...
                </span>
              </div>
            ) : (
              <CategoryStep
                categories={availableCategories}
                selected={selectedCategories}
                onToggle={handleCategoryToggle}
                onNext={handleCategoriesNext}
                onBack={() => setStep("request")}
              />
            )}
          </>
        )}

        {/* Step 3: Dashboard */}
        {step === "dashboard" && (
          <>
            <h2 className="text-2xl font-semibold text-gray-700 mb-6">
              Dashboard
            </h2>

            <div className="grid grid-cols-3 gap-4 mb-6">
              <Card>
                <p className="text-gray-600 mb-2">Total Files</p>
                <p className="text-3xl font-bold text-gray-900">
                  {files.length.toLocaleString()}
                </p>
              </Card>
              <Card>
                <p className="text-gray-600 mb-2">Categories</p>
                <p className="text-3xl font-bold text-gray-900">
                  {selectedCategories.length}
                </p>
              </Card>
              <Card>
                <p className="text-gray-600 mb-2">Expected Zip Size</p>
                <p className="text-3xl font-bold text-gray-900">
                  ~{Math.round(files.length * 0.85).toLocaleString()} files
                </p>
              </Card>
            </div>

            <Card className="mb-6">
              <h3 className="text-xl font-bold text-gray-900 mb-4">Files</h3>
              <div className="overflow-x-auto">
                <table className="w-full text-left">
                  <thead>
                    <tr className="border-b">
                      <th className="p-2">Filename</th>
                      <th className="p-2">Category</th>
                      <th className="p-2">Path</th>
                    </tr>
                  </thead>
                  <tbody>
                    {files.slice(0, 50).map((file, idx) => (
                      <tr key={idx} className="border-b hover:bg-gray-50">
                        <td className="p-2 font-medium">{file.filename}</td>
                        <td className="p-2">
                          <span className="px-2 py-1 bg-blue-100 text-blue-800 rounded text-sm capitalize">
                            {file.category}
                          </span>
                        </td>
                        <td
                          className="p-2 text-sm text-gray-600 truncate max-w-md"
                          title={file.path}
                        >
                          {file.path}
                        </td>
                      </tr>
                    ))}
                  </tbody>
                </table>
                {files.length > 50 && (
                  <p className="p-4 text-gray-600 text-center">
                    Showing first 50 of {files.length.toLocaleString()} files
                  </p>
                )}
              </div>
            </Card>

            <div className="flex justify-between">
              <Button
                onClick={() => setStep("categories")}
                color="gray"
                size="lg"
              >
                Back
              </Button>
            </div>
          </>
        )}
      </div>
    </div>
  );
}

export default App;
