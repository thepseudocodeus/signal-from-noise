# Desktop Application Implementation Guide

**Version:** 1.0  
**Date:** 2026-01-XX  
**For:** Cursor Agent Building Wails + React Application

---

## 1. Project Setup

### 1.1 Initialize Wails Project

```bash
# Install Wails CLI
go install github.com/wailsapp/wails/v2/cmd/wails@latest

# Create new project
wails init -n email-query-app -t react-ts

# Install dependencies
cd email-query-app
npm install
```

### 1.2 Install UI Dependencies

```bash
# Tailwind CSS
npm install -D tailwindcss postcss autoprefixer
npx tailwindcss init -p

# Flowbite React
npm install flowbite flowbite-react

# Additional utilities
npm install date-fns  # For date handling
npm install react-icons  # For icons
```

### 1.3 Project Structure

```
email-query-app/
├── app/                    # Go backend
│   ├── main.go
│   ├── queries/           # Query engine
│   ├── data/              # Data access layer
│   └── models/            # Data models
├── frontend/              # React frontend
│   ├── src/
│   │   ├── components/    # React components
│   │   │   ├── steps/    # Step components
│   │   │   └── shared/   # Shared components
│   │   ├── hooks/        # Custom hooks
│   │   ├── types/        # TypeScript types
│   │   └── utils/        # Utilities
│   └── tailwind.config.js
└── wails.json
```

---

## 2. Data Models

### 2.1 Information Request Model

```typescript
// frontend/src/types/information-request.ts

export interface InformationRequest {
  id: string;              // e.g., "IR-001"
  title: string;
  description: string;
  dataRequirements: string[];
  outputFormat: string;
  estimatedTime: number;   // seconds
}

export const INFORMATION_REQUESTS: InformationRequest[] = [
  {
    id: "IR-001",
    title: "Email Communication Timeline",
    description: "All emails within date range, ordered chronologically",
    dataRequirements: ["email_parquet"],
    outputFormat: "CSV with timestamp, from, to, subject",
    estimatedTime: 15
  },
  // ... 19 more requests
];
```

### 2.2 User Selection Model

```typescript
// frontend/src/types/query.ts

export interface UserQuery {
  informationRequest: InformationRequest;
  dateRange: {
    start: Date;
    end: Date;
  };
  categories: Category[];
  keywords: string[];
  filters: Record<string, any>;
}

export type Category = "Claims" | "Emails" | "Other";
```

### 2.3 Go Backend Models

```go
// app/models/query.go

type InformationRequest struct {
    ID               string   `json:"id"`
    Title            string   `json:"title"`
    Description      string   `json:"description"`
    DataRequirements []string `json:"data_requirements"`
    OutputFormat     string   `json:"output_format"`
}

type UserQuery struct {
    InformationRequest InformationRequest `json:"information_request"`
    DateRange          DateRange          `json:"date_range"`
    Categories         []string           `json:"categories"`
    Keywords           []string           `json:"keywords"`
    Filters            map[string]any     `json:"filters"`
}

type DateRange struct {
    Start time.Time `json:"start"`
    End   time.Time `json:"end"`
}
```

---

## 3. Component Implementation

### 3.1 Step 1: Information Request Selector

```tsx
// frontend/src/components/steps/InformationRequestStep.tsx

import { Card } from 'flowbite-react';
import { InformationRequest } from '@/types/information-request';

interface Props {
  requests: InformationRequest[];
  selected: InformationRequest | null;
  onSelect: (request: InformationRequest) => void;
  progress: number;
}

export function InformationRequestStep({ 
  requests, 
  selected, 
  onSelect, 
  progress 
}: Props) {
  return (
    <div className="space-y-6">
      <ProgressIndicator current={1} total={5} />
      
      <h2 className="text-2xl font-bold">
        What information do you need?
      </h2>
      
      <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
        {requests.map((request) => (
          <Card
            key={request.id}
            className={`cursor-pointer transition-all ${
              selected?.id === request.id
                ? 'ring-2 ring-blue-500'
                : 'hover:shadow-lg'
            }`}
            onClick={() => onSelect(request)}
          >
            <h3 className="text-lg font-semibold">{request.title}</h3>
            <p className="text-sm text-gray-600">{request.description}</p>
          </Card>
        ))}
      </div>
      
      <button
        disabled={!selected}
        className="btn-primary"
        onClick={() => {/* Navigate to next step */}}
      >
        Next
      </button>
    </div>
  );
}
```

### 3.2 Step 2: Date Range Selector

