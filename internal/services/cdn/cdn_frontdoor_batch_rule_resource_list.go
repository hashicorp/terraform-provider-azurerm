// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package cdn

import (
	"context"
	"fmt"

	"github.com/hashicorp/go-azure-helpers/framework/typehelpers"
	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	legacyrulesets "github.com/hashicorp/go-azure-sdk/resource-manager/cdn/2024-02-01/rulesets"
	"github.com/hashicorp/go-azure-sdk/resource-manager/cdn/2025-12-01/rules"
	"github.com/hashicorp/terraform-plugin-framework/list"
	listschema "github.com/hashicorp/terraform-plugin-framework/list/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type CdnFrontDoorBatchRuleListResource struct{}

type CdnFrontDoorBatchRuleListModel struct {
	CdnFrontDoorRuleSetID types.String `tfsdk:"cdn_frontdoor_rule_set_id"`
}

var _ sdk.FrameworkListWrappedResource = new(CdnFrontDoorBatchRuleListResource)

func (CdnFrontDoorBatchRuleListResource) Metadata(_ context.Context, _ resource.MetadataRequest, response *resource.MetadataResponse) {
	response.TypeName = CdnFrontDoorBatchRuleResource{}.ResourceType()
}

func (CdnFrontDoorBatchRuleListResource) ResourceFunc() *pluginsdk.Resource {
	return sdk.WrappedResource(CdnFrontDoorBatchRuleResource{})
}

func (CdnFrontDoorBatchRuleListResource) ListResourceConfigSchema(_ context.Context, _ list.ListResourceSchemaRequest, response *list.ListResourceSchemaResponse) {
	response.Schema = listschema.Schema{
		Attributes: map[string]listschema.Attribute{
			"cdn_frontdoor_rule_set_id": listschema.StringAttribute{
				Required: true,
				Validators: []validator.String{
					typehelpers.WrappedStringValidator{
						Func: rules.ValidateRuleSetID,
					},
				},
			},
		},
	}
}

func (CdnFrontDoorBatchRuleListResource) List(ctx context.Context, request list.ListRequest, stream *list.ListResultsStream, metadata sdk.ResourceMetadata) {
	client := metadata.Client.Cdn.FrontDoorRuleSetsClient_v2025_12_01

	var data CdnFrontDoorBatchRuleListModel
	diags := request.Config.Get(ctx, &data)
	if diags.HasError() {
		stream.Results = list.ListResultsStreamDiagnostics(diags)
		return
	}

	ruleSetID, err := rules.ParseRuleSetID(data.CdnFrontDoorRuleSetID.ValueString())
	if err != nil {
		sdk.SetResponseErrorDiagnostic(stream, "parsing `cdn_frontdoor_rule_set_id`", err)
		return
	}

	resource := CdnFrontDoorBatchRuleResource{}

	legacyRuleSetID := legacyrulesets.NewRuleSetID(ruleSetID.SubscriptionId, ruleSetID.ResourceGroupName, ruleSetID.ProfileName, ruleSetID.RuleSetName)
	resp, err := client.Get(ctx, legacyRuleSetID)
	if err != nil {
		sdk.SetResponseErrorDiagnostic(stream, fmt.Sprintf("retrieving `%s`", resource.ResourceType()), err)
		return
	}

	stream.Results = func(push func(list.ListResult) bool) {
		if resp.Model == nil || resp.Model.Properties == nil || resp.Model.Properties.BatchMode == nil || !*resp.Model.Properties.BatchMode || resp.Model.Properties.Rules == nil || len(pointer.From(resp.Model.Properties.Rules)) == 0 {
			return
		}

		result := request.NewListResult(ctx)
		result.DisplayName = ruleSetID.RuleSetName

		meta := sdk.NewResourceMetaData(metadata.Client, resource)
		meta.ResourceData.SetId(ruleSetID.ID())

		if err := resource.flatten(meta, ruleSetID, resp.Model); err != nil {
			sdk.SetErrorDiagnosticAndPushListResult(result, push, fmt.Sprintf("encoding `%s` resource data", resource.ResourceType()), err)
			return
		}

		sdk.EncodeListResult(ctx, meta.ResourceData, &result)
		if result.Diagnostics.HasError() {
			push(result)
			return
		}

		if !push(result) {
			return
		}
	}
}
