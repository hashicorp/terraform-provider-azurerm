package namespaces

import (
	"github.com/hashicorp/terraform-provider-azurerm/internal/identity"
)

type EHNamespace struct {
	Id         *string                          `json:"id,omitempty"`
	Identity   *identity.SystemAssignedIdentity `json:"identity,omitempty"`
	Location   *string                          `json:"location,omitempty"`
	Name       *string                          `json:"name,omitempty"`
	Properties *EHNamespaceProperties           `json:"properties,omitempty"`
	Sku        *Sku                             `json:"sku,omitempty"`
	Tags       *map[string]string               `json:"tags,omitempty"`
	Type       *string                          `json:"type,omitempty"`
}
