# Calendar Component Fix: Expert Debugging Process

## Problem

Calendar component (Datepicker) doesn't display anything.

## Expert Approach: Systematic Hypothesis Testing

### Step 1: State the Assumptions

**Assumption 1**: Flowbite React Datepicker works with default setup

- **Test**: Check documentation and examples
- **Result**: ✓ Dependencies installed, Tailwind configured

**Assumption 2**: We're using the component API correctly

- **Test**: Check our implementation vs. Flowbite docs
- **Result**: ❌ **FOUND THE ISSUE**

### Step 2: Identify the Root Cause

**Hypothesis**: Wrong prop names or value format

**Evidence Found:**

1. We used `onSelectedDateChanged` → Should be `onChange`
2. We used `value={date?.toLocaleDateString() || ""}` → Should be `value={date || null}` (Date object, not string)

**From Flowbite React Documentation:**

- `value` prop accepts: `Date | null` (not string)
- Event handler is: `onChange` (not `onSelectedDateChanged`)

### Step 3: Fix the Implementation

**Before (Incorrect):**

```tsx
<Datepicker
  value={dateRange.start?.toLocaleDateString() || ""} // ❌ String
  onSelectedDateChanged={handleStartDateChange} // ❌ Wrong prop name
/>
```

**After (Correct):**

```tsx
<Datepicker
  value={dateRange.start || null} // ✅ Date object or null
  onChange={handleStartDateChange} // ✅ Correct prop name
/>
```

### Step 4: Verify the Fix

**What We Changed:**

1. ✅ Changed `onSelectedDateChanged` → `onChange`
2. ✅ Changed `value` from string to Date object or null
3. ✅ Updated handler to accept `Date | null`

**Why This Works:**

- Flowbite React Datepicker expects Date objects, not strings
- The `onChange` prop is the correct event handler name
- Passing `null` clears the date (as per Flowbite API)

## Expert Lessons

1. **Verify Assumptions**: Don't assume API - check documentation
2. **Test Hypotheses**: Systematically test each assumption
3. **Check Types**: Type mismatches (string vs Date) cause silent failures
4. **Read Documentation**: Flowbite docs clearly state the API

## Prevention

**For Future Components:**

1. Always check official documentation first
2. Verify prop types match expected types
3. Test with minimal example from docs
4. Use TypeScript to catch type mismatches
