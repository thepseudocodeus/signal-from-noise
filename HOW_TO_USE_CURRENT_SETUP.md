# How to Use Current Setup for Email Query Demo
## Connecting to Data Lake on External Drive

**Data Lake Location:** `/Volumes/DirectFileStore/DELEMAR/`
**Email Data Path:** `/Volumes/DirectFileStore/DELEMAR/unprocessed/emails_final/`

---

## Current Setup Status

✅ **Completed:**
- Wails + Go + React + TypeScript setup
- Tailwind CSS + Flowbite React configured
- TypeScript conversion complete (100% type coverage)
- Build system working

⏳ **To Build:**
- Data lake connection service
- Email query processing
- Typeform-inspired step-by-step UI
- Frontend-backend integration

---

## Architecture Overview

### Current Structure
```
signal-from-noise/
├── .env                          # Data lake paths (configured)
├── app.go                        # App struct (basic)
├── main.go                       # Wails entry point
├── frontend/
│   └── src/
│       ├── App.tsx              # Current: Simple greet demo
│       └── main.tsx             # Entry point
```

### Target Structure (For Demo)
```
signal-from-noise/
├── config/
│   └── config.go                # Read .env, manage config
├── datalake/
│   └── datalake.go              # Data lake connection & discovery
├── queries/
│   └── email_query.go           # Email query processing
├── app.go                       # Enhanced with query methods
└── frontend/
    └── src/
        ├── App.tsx              # Step-by-step flow
        ├── components/
        │   ├── steps/           # Step components
        │   └── shared/          # Shared components
        ├── types/
        │   └── query.ts         # TypeScript types
        └── hooks/
            └── useQuery.ts      # Query state management
```

---

## Data Lake Structure (From .env)

```
/Volumes/DirectFileStore/DELEMAR/
├── unprocessed/
│   ├── emails_final/            # ← Processed email Parquet files (USE THIS)
│   ├── emails_in_process/
│   └── Isaiah.Delemar@sol.doi.gov-olderthan1year.pst
└── [other data lake directories]
```

**For Email Demo:** Focus on `/Volumes/DirectFileStore/DELEMAR/unprocessed/emails_final/`

---

## Implementation Roadmap

### Step 1: Backend Foundation
**Goal:** Connect to data lake and discover email files

**Files to Create:**
1. `config/config.go` - Read .env, provide configuration
2. `datalake/datalake.go` - Data lake connection service
3. Enhance `app.go` - Add data lake methods

**Key Functions:**
```go
// config/config.go
func LoadConfig() (*Config, error)
func GetDataLakePath() string
func GetEmailFinalPath() string

// datalake/datalake.go
func DiscoverEmailFiles(path string) ([]string, error)
func ValidateDataLake(path string) error

// app.go
func (a *App) GetDataLakeStatus() (string, error)
func (a *App) ListEmailFiles() ([]string, error)
```

### Step 2: Email Query Processing
**Goal:** Process email Parquet files based on query parameters

**Files to Create:**
1. `queries/email_query.go` - Email query engine
2. `models/query.go` - Query data models

**Key Functions:**
```go
// queries/email_query.go
func ExecuteEmailQuery(query EmailQuery) (*QueryResult, error)
func FilterByDateRange(files []string, startDate, endDate time.Time) ([]EmailRecord, error)
func FilterByKeywords(records []EmailRecord, keywords []string) ([]EmailRecord, error)

// models/query.go
type EmailQuery struct {
    DateRange DateRange
    Keywords  []string
    QueryType string
}
```

### Step 3: Frontend Structure
**Goal:** Create Typeform-inspired step-by-step UI

**Files to Create:**
1. `frontend/src/types/query.ts` - TypeScript types
2. `frontend/src/components/steps/` - Step components
3. `frontend/src/components/shared/ProgressIndicator.tsx`
4. `frontend/src/hooks/useQuery.ts` - Query state
5. Update `frontend/src/App.tsx` - Multi-step flow

**Step Components:**
1. `WelcomeStep.tsx` - Welcome screen
2. `QueryTypeStep.tsx` - Select query type (Email Timeline, Keyword Search, etc.)
3. `DateRangeStep.tsx` - Select date range (Flowbite DatePicker)
4. `KeywordsStep.tsx` - Enter keywords (optional)
5. `ReviewStep.tsx` - Review selections
6. `ProcessingStep.tsx` - Show progress
7. `ResultsStep.tsx` - Display results

