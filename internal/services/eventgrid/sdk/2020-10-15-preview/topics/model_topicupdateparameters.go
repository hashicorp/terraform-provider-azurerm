package topics

import (
	"github.com/hashicorp/terraform-provider-azurerm/internal/identity"
)

type TopicUpdateParameters struct {
	Identity   *identity.SystemUserAssignedIdentityMap `json:"identity,omitempty"`
	Properties *TopicUpdateParameterProperties         `json:"properties,omitempty"`
	Sku        *ResourceSku                            `json:"sku,omitempty"`
	Tags       *map[string]string                      `json:"tags,omitempty"`
}
