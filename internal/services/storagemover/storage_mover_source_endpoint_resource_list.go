// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package storagemover

import (
	"context"
	"fmt"

	"github.com/hashicorp/go-azure-helpers/framework/typehelpers"
	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/storagemover/2025-07-01/endpoints"
	"github.com/hashicorp/go-azure-sdk/resource-manager/storagemover/2025-07-01/storagemovers"
	"github.com/hashicorp/terraform-plugin-framework/list"
	"github.com/hashicorp/terraform-plugin-framework/list/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type StorageMoverSourceEndpointListResource struct{}

type StorageMoverSourceEndpointListModel struct {
	StorageMoverId types.String `tfsdk:"storage_mover_id"`
}

var _ sdk.FrameworkListWrappedResource = new(StorageMoverSourceEndpointListResource)

func (StorageMoverSourceEndpointListResource) Metadata(_ context.Context, _ resource.MetadataRequest, response *resource.MetadataResponse) {
	response.TypeName = StorageMoverSourceEndpointResource{}.ResourceType()
}

func (StorageMoverSourceEndpointListResource) ResourceFunc() *pluginsdk.Resource {
	return sdk.WrappedResource(StorageMoverSourceEndpointResource{})
}

func (StorageMoverSourceEndpointListResource) ListResourceConfigSchema(_ context.Context, _ list.ListResourceSchemaRequest, response *list.ListResourceSchemaResponse) {
	response.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"storage_mover_id": schema.StringAttribute{
				Required: true,
				Validators: []validator.String{
					typehelpers.WrappedStringValidator{Func: storagemovers.ValidateStorageMoverID},
				},
			},
		},
	}
}

func (StorageMoverSourceEndpointListResource) List(ctx context.Context, request list.ListRequest, stream *list.ListResultsStream, metadata sdk.ResourceMetadata) {
	client := metadata.Client.StorageMover.EndpointsClient

	var data StorageMoverSourceEndpointListModel
	diags := request.Config.Get(ctx, &data)
	if diags.HasError() {
		stream.Results = list.ListResultsStreamDiagnostics(diags)
		return
	}

	storageMoverID, err := storagemovers.ParseStorageMoverID(data.StorageMoverId.ValueString())
	if err != nil {
		sdk.SetResponseErrorDiagnostic(stream, fmt.Sprintf("parsing Storage Mover ID for `%s`", StorageMoverSourceEndpointResource{}.ResourceType()), err)
		return
	}

	resp, err := client.ListComplete(ctx, endpoints.NewStorageMoverID(storageMoverID.SubscriptionId, storageMoverID.ResourceGroupName, storageMoverID.StorageMoverName))
	if err != nil {
		sdk.SetResponseErrorDiagnostic(stream, fmt.Sprintf("listing `%s`", StorageMoverSourceEndpointResource{}.ResourceType()), err)
		return
	}

	resource := StorageMoverSourceEndpointResource{}
	stream.Results = func(push func(list.ListResult) bool) {
		for _, item := range resp.Items {
			if _, ok := item.Properties.(endpoints.NfsMountEndpointProperties); !ok {
				continue
			}

			result := request.NewListResult(ctx)
			result.DisplayName = pointer.From(item.Name)

			id, err := endpoints.ParseEndpointIDInsensitively(pointer.From(item.Id))
			if err != nil {
				sdk.SetErrorDiagnosticAndPushListResult(result, push, "parsing Storage Mover Source Endpoint ID", err)
				return
			}

			meta := sdk.NewResourceMetaData(metadata.Client, resource)
			meta.SetID(id)

			if err := resource.flatten(meta, id, &item); err != nil {
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
}
