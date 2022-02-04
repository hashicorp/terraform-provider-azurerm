package cognitiveservicesaccounts

import (
	"github.com/hashicorp/terraform-provider-azurerm/internal/identity"
)

type Account struct {
	Etag       *string                                 `json:"etag,omitempty"`
	Id         *string                                 `json:"id,omitempty"`
	Identity   *identity.SystemUserAssignedIdentityMap `json:"identity,omitempty"`
	Kind       *string                                 `json:"kind,omitempty"`
	Location   *string                                 `json:"location,omitempty"`
	Name       *string                                 `json:"name,omitempty"`
	Properties *AccountProperties                      `json:"properties,omitempty"`
	Sku        *Sku                                    `json:"sku,omitempty"`
	SystemData *SystemData                             `json:"systemData,omitempty"`
	Tags       *map[string]string                      `json:"tags,omitempty"`
	Type       *string                                 `json:"type,omitempty"`
}
