package convert

import (
	"errors"
	"fmt"
	"reflect"
	"slices"
	"strings"

	"github.com/aiven/go-api-schemas/internal/pkg/types"
	"github.com/aiven/go-api-schemas/internal/pkg/util"
)

// ConversionError represents a detailed error during conversion.
type ConversionError struct {
	Path        string
	FullPath    []string
	Expected    string
	ActualValue interface{}
	ActualType  string
}

func (e *ConversionError) Error() string {
	fullPath := strings.Join(e.FullPath, ".")

	return fmt.Sprintf(
		"error '%s': expected %s, got value: '%v' type: '%s'",
		fullPath, e.Expected, e.ActualValue, e.ActualType,
	)
}

func (e *ConversionError) appendPath(segment string) {
	e.FullPath = append([]string{segment}, e.FullPath...)
}

// UserConfigSchema converts a map[string]interface{} to UserConfigSchema.
// nolint:lll,nestif,goconst,funlen // This function is long, but it's a conversion function.
func UserConfigSchema(v map[string]interface{}) (*types.UserConfigSchema, error) {
	userConfigSchema := &types.UserConfigSchema{}
	elem := reflect.ValueOf(userConfigSchema).Elem()
	typeOfT := elem.Type()

	// Map of functions to handle different field types
	fieldHandlers := map[reflect.Kind]func(map[string]interface{}, reflect.Value, int, interface{}, string) error{
		reflect.String:    setStringField,    // Title, Description, Pattern, UserError
		reflect.Ptr:       setPointerField,   // Items, Minimum, Maximum, MinLength, MaxLength, MaxItems
		reflect.Interface: setInterfaceField, // Type, Default, Example
		reflect.Map:       setMapField,       // Properties
		reflect.Slice:     setSliceField,     // Required, OneOf, Enum
		reflect.Bool:      setBoolField,      // CreateOnly, Secure
	}

	// Start from 1 to skip the first inlined field
	for i := 1; i < elem.NumField(); i++ {
		field := typeOfT.Field(i)

		// Use json tag to infer the field name and remove `,omitempty` suffix if present
		fieldName := strings.TrimSuffix(field.Tag.Get("json"), ",omitempty")

		value, ok := v[fieldName]
		if !ok || value == nil {
			continue
		}

		// Use the map of functions to handle the field type
		if handler, ok := fieldHandlers[field.Type.Kind()]; ok {
			if err := handler(v, elem, i, value, fieldName); err != nil {
				return nil, err
			}
		}
	}

	return userConfigSchema, nil
}

// extractSchemaTypes extracts and processes the schema types from the input map.
func extractSchemaTypes(v map[string]interface{}) ([]string, error) {
	var schemaTypes []string

	if typeValue, ok := v["type"]; ok && typeValue != nil {
		switch t := typeValue.(type) {
		case string:
			schemaTypes = []string{t}
		case []interface{}:
			for idx, item := range t {
				if str, ok := item.(string); ok {
					schemaTypes = append(schemaTypes, str)
				} else {
					return nil, &ConversionError{
						Expected:    "string",
						ActualValue: item,
						ActualType:  fmt.Sprintf("%T", item),
						FullPath:    []string{fmt.Sprintf("type[%d]", idx)},
					}
				}
			}
		default:
			return nil, &ConversionError{
				Expected:    "string or []interface{}",
				ActualValue: typeValue,
				ActualType:  fmt.Sprintf("%T", typeValue),
				FullPath:    []string{"type"},
			}
		}
	}

	return schemaTypes, nil
}

// setStringField sets the value of a string field in the UserConfigSchema struct.
func setStringField(_ map[string]interface{}, elem reflect.Value, i int, value interface{}, fieldName string) error {
	if str, ok := value.(string); ok {
		elem.Field(i).SetString(str)

		return nil
	}

	return &ConversionError{
		Expected:    "string",
		ActualValue: value,
		ActualType:  fmt.Sprintf("%T", value),
		FullPath:    []string{fieldName},
	}
}

