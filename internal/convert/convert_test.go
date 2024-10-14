// Package convert is the package that contains the convert functionality.
package convert

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/stretchr/testify/assert"

	"github.com/aiven/go-api-schemas/internal/pkg/types"
	"github.com/aiven/go-api-schemas/internal/pkg/util"
)

// TestUserConfigSchema tests the UserConfigSchema function.
// nolint:funlen,lll // This function is long, but it's a test function.
// // These lines are long, but they're test data.
func TestUserConfigSchema(t *testing.T) {
	type args struct {
		v map[string]any
	}

	tests := []struct {
		name    string
		args    args
		want    *types.UserConfigSchema
		wantErr error
	}{
		{
			name: "basic",
			args: args{
				v: map[string]any{
					"title":                "",
					"description":          "",
					"type":                 nil,
					"default":              nil,
					"required":             nil,
					"properties":           nil,
					"additionalProperties": nil,
					"items":                nil,
					"oneOf":                nil,
					"enum":                 nil,
					"minimum":              nil,
					"maximum":              nil,
					"minLength":            nil,
					"maxLength":            nil,
					"maxItems":             nil,
					"createOnly":           false,
					"pattern":              "",
					"example":              nil,
					"user_error":           "",
				},
			},
			want: &types.UserConfigSchema{
				UserConfigSchemaDeprecationInfo: types.UserConfigSchemaDeprecationInfo{},
				Title:                           "",
				Description:                     "",
				Type:                            nil,
				Default:                         nil,
				Required:                        nil,
				Properties:                      nil,
				Items:                           nil,
				OneOf:                           nil,
				Enum:                            nil,
				Minimum:                         nil,
				Maximum:                         nil,
				MinLength:                       nil,
				MaxLength:                       nil,
				MaxItems:                        nil,
				CreateOnly:                      false,
				Pattern:                         "",
				Example:                         nil,
				UserError:                       "",
			},
			wantErr: nil,
		},
		{
			name: "filled",
			args: args{
				v: map[string]any{
					"title":       "",
					"description": "",
					"type":        "object",
					"default":     nil,
					"required": []string{
						"datadog_api_key",
					},
					"properties": map[string]any{
						"max_partition_contexts": map[string]any{
							"title":                "Maximum number of partition contexts to send",
							"description":          "",
							"type":                 "integer",
							"default":              nil,
							"required":             nil,
							"properties":           nil,
							"additionalProperties": nil,
							"items":                nil,
							"oneOf":                nil,
							"enum":                 nil,
							"minimum":              200.0,
							"maximum":              200000.0,
							"minLength":            nil,
							"maxLength":            nil,
							"maxItems":             nil,
							"createOnly":           false,
							"pattern":              "",
							"example":              "32000",
							"user_error":           "",
						},
						"site": map[string]any{
							"title":                "Datadog intake site. Defaults to datadoghq.com",
							"description":          "",
							"type":                 "string",
							"default":              nil,
							"required":             nil,
							"properties":           nil,
							"additionalProperties": nil,
							"items":                nil,
							"oneOf":                nil,
							"enum": []interface{}{
								"datadoghq.com",
								"datadoghq.eu",
							},
							"minimum":    nil,
							"maximum":    nil,
							"minLength":  nil,
							"maxLength":  nil,
							"maxItems":   nil,
							"createOnly": false,
							"pattern":    "",
							"example":    "datadoghq.com",
							"user_error": "",
						},
						"datadog_api_key": map[string]any{
							"title":                "Datadog API key",
							"description":          "",
							"type":                 "string",
							"default":              nil,
							"required":             nil,
							"properties":           nil,
							"additionalProperties": nil,
							"items":                nil,
							"oneOf":                nil,
							"enum":                 nil,
							"minimum":              nil,
							"maximum":              nil,
							"minLength":            nil,
							"maxLength":            nil,
							"maxItems":             nil,
							"createOnly":           false,
							"pattern":              "^[A-Za-z0-9]{32}$",
							"example":              "848f30907c15c55d601fe45487cce9b6",
							"user_error":           "Must consist of alpha-numeric characters and contain 32 characters",
						},
						"datadog_tags": map[string]any{
							"title":                "Custom tags provided by user",
							"description":          "",
							"type":                 "array",
							"default":              nil,
							"required":             nil,
							"properties":           nil,
							"additionalProperties": nil,
							"items":                nil,
							"oneOf":                nil,
							"enum":                 nil,
							"minimum":              nil,
							"maximum":              nil,
							"minLength":            nil,
							"maxLength":            nil,
							"maxItems":             nil,
							"createOnly":           false,
							"pattern":              "",
							"example": []interface{}{
								map[string]interface{}{"tag": "foo"},
								map[string]interface{}{
									"comment": "Useful tag",
									"tag":     "bar:buzz",
								},
							},
							"user_error": "",
						},
						"disable_consumer_stats": map[string]any{
							"title":                "Disable consumer group metrics",
							"description":          "",
							"type":                 "boolean",
							"default":              nil,
							"required":             nil,
							"properties":           nil,
							"additionalProperties": nil,
							"items":                nil,
							"oneOf":                nil,
							"enum":                 nil,
							"minimum":              nil,
							"maximum":              nil,
							"minLength":            nil,
							"maxLength":            nil,
							"maxItems":             nil,
							"createOnly":           false,
							"pattern":              "",
							"example":              true,
							"user_error":           "",
						},
						"kafka_consumer_check_instances": map[string]any{
							"title":                "Number of separate instances to fetch kafka consumer statistics with",
							"description":          "",
							"type":                 "integer",
							"default":              nil,
							"required":             nil,
							"properties":           nil,
							"additionalProperties": nil,
							"items":                nil,
							"oneOf":                nil,
							"enum":                 nil,
							"minimum":              1.0,
							"maximum":              100.0,
							"minLength":            nil,
							"maxLength":            nil,
							"maxItems":             nil,
							"createOnly":           false,
							"pattern":              "",
							"example":              "8",
							"user_error":           "",
						},
						"kafka_consumer_stats_timeout": map[string]any{
							"title":                "Number of seconds that datadog will wait to get consumer statistics from brokers",
							"description":          "",
							"type":                 "integer",
							"default":              nil,
							"required":             nil,
							"properties":           nil,
							"additionalProperties": nil,
							"items":                nil,
							"oneOf":                nil,
							"enum":                 nil,
							"minimum":              2.0,
							"maximum":              600.0,
							"minLength":            nil,
							"maxLength":            nil,
							"maxItems":             nil,
							"createOnly":           false,
							"pattern":              "",
							"example":              "60",
							"user_error":           "",
						},
					},
					"additionalProperties": false,
					"items":                nil,
					"oneOf":                nil,
					"enum":                 nil,
					"minimum":              nil,
					"maximum":              nil,
					"minLength":            nil,
					"maxLength":            nil,
					"maxItems":             nil,
					"createOnly":           false,
					"pattern":              "",
					"example":              nil,
					"user_error":           "",
				},
			},
			want: &types.UserConfigSchema{
				UserConfigSchemaDeprecationInfo: types.UserConfigSchemaDeprecationInfo{},
				Title:                           "",
				Description:                     "",
				Type:                            "object",
				Default:                         nil,
				Required: []string{
					"datadog_api_key",
				},
				Properties: map[string]types.UserConfigSchema{
					"datadog_api_key": {
						UserConfigSchemaDeprecationInfo: types.UserConfigSchemaDeprecationInfo{},
						Title:                           "Datadog API key",
						Description:                     "",
						Type:                            "string",
						Default:                         nil,
						Required:                        nil,
						Properties:                      nil,
						Items:                           nil,
						OneOf:                           nil,
						Enum:                            nil,
						Minimum:                         nil,
						Maximum:                         nil,
						MinLength:                       nil,
						MaxLength:                       nil,
						MaxItems:                        nil,
						CreateOnly:                      false,
						Pattern:                         "^[A-Za-z0-9]{32}$",
						Example:                         "848f30907c15c55d601fe45487cce9b6",
						UserError:                       "Must consist of alpha-numeric characters and contain 32 characters",
					},
					"datadog_tags": {
						UserConfigSchemaDeprecationInfo: types.UserConfigSchemaDeprecationInfo{},
						Title:                           "Custom tags provided by user",
						Description:                     "",
						Type:                            "array",
						Default:                         nil,
						Required:                        nil,
						Properties:                      nil,
						Items:                           nil,
						OneOf:                           nil,
						Enum:                            nil,
						Minimum:                         nil,
						Maximum:                         nil,
						MinLength:                       nil,
						MaxLength:                       nil,
						MaxItems:                        nil,
						CreateOnly:                      false,
						Pattern:                         "",
						Example: []interface{}{
							map[string]interface{}{"tag": "foo"},
							map[string]interface{}{
								"comment": "Useful tag",
								"tag":     "bar:buzz",
							},
						},
						UserError: "",
					},
					"disable_consumer_stats": {
						UserConfigSchemaDeprecationInfo: types.UserConfigSchemaDeprecationInfo{},
						Title:                           "Disable consumer group metrics",
						Description:                     "",
						Type:                            "boolean",
						Default:                         nil,
						Required:                        nil,
						Properties:                      nil,
						Items:                           nil,
						OneOf:                           nil,
						Enum:                            nil,
						Minimum:                         nil,
						Maximum:                         nil,
						MinLength:                       nil,
						MaxLength:                       nil,
						MaxItems:                        nil,
						CreateOnly:                      false,
						Pattern:                         "",
						Example:                         true,
						UserError:                       "",
					},
					"kafka_consumer_check_instances": {
						UserConfigSchemaDeprecationInfo: types.UserConfigSchemaDeprecationInfo{},
						Title:                           "Number of separate instances to fetch kafka consumer statistics with",
						Description:                     "",
						Type:                            "integer",
						Default:                         nil,
						Required:                        nil,
						Properties:                      nil,
						Items:                           nil,
						OneOf:                           nil,
						Enum:                            nil,
						Minimum:                         util.Ref(1.0),
						Maximum:                         util.Ref(100.0),
						MinLength:                       nil,
						MaxLength:                       nil,
						MaxItems:                        nil,
						CreateOnly:                      false,
						Pattern:                         "",
						Example:                         "8",
						UserError:                       "",
					},
					"kafka_consumer_stats_timeout": {
						UserConfigSchemaDeprecationInfo: types.UserConfigSchemaDeprecationInfo{},
						Title:                           "Number of seconds that datadog will wait to get consumer statistics from brokers",
						Description:                     "",
						Type:                            "integer",
						Default:                         nil,
						Required:                        nil,
						Properties:                      nil,
						Items:                           nil,
						OneOf:                           nil,
						Enum:                            nil,
						Minimum:                         util.Ref(2.0),
						Maximum:                         util.Ref(600.0),
						MinLength:                       nil,
						MaxLength:                       nil,
						MaxItems:                        nil,
						CreateOnly:                      false,
						Pattern:                         "",
						Example:                         "60",
						UserError:                       "",
					},
					"max_partition_contexts": {
						UserConfigSchemaDeprecationInfo: types.UserConfigSchemaDeprecationInfo{},
						Title:                           "Maximum number of partition contexts to send",
						Description:                     "",
						Type:                            "integer",
						Default:                         nil,
						Required:                        nil,
						Properties:                      nil,
						Items:                           nil,
						OneOf:                           nil,
						Enum:                            nil,
						Minimum:                         util.Ref(200.0),
						Maximum:                         util.Ref(200000.0),
						MinLength:                       nil,
						MaxLength:                       nil,
						MaxItems:                        nil,
						CreateOnly:                      false,
						Pattern:                         "",
						Example:                         "32000",
						UserError:                       "",
					},
					"site": {
						UserConfigSchemaDeprecationInfo: types.UserConfigSchemaDeprecationInfo{},
						Title:                           "Datadog intake site. Defaults to datadoghq.com",
						Description:                     "",
						Type:                            "string",
						Default:                         nil,
						Required:                        nil,
						Properties:                      nil,
						Items:                           nil,
						OneOf:                           nil,
						Enum: []types.UserConfigSchemaEnumValue{
							{Value: "datadoghq.com"},
							{Value: "datadoghq.eu"},
						},
						Minimum:    nil,
						Maximum:    nil,
						MinLength:  nil,
						MaxLength:  nil,
						MaxItems:   nil,
						CreateOnly: false,
						Pattern:    "",
						Example:    "datadoghq.com",
						UserError:  "",
					},
				},
				Items:      nil,
				OneOf:      nil,
				Enum:       nil,
				Minimum:    nil,
				Maximum:    nil,
				MinLength:  nil,
				MaxLength:  nil,
				MaxItems:   nil,
				CreateOnly: false,
				Pattern:    "",
				Example:    nil,
				UserError:  "",
			},
			wantErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := UserConfigSchema(tt.args.v)
			if !cmp.Equal(err, tt.wantErr) {
				t.Errorf("UserConfigSchema() error = %v, wantErr %v", err, tt.wantErr)

				return
			}

			if !cmp.Equal(got, tt.want) {
				t.Errorf(cmp.Diff(tt.want, got))
			}
		})
	}
}

