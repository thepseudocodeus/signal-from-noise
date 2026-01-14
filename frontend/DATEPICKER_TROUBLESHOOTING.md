# Datepicker Troubleshooting Guide

## Expert Systematic Debugging

### Step 1: Verify the Component Renders

**Test in Browser Console:**

```javascript
// Check if Datepicker element exists
document.querySelector("#start-date");
// Should return the input element

// Check if calendar popup exists (after clicking)
document.querySelector('[role="dialog"]');
// Should return calendar popup when open
```

### Step 2: Check for Common Issues

#### Issue 1: Calendar Renders But Is Hidden

**Symptoms:**

- Input field is visible
- Clicking does nothing
- No errors in console

**Possible Causes:**

- Z-index too low (calendar behind other elements)
- Parent has `overflow: hidden`
- Calendar positioned off-screen

**Fix Applied:**

- Added `position: relative` and `zIndex: 10` to parent divs
- Removed inline z-index from Datepicker (let Flowbite handle it)

#### Issue 2: JavaScript Not Initializing

**Symptoms:**

- Input field visible
- Clicking doesn't trigger calendar
- Console errors about missing functions

**Check:**

- Flowbite JS is loaded (check network tab)
- No conflicting datepicker libraries
- React version compatibility

#### Issue 3: CSS Not Loading

**Symptoms:**

- Input looks unstyled
- Calendar might render but be invisible

**Check:**

- Tailwind CSS is loading
- Flowbite plugin is in tailwind.config.js ✅
- No CSS conflicts

### Step 3: Test Minimal Example

Add this to a page temporarily to test:

```tsx
import { DatePickerTest } from "./components/shared/DatePickerTest";

// In your component:
<DatePickerTest />;
```

If this works → problem is in DateRangeSelector
If this doesn't work → problem is in setup

### Step 4: Browser-Specific Checks

**Chrome/Edge:**

- Check DevTools → Elements → see if calendar div appears on click
- Check Console for errors

**Firefox:**

- Same checks as Chrome

**Safari:**

- May have different z-index behavior

### Step 5: Verify Flowbite Version

Check if there are known issues with 0.12.16:

```bash
npm list flowbite-react
```

Consider updating if issues persist:

```bash
npm update flowbite-react
```

## Current Implementation

✅ Using correct props (`value`, `onChange`)
✅ Using Date objects (not strings)
✅ Added positioning for z-index
✅ Created minimal test component

## Next Steps

1. **Test DatePickerTest component** - Does minimal version work?
2. **Check browser console** - Any errors?
3. **Inspect DOM on click** - Does calendar element appear?
4. **Check CSS** - Is calendar visible but positioned wrong?

## Expert Principle

**Isolate the problem:**

- If minimal test works → our code has the issue
- If minimal test fails → setup/configuration has the issue

This narrows down where to look.
