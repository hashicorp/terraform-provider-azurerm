package domains

import (
	"github.com/hashicorp/terraform-provider-azurerm/internal/identity"
)

type Domain struct {
	Id         *string                                 `json:"id,omitempty"`
	Identity   *identity.SystemUserAssignedIdentityMap `json:"identity,omitempty"`
	Location   string                                  `json:"location"`
	Name       *string                                 `json:"name,omitempty"`
	Properties *DomainProperties                       `json:"properties,omitempty"`
	Sku        *ResourceSku                            `json:"sku,omitempty"`
	SystemData *SystemData                             `json:"systemData,omitempty"`
	Tags       *map[string]string                      `json:"tags,omitempty"`
	Type       *string                                 `json:"type,omitempty"`
}
