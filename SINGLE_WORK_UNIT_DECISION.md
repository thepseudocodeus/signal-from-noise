# Single Work Unit Decision Framework
## Expert Approach: Choose the ONE Thing

**Time Available:** 1 unit of work
**Principle:** Work on the 1 thing that reduces possible choices the most
**Goal:** Maximum uncertainty reduction, maximum confidence increase

---

## Decision Framework

### Step 1: Identify Uncertainties

**Current Uncertainties (Ranked by Impact):**

1. **Does the demo actually work?** (Highest Impact)
   - Risk: Demo might not run
   - Impact: Can't demonstrate value
   - Uncertainty: High (we haven't tested runtime)

2. **Does Wails communication work?** (High Impact)
   - Risk: Frontend can't call backend
   - Impact: Core functionality broken
   - Uncertainty: Medium (code compiles, but runtime untested)

3. **Does data lake connection work?** (High Impact)
   - Risk: Can't access USB drive
   - Impact: Can't demonstrate data access
   - Uncertainty: Medium (path configured, but untested)

4. **What does client want next?** (Medium Impact - if meeting happened)
   - Risk: Working on wrong thing
   - Impact: Misaligned priorities
   - Uncertainty: Depends on meeting outcome

---

### Step 2: Which Reduces Choices Most?

**Option A: Verify Demo Works**
- **Uncertainty Reduced:** Does the system work at runtime?
- **Choices Reduced:** Eliminates "does it work?" uncertainty
- **Impact:** High - enables demonstration, proves capability
- **Time:** ~15-30 minutes
- **Success Criteria:** Demo runs, shows UI, connects to backend

**Option B: Document Meeting (if happened)**
- **Uncertainty Reduced:** What does client want?
- **Choices Reduced:** Eliminates "what to work on?" uncertainty
- **Impact:** Medium-High - directs future work
- **Time:** ~15 minutes
- **Success Criteria:** Meeting outcomes documented

**Option C: Add One More Feature**
- **Uncertainty Reduced:** Can we add features?
- **Choices Reduced:** Low - doesn't reduce core uncertainty
- **Impact:** Low - doesn't address core risks
- **Time:** Variable
- **Success Criteria:** Feature added

---

### Step 3: Expert Decision

**Highest Impact = Verify Demo Works**

**Why:**
- Reduces the biggest uncertainty (does it work?)
- Enables demonstration (increases confidence)
- Proves capability (validates assumptions)
- Low risk (just testing)
- Quick win (15-30 minutes)

---

## Recommended Work Unit: Verify Demo Works

### Objective
Verify that `wails dev` starts successfully and the application runs, reducing uncertainty about whether the system works at runtime.

### Why This Over Others
- **Reduces most uncertainty:** Answers "does it work?"
- **Enables demonstration:** Can show working system
- **Proves capability:** Validates Wails integration
- **Quick win:** 15-30 minutes
- **Low risk:** Just testing, no breaking changes

### Success Criteria
- ✅ `wails dev` starts without errors
- ✅ Application window opens
- ✅ Splash page displays
- ✅ Data lake status shows (even if "unavailable")
- ✅ No critical errors in console

### Procedure
1. Run `wails dev`
2. Observe startup (errors? warnings?)
3. Check if application opens
4. Test splash page loads
5. Document results (works/doesn't work/partially works)

### Deliverable
- Status: Demo works / doesn't work / partially works
- Documentation: What worked, what didn't
- Next steps: Based on results

---

## Alternative: If Meeting Already Happened

**If client meeting already occurred:**
- **Work Unit:** Document meeting outcomes
- **Why:** Captures data, reduces "what does client want?" uncertainty
- **Time:** 15 minutes
- **Deliverable:** Meeting debrief document

---

## Decision Matrix

| Option | Uncertainty Reduced | Time | Risk | Impact | Choice |
|--------|-------------------|------|------|--------|--------|
| Verify Demo | High (does it work?) | 15-30 min | Low | High | ✅ **BEST** |
| Document Meeting | Medium (what does client want?) | 15 min | Low | Medium-High | ⚠️ If meeting happened |
| Add Feature | Low | Variable | Medium | Low | ❌ Not recommended |

---

## Recommendation

**Do This:** Verify Demo Works

**Steps:**
1. Run `wails dev`
2. Document what happens
3. Record results
4. Update status

**Time:** 15-30 minutes
**Impact:** High - reduces biggest uncertainty
**Risk:** Low - just testing

---

**Expert Principle Applied:** Work on the 1 thing that reduces possible choices the most.

**In This Case:** Verifying the demo works reduces the most uncertainty (does it work?) and enables future demonstrations.
