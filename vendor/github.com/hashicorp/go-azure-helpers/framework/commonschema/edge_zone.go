package commonschema

import (
	"github.com/hashicorp/go-azure-helpers/framework/planmodifiers/casing"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	datasourceschema "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	resourceschema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

func EdgeZoneDataSourceAttribute() datasourceschema.StringAttribute {
	return datasourceschema.StringAttribute{
		Computed: true,
	}
}

func EdgeZoneComputedAttribute() resourceschema.StringAttribute {
	return resourceschema.StringAttribute{
		Computed: true,
	}
}

func EdgeZoneOptionalAttribute() resourceschema.StringAttribute {
	return resourceschema.StringAttribute{
		Optional: true,
		Computed: true,
		Validators: []validator.String{
			stringvalidator.LengthAtLeast(1),
		},
		PlanModifiers: []planmodifier.String{
			casing.NormaliseLocationStringPlanModifier(),
		},
	}
}

func EdgeZoneOptionalRequiresReplaceAttribute() resourceschema.StringAttribute {
	return resourceschema.StringAttribute{
		Optional: true,
		Computed: true,
		Validators: []validator.String{
			stringvalidator.LengthAtLeast(1),
		},
		PlanModifiers: []planmodifier.String{
			stringplanmodifier.RequiresReplace(),
			casing.NormaliseLocationStringPlanModifier(),
		},
	}
}
