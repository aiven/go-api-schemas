// Package types contains the types of the application.
package types

// UserConfigSchemaDeprecationInfo is a struct that contains the deprecation info for a user config schema entry.
type UserConfigSchemaDeprecationInfo struct {
	IsDeprecated      bool   `yaml:"is_deprecated,omitempty"`
	DeprecationNotice string `yaml:"deprecation_notice,omitempty"`
}

// UserConfigSchemaEnumValue is a struct that contains the enum value for a user config schema entry.
type UserConfigSchemaEnumValue struct {
	UserConfigSchemaDeprecationInfo `yaml:",inline"`

	Value interface{} `yaml:"value"`
}

// UserConfigSchema represents an output schema for the user config.
type UserConfigSchema struct {
	UserConfigSchemaDeprecationInfo `yaml:",inline"`

	// `json` tag identifies fields in API responses. Field casing is inconsistent,
	// e.g., minLength (camelCase) and user_error (snake_case).
	Title       string                      `json:"title" yaml:"title,omitempty"`
	Description string                      `json:"description" yaml:"description,omitempty"`
	Type        interface{}                 `json:"type,omitempty" yaml:"type,omitempty"`
	Default     interface{}                 `json:"default,omitempty" yaml:"default,omitempty"`
	Required    []string                    `json:"required" yaml:"required,omitempty"`
	Properties  map[string]UserConfigSchema `json:"properties" yaml:"properties,omitempty"`
	Items       *UserConfigSchema           `json:"items,omitempty" yaml:"items,omitempty"`
	OneOf       []UserConfigSchema          `json:"oneOf" yaml:"one_of,omitempty"`
	Enum        []UserConfigSchemaEnumValue `json:"enum" yaml:"enum,omitempty"`
	Minimum     *float64                    `json:"minimum,omitempty" yaml:"minimum,omitempty"`
	Maximum     *float64                    `json:"maximum,omitempty" yaml:"maximum,omitempty"`
	MinLength   *int                        `json:"minLength,omitempty" yaml:"min_length,omitempty"`
	MaxLength   *int                        `json:"maxLength,omitempty" yaml:"max_length,omitempty"`
	MaxItems    *int                        `json:"maxItems,omitempty" yaml:"max_items,omitempty"`
	CreateOnly  bool                        `json:"createOnly" yaml:"create_only,omitempty"`
	Pattern     string                      `json:"pattern" yaml:"pattern,omitempty"`
	Example     interface{}                 `json:"example,omitempty" yaml:"example,omitempty"`
	UserError   string                      `json:"user_error" yaml:"user_error,omitempty"`
	Secure      bool                        `json:"_secure" yaml:"_secure,omitempty"`
}

// GenerationResult represents the result of a generation.
type GenerationResult map[int]map[string]UserConfigSchema

// ReadResult represents the result of a read.
type ReadResult map[int]map[string]UserConfigSchema

// DiffResult represents the result of a diff.
type DiffResult map[int]map[string]UserConfigSchema

const (
	// KeyServiceTypes is the key for the service types.
	KeyServiceTypes int = iota

	// KeyIntegrationTypes is the key for the integration types.
	KeyIntegrationTypes

	// KeyIntegrationEndpointTypes is the key for the integration endpoint types.
	KeyIntegrationEndpointTypes
)
