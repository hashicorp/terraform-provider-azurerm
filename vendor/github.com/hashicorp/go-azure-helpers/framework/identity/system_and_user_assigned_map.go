package identity

import (
	"context"
	"encoding/json"
	"fmt"

	rmidentity "github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var _ json.Marshaler = &SystemAndUserAssignedMap{}

type SystemAndUserAssignedMap struct {
	*rmidentity.SystemAndUserAssignedMap
}

func (s *SystemAndUserAssignedMap) MarshalJSON() ([]byte, error) {
	identityType := TypeNone
	userAssignedIdentityIds := map[string]rmidentity.UserAssignedIdentityDetails{}

	switch s.Type {
	case rmidentity.TypeSystemAssigned:
		identityType = TypeSystemAssigned
	case rmidentity.TypeUserAssigned:
		identityType = TypeUserAssigned
	case rmidentity.TypeSystemAssignedUserAssigned:
		identityType = TypeSystemAssignedUserAssigned
	}

	if identityType != TypeNone {
		userAssignedIdentityIds = s.IdentityIds
	}

	out := map[string]interface{}{
		"type":                   string(identityType),
		"userAssignedIdentities": nil,
	}
	if len(userAssignedIdentityIds) > 0 {
		out["userAssignedIdentities"] = userAssignedIdentityIds
	}

	return json.Marshal(out)
}

func ExpandSystemAndUserAssignedMap(input types.List) (result *rmidentity.SystemAndUserAssignedMap, diags diag.Diagnostics) {
	bg := context.Background()
	if input.IsNull() || input.IsUnknown() {
		return &rmidentity.SystemAndUserAssignedMap{
			Type:        rmidentity.TypeNone,
			IdentityIds: nil,
		}, nil
	}

	identities := make([]ModelSystemAssignedUserAssigned, 0)

	diags = input.ElementsAs(bg, &identities, false)
	if diags.HasError() {
		return nil, diags
	}

	id := identities[0]
	ids := make([]string, 0)
	diags = id.IdentityIds.ElementsAs(bg, &ids, false)

	identityIds := make(map[string]rmidentity.UserAssignedIdentityDetails, len(ids))
	for _, v := range ids {
		identityIds[v] = rmidentity.UserAssignedIdentityDetails{
			// intentionally empty since the expand shouldn't send these values
		}
	}
	if len(identityIds) > 0 && (id.Type.ValueString() != string(TypeSystemAssignedUserAssigned) && id.Type.ValueString() != string(TypeUserAssigned)) {
		diags.AddError("identity error", fmt.Sprintf("`identity_ids` can only be specified when `type` is set to %q or %q", TypeSystemAssignedUserAssigned, TypeUserAssigned))
		return nil, diags
	}

	return &rmidentity.SystemAndUserAssignedMap{
		Type:        rmidentity.Type(id.Type.ValueString()),
		IdentityIds: identityIds,
	}, nil
}

func FlattenSystemAndUserAssignedMap(input *rmidentity.SystemAndUserAssignedMap) (result types.List, diags diag.Diagnostics) {
	bg := context.Background()
	diags = make([]diag.Diagnostic, 0)

	if input == nil {
		result = types.ListNull(types.ObjectType{}.WithAttributeTypes(IdentityAttr))
		return
	}

	return types.ListValueFrom(bg, types.ObjectType{}.WithAttributeTypes(IdentityAttr), input)
}
