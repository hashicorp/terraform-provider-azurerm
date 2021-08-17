package account

import (
	"github.com/hashicorp/terraform-provider-azurerm/internal/identity"
)

type Account struct {
	Id         *string                          `json:"id,omitempty"`
	Identity   *identity.SystemAssignedIdentity `json:"identity,omitempty"`
	Location   *string                          `json:"location,omitempty"`
	Name       *string                          `json:"name,omitempty"`
	Properties *AccountProperties               `json:"properties,omitempty"`
	Sku        *AccountSku                      `json:"sku,omitempty"`
	Tags       *map[string]string               `json:"tags,omitempty"`
	Type       *string                          `json:"type,omitempty"`
}
