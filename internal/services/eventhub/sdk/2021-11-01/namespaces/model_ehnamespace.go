package namespaces

import (
	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
)

type EHNamespace struct {
	Id         *string                            `json:"id,omitempty"`
	Identity   *identity.SystemAndUserAssignedMap `json:"identity,omitempty"`
	Location   *string                            `json:"location,omitempty"`
	Name       *string                            `json:"name,omitempty"`
	Properties *EHNamespaceProperties             `json:"properties,omitempty"`
	Sku        *Sku                               `json:"sku,omitempty"`
	SystemData *SystemData                        `json:"systemData,omitempty"`
	Tags       *map[string]string                 `json:"tags,omitempty"`
	Type       *string                            `json:"type,omitempty"`
}