type args struct {
	v map[string]any
}

// Helper function to create a test case for TestUserConfigSchemaErrors
func createTestCase(name string, specificFields map[string]interface{}, wantErr string) struct {
	name    string
	args    args
	wantErr string
} {
	// Shared fields
	fields := map[string]interface{}{
		"title":       "Test",
		"description": "Test description",
	}

	// Override shared fields with specific fields
	for k, v := range specificFields {
		fields[k] = v
	}

	return struct {
		name    string
		args    args
		wantErr string
	}{
		name: name,
		args: args{
			v: fields,
		},
		wantErr: wantErr,
	}
}

// TestUserConfigSchemaErrors tests the UserConfigSchema function for error cases.
// nolint:funlen,lll // This function is long, but it's a test function.
func TestUserConfigSchemaErrors(t *testing.T) {
	type testCase struct {
		name    string
		args    args
		wantErr string
	}

	tests := []testCase{
		createTestCase("invalid property item type", map[string]interface{}{
			"type": "object",
			"properties": map[string]interface{}{
				"invalid_property": "invalid",
			},
		}, "error converting property invalid_property: expected map[string]interface{}, got string"),
		createTestCase("invalid oneOf item type", map[string]interface{}{
			"type": "object",
			"oneOf": []interface{}{
				"invalid",
			},
		}, "error converting slice item at index 0: expected map[string]interface{}, got string"),
		createTestCase("error converting oneOf item", map[string]interface{}{
			"type": "object",
			"oneOf": []interface{}{
				map[string]interface{}{
					"type": "object",
					"properties": map[string]interface{}{
						"invalid_nested_oneOf_item": "invalid",
					},
				},
			},
		}, "error converting slice item at index 0: error converting property invalid_nested_oneOf_item: expected map[string]interface{}, got string"),
		createTestCase("invalid maximum type", map[string]interface{}{
			"type":    "integer",
			"maximum": "invalid",
		}, "error converting maximum: expected float64, got string"),
		createTestCase("error converting nested property", map[string]interface{}{
			"type": "object",
			"properties": map[string]interface{}{
				"nested_property": map[string]interface{}{
					"type": "object",
					"properties": map[string]interface{}{
						"invalid_nested_property": "invalid",
					},
				},
			},
		}, "error converting property nested_property: error converting property invalid_nested_property: expected map[string]interface{}, got string"),
		createTestCase("invalid required type", map[string]interface{}{
			"type":     "object",
			"required": "invalid",
		}, "error converting required fields: expected []interface{} or []string, got string"),
		createTestCase("invalid enum type", map[string]interface{}{
			"type": "string",
			"enum": []interface{}{
				map[string]interface{}{},
			},
		}, "error converting enum value: expected string or number, got map[string]interface {}"),
		createTestCase("invalid minLength type", map[string]interface{}{
			"type":      "string",
			"minLength": "invalid",
		}, "error converting minLength: expected float64, got string"),
		createTestCase("invalid pattern type", map[string]interface{}{
			"type":    "string",
			"pattern": 123,
		}, "error converting pattern: expected string, got int"),
		createTestCase("invalid secure type", map[string]interface{}{
			"type":    "string",
			"_secure": "invalid",
		}, "error converting _secure: expected bool, got string"),
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := UserConfigSchema(tt.args.v)
			if err != nil {
				assert.EqualError(t, err, tt.wantErr)
			} else {
				t.Errorf("expected error but got none")
			}
		})
	}
}
