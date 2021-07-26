package identity

import (
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/pluginsdk"
)

type identityType string

const (
	none                       identityType = "None"
	systemAssigned             identityType = "SystemAssigned"
	userAssigned               identityType = "UserAssigned"
	systemAssignedUserAssigned identityType = "SystemAssigned, UserAssigned"
)

type ExpandedConfig struct {
	// Type is the type of User Assigned Identity, either `None`, `SystemAssigned`, `UserAssigned`
	// or `SystemAssigned, UserAssigned`
	Type                    identityType
	PrincipalId             *string
	TenantId                *string
	UserAssignedIdentityIds *[]string
}

type Identity interface {
	Expand(input []interface{}) (*ExpandedConfig, error)
	Flatten(input *ExpandedConfig) []interface{}
	Schema() *pluginsdk.Schema
}
