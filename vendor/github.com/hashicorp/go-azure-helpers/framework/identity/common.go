// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package identity

import (
	"context"

	"github.com/hashicorp/go-azure-helpers/framework/typehelpers"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
	"github.com/hashicorp/terraform-plugin-framework-validators/listvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	datasourceschema "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type IdentityModel struct {
	Type        types.String                         `tfsdk:"type"`
	IdentityIDs typehelpers.SetValueOf[types.String] `tfsdk:"identity_ids" convert:"IdentityIds"`
	PrincipalID types.String                         `tfsdk:"principal_id" convert:"PrincipalId"`
	TenantID    types.String                         `tfsdk:"tenant_id" convert:"TenantId"`
}

type SystemIdentityModel struct {
	Type        types.String `tfsdk:"type"`
	PrincipalID types.String `tfsdk:"principal_id" convert:"PrincipalId"`
	TenantID    types.String `tfsdk:"tenant_id" convert:"TenantId"`
}

// IdentityResourceAttributeSchema returns a Framework resource attribute schema
// intended for use with TF ProtoV6
func IdentityResourceAttributeSchema(ctx context.Context, validTypes ...identity.Type) schema.ListNestedAttribute {
	typeValidatorList := buildTypesValidatorList(validTypes)

	return schema.ListNestedAttribute{
		CustomType: typehelpers.NewListNestedObjectTypeOf[IdentityModel](ctx),
		NestedObject: schema.NestedAttributeObject{
			Attributes: map[string]schema.Attribute{
				"type": schema.StringAttribute{
					Required: true,
					Validators: []validator.String{
						stringvalidator.OneOf(
							typeValidatorList...,
						),
					},
				},

				"identity_ids": schema.ListAttribute{
					ElementType: types.StringType,
					Optional:    true,
					Validators: []validator.List{
						listvalidator.ValueStringsAre(
							typehelpers.WrappedStringValidator{
								Func: commonids.ValidateUserAssignedIdentityID,
							},
						),
					},
				},

				"principal_id": schema.StringAttribute{
					Computed: true,
					PlanModifiers: []planmodifier.String{
						stringplanmodifier.UseStateForUnknown(),
					},
				},

				"tenant_id": schema.StringAttribute{
					Computed: true,
					PlanModifiers: []planmodifier.String{
						stringplanmodifier.UseStateForUnknown(),
					},
				},
			},
		},
		Validators: []validator.List{
			listvalidator.SizeAtMost(1),
			// TODO - listvalidator for checking type and identityIds validation
		},
	}
}

// IdentityDataSourceAttributeSchema returns a Framework resource attribute schema
// intended for use with TF ProtoV6
func IdentityDataSourceAttributeSchema(ctx context.Context) schema.ListNestedAttribute {
	return schema.ListNestedAttribute{
		CustomType: typehelpers.NewListNestedObjectTypeOf[IdentityModel](ctx),
		NestedObject: schema.NestedAttributeObject{
			Attributes: map[string]schema.Attribute{
				"type": schema.StringAttribute{
					Computed: true,
				},

				"identity_ids": schema.ListAttribute{
					ElementType: types.StringType,
					Computed:    true,
				},

				"principal_id": schema.StringAttribute{
					Computed: true,
				},

				"tenant_id": schema.StringAttribute{
					Computed: true,
				},
			},
		},
		Validators: []validator.List{
			listvalidator.SizeAtMost(1),
		},
	}
}

// IdentityResourceBlockSchema returns a Framework resource Block schema for use with resources that
// require a Resource Manager Identity block. Valid identity.Type values can be specified as an optional
// []identity.Type{} parameter to the function. omitting or empty list will result "UserAssigned", "SystemAssigned", and
// "SystemAssigned, UserAssigned" being valid
func IdentityResourceBlockSchema(ctx context.Context, validTypes ...identity.Type) schema.ListNestedBlock {
	typeValidatorList := buildTypesValidatorList(validTypes)

	return schema.ListNestedBlock{
		CustomType: typehelpers.NewListNestedObjectTypeOf[IdentityModel](ctx),
		NestedObject: schema.NestedBlockObject{
			Attributes: map[string]schema.Attribute{
				"type": schema.StringAttribute{
					Required: true,
					Validators: []validator.String{
						stringvalidator.OneOf(
							typeValidatorList...,
						),
					},
				},

				"identity_ids": schema.ListAttribute{
					ElementType: types.StringType,
					Optional:    true,
					Validators: []validator.List{
						listvalidator.ValueStringsAre(
							typehelpers.WrappedStringValidator{
								Func: commonids.ValidateUserAssignedIdentityID,
							},
						),
					},
				},

				"principal_id": schema.StringAttribute{
					Computed: true,
					PlanModifiers: []planmodifier.String{
						stringplanmodifier.UseStateForUnknown(),
					},
				},

				"tenant_id": schema.StringAttribute{
					Computed: true,
					PlanModifiers: []planmodifier.String{
						stringplanmodifier.UseStateForUnknown(),
					},
				},
			},
		},
		Validators: []validator.List{
			listvalidator.SizeAtMost(1),
			// TODO - listvalidator for checking type and identityIds validation
		},
	}
}

// IdentityDataSourceBlockSchema returns a Framework datasource Block schema for use with resources that
// require a Resource Manager SystemAndUserAssignedMap
func IdentityDataSourceBlockSchema(ctx context.Context) datasourceschema.ListNestedBlock {
	return datasourceschema.ListNestedBlock{
		CustomType: typehelpers.NewListNestedObjectTypeOf[IdentityModel](ctx),
		NestedObject: datasourceschema.NestedBlockObject{
			Attributes: map[string]datasourceschema.Attribute{
				"type": schema.StringAttribute{
					Computed: true,
				},

				"identity_ids": schema.ListAttribute{
					ElementType: types.StringType,
					Computed:    true,
				},

				"principal_id": schema.StringAttribute{
					Computed: true,
				},

				"tenant_id": schema.StringAttribute{
					Computed: true,
				},
			},
		},
	}
}

func buildTypesValidatorList(input []identity.Type) (result []string) {
	if len(input) == 0 {
		result = []string{
			// string(identity.TypeNone), // TODO - investigate re-introducing allowing explicit "None" config
			string(identity.TypeUserAssigned),
			string(identity.TypeSystemAssigned),
			string(identity.TypeSystemAssignedUserAssigned),
		}
	} else {
		for _, v := range input {
			result = append(result, string(v))
		}
	}

	return
}
