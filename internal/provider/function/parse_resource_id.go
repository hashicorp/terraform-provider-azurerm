// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package function

import (
	"context"
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/recaser"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/function"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type ParseResourceIDFunction struct{}

var _ function.Function = ParseResourceIDFunction{}

var idParseResultTypes = map[string]attr.Type{
	"resource_name":       types.StringType,
	"resource_provider":   types.StringType,
	"resource_group_name": types.StringType,
	"resource_type":       types.StringType,
	"resource_scope":      types.StringType,
	"full_resource_type":  types.StringType,
	"subscription_id":     types.StringType,
	"parent_resources":    types.MapType{}.WithElementType(types.StringType),
}

func NewParseResourceIDFunction() function.Function {
	return &ParseResourceIDFunction{}
}

func (p ParseResourceIDFunction) Metadata(_ context.Context, _ function.MetadataRequest, response *function.MetadataResponse) {
	response.Name = "parse_resource_id"
}

func (p ParseResourceIDFunction) Definition(_ context.Context, _ function.DefinitionRequest, response *function.DefinitionResponse) {
	response.Definition = function.Definition{
		Summary:             "parse_resource_id",
		Description:         "Parses an Azure Resource Manager ID and exposes the contained information",
		MarkdownDescription: "Parses an Azure Resource Manager ID and exposes the contained information",
		Parameters: []function.Parameter{
			function.StringParameter{
				Name:                "id",
				Description:         "Resource ID",
				MarkdownDescription: "Resource ID",
			},
		},
		Return: function.ObjectReturn{
			AttributeTypes: idParseResultTypes,
		},
	}
}

func (p ParseResourceIDFunction) Run(ctx context.Context, request function.RunRequest, response *function.RunResponse) {
	var id string

	response.Error = function.ConcatFuncErrors(request.Arguments.Get(ctx, &id))

	if response.Error != nil {
		return
	}

	if len(id) == 0 {
		response.Error = function.NewFuncError("Got empty ID")
		return
	}

	// These outputs should always have a value
	output := map[string]attr.Value{
		"resource_name":       types.StringValue(""),
		"resource_provider":   types.StringValue(""),
		"resource_group_name": types.StringValue(""),
		"resource_type":       types.StringValue(""),
		"resource_scope":      types.StringValue(""),
		"full_resource_type":  types.StringValue(""),
		"subscription_id":     types.StringValue(""),
		"parent_resources":    types.Map{},
	}

	idType := recaser.ResourceIdTypeFromResourceId(id)
	if idType == nil {
		response.Error = function.NewFuncError(fmt.Sprintf("could not determine resource ID type from %s, ID may be malformed or currently not supported in the provider", id))
		return
	}

	parser := resourceids.NewParserFromResourceIdType(idType)
	parsed, err := parser.Parse(id, true)
	if err != nil {
		response.Error = function.NewFuncError(fmt.Sprintf("Parsing Resource ID Error: %s", err))
		return
	}

	err = idType.FromParseResult(*parsed)
	if err != nil {
		response.Error = function.NewFuncError(fmt.Sprintf("Expanding Parsed Resource ID Error: %s", err))
		return
	}

	s := idType.Segments()
	numSegments := len(s)
	pTemp := ""
	fullResourceType := ""
	parentMap := map[string]string{}
	for k, v := range s {
		switch v.Type {
		case resourceids.ResourceGroupSegmentType:
			output["resource_group_name"] = types.StringValue(parsed.Parsed[v.Name])

		case resourceids.ResourceProviderSegmentType:
			output["resource_provider"] = types.StringPointerValue(v.FixedValue)
			fullResourceType = pointer.From(v.FixedValue)

		case resourceids.SubscriptionIdSegmentType:
			output["subscription_id"] = types.StringValue(parsed.Parsed["subscriptionId"])

		case resourceids.StaticSegmentType:
			switch {
			case k == (numSegments - 2):
				{
					output["resource_type"] = types.StringPointerValue(v.FixedValue)
					fullResourceType = fmt.Sprintf("%s/%s", fullResourceType, pointer.From(v.FixedValue))
				}
			case v.FixedValue != nil && *v.FixedValue != "subscriptions" && *v.FixedValue != "resourceGroups" && *v.FixedValue != "providers":
				{
					pTemp = parsed.Parsed[v.Name]
					fullResourceType = fmt.Sprintf("%s/%s", fullResourceType, pointer.From(v.FixedValue))
				}
			}

		case resourceids.UserSpecifiedSegmentType:
			if k == (numSegments - 1) {
				output["resource_name"] = types.StringValue(parsed.Parsed[v.Name])
			} else {
				parentMap[pTemp] = parsed.Parsed[v.Name]
			}
		case resourceids.ScopeSegmentType:
			output["resource_scope"] = types.StringValue(parsed.Parsed[v.Name])
		}
	}
	parentMapValue, diags := types.MapValueFrom(ctx, types.StringType, parentMap)
	if diags.HasError() {
		response.Error = function.NewFuncError("failed to flatten parent resources")
		return
	}
	output["parent_resources"] = parentMapValue
	output["full_resource_type"] = types.StringValue(fullResourceType)

	result, diags := types.ObjectValue(idParseResultTypes, output)
	if diags.HasError() {
		response.Error = function.ConcatFuncErrors(response.Error, function.FuncErrorFromDiags(ctx, diags))
		return
	}

	response.Error = function.ConcatFuncErrors(response.Result.Set(ctx, result))
}
