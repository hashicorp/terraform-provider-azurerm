package commonschema

import (
	"github.com/hashicorp/go-azure-helpers/framework/typehelpers"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourcegroups"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

func ResourceGroupNameAttribute() schema.StringAttribute {
	return schema.StringAttribute{
		Required:            true,
		Description:         "The name of the resource group",
		MarkdownDescription: "The name of the resource group",
		Validators: []validator.String{
			typehelpers.WrappedStringValidator{
				Func: resourcegroups.ValidateName,
			},
		},
		PlanModifiers: []planmodifier.String{
			stringplanmodifier.RequiresReplace(),
		},
	}
}