### Step 4: Integration
**Goal:** Connect frontend to backend

**Tasks:**
1. Add Wails methods to `app.go`
2. Generate TypeScript bindings (`wails generate module`)
3. Create service layer in frontend
4. Connect step components to backend
5. Test end-to-end flow

---

## Quick Start Implementation Plan

### Phase 1: Make It Work (Demo-Ready)

**1. Backend - Data Lake Connection**
```go
// config/config.go
package config

import (
    "os"
    "path/filepath"
)

type Config struct {
    DataLakePath    string
    EmailFinalPath  string
}

func LoadConfig() (*Config, error) {
    dataLake := os.Getenv("DATA_LAKE")
    if dataLake == "" {
        return nil, fmt.Errorf("DATA_LAKE not set")
    }

    return &Config{
        DataLakePath:   dataLake,
        EmailFinalPath: filepath.Join(dataLake, "unprocessed/emails_final"),
    }, nil
}
```

**2. Backend - Email Query (Simplified for Demo)**
```go
// For demo: Start with CSV/JSON files if Parquet is complex
// Or use basic Parquet library
```

**3. Frontend - Step-by-Step UI**
```tsx
// App.tsx - Multi-step flow
const [currentStep, setCurrentStep] = useState(0);
const [queryData, setQueryData] = useState<QueryData>({});

const steps = [
  WelcomeStep,
  QueryTypeStep,
  DateRangeStep,
  KeywordsStep,
  ReviewStep,
  ProcessingStep,
  ResultsStep
];
```

---

## Key Decisions for Demo

### 1. Parquet vs CSV/JSON
**Question:** Are the email files in Parquet format or CSV/JSON?

**For Demo:**
- **If Parquet:** Use Go Parquet library (e.g., `github.com/apache/arrow/go/v14/parquet`)
- **If CSV/JSON:** Use standard Go file reading (simpler for demo)

**Recommendation:** Start with CSV/JSON if available, add Parquet support if needed.

### 2. Query Scope
**For Demo:** Start simple
- Query type: "Email Timeline" (all emails in date range)
- Filters: Date range + optional keywords
- Output: Display first 100 results (or summary)

### 3. Performance
**For Demo:**
- Process sample files (first N files)
- Return first 100-1000 records
- Full processing can come later

---

## Testing Strategy

### 1. Data Lake Connection
```bash
# Test that external drive is accessible
ls /Volumes/DirectFileStore/DELEMAR/unprocessed/emails_final/
```

### 2. Backend Testing
```go
// Test data lake connection
config, _ := config.LoadConfig()
files, _ := datalake.DiscoverEmailFiles(config.EmailFinalPath)
fmt.Println("Found", len(files), "email files")
```

### 3. Frontend Testing
```bash
# Run dev server
cd frontend && npm run dev

# Or run Wails dev
wails dev
```

### 4. End-to-End Testing
- Select query type
- Set date range
- Add keywords (optional)
- Execute query
- Verify results

---

## Next Immediate Steps

1. **Verify Data Lake Structure**
   ```bash
   ls -la /Volumes/DirectFileStore/DELEMAR/unprocessed/emails_final/
   ```
   - Check file format (Parquet, CSV, JSON, etc.)
   - Check file count and structure

2. **Create Config Package**
   - Read .env file
   - Provide data lake paths

3. **Create Data Lake Service**
   - Connect to external drive
   - Discover email files
   - Test connection

4. **Build First Step Component**
   - Welcome/Query Type selection
   - Test Flowbite components

5. **Connect Frontend to Backend**
   - Add Wails methods
   - Test data flow

---

## Success Criteria for Demo

**Demo is ready when:**
- ✅ Application reads data lake path from .env
- ✅ Can discover email files on external drive
- ✅ User can select query parameters via Typeform-inspired UI
- ✅ Application processes email data (even if simplified)
- ✅ Results are displayed in the app
- ✅ User experience is smooth and engaging

---

## Questions to Clarify

1. **File Format:** What format are the email files? (Parquet, CSV, JSON, etc.)
2. **File Structure:** What's the structure of email files? (columns, schema)
3. **Demo Scope:** How many files should we process for demo? (sample vs all)
4. **Output Format:** What format for results? (CSV display, JSON, download?)
5. **Query Types:** Which query types to demo? (Timeline, Keyword Search, Summary?)

---

**Ready to start building!** Let's begin with Step 1: Config and Data Lake Connection.
