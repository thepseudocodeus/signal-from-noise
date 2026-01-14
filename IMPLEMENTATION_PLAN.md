# Implementation Plan: Working Production Request Application

**Goal:** Create a working application where users select criteria, view paginated files, select/deselect files, and create a zip file named by production request.

**Principle:** Make it work → Make it right → Make it fast

---

## Step 1: Database Setup (Make it Work)
**Goal:** Create SQLite database with mock data representing real-world issues

**Tasks:**
1. Create Go database package (`database/`)
2. Initialize SQLite database with schema:
   - `files` table: id, path, directory, category (email/claim/other), date, size, privileged (boolean), duplicate_hash
   - `production_requests` table: id, title, description
3. Create mock data generator:
   - Duplicate files (same content, different paths)
   - Files in directories with "Claim" in name
   - Email files with privileged flag (client-lawyer communications)
   - Files across different date ranges
   - Mix of categories
4. Add Go methods to query database:
   - `GetFiles(filters)` - with date range, category, exclude privileged
   - `GetFileCount(filters)` - for pagination
   - `GetFileByID(id)` - for individual file access

**Validation:** Database initializes, contains mock data, queries return results

---

## Step 2: Backend API Methods (Make it Work)
**Goal:** Expose database queries through Wails bindings

**Tasks:**
1. Add methods to `app.go`:
   - `SearchFiles(productionRequestID, dateStart, dateEnd, categories, excludePrivileged)` → returns file list with pagination
   - `GetFileCount(filters)` → returns count
   - `CreateZipFile(productionRequestID, fileIDs)` → creates zip, returns path
2. Add transparent logging for all operations
3. Return clear error messages

**Validation:** Methods are callable from frontend, return expected data

---

## Step 3: Update Frontend Flow (Make it Work)
**Goal:** Integrate category and date range steps, add file selection step

**Tasks:**
1. Update `App.tsx` to include:
   - Category step (use existing `CategoryStep`)
   - Date range step (use existing `DateRangeStep` with Flowbite calendar)
   - New `FileSelectionStep` component for paginated file list
2. Update `QueryState` type to include:
   - `dateRange: DateRange` (already exists)
   - `categories: DataCategory[]` (already exists)
3. Flow: Production Request → Category → Date Range → File Selection → Zip Creation

**Validation:** User can navigate through all steps, selections are preserved

---

## Step 4: File Selection Component (Make it Work)
**Goal:** Paginated file list with selection capabilities

**Tasks:**
1. Create `FileSelectionStep.tsx`:
   - Call backend `SearchFiles` with filters
   - Display paginated table using Flowbite Table component
   - Checkboxes for each file
   - "Select All" / "Deselect All" buttons
   - Remove individual items button
   - Pagination controls
   - Show file count, total size
2. State management:
   - Track selected file IDs
   - Track current page
   - Load files on mount and filter changes

**Validation:** Files display, selection works, pagination works

---

## Step 5: Zip Creation (Make it Work)
**Goal:** Create zip file from selected files

**Tasks:**
1. Add "Create Zip" button to `FileSelectionStep`
2. Call backend `CreateZipFile` with production request ID and selected file IDs
3. Show progress indicator during zip creation
4. Display success message with zip file path
5. Backend implementation:
   - Read files from database
   - Create zip file in temp directory
   - Name: `{productionRequestID}_{timestamp}.zip`
   - Return path to frontend

**Validation:** Zip file is created, contains selected files, named correctly

---

## Step 6: Theme Revert (Make it Right)
**Goal:** Use Flowbite default theme

**Tasks:**
1. Remove custom theme overrides
2. Use Flowbite default styling
3. Ensure components use default Flowbite appearance

**Validation:** UI uses Flowbite default theme

---

## Step 7: Error Handling & Transparency (Make it Right)
**Goal:** Clear error messages and transparent backend operations

**Tasks:**
1. Add error handling to all frontend API calls
2. Display user-friendly error messages
3. Add backend logging for all operations:
   - Log query parameters
   - Log file counts
   - Log zip creation progress
   - Log errors with context
4. Add loading states to UI

**Validation:** Errors are clear, operations are logged, user sees progress

---

## Step 8: Testing & Refinement (Make it Fast)
**Goal:** Ensure everything works end-to-end

**Tasks:**
1. Test complete flow:
   - Select production request
   - Select categories
   - Select date range
   - View files
   - Select files
   - Create zip
2. Test edge cases:
   - No files found
   - All files selected
   - Empty selection
   - Large file lists
3. Performance check:
   - Pagination works smoothly
   - Zip creation doesn't block UI

**Validation:** Complete flow works, edge cases handled, performance acceptable

---

## File Structure Changes

### New Files:
- `database/database.go` - Database initialization and queries
- `database/schema.sql` - Database schema
- `database/mockdata.go` - Mock data generation
- `frontend/src/components/steps/FileSelectionStep.tsx` - File selection UI

### Modified Files:
- `app.go` - Add new API methods
- `frontend/src/App.tsx` - Update flow to include all steps
- `frontend/src/types/query.ts` - Ensure types are complete
- `go.mod` - Add SQLite driver dependency

---

## Dependencies

### Go:
- `github.com/mattn/go-sqlite3` - SQLite driver
- `archive/zip` - Zip file creation (standard library)
- `path/filepath` - File path handling (standard library)

### Frontend:
- Flowbite React components (already installed)
- Existing React Router setup

---

## Success Criteria

1. ✅ User can select production request
2. ✅ User can select categories (Email, Claims, Other)
3. ✅ User can select date range using calendar
4. ✅ User sees paginated list of files matching criteria
5. ✅ User can select/deselect individual files or all files
6. ✅ User can remove selected files
7. ✅ User can create zip file named by production request
8. ✅ Backend operations are logged transparently
9. ✅ Application uses Flowbite default theme
10. ✅ All operations work end-to-end

---

## Next Steps After This Works

- Optimize database queries
- Add caching
- Improve UI/UX
- Add more filters
- Performance optimizations
