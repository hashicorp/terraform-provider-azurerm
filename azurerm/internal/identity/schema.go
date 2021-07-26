package identity

import (
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/pluginsdk"
)

const (
	none                       = "None"
	systemAssigned             = "SystemAssigned"
	userAssigned               = "UserAssigned"
	systemAssignedUserAssigned = "SystemAssigned, UserAssigned"
)

type ExpandedConfig struct {
	// Type is the type of User Assigned Identity, either `None`, `SystemAssigned`, `UserAssigned`
	// or `SystemAssigned, UserAssigned`
	Type                    string
	PrincipalId             *string
	TenantId                *string
	UserAssignedIdentityIds *[]string
}

type Identity interface {
	Expand(input []interface{}) (*ExpandedConfig, error)
	Flatten(input *ExpandedConfig) []interface{}
	Schema() *pluginsdk.Schema
}
