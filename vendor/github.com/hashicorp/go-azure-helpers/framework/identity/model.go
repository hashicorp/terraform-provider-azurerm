package identity

import (
	"encoding/json"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type Type string

const (
	TypeNone                       Type = "None"
	TypeSystemAssigned             Type = "SystemAssigned"
	TypeUserAssigned               Type = "UserAssigned"
	TypeSystemAssignedUserAssigned Type = "SystemAssigned, UserAssigned"

	// this is an internal-only type to transform the legacy API value to the type we want to expose
	typeLegacySystemAssignedUserAssigned Type = "SystemAssigned,UserAssigned"
)

type ModelUserAssigned struct {
	Type        types.String   `tfsdk:"type"`
	IdentityIds types.ListType `tfsdk:"identity_ids"`
}

var ModelUserAssignedAttr = map[string]attr.Type{
	"type":         types.StringType,
	"identity_ids": types.ListType{}.WithElementType(types.StringType),
}

type ModelSystemAssigned struct {
	Type        types.String `tfsdk:"type"`
	PrincipalId types.String `tfsdk:"principal_id"`
	TenantId    types.String `tfsdk:"tenant_id"`
}

var ModelSystemAssignedAttr = map[string]attr.Type{
	"type":         types.StringType,
	"principal_id": types.StringType,
	"tenant_id":    types.StringType,
}

type ModelSystemAssignedUserAssigned struct {
	Type        types.String `tfsdk:"type"`
	PrincipalId types.String `tfsdk:"principal_id"`
	TenantId    types.String `tfsdk:"tenant_id"`
	IdentityIds types.List   `tfsdk:"identity_ids"`
}

var ModelSystemAssignedUserAssignedAttr = map[string]attr.Type{
	"type":         types.StringType,
	"principal_id": types.StringType,
	"tenant_id":    types.StringType,
	"identity_ids": types.ListType{}.WithElementType(types.StringType),
}

var _ json.Marshaler = UserAssignedIdentityDetails{}

type UserAssignedIdentityDetails struct {
	ClientId    *string `json:"clientId,omitempty"`
	PrincipalId *string `json:"principalId,omitempty"`
}

func (u UserAssignedIdentityDetails) MarshalJSON() ([]byte, error) {
	// none of these properties can be set, so we'll just flatten an empty struct
	return json.Marshal(map[string]interface{}{})
}
