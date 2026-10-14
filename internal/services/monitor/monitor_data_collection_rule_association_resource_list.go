// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package monitor

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/framework/typehelpers"
	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/insights/2023-03-11/datacollectionruleassociations"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/list"
	"github.com/hashicorp/terraform-plugin-framework/list/schema"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type (
	DataCollectionRuleAssociationListResource struct{}
	DataCollectionRuleAssociationListModel    struct {
		DataCollectionEndpointId types.String `tfsdk:"data_collection_endpoint_id"`
		DataCollectionRuleId     types.String `tfsdk:"data_collection_rule_id"`
		TargetResourceId         types.String `tfsdk:"target_resource_id"`
	}
)

var (
	_ sdk.FrameworkListWrappedResource           = new(DataCollectionRuleAssociationListResource)
	_ sdk.FrameworkListWrappedResourceWithConfig = new(DataCollectionRuleAssociationListResource)
)

func (DataCollectionRuleAssociationListResource) Metadata(_ context.Context, _ resource.MetadataRequest, response *resource.MetadataResponse) {
	response.TypeName = DataCollectionRuleAssociationResource{}.ResourceType()
}

func (DataCollectionRuleAssociationListResource) ResourceFunc() *pluginsdk.Resource {
	return sdk.WrappedResource(DataCollectionRuleAssociationResource{})
}

func (DataCollectionRuleAssociationListResource) ListResourceConfigSchema(_ context.Context, _ list.ListResourceSchemaRequest, response *list.ListResourceSchemaResponse) {
	response.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"data_collection_endpoint_id": schema.StringAttribute{
				Optional: true,
				Validators: []validator.String{
					typehelpers.WrappedStringValidator{
						Func: datacollectionruleassociations.ValidateDataCollectionEndpointID,
					},
					stringvalidator.ExactlyOneOf(
						path.MatchRoot("data_collection_endpoint_id"),
						path.MatchRoot("data_collection_rule_id"),
						path.MatchRoot("target_resource_id"),
					),
				},
			},
			"data_collection_rule_id": schema.StringAttribute{
				Optional: true,
				Validators: []validator.String{
					typehelpers.WrappedStringValidator{
						Func: datacollectionruleassociations.ValidateDataCollectionRuleID,
					},
					stringvalidator.ExactlyOneOf(
						path.MatchRoot("data_collection_endpoint_id"),
						path.MatchRoot("data_collection_rule_id"),
						path.MatchRoot("target_resource_id"),
					),
				},
			},
			"target_resource_id": schema.StringAttribute{
				Optional: true,
				Validators: []validator.String{
					typehelpers.WrappedStringValidator{
						Func: commonids.ValidateScopeID,
					},
					stringvalidator.ExactlyOneOf(
						path.MatchRoot("data_collection_endpoint_id"),
						path.MatchRoot("data_collection_rule_id"),
						path.MatchRoot("target_resource_id"),
					),
				},
			},
		},
	}
}

func (DataCollectionRuleAssociationListResource) List(ctx context.Context, request list.ListRequest, stream *list.ListResultsStream, metadata sdk.ResourceMetadata) {
	client := metadata.Client.Monitor.DataCollectionRuleAssociationsClient

	var data DataCollectionRuleAssociationListModel
	diags := request.Config.Get(ctx, &data)
	if diags.HasError() {
		stream.Results = list.ListResultsStreamDiagnostics(diags)
		return
	}

	results := make([]datacollectionruleassociations.DataCollectionRuleAssociationProxyOnlyResource, 0)

	switch {
	case !data.TargetResourceId.IsNull():
		targetResourceId := data.TargetResourceId.ValueString()
		scopeId := commonids.NewScopeID(targetResourceId)

		resp, err := client.ListByResourceComplete(ctx, scopeId)
		if err != nil {
			sdk.SetResponseErrorDiagnostic(stream, fmt.Sprintf("listing data collection rule associations for target resource `%s`", targetResourceId), err)
			return
		}
		results = resp.Items

	case !data.DataCollectionRuleId.IsNull():
		dcrIdStr := data.DataCollectionRuleId.ValueString()
		ruleId, err := datacollectionruleassociations.ParseDataCollectionRuleID(dcrIdStr)
		if err != nil {
			sdk.SetResponseErrorDiagnostic(stream, fmt.Sprintf("parsing Data Collection Rule ID `%s`", dcrIdStr), err)
			return
		}
		resp, err := client.ListByRuleComplete(ctx, *ruleId)
		if err != nil {
			sdk.SetResponseErrorDiagnostic(stream, fmt.Sprintf("listing data collection rule associations for data collection rule `%s`", dcrIdStr), err)
			return
		}
		results = resp.Items

	case !data.DataCollectionEndpointId.IsNull():
		dceIdStr := data.DataCollectionEndpointId.ValueString()
		endpointId, err := datacollectionruleassociations.ParseDataCollectionEndpointID(dceIdStr)
		if err != nil {
			sdk.SetResponseErrorDiagnostic(stream, fmt.Sprintf("parsing Data Collection Endpoint ID `%s`", dceIdStr), err)
			return
		}
		resp, err := client.ListByDataCollectionEndpointComplete(ctx, *endpointId)
		if err != nil {
			sdk.SetResponseErrorDiagnostic(stream, fmt.Sprintf("listing data collection rule associations for data collection endpoint `%s`", dceIdStr), err)
			return
		}
		results = resp.Items
	}

	r := DataCollectionRuleAssociationResource{}

	stream.Results = func(push func(list.ListResult) bool) {
		for _, item := range results {
			result := request.NewListResult(ctx)
			result.DisplayName = pointer.From(item.Name)

			rawId := pointer.From(item.Id)
			if !strings.HasPrefix(rawId, "/") {
				rawId = "/" + rawId
			}

			id, err := datacollectionruleassociations.ParseScopedDataCollectionRuleAssociationIDInsensitively(rawId)
			if err != nil {
				sdk.SetErrorDiagnosticAndPushListResult(result, push, fmt.Sprintf("parsing Data Collection Rule Association ID for `%s`", r.ResourceType()), err)
				return
			}

			rmd := sdk.NewResourceMetaData(metadata.Client, r)
			rmd.SetID(id)

			if err := r.flatten(rmd, id, &item); err != nil {
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
