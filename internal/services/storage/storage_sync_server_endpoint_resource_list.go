// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package storage

import (
	"context"
	"fmt"

	"github.com/hashicorp/go-azure-helpers/framework/typehelpers"
	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/storagesync/2020-03-01/serverendpointresource"
	"github.com/hashicorp/terraform-plugin-framework/list"
	"github.com/hashicorp/terraform-plugin-framework/list/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type (
	SyncServerEndpointListResource struct{}
	SyncServerEndpointListModel    struct {
		StorageSyncGroupId types.String `tfsdk:"storage_sync_group_id"`
	}
)

var _ sdk.FrameworkListWrappedResource = new(SyncServerEndpointListResource)

func (SyncServerEndpointListResource) Metadata(_ context.Context, _ resource.MetadataRequest, response *resource.MetadataResponse) {
	response.TypeName = SyncServerEndpointResource{}.ResourceType()
}

func (SyncServerEndpointListResource) ResourceFunc() *pluginsdk.Resource {
	return sdk.WrappedResource(SyncServerEndpointResource{})
}

func (SyncServerEndpointListResource) ListResourceConfigSchema(_ context.Context, _ list.ListResourceSchemaRequest, response *list.ListResourceSchemaResponse) {
	response.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"storage_sync_group_id": schema.StringAttribute{
				Required: true,
				Validators: []validator.String{
					typehelpers.WrappedStringValidator{Func: serverendpointresource.ValidateSyncGroupID},
				},
			},
		},
	}
}

func (SyncServerEndpointListResource) List(ctx context.Context, request list.ListRequest, stream *list.ListResultsStream, metadata sdk.ResourceMetadata) {
	client := metadata.Client.Storage.SyncServerEndpointsClient

	var data SyncServerEndpointListModel
	diags := request.Config.Get(ctx, &data)
	if diags.HasError() {
		stream.Results = list.ListResultsStreamDiagnostics(diags)
		return
	}

	syncGroupId, err := serverendpointresource.ParseSyncGroupID(data.StorageSyncGroupId.ValueString())
	if err != nil {
		sdk.SetResponseErrorDiagnostic(stream, fmt.Sprintf("parsing storage sync group id for `%s`", SyncServerEndpointResource{}.ResourceType()), err)
		return
	}

	resp, err := client.ServerEndpointsListBySyncGroup(ctx, *syncGroupId)
	if err != nil {
		sdk.SetResponseErrorDiagnostic(stream, fmt.Sprintf("listing `%s`", SyncServerEndpointResource{}.ResourceType()), err)
		return
	}

	results := make([]serverendpointresource.ServerEndpoint, 0)
	if model := resp.Model; model != nil {
		results = pointer.From(model.Value)
	}

	r := SyncServerEndpointResource{}

	stream.Results = func(push func(list.ListResult) bool) {
		for _, serverEndpoint := range results {
			result := request.NewListResult(ctx)
			result.DisplayName = pointer.From(serverEndpoint.Name)

			id, err := serverendpointresource.ParseServerEndpointIDInsensitively(pointer.From(serverEndpoint.Id))
			if err != nil {
				sdk.SetErrorDiagnosticAndPushListResult(result, push, "parsing Storage Sync Server Endpoint ID", err)
				return
			}

			rmd := sdk.NewResourceMetaData(metadata.Client, r)
			rmd.SetID(id)

			if err := r.flatten(rmd, id, &serverEndpoint); err != nil {
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
