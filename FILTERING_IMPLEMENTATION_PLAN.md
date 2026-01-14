# Filtering Implementation Plan: Expert Approach

## Goal
Add Topic, People, and Sentiment filtering after date range selection to incrementally reduce the problem space.

## Expert Principles Applied

1. **Incremental Complexity Reduction**: Each filter reduces the universe of possible files
2. **Explicit Assumptions**: Assert what must be true at each step
3. **Incremental Query Building**: Build SQL query as selections are made
4. **Precomputed Sets**: Internal/external email sets computed once
5. **Make it Work → Make it Right → Make it Fast**

## Implementation Steps

### Step 1: Database Schema Extension (Make it Work)
- Add email fields: `subject`, `from_email`, `to_email`, `sentiment`, `is_internal`
- Update File struct to include these fields
- Migration-friendly: new columns, existing data gets defaults

### Step 2: Mock Data Enhancement (Make it Work)
- Generate email subjects with topics
- Generate internal/external email addresses
- Assign sentiment values
- Precompute internal vs external sets

### Step 3: Query Builder (Make it Right)
- Incremental SQL query construction
- Each filter adds WHERE clauses
- Transparent logging of query building
- Assertions verify query validity

### Step 4: Backend Methods (Make it Right)
- `GetTopics()` - Extract unique topics from email subjects
- `GetPeople()` - Get internal/external email lists
- `GetSentimentOptions()` - Return available sentiments
- All with transparent logging

### Step 5: Frontend Components (Make it Work)
- `TopicSelectionStep` - Multi-select from topics
- `PeopleSelectionStep` - Select internal/external/individuals
- `SentimentSelectionStep` - Radio buttons for sentiment
- All using Flowbite components

### Step 6: Flow Integration (Make it Right)
- Add steps after date-range
- Update QueryState to track new filters
- Incremental query building as selections made

### Step 7: Query Execution (Make it Fast)
- Optimize SQL with proper indexes
- Use precomputed sets where possible
- Log query performance

## Assumptions to Assert

1. **Email subjects contain topic information** - Can extract topics
2. **FROM/TO fields identify people** - Can classify internal/external
3. **Sentiment can be determined** - Available in data or computed
4. **Each filter reduces result set** - Query gets more specific
5. **SQL query is valid** - Each addition maintains query correctness

## Success Criteria

- ✅ User can select topics from email subjects
- ✅ User can filter by internal/external people
- ✅ User can filter by sentiment
- ✅ Each selection reduces file count
- ✅ Query builds incrementally
- ✅ All operations logged transparently
- ✅ Assumptions asserted at each step
