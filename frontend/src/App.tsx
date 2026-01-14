import { useEffect, useState } from "react";
import { Button, Card, Label, Select, Spinner } from "flowbite-react";
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

  const [requests, setRequests] = useState<ProductionRequest[]>([]);
  const [selectedRequest, setSelectedRequest] = useState<number | null>(null);

  const [categories, setCategories] = useState<Record<string, number>>({});
  const [selectedCategories, setSelectedCategories] = useState<DataCategory[]>(
    []
  );

  const [files, setFiles] = useState<any[]>([]);
  const [loading, setLoading] = useState(false);

  /* -------------------------
     Load production requests
  -------------------------- */
  useEffect(() => {
    const load = async () => {
      let attempts = 0;

      while (attempts < 50) {
        if ((window as any)?.go?.main?.App) break;
        await new Promise((r) => setTimeout(r, 100));
        attempts++;
      }

      setLoading(true);
      try {
        const { GetProductionRequests } = await import(
          "../wailsjs/go/main/App"
        );
        const data = await GetProductionRequests();

        setRequests(
          (data || []).map((r: any) => ({
            id: r.id,
            title: r.title || `REQUEST FOR PRODUCTION NO: ${r.id}`,
            description: r.description || "",
          }))
        );
      } catch (err) {
        console.error(err);
      } finally {
        setLoading(false);
      }
    };

    load();
  }, []);

  /* -------------------------
     Load categories on step
  -------------------------- */
  useEffect(() => {
    if (step !== "categories") return;

    const load = async () => {
      setLoading(true);
      try {
        const { GetCategories } = await import("../wailsjs/go/main/App");
        const data = await GetCategories();
        setCategories(data || {});
      } catch (err) {
        console.error(err);
      } finally {
        setLoading(false);
      }
    };

    load();
  }, [step]);

  /* -------------------------
     Navigation handlers
  -------------------------- */
  const nextFromRequest = () => {
    if (selectedRequest !== null) {
      setStep("categories");
    }
  };

  const backToRequest = () => {
    setStep("request");
  };

  const nextFromCategories = async () => {
    if (selectedCategories.length === 0) return;

    setLoading(true);
    try {
      const { GetFiles } = await import("../wailsjs/go/main/App");

      const backendCats = selectedCategories.map((c) =>
        c === "Email" ? "email" : c === "Claims" ? "claim" : "other"
      );

      const data = await GetFiles(backendCats);
      setFiles(data || []);
      setStep("dashboard");
    } catch (err) {
      console.error(err);
    } finally {
      setLoading(false);
    }
  };

  const backToCategories = () => {
    setStep("categories");
  };

  /* -------------------------
     Derived data
  -------------------------- */
  const selectedRequestObj = requests.find((r) => r.id === selectedRequest);

  const availableCategories: DataCategory[] =
    Object.keys(categories).length > 0
      ? Object.keys(categories).map((c) =>
          c === "email" ? "Email" : c === "claim" ? "Claims" : "Other"
        )
      : ["Email", "Claims", "Other"];

  /* -------------------------
     Render
  -------------------------- */
  return (
    <div className="min-h-screen bg-gray-50 p-8">
      <div className="max-w-4xl mx-auto">
        <h1 className="text-4xl font-bold text-gray-900 mb-8">
          Signal from Noise
        </h1>

        {/* STEP 1 — Request */}
        {step === "request" && (
          <div className="max-w-2xl mx-auto">
            <ProgressIndicator current={1} total={3} />

            {loading ? (
              <div className="flex justify-center py-12">
                <Spinner size="xl" />
              </div>
            ) : (
              <Card className="mt-8">
                <Label htmlFor="request" value="Production request" />
                <Select
                  id="request"
                  className="mt-2"
                  value={selectedRequest ?? ""}
                  onChange={(e) =>
                    setSelectedRequest(
                      e.target.value ? Number(e.target.value) : null
                    )
                  }
                >
                  <option value="">Choose a request…</option>
                  {requests.map((r) => (
                    <option key={r.id} value={r.id}>
                      {r.title}
                    </option>
                  ))}
                </Select>

                {selectedRequestObj && (
                  <div className="mt-4 p-4 bg-blue-50 rounded-lg">
                    <p className="font-semibold">{selectedRequestObj.title}</p>
                    <p className="text-sm text-gray-700">
                      {selectedRequestObj.description}
                    </p>
                  </div>
                )}
              </Card>
            )}

            <div className="flex justify-end mt-8">
              <Button
                onClick={nextFromRequest}
                disabled={selectedRequest === null}
              >
                Next
              </Button>
            </div>
          </div>
        )}

        {/* STEP 2 — Categories */}
        {step === "categories" && (
          <>
            {loading && Object.keys(categories).length === 0 ? (
              <div className="flex justify-center py-12">
                <Spinner size="xl" />
              </div>
            ) : (
              <CategoryStep
                categories={availableCategories}
                selected={selectedCategories}
                onToggle={(c) =>
                  setSelectedCategories((prev) =>
                    prev.includes(c)
                      ? prev.filter((x) => x !== c)
                      : [...prev, c]
                  )
                }
                onBack={backToRequest}
                onNext={nextFromCategories}
              />
            )}
          </>
        )}

        {/* STEP 3 — Dashboard */}
        {step === "dashboard" && (
          <>
            <ProgressIndicator current={3} total={3} />

            <div className="mt-8 grid grid-cols-3 gap-4">
              <Card>
                <p>Total Files</p>
                <p className="text-3xl font-bold">{files.length}</p>
              </Card>
              <Card>
                <p>Categories</p>
                <p className="text-3xl font-bold">
                  {selectedCategories.length}
                </p>
              </Card>
              <Card>
                <p>Estimated Size</p>
                <p className="text-3xl font-bold">
                  ~{Math.round(files.length * 0.85)}
                </p>
              </Card>
            </div>

            <div className="mt-8">
              <Button color="gray" onClick={backToCategories}>
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
