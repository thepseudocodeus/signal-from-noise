# Research Paper Approach: Signal from Noise Application
## Current State, Research Questions, and Validation Plan

**Date:** 2026-01-13
**Approach:** Expert-guided research methodology (Knuth, Simons, Renaissance, Jane Street)
**Status:** Initial Implementation Complete - Ready for Validation Phase

---

## Abstract

This document treats the Signal from Noise desktop application development as a research project. We document the current state, define research questions, state explicit assumptions, create testable hypotheses, and design validation experiments.

---

## 1. Current State Documentation

### 1.1 What We Have Built

**Backend (Go + Wails):**
- ✅ Configuration system (`config/config.go`) - Reads .env, validates paths
- ✅ Data lake service (`datalake/datalake.go`) - Connects to USB external drive, discovers files
- ✅ App integration (`app.go`) - Wails methods exposed to frontend
- ✅ Methods: `GetDataLakeStatus()`, `GetEmailFileCount()`, `ListEmailFiles()`, `GetEmailPath()`, `GetDataLakePath()`

**Frontend (React + TypeScript + Flowbite):**
- ✅ Typeform-inspired step-by-step UI
- ✅ Splash page with Wails backend connection
- ✅ Production request selection step
- ✅ Year range selection step
- ✅ Results display step
- ✅ Progress indicators
- ✅ TypeScript types defined

**Integration:**
- ✅ Wails bindings generated (TypeScript ↔ Go)
- ✅ Frontend calls backend methods
- ✅ Data flows from USB drive → Go backend → React frontend

### 1.2 What Works

**Verified:**
- ✅ Code compiles (Go and TypeScript)
- ✅ Build succeeds
- ✅ Wails bindings generated correctly
- ✅ Type safety (100% type coverage)
- ✅ USB drive path configured (`.env` file)

**Not Yet Verified:**
- ⏳ Runtime: Does `wails dev` start successfully?
- ⏳ Backend: Do backend methods execute correctly?
- ⏳ Data: Does data lake connection work with real USB drive?
- ⏳ UI: Does UI render and function correctly?
- ⏳ Communication: Does frontend-backend communication work?

### 1.3 Known Limitations

1. **Simplified Flow:** Only 3 steps (splash → production request → year range → results)
2. **Limited Production Requests:** 3-4 demo requests (not all 20)
3. **Placeholder Data:** Results step shows simulated data
4. **No Query Execution:** No actual data processing yet
5. **No File Processing:** No Parquet/CSV reading implemented
6. **No Progress Updates:** File counts don't update after filtering

---

## 2. Research Questions

### 2.1 Core Research Questions

**RQ1: Does Wails Backend-Frontend Communication Work?**
- **Question:** Can the React frontend successfully call Go backend methods?
- **Sub-questions:**
  - Do Wails bindings work correctly?
  - Can frontend get data lake status from backend?
  - Can frontend get file counts from backend?
  - Are errors handled gracefully?

**RQ2: Does Data Lake Connection Work?**
- **Question:** Can the application successfully connect to and read from the USB external drive?
- **Sub-questions:**
  - Is the USB drive path correctly configured?
  - Can the application detect if the drive is mounted?
  - Can the application discover email files on the drive?
  - Are file counts accurate?

**RQ3: Is the UI Flow Usable?**
- **Question:** Does the Typeform-inspired step-by-step interface provide a good user experience?
- **Sub-questions:**
  - Is navigation intuitive?
  - Are progress indicators helpful?
  - Is validation working correctly?
  - Is the flow engaging?

**RQ4: Can We Expand the System Incrementally?**
- **Question:** Is the architecture suitable for incremental expansion?
- **Sub-questions:**
  - Can we add more steps easily?
  - Can we add more Wails methods easily?
  - Can we add more filtering options easily?
  - Does the code structure support expansion?

### 2.2 Validation Research Questions