// maxSafeNumber the last number in double precision floating point that can make valid comparison
const maxSafeNumber = float64(1<<53 - 1)

// setPointerField sets the value of a pointer field in the UserConfigSchema struct.
func setPointerField(input map[string]interface{}, elem reflect.Value, i int, value interface{}, fieldName string) error { // nolint:lll
	// Extract and process the schema types
	schemaTypes, err := extractSchemaTypes(input)
	if err != nil {
		return err
	}

	if f, ok := value.(float64); ok {
		switch fieldName {
		case "minimum", "maximum": // nolint:goconst
			if fieldName == "maximum" && slices.Contains(schemaTypes, "integer") && f >= maxSafeNumber {
				f = maxSafeNumber
			}

			elem.Field(i).Set(reflect.ValueOf(&f))
		case "minLength", "maxLength", "maxItems":
			elem.Field(i).Set(reflect.ValueOf(util.Ref(int(f))))
		default:
			return &ConversionError{
				Expected:    "float64",
				ActualValue: value,
				ActualType:  fmt.Sprintf("%T", value),
				FullPath:    []string{fieldName},
			}
		}

		return nil
	} else if _, ok := value.(string); ok {
		return &ConversionError{
			Expected:    "float64",
			ActualValue: value,
			ActualType:  fmt.Sprintf("%T", value),
			FullPath:    []string{fieldName},
		}
	}

	return nil
}

func isValidType(s string) bool {
	// Valid JSON schema types: https://json-schema.org/understanding-json-schema/reference/type
	validTypes := []string{"string", "number", "integer", "object", "array", "boolean", "null"}
	for _, t := range validTypes {
		if s == t {
			return true
		}
	}

	return false
}

// setInterfaceField sets the value of an interface field in the UserConfigSchema struct.
func setInterfaceField(_ map[string]interface{}, elem reflect.Value, i int, value interface{}, fieldName string) error {
	switch fieldName {
	case "type":
		switch v := value.(type) {
		case string:
			elem.Field(i).Set(reflect.ValueOf(v))
		case []interface{}:
			typeSlice := make([]string, len(v))

			for idx, item := range v {
				if s, ok := item.(string); ok {
					if isValidType(s) {
						typeSlice[idx] = s
					} else {
						return &ConversionError{
							Expected:    "valid type",
							ActualValue: s,
							ActualType:  "string",
							FullPath:    []string{fmt.Sprintf("%s[%d]", fieldName, idx)},
						}
					}
				} else {
					return &ConversionError{
						Expected:    "string",
						ActualValue: item,
						ActualType:  fmt.Sprintf("%T", item),
						FullPath:    []string{fmt.Sprintf("%s[%d]", fieldName, idx)},
					}
				}
			}

			elem.Field(i).Set(reflect.ValueOf(typeSlice))
		default:
			return &ConversionError{
				Expected:    "string or []interface{}",
				ActualValue: value,
				ActualType:  fmt.Sprintf("%T", value),
				FullPath:    []string{fieldName},
			}
		}
	case "default", "example":
		elem.Field(i).Set(reflect.ValueOf(value))
	}

	return nil
}

// setMapField sets the value of a map field in the UserConfigSchema struct.
func setMapField(_ map[string]interface{}, elem reflect.Value, i int, value interface{}, fieldName string) error {
	if m, ok := value.(map[string]interface{}); ok { //nolint:nestif
		properties := make(map[string]types.UserConfigSchema)

		for k, v := range m {
			if subMap, ok := v.(map[string]interface{}); ok {
				propSchema, err := UserConfigSchema(subMap)
				if err != nil {
					var convErr *ConversionError
					if errors.As(err, &convErr) {
						convErr.FullPath = append([]string{fieldName, k}, convErr.FullPath...)
					}

					return err
				}

				properties[k] = *propSchema
			} else {
				return &ConversionError{
					Expected:    "map[string]interface{}",
					ActualValue: v,
					ActualType:  fmt.Sprintf("%T", v),
					FullPath:    []string{fieldName, k},
				}
			}
		}

		elem.Field(i).Set(reflect.ValueOf(properties))

		return nil
	}

	return &ConversionError{
		Expected:    "map[string]interface{}",
		ActualValue: value,
		ActualType:  fmt.Sprintf("%T", value),
		FullPath:    []string{fieldName},
	}
}

