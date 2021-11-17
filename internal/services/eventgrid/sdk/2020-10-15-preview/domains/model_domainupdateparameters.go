package domains

import (
	"github.com/hashicorp/terraform-provider-azurerm/internal/identity"
)

type DomainUpdateParameters struct {
	Identity   *identity.SystemUserAssignedIdentityMap `json:"identity,omitempty"`
	Properties *DomainUpdateParameterProperties        `json:"properties,omitempty"`
	Sku        *ResourceSku                            `json:"sku,omitempty"`
	Tags       *map[string]string                      `json:"tags,omitempty"`
}
