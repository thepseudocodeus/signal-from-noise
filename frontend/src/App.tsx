// App.tsx
// =============================================================================
// Signal from Noise — Application Orchestrator
// =============================================================================
//
// This file is written in a *literate programming* style.
// It should read like a short research paper explaining:
//
//   • What the application does
//   • What state exists
//   • What events can happen
//   • How state changes in response to events
//   • Which parts are pure and which parts perform IO
//
// ----------------------------------------------------------------------------
// PURPOSE
// ----------------------------------------------------------------------------
// A deterministic 3-step workflow:
//
//   1) User selects a Production Request
//   2) User selects one or more Data Categories
//   3) Application loads files and renders a dashboard
//
// ----------------------------------------------------------------------------
// MENTAL MODEL (Haskell / Elm / State Machine)
// ----------------------------------------------------------------------------
//
//   The runtime owns an implicit `while True:` loop.
//   On each iteration:
//
//     Input   -> Msg
//     Process -> update(Model, Msg) => Model'
//     Output  -> view(Model')
//
//   React is responsible for:
//     • storing Model between iterations
//     • re-running the view after each update
//
//   We NEVER manually refresh the UI.
//
// ----------------------------------------------------------------------------
// INVARIANTS (must always hold)
// ----------------------------------------------------------------------------
//
//   • step === "categories"  => selectedRequestId !== null
//   • step === "dashboard"   => files is the current result set
//   • loading flags describe *what the app is waiting for*
//   • buttons are DISABLED when invalid, not hidden
//
// ----------------------------------------------------------------------------
// OWNERSHIP
// ----------------------------------------------------------------------------
//
//   • Backend (Wails / Go):
//       - GetProductionRequests
//       - GetCategories
//       - GetFiles
//
//   • App.tsx:
//       - Orchestrates workflow
//       - Owns state machine
//       - Performs effects
//
//   • Child components:
//       - Pure views
//       - Emit events via callbacks
//       - Never fetch data directly
//
// ----------------------------------------------------------------------------
// DEBUGGING STRATEGY
// ----------------------------------------------------------------------------
//
//   • Log every Msg
//   • Log state transitions (summary only)
//   • Log backend request/response boundaries
//   • If UI is wrong, state is wrong — inspect Model
//
// ----------------------------------------------------------------------------
// TODO / OPEN QUESTIONS
// ----------------------------------------------------------------------------
// [ ] Confirm first-screen button visibility logic
// [ ] Validate category mapping assumptions
// [ ] Tighten FileItem shape once backend is stable
//
// App.tsx (literate style)
// Goal: a deterministic 3-step flow:
//   1) choose a production request
//   2) choose categories
//   3) load files and render a dashboard
//
// Mental model (React):
//   - State is our "data structure" (the current truth).
//   - Handlers are our "process" (they update state, and may call the backend).
//   - Render is our "output" (a pure projection of state).
//
// Crucially:
//   - We do NOT manually "refresh the UI".
//   - We update state, and React re-renders automatically.
//
// [ ] TODO: walk through flow to confirm logic
// [ ] TODO: fix button loading bug
// START APP
// (runtime loop)
//   ├─ user clicks / async completes
//   ├─ React calls your handler with an event
//   ├─ you compute next_state = transition(state, event)
//   ├─ React stores next_state
//   ├─ React calls render(next_state)
//   └─ wait for next event
//   └─ UNLESS stop
// END APP
//
// =============================================================================

// ----------------------------------------------------------------------------
// Imports: framework
// ----------------------------------------------------------------------------
import { useEffect, useReducer, useRef } from "react";

// ----------------------------------------------------------------------------
// Imports: UI components (pure views)
// ----------------------------------------------------------------------------
import { Button, Card, Spinner } from "flowbite-react";
import { ProductionRequestSelectionPage } from "./components/pages/ProductionRequestSelectionPage";
import { CategoryStep } from "./components/steps/CategoryStep";

// ----------------------------------------------------------------------------
// Imports: domain types
// ----------------------------------------------------------------------------
import { DataCategory } from "./types/query";

// ----------------------------------------------------------------------------
// DOMAIN TYPES
// ----------------------------------------------------------------------------

export type Step = "request" | "categories" | "dashboard";

export type FileItem = {
  id?: string;
  filename?: string;
  path?: string;
  category?: string;
  // NOTE: backend may send additional fields; tighten later
  [k: string]: unknown;
};

// ----------------------------------------------------------------------------
// MODEL (Application State)
// ----------------------------------------------------------------------------

