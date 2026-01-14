package app

import (
	"fmt"
)

// Mode represents the operational mode of the application
// This makes the system's "physics" explicit - what mode determines what operations are valid
type Mode string

const (
	// DatabaseMode indicates the app is using SQLite database for all operations
	// In this mode, file operations use the database, not the data lake
	DatabaseMode Mode = "database"

	// DataLakeMode indicates the app is using the data lake for file operations
	// In this mode, file operations access the data lake directly
	DataLakeMode Mode = "hybrid" // "hybrid" because we still use database for queries, but data lake for file access
)

// String returns the string representation of the mode
func (m Mode) String() string {
	return string(m)
}

// IsValid checks if the mode is a valid mode
func (m Mode) IsValid() bool {
	return m == DatabaseMode || m == DataLakeMode
}

// ModeInfo provides information about what a mode enables
type ModeInfo struct {
	Mode        Mode
	Description string
	EnabledFeatures []string
}

// GetModeInfo returns information about what the mode enables
func GetModeInfo(mode Mode) ModeInfo {
	switch mode {
	case DatabaseMode:
		return ModeInfo{
			Mode:        DatabaseMode,
			Description: "Database-only mode: All file operations use SQLite database",
			EnabledFeatures: []string{
				"File search and filtering via database",
				"Production request zip creation",
				"Category and date range filtering",
				"Privileged file exclusion",
			},
		}
	case DataLakeMode:
		return ModeInfo{
			Mode:        DataLakeMode,
			Description: "Hybrid mode: Database for queries, data lake for file access",
			EnabledFeatures: []string{
				"All database mode features",
				"Direct data lake file access",
				"Email file counting from data lake",
				"Data lake status checking",
			},
		}
	default:
		return ModeInfo{
			Mode:        mode,
			Description: "Unknown mode",
			EnabledFeatures: []string{},
		}
	}
}

// ValidateModeTransition checks if transitioning from one mode to another is valid
// This enforces the "physics" - what mode transitions make sense
func ValidateModeTransition(from Mode, to Mode) error {
	// All transitions are currently valid
	// In the future, we might restrict certain transitions
	if !to.IsValid() {
		return fmt.Errorf("invalid target mode: %s", to)
	}
	return nil
}
