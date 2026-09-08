// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package storagemover

import (
	"context"
	"fmt"

	"github.com/hashicorp/go-azure-helpers/framework/typehelpers"
	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/storagemover/2025-07-01/projects"
	"github.com/hashicorp/go-azure-sdk/resource-manager/storagemover/2025-07-01/storagemovers"
	"github.com/hashicorp/terraform-plugin-framework/list"
	"github.com/hashicorp/terraform-plugin-framework/list/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type StorageMoverProjectListResource struct{}

type StorageMoverProjectListModel struct {
	StorageMoverId types.String `tfsdk:"storage_mover_id"`
}

var _ sdk.FrameworkListWrappedResource = new(StorageMoverProjectListResource)

func (StorageMoverProjectListResource) Metadata(_ context.Context, _ resource.MetadataRequest, response *resource.MetadataResponse) {
	response.TypeName = StorageMoverProjectResource{}.ResourceType()
}

func (StorageMoverProjectListResource) ResourceFunc() *pluginsdk.Resource {
	return sdk.WrappedResource(StorageMoverProjectResource{})
}

func (StorageMoverProjectListResource) ListResourceConfigSchema(_ context.Context, _ list.ListResourceSchemaRequest, response *list.ListResourceSchemaResponse) {
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

func (StorageMoverProjectListResource) List(ctx context.Context, request list.ListRequest, stream *list.ListResultsStream, metadata sdk.ResourceMetadata) {
	client := metadata.Client.StorageMover.ProjectsClient

	var data StorageMoverProjectListModel
	diags := request.Config.Get(ctx, &data)
	if diags.HasError() {
		stream.Results = list.ListResultsStreamDiagnostics(diags)
		return
	}

	storageMoverID, err := storagemovers.ParseStorageMoverID(data.StorageMoverId.ValueString())
	if err != nil {
		sdk.SetResponseErrorDiagnostic(stream, fmt.Sprintf("parsing Storage Mover ID for `%s`", StorageMoverProjectResource{}.ResourceType()), err)
		return
	}

	resp, err := client.ListComplete(ctx, projects.NewStorageMoverID(storageMoverID.SubscriptionId, storageMoverID.ResourceGroupName, storageMoverID.StorageMoverName))
	if err != nil {
		sdk.SetResponseErrorDiagnostic(stream, fmt.Sprintf("listing `%s`", StorageMoverProjectResource{}.ResourceType()), err)
		return
	}

	resource := StorageMoverProjectResource{}
	stream.Results = func(push func(list.ListResult) bool) {
		for _, item := range resp.Items {
			result := request.NewListResult(ctx)
			result.DisplayName = pointer.From(item.Name)

			id, err := projects.ParseProjectIDInsensitively(pointer.From(item.Id))
			if err != nil {
				sdk.SetErrorDiagnosticAndPushListResult(result, push, "parsing Storage Mover Project ID", err)
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
