package frontend

import (
	"fmt"
	"regexp"
)

// Field represents a form field.
type Field struct {
	Name        string
	Value       interface{}
	Label       string
	Placeholder string
	Error       string
	Required    bool
	Valid       bool
}

// Form represents a form with fields.
type Form struct {
	Fields map[string]*Field
	Valid  bool
	Errors map[string]string
}

// New creates a new form.
func New() *Form {
	return &Form{
		Fields: make(map[string]*Field),
		Errors: make(map[string]string),
	}
}

// AddField adds a field to the form.
func (f *Form) AddField(name, label, placeholder string, required bool) {
	f.Fields[name] = &Field{
		Name:        name,
		Label:       label,
		Placeholder: placeholder,
		Required:    required,
		Valid:       true,
	}
}

// SetValue sets a field value.
func (f *Form) SetValue(name string, value interface{}) {
	if field, ok := f.Fields[name]; ok {
		field.Value = value
	}
}

// Validate validates all fields.
func (f *Form) Validate() bool {
	f.Valid = true
	f.Errors = make(map[string]string)

	for name, field := range f.Fields {
		if field.Required {
			if field.Value == nil || field.Value == "" {
				f.Errors[name] = fmt.Sprintf("%s is required", field.Label)
				field.Valid = false
				f.Valid = false
				continue
			}
		}

		field.Valid = true
	}

	return f.Valid
}

// ErrorsToMap converts errors to map.
func (f *Form) ErrorsToMap() map[string]string {
	return f.Errors
}

// ValidationRule represents a validation rule.
type ValidationRule func(interface{}) error

// RequiredRule returns a required validation rule.
func RequiredRule(message string) ValidationRule {
	return func(value interface{}) error {
		if value == nil || value == "" {
			return fmt.Errorf("%s", message)
		}
		return nil
	}
}

// EmailRule returns an email validation rule.
func EmailRule() ValidationRule {
	return func(value interface{}) error {
		s, ok := value.(string)
		if !ok {
			return nil
		}

		pattern := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
		matched, _ := regexp.MatchString(pattern, s)
		if !matched {
			return fmt.Errorf("must be a valid email address")
		}
		return nil
	}
}

// MinLengthRule returns a minimum length rule.
func MinLengthRule(min int) ValidationRule {
	return func(value interface{}) error {
		s, ok := value.(string)
		if !ok {
			return nil
		}

		if len(s) < min {
			return fmt.Errorf("must be at least %d characters", min)
		}
		return nil
	}
}

// MaxLengthRule returns a maximum length rule.
func MaxLengthRule(max int) ValidationRule {
	return func(value interface{}) error {
		s, ok := value.(string)
		if !ok {
			return nil
		}

		if len(s) > max {
			return fmt.Errorf("must be at most %d characters", max)
		}
		return nil
	}
}

// PatternRule returns a pattern matching rule.
func PatternRule(pattern, message string) ValidationRule {
	re := regexp.MustCompile(pattern)
	return func(value interface{}) error {
		s, ok := value.(string)
		if !ok {
			return nil
		}

		if !re.MatchString(s) {
			return fmt.Errorf("%s", message)
		}
		return nil
	}
}

// MinRule returns a minimum value rule.
func MinRule(min int) ValidationRule {
	return func(value interface{}) error {
		switch v := value.(type) {
		case int:
			if v < min {
				return fmt.Errorf("must be at least %d", min)
			}
		case float64:
			if v < float64(min) {
				return fmt.Errorf("must be at least %d", min)
			}
		}
		return nil
	}
}

// MaxRule returns a maximum value rule.
func MaxRule(max int) ValidationRule {
	return func(value interface{}) error {
		switch v := value.(type) {
		case int:
			if v > max {
				return fmt.Errorf("must be at most %d", max)
			}
		case float64:
			if v > float64(max) {
				return fmt.Errorf("must be at most %d", max)
			}
		}
		return nil
	}
}

// ApplyRule applies a validation rule to a field.
func (f *Form) ApplyRule(fieldName string, rule ValidationRule) {
	if field, ok := f.Fields[fieldName]; ok {
		if err := rule(field.Value); err != nil {
			f.Errors[fieldName] = err.Error()
			field.Error = err.Error()
			field.Valid = false
		}
	}
}

// ValidateWithRules validates fields with specific rules.
func (f *Form) ValidateWithRules(fieldName string, rules ...ValidationRule) {
	if field, ok := f.Fields[fieldName]; ok {
		for _, rule := range rules {
			if err := rule(field.Value); err != nil {
				f.Errors[fieldName] = err.Error()
				field.Error = err.Error()
				field.Valid = false
				break
			}
		}
	}
}
