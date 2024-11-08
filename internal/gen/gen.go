// Package gen is the package that contains the generation logic.
package gen

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"slices"
	"strings"

	"github.com/huandu/xstrings"

	"github.com/aiven/go-api-schemas/internal/pkg/types"
)

type doc struct {
	// openapi-uc.json
	Components struct {
		Schemas map[string]*schema `json:"schemas"`
	} `json:"components"`

	// Legacy user config files
	legacyDoc
}

type schema struct {
	Title       string             `json:"title,omitempty"`
	Description string             `json:"description,omitempty"`
	Type        any                `json:"type,omitempty"`
	Default     interface{}        `json:"default,omitempty"`
	Required    []string           `json:"required,omitempty"`
	Properties  map[string]*schema `json:"properties,omitempty"`
	Items       *schema            `json:"items,omitempty"`
	AnyOf       []*schema          `json:"anyOf,omitempty"`
	OneOf       []*schema          `json:"oneOf,omitempty"`
	Enum        []any              `json:"enum,omitempty"`
	Minimum     *float64           `json:"minimum,omitempty"`
	Maximum     *float64           `json:"maximum,omitempty"`
	MinLength   *int               `json:"minLength,omitempty"`
	MaxLength   *int               `json:"maxLength,omitempty"`
	MaxItems    *int               `json:"maxItems,omitempty"`
	Pattern     string             `json:"pattern,omitempty"`
	Example     any                `json:"example,omitempty"`
	Nullable    bool               `json:"nullable,omitempty"`
	UserError   string             `json:"user_error,omitempty"`
	CreateOnly  bool               `json:"createOnly,omitempty"`
	Secure      bool               `json:"_secure,omitempty"`

	// Openapi-uc
	XUserError  string `json:"x-user_error,omitempty"`
	XCreateOnly bool   `json:"x-createOnly,omitempty"`
	XSecure     bool   `json:"x-_secure,omitempty"`

	// Internal
	isRequired bool // true, when a field is on the parent's Required list
}

// maxSafeInteger
// https://developer.mozilla.org/en-US/docs/Web/JavaScript/Reference/Global_Objects/Number/MAX_SAFE_INTEGER
const maxSafeInteger = float64(1<<53 - 1) // Or 9007199254740991.0

func isSafeInt(v float64) bool {
	// Maximum is float64, it can't fit more that 2^53-1
	// https://developer.mozilla.org/en-US/docs/Web/JavaScript/Reference/Global_Objects/Number/MAX_SAFE_INTEGER
	return v < maxSafeInteger
}

// reUserConfigKey finds user config keys.
var reUserConfigKey = regexp.MustCompile(`^(Service|IntegrationEndpoint|Integration)([0-9a-zA-Z_]+)UserConfig$`)

func fromFile(fileName string) (types.GenerationResult, error) {
	b, err := os.ReadFile(filepath.Clean(fileName))
	if err != nil {
		return nil, err
	}

	d := new(doc)
	err = json.Unmarshal(b, &d)
	if err != nil {
		return nil, err
	}

	legacyToComponents(d)

	kinds := map[string]int{
		"Service":             types.KeyServiceTypes,
		"Integration":         types.KeyIntegrationTypes,
		"IntegrationEndpoint": types.KeyIntegrationEndpointTypes,
	}

	// New openapi-uc schema file
	result := make(types.GenerationResult)
	for _, v := range kinds {
		result[v] = make(map[string]types.UserConfigSchema)
	}

	for k, v := range d.Components.Schemas {
		match := reUserConfigKey.FindStringSubmatch(k)
		if len(match) == 0 {
			continue
		}

		kind := kinds[match[1]]
		name := match[2]
		if !strings.HasPrefix(name, "m3") {
			// m3 is a special case
			name = xstrings.ToSnakeCase(match[2])
		}

		uc, err := toUserConfig(v)
		if err != nil {
			return nil, fmt.Errorf("failed to convert %s %s: %w", match[1], match[2], err)
		}

		result[kind][name] = *uc
	}

	return result, nil
}

