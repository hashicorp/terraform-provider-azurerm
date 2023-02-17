package providerjson

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-azurerm/internal/provider"
)

type ProviderJSON schema.Provider

type SchemaJSON struct {
	Type        string      `json:"type,omitempty"`
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
	ProviderSchema *ProviderSchemaJSON `json:"providerSchema,omitempty"`
}
