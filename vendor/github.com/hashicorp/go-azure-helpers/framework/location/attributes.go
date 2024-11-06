package location

import (
	"github.com/hashicorp/go-azure-helpers/framework/typehelpers"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

func LocationAttribute() schema.StringAttribute {
	return schema.StringAttribute{
		Required:            true,
		Description:         "",
		MarkdownDescription: "",
		Validators: []validator.String{
			typehelpers.WrappedStringValidator{
				Func: location.EnhancedValidate,
			},
		},
		PlanModifiers: []planmodifier.String{
			stringplanmodifier.RequiresReplace(),
		},
	}
}