**RQ5: What Performance Characteristics Does the System Have?**
- **Question:** How does the system perform with real data?
- **Sub-questions:**
  - How long does it take to discover files?
  - How long does it take to get file counts?
  - Is the UI responsive?
  - Are there memory/CPU issues?

**RQ6: What Error Cases Need Handling?**
- **Question:** What failures can occur and how should they be handled?
- **Sub-questions:**
  - What happens if USB drive is disconnected?
  - What happens if no files are found?
  - What happens if backend methods fail?
  - Are error messages helpful?

---

## 3. Explicit Assumptions

### 3.1 Technical Assumptions

**A1: Wails Works as Expected**
- **Assumption:** Wails v2 successfully bridges Go backend and React frontend
- **Risk:** Medium (Wails is mature, but integration could have issues)
- **Validation:** Test `wails dev` and verify method calls work

**A2: USB Drive is Accessible**
- **Assumption:** USB external drive is mounted at `/Volumes/DirectFileStore/DELEMAR/`
- **Risk:** Medium (Drive might not be mounted, permissions might be wrong)
- **Validation:** Test with real USB drive, handle unmounted case

**A3: File Discovery Works**
- **Assumption:** `DiscoverEmailFiles()` correctly finds email files
- **Risk:** Low (Standard file system operations)
- **Validation:** Verify file counts match actual files

**A4: Typeform-Inspired UI is Effective**
- **Assumption:** Step-by-step interface is intuitive for users
- **Risk:** Medium (User experience is subjective)
- **Validation:** User testing, but for demo: verify navigation works

### 3.2 Data Assumptions

**A5: Email Files Exist**
- **Assumption:** Email files exist in `emails_final` directory
- **Risk:** High (Directory might be empty - we saw it was empty earlier)
- **Validation:** Check if directory has files, handle empty case

**A6: File Formats are Standard**
- **Assumption:** Files are in expected formats (Parquet, CSV, JSON)
- **Risk:** Medium (File formats might vary)
- **Validation:** Check file extensions, handle unknown formats

### 3.3 User Experience Assumptions

**A7: Users Understand the Flow**
- **Assumption:** Users can navigate the step-by-step interface
- **Risk:** Low (Clear labels, progress indicators)
- **Validation:** Verify navigation works, error messages are clear

**A8: Year Selection is Sufficient**
- **Assumption:** Year selection (not full dates) is adequate for filtering
- **Risk:** Low (Simpler is better for demo)
- **Validation:** Verify year filtering works correctly

---

## 4. Testable Hypotheses

### 4.1 Core Hypotheses

**H1: Wails Communication Hypothesis**
- **Hypothesis:** Frontend can successfully call backend methods and receive responses
- **Test:** Call `GetDataLakeStatus()` and `GetEmailFileCount()` from frontend
- **Expected:** Methods return status string and number, no errors
- **Failure Criteria:** Methods fail, errors occur, no response

**H2: Data Lake Connection Hypothesis**
- **Hypothesis:** Application can connect to USB drive and discover files
- **Test:** Call backend methods with USB drive connected
- **Expected:** Status returns "available", file count > 0 (if files exist)
- **Failure Criteria:** Status returns "unavailable", file count errors

**H3: UI Flow Hypothesis**
- **Hypothesis:** Users can navigate through all steps successfully
- **Test:** Manually test navigation: splash → production request → year range → results
- **Expected:** All steps render, navigation works, validation works
- **Failure Criteria:** Steps don't render, navigation breaks, validation fails

### 4.2 Performance Hypotheses

**H4: File Discovery Performance Hypothesis**
- **Hypothesis:** File discovery completes in < 5 seconds for typical datasets
- **Test:** Measure time to discover files
- **Expected:** File discovery < 5 seconds
- **Failure Criteria:** File discovery > 10 seconds

