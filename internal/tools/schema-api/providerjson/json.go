package providerjson

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-azurerm/internal/provider"
)

type ProviderJSON schema.Provider

type SchemaJSON struct {
	Type        string      `json:"type,omitempty"`
	ConfigMode  string      `json:"config_mode,omitempty"`
	Optional    bool        `json:"optional,omitempty"`
	Required    bool        `json:"required,omitempty"`
	Default     interface{} `json:"default,omitempty"`
	Description string      `json:"description,omitempty"`
	Computed    bool        `json:"computed,omitempty"`
	ForceNew    bool        `json:"force_new,omitempty"`
	Elem        interface{} `json:"elem,omitempty"`
	MaxItems    int         `json:"max_items,omitempty"`
	MinItems    int         `json:"min_items,omitempty"`
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
