package convert

import (
	"fmt"
	"reflect"
	"slices"
	"strings"

	"github.com/aiven/go-api-schemas/internal/pkg/types"
	"github.com/aiven/go-api-schemas/internal/pkg/util"
)

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
			for _, item := range t {
				if str, ok := item.(string); ok {
					schemaTypes = append(schemaTypes, str)
				} else {
					return nil, fmt.Errorf("error converting type: expected string in slice, got %T", item)
				}
			}
		default:
			return nil, fmt.Errorf("error converting type: expected string or []interface{}, got %T", typeValue)
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

	return fmt.Errorf("error converting %s: expected string, got %T", fieldName, value)
}

// maxSafeNumber the last number in double precision floating point that can make valid comparison
const maxSafeNumber = float64(1<<53 - 1)

// setPointerField sets the value of a pointer field in the UserConfigSchema struct.
func setPointerField(input map[string]interface{}, elem reflect.Value, i int, value interface{}, fieldName string) error {
	// Extract and process the schema types
	schemaTypes, err := extractSchemaTypes(input)
	if err != nil {
		return err
	}

	if f, ok := value.(float64); ok {
		switch fieldName {
		case "minimum", "maximum":
			if fieldName == "maximum" && slices.Contains(schemaTypes, "integer") && f >= maxSafeNumber {
				f = maxSafeNumber
			}

			elem.Field(i).Set(reflect.ValueOf(&f))
		case "minLength", "maxLength", "maxItems":
			elem.Field(i).Set(reflect.ValueOf(util.Ref(int(f))))
		default:
			return fmt.Errorf("unsupported pointer field: %s", fieldName)
		}

		return nil
	} else if _, ok := value.(string); ok {
		return fmt.Errorf("error converting %s: expected float64, got %T", fieldName, value)
	}

	return nil
}

// setInterfaceField sets the value of an interface field in the UserConfigSchema struct.
func setInterfaceField(_ map[string]interface{}, elem reflect.Value, i int, value interface{}, fieldName string) error {
	switch fieldName {
	case "type", "default", "example":
		elem.Field(i).Set(reflect.ValueOf(value))
	}

	return nil
}

// setMapField sets the value of a map field in the UserConfigSchema struct.
func setMapField(_ map[string]interface{}, elem reflect.Value, i int, value interface{}, _ string) error {
	if m, ok := value.(map[string]interface{}); ok {
		properties := make(map[string]types.UserConfigSchema)

		for k, v := range m {
			if subMap, ok := v.(map[string]interface{}); ok {
				propSchema, err := UserConfigSchema(subMap)
				if err != nil {
					return fmt.Errorf("error converting property %s: %w", k, err)
				}

				properties[k] = *propSchema
			} else {
				return fmt.Errorf("error converting property %s: expected map[string]interface{}, got %T", k, v)
			}
		}

		elem.Field(i).Set(reflect.ValueOf(properties))

		return nil
	}

	return fmt.Errorf("error converting properties: expected map[string]interface{}, got %T", value)
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
		return fmt.Errorf("unsupported slice field: %s", fieldName)
	}
}

func setRequiredField(elem reflect.Value, i int, value interface{}) error {
	switch s := value.(type) {
	case []string:
		elem.Field(i).Set(reflect.ValueOf(s))
	case []interface{}:
		required := make([]string, len(s))
		for i, v := range s {
			if str, ok := v.(string); ok {
				required[i] = str
			} else {
				return fmt.Errorf("error converting required field at index %d: expected string, got %T", i, v)
			}
		}
		elem.Field(i).Set(reflect.ValueOf(required))
	default:
		return fmt.Errorf("error converting required fields: expected []interface{} or []string, got %T", value)
	}
	return nil
}

func setOneOfField(elem reflect.Value, i int, value interface{}) error {
	if s, ok := value.([]interface{}); ok {
		var slice []types.UserConfigSchema
		for i, v := range s {
			if itemMap, ok := v.(map[string]interface{}); ok {
				itemSchema, err := UserConfigSchema(itemMap)
				if err != nil {
					return fmt.Errorf("error converting slice item at index %d: %w", i, err)
				}
				slice = append(slice, *itemSchema)
			} else {
				return fmt.Errorf("error converting slice item at index %d: expected map[string]interface{}, got %T", i, v)
			}
		}
		elem.Field(i).Set(reflect.ValueOf(slice))
	}
	return nil
}

func setEnumField(elem reflect.Value, i int, value interface{}) error {
	if s, ok := value.([]interface{}); ok {
		var enumValues []types.UserConfigSchemaEnumValue
		for _, v := range s {
			switch enumValue := v.(type) {
			case string:
				enumValues = append(enumValues, types.UserConfigSchemaEnumValue{Value: enumValue})
			case int, int64, float64:
				enumValues = append(enumValues, types.UserConfigSchemaEnumValue{Value: enumValue})
			default:
				return fmt.Errorf("error converting enum value: expected string or number, got %T", v)
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

	return fmt.Errorf("error converting %s: expected bool, got %T", fieldName, value)
}
