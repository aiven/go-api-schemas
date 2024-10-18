// Package gen is the package that contains the generation logic.
package gen

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/aiven/go-api-schemas/internal/pkg/types"
	"github.com/aiven/go-api-schemas/internal/pkg/util"
)

var logger *util.Logger // nolint // Used in the setup function

type doc struct {
	Components struct {
		Schemas map[string]*schema `json:"schemas"`
	} `json:"components"`
}

type schema struct {
	Title       string             `json:"title,omitempty"`
	Description string             `json:"description,omitempty"`
	Type        string             `json:"type,omitempty"`
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
	UserError   string             `json:"x-user_error,omitempty"`
	CreateOnly  bool               `json:"x-createOnly,omitempty"`
	Secure      bool               `json:"x-_secure,omitempty"`
}

const maxSafeInteger = 9007199254740991.0

func isSafeInt(v float64) bool {
	// Maximum is float64, it can't fit more that 2^53-1
	// https://developer.mozilla.org/en-US/docs/Web/JavaScript/Reference/Global_Objects/Number/MAX_SAFE_INTEGER
	return v < maxSafeInteger
}

// reUserConfigKey finds user config keys.
var reUserConfigKey = regexp.MustCompile(`^(Service|IntegrationEndpoint|Integration)([0-9a-zA-Z_]+)UserConfig$`)

func fromSpec(fileName string) (types.GenerationResult, error) {
	b, err := os.ReadFile(filepath.Clean(fileName))
	if err != nil {
		return nil, err
	}

	d := new(doc)
	err = json.Unmarshal(b, &d)
	if err != nil {
		return nil, err
	}

	kinds := map[string]int{
		"Service":             types.KeyServiceTypes,
		"Integration":         types.KeyIntegrationTypes,
		"IntegrationEndpoint": types.KeyIntegrationEndpointTypes,
	}

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
		name := toSnakeCase(match[2])

		// Handle a special case
		name = strings.ReplaceAll(name, "m3_", "m3")

		uc, err := toUserConfig(v)
		if err != nil {
			return nil, fmt.Errorf("failed to convert %s %s: %w", match[1], match[2], err)
		}

		result[kind][name] = *uc
	}

	return result, nil
}

func toUserConfig(src *schema) (*types.UserConfigSchema, error) { // nolint: funlen
	uc := types.UserConfigSchema{
		Properties:  make(map[string]types.UserConfigSchema),
		Title:       src.Title,
		Description: src.Description,
		Required:    src.Required,
		Minimum:     src.Minimum,
		MinLength:   src.MinLength,
		MaxLength:   src.MaxLength,
		MaxItems:    src.MaxItems,
		CreateOnly:  src.CreateOnly,
		Pattern:     src.Pattern,
		UserError:   src.UserError,
		Secure:      src.Secure,
		Example:     formatValue(src.Type, src.Example),
		Default:     formatValue(src.Type, src.Default),
	}

	if src.Maximum != nil {
		if isSafeInt(*src.Maximum) {
			uc.Maximum = src.Maximum
		}
	}

	// Collects all the types
	kinds := make([]string, 0, 1)
	if src.Type != "" {
		kinds = append(kinds, src.Type)
	}

	if src.Nullable {
		kinds = append(kinds, "null")
		uc.Type = kinds
	}

	for _, v := range src.AnyOf {
		if v.Type != "" {
			kinds = append(kinds, v.Type)
		}
	}

	switch len(kinds) {
	case 1:
		uc.Type = kinds[0]
	case 0:
		// fixme: why it has empty type?
		uc.Type = "string"
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
		child, err := toUserConfig(v)
		if err != nil {
			return nil, err
		}
		uc.Properties[k] = *child
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
func Run(l *util.Logger, fileName string) (types.GenerationResult, error) {
	logger = l
	return fromSpec(fileName)
}

var reUnderscore = regexp.MustCompile(`_+`)

// toSnakeCase converts a string to snake case.
// strcase fails with S3 like strings https://github.com/iancoleman/strcase/issues/42
func toSnakeCase(s string) string {
	var r string
	for i, v := range s {
		if i > 0 && v >= 'A' && v <= 'Z' {
			r += "_"
		}
		r += string(v)
	}
	return reUnderscore.ReplaceAllString(strings.ToLower(r), "_")
}
