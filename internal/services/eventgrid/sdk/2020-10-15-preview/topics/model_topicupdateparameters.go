package topics

import (
	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
)

type TopicUpdateParameters struct {
	Identity   *identity.SystemUserAssignedList `json:"identity,omitempty"`
	Properties *TopicUpdateParameterProperties  `json:"properties,omitempty"`
	Sku        *ResourceSku                     `json:"sku,omitempty"`
	Tags       *map[string]string               `json:"tags,omitempty"`
}
