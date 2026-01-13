# Working Agreement: Signal from Noise Application
## Desktop Email Data Query Interface Development

**Version:** 1.0
**Date:** 2026-01-13
**Context:** Based on research specifications (00a, 00b)

---

## Project Context

### Client Problem
Build a desktop application that allows users to query email data from a large data lake (10s of GB, tens of thousands of files) to answer 20 specific information requests. The application must:
- Process Parquet files from PST conversion pipeline
- Provide intuitive, Typeform-inspired step-by-step interface
- Generate filtered, structured ZIP file outputs
- Be engaging, enjoyable, and intuitive

### Technology Stack
- **Backend:** Wails v2 + Go (data processing, Parquet reading, ZIP generation)
- **Frontend:** React + TypeScript + Tailwind CSS + Flowbite React
- **Design:** Typeform-inspired (one question at a time, smooth transitions)

---

## Our Working Approach

### 1. Development Philosophy (Make It Work → Make It Right → Make It Fast)

**Phase 1: Make It Work**
- Get basic functionality working
- Focus on core user flow
- Accept technical debt temporarily
- **Goal:** Working demo for client

**Phase 2: Make It Right**
- Refactor for maintainability
- Add proper error handling
- Implement guarantees (as per spec)
- **Goal:** Production-quality code

**Phase 3: Make It Fast**
- Optimize performance
- Add caching
- Stream processing
- **Goal:** Handle 10s of GB efficiently

### 2. Problem Definition (Knuth-Style)

Before any implementation:
1. **Define the problem clearly:** What are we building?
2. **State assumptions explicitly:** What do we assume?
3. **Create guarantees:** What must be true?
4. **Plan tests:** How do we verify?

**Reference:** `research/00a_desktop_application_specification.md` for full problem definition

### 3. Modular Development

**Work function-by-function, module-by-module:**
- Each module has a clear responsibility
- Each function has a single purpose
- Dependencies are explicit
- Tests accompany implementation

**Module Structure:**
```
frontend/
  src/
    components/
      steps/          # Step-by-step UI components
      shared/         # Shared UI components
    hooks/            # Custom React hooks
    types/            # TypeScript type definitions
    services/         # API/service layer
    utils/            # Utility functions

app/                  # Go backend
  queries/            # Query engine
  data/               # Data access layer
  models/             # Data models
```

### 4. Planning First (NASA/Senior Engineer Approach)

**For each module:**
1. **Plan:** Create implementation plan with milestones
2. **Implement:** Build function by function
3. **Verify:** Test each milestone
4. **Refine:** Adjust based on learnings

**Never skip planning phase** - this prevents rework and ensures we build the right thing.

### 5. Explicit Guarantees (NASA Principles)

**For each feature:**
- **Preconditions:** What must be true before?
- **Postconditions:** What will be true after?
- **Invariants:** What stays true throughout?
- **Assertions:** How do we check?

**Reference:** `research/00a_desktop_application_specification.md` Section 6 (Guarantees and Validation)

**Example:**
- **G2:** Date range is valid
  - **Assertion:** `startDate <= endDate && startDate >= minDate && endDate <= maxDate`
  - **Validation:** At Step 2 (date range selection)

### 6. Rapid Application Development with Code Generation

**Use Flowbite components as building blocks:**
- Leverage pre-built components (Cards, DatePicker, Checkboxes, etc.)
- Customize as needed
- Generate UI quickly using Flowbite patterns

