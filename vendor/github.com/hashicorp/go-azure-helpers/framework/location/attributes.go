// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package location

import (
	"github.com/hashicorp/go-azure-helpers/framework/typehelpers"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	datasourceschema "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

func LocationAttribute() schema.StringAttribute {
	return schema.StringAttribute{
		Required: true,
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

func LocationComputedAttribute() schema.StringAttribute {
	return schema.StringAttribute{
		Computed: true,
	}
}

func LocationDataSourceAttribute() datasourceschema.StringAttribute {
	return datasourceschema.StringAttribute{
		Required: true,
		Validators: []validator.String{
			typehelpers.WrappedStringValidator{
				Func: location.EnhancedValidate,
			},
		},
	}
}

func LocationRelocatableAttribute() schema.StringAttribute {
	return schema.StringAttribute{
		Required: true,
		Validators: []validator.String{
			typehelpers.WrappedStringValidator{
				Func: location.EnhancedValidate,
			},
		},
	}
}
