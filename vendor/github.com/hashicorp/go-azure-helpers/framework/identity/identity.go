package identity

import (
	"github.com/hashicorp/go-azure-helpers/framework/typehelpers"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
	"github.com/hashicorp/terraform-plugin-framework-validators/listvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func IdentitySchemaAttribute() schema.ListNestedAttribute {
	return schema.ListNestedAttribute{
		NestedObject: schema.NestedAttributeObject{
			Attributes: map[string]schema.Attribute{
				"type": schema.StringAttribute{
					Required: true,
					Validators: []validator.String{
						stringvalidator.OneOf(
							string(identity.TypeUserAssigned),
							string(identity.TypeSystemAssigned),
							string(identity.TypeSystemAssignedUserAssigned),
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
					Default: typehelpers.WrappedStringDefault{
						Value: "",
					},
				},
				"tenant_id": schema.StringAttribute{
					Computed: true,
					PlanModifiers: []planmodifier.String{
						stringplanmodifier.UseStateForUnknown(),
					},
					Default: typehelpers.WrappedStringDefault{
						Value: "",
					},
				},
			},
		},
		Validators: []validator.List{
			listvalidator.SizeAtMost(1),
		},
	}
}

func IdentitySchemaBlock() schema.ListNestedBlock {
	return schema.ListNestedBlock{
		NestedObject: schema.NestedBlockObject{
			Attributes: map[string]schema.Attribute{
				"type": schema.StringAttribute{
					Required: true,
					Validators: []validator.String{
						stringvalidator.OneOf(
							string(identity.TypeUserAssigned),
							string(identity.TypeSystemAssigned),
							string(identity.TypeSystemAssignedUserAssigned),
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
					Default: typehelpers.WrappedStringDefault{
						Value: "",
					},
				},
				"tenant_id": schema.StringAttribute{
					Computed: true,
					PlanModifiers: []planmodifier.String{
						stringplanmodifier.UseStateForUnknown(),
					},
					Default: typehelpers.WrappedStringDefault{
						Value: "",
					},
				},
			},
		},
		Validators: []validator.List{
			listvalidator.SizeAtMost(1),
		},
	}
}

type Identity struct {
	Type        types.String `tfsdk:"type"`
	IdentityIDs types.List   `tfsdk:"identity_ids"`
	PrincipalID types.String `tfsdk:"principal_id"`
	TenantID    types.String `tfsdk:"tenant_id"`
}

var IdentityAttr = map[string]attr.Type{
	"type":         types.StringType,
	"identity_ids": types.ListType{}.WithElementType(types.StringType),
	"principal_id": types.StringType,
	"tenant_id":    types.StringType,
}
