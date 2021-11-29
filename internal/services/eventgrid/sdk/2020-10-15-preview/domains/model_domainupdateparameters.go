package domains

import (
	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
)

type DomainUpdateParameters struct {
	Identity   *identity.SystemUserAssignedList `json:"identity,omitempty"`
	Properties *DomainUpdateParameterProperties `json:"properties,omitempty"`
	Sku        *ResourceSku                     `json:"sku,omitempty"`
	Tags       *map[string]string               `json:"tags,omitempty"`
}
