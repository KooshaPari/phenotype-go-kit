# Contributing to phenotype-go-kit

Thank you for your interest in contributing to phenotype-go-kit!

## Development Setup

1. **Prerequisites**
   - Go 1.21 or later
   - Git

2. **Clone the repository**
   ```bash
   git clone https://github.com/KooshaPari/phenotype-go-kit
   cd phenotype-go-kit
   ```

3. **Install dependencies**
   ```bash
   go mod download
   ```

4. **Run tests**
   ```bash
   go test ./... -short
   ```

5. **Run linter**
   ```bash
   golangci-lint run
   ```

## Pull Request Process

1. Fork the repository and create a feature branch from `main`.
2. Follow the commit message format: `<type>(<scope>): <description>`
   - Types: `feat`, `fix`, `chore`, `docs`, `refactor`, `test`, `ci`
3. Ensure all tests pass and linter is clean.
4. Update documentation if adding new features.
5. Submit a pull request with a clear description of changes.

## Code Standards

- Write idiomatic Go code
- Add tests for new functionality
- Update doc comments for public APIs
- Keep functions focused and small

## Reporting Issues

- Use GitHub Issues for bug reports and feature requests
- Include Go version and relevant code snippets
- For security issues, see SECURITY.md
