// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package cosmos

import (
	"context"
	"fmt"

	"github.com/hashicorp/go-azure-helpers/framework/typehelpers"
	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/cosmosdb/2025-10-15/fleets"
	"github.com/hashicorp/terraform-plugin-framework/list"
	"github.com/hashicorp/terraform-plugin-framework/list/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type CosmosDbFleetspaceListResource struct{}

var _ sdk.FrameworkListWrappedResource = new(CosmosDbFleetspaceListResource)

type CosmosDbFleetspaceListModel struct {
	FleetId types.String `tfsdk:"fleet_id"`
}

func (CosmosDbFleetspaceListResource) ListResourceConfigSchema(_ context.Context, _ list.ListResourceSchemaRequest, response *list.ListResourceSchemaResponse) {
	response.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"fleet_id": schema.StringAttribute{
				Required: true,
				Validators: []validator.String{
					typehelpers.WrappedStringValidator{
						Func: fleets.ValidateFleetID,
					},
				},
			},
		},
	}
}

func (CosmosDbFleetspaceListResource) ResourceFunc() *pluginsdk.Resource {
	return sdk.WrappedResource(CosmosDbFleetspaceResource{})
}

func (CosmosDbFleetspaceListResource) Metadata(_ context.Context, _ resource.MetadataRequest, response *resource.MetadataResponse) {
	response.TypeName = CosmosDbFleetspaceResource{}.ResourceType()
}

func (CosmosDbFleetspaceListResource) List(ctx context.Context, request list.ListRequest, stream *list.ListResultsStream, metadata sdk.ResourceMetadata) {
	client := metadata.Client.Cosmos.FleetsClient

	var data CosmosDbFleetspaceListModel
	diags := request.Config.Get(ctx, &data)
	if diags.HasError() {
		stream.Results = list.ListResultsStreamDiagnostics(diags)
		return
	}

	var results []fleets.FleetspaceResource
	r := CosmosDbFleetspaceResource{}

	fleetId, _ := fleets.ParseFleetID(data.FleetId.ValueString())
	resp, err := client.FleetspaceListComplete(ctx, pointer.From(fleetId))
	if err != nil {
		sdk.SetResponseErrorDiagnostic(stream, fmt.Sprintf("listing `%s`", r.ResourceType()), err)
		return
	}

	results = resp.Items

	stream.Results = func(push func(list.ListResult) bool) {
		for _, cosmosDbFleetspaceResult := range results {
			result := request.NewListResult(ctx)
			result.DisplayName = pointer.From(cosmosDbFleetspaceResult.Name)

			id, err := fleets.ParseFleetspaceID(pointer.From(cosmosDbFleetspaceResult.Id))
			if err != nil {
				sdk.SetErrorDiagnosticAndPushListResult(result, push, "parsing Fleetspace ID", err)
				return
			}

			rmd := sdk.NewResourceMetaData(metadata.Client, r)
			rmd.SetID(id)

			if err := r.flatten(rmd, id, &cosmosDbFleetspaceResult); err != nil {
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
