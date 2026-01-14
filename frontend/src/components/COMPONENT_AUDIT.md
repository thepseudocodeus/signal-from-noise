# Component Audit: Flowbite Default Components

## Assumption
**We use Flowbite React default components for all UI elements, with only app-specific content/logic as custom code.**

## Component Registry

### Flowbite Components Used ✓

| Component | Source | Used In | Status |
|-----------|--------|---------|--------|
| `Button` | flowbite-react | All step components | ✓ Default |
| `Card` | flowbite-react | ProductionRequestStep, ResultsStep, SplashStep | ✓ Default |
| `Checkbox` | flowbite-react | CategoryStep, FileSelectionStep | ✓ Default |
| `Label` | flowbite-react | DateRangeStep, CategoryStep | ✓ Default |
| `Datepicker` | flowbite-react | DateRangeStep | ✓ Default |
| `Table` | flowbite-react | FileSelectionStep | ✓ Default |
| `Badge` | flowbite-react | FileSelectionStep, ResultsStep | ✓ Default |
| `Spinner` | flowbite-react | FileSelectionStep | ✓ Default |
| `Alert` | flowbite-react | FileSelectionStep | ✓ Default |
| `TextInput` | flowbite-react | KeywordsStep | ✓ Default |
| `Select` | flowbite-react | YearRangeStep | ✓ Default |

### Custom Components

| Component | Location | Reason | Flowbite Alternative? |
|-----------|----------|--------|----------------------|
| `ProgressIndicator` | shared/ProgressIndicator.tsx | Custom progress bar | Could use Flowbite Progress, but custom is simpler |

### Custom Layout/Styling

- Layout divs with Tailwind classes (standard practice)
- Custom styling for app-specific needs (acceptable)

## Verification Rules

1. **All UI components must come from `flowbite-react`**
2. **No custom component replacements for Flowbite components**
3. **Custom components only for app-specific logic (like ProgressIndicator)**
4. **All Flowbite components use default theme (no custom overrides)**

## How to Verify

Run the component audit script to verify this assumption.
