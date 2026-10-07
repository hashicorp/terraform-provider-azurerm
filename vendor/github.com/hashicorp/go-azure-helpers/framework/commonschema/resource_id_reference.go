package commonschema

import (
	"github.com/hashicorp/go-azure-helpers/framework/validators"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
	resourceschema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

func ResourceIDReferenceRequired(id resourceids.ResourceId) resourceschema.StringAttribute {
	return resourceschema.StringAttribute{
		Required: true,
		Validators: []validator.String{
			validators.AzureResourceManagerId(id),
		},
	}
}

func ResourceIDReferenceRequiredForceNew(id resourceids.ResourceId) resourceschema.StringAttribute {
	return resourceschema.StringAttribute{
		Required: true,
		Validators: []validator.String{
			validators.AzureResourceManagerId(id),
		},
		PlanModifiers: []planmodifier.String{
			stringplanmodifier.RequiresReplace(),
		},
	}
}
