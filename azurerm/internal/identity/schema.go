package identity

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

const none = "None"
const systemAssigned = "SystemAssigned"
const systemAssignedUserAssigned = "SystemAssigned, UserAssigned"
const userAssigned = "UserAssigned"

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
	Schema() *schema.Schema
}
