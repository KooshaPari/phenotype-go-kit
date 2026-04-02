// Package docs provides documentation generation utilities.
package docs

import (
	"bytes"
	"fmt"
	"strings"
)

// APIDoc represents API documentation.
type APIDoc struct {
	Title       string
	Description string
	Version     string
	BasePath    string
	Endpoints   []Endpoint
}

// Endpoint represents an API endpoint.
type Endpoint struct {
	Path        string
	Method      string
	Summary     string
	Description string
	Parameters  []Parameter
	RequestBody interface{}
	Response    map[string]interface{}
	Security    []string
}

// Parameter represents an API parameter.
type Parameter struct {
	Name        string
	In          string // path, query, header, body
	Type        string
	Required    bool
	Description string
}

// GenerateOpenAPI generates OpenAPI 3.0 specification.
func (a *APIDoc) GenerateOpenAPI() string {
	var buf bytes.Buffer

	buf.WriteString("openapi: 3.0.0\n")
	buf.WriteString(fmt.Sprintf("info:\n  title: %s\n  description: %s\n  version: %s\n\n", a.Title, a.Description, a.Version))
	buf.WriteString(fmt.Sprintf("servers:\n  - url: %s\n\n", a.BasePath))
	buf.WriteString("paths:\n")

	for _, ep := range a.Endpoints {
		buf.WriteString(fmt.Sprintf("  %s:\n    %s:\n      summary: %s\n      description: %s\n", ep.Path, strings.ToLower(ep.Method), ep.Summary, ep.Description))

		if len(ep.Parameters) > 0 {
			buf.WriteString("      parameters:\n")
			for _, p := range ep.Parameters {
				required := "false"
				if p.Required {
					required = "true"
				}
				buf.WriteString(fmt.Sprintf("        - name: %s\n          in: %s\n          required: %s\n          schema:\n            type: %s\n          description: %s\n", p.Name, p.In, required, p.Type, p.Description))
			}
		}

		if ep.RequestBody != nil {
			buf.WriteString("      requestBody:\n        content:\n          application/json:\n            schema:\n              type: object\n")
		}

		if len(ep.Response) > 0 {
			buf.WriteString("      responses:\n")
			for code, resp := range ep.Response {
				buf.WriteString(fmt.Sprintf("        '%s':\n          description: %v\n", code, resp))
			}
		}

		if len(ep.Security) > 0 {
			buf.WriteString("      security:\n")
			for _, sec := range ep.Security {
				buf.WriteString(fmt.Sprintf("        - %s: []\n", sec))
			}
		}

		buf.WriteString("\n")
	}

	return buf.String()
}

// MarkdownDocs generates Markdown documentation.
func (a *APIDoc) MarkdownDocs() string {
	var buf bytes.Buffer

	buf.WriteString(fmt.Sprintf("# %s\n\n", a.Title))
	buf.WriteString(fmt.Sprintf("%s\n\n", a.Description))
	buf.WriteString(fmt.Sprintf("**Version:** %s\n\n", a.Version))
	buf.WriteString(fmt.Sprintf("**Base URL:** %s\n\n", a.BasePath))

	buf.WriteString("## Endpoints\n\n")

	for _, ep := range a.Endpoints {
		buf.WriteString(fmt.Sprintf("### %s %s\n\n", ep.Method, ep.Path))
		buf.WriteString(fmt.Sprintf("**%s**\n\n%s\n\n", ep.Summary, ep.Description))

		if len(ep.Parameters) > 0 {
			buf.WriteString("**Parameters:**\n\n")
			buf.WriteString("| Name | In | Type | Required | Description |\n")
			buf.WriteString("|------|----|------|----------|-------------|\n")
			for _, p := range ep.Parameters {
				required := "No"
				if p.Required {
					required = "Yes"
				}
				buf.WriteString(fmt.Sprintf("| %s | %s | %s | %s | %s |\n", p.Name, p.In, p.Type, required, p.Description))
			}
			buf.WriteString("\n")
		}

		if len(ep.Response) > 0 {
			buf.WriteString("**Responses:**\n\n")
			for code, resp := range ep.Response {
				buf.WriteString(fmt.Sprintf("- `%s`: %v\n", code, resp))
			}
			buf.WriteString("\n")
		}
	}

	return buf.String()
}

// NewAPIDoc creates a new API documentation.
func NewAPIDoc(title, description, version, basePath string) *APIDoc {
	return &APIDoc{
		Title:       title,
		Description: description,
		Version:     version,
		BasePath:    basePath,
		Endpoints:   make([]Endpoint, 0),
	}
}

// AddEndpoint adds an endpoint.
func (a *APIDoc) AddEndpoint(ep Endpoint) {
	a.Endpoints = append(a.Endpoints, ep)
}
