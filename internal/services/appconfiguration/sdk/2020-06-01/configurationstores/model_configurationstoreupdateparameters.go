package configurationstores

import (
	"github.com/hashicorp/terraform-provider-azurerm/internal/identity"
)

type ConfigurationStoreUpdateParameters struct {
	Identity   *identity.SystemUserAssignedIdentityMap       `json:"identity,omitempty"`
	Properties *ConfigurationStorePropertiesUpdateParameters `json:"properties,omitempty"`
	Sku        *Sku                                          `json:"sku,omitempty"`
	Tags       *map[string]string                            `json:"tags,omitempty"`
}