func toUserConfig(src *schema) (*types.UserConfigSchema, error) { // nolint: funlen
	normTypes, err := normalizeType(src)
	if err != nil {
		return nil, err
	}

	uc := types.UserConfigSchema{
		Properties:  make(map[string]types.UserConfigSchema),
		Title:       src.Title,
		Description: src.Description,
		Required:    src.Required,
		Minimum:     src.Minimum,
		MinLength:   src.MinLength,
		MaxLength:   src.MaxLength,
		MaxItems:    src.MaxItems,
		Pattern:     src.Pattern,
		UserError:   or(src.XUserError, src.UserError),
		Secure:      or(src.XSecure, src.Secure),
		CreateOnly:  or(src.XCreateOnly, src.CreateOnly),
		Example:     formatValue(normTypes[0], src.Example),
		Default:     formatValue(normTypes[0], src.Default),
	}

	if src.Maximum != nil {
		if isSafeInt(*src.Maximum) {
			uc.Maximum = src.Maximum
		}
	}

	// Collects all the types
	kinds := make([]string, 0, 1)
	if normTypes != nil {
		kinds = append(kinds, normTypes...)
	}

	switch len(kinds) {
	case 1:
		uc.Type = kinds[0]
	default:
		uc.Type = kinds
	}

	for _, v := range src.Enum {
		if fmt.Sprintf("%v", v) != "" {
			uc.Enum = append(uc.Enum, types.UserConfigSchemaEnumValue{Value: v})
		}
	}

	if src.Items != nil {
		item, err := toUserConfig(src.Items)
		if err != nil {
			return nil, err
		}
		uc.Items = item
	}

	for k, v := range src.Properties {
		v.isRequired = slices.Contains(src.Required, k)
		child, err := toUserConfig(v)
		if err != nil {
			return nil, err
		}
		uc.Properties[k] = *child
	}

	if len(src.OneOf) != 0 {
		// Must remove the type if there is a oneOf
		uc.Type = nil
	}

	for _, v := range src.OneOf {
		child, err := toUserConfig(v)
		if err != nil {
			return nil, err
		}
		uc.OneOf = append(uc.OneOf, *child)
	}

	return &uc, nil
}

func distinct[T comparable](list []T) []T {
	seen := make(map[T]bool)
	result := make([]T, 0, len(list))
	for _, v := range list {
		if !seen[v] {
			seen[v] = true
			result = append(result, v)
		}
	}
	return result
}

// or returns the first non-zero value.
func or[T comparable](args ...T) T {
	var zero T
	for _, a := range args {
		if a != zero {
			return a
		}
	}
	return zero
}

func normalizeType(s *schema) ([]string, error) {
	result := make([]string, 0)
	switch t := s.Type.(type) {
	case string:
		return []string{t}, nil
	case []any:
		for _, v := range t {
			if v == "null" {
				s.Nullable = true
			} else {
				result = append(result, fmt.Sprintf("%v", v))
			}
		}
	case nil:
	default:
		return nil, fmt.Errorf("unknown type %T", s.Type)
	}

	// Golang supports null for required fields only
	if s.isRequired && s.Nullable {
		result = append(result, "null")
	}

	for _, v := range s.AnyOf {
		if v.Type != nil {
			vType, err := normalizeType(v)
			if err != nil {
				return nil, err
			}

			result = append(result, vType...)
		}
	}

	if len(result) == 0 {
		result = append(result, "string")
	}

	return distinct(result), nil
}

func formatValue(t string, v any) any {
	s := fmt.Sprintf("%v", v)
	if v == nil || s == "" {
		return nil
	}

	switch t {
	case "integer":
		i, ok := v.(float64)
		if ok && isSafeInt(i) {
			return fmt.Sprintf("%d", int(i))
		}
		return nil
	case "number":
		return fmt.Sprintf("%.1f", v)
	case "boolean", "array", "object":
		return v
	}

	return s
}

// Run executes the generation process.
func Run(fileNames ...string) (types.GenerationResult, error) {
	result := make(types.GenerationResult)
	for _, fileName := range fileNames {
		subResult, err := fromFile(fileName)
		if err != nil {
			return nil, err
		}

		for kind, value := range subResult {
			if len(value) == 0 {
				continue
			}

			if _, ok := result[kind]; !ok {
				result[kind] = value
				continue
			}

			for k, v := range value {
				result[kind][k] = v
			}
		}
	}
	return result, nil
}
