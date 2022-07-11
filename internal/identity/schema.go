package identity

import (
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type Type string

type ExpandedConfig struct {
	// Type is the type of User Assigned Identity, either `None`, `SystemAssigned`, `UserAssigned`
	// or `SystemAssigned, UserAssigned`
	Type                    Type     `tfschema:"type"`
	PrincipalId             string   `tfschema:"principal_id"`
	TenantId                string   `tfschema:"tenant_id"`
	UserAssignedIdentityIds []string `tfschema:"identity_ids"`
}

type Identity interface {
	Expand(input []interface{}) (*ExpandedConfig, error)
	Flatten(input *ExpandedConfig) []interface{}
	Schema() *pluginsdk.Schema
}
