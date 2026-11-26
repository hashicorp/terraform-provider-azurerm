package commonschema

import (
	"github.com/hashicorp/go-azure-helpers/framework/typehelpers"
	"github.com/hashicorp/terraform-plugin-framework-validators/setvalidator"
	datasourceschema "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	resourceschema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/setplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

func ZonesMultipleOptional() resourceschema.SetAttribute {
	return resourceschema.SetAttribute{
		ElementType: typehelpers.SetOfStringType,
		Optional:    true,
		Validators: []validator.Set{
			setvalidator.SizeAtLeast(1),
		},
	}
}

func ZonesMultipleOptionalForceNew() resourceschema.SetAttribute {
	return resourceschema.SetAttribute{
		ElementType: typehelpers.SetOfStringType,
		Optional:    true,
		Validators: []validator.Set{
			setvalidator.SizeAtLeast(1),
		},
		PlanModifiers: []planmodifier.Set{
			setplanmodifier.RequiresReplace(),
		},
	}
}

func ZonesMultipleRequired() resourceschema.SetAttribute {
	return resourceschema.SetAttribute{
		ElementType: typehelpers.SetOfStringType,
		Required:    true,
		Validators: []validator.Set{
			setvalidator.SizeAtLeast(1),
		},
	}
}

func ZonesMultipleRequiredForceNew() resourceschema.SetAttribute {
	return resourceschema.SetAttribute{
		ElementType: typehelpers.SetOfStringType,
		Required:    true,
		Validators: []validator.Set{
			setvalidator.SizeAtLeast(1),
		},
		PlanModifiers: []planmodifier.Set{
			setplanmodifier.RequiresReplace(),
		},
	}
}

func ZonesMultipleComputed() resourceschema.SetAttribute {
	return resourceschema.SetAttribute{
		ElementType: typehelpers.SetOfStringType,
		Computed:    true,
	}
}

func ZonesMultipleDataSource() datasourceschema.SetAttribute {
	return datasourceschema.SetAttribute{
		ElementType: typehelpers.SetOfStringType,
		Computed:    true,
	}
}
