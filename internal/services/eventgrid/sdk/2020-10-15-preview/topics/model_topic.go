package topics

import (
	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
)

type Topic struct {
	ExtendedLocation *ExtendedLocation                `json:"extendedLocation,omitempty"`
	Id               *string                          `json:"id,omitempty"`
	Identity         *identity.SystemUserAssignedList `json:"identity,omitempty"`
	Kind             *ResourceKind                    `json:"kind,omitempty"`
	Location         string                           `json:"location"`
	Name             *string                          `json:"name,omitempty"`
	Properties       *TopicProperties                 `json:"properties,omitempty"`
	Sku              *ResourceSku                     `json:"sku,omitempty"`
	SystemData       *SystemData                      `json:"systemData,omitempty"`
	Tags             *map[string]string               `json:"tags,omitempty"`
	Type             *string                          `json:"type,omitempty"`
}
