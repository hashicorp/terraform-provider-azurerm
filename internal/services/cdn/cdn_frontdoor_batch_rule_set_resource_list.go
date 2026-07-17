// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package cdn

import (
	"context"
	"fmt"

	"github.com/hashicorp/go-azure-helpers/framework/typehelpers"
	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/cdn/2025-12-01/profiles"
	"github.com/hashicorp/go-azure-sdk/resource-manager/cdn/2025-12-01/rulesets"
	"github.com/hashicorp/terraform-plugin-framework/list"
	listschema "github.com/hashicorp/terraform-plugin-framework/list/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type CdnFrontDoorBatchRuleSetListResource struct{}

type CdnFrontDoorBatchRuleSetListModel struct {
	CdnFrontDoorProfileID types.String `tfsdk:"cdn_frontdoor_profile_id"`
}

var _ sdk.FrameworkListWrappedResource = new(CdnFrontDoorBatchRuleSetListResource)

func (CdnFrontDoorBatchRuleSetListResource) Metadata(_ context.Context, _ resource.MetadataRequest, response *resource.MetadataResponse) {
	response.TypeName = CdnFrontDoorBatchRuleSetResource{}.ResourceType()
}

func (CdnFrontDoorBatchRuleSetListResource) ResourceFunc() *pluginsdk.Resource {
	return sdk.WrappedResource(CdnFrontDoorBatchRuleSetResource{})
}

func (CdnFrontDoorBatchRuleSetListResource) ListResourceConfigSchema(_ context.Context, _ list.ListResourceSchemaRequest, response *list.ListResourceSchemaResponse) {
	response.Schema = listschema.Schema{
		Attributes: map[string]listschema.Attribute{
			"cdn_frontdoor_profile_id": listschema.StringAttribute{
				Required: true,
				Validators: []validator.String{
					typehelpers.WrappedStringValidator{
						Func: rulesets.ValidateProfileID,
					},
				},
			},
		},
	}
}

func (CdnFrontDoorBatchRuleSetListResource) List(ctx context.Context, request list.ListRequest, stream *list.ListResultsStream, metadata sdk.ResourceMetadata) {
	client := metadata.Client.Cdn.FrontDoorRuleSetsClient

	var data CdnFrontDoorBatchRuleSetListModel
	diags := request.Config.Get(ctx, &data)
	if diags.HasError() {
		stream.Results = list.ListResultsStreamDiagnostics(diags)
		return
	}

	profileID, err := profiles.ParseProfileID(data.CdnFrontDoorProfileID.ValueString())
	if err != nil {
		sdk.SetResponseErrorDiagnostic(stream, "parsing `cdn_frontdoor_profile_id`", err)
		return
	}

	r := CdnFrontDoorBatchRuleSetResource{}

	resp, err := client.ListByProfileComplete(ctx, rulesets.NewProfileID(profileID.SubscriptionId, profileID.ResourceGroupName, profileID.ProfileName))
	if err != nil {
		sdk.SetResponseErrorDiagnostic(stream, fmt.Sprintf("retrieving `%s`", r.ResourceType()), err)
		return
	}

	deadline, ok := ctx.Deadline()
	if !ok {
		sdk.SetResponseErrorDiagnostic(stream, "internal-error", "context had no deadline")
		return
	}

	stream.Results = func(push func(list.ListResult) bool) {
		ctx, cancel := context.WithDeadline(context.Background(), deadline)
		defer cancel()

		for _, item := range resp.Items {
			result := request.NewListResult(ctx)
			result.DisplayName = pointer.From(item.Name)

			ruleSetID, err := rulesets.ParseRuleSetIDInsensitively(pointer.From(item.Id))
			if err != nil {
				sdk.SetErrorDiagnosticAndPushListResult(result, push, "parsing Rule Set ID", err)
				return
			}

			detailedResp, err := client.Get(ctx, *ruleSetID)
			if err != nil {
				sdk.SetErrorDiagnosticAndPushListResult(result, push, fmt.Sprintf("retrieving `%s`", r.ResourceType()), err)
				return
			}
			if detailedResp.Model == nil || detailedResp.Model.Properties == nil || detailedResp.Model.Properties.BatchMode == nil || !*detailedResp.Model.Properties.BatchMode {
				continue
			}

			rmd := sdk.NewResourceMetaData(metadata.Client, r)
			rmd.ResourceData.SetId(ruleSetID.ID())

			if err := r.flatten(rmd, ruleSetID, detailedResp.Model); err != nil {
				sdk.SetErrorDiagnosticAndPushListResult(result, push, fmt.Sprintf("encoding `%s` resource data", r.ResourceType()), err)
				return
			}

			sdk.EncodeListResult(ctx, rmd.ResourceData, &result)
			if result.Diagnostics.HasError() {
				push(result)
				return
			}

			if !push(result) {
				return
			}
		}
	}
}
