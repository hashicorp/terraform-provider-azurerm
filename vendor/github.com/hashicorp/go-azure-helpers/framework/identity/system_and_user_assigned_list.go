// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package identity

import (
	"context"

	"github.com/hashicorp/go-azure-helpers/framework/convert"
	"github.com/hashicorp/go-azure-helpers/framework/typehelpers"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func ExpandToSystemAndUserAssignedList(ctx context.Context, input typehelpers.ListNestedObjectValueOf[IdentityModel], result *identity.SystemAndUserAssignedList, diags *diag.Diagnostics) {
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
		convert.Expand(ctx, ident, result, diags)
	}
}

func FlattenFromSystemAndUserAssignedList(ctx context.Context, input *identity.SystemAndUserAssignedList, result *typehelpers.ListNestedObjectValueOf[IdentityModel], diags *diag.Diagnostics) {
	if input == nil {
		r := typehelpers.NewListNestedObjectValueOfNull[IdentityModel](ctx)
		*result = r

		return
	}

	flat := IdentityModel{
		IdentityIDs: typehelpers.NewSetValueOfNull[types.String](ctx),
	}

	convert.Flatten(ctx, input, &flat, diags)
	list, d := typehelpers.NewListNestedObjectValueOfValueSlice(ctx, []IdentityModel{flat})
	if d.HasError() {
		diags.Append(d...)
		return
	}

	*result = list
}