```tsx
// frontend/src/components/steps/DateRangeStep.tsx

import { Datepicker } from 'flowbite-react';
import { useState } from 'react';

interface Props {
  minDate: Date;
  maxDate: Date;
  startDate: Date | null;
  endDate: Date | null;
  onStartChange: (date: Date) => void;
  onEndChange: (date: Date) => void;
  progress: number;
}

export function DateRangeStep({
  minDate,
  maxDate,
  startDate,
  endDate,
  onStartChange,
  onEndChange,
  progress
}: Props) {
  const [error, setError] = useState<string | null>(null);
  
  const validateRange = (start: Date, end: Date) => {
    if (end < start) {
      setError("End date must be after start date");
      return false;
    }
    setError(null);
    return true;
  };
  
  return (
    <div className="space-y-6">
      <ProgressIndicator current={2} total={5} />
      
      <h2 className="text-2xl font-bold">
        Select date range
      </h2>
      
      <div className="grid grid-cols-2 gap-4">
        <div>
          <label>Start Date</label>
          <Datepicker
            minDate={minDate}
            maxDate={maxDate}
            value={startDate?.toLocaleDateString()}
            onSelectedDateChanged={(date) => {
              onStartChange(date);
              if (endDate) validateRange(date, endDate);
            }}
          />
        </div>
        
        <div>
          <label>End Date</label>
          <Datepicker
            minDate={startDate || minDate}
            maxDate={maxDate}
            value={endDate?.toLocaleDateString()}
            onSelectedDateChanged={(date) => {
              onEndChange(date);
              if (startDate) validateRange(startDate, date);
            }}
          />
        </div>
      </div>
      
      {error && <div className="text-red-500">{error}</div>}
      
      <button
        disabled={!startDate || !endDate || !!error}
        className="btn-primary"
      >
        Next
      </button>
    </div>
  );
}
```

### 3.3 Step 3: Category Selector

```tsx
// frontend/src/components/steps/CategoryStep.tsx

import { Checkbox } from 'flowbite-react';
import { Category } from '@/types/query';

interface Props {
  categories: Category[];
  selected: Category[];
  onToggle: (category: Category) => void;
  progress: number;
  minSelection: number;
}

export function CategoryStep({
  categories,
  selected,
  onToggle,
  progress,
  minSelection
}: Props) {
  const isValid = selected.length >= minSelection;
  
  return (
    <div className="space-y-6">
      <ProgressIndicator current={3} total={5} />
      
      <h2 className="text-2xl font-bold">
        Select categories
      </h2>
      
      <div className="space-y-4">
        {categories.map((category) => (
          <div key={category} className="flex items-center">
            <Checkbox
              id={category}
              checked={selected.includes(category)}
              onChange={() => onToggle(category)}
            />
            <label htmlFor={category} className="ml-2">
              {category}
            </label>
          </div>
        ))}
      </div>
      
      {!isValid && (
        <div className="text-yellow-500">
          Please select at least {minSelection} category
        </div>
      )}
      
      <button
        disabled={!isValid}
        className="btn-primary"
      >
        Next
      </button>
    </div>
  );
}
```

---

## 4. Backend Implementation

### 4.1 Query Engine

```go
// app/queries/engine.go

package queries

import (
    "context"
    "time"
)

type QueryEngine struct {
    dataPath string
}

func NewQueryEngine(dataPath string) *QueryEngine {
    return &QueryEngine{dataPath: dataPath}
}

func (e *QueryEngine) ExecuteQuery(ctx context.Context, query UserQuery) (*QueryResult, error) {
    // 1. Validate query
    if err := e.validateQuery(query); err != nil {
        return nil, err
    }
    
    // 2. Determine data sources
    sources := e.determineSources(query)
    
    // 3. Process each source
    results := make([]Record, 0)
    for _, source := range sources {
        records, err := e.processSource(ctx, source, query)
        if err != nil {
            return nil, err
        }
        results = append(results, records...)
    }
    
    // 4. Aggregate and order
    aggregated := e.aggregate(results, query)
    
    // 5. Generate output
    return e.generateOutput(aggregated, query)
}

func (e *QueryEngine) validateQuery(query UserQuery) error {
    // Guarantee G2: Date range validation
    if query.DateRange.End.Before(query.DateRange.Start) {
        return fmt.Errorf("end date must be after start date")
    }
    
    // Guarantee G3: At least one category
    if len(query.Categories) == 0 {
        return fmt.Errorf("at least one category must be selected")
    }
    
    return nil
}
```

### 4.2 Parquet Reader

```go
// app/data/parquet_reader.go

package data

import (
    "github.com/apache/arrow/go/v14/parquet"
    "github.com/apache/arrow/go/v14/parquet/file"
)

type ParquetReader struct {
    filePath string
}

func (r *ParquetReader) ReadWithFilters(
    ctx context.Context,
    filters []Filter,
    dateRange DateRange,
) ([]Record, error) {
    // Open Parquet file
    reader, err := file.OpenParquetFile(r.filePath, false)
    if err != nil {
        return nil, err
    }
    defer reader.Close()
    
    // Apply predicate pushdown (date range)
    // Process in chunks
    // Apply filters
    // Return records
}
```

