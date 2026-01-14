# Logging and Assertions Implementation

## Philosophy

Following NASA and expert software engineering practices, this implementation makes the "physics" of the system explicit through:

1. **Assertions** - Test assumptions that must be true for functions to work
2. **Transparent Logging** - Make operations visible to developers
3. **Invariant Checking** - Verify properties that must always hold
4. **Hypothesis Testing** - Each function tests assumptions about what should happen

## How It Works

### Assertions (`assert/assert.go`)

Assertions test assumptions that must be true for the function to work correctly. If an assumption fails, the program panics with a clear message explaining what assumption was violated.

**Example:**

```go
// ASSUMPTION: Database connection must exist for queries to execute
// If db is nil, the database was not properly initialized
assert.ThatNotNil(d.db, "database connection must exist for queries to execute")
```

**When to use:**

- Preconditions that must be true before function execution
- Postconditions that must be true after function execution
- Invariants that must always hold
- Data integrity checks

### Logging (`logging/logger.go`)

Logging makes operations transparent by documenting:

- What operation is starting/ending
- What assumptions are being tested
- What state the system is in
- What transitions occur
- What errors happen and why

**Example:**

```go
op := logging.StartOperation("SearchFiles", map[string]interface{}{
    "date_start": filters.DateStart,
    "categories": filters.Categories,
})
defer op.EndOperation()
```

**Logging Functions:**

- `StartOperation()` - Begin logging an operation
- `EndOperation()` - End operation with duration
- `LogAssumption()` - Document an assumption being tested
- `LogInvariant()` - Check an invariant
- `LogState()` - Log current system state
- `LogError()` - Log errors with context
- `LogQuery()` - Log database queries
- `LogResult()` - Log query results
- `LogCheckpoint()` - Log execution checkpoint

## Implementation Pattern

Each function follows this pattern:

1. **Start Operation Logging**

   ```go
   op := logging.StartOperation("FunctionName", context)
   defer op.EndOperation()
   ```

2. **Test Preconditions (Assertions)**

   ```go
   // ASSUMPTION: [What must be true]
   // [Why it must be true]
   assert.That(condition, "assumption description")
   ```

3. **Log Assumptions Being Tested**

   ```go
   logging.LogAssumption("assumption description", context)
   ```

4. **Log State/Checkpoints**

   ```go
   logging.LogCheckpoint("location", data)
   ```

5. **Log Errors with Context**

   ```go
   logging.LogError("operation", err, context)
   ```

6. **Test Postconditions (Assertions)**
   ```go
   // ASSUMPTION: [What must be true after execution]
   assert.That(result != nil, "result must be non-nil")
   ```

## What Assumptions Are Tested

### Database Operations

1. **Connection Validity**

   - Database connection exists and is valid
   - Database can be pinged

2. **Query Validity**

   - SQL queries are well-formed
   - Parameters match query structure
   - Results match expected structure

3. **Data Integrity**

   - File IDs are positive
   - Counts are non-negative
   - Date ranges are valid (start <= end)

4. **Pagination**
   - Page numbers are positive
   - Page sizes are positive
   - Offsets are non-negative
   - Returned items don't exceed page size

### File Operations

1. **File Selection**

   - At least one file selected
   - All file IDs exist in database
   - File paths are non-empty

2. **Zip Creation**
   - Output directory is writable
   - Zip file can be created
   - All files can be added to zip
   - Manifest can be written

### API Operations

1. **Request Validity**

   - Production request ID is non-empty
   - Date formats are valid
   - Date ranges are logical
   - Categories are valid

2. **State Consistency**
   - Database is initialized
   - Results match requests
   - Transformations preserve data

## Benefits

1. **Fail Fast** - Problems are detected immediately when assumptions are violated
2. **Clear Errors** - Assertion failures explain exactly what assumption was wrong
3. **Transparent Operations** - Logging shows what the system is doing
4. **Documentation** - Assertions document what must be true for code to work
5. **Testing** - Running the app tests all assumptions (falsifies hypotheses)

## Example Output

When running the application, you'll see logs like:

```
▶️  START: SearchFiles
   date_start: 2022-01-01T00:00:00Z
   categories: [email claim]
   page: 1
🔬 TESTING ASSUMPTION: Request contains valid filter parameters
   has_date_start: true
   category_count: 2
🔒 CHECKING INVARIANT: pagination_parameters
   page: 1
   page_size: 50
📝 QUERY: SELECT COUNT(*) FROM files WHERE date >= ? AND category IN (?,?)
   args_count: 3
📦 RESULT [SearchFiles]: 150 items
📍 CHECKPOINT: SearchFiles
   total_count: 150
   current_page: 1
   offset: 0
✅ END: SearchFiles (duration: 45ms)
   Result: map[files_returned:50 total_count:150 page:1 total_pages:3]
```

If an assumption fails:

```
❌ ASSERTION FAILED: database connection must exist for queries to execute
   Location: app.go:123 in SearchFiles()
panic: ASSERTION FAILED: database connection must exist for queries to execute (at app.go:123)
```

## Running as an Experiment

As you mentioned, running the application is the experiment that falsifies hypotheses about the "invariant physics" of the problems. Each assertion tests a hypothesis:

- **Hypothesis**: "Database connection exists when SearchFiles is called"
- **Test**: `assert.ThatNotNil(a.db, ...)`
- **Falsification**: If assertion fails, hypothesis is false - database wasn't initialized

This makes the code self-documenting and self-testing. Every function explicitly states what must be true for it to work.
