// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package storage_sync_server_endpoint

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
	ListResource                struct{}
	SyncServerEndpointListModel struct {
		StorageSyncGroupId types.String `tfsdk:"storage_sync_group_id"`
	}
)

var _ sdk.FrameworkListWrappedResource = new(ListResource)

func (ListResource) Metadata(_ context.Context, _ resource.MetadataRequest, response *resource.MetadataResponse) {
	response.TypeName = Resource{}.ResourceType()
}

func (ListResource) ResourceFunc() *pluginsdk.Resource {
	return sdk.WrappedResource(Resource{})
}

func (ListResource) ListResourceConfigSchema(_ context.Context, _ list.ListResourceSchemaRequest, response *list.ListResourceSchemaResponse) {
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

func (ListResource) List(ctx context.Context, request list.ListRequest, stream *list.ListResultsStream, metadata sdk.ResourceMetadata) {
	client := metadata.Client.Storage.SyncServerEndpointsClient

	var data SyncServerEndpointListModel
	diags := request.Config.Get(ctx, &data)
	if diags.HasError() {
		stream.Results = list.ListResultsStreamDiagnostics(diags)
		return
	}

	syncGroupId, err := serverendpointresource.ParseSyncGroupID(data.StorageSyncGroupId.ValueString())
	if err != nil {
		sdk.SetResponseErrorDiagnostic(stream, fmt.Sprintf("parsing storage sync group id for `%s`", Resource{}.ResourceType()), err)
		return
	}

	resp, err := client.ServerEndpointsListBySyncGroup(ctx, *syncGroupId)
	if err != nil {
		sdk.SetResponseErrorDiagnostic(stream, fmt.Sprintf("listing `%s`", Resource{}.ResourceType()), err)
		return
	}

	results := make([]serverendpointresource.ServerEndpoint, 0)
	if model := resp.Model; model != nil {
		results = pointer.From(model.Value)
	}

	r := Resource{}

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
