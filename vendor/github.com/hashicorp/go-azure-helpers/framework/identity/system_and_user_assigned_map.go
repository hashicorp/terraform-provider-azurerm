// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package identity

import (
	"context"

	"github.com/hashicorp/go-azure-helpers/framework/typehelpers"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func ExpandToSystemAndUserAssignedMap(ctx context.Context, input typehelpers.ListNestedObjectValueOf[IdentityModel], result *identity.SystemAndUserAssignedMap, diags *diag.Diagnostics) {
	if result == nil {
		diags.AddError("Expanding identity", "could not expand identity as target was a nil pointer")
		return
	}

	if input.IsNull() || input.IsUnknown() || len(input.Elements()) == 0 {
		result.Type = identity.TypeNone
		result.IdentityIds = nil
		result.PrincipalId = ""
		result.TenantId = ""

		return
	}

	identityList := make([]IdentityModel, len(input.Elements()))

	d := input.ElementsAs(ctx, &identityList, true)
	if d.HasError() {
		diags.Append(d...)
		return
	}

	if len(identityList) == 1 {
		ident := identityList[0]

		res := identity.SystemAndUserAssignedMap{}

		res.Type = identity.Type(ident.Type.ValueString())
		res.PrincipalId = ident.PrincipalID.ValueString()
		res.TenantId = ident.TenantID.ValueString()

		// convert identities from list to map construct
		identities := map[string]identity.UserAssignedIdentityDetails{}
		idList := make([]string, 0)
		ident.IdentityIDs.ElementsAs(ctx, &idList, false)

		for _, id := range idList {
			identities[id] = identity.UserAssignedIdentityDetails{}
		}

		res.IdentityIds = identities
		*result = res
	}
}

func FlattenFromSystemAndUserAssignedMap(ctx context.Context, input *identity.SystemAndUserAssignedMap, result *typehelpers.ListNestedObjectValueOf[IdentityModel], diags *diag.Diagnostics) {
	if input == nil {
		r := typehelpers.NewListNestedObjectValueOfNull[IdentityModel](ctx)
		*result = r

		return
	}

	i := *input

	ident := IdentityModel{
		Type:        types.StringValue(string(i.Type)),
		PrincipalID: types.StringValue(i.PrincipalId),
		TenantID:    types.StringValue(i.TenantId),
	}

	if len(i.IdentityIds) > 0 {
		ids := make([]attr.Value, 0)
		for id := range i.IdentityIds {
			ids = append(ids, types.StringValue(id))
		}

		ident.IdentityIDs, *diags = typehelpers.NewSetValueOf[types.String](ctx, ids)
		if diags.HasError() {
			return
		}
	} else {
		ident.IdentityIDs = typehelpers.NewSetValueOfNull[types.String](ctx)
	}

	r, d := typehelpers.NewListNestedObjectValueOfValueSlice(ctx, []IdentityModel{ident})
	if d.HasError() {
		diags.Append(d...)
		return
	}

	*result = r
}
