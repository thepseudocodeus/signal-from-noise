# Operational Mode Design: Expert Thinking

## The Problem We Solved

Instead of hiding configuration failures or making them "optional," we made the system's operational state **explicit** through mode declaration. This follows expert software engineering principles.

## Key Principles Applied

### 1. **Explicit State Declaration**

**What experts do:** Make system state explicit, not implicit.

**Before:**

- Config fails silently → unclear what mode we're in
- Operations might work or might not → unpredictable
- No way to know what features are available

**After:**

- Mode is explicitly declared: `DatabaseMode` or `DataLakeMode`
- Mode determines what operations are valid
- System state is queryable: `GetOperationalMode()`

**Why this matters:**

- **Predictability**: You always know what mode you're in
- **Debugging**: Problems are immediately clear (wrong mode for operation)
- **Documentation**: Mode tells you what the system can do

### 2. **Design by Contract with Assertions**

**What experts do:** Assertions test assumptions about what must be true.

**Example:**

```go
// ASSUMPTION: Data lake operations require DataLakeMode
// If we're in DatabaseMode, this operation is invalid
assert.That(a.mode == app.DataLakeMode,
    "GetEmailFileCount requires DataLakeMode (operation-mode consistency check)")
```

**What this teaches:**

- **Fail Fast**: If mode doesn't match operation, fail immediately with clear message
- **Self-Documenting**: Assertion comment explains the assumption
- **Hypothesis Testing**: Running the app tests the hypothesis "mode matches operation"

**The Physics:**

- Each operation has a **mode requirement**
- Assertions verify the **mode matches the operation intent**
- If assertion fails, our hypothesis about the system state is wrong

### 3. **Mode-Based Feature Gating**

**What experts do:** Features are explicitly enabled/disabled based on mode.

**Implementation:**

```go
type ModeInfo struct {
    Mode            Mode
    Description     string
    EnabledFeatures []string
}
```

**Why this matters:**

- **Transparency**: You know exactly what each mode enables
- **User Communication**: Can tell users what features are available
- **Testing**: Can verify mode enables expected features

### 4. **State Transitions Are Logged**

**What experts do:** Log state transitions explicitly.

**Example:**

```go
logging.LogTransition("unknown", string(app.DatabaseMode),
    "data lake configuration unavailable")
```

**Why this matters:**

- **Audit Trail**: You can see how the system got into its current state
- **Debugging**: Understand why mode was chosen
- **Transparency**: Makes decision-making process visible

### 5. **Mode Consistency Checks**

**What experts do:** Verify that mode matches system resources.

**Example:**

```go
// ASSUMPTION: In DataLakeMode, data lake service must be initialized
// If service is nil, mode determination logic was incorrect
assert.ThatNotNil(a.dataLakeService,
    "data lake service must be initialized in DataLakeMode (mode consistency check)")
```

**What this teaches:**

- **Invariant Checking**: Mode and resources must be consistent
- **Fail Fast**: If inconsistent, fail immediately (don't fail later mysteriously)
- **Self-Correcting**: System detects its own inconsistencies

## The Expert Mindset

### Thinking in Terms of "Physics"

Experts think: "What are the invariant laws that govern this system?"

**For our system:**

1. **Mode determines available operations** (law of operations)
2. **Mode must match resources** (law of consistency)
3. **Operations must match mode** (law of validity)

**Assertions test these laws:**

- If law is violated → assertion fails → hypothesis is wrong
- Running the app = experiment that tests these laws

### Making Assumptions Explicit

**Before (implicit):**

- "Maybe config works, maybe it doesn't"
- "Maybe this operation works, maybe it doesn't"

**After (explicit):**

- "We are in DatabaseMode" (explicit state)
- "This operation requires DataLakeMode" (explicit requirement)
- Assertion verifies: mode matches requirement (explicit check)

### Fail Fast, Fail Clearly

**Expert principle:** Fail immediately with a clear message, not later with a mysterious error.

**Our implementation:**

- Assertion fails → panic with clear message
- Message explains: what assumption failed, where, why it matters
- No silent failures, no mysterious errors later

## How to Use This Pattern

### 1. Declare Modes Explicitly

```go
type Mode string
const (
    DatabaseMode Mode = "database"
    DataLakeMode Mode = "hybrid"
)
```

### 2. Determine Mode Based on Resources

```go
if hasDataLake {
    a.mode = app.DataLakeMode
} else {
    a.mode = app.DatabaseMode
}
```

### 3. Assert Mode Requirements in Operations

```go
// ASSUMPTION: [What mode is required]
assert.That(a.mode == app.DataLakeMode,
    "[operation] requires DataLakeMode (operation-mode consistency check)")
```

### 4. Log Mode Transitions

```go
logging.LogTransition(from, to, reason)
```

### 5. Verify Mode Consistency

```go
// ASSUMPTION: Mode and resources must match
assert.ThatNotNil(a.dataLakeService,
    "data lake service must exist in DataLakeMode")
```

## Benefits

1. **Predictability**: Always know what mode you're in
2. **Debugging**: Clear errors when mode/operation mismatch
3. **Documentation**: Mode tells you what's available
4. **Testing**: Running app tests mode assumptions
5. **Transparency**: Logs show mode decisions and transitions

## What You Learn

1. **Explicit > Implicit**: Make state explicit, not hidden
2. **Assertions Test Physics**: Each assertion tests a system law
3. **Fail Fast**: Detect problems immediately, not later
4. **State Transitions Matter**: Log how system changes state
5. **Consistency is Key**: Mode and resources must match

This is how experts think: make the system's "physics" explicit, test assumptions with assertions, and fail fast with clear messages.
