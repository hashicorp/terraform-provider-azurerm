package commonschema

import (
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	datasourceschema "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	resourceschema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

func ZoneSingleOptional() resourceschema.StringAttribute {
	return resourceschema.StringAttribute{
		Optional: true,
		Validators: []validator.String{
			stringvalidator.LengthAtLeast(1),
		},
	}
}

func ZoneSingleOptionalForceNew() resourceschema.StringAttribute {
	return resourceschema.StringAttribute{
		Optional: true,
		PlanModifiers: []planmodifier.String{
			stringplanmodifier.RequiresReplace(),
		},
		Validators: []validator.String{
			stringvalidator.LengthAtLeast(1),
		},
	}
}

func ZoneSingleRequired() resourceschema.StringAttribute {
	return resourceschema.StringAttribute{
		Required: true,
		Validators: []validator.String{
			stringvalidator.LengthAtLeast(1),
		},
	}
}

func ZoneSingleRequiredForceNew() resourceschema.StringAttribute {
	return resourceschema.StringAttribute{
		Required: true,
		PlanModifiers: []planmodifier.String{
			stringplanmodifier.RequiresReplace(),
		},
		Validators: []validator.String{
			stringvalidator.LengthAtLeast(1),
		},
	}
}

func ZoneSingleOptionalComputed() resourceschema.StringAttribute {
	return resourceschema.StringAttribute{
		Optional: true,
		Computed: true,
		Validators: []validator.String{
			stringvalidator.LengthAtLeast(1),
		},
	}
}
func ZoneSingleComputed() resourceschema.StringAttribute {
	return resourceschema.StringAttribute{
		Computed: true,
	}
}
func ZoneSingleDataSource() datasourceschema.StringAttribute {
	return datasourceschema.StringAttribute{
		Computed: true,
	}
}
