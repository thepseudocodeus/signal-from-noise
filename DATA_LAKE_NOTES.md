# Data Lake Notes
## TB USB External Drive

**Type:** USB External Drive (TB scale)
**Mount Point:** `/Volumes/DirectFileStore/DELEMAR/`
**Status:** Physical external drive (not cloud)

---

## Important Considerations

### 1. USB Drive Performance
- **USB speed:** Limited by USB interface (USB 3.0 typical)
- **Large datasets:** TB scale means large file counts
- **Performance:** Slower than internal SSD
- **Strategy:** Process in chunks, stream data, don't load all into memory

### 2. Drive Availability
- **Mounting:** Drive must be mounted before app can access
- **Disconnection:** USB drive can be disconnected
- **Error handling:** Must handle cases where drive is not available
- **Validation:** Check drive availability before processing

### 3. File System
- **Format:** Check file system format (APFS, HFS+, exFAT, etc.)
- **Permissions:** Ensure app has read permissions
- **Path stability:** USB mount points may vary (check mount point)

### 4. Demo Strategy
- **Sample data:** For demo, process sample files (first N files)
- **Streaming:** Stream processing (don't load all data)
- **Progress:** Show progress for long operations
- **Cancellation:** Allow user to cancel long operations

---

## Implementation Approach

### Backend Considerations

```go
// Check drive availability
func ValidateDataLake(path string) error {
    // Check if path exists
    // Check if path is readable
    // Check if it's a directory
    // Return error if not available
}

// Discover files efficiently
func DiscoverEmailFiles(path string) ([]string, error) {
    // Walk directory (efficiently)
    // Filter for email files (.parquet, .csv, .json)
    // Return file list
    // Don't load file contents here
}

// Process files in chunks
func ProcessEmailFiles(files []string, query EmailQuery) (*QueryResult, error) {
    // Process files one at a time
    // Stream results
    // Update progress
    // Handle errors gracefully
}
```

### Frontend Considerations

```tsx
// Check drive status
const [driveStatus, setDriveStatus] = useState<'checking' | 'available' | 'unavailable'>('checking');

// Show drive status
if (driveStatus === 'unavailable') {
  return <ErrorMessage message="External drive not available. Please ensure drive is connected." />;
}

// Progress indicator for long operations
<ProgressBar value={progressPercent} label={`Processing ${currentFile} of ${totalFiles}`} />
```

---

## Error Handling

### Drive Not Available
- **Error:** Drive not mounted or disconnected
- **User message:** "External drive not available. Please ensure the drive is connected and mounted."
- **Action:** Check drive status, show error, allow retry

### Drive Slow
- **Warning:** Large dataset, processing may take time
- **User message:** "Processing large dataset. This may take several minutes."
- **Action:** Show progress, allow cancellation

### Permission Issues
- **Error:** No read permission
- **User message:** "Permission denied. Please check file permissions."
- **Action:** Check permissions, show error

---

## Demo Optimization

### For Demo Performance

1. **Sample Files**
   - Process first 10-50 files (not all)
   - Show "Processing sample" indicator
   - Full processing can come later

2. **Sample Results**
   - Return first 100-1000 records
   - Show "Showing first N results" message
   - Full results can be exported/downloaded

3. **Caching**
   - Cache file list (don't re-scan on every query)
   - Cache file metadata
   - Invalidate cache when needed

4. **Progress Feedback**
   - Show file count discovered
   - Show files processed
   - Show estimated time remaining

---

## Testing Strategy

### 1. Drive Availability
```bash
# Check if drive is mounted
ls /Volumes/DirectFileStore/DELEMAR/

# Check permissions
ls -la /Volumes/DirectFileStore/DELEMAR/unprocessed/emails_final/
```

### 2. File Discovery
```go
// Test file discovery
files, err := datalake.DiscoverEmailFiles(config.EmailFinalPath)
fmt.Printf("Found %d email files\n", len(files))
```

### 3. Performance Testing
```go
// Test processing time
start := time.Now()
results, err := ProcessEmailFiles(files[:10], query) // First 10 files
duration := time.Since(start)
fmt.Printf("Processed 10 files in %v\n", duration)
```

---

## Recommendations

### For Demo
1. **Sample Processing:** Process first 10-50 files for demo
2. **Sample Results:** Show first 100-1000 records
3. **Progress Indicators:** Show progress clearly
4. **Error Handling:** Handle drive unavailability gracefully

### For Production
1. **Streaming:** Process files in streams
2. **Parallel Processing:** Process multiple files concurrently (if USB supports)
3. **Caching:** Cache file lists and metadata
4. **Resume:** Allow resuming interrupted queries
5. **Background Processing:** Process in background thread

---

**Key Takeaway:** USB drive means we need to be mindful of performance and handle drive availability gracefully. For demo, sampling is key.
