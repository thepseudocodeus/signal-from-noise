# Component Usage Policy: Making Assumptions Explicit

## The Assumption

**We use Flowbite React default components for all UI elements.**

This assumption must be:
1. **Explicit** - Documented and clear
2. **Verifiable** - Can be checked automatically
3. **Enforced** - Violations are caught early

## Why This Matters

### Expert Principle: Make Assumptions Explicit

**Before (Implicit):**
- "We probably use Flowbite components"
- "I think this is from Flowbite"
- "It should work like Flowbite"

**After (Explicit):**
- "All UI components come from flowbite-react" (documented)
- "Component audit verifies this" (verifiable)
- "Violations fail the build" (enforced)

## Policy Rules

### ✅ Allowed

1. **Flowbite React Components**
   - Import from `flowbite-react`
   - Use default styling
   - Use default props

2. **Custom Components for App Logic**
   - Components that don't have Flowbite equivalents
   - App-specific business logic components
   - Must be documented in `COMPONENT_AUDIT.md`

3. **Layout/Styling**
   - Tailwind utility classes for layout
   - Custom CSS for app-specific needs
   - Flowbite theme customization (if needed)

### ❌ Not Allowed

1. **Other UI Libraries**
   - No Material-UI
   - No Ant Design
   - No React Bootstrap
   - No custom replacements for Flowbite components

2. **Custom Flowbite Component Replacements**
   - Don't create custom Button if Flowbite Button exists
   - Don't create custom Card if Flowbite Card exists
   - Use Flowbite components, customize with props/styling

## Verification

### Manual Audit
See `COMPONENT_AUDIT.md` for current component registry.

### Automated Verification
Run the component audit script:
```bash
npx tsx src/components/verify-flowbite-usage.ts
```

### CI/CD Integration
Add to build process to fail if violations are found.

## How to Add a New Component

1. **Check Flowbite First**
   - Does Flowbite have this component?
   - If yes, use Flowbite component
   - If no, proceed to step 2

2. **Create Custom Component**
   - Document why it's custom in `COMPONENT_AUDIT.md`
   - Add to `ALLOWED_CUSTOM_COMPONENTS` in audit script
   - Use Flowbite styling patterns

3. **Verify**
   - Run audit script
   - Ensure no violations

## Benefits

1. **Consistency**: All components follow same design system
2. **Maintainability**: One source of truth (Flowbite)
3. **Predictability**: Components behave as documented
4. **Testing**: Can verify assumption automatically
5. **Onboarding**: New developers know what to use

## Expert Thinking Applied

This policy makes the assumption explicit and verifiable:

- **Assumption**: "We use Flowbite defaults"
- **Test**: Component audit script
- **Enforcement**: Fails build if violated
- **Documentation**: Clear policy and registry

Just like our mode system and assertions, this makes the "physics" of our component usage explicit and testable.
