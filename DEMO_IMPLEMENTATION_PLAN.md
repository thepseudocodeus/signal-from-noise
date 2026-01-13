# Demo Implementation Plan
## Email Query Interface - Client Demo

**Goal:** Build a working demo that allows the client to query email data from a data lake on an external drive using a Typeform-inspired desktop interface.

**Focus:** Email capabilities (foundation for broader data lake orchestration)

---

## Problem Definition (Knuth-Style)

### What We Are Given:
- **Data Lake Location:** External drive (path in .env)
- **Email Data:** Parquet files from PST conversion pipeline
- **User Need:** Query email data to answer information requests
- **Current Setup:** Wails + Go + React + TypeScript + Tailwind + Flowbite

### What We Must Produce:
- **Demo Application:** Typeform-inspired step-by-step interface
- **Core Functionality:**
  - Connect to data lake on external drive
  - Allow user to select query parameters (date range, keywords, categories)
  - Process email Parquet files
  - Generate output (CSV/JSON for demo - ZIP later)
- **User Experience:** Engaging, intuitive, one-question-at-a-time flow

### Constraints:
- **Demo Focus:** Email capabilities (claims/other later)
- **Time:** Need working demo quickly
- **Data:** Large data lake on external drive
- **Performance:** Must handle large datasets efficiently

---

## Implementation Strategy: "Make It Work" Phase

### Phase 1: Foundation (Core Setup)
1. **Environment Configuration**
   - Read data lake path from .env
   - Set up Go backend to access external drive
   - Validate data lake accessibility

2. **Basic Backend Structure**
   - Data lake connection/service
   - Email Parquet file discovery
   - Basic query structure

3. **Frontend Foundation**
   - Replace current App.tsx with step-by-step flow
   - Set up routing/state for multi-step process
   - Basic Flowbite components

### Phase 2: Core UI (Step-by-Step Flow)
1. **Step 1: Welcome/Information Request Selection**
   - Simple email query options for demo
   - Typeform-inspired card selection

2. **Step 2: Date Range Selection**
   - Flowbite DatePicker
   - Date range validation

3. **Step 3: Keywords (Optional)**
   - Input field for keywords
   - Tag-based interface

4. **Step 4: Review & Confirm**
   - Summary of selections
   - Confirm button

5. **Step 5: Processing & Results**
   - Progress indicator
   - Display results (simplified for demo)

### Phase 3: Backend Integration
1. **Email Data Processing**
   - Connect to data lake
   - Discover Parquet files
   - Basic filtering (date range, keywords)
   - Generate output (CSV/JSON for demo)

2. **Wails Bridge**
   - Connect frontend to backend
   - Pass query parameters
   - Return results

---

## Technical Approach

### Data Lake Connection
```go
// Read from .env or config
dataLakePath := os.Getenv("DATA_LAKE_PATH")
// Validate path exists and is accessible
// Discover Parquet files in email directory
```

### Email Query (Simplified for Demo)
```go
// For demo: Focus on basic email queries
// 1. Find email Parquet files
// 2. Filter by date range
// 3. Filter by keywords (if provided)
// 4. Return sample results (first N records for demo)
```

### Frontend Flow
```
Welcome → Select Query Type → Date Range → Keywords → Review → Process → Results
```

---

## File Structure

```
signal-from-noise/
├── .env                          # Data lake path
├── app.go                        # Main app struct
├── main.go                       # Entry point
├── config/
│   └── config.go                # Configuration (read .env)
├── datalake/
│   └── datalake.go              # Data lake connection & discovery
├── queries/
│   └── email_query.go           # Email query processing
└── frontend/
    └── src/
        ├── App.tsx              # Main app (step-by-step flow)
        ├── components/
        │   ├── steps/
        │   │   ├── WelcomeStep.tsx
        │   │   ├── QueryTypeStep.tsx
        │   │   ├── DateRangeStep.tsx
        │   │   ├── KeywordsStep.tsx
        │   │   ├── ReviewStep.tsx
        │   │   └── ResultsStep.tsx
        │   └── shared/
        │       └── ProgressIndicator.tsx
        ├── types/
        │   └── query.ts         # TypeScript types
        └── hooks/
            └── useQuery.ts      # Query state management
```

---

## Demo Query Types (Simplified)

For the demo, focus on 2-3 simple email query types:

1. **Email Timeline**
   - All emails in date range
   - Ordered chronologically

2. **Keyword Search**
   - Emails containing keywords
   - In date range

3. **Email Summary**
   - Count of emails
   - Date range statistics

---

## Implementation Order

### Step 1: Environment & Data Lake Connection
- [ ] Read .env file
- [ ] Create config package
- [ ] Create data lake service
- [ ] Test connection to external drive
- [ ] Discover Parquet files

### Step 2: Backend Query Structure
- [ ] Define query models (Go)
- [ ] Create email query service
- [ ] Basic Parquet file reading (or CSV if easier for demo)
- [ ] Filter by date range
- [ ] Return sample results

### Step 3: Frontend Foundation
- [ ] Install date-fns for date handling
- [ ] Create step-by-step component structure
- [ ] Set up state management
- [ ] Create ProgressIndicator component

### Step 4: Step Components
- [ ] WelcomeStep
- [ ] QueryTypeStep
- [ ] DateRangeStep
- [ ] KeywordsStep
- [ ] ReviewStep
- [ ] ResultsStep

### Step 5: Integration
- [ ] Connect frontend to backend (Wails)
- [ ] Test end-to-end flow
- [ ] Add error handling
- [ ] Add loading states

### Step 6: Polish
- [ ] Styling refinements
- [ ] Animations/transitions
- [ ] Error messages
- [ ] Success states

---

## Success Criteria

**Demo is ready when:**
- ✅ Application connects to data lake on external drive
- ✅ User can select query parameters (date range, keywords)
- ✅ Application processes email data
- ✅ Results are displayed
- ✅ UI is Typeform-inspired (one question at a time)
- ✅ Flow is smooth and engaging

---

## Next Steps

1. **Read .env file** to understand data lake structure
2. **Set up data lake connection** in Go backend
3. **Create simplified query structure** for email demo
4. **Build step-by-step UI** with Flowbite components
5. **Connect frontend to backend** via Wails
6. **Test with real data** from external drive
