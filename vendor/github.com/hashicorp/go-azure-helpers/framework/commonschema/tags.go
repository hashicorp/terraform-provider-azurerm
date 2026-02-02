// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package commonschema

import (
	"context"

	"github.com/hashicorp/go-azure-helpers/framework/typehelpers"
	"github.com/hashicorp/terraform-plugin-framework-validators/mapvalidator"
	datasourceschema "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	resourceschema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

func TagsResourceAttribute(ctx context.Context) resourceschema.MapAttribute {
	return resourceschema.MapAttribute{
		CustomType:          typehelpers.NewMapTypeOf[types.String](ctx),
		ElementType:         basetypes.StringType{},
		Optional:            true,
		Description:         "A map of tags to be assigned to the resource",
		MarkdownDescription: "A map of tags to be assigned to the resource",
		Validators: []validator.Map{
			mapvalidator.SizeAtLeast(1),
		},
	}
}

func TagsDataSourceAttribute(ctx context.Context) datasourceschema.MapAttribute {
	return datasourceschema.MapAttribute{
		CustomType:          typehelpers.NewMapTypeOf[types.String](ctx),
		ElementType:         basetypes.StringType{},
		Optional:            true,
		Description:         "A map of tags assigned to the resource",
		MarkdownDescription: "A map of tags assigned to the resource",
		Validators: []validator.Map{
			mapvalidator.SizeAtLeast(1),
		},
	}
}

func ExpandTags(input types.Map) (result *map[string]string, diags diag.Diagnostics) {
	if input.IsNull() || input.IsUnknown() {
		return
	}

	diags = input.ElementsAs(context.Background(), &result, false)

	return
}

func FlattenTags(tags *map[string]string) (result basetypes.MapValue, diags diag.Diagnostics) {
	if tags == nil {
		return basetypes.NewMapNull(basetypes.StringType{}), nil
	}

	return types.MapValueFrom(context.Background(), basetypes.StringType{}, tags)
}
