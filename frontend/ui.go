package frontend

import (
	"fmt"
	"html"
	"strings"
	"time"
)

// EscapeHTML escapes HTML content.
func EscapeHTML(s string) string {
	return html.EscapeString(s)
}

// SafeHTML marks HTML as safe.
type SafeHTML string

func (s SafeHTML) String() string {
	return string(s)
}

// FormatDate formats a date.
func FormatDate(t time.Time, layout string) string {
	if layout == "" {
		layout = "Jan 2, 2006"
	}
	return t.Format(layout)
}

// FormatRelativeTime returns relative time string.
func FormatRelativeTime(t time.Time) string {
	now := time.Now()
	diff := now.Sub(t)

	if diff < time.Minute {
		return "just now"
	}

	if diff < time.Hour {
		mins := int(diff.Minutes())
		return fmt.Sprintf("%d minute%s ago", mins, suffix(mins))
	}

	if diff < 24*time.Hour {
		hours := int(diff.Hours())
		return fmt.Sprintf("%d hour%s ago", hours, suffix(hours))
	}

	if diff < 30*24*time.Hour {
		days := int(diff.Hours() / 24)
		return fmt.Sprintf("%d day%s ago", days, suffix(days))
	}

	if diff < 365*24*time.Hour {
		months := int(diff.Hours() / 24 / 30)
		return fmt.Sprintf("%d month%s ago", months, suffix(months))
	}

	years := int(diff.Hours() / 24 / 365)
	return fmt.Sprintf("%d year%s ago", years, suffix(years))
}

func suffix(n int) string {
	if n == 1 {
		return ""
	}
	return "s"
}

// Truncate truncates a string.
func Truncate(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen-3] + "..."
}

// Slugify converts a string to a URL-friendly slug.
func Slugify(s string) string {
	s = strings.ToLower(s)
	s = strings.ReplaceAll(s, " ", "-")

	var result strings.Builder
	for _, r := range s {
		if (r >= 'a' && r <= 'z') || (r >= '0' && r <= '9') || r == '-' {
			result.WriteRune(r)
		}
	}

	return strings.Trim(result.String(), "-")
}

// ClassNames joins CSS class names.
func ClassNames(classes ...string) string {
	var result []string
	for _, c := range classes {
		if c != "" {
			result = append(result, c)
		}
	}
	return strings.Join(result, " ")
}

// StyleFromMap converts a map to CSS style string.
func StyleFromMap(styles map[string]string) string {
	if len(styles) == 0 {
		return ""
	}

	var result []string
	for k, v := range styles {
		result = append(result, fmt.Sprintf("%s: %s", k, v))
	}

	return strings.Join(result, "; ")
}

// FormatBytes formats bytes to human readable.
func FormatBytes(bytes int64) string {
	const unit = 1024
	if bytes < unit {
		return fmt.Sprintf("%d B", bytes)
	}

	div, exp := int64(unit), 0
	for n := bytes / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}

	return fmt.Sprintf("%.1f %cB", float64(bytes)/float64(div), "KMGTPE"[exp])
}

// Pluralize returns singular or plural form.
func Pluralize(count int, singular, plural string) string {
	if count == 1 {
		return singular
	}
	return plural
}

// FormatNumber formats a number with separators.
func FormatNumber(n int) string {
	s := fmt.Sprintf("%d", n)
	var result strings.Builder

	for i, r := range s {
		if i > 0 && (len(s)-i)%3 == 0 {
			result.WriteRune(',')
		}
		result.WriteRune(r)
	}

	return result.String()
}

// QueryBuilder builds URL query strings.
type QueryBuilder struct {
	params map[string][]string
}

// NewQueryBuilder creates a new query builder.
func NewQueryBuilder() *QueryBuilder {
	return &QueryBuilder{
		params: make(map[string][]string),
	}
}

// Add adds a query parameter.
func (qb *QueryBuilder) Add(key, value string) {
	qb.params[key] = append(qb.params[key], value)
}

// Build returns the query string.
func (qb *QueryBuilder) Build() string {
	if len(qb.params) == 0 {
		return ""
	}

	var pairs []string
	for k, vs := range qb.params {
		for _, v := range vs {
			pairs = append(pairs, fmt.Sprintf("%s=%s", k, v))
		}
	}

	return strings.Join(pairs, "&")
}

// HasPrefix checks string prefix.
func HasPrefix(s, prefix string) bool {
	return strings.HasPrefix(s, prefix)
}

// HasSuffix checks string suffix.
func HasSuffix(s, suffix string) bool {
	return strings.HasSuffix(s, suffix)
}

// Contains checks substring.
func Contains(s, substr string) bool {
	return strings.Contains(s, substr)
}

// ToCamelCase converts string to camelCase.
func ToCamelCase(s string) string {
	parts := strings.Split(s, "_")

	var result string
	for i, part := range parts {
		if i == 0 {
			result += strings.ToLower(part)
		} else {
			result += strings.Title(strings.ToLower(part))
		}
	}

	return result
}

// ToSnakeCase converts string to snake_case.
func ToSnakeCase(s string) string {
	var result strings.Builder

	for i, r := range s {
		if i > 0 && r >= 'A' && r <= 'Z' {
			result.WriteRune('_')
		}
		result.WriteRune(r)
	}

	return strings.ToLower(result.String())
}