---

## 5. Guarantee System Integration

### 5.1 Frontend Guarantees

```typescript
// frontend/src/hooks/useGuarantees.ts

import { Guarantee } from '@/types/guarantees';

export function useStepGuarantees(step: number, data: any) {
  const guarantees: Guarantee[] = [];
  
  if (step === 1) {
    guarantees.push({
      id: 'G1_request_selected',
      condition: () => data.selectedRequest !== null,
      expect: true,
      description: 'Information request is selected'
    });
  }
  
  if (step === 2) {
    guarantees.push({
      id: 'G2_date_range_valid',
      condition: () => {
        const { start, end } = data.dateRange;
        return start <= end && start >= data.minDate && end <= data.maxDate;
      },
      expect: true,
      description: 'Date range is valid'
    });
  }
  
  // Check guarantees
  const violations = guarantees.filter(g => !g.condition());
  
  return { guarantees, violations, allPass: violations.length === 0 };
}
```

### 5.2 Backend Guarantees

```go
// app/queries/guarantees.go

package queries

type Guarantee struct {
    ID          string
    Description string
    Condition   func() bool
    Expect      bool
}

func (e *QueryEngine) validateGuarantees(query UserQuery) []GuaranteeViolation {
    guarantees := []Guarantee{
        {
            ID:          "G5_output_relevant",
            Description: "Output contains relevant data",
            Condition:   func() bool { return e.checkRelevance(query) },
            Expect:      true,
        },
        // ... more guarantees
    }
    
    violations := []GuaranteeViolation{}
    for _, g := range guarantees {
        if g.Condition() != g.Expect {
            violations = append(violations, GuaranteeViolation{
                Guarantee: g,
                Message:   fmt.Sprintf("Guarantee %s violated", g.ID),
            })
        }
    }
    
    return violations
}
```

---

## 6. Testing Strategy

### 6.1 Component Tests

```typescript
// frontend/src/components/steps/__tests__/InformationRequestStep.test.tsx

import { render, screen } from '@testing-library/react';
import { InformationRequestStep } from '../InformationRequestStep';

describe('InformationRequestStep', () => {
  it('validates guarantee G1: request must be selected', () => {
    const requests = [/* ... */];
    const { getByText } = render(
      <InformationRequestStep
        requests={requests}
        selected={null}
        onSelect={jest.fn()}
        progress={1/5}
      />
    );
    
    const nextButton = getByText('Next');
    expect(nextButton).toBeDisabled(); // Guarantee: can't proceed without selection
  });
});
```

### 6.2 Integration Tests

```go
// app/queries/engine_test.go

func TestQueryEngine_ValidateGuarantees(t *testing.T) {
    engine := NewQueryEngine("/test/data")
    query := UserQuery{
        DateRange: DateRange{
            Start: time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
            End:   time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC), // Invalid!
        },
    }
    
    violations := engine.validateGuarantees(query)
    
    // Guarantee G2 should be violated
    assert.Contains(t, violations, "G2_date_range_valid")
}
```

---

## 7. Expert Teaching Moments

### 7.1 Assumption Falsification

**When building, ask:**
- "What assumption am I making here?"
- "How can I test this assumption?"
- "What would prove this assumption wrong?"

**Example:**
- **Assumption:** Users want to see all 20 requests at once
- **Test:** A/B test with pagination vs. all-at-once
- **Falsification:** If users prefer pagination, assumption is false

### 7.2 Guarantee-Driven Development

**For each feature:**
1. **Define guarantee:** What must be true?
2. **Encode as assertion:** How do we check it?
3. **Test guarantee:** Does it hold?
4. **Learn from failures:** If it fails, what did we learn?

**Example:**
- **Feature:** Date range selection
- **Guarantee:** End date >= start date
- **Assertion:** `endDate >= startDate`
- **Test:** Try invalid range, guarantee should fail
- **Learning:** If guarantee fails, UX needs improvement

### 7.3 Systematic Approach

**Work in phases:**
1. **Research:** Understand the problem
2. **Specify:** Define what we're building
3. **Implement:** Build it
4. **Validate:** Test assumptions
5. **Refine:** Improve based on learning

**Don't skip phases:**
- Don't implement without specification
- Don't validate without implementation
- Don't refine without validation

---

## 8. Data Flow Architecture

### 8.1 Frontend → Backend Communication

```typescript
// frontend/src/services/queryService.ts

import { invoke } from '@wailsapp/runtime';

export async function executeQuery(query: UserQuery): Promise<QueryResult> {
  // Wails automatically handles Go ↔ TypeScript communication
  const result = await invoke('ExecuteQuery', query);
  return result as QueryResult;
}
```