**H5: UI Responsiveness Hypothesis**
- **Hypothesis:** UI remains responsive during backend operations
- **Test:** Monitor UI while backend methods execute
- **Expected:** UI doesn't freeze, loading states show
- **Failure Criteria:** UI freezes, no feedback

---

## 5. Validation Experiments

### 5.1 Experiment 1: Wails Communication Test

**Objective:** Verify frontend-backend communication works

**Setup:**
1. Start application: `wails dev`
2. Open browser/devtools
3. Navigate to splash page

**Procedure:**
1. Observe splash page loads
2. Check if data lake status displays
3. Check if file count displays
4. Check browser console for errors
5. Check backend logs for method calls

**Success Criteria:**
- ✅ Splash page loads without errors
- ✅ Data lake status displays (available/unavailable)
- ✅ File count displays (number)
- ✅ No console errors
- ✅ Backend methods are called

**Failure Handling:**
- If errors occur, document error messages
- Check Wails bindings are generated correctly
- Verify backend methods are exported correctly
- Check for type mismatches

**Data Collection:**
- Screenshot of splash page
- Console logs
- Backend logs
- Error messages (if any)
- Response times

### 5.2 Experiment 2: Data Lake Connection Test

**Objective:** Verify USB drive connection works

**Setup:**
1. Ensure USB drive is mounted
2. Verify `.env` file has correct path
3. Start application

**Procedure:**
1. Call `GetDataLakeStatus()` from frontend
2. Call `GetEmailFileCount()` from frontend
3. Verify status is "available"
4. Verify file count is returned
5. Test with drive unmounted (if possible)

**Success Criteria:**
- ✅ Status returns "available" when drive mounted
- ✅ File count returns correct number (or 0 if no files)
- ✅ Error handling works when drive unmounted
- ✅ Error messages are clear

**Failure Handling:**
- Document error messages
- Check file permissions
- Verify path is correct
- Test with different drive states

**Data Collection:**
- Status responses
- File counts
- Error messages
- Drive mount status
- File system permissions

### 5.3 Experiment 3: UI Flow Test

**Objective:** Verify step-by-step navigation works

**Setup:**
1. Start application
2. Have test user (or developer) navigate

**Procedure:**
1. Start at splash page
2. Click "Start Search"
3. Select production request
4. Select year range
5. View results
6. Test back navigation
7. Test validation (try to proceed without selections)

**Success Criteria:**
- ✅ All steps render correctly
- ✅ Navigation works (forward and back)
- ✅ Validation prevents invalid progress
- ✅ Progress indicators update correctly
- ✅ UI is responsive and smooth

**Failure Handling:**
- Document UI bugs
- Check for layout issues
- Verify Flowbite components work
- Test on different screen sizes

**Data Collection:**
- Screenshots of each step
- Navigation timing
- User feedback (if available)
- UI bugs/issues

### 5.4 Experiment 4: Architecture Validation

**Objective:** Verify architecture supports incremental expansion

**Setup:**
1. Review code structure
2. Plan to add one new step

**Procedure:**
1. Add new step component (e.g., Data Category step)
2. Add step to main App.tsx
3. Update types if needed
4. Test new step integrates correctly

**Success Criteria:**
- ✅ New step integrates easily
- ✅ No breaking changes to existing steps
- ✅ Code structure remains clean
- ✅ Types are correct

**Failure Handling:**
- Document integration issues
- Refactor if needed
- Update architecture if problems found

**Data Collection:**
- Code changes required
- Time to add new step
- Any refactoring needed

---

## 6. Measurement Plan

### 6.1 Quantitative Metrics

**Performance Metrics:**
- File discovery time (seconds)
- Method call latency (milliseconds)
- UI render time (milliseconds)
- File count accuracy (actual vs. reported)

**Code Metrics:**
- Lines of code added/modified
- Type coverage (maintain 100%)
- Build time (seconds)
- Bundle size (KB)

**Functional Metrics:**
- Number of steps implemented
- Number of Wails methods exposed
- Number of production requests
- Error rate (errors per test)