// setSliceField sets the value of a slice field in the UserConfigSchema struct.
func setSliceField(_ map[string]interface{}, elem reflect.Value, i int, value interface{}, fieldName string) error {
	switch fieldName {
	case "required":
		return setRequiredField(elem, i, value)
	case "oneOf":
		return setOneOfField(elem, i, value)
	case "enum":
		return setEnumField(elem, i, value)
	default:
		return &ConversionError{
			Expected:    "slice",
			ActualValue: value,
			ActualType:  fmt.Sprintf("%T", value),
			FullPath:    []string{fieldName},
		}
	}
}

func setRequiredField(elem reflect.Value, i int, value interface{}) error {
	switch s := value.(type) {
	case []string:
		elem.Field(i).Set(reflect.ValueOf(s))
	case []interface{}:
		required := make([]string, len(s))

		for idx, v := range s {
			if str, ok := v.(string); ok {
				required[idx] = str
			} else {
				return &ConversionError{
					Expected:    "string",
					ActualValue: v,
					ActualType:  fmt.Sprintf("%T", v),
					FullPath:    []string{fmt.Sprintf("required[%d]", idx)},
				}
			}
		}

		elem.Field(i).Set(reflect.ValueOf(required))
	default:
		return &ConversionError{
			Expected:    "[]interface{} or []string",
			ActualValue: value,
			ActualType:  fmt.Sprintf("%T", value),
			FullPath:    []string{"required"},
		}
	}

	return nil
}

func setOneOfField(elem reflect.Value, i int, value interface{}) error {
	if s, ok := value.([]interface{}); ok { //nolint:nestif
		var slice []types.UserConfigSchema

		for idx, v := range s {
			if itemMap, ok := v.(map[string]interface{}); ok {
				itemSchema, err := UserConfigSchema(itemMap)
				if err != nil {
					var convErr *ConversionError
					if errors.As(err, &convErr) {
						convErr.appendPath(fmt.Sprintf("oneOf[%d]", idx))

						return convErr
					}

					return err
				}

				slice = append(slice, *itemSchema)
			} else {
				return &ConversionError{
					Expected:    "map[string]interface{}",
					ActualValue: v,
					ActualType:  fmt.Sprintf("%T", v),
					FullPath:    []string{fmt.Sprintf("oneOf[%d]", idx)},
				}
			}
		}

		elem.Field(i).Set(reflect.ValueOf(slice))
	}

	return nil
}

func setEnumField(elem reflect.Value, i int, value interface{}) error {
	if s, ok := value.([]interface{}); ok {
		var enumValues []types.UserConfigSchemaEnumValue

		for idx, v := range s {
			switch enumValue := v.(type) {
			case string:
				enumValues = append(enumValues, types.UserConfigSchemaEnumValue{Value: enumValue})
			case int, int64, float64:
				enumValues = append(enumValues, types.UserConfigSchemaEnumValue{Value: enumValue})
			default:
				return &ConversionError{
					Expected:    "string or number",
					ActualValue: v,
					ActualType:  fmt.Sprintf("%T", v),
					FullPath:    []string{fmt.Sprintf("enum[%d]", idx)},
				}
			}
		}

		elem.Field(i).Set(reflect.ValueOf(enumValues))
	}

	return nil
}

// setBoolField sets the value of a boolean field in the UserConfigSchema struct.
func setBoolField(_ map[string]interface{}, elem reflect.Value, i int, value interface{}, fieldName string) error {
	if b, ok := value.(bool); ok {
		elem.Field(i).SetBool(b)

		return nil
	}

	return &ConversionError{
		Expected:    "bool",
		ActualValue: value,
		ActualType:  fmt.Sprintf("%T", value),
		FullPath:    []string{fieldName},
	}
}
