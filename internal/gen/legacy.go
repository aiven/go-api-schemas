package gen

import (
	"fmt"
)

type legacyDoc struct {
	ServiceTypes     map[string]legacySchema `json:"service_types"`
	IntegrationTypes []struct {
		legacySchema
		Key string `json:"integration_type"`
	} `json:"integration_types"`
	IntegrationEndpointTypes []struct {
		legacySchema
		Key string `json:"endpoint_type"`
	} `json:"endpoint_types"`
}

type legacySchema struct {
	UserConfigSchema *schema `json:"user_config_schema"`
}

// legacyToComponents adapts legacy schemas to the new format.
func legacyToComponents(d *doc) {
	if d.Components.Schemas == nil {
		d.Components.Schemas = make(map[string]*schema)
	}

	for k, v := range d.ServiceTypes {
		key := fmt.Sprintf("Service%sUserConfig", k)
		d.Components.Schemas[key] = v.UserConfigSchema
	}

	for _, v := range d.IntegrationTypes {
		key := fmt.Sprintf("Integration%sUserConfig", v.Key)
		d.Components.Schemas[key] = v.UserConfigSchema
	}

	for _, v := range d.IntegrationEndpointTypes {
		key := fmt.Sprintf("IntegrationEndpoint%sUserConfig", v.Key)
		d.Components.Schemas[key] = v.UserConfigSchema
	}
}