export type Model = {
  step: Step;

  // Selections
  selectedRequestId: number | null;
  selectedCategories: DataCategory[];

  // Data
  files: FileItem[];

  // Loading flags describe *what we are waiting for*
  loading: {
    requests: boolean;
    categories: boolean;
    files: boolean;
  };

  // Error boundary (user-visible)
  error?: string;
};

// ----------------------------------------------------------------------------
// MESSAGES (Events)
// ----------------------------------------------------------------------------

export type Msg =
  | { type: "AppStarted" }
  | { type: "RequestSelected"; id: number }
  | { type: "CategoriesLoading" }
  | { type: "CategoriesLoaded"; categories: DataCategory[] }
  | { type: "CategoryToggled"; category: DataCategory }
  | { type: "FilesLoading" }
  | { type: "FilesLoaded"; files: FileItem[] }
  | { type: "BackToRequest" }
  | { type: "BackToCategories" }
  | { type: "Error"; message: string };

// ----------------------------------------------------------------------------
// LOGGING (Boundary-aware observability)
// ----------------------------------------------------------------------------
// We need transparency to make apps intelligent.
// This utility will:
// - log boundaries between phases
// - log inputs to backend
// - log shape of outputs from backend

const DEBUG = true;

function logPhase(phase: string, payload?: unknown) {
  if (!DEBUG) return;
  console.groupCollapsed(`[App] ${phase}`);
  if (payload !== undefined) console.log(payload);
  console.groupEnd();
}

function assert(condition: unknown, message: string): asserts condition {
  if (!condition) {
    throw new Error(`[ASSERTION FAILED] ${message}`);
  }
}

function assertInvariants(model: Model) {
  // Step invariants
  if (model.step === "categories") {
    assert(
      model.selectedRequestId !== null,
      "categories step requires selectedRequestId"
    );
  }

  if (model.step === "dashboard") {
    assert(Array.isArray(model.files), "dashboard step requires files array");
  }

  // Loading invariants
  if (model.loading.files) {
    assert(model.step === "dashboard", "loading.files implies dashboard step");
  }
}

// ----------------------------------------------------------------------------
// STATE
// ----------------------------------------------------------------------------

/*
We will borrow from Haskell to build a single source of state transition true = the state machine.

Flow:
- prepare the model with initial state (initModel)
- create a reducer function that will transition the model based on the message


*/

export const initModel: Model = {
  step: "request",
  selectedRequestId: null,
  selectedCategories: [],
  files: [],
  loading: {
    requests: true, // we load requests on startup
    categories: false,
    files: false,
  },
  error: undefined,
};

export function update(model: Model, msg: Msg): Model {
  switch (msg.type) {
    /* ---------------------------------------------------------------------- */
    /* App lifecycle                                                           */
    /* ---------------------------------------------------------------------- */

    case "AppStarted":
      // Startup is when we begin loading requests.
      return {
        ...model,
        step: "request",
        loading: { ...model.loading, requests: true },
        error: undefined,
      };

    /* ---------------------------------------------------------------------- */
    /* Request selection                                                       */
    /* ---------------------------------------------------------------------- */

    case "RequestSelected":
      // Selecting a request moves us toward categories (after categories load).
      return {
        ...model,
        selectedRequestId: msg.id,
        selectedCategories: [],
        files: [],
        step: "categories",
        loading: { ...model.loading, categories: true, files: false },
        error: undefined,
      };

    /* ---------------------------------------------------------------------- */
    /* Categories                                                              */
    /* ---------------------------------------------------------------------- */

    case "CategoriesLoading":
      return {
        ...model,
        step: "categories",
        loading: { ...model.loading, categories: true },
        error: undefined,
      };

    case "CategoriesLoaded":
      // Categories are loaded; user can now choose them.
      return {
        ...model,
        step: "categories",
        loading: { ...model.loading, categories: false },
        error: undefined,
        // Note: categories list itself can be stored elsewhere if needed.
        // If you keep categories in Model, add it to Model and set it here.
      };

    case "CategoryToggled": {
      // Toggle in a purely functional way.
      const nextSelected = model.selectedCategories.includes(msg.category)
        ? model.selectedCategories.filter((c) => c !== msg.category)
        : [...model.selectedCategories, msg.category];

      return {
        ...model,
        selectedCategories: nextSelected,
        error: undefined,
      };
    }
  }
}

