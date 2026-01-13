package config

import (
	"fmt"
	"os"
	"path/filepath"
)

// Config holds the application configuration
type Config struct {
	DataLakePath     string
	EmailRawPath     string
	EmailProcessingPath string
	EmailFinalPath   string
}

// LoadConfig loads configuration from environment variables
func LoadConfig() (*Config, error) {
	dataLake := os.Getenv("DATA_LAKE")
	if dataLake == "" {
		return nil, fmt.Errorf("DATA_LAKE environment variable not set")
	}

	// Validate data lake path exists
	if _, err := os.Stat(dataLake); os.IsNotExist(err) {
		return nil, fmt.Errorf("data lake path does not exist: %s", dataLake)
	}

	cfg := &Config{
		DataLakePath:        dataLake,
		EmailRawPath:        filepath.Join(dataLake, "unprocessed", "Isaiah.Delemar@sol.doi.gov-olderthan1year.pst"),
		EmailProcessingPath: filepath.Join(dataLake, "unprocessed", "emails_in_process"),
		EmailFinalPath:      filepath.Join(dataLake, "unprocessed", "emails_final"),
	}

	return cfg, nil
}

// GetDataLakePath returns the data lake root path
func (c *Config) GetDataLakePath() string {
	return c.DataLakePath
}

// GetEmailFinalPath returns the path to processed email files
func (c *Config) GetEmailFinalPath() string {
	return c.EmailFinalPath
}
