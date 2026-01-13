# Desktop Application Specification: Email Data Query Interface

**Version:** 1.0  
**Date:** 2026-01-XX  
**Status:** Specification Phase  
**Target:** Wails + Go + React + Tailwind + Flowbite Desktop Application

---

## 1. Problem Definition (Knuth-Style)

### 1.1 What We Are Given

**Input Domain:**
- **Data Lake:** Full data lake of information (tens of thousands of files, 10s of GB)
- **Email Results:** Final email results from PST conversion pipeline (Parquet files)
- **User Intent:** One of 20 information requests that must be answered
- **User Constraints:** Date ranges, keywords, categories (Claims/Emails/Other)
- **Data Characteristics:** Noisy, chaotic, unstructured → needs filtering/ordering/exclusion

**Constraints:**
- Data is large (10s of GB)
- Data is unstructured/chaotic
- User needs to navigate complex filtering
- Must produce specific outputs for 20 information requests
- Must be enjoyable, engaging, intuitive (Typeform-inspired)

### 1.2 What We Must Produce

**Output Domain:**
- **Zip File:** Named `{information_request_id}_{timestamp}.zip`
- **Contents:** Filtered, ordered, structured data relevant to the information request
- **Format:** Structured, valuable, insightful information

**Output Requirements:**
1. **Relevance:** Data must answer the specific information request
2. **Structure:** Data must be organized and queryable
3. **Completeness:** All relevant data included (within constraints)
4. **Usability:** Output must be immediately useful
5. **Traceability:** Must be able to trace back to source data

### 1.3 Transformation Function

```
T: (DataLake, EmailResults, UserSelections) → ZipFile

where:
UserSelections = {
    information_request: InformationRequest (1 of 20),
    date_range: (start_date, end_date),
    categories: Set[Category],  // Claims, Emails, Other
    keywords: List[String],
    filters: FilterSet,
    ordering: Ordering
}
```

**Invariants:**
- `relevance(output, information_request) = True`
- `completeness(output, filters) = True`
- `traceability(output, source_data) = True`

---

## 2. System Architecture

### 2.1 Technology Stack

**Frontend:**
- **Framework:** React (TypeScript)
- **UI Library:** Flowbite React
- **Styling:** Tailwind CSS
- **Inspiration:** Typeform (step-by-step, engaging, intuitive)

**Backend:**
- **Framework:** Wails v2
- **Language:** Go
- **Data Processing:** Go + Parquet readers

**Data:**
- **Input:** Parquet files from PST conversion pipeline
- **Query Engine:** Efficient filtering/aggregation
- **Output:** Zip files with structured data

### 2.2 Application Flow

```
┌─────────────────┐
│  Welcome Screen │
└────────┬────────┘
         │
         ▼
┌─────────────────────────┐
│ Step 1: Select          │
│ Information Request     │
│ (1 of 20)              │
└────────┬────────────────┘
         │
         ▼
┌─────────────────────────┐
│ Step 2: Select          │
│ Date Range              │
│ (Calendar Component)    │
└────────┬────────────────┘
         │
         ▼
┌─────────────────────────┐
│ Step 3: Select          │
│ Categories              │
│ (Claims/Emails/Other)   │
└────────┬────────────────┘
         │
         ▼
┌─────────────────────────┐
│ Step 4: Keywords        │
│ (Optional)              │
└────────┬────────────────┘
         │
         ▼
┌─────────────────────────┐
│ Step 5: Review &        │
│ Confirm                 │
└────────┬────────────────┘
         │
         ▼
┌─────────────────────────┐
│ Processing...           │
│ (Progress Indicator)    │
└────────┬────────────────┘
         │
         ▼
┌─────────────────────────┐
│ Download Zip File       │
│ {request_id}_{ts}.zip  │
└─────────────────────────┘
```

---

## 3. Information Requests (20 Items)

### 3.1 Request Structure

