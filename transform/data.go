package transform

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strings"
)

// Mapper provides data transformation between structs.
type Mapper struct {
	tagName string
}

// New creates a new mapper.
func New(tagName string) *Mapper {
	if tagName == "" {
		tagName = "json"
	}
	return &Mapper{tagName: tagName}
}

// Map transforms from one struct to another.
func (m *Mapper) Map(src, dest interface{}) error {
	srcVal := reflect.ValueOf(src)
	destVal := reflect.ValueOf(dest)

	if srcVal.Kind() != reflect.Ptr || destVal.Kind() != reflect.Ptr {
		return fmt.Errorf("src and dest must be pointers")
	}

	srcVal = srcVal.Elem()
	destVal = destVal.Elem()

	if srcVal.Kind() != reflect.Struct || destVal.Kind() != reflect.Struct {
		return fmt.Errorf("src and dest must be structs")
	}

	return m.mapValues(srcVal, destVal)
}

func (m *Mapper) mapValues(src, dest reflect.Value) error {
	destType := dest.Type()

	for i := 0; i < destType.NumField(); i++ {
		field := destType.Field(i)
		tag := field.Tag.Get(m.tagName)

		if tag == "-" {
			continue
		}

		fieldName := field.Name
		if tag != "" {
			fieldName = strings.Split(tag, ",")[0]
		}

		srcField := m.findField(src, fieldName)
		if !srcField.IsValid() {
			continue
		}

		if !srcField.Type().AssignableTo(field.Type) {
			continue
		}

		dest.Field(i).Set(srcField)
	}

	return nil
}

func (m *Mapper) findField(v reflect.Value, name string) reflect.Value {
	for i := 0; i < v.Type().NumField(); i++ {
		field := v.Type().Field(i)

		tag := field.Tag.Get(m.tagName)
		tagName := strings.Split(tag, ",")[0]

		if tagName == name || field.Name == name {
			return v.Field(i)
		}
	}

	return reflect.Value{}
}

// ToMap converts a struct to a map.
func ToMap(src interface{}, tagName string) (map[string]interface{}, error) {
	result := make(map[string]interface{})

	v := reflect.ValueOf(src)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	if v.Kind() != reflect.Struct {
		return nil, fmt.Errorf("src must be a struct")
	}

	t := v.Type()

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		tag := field.Tag.Get(tagName)

		if tag == "-" || tag == "" {
			continue
		}

		name := strings.Split(tag, ",")[0]
		result[name] = v.Field(i).Interface()
	}

	return result, nil
}

// FromMap converts a map to a struct.
func FromMap(src map[string]interface{}, dest interface{}, tagName string) error {
	v := reflect.ValueOf(dest)
	if v.Kind() != reflect.Ptr {
		return fmt.Errorf("dest must be a pointer")
	}

	v = v.Elem()
	if v.Kind() != reflect.Struct {
		return fmt.Errorf("dest must be a struct")
	}

	t := v.Type()

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		tag := field.Tag.Get(tagName)

		if tag == "-" || tag == "" {
			continue
		}

		name := strings.Split(tag, ",")[0]

		if val, ok := src[name]; ok {
			fieldVal := v.Field(i)
			if fieldVal.CanSet() {
				valVal := reflect.ValueOf(val)
				fieldType := field.Type
				if valVal.Type().AssignableTo(fieldType) {
					fieldVal.Set(valVal)
				}
			}
		}
	}

	return nil
}

// Convert converts between types using JSON marshaling.
func Convert(src interface{}, dest interface{}) error {
	data, err := json.Marshal(src)
	if err != nil {
		return err
	}
	return json.Unmarshal(data, dest)
}

// Merge merges source into destination (src overwrites dest).
func Merge(dest, src interface{}) error {
	destMap, err := ToMap(dest, "json")
	if err != nil {
		return err
	}

	srcMap, err := ToMap(src, "json")
	if err != nil {
		return err
	}

	for k, v := range srcMap {
		destMap[k] = v
	}

	return FromMap(destMap, dest, "json")
}

// Flatten flattens nested maps.
func Flatten(src map[string]interface{}, prefix string) map[string]interface{} {
	result := make(map[string]interface{})

	for k, v := range src {
		key := k
		if prefix != "" {
			key = prefix + "." + k
		}

		if m, ok := v.(map[string]interface{}); ok {
			sub := Flatten(m, key)
			for kk, vv := range sub {
				result[kk] = vv
			}
		} else {
			result[key] = v
		}
	}

	return result
}

// Unflatten unflattens a flat map.
func Unflatten(src map[string]interface{}) map[string]interface{} {
	result := make(map[string]interface{})

	for k, v := range src {
		keys := strings.Split(k, ".")

		current := result
		for i := 0; i < len(keys)-1; i++ {
			if _, ok := current[keys[i]]; !ok {
				current[keys[i]] = make(map[string]interface{})
			}
			current = current[keys[i]].(map[string]interface{})
		}

		current[keys[len(keys)-1]] = v
	}

	return result
}
