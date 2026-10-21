// Copyright IBM Corp. 2018, 2025
// SPDX-License-Identifier: MPL-2.0

package commonschema

import (
	"context"

	"github.com/hashicorp/go-azure-helpers/framework/convert"
	"github.com/hashicorp/go-azure-helpers/framework/typehelpers"
	"github.com/hashicorp/terraform-plugin-framework-validators/mapvalidator"
	datasourceschema "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	resourceschema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func TagsResourceAttribute(ctx context.Context) resourceschema.MapAttribute {
	return resourceschema.MapAttribute{
		CustomType:          typehelpers.NewMapTypeOf[types.String](ctx),
		ElementType:         types.StringType,
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
		ElementType:         types.StringType,
		Optional:            true,
		Description:         "A map of tags assigned to the resource",
		MarkdownDescription: "A map of tags assigned to the resource",
	}
}

func ExpandTags(ctx context.Context, input typehelpers.MapValueOf[types.String], diags *diag.Diagnostics) *map[string]string {
	if input.IsNull() || input.IsUnknown() {
		return nil
	}

	result := make(map[string]string)

	convert.Expand(ctx, input, &result, diags)

	return &result
}

func FlattenTags(ctx context.Context, tags *map[string]string, diags *diag.Diagnostics) typehelpers.MapValueOf[types.String] {
	result := typehelpers.NewMapValueOfNull[types.String](ctx)
	if tags == nil {
		return result
	}

	convert.Flatten(ctx, tags, &result, diags)

	return result
}