Each information request has:
- **ID:** Unique identifier (e.g., "IR-001")
- **Title:** Human-readable title
- **Description:** What information is needed
- **Data Requirements:** What data sources needed
- **Output Format:** Expected output structure
- **Validation Criteria:** How to verify completeness

### 3.2 Example Requests (Template)

**IR-001: Email Communication Timeline**
- **Description:** All emails within date range, ordered chronologically
- **Data:** Email Parquet files
- **Output:** CSV with columns: timestamp, from, to, subject, body_hash
- **Filters:** Date range, category (Emails)

**IR-002: Claims-Related Communications**
- **Description:** All emails/communications related to claims
- **Data:** Email Parquet files + Claims data
- **Output:** Structured data linking emails to claims
- **Filters:** Category (Claims), keywords, date range

**IR-003: [Template for remaining 18]**
- **Description:** [To be defined]
- **Data:** [To be defined]
- **Output:** [To be defined]
- **Filters:** [To be defined]

**Note:** The actual 20 requests should be defined based on client requirements. This structure provides the template.

---

## 4. User Interface Specification

### 4.1 Design Principles

**Typeform-Inspired:**
- **One question at a time:** Focus user attention
- **Progress indicator:** Show completion status
- **Smooth transitions:** Engaging animations
- **Clear CTAs:** Obvious next steps
- **Visual feedback:** Immediate response to actions

**Expert Principles:**
- **Explicit:** Clear what each step does
- **Validated:** User selections validated before proceeding
- **Traceable:** Can see what was selected
- **Reproducible:** Same selections = same output

### 4.2 Step-by-Step Flow

#### Step 1: Information Request Selection

**Component:** Flowbite Card Grid or Select Component

**Requirements:**
- Display all 20 information requests
- Visual cards with title and brief description
- Single selection (radio-style)
- Progress: "Step 1 of 5"
- Next button (disabled until selection)

**Validation:**
- Must select exactly one request
- Selection stored in state

**UI Elements:**
```tsx
<InformationRequestSelector
  requests={informationRequests}
  selected={selectedRequest}
  onSelect={handleSelect}
  progress={1/5}
/>
```

#### Step 2: Date Range Selection

**Component:** Flowbite Calendar/DatePicker

**Requirements:**
- Constrain to valid date range (from data)
- Show min/max dates available
- Visual calendar interface
- Start date and end date selection
- Clear indication of selected range
- Validation: end_date >= start_date

**UI Elements:**
```tsx
<DateRangeSelector
  minDate={dataMinDate}
  maxDate={dataMaxDate}
  startDate={startDate}
  endDate={endDate}
  onStartChange={handleStartDate}
  onEndChange={handleEndDate}
  progress={2/5}
/>
```

#### Step 3: Category Selection

**Component:** Flowbite Checkbox Group

**Requirements:**
- Multiple selection allowed
- Options: Claims, Emails, Other
- Visual checkboxes with icons
- At least one must be selected
- Show count of selected items

**UI Elements:**
```tsx
<CategorySelector
  categories={['Claims', 'Emails', 'Other']}
  selected={selectedCategories}
  onToggle={handleCategoryToggle}
  progress={3/5}
  minSelection={1}
/>
```

#### Step 4: Keywords (Optional)

**Component:** Flowbite Input with Tags

**Requirements:**
- Multi-keyword input
- Tag-based interface
- Add/remove keywords easily
- Optional step (can skip)
- Show keyword count

**UI Elements:**
```tsx
<KeywordSelector
  keywords={keywords}
  onAdd={handleAddKeyword}
  onRemove={handleRemoveKeyword}
  progress={4/5}
  optional={true}
/>
```

#### Step 5: Review & Confirm

**Component:** Summary Card

**Requirements:**
- Show all selections
- Allow editing (go back)
- Show estimated result size
- Confirm button
- Progress: "Step 5 of 5"

