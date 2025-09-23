// Package gen is the package that contains the generation logic.
package gen

import (
	"encoding/json"
	"fmt"
	"log"
	"maps"
	"os"
	"path/filepath"
	"regexp"
	"slices"
	"sort"
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
	Minimum     *json.Number       `json:"minimum,omitempty"`
	Maximum     *json.Number       `json:"maximum,omitempty"`
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

		if kind == types.KeyServiceTypes && name == "opensearch" {
			delete(uc.Properties, "custom_repos")
			log.Printf("Removed `custom_repos` from opensearch, because of `one_of`")
		}

		// autoscale_cpu_fraction is an invalid property, it shouldn't be there.
		// todo: remove this when the schema is fixed
		if kind == types.KeyIntegrationEndpointTypes && name == "autoscaler" {
			if autoscaling, ok := uc.Properties["autoscaling"]; ok && autoscaling.Items != nil {
				if t, ok := autoscaling.Items.Properties["type"]; ok {
					t.Enum = filterEnums(t.Enum, func(v any) bool {
						return fmt.Sprint(v) != "autoscale_cpu_fraction"
					})

					// fixme: properties are not pointers
					autoscaling.Items.Properties["type"] = t
				}
			}
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

	// Convert "null" type to a "nullable" boolean field.
	// It is easier to ask a bool field and cast a single "type" instead of a list.
	normTypes = slices.DeleteFunc(normTypes, func(s string) bool {
		if s != "null" {
			return false
		}
		src.Nullable = true
		return true
	})

	// Only mark required fields as nullable since optional fields can't be null in Go.
	src.Nullable = src.isRequired && src.Nullable

	uc := types.UserConfigSchema{
		Properties:  make(map[string]types.UserConfigSchema),
		Title:       normalizeWhitespace(src.Title),
		Description: normalizeWhitespace(src.Description),
		Required:    src.Required,
		Minimum:     src.Minimum,
		Maximum:     src.Maximum,
		MinLength:   src.MinLength,
		MaxLength:   src.MaxLength,
		MaxItems:    src.MaxItems,
		Pattern:     src.Pattern,
		Nullable:    src.Nullable,
		UserError:   or(src.XUserError, src.UserError),
		Secure:      or(src.XSecure, src.Secure),
		CreateOnly:  or(src.XCreateOnly, src.CreateOnly),
		Example:     formatValue(normTypes[0], src.Example),
		Default:     formatValue(normTypes[0], src.Default),
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

	// Sorts enum values for consistent output
	sort.Slice(uc.Enum, func(i, j int) bool {
		return fmt.Sprint(uc.Enum[i].Value) < fmt.Sprint(uc.Enum[j].Value)
	})

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
	result := make(map[string]bool)
	switch t := s.Type.(type) {
	case string:
		return []string{t}, nil
	case []any:
		for _, v := range t {
			result[fmt.Sprintf("%v", v)] = true
		}
	case nil:
	default:
		return nil, fmt.Errorf("unknown type %T", s.Type)
	}

	for _, v := range s.AnyOf {
		if v.Type != nil {
			vType, err := normalizeType(v)
			if err != nil {
				return nil, err
			}

			for _, t := range vType {
				result[t] = true
			}
		}
	}

	// Fixes the case when there is no type defined
	if len(result) == 0 {
		result["string"] = true
	}

	// Prioritize number over integer, because number is more general, and Go can't have both.
	if result["number"] && result["integer"] {
		delete(result, "integer")
	}

	keys := slices.Collect(maps.Keys(result))
	slices.Sort(keys)
	return keys, nil
}

func formatValue(t string, v any) any {
	if v == nil {
		return nil
	}

	s := fmt.Sprintf("%v", v)
	if s == "" {
		return nil
	}

	switch t {
	case "integer":
		n := json.Number(s)
		_, err := n.Int64()
		if err != nil {
			return nil
		}
		return n.String()
	case "number":
		n := json.Number(s)
		_, err := n.Float64()
		if err != nil {
			return nil
		}
		return n.String()
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

	removeBlockedFields(result)
	return result, nil
}

func removeBlockedFields(result types.GenerationResult) {
	services, ok := result[types.KeyServiceTypes]
	if ok {
		if v, ok := services["opensearch"]; ok {
			delete(v.Properties, "elasticsearch_version")
		}
	}
}

var reWhitespace = regexp.MustCompile(`\s+`)

func normalizeWhitespace(s string) string {
	return strings.TrimSpace(reWhitespace.ReplaceAllString(s, " "))
}

func filterEnums(enums []types.UserConfigSchemaEnumValue, keep func(v any) bool) []types.UserConfigSchemaEnumValue {
	result := make([]types.UserConfigSchemaEnumValue, 0, len(enums))
	for _, v := range enums {
		if keep(v.Value) {
			result = append(result, v)
		}
	}
	return result
}