**Code generation principles:**
- Use templates where appropriate
- Generate repetitive code (don't write manually)
- Focus on business logic, not boilerplate

### 7. Learning and Iteration

**Test assumptions continuously:**
- Define assumptions explicitly
- Test assumptions with data
- Learn from failures
- Refine understanding

**Example workflow:**
1. **Assumption:** Users prefer 5 steps
2. **Test:** Build 5-step UI, gather feedback
3. **Learn:** If users find it too many, adjust
4. **Refine:** Update assumption and implementation

---

## Development Workflow

### Step-by-Step Process

**For each feature/module:**

1. **Define the problem** (Knuth-style)
   - What are we building?
   - What are the inputs/outputs?
   - What are the constraints?

2. **Plan the implementation** (NASA-style)
   - Break into milestones
   - Define function signatures
   - Identify dependencies
   - Plan tests

3. **Implement milestone-by-milestone** (Make it work)
   - Work function by function
   - Get each milestone working
   - Test as we go

4. **Refactor** (Make it right)
   - Clean up technical debt
   - Add proper error handling
   - Implement guarantees

5. **Optimize** (Make it fast)
   - Profile performance
   - Optimize bottlenecks
   - Add caching where needed

### Code Review Checklist

Before marking a module complete:
- [ ] Problem clearly defined
- [ ] Plan documented
- [ ] Guarantees defined
- [ ] Tests written
- [ ] Error handling implemented
- [ ] TypeScript types complete
- [ ] No `any` types (unless necessary)
- [ ] Code follows single responsibility
- [ ] Functions are testable
- [ ] Documentation updated

---

## Key Principles

### 1. Explicit is Better Than Implicit
- Clear function names
- Explicit types
- Explicit guarantees
- Explicit error handling

### 2. Single Responsibility
- Each module does one thing
- Each function does one thing
- Each component does one thing

### 3. Test-Driven Learning
- Test assumptions
- Test guarantees
- Test edge cases
- Learn from test failures

### 4. Modular Architecture
- Clear module boundaries
- Explicit dependencies
- Testable in isolation
- Composable for larger features

### 5. Data-Driven Decisions
- Measure performance
- Track user behavior
- Validate assumptions with data
- Refine based on metrics

---

## Communication Style

### When Asking Questions
- **Be specific:** What exactly are we trying to understand?
- **State assumptions:** What do we assume?
- **Propose solution:** What do we think the answer might be?

### When Implementing
- **Show progress:** Share what's working
- **Show problems:** Share what's not working
- **Show learnings:** Share what we learned

### When Planning
- **Define scope:** What are we building?
- **Break down:** How do we break it into steps?
- **Identify risks:** What could go wrong?

---

## Project-Specific Guidelines

### Frontend (React + TypeScript + Flowbite)

**Component Structure:**
- Use Flowbite React components as base
- Customize styling with Tailwind
- Follow Typeform-inspired design principles
- One question at a time
- Smooth transitions
- Clear progress indicators

**State Management:**
- React hooks for local state
- Context for shared state (if needed)
- Explicit types for all state

**TypeScript:**
- 100% type coverage (target)
- No `any` types (unless necessary)
- Strict mode enabled
- Explicit function signatures

### Backend (Go)

**Structure:**
- Query engine in `app/queries/`
- Data access in `app/data/`
- Models in `app/models/`
- Main app logic in `app/`

**Guarantees:**
- Validate all inputs
- Check guarantees explicitly
- Return clear errors
- Log for learning

**Performance:**
- Stream processing (don't load all into memory)
- Predicate pushdown (filter at Parquet level)
- Parallel processing where possible
- Profile and optimize bottlenecks

### Data Processing

**Parquet Files:**
- Read with predicate pushdown (date range)
- Filter by category
- Filter by keywords
- Stream results (don't load all into memory)
- Aggregate and order
- Package into ZIP

**ZIP Generation:**
- Include manifest.json
- Include data.csv (or Parquet if large)
- Include summary.json
- Include readme.txt
- Name: `{information_request_id}_{timestamp}.zip`

---

## Success Metrics

### Development Metrics
- **Type coverage:** >95%
- **Test coverage:** >80% (target)
- **Build time:** <2 seconds (baseline established: ~1.13s)
- **Zero compilation errors**
- **Zero type errors**

### Application Metrics (Post-Demo)
- **User completion time:** <2 minutes (90% of users)
- **User satisfaction:** >=4.5/5
- **Query processing time:** <30 seconds (90% of queries)
- **Memory usage:** <2GB during processing

---

## Reference Documents

- **Specification:** `research/00a_desktop_application_specification.md`
- **Implementation Guide:** `research/00b_application_implementation_guide.md`
- **TypeScript Conversion:** `JSTOTS.md`
- **TypeScript Verification:** `TS_VERIFICATION_RESULTS.md`

---

## Next Steps

1. **Review this agreement:** Ensure we're aligned
2. **Define first module:** What do we build first?
3. **Plan first module:** Break into milestones
4. **Implement first module:** Make it work, make it right, make it fast
5. **Iterate:** Learn, refine, repeat

---

**Status:** Active working agreement
**Last Updated:** 2026-01-13
**Review Frequency:** Update as project evolves