**UI Elements:**
```tsx
<ReviewSummary
  request={selectedRequest}
  dateRange={dateRange}
  categories={selectedCategories}
  keywords={keywords}
  estimatedSize={estimatedSize}
  onConfirm={handleConfirm}
  onEdit={handleEdit}
  progress={5/5}
/>
```

#### Processing Screen

**Component:** Progress Indicator

**Requirements:**
- Show processing status
- Progress bar
- Estimated time remaining
- Cancel option (if possible)

**UI Elements:**
```tsx
<ProcessingScreen
  status={processingStatus}
  progress={progressPercent}
  estimatedTime={estimatedSeconds}
  onCancel={handleCancel}
/>
```

#### Download Screen

**Component:** Success Card with Download Button

**Requirements:**
- Show success message
- Display zip file name
- Download button
- Option to start new query
- Show file size

**UI Elements:**
```tsx
<DownloadScreen
  fileName={zipFileName}
  fileSize={fileSizeBytes}
  onDownload={handleDownload}
  onNewQuery={handleNewQuery}
/>
```

---

## 5. Data Processing Specification

### 5.1 Query Engine

**Purpose:** Efficiently filter and aggregate data from Parquet files

**Requirements:**
- **Performance:** Handle 10s of GB efficiently
- **Filtering:** Apply date range, categories, keywords
- **Aggregation:** Group and order results
- **Streaming:** Process in chunks (don't load all into memory)

**Algorithm:**
```
1. Load information request definition
2. Determine required data sources
3. For each Parquet file:
   a. Read with predicate pushdown (date range)
   b. Filter by category
   c. Filter by keywords (if provided)
   d. Apply request-specific filters
   e. Stream results to output
4. Aggregate and order results
5. Package into zip file
```

### 5.2 Output Generation

**Format:** Zip file containing:
- **manifest.json:** Metadata about the query
- **data.csv:** Main data (or Parquet if large)
- **summary.json:** Query summary and statistics
- **readme.txt:** Description of contents

**Naming:** `{information_request_id}_{timestamp}.zip`

**Example:** `IR-001_20260113_143022.zip`

---

## 6. Guarantees and Validation

### 6.1 User Input Guarantees

**G1:** Information request is valid
- **Assertion:** `selectedRequest in validRequests`
- **Validation:** At Step 1

**G2:** Date range is valid
- **Assertion:** `startDate <= endDate && startDate >= minDate && endDate <= maxDate`
- **Validation:** At Step 2

**G3:** At least one category selected
- **Assertion:** `selectedCategories.length >= 1`
- **Validation:** At Step 3

**G4:** Keywords are valid (if provided)
- **Assertion:** `keywords.every(k => k.length > 0 && k.length <= 100)`
- **Validation:** At Step 4

### 6.2 Data Processing Guarantees

**G5:** Output contains relevant data
- **Assertion:** `output.relevance_score >= threshold`
- **Validation:** After processing

**G6:** Output is complete (within constraints)
- **Assertion:** `output.completeness >= 0.95`
- **Validation:** After processing

**G7:** Output is traceable
- **Assertion:** `output.manifest.source_files.length > 0`
- **Validation:** After processing

---

## 7. Research Questions

### 7.1 UX Research Questions

**RQ1:** What is the optimal number of steps for user engagement?
- **Hypothesis:** 5 steps is optimal (not too few, not too many)
- **Test:** A/B test with 3, 5, 7 steps

**RQ2:** What calendar component provides best date selection UX?
- **Hypothesis:** Flowbite calendar with range selection
- **Test:** Compare with other calendar libraries

**RQ3:** How do users prefer to select categories?
- **Hypothesis:** Checkboxes with visual icons
- **Test:** Compare with dropdown, radio buttons

### 7.2 Performance Research Questions

**RQ4:** What is the optimal chunk size for Parquet reading?
- **Hypothesis:** 10,000 rows per chunk
- **Test:** Benchmark different chunk sizes

**RQ5:** How does predicate pushdown affect query performance?
- **Hypothesis:** 10× improvement with pushdown
- **Test:** Compare with/without pushdown

---

## 8. Assumptions

### 8.1 UX Assumptions

**A1:** Users prefer step-by-step over single-page forms
- **Rationale:** Typeform success, cognitive load reduction
- **Risk:** Medium
- **Validation:** User testing

**A2:** Calendar component is intuitive for date selection
- **Rationale:** Common UI pattern
- **Risk:** Low
- **Validation:** Usability testing

**A3:** Multiple category selection is preferred
- **Rationale:** Flexibility, common pattern
- **Risk:** Low
- **Validation:** User feedback

### 8.2 Technical Assumptions

**A4:** Parquet files can be efficiently queried with Go
- **Rationale:** Parquet is columnar, efficient
- **Risk:** Medium
- **Validation:** Performance benchmarks

**A5:** Wails provides sufficient performance for data processing
- **Rationale:** Go backend, efficient
- **Risk:** Low
- **Validation:** Load testing

---

## 9. Expected Outcomes

### 9.1 User Experience Outcomes

**Prediction P1:** Users complete queries in < 2 minutes
- **Measurement:** Time from start to download
- **Target:** 90% of users < 2 minutes

**Prediction P2:** User satisfaction score >= 4.5/5
- **Measurement:** Post-query survey
- **Target:** 4.5/5 average

### 9.2 Technical Outcomes

**Prediction P3:** Query processing time < 30 seconds for typical queries
- **Measurement:** Time from confirm to zip ready
- **Target:** < 30 seconds for 90% of queries

**Prediction P4:** Memory usage < 2GB during processing
- **Measurement:** Peak memory during query
- **Target:** < 2GB

---

## 10. Implementation Phases

### Phase 1: Foundation (Week 1-2)
- Wails project setup
- React + Tailwind + Flowbite setup
- Basic routing and navigation
- Information request data structure

### Phase 2: Core UI (Week 3-4)
- Step 1: Information request selection
- Step 2: Date range selection
- Step 3: Category selection
- Step 4: Keywords (optional)
- Step 5: Review & confirm

### Phase 3: Backend (Week 5-6)
- Go data processing engine
- Parquet file reading
- Filtering and aggregation
- Zip file generation

### Phase 4: Integration (Week 7)
- Connect frontend to backend
- End-to-end testing
- Performance optimization

### Phase 5: Polish (Week 8)
- UX refinements
- Error handling
- Loading states
- Documentation

---

## 11. Expert Principles Applied

### 11.1 Knuth Principles

- **Mathematical Rigor:** Formal specification of data transformations
- **Algorithm Correctness:** Guarantees ensure correctness
- **Documentation:** Complete specifications before implementation

### 11.2 Simons Principles

- **Systematic Approach:** Step-by-step, methodical
- **Quantitative Validation:** Measure everything (time, satisfaction, performance)
- **Reproducibility:** Same inputs = same outputs

### 11.3 Jane Street Principles

- **Type Safety:** TypeScript for frontend, strong types in Go
- **Explicit Contracts:** Guarantees encode assumptions
- **Defensive Programming:** Validate all inputs, handle errors gracefully

---

## 12. Learning Opportunities

### 12.1 Assumption Falsification

**When a guarantee fails:**
- Document what failed
- Analyze why
- Update assumptions
- Refine understanding

**Example:**
- **Assumption:** Users prefer 5 steps
- **Reality:** Users find 5 steps too many
- **Learning:** Reduce to 3-4 steps
- **Update:** Refine UX assumptions

### 12.2 Continuous Improvement

**Metrics to Track:**
- Query completion time
- User satisfaction
- Error rates
- Performance metrics

**Use metrics to:**
- Validate assumptions
- Identify improvements
- Refine hypotheses

---

**Document Status:** Specification complete. Ready for implementation planning.
