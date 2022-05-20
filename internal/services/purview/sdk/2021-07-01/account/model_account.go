package account

import (
	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
)

type Account struct {
	Id         *string                           `json:"id,omitempty"`
	Identity   *identity.SystemOrUserAssignedMap `json:"identity,omitempty"`
	Location   *string                           `json:"location,omitempty"`
	Name       *string                           `json:"name,omitempty"`
	Properties *AccountProperties                `json:"properties,omitempty"`
	Sku        *AccountSku                       `json:"sku,omitempty"`
	SystemData *SystemData                       `json:"systemData,omitempty"`
	Tags       *map[string]string                `json:"tags,omitempty"`
	Type       *string                           `json:"type,omitempty"`
}