function updateWithAssertions(model: Model, msg: Msg): Model {
  const next = update(model, msg);
  if (DEBUG) assertInvariants(next);
  return next;
}

/* -------------------------------------------------------------------------- */
/* 3) The App: state + processes + output                                      */
/* -------------------------------------------------------------------------- */

/*

App story

@load browser loads main.tsx and then launches this App component
- compoents are imported to be rendered
- data is defined in strutures
- helper functions are defined to keep single responsibility

*/

export default function App() {
  // [ ] TODO: confirm tracking added to every component to track flow
  console.log("FRONTEND: App render");
  /* ------------------------------------------------------------------------ */
  /* State: the "data structure"                                               */
  /* ------------------------------------------------------------------------ */

  const [step, setStep] = useState<Step>("request");

  const [requests, setRequests] = useState<ProductionRequest[]>([]);
  const [selectedRequest, setSelectedRequest] =
    useState<ProductionRequest | null>(null);

  const [categories, setCategories] = useState<DataCategory[]>([
    "Email",
    "Claims",
    "Other",
  ]);
  const [selectedCategories, setSelectedCategories] = useState<DataCategory[]>(
    []
  );

  const [files, setFiles] = useState<FileItem[]>([]);

  const [loadingRequests, setLoadingRequests] = useState(true);
  const [loadingCategories, setLoadingCategories] = useState(false);
  const [loadingFiles, setLoadingFiles] = useState(false);

  // A ref to prevent stale async responses from overwriting newer state.
  // This is a common cause of “React isn’t working as expected”.
  const latestFilesReq = useRef(0);

  /* ------------------------------------------------------------------------ */
  /* Side effect: load initial requests once on mount                          */
  /* ------------------------------------------------------------------------ */

  useEffect(() => {
    // This runs once after the first render (like "on start").
    // It is "process" that produces state.
    const loadRequests = async () => {
      logPhase("EFFECT: loadRequests (start)");
      setLoadingRequests(true);

      try {
        const { GetProductionRequests } = await import(
          "../wailsjs/go/main/App"
        );
        const res = await GetProductionRequests();

        logPhase("EFFECT: loadRequests (result)", {
          isArray: Array.isArray(res),
          length: Array.isArray(res) ? res.length : "n/a",
          first: Array.isArray(res) ? res[0] : res,
        });

        setRequests(res || []);
      } catch (err) {
        logPhase("EFFECT: loadRequests (error)", err);
        setRequests([]);
      } finally {
        setLoadingRequests(false);
        logPhase("EFFECT: loadRequests (done)");
      }
    };

    loadRequests();
  }, []);

  /* ------------------------------------------------------------------------ */
  /* Process: load categories (called when moving request -> categories)       */
  /* ------------------------------------------------------------------------ */

  const loadCategories = async () => {
    logPhase("PROCESS: loadCategories (start)");
    setLoadingCategories(true);

    try {
      const { GetCategories } = await import("../wailsjs/go/main/App");
      const res = await GetCategories();

      // Support either ["email","claim"] or {email:..., claim:...}
      const rawCats: string[] = Array.isArray(res)
        ? res
        : res && typeof res === "object"
        ? Object.keys(res as Record<string, unknown>)
        : [];

      const mapped: DataCategory[] = rawCats.map((cat) => {
        const c = String(cat).toLowerCase();
        if (c === "email") return "Email";
        if (c === "claim" || c === "claims") return "Claims";
        return "Other";
      });

      logPhase("PROCESS: loadCategories (result)", {
        rawCats,
        mapped,
      });

      setCategories(mapped.length ? mapped : ["Email", "Claims", "Other"]);
    } catch (err) {
      logPhase("PROCESS: loadCategories (error)", err);
      setCategories(["Email", "Claims", "Other"]);
    } finally {
      setLoadingCategories(false);
      logPhase("PROCESS: loadCategories (done)");
    }
  };

  /* ------------------------------------------------------------------------ */
  /* Process: start a request (input -> process -> state update -> output)     */
  /* ------------------------------------------------------------------------ */

  const handleStartRequest = async (requestNumber: number) => {
    logPhase("INPUT: user started request", { requestNumber });

    const chosen =
      requests.find((r) => r.id === requestNumber) ||
      ({
        id: requestNumber,
        title: `Production Request ${requestNumber}`,
        description: "",
      } as ProductionRequest);

    // Process = state updates
    setSelectedRequest(chosen);

    // More process: fetch categories
    await loadCategories();

    // Output will follow automatically via render, because state changed
    setStep("categories");

    logPhase("STATE: moved to categories", { selectedRequest: chosen });
  };

  /* ------------------------------------------------------------------------ */
  /* Process: toggle category selection                                        */
  /* ------------------------------------------------------------------------ */

  const handleCategoryToggle = (category: DataCategory) => {
    logPhase("INPUT: toggle category", { category });

    setSelectedCategories((prev) => {
      const next = prev.includes(category)
        ? prev.filter((c) => c !== category)
        : [...prev, category];

      logPhase("STATE: selectedCategories updated", { prev, next });
      return next;
    });
  };

  /* ------------------------------------------------------------------------ */
  /* Process: load files (categories -> dashboard)                             */
  /* ------------------------------------------------------------------------ */

  const handleNextFromCategories = async () => {
    if (!selectedCategories.length) return;

    latestFilesReq.current += 1;
    const myReq = latestFilesReq.current;

    logPhase("INPUT: user requested files for categories", {
      selectedCategories,
      myReq,
    });

    setLoadingFiles(true);

    try {
      const { GetFiles } = await import("../wailsjs/go/main/App");

      const catStrings = selectedCategories.map((c) =>
        c === "Email" ? "email" : c === "Claims" ? "claim" : "other"
      );

      logPhase("PROCESS: GetFiles called with", { catStrings, myReq });

      const res = await GetFiles(catStrings);

      logPhase("PROCESS: GetFiles returned", {
        isArray: Array.isArray(res),
        length: Array.isArray(res) ? res.length : "n/a",
        first: Array.isArray(res) ? res[0] : res,
        myReq,
      });

      // Ignore stale response
      if (myReq !== latestFilesReq.current) {
        logPhase("PROCESS: stale GetFiles response ignored", { myReq });
        return;
      }

      // State update: files now become the truth for the dashboard render
      setFiles((Array.isArray(res) ? res : []) as FileItem[]);
      setStep("dashboard");

      logPhase("STATE: moved to dashboard", {
        filesCount: Array.isArray(res) ? res.length : 0,
      });
    } catch (err) {
      logPhase("PROCESS: GetFiles error", err);
      alert("Error loading files: " + String(err));
    } finally {
      if (myReq === latestFilesReq.current) {
        setLoadingFiles(false);
      }
      logPhase("PROCESS: load files (done)", { myReq });
    }
  };

  /* ------------------------------------------------------------------------ */
  /* Output: render (a pure projection of current state)                       */
  /* ------------------------------------------------------------------------ */

  // Helpful render-phase log: confirms which step is actually rendering.
  // (Optional; remove when stable.)
  logPhase("RENDER", {
    step,
    requests: requests.length,
    selectedRequest: selectedRequest?.id ?? null,
    selectedCategories,
    files: files.length,
    loadingRequests,
    loadingCategories,
    loadingFiles,
  });

  return (
    <div className="min-h-screen bg-gray-50 p-8">
      <div className="max-w-4xl mx-auto">
        <h1 className="text-4xl font-bold text-gray-900 mb-8">
          Signal from Noise
        </h1>

        {step === "request" && (
          <>
            {loadingRequests ? (
              <div className="flex justify-center items-center py-12 space-x-4">
                <Spinner size="xl" />
                <span className="text-gray-600">Loading requests...</span>
              </div>
            ) : (
              <ProductionRequestSelectionPage onStart={handleStartRequest} />
            )}
          </>
        )}

        {step === "categories" && (
          <CategoryStep
            categories={categories}
            selected={selectedCategories}
            onToggle={handleCategoryToggle}
            onNext={handleNextFromCategories}
            onBack={() => setStep("request")}
            minSelection={1}
            loading={loadingCategories as any} // remove if CategoryStep doesn't accept it
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
                        {files.slice(0, 50).map((file) => (
                          <tr
                            key={
                              file.path ??
                              file.id ??
                              `${file.category ?? "unknown"}:${
                                file.filename ?? "unknown"
                              }`
                            }
                            className="border-b hover:bg-gray-50"
                          >
                            <td className="p-2 font-medium">
                              {file.filename ?? "(missing filename)"}
                            </td>
                            <td className="p-2">
                              <span className="px-2 py-1 bg-blue-100 text-blue-800 rounded text-sm capitalize">
                                {file.category ?? "unknown"}
                              </span>
                            </td>
                            <td
                              className="p-2 text-sm text-gray-600 truncate max-w-md"
                              title={file.path ?? ""}
                            >
                              {file.path ?? "(missing path)"}
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
