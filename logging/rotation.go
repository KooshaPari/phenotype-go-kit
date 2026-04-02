package logging

import (
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// LogRotationConfig holds configuration for log rotation.
type LogRotationConfig struct {
	Directory  string `yaml:"directory" json:"directory"`
	MaxSizeMB  int    `yaml:"max_size_mb" json:"max_size_mb"`
	MaxAgeDays int    `yaml:"max_age_days" json:"max_age_days"`
	MaxBackups int    `yaml:"max_backups" json:"max_backups"`
	Compress   bool   `yaml:"compress" json:"compress"`
	Pattern    string `yaml:"pattern" json:"pattern"` // e.g., "app-*.log"
}

// RotationPolicy defines when to rotate logs.
type RotationPolicy struct {
	MaxSizeBytes int64
	MaxAgeDays   int
	Schedule     string // "daily", "hourly", or cron expression
}

// DefaultRotationPolicy returns the default rotation policy.
func DefaultRotationPolicy() RotationPolicy {
	return RotationPolicy{
		MaxSizeBytes: 100 * 1024 * 1024, // 100MB
		MaxAgeDays:   30,
		Schedule:     "daily",
	}
}

// RetentionPolicy defines log retention rules.
type RetentionPolicy struct {
	KeepDays         int      `yaml:"keep_days" json:"keep_days"`
	KeepBackups      int      `yaml:"keep_backups" json:"keep_backups"`
	CompressOldLogs  bool     `yaml:"compress_old_logs" json:"compress_old_logs"`
	ExcludedPatterns []string `yaml:"excluded_patterns" json:"excluded_patterns"` // e.g., ["debug.log"]
}

// LogManager handles log file management.
type LogManager struct {
	config          LogRotationConfig
	retentionPolicy RetentionPolicy
}

// NewLogManager creates a new log manager.
func NewLogManager(cfg LogRotationConfig, retention RetentionPolicy) *LogManager {
	return &LogManager{
		config:          cfg,
		retentionPolicy: retention,
	}
}

// Cleanup removes old log files based on retention policy.
func (m *LogManager) Cleanup() error {
	if m.config.Directory == "" {
		return nil
	}

	entries, err := os.ReadDir(m.config.Directory)
	if err != nil {
		return fmt.Errorf("failed to read log directory: %w", err)
	}

	cutoff := time.Now().AddDate(0, 0, -m.retentionPolicy.KeepDays)
	backupCount := make(map[string]int)

	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}

		// Check age
		info, err := entry.Info()
		if err != nil {
			continue
		}

		// Clean up old backups
		name := info.Name()
		datePart := extractDate(name)
		if datePart != "" {
			backupCount[datePart]++
			if backupCount[datePart] > m.retentionPolicy.KeepBackups {
				os.RemoveAll(filepath.Join(m.config.Directory, name))
			}
		}

		// Compress old logs if needed
		if m.retentionPolicy.CompressOldLogs && info.ModTime().Before(cutoff) {
			if filepath.Ext(name) == ".log" {
				m.compressLogFile(filepath.Join(m.config.Directory, name))
			}
		}
	}

	return nil
}

// compressLogFile compresses a log file using gzip.
func (m *LogManager) compressLogFile(path string) error {
	// In production, implement actual compression
	// For now, this is a stub
	return nil
}

func extractDate(name string) string {
	// Extract date portion from filename like "app-2026-03-24.log"
	// This is a simplified implementation
	return ""
}

// LogFileInfo represents information about a log file.
type LogFileInfo struct {
	Name         string
	Path         string
	SizeBytes    int64
	ModTime      time.Time
	IsCompressed bool
}

// ListLogFiles returns all log files in the directory.
func (m *LogManager) ListLogFiles() ([]LogFileInfo, error) {
	if m.config.Directory == "" {
		return nil, nil
	}

	entries, err := os.ReadDir(m.config.Directory)
	if err != nil {
		return nil, err
	}

	var files []LogFileInfo
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		info, err := entry.Info()
		if err != nil {
			continue
		}

		files = append(files, LogFileInfo{
			Name:         entry.Name(),
			Path:         filepath.Join(m.config.Directory, entry.Name()),
			SizeBytes:    info.Size(),
			ModTime:      info.ModTime(),
			IsCompressed: filepath.Ext(entry.Name()) == ".gz",
		})
	}

	return files, nil
}

// GetTotalSize returns total size of all log files.
func (m *LogManager) GetTotalSize() (int64, error) {
	files, err := m.ListLogFiles()
	if err != nil {
		return 0, err
	}

	var total int64
	for _, f := range files {
		total += f.SizeBytes
	}

	return total, nil
}

// RotationConfigToYAML generates a sample log rotation config.
func RotationConfigToYAML() string {
	return `# Log Rotation Configuration
# Example for use with logrotate or similar

# Rotate logs daily
daily

# Keep 30 days of logs
rotate 30

# Maximum log file size 100MB
maxsize 100M

# Compress old logs
compress
delaycompress

# Don't error if log file is missing
missingok

# Don't rotate empty logs
notifempty

# Create new log files with these permissions
create 0644 root root

# Sample pattern for application logs
# /var/log/phenotype/*.log {
#     daily
#     rotate 30
#     maxsize 100M
#     compress
#     delaycompress
#     missingok
#     notifempty
# }
`
}
