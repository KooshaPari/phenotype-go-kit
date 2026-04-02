package docs

import (
	"bytes"
	"fmt"
	"time"
)

// ArchitectureDoc represents architecture documentation.
type ArchitectureDoc struct {
	Title       string
	LastUpdated time.Time
	Sections    []Section
}

// Section represents a documentation section.
type Section struct {
	Title    string
	Level    int
	Content  string
	Children []Section
}

// NewArchitectureDoc creates architecture documentation.
func NewArchitectureDoc(title string) *ArchitectureDoc {
	return &ArchitectureDoc{
		Title:       title,
		LastUpdated: time.Now(),
		Sections:    make([]Section, 0),
	}
}

// AddSection adds a section.
func (a *ArchitectureDoc) AddSection(title string, level int, content string) {
	a.Sections = append(a.Sections, Section{
		Title:   title,
		Level:   level,
		Content: content,
	})
}

// GenerateMarkdown creates Markdown documentation.
func (a *ArchitectureDoc) GenerateMarkdown() string {
	var buf bytes.Buffer

	buf.WriteString(fmt.Sprintf("# %s\n\n", a.Title))
	buf.WriteString(fmt.Sprintf("*Last Updated: %s*\n\n", a.LastUpdated.Format("2006-01-02")))

	for _, s := range a.Sections {
		level := ""
		for i := 0; i < s.Level; i++ {
			level += "#"
		}
		buf.WriteString(fmt.Sprintf("%s %s\n\n%s\n\n", level, s.Title, s.Content))
	}

	return buf.String()
}

// ArchitectureOverview returns the main architecture document.
func ArchitectureOverview() *ArchitectureDoc {
	doc := NewArchitectureDoc("Phenotype Go Kit Architecture")

	doc.AddSection("Overview", 1, "Phenotype Go Kit is a comprehensive backend infrastructure library providing observability, database scaling, API gateway, microservices communication, and CI/CD pipeline support.")

	doc.AddSection("System Architecture", 2, "The system follows hexagonal architecture principles with clear separation between domain, application, and infrastructure layers.")

	doc.AddSection("Core Components", 2, "")
	doc.AddSection("Jobs", 3, "Background job processing with queue management, webhook delivery, and scheduled tasks.")
	doc.AddSection("Observability", 3, "Structured logging, OpenTelemetry tracing, Prometheus metrics, and health checks.")
	doc.AddSection("Database", 3, "Connection pooling, Redis caching, query optimization, and migrations.")
	doc.AddSection("API Gateway", 3, "JWT authentication, API keys, rate limiting, CORS, and OAuth2 support.")
	doc.AddSection("Microservices", 3, "Event bus, service discovery, circuit breaker, and retry mechanisms.")
	doc.AddSection("Data Layer", 3, "Validation, transformation, and repository pattern implementations.")
	doc.AddSection("Frontend", 3, "HTTP client, state management, form validation, and UI utilities.")
	doc.AddSection("CI/CD", 3, "Pipeline definition, Docker configuration, Kubernetes/Helm deployment, and secrets management.")

	return doc
}