### 6.2 Qualitative Observations

**User Experience:**
- Navigation intuitiveness
- UI clarity
- Error message helpfulness
- Overall flow smoothness

**Code Quality:**
- Code organization
- Type safety
- Error handling
- Maintainability

---

## 7. Next Steps (Research-Driven)

### 7.1 Immediate Next Steps (Validation Phase)

1. **Run Experiment 1: Wails Communication Test**
   - Start application: `wails dev`
   - Verify frontend-backend communication
   - Document results

2. **Run Experiment 2: Data Lake Connection Test**
   - Test with USB drive connected
   - Verify file discovery works
   - Document results

3. **Run Experiment 3: UI Flow Test**
   - Navigate through all steps
   - Test validation
   - Document results

4. **Document Findings**
   - What works
   - What doesn't work
   - What needs fixing
   - What we learned

### 7.2 Incremental Expansion (Based on Results)

**If Validation Succeeds:**
- Add more Wails methods (query execution)
- Add more steps (Data Category, Claims, etc.)
- Add progress updates
- Expand production requests

**If Validation Reveals Issues:**
- Fix critical issues first
- Document what didn't work
- Refactor if needed
- Re-test

### 7.3 Research Continuation

**Ongoing:**
- Document assumptions that proved false
- Document learnings
- Update hypotheses based on results
- Refine architecture based on findings

---

## 8. Research Paper Structure (For Documentation)

### 8.1 Sections to Maintain

1. **Introduction:** Problem statement, goals
2. **Methodology:** How we're building (incremental, research-driven)
3. **Current State:** What we have (this document)
4. **Research Questions:** What we're trying to answer
5. **Assumptions:** What we assume (explicit)
6. **Hypotheses:** What we expect (testable)
7. **Experiments:** How we test (validation)
8. **Results:** What we found (after validation)
9. **Discussion:** What we learned (after validation)
10. **Next Steps:** What's next (based on results)

### 8.2 Documentation Standards

**For Each Experiment:**
- Objective (what we're testing)
- Setup (what we need)
- Procedure (what we do)
- Success Criteria (what success looks like)
- Results (what we found)
- Discussion (what we learned)

**For Each Assumption:**
- Statement (what we assume)
- Risk (how likely to be wrong)
- Validation (how we test it)
- Result (was it correct?)

---

## 9. Expert Principles Applied

### 9.1 Knuth's Principles

- **Mathematical Rigor:** Formal problem definition, explicit assumptions
- **Documentation:** Everything documented as research
- **Correctness:** Guarantees and validation
- **Learning:** Treat failures as learning opportunities

### 9.2 Simons/Renaissance Principles

- **Quantitative Validation:** Measure everything
- **Data-Driven:** Base decisions on data, not assumptions
- **Systematic:** Structured approach to validation
- **Reproducibility:** Document so results are reproducible

### 9.3 Jane Street Principles

- **Type Safety:** TypeScript for correctness
- **Explicit Contracts:** Clear method signatures
- **Defensive Programming:** Error handling, validation
- **Testing:** Test assumptions explicitly

---

## 10. Conclusion (Current State)

**What We Have:**
- Working code (compiles, builds)
- Wails integration (bindings generated)
- UI structure (step-by-step flow)
- Backend structure (data lake connection)

**What We Need to Validate:**
- Does it actually work at runtime?
- Does data flow correctly?
- Is the UX good?
- Can we expand it?

**What We're Doing Next:**
1. Run validation experiments
2. Document results
3. Fix issues
4. Expand incrementally

**Research Approach:**
- Treat as research project
- Document everything
- Test assumptions
- Learn from results
- Iterate based on data

---

**Status:** Ready for Validation Phase
**Next Action:** Run Experiment 1 (Wails Communication Test)
**Expected Outcome:** Documentation of what works, what doesn't, and next steps