```go
// app/main.go

func (a *App) ExecuteQuery(query models.UserQuery) (*models.QueryResult, error) {
    engine := queries.NewQueryEngine(a.config.DataPath)
    return engine.ExecuteQuery(context.Background(), query)
}
```

### 8.2 Data Processing Flow

```
User Selections (Frontend)
    ↓
Wails Bridge
    ↓
Go Backend (Query Engine)
    ↓
Parquet Files (Data Lake)
    ↓
Filtered/Aggregated Data
    ↓
Zip File Generation
    ↓
Return to Frontend
    ↓
Download to User
```

---

## 9. Performance Considerations

### 9.1 Frontend Performance

- **Lazy loading:** Load components on demand
- **Memoization:** Cache expensive computations
- **Virtual scrolling:** For long lists
- **Debouncing:** For search inputs

### 9.2 Backend Performance

- **Streaming:** Process data in chunks
- **Predicate pushdown:** Filter at Parquet level
- **Parallel processing:** Process multiple files concurrently
- **Caching:** Cache query results when possible

---

## 10. Error Handling

### 10.1 Frontend Error Handling

```typescript
// frontend/src/hooks/useQuery.ts

export function useQuery() {
  const [error, setError] = useState<Error | null>(null);
  const [loading, setLoading] = useState(false);
  
  const execute = async (query: UserQuery) => {
    try {
      setLoading(true);
      setError(null);
      
      // Validate guarantees first
      const violations = validateGuarantees(query);
      if (violations.length > 0) {
        throw new GuaranteeViolationError(violations);
      }
      
      const result = await executeQuery(query);
      return result;
    } catch (err) {
      setError(err);
      throw err;
    } finally {
      setLoading(false);
    }
  };
  
  return { execute, error, loading };
}
```

### 10.2 Backend Error Handling

```go
// app/queries/engine.go

func (e *QueryEngine) ExecuteQuery(ctx context.Context, query UserQuery) (*QueryResult, error) {
    // Validate guarantees
    violations := e.validateGuarantees(query)
    if len(violations) > 0 {
        return nil, &GuaranteeViolationError{
            Violations: violations,
            Message:    "Query violates guarantees",
        }
    }
    
    // Execute with error handling
    result, err := e.processQuery(ctx, query)
    if err != nil {
        // Log error for learning
        log.Printf("Query failed: %v", err)
        return nil, fmt.Errorf("query execution failed: %w", err)
    }
    
    return result, nil
}
```

---

## 11. Implementation Checklist

### Phase 1: Foundation
- [ ] Wails project setup
- [ ] React + TypeScript + Tailwind configuration
- [ ] Flowbite integration
- [ ] Basic routing
- [ ] Information request data structure

### Phase 2: UI Components
- [ ] Step 1: Information request selector
- [ ] Step 2: Date range selector
- [ ] Step 3: Category selector
- [ ] Step 4: Keywords input
- [ ] Step 5: Review & confirm
- [ ] Processing screen
- [ ] Download screen
- [ ] Progress indicator component

### Phase 3: Backend
- [ ] Query engine structure
- [ ] Parquet file reader
- [ ] Filter implementation
- [ ] Aggregation logic
- [ ] Zip file generator
- [ ] Guarantee validation

### Phase 4: Integration
- [ ] Wails bridge setup
- [ ] Frontend-backend communication
- [ ] End-to-end data flow
- [ ] Error handling
- [ ] Loading states

### Phase 5: Testing
- [ ] Component tests
- [ ] Integration tests
- [ ] Guarantee validation tests
- [ ] Performance tests
- [ ] User acceptance tests

---

## 12. Key Principles to Follow

### 12.1 Research-Driven Development

**Before implementing:**
- Define the problem clearly
- State assumptions explicitly
- Create guarantees (assertions)
- Plan how to test assumptions

**During implementation:**
- Encode assumptions as guarantees
- Test guarantees continuously
- Learn from failures

**After implementation:**
- Validate assumptions
- Document what was learned
- Refine understanding

### 12.2 Single Responsibility

**Each component does one thing:**
- `InformationRequestStep`: Select request (only)
- `DateRangeStep`: Select dates (only)
- `CategoryStep`: Select categories (only)

**Each function does one thing:**
- `validateQuery`: Validate (only)
- `executeQuery`: Execute (only)
- `generateOutput`: Generate (only)

### 12.3 Explicit Contracts

**Every function has:**
- Preconditions (what must be true before)
- Postconditions (what will be true after)
- Guarantees (assertions that encode assumptions)

**Example:**
```typescript
function processQuery(query: UserQuery): QueryResult {
  // Precondition
  assert(query.dateRange.start <= query.dateRange.end);
  
  // Function body
  const result = ...;
  
  // Postcondition
  assert(result.records.length > 0);
  
  return result;
}
```

---

**Document Status:** Implementation guide complete. Ready for development.
