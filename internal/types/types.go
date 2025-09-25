package types

import "encoding/json"

// UserConfigSchemaDeprecationInfo is a struct that contains the deprecation info for a user config schema entry.
type UserConfigSchemaDeprecationInfo struct {
	IsDeprecated      bool   `yaml:"is_deprecated,omitempty"`
	DeprecationNotice string `yaml:"deprecation_notice,omitempty"`
}

// Deprecate sets the deprecation info for a user config schema entry.
func (u *UserConfigSchemaDeprecationInfo) Deprecate(msg string) {
	u.IsDeprecated = true
	u.DeprecationNotice = msg
}

// UserConfigSchemaEnumValue is a struct that contains the enum value for a user config schema entry.
type UserConfigSchemaEnumValue struct {
	UserConfigSchemaDeprecationInfo `yaml:",inline"`

	Value interface{} `yaml:"value"`
}

// UserConfigSchema represents an output schema for the user config.
type UserConfigSchema struct {
	UserConfigSchemaDeprecationInfo `yaml:",inline"`

	Title       string                       `yaml:"title,omitempty"`
	Description string                       `yaml:"description,omitempty"`
	Type        interface{}                  `yaml:"type,omitempty"`
	Default     interface{}                  `yaml:"default,omitempty"`
	Required    []string                     `yaml:"required,omitempty"`
	Properties  map[string]*UserConfigSchema `yaml:"properties,omitempty"`
	Items       *UserConfigSchema            `yaml:"items,omitempty"`
	OneOf       []*UserConfigSchema          `yaml:"one_of,omitempty"`
	Enum        []*UserConfigSchemaEnumValue `yaml:"enum,omitempty"`
	Minimum     *json.Number                 `yaml:"minimum,omitempty"`
	Maximum     *json.Number                 `yaml:"maximum,omitempty"`
	MinLength   *int                         `yaml:"min_length,omitempty"`
	MaxLength   *int                         `yaml:"max_length,omitempty"`
	MaxItems    *int                         `yaml:"max_items,omitempty"`
	CreateOnly  bool                         `yaml:"create_only,omitempty"`
	Pattern     string                       `yaml:"pattern,omitempty"`
	Example     interface{}                  `yaml:"example,omitempty"`
	UserError   string                       `yaml:"user_error,omitempty"`
	Secure      bool                         `yaml:"_secure,omitempty"`
	Nullable    bool                         `yaml:"nullable,omitempty"`
}

type SchemaType int

const (
	ServiceSchemaType SchemaType = iota
	IntegrationSchemaType
	IntegrationEndpointSchemaType
)

func GetSchemaTypes() []SchemaType {
	return []SchemaType{
		ServiceSchemaType,
		IntegrationSchemaType,
		IntegrationEndpointSchemaType,
	}
}

// GenerationResult result of newly generated schemas from an OpenAPI file.
type GenerationResult map[SchemaType]map[string]*UserConfigSchema

// ReadResult the result of GenerationResult that was read from a file.
type ReadResult map[SchemaType]map[string]*UserConfigSchema

// DiffResult the result of comparing ReadResult and GenerationResult.
type DiffResult map[SchemaType]map[string]*UserConfigSchema
