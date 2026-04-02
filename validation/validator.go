package validation

import (
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"strings"
)

// Validator provides data validation.
type Validator struct {
	rules map[string][]Rule
}

// Rule represents a validation rule.
type Rule struct {
	Name     string
	Message  string
	Validate func(interface{}) bool
}

// New creates a new validator.
func New() *Validator {
	return &Validator{
		rules: make(map[string][]Rule),
	}
}

// AddRule adds a validation rule for a field.
func (v *Validator) AddRule(field string, rule Rule) {
	v.rules[field] = append(v.rules[field], rule)
}

// Validate validates a struct.
func (v *Validator) Validate(data interface{}) (map[string][]string, bool) {
	errors := make(map[string][]string)

	rv := reflect.ValueOf(data)
	if rv.Kind() == reflect.Ptr {
		rv = rv.Elem()
	}

	if rv.Kind() != reflect.Struct {
		return nil, false
	}

	rt := rv.Type()

	for fieldName, rules := range v.rules {
		var fieldValue reflect.Value
		var found bool

		for i := 0; i < rt.NumField(); i++ {
			field := rt.Field(i)
			jsonTag := field.Tag.Get("json")
			fieldNameFromTag := strings.Split(jsonTag, ",")[0]

			if fieldNameFromTag == fieldName || field.Name == fieldName {
				fieldValue = rv.Field(i)
				found = true
				break
			}
		}

		if !found {
			continue
		}

		fieldVal := fieldValue.Interface()

		for _, rule := range rules {
			if !rule.Validate(fieldVal) {
				errors[fieldName] = append(errors[fieldName], rule.Message)
			}
		}
	}

	return errors, len(errors) == 0
}

// Common validation rules.
var (
	Required = Rule{
		Name:    "required",
		Message: "is required",
		Validate: func(v interface{}) bool {
			if v == nil {
				return false
			}
			if s, ok := v.(string); ok {
				return strings.TrimSpace(s) != ""
			}
			return true
		},
	}

	Email = Rule{
		Name:    "email",
		Message: "must be a valid email address",
		Validate: func(v interface{}) bool {
			s, ok := v.(string)
			if !ok {
				return false
			}
			pattern := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
			matched, _ := regexp.MatchString(pattern, s)
			return matched
		},
	}
)

func MinLength(n int) Rule {
	return Rule{
		Name:    "min_length",
		Message: fmt.Sprintf("must be at least %d characters", n),
		Validate: func(v interface{}) bool {
			s, ok := v.(string)
			if !ok {
				return false
			}
			return len(s) >= n
		},
	}
}

func MaxLength(n int) Rule {
	return Rule{
		Name:    "max_length",
		Message: fmt.Sprintf("must be at most %d characters", n),
		Validate: func(v interface{}) bool {
			s, ok := v.(string)
			if !ok {
				return false
			}
			return len(s) <= n
		},
	}
}

func Min(n int) Rule {
	return Rule{
		Name:    "min",
		Message: fmt.Sprintf("must be at least %d", n),
		Validate: func(v interface{}) bool {
			switch val := v.(type) {
			case int:
				return val >= n
			case int64:
				return val >= int64(n)
			case float64:
				return val >= float64(n)
			case string:
				i, err := strconv.Atoi(val)
				return err == nil && i >= n
			}
			return false
		},
	}
}

func Max(n int) Rule {
	return Rule{
		Name:    "max",
		Message: fmt.Sprintf("must be at most %d", n),
		Validate: func(v interface{}) bool {
			switch val := v.(type) {
			case int:
				return val <= n
			case int64:
				return val <= int64(n)
			case float64:
				return val <= float64(n)
			case string:
				i, err := strconv.Atoi(val)
				return err == nil && i <= n
			}
			return false
		},
	}
}

func Pattern(p string, msg string) Rule {
	return Rule{
		Name:    "pattern",
		Message: msg,
		Validate: func(v interface{}) bool {
			s, ok := v.(string)
			if !ok {
				return false
			}
			matched, _ := regexp.MatchString(p, s)
			return matched
		},
	}
}

func In(valid ...string) Rule {
	return Rule{
		Name:    "in",
		Message: fmt.Sprintf("must be one of: %s", strings.Join(valid, ", ")),
		Validate: func(v interface{}) bool {
			s, ok := v.(string)
			if !ok {
				return false
			}
			for _, x := range valid {
				if s == x {
					return true
				}
			}
			return false
		},
	}
}
