package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// ProductionRequest represents a production request
type ProductionRequest struct {
	ID          int    `json:"id"`
	Title       string `json:"title,omitempty"` // Generated, so omitempty
	Description string `json:"description"`
}

// App struct
type App struct {
	ctx      context.Context
	manifest map[string]string
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{
		manifest: make(map[string]string),
	}
}

// startup is called when the app starts
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
	// Load manifest will happen lazily when needed
}

// loadManifest loads the JSON manifest file
func (a *App) loadManifest() error {
	if len(a.manifest) > 0 {
		return nil // Already loaded
	}

	// Try database folder first, then frontend
	paths := []string{
		"database/raw_delemar_manifest.json",
		"frontend/src/data/raw_delemar_manifest.json",
	}

	var data []byte
	var err error
	for _, path := range paths {
		data, err = os.ReadFile(path)
		if err == nil {
			break
		}
	}

	if err != nil {
		return fmt.Errorf("could not load manifest: %w", err)
	}

	if err := json.Unmarshal(data, &a.manifest); err != nil {
		return fmt.Errorf("could not parse manifest: %w", err)
	}

	return nil
}

// Greet returns a greeting for the given name
func (a *App) Greet(name string) string {
	return "Hello " + name + ", It's show time!"
}

// GetProductionRequests loads and returns production requests from JSON
func (a *App) GetProductionRequests() ([]ProductionRequest, error) {
	// Try database folder first, then frontend
	paths := []string{
		"database/production_requests.json",
		"frontend/src/data/production_requests.json",
	}

	var data []byte
	var err error
	for _, path := range paths {
		data, err = os.ReadFile(path)
		if err == nil {
			break
		}
	}

	if err != nil {
		return nil, fmt.Errorf("could not load production requests: %w", err)
	}

	var requests []ProductionRequest
	if err := json.Unmarshal(data, &requests); err != nil {
		return nil, fmt.Errorf("could not parse production requests: %w", err)
	}

	// Generate title from ID (ensure it's always set)
	for i := range requests {
		if requests[i].Title == "" {
			requests[i].Title = fmt.Sprintf("REQUEST FOR PRODUCTION NO: %d", requests[i].ID)
		}
	}

	return requests, nil
}

// GetCategories returns available categories with counts based on file paths
func (a *App) GetCategories() (map[string]int, error) {
	if err := a.loadManifest(); err != nil {
		return nil, err
	}

	counts := map[string]int{
		"email":  0,
		"claim":  0,
		"other":  0,
	}

	for path := range a.manifest {
		pathLower := strings.ToLower(path)

		// Skip hidden/system files
		if strings.HasPrefix(path, ".") {
			continue
		}

		// Skip archives
		ext := strings.ToLower(filepath.Ext(path))
		if ext == ".zip" || ext == ".pst" || ext == ".zst" {
			continue
		}

		// Categorize by path patterns
		if strings.Contains(pathLower, "email") || strings.Contains(pathLower, ".eml") || strings.Contains(pathLower, "mail") {
			counts["email"]++
		} else if strings.Contains(pathLower, "claim") {
			counts["claim"]++
		} else {
			counts["other"]++
		}
	}

	return counts, nil
}

// GetFileCount returns total file count after applying default filters
func (a *App) GetFileCount() (int, error) {
	if err := a.loadManifest(); err != nil {
		return 0, err
	}

	count := 0
	for path := range a.manifest {
		// Apply default filters
		if a.shouldExclude(path) {
			continue
		}
		count++
	}

	return count, nil
}

// GetFiles returns filtered files based on categories
func (a *App) GetFiles(categories []string) ([]map[string]interface{}, error) {
	if err := a.loadManifest(); err != nil {
		return nil, err
	}

	var files []map[string]interface{}

	for path, hash := range a.manifest {
		// Apply default filters
		if a.shouldExclude(path) {
			continue
		}

		// Categorize
		category := a.categorizePath(path)

		// Filter by selected categories
		if len(categories) > 0 {
			found := false
			for _, cat := range categories {
				if cat == category {
					found = true
					break
				}
			}
			if !found {
				continue
			}
		}

		files = append(files, map[string]interface{}{
			"path":     path,
			"hash":     hash,
			"category": category,
			"filename": filepath.Base(path),
			"size":     0, // Will add later if needed
		})
	}

	return files, nil
}

// categorizePath determines category from file path
func (a *App) categorizePath(path string) string {
	pathLower := strings.ToLower(path)

	if strings.Contains(pathLower, "email") || strings.Contains(pathLower, ".eml") || strings.Contains(pathLower, "mail") {
		return "email"
	} else if strings.Contains(pathLower, "claim") {
		return "claim"
	}
	return "other"
}

// shouldExclude applies default exclusion rules
func (a *App) shouldExclude(path string) bool {
	// Exclude hidden files
	if strings.HasPrefix(filepath.Base(path), ".") {
		return true
	}

	// Exclude archives
	ext := strings.ToLower(filepath.Ext(path))
	if ext == ".zip" || ext == ".pst" || ext == ".zst" {
		return true
	}

	return false
}
