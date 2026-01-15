// ../App.tsx
import { Button, Card, Spinner } from "flowbite-react";
import { useEffect, useRef, useState } from "react";
import { CategoryStep } from "./components/steps/CategoryStep";
import { ProductionRequestStep } from "./components/steps/ProductionRequestStep";
import { DataCategory } from "./types/query";

type Step = "request" | "categories" | "dashboard";

type ProductionRequest = {
  id: number;
  title: string;
  description: string;
};

// [ ] TODO: Update with go structure
type FileItem = {
  id?: string;
  filename?: string;
  path?: string;
  category?: string;
  [k: string]: any;
};

function App() {
  const [step, setStep] = useState<Step>("request");
  const [requests, setRequests] = useState<ProductionRequest[]>([]);
  const [selectedRequest, setSelectedRequest] =
    useState<ProductionRequest | null>(null);

  const [selectedCategories, setSelectedCategories] = useState<DataCategory[]>(
    []
  );
  const [categories, setCategories] = useState<DataCategory[]>([
    "Email",
    "Claims",
    "Other",
  ]);

  const [files, setFiles] = useState<FileItem[]>([]);
  const [loadingRequests, setLoadingRequests] = useState(true);
  const [loadingCategories, setLoadingCategories] = useState(false);
  const [loadingFiles, setLoadingFiles] = useState(false);

  // Guard against stale async responses (better than reqSeq state)
  // [ ] TODO: confirm this is connected
  const latestFilesReq = useRef(0);

  // --- Load requests from Wails / JSON ---
  // [ ] TODO: map flow from ProductionRequestSelectionPage
  useEffect(() => {
    const loadRequests = async () => {
      setLoadingRequests(true);
      try {
        const { GetProductionRequests } = await import(
          "../wailsjs/go/main/App"
        );
        const res = await GetProductionRequests();
        setRequests(res || []);
      } catch (err) {
        console.error(err);
        setRequests([]);
      } finally {
        setLoadingRequests(false);
      }
    };
    loadRequests();
  }, []);

  // --- Load categories ---
  const loadCategories = async () => {
    setLoadingCategories(true);
    try {
      const { GetCategories } = await import("../wailsjs/go/main/App");
      const res = await GetCategories();

      const rawCats: string[] = Array.isArray(res)
        ? res
        : res && typeof res === "object"
        ? Object.keys(res)
        : [];
      // Map backend strings to DataCategory
      const mapped: DataCategory[] = rawCats.map((cat) => {
        const c = String(cat).toLowerCase();
        if (c === "email") return "Email";
        if (c === "claim" || c === "claims") return "Claims";
        return "Other";
      });
      setCategories(mapped.length ? mapped : ["Email", "Claims", "Other"]);
    } catch (err) {
      console.error(err);
      setCategories(["Email", "Claims", "Other"]);
    } finally {
      setLoadingCategories(false);
    }
  };

  // --- Start from ProductionRequestSelectionPage ---
  // This page only provides a request number. We map it to selectedRequest.
  // [ ] TODO: confirm logic

  const handleStartRequest = async (requestNumber: number) => {
    // Prefer real request object if you have it; fallback to a minimal placeholder
    const chosen =
      requests.find((r) => r.id === requestNumber) ||
      ({
        id: requestNumber,
        title: `Production Request ${requestNumber}`,
        description: "",
      } as ProductionRequest);

    setSelectedRequest(chosen);

    // Proceed to categories
    await loadCategories();
    setStep("categories");
  };

  // --- Handle Next from request page ---
  // const handleNextFromRequest = async () => {
  //   if (!selectedRequest) return;
  //   await loadCategories();
  //   setStep("categories");
  // };

  const [reqSeq, setReqSeq] = useState(0);

  // --- Handle Next from category page ---
  const handleNextFromCategories = async () => {
    if (!selectedCategories.length) return;

    const mySeq = reqSeq + 1;
    // setReqSeq(mySeq);
    const myReq = latestFilesReq.current;

    setLoadingFiles(true);
    try {
      const { GetFiles } = await import("../wailsjs/go/main/App");
      // Backend expects lowercase strings
      const catStrings = selectedCategories.map((c) =>
        c === "Email" ? "email" : c === "Claims" ? "claim" : "other"
      );

      // [ ] TODO: checkpoint
      console.log("Selected categories (UI):", selectedCategories);
      console.log("Selected categories (backend):", catStrings);

      const res = await GetFiles(catStrings);

      /*
      Confirm file object structure
      */
     console.log("GetFiles response isArray:", Array.isArray(res));
     console.log("GetFiles length:", Array.isArray(res) ? res.length : "n/a");
     console.log("GetFiles[0]:", Array.isArray(res) ? res[0] : res);

      setFiles((prev) => {
        if (mySeq !== reqSeq + 1) return prev;
        return res || [];
      });
      setStep("dashboard");
    } catch (err) {
      console.error(err);
      alert("Error loading files: " + err);
    } finally {
      setLoadingFiles(false);
    }
  };

  const handleCategoryToggle = (category: DataCategory) => {
    setSelectedCategories((prev) =>
      prev.includes(category)
        ? prev.filter((c) => c !== category)
        : [...prev, category]
    );
  };

  // --- Step Rendering ---
  return (
    <div className="min-h-screen bg-gray-50 p-8">
      <div className="max-w-4xl mx-auto">
        <h1 className="text-4xl font-bold text-gray-900 mb-8">
          Signal from Noise
        </h1>

        {step === "request" && (
          <ProductionRequestStep
            requests={requests}
            selected={selectedRequest}
            onSelect={setSelectedRequest}
            onNext={handleNextFromRequest}
            loading={loadingRequests}
          />
        )}

        {step === "categories" && (
          <CategoryStep
            categories={categories}
            selected={selectedCategories}
            onToggle={handleCategoryToggle}
            onNext={handleNextFromCategories}
            onBack={() => setStep("request")}
            minSelection={1}
          />
        )}

        {step === "dashboard" && (
          <div className="max-w-3xl mx-auto">
            <h2 className="text-2xl font-semibold text-gray-700 mb-6">
              Dashboard
            </h2>

            {loadingFiles ? (
              <div className="flex justify-center items-center py-12 space-x-4">
                <Spinner size="xl" />
                <span className="text-gray-600">Loading files...</span>
              </div>
            ) : (
              <>
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
                  <h3 className="text-xl font-bold text-gray-900 mb-4">
                    Files
                  </h3>
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
                          <tr
                            key={
                              file.path ??
                              file.id ??
                              `${file.category}:${file.filename}`
                            }
                            className="border-b hover:bg-gray-50"
                          >
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
                        Showing first 50 of {files.length.toLocaleString()}{" "}
                        files
                      </p>
                    )}
                  </div>
                </Card>

                <div className="flex justify-between">
                  <Button color="gray" onClick={() => setStep("categories")}>
                    Back
                  </Button>
                </div>
              </>
            )}
          </div>
        )}
      </div>
    </div>
  );
}

export default App;
