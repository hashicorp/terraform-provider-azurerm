// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package providerjson

import (
	"encoding/json"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-azurerm/internal/provider"
)

const (
	SchemaTypeSet    = "TypeSet"
	SchemaTypeList   = "TypeList"
	SchemaTypeInt    = "TypeInt"
	SchemaTypeString = "String"
	SchemaTypeBool   = "Bool"
	SchemaTypeFloat  = "Float"
)

type ProviderJSON schema.Provider

type SchemaJSON struct {
	Type        string      `json:"type,omitempty"` // TODO - Needs to be interface{}
	ConfigMode  string      `json:"configMode,omitempty"`
	Optional    bool        `json:"optional,omitempty"`
	Required    bool        `json:"required,omitempty"`
	Default     interface{} `json:"default,omitempty"`
	Description string      `json:"description,omitempty"`
	Computed    bool        `json:"computed,omitempty"`
	ForceNew    bool        `json:"forceNew,omitempty"`
	Elem        interface{} `json:"elem,omitempty"`
	MaxItems    int         `json:"maxItems,omitempty"`
	MinItems    int         `json:"minItems,omitempty"`
}

func (b *SchemaJSON) UnmarshalJSON(body []byte) error {
	var m map[string]interface{}
	err := json.Unmarshal(body, &m)
	if err != nil {
		return err
	}
	b.Type, _ = m["type"].(string)
	b.ConfigMode, _ = m["configMode"].(string)
	b.Optional, _ = m["optional"].(bool)
	b.Required, _ = m["required"].(bool)
	b.Description, _ = m["description"].(string)
	b.Computed, _ = m["computed"].(bool)
	b.ForceNew, _ = m["forceNew"].(bool)
	if max, ok := m["maxItems"].(float64); ok {
		b.MaxItems = int(max)
	}
	if min, ok := m["minItems"].(float64); ok {
		b.MaxItems = int(min)
	}

	if def, ok := m["default"]; ok && def != nil {
		switch def.(type) {
		case string:
			b.Default = def
		case bool:
			b.Default = def
		case int:
			b.Default = def
		case float32:
			b.Default = def
		case float64:
			b.Default = def
		}
	}

	if e, ok := m["elem"]; ok && e != nil {
		elem := e.(map[string]interface{})
		if schema, ok := elem["schema"]; ok {
			b.Elem = ResourceFromMap(schema.(map[string]interface{}))
		}
		if t, ok := elem["type"]; ok {
			b.Elem = t.(string)
		}
	}

	return nil
}

type ResourceJSON struct {
	Schema   map[string]SchemaJSON `json:"schema"`
	Timeouts *ResourceTimeoutJSON  `json:"timeouts,omitempty"`
}

type ResourceTimeoutJSON struct {
	Create int `json:"create,omitempty"`
	Read   int `json:"read,omitempty"`
	Delete int `json:"delete,omitempty"`
	Update int `json:"update,omitempty"`
}

func LoadData() *ProviderJSON {
	p := provider.AzureProvider()
	return (*ProviderJSON)(p)
}

type ProviderSchemaJSON struct {
	Schema         map[string]SchemaJSON   `json:"schema"`
	ResourcesMap   map[string]ResourceJSON `json:"resources,omitempty"`
	DataSourcesMap map[string]ResourceJSON `json:"dataSources,omitempty"`
}

type ProviderWrapper struct {
	ProviderName   string              `json:"providerName"`
	SchemaVersion  string              `json:"schemaVersion"`
	ProviderSchema *ProviderSchemaJSON `json:"providerSchema,omitempty"`
}
