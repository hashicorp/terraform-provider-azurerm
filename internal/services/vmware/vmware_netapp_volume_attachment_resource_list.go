// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package vmware

import (
	"context"
	"fmt"

	"github.com/hashicorp/go-azure-helpers/framework/typehelpers"
	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/vmware/2022-05-01/datastores"
	"github.com/hashicorp/terraform-plugin-framework/list"
	"github.com/hashicorp/terraform-plugin-framework/list/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/vmware/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type NetappFileVolumeAttachmentListResource struct{}

type NetappFileVolumeAttachmentListModel struct {
	VmwareClusterId types.String `tfsdk:"vmware_cluster_id"`
}

var _ sdk.FrameworkListWrappedResource = new(NetappFileVolumeAttachmentListResource)

func (NetappFileVolumeAttachmentListResource) Metadata(_ context.Context, _ resource.MetadataRequest, response *resource.MetadataResponse) {
	response.TypeName = NetappFileVolumeAttachmentResource{}.ResourceType()
}

func (NetappFileVolumeAttachmentListResource) ResourceFunc() *pluginsdk.Resource {
	return sdk.WrappedResource(NetappFileVolumeAttachmentResource{})
}

func (NetappFileVolumeAttachmentListResource) ListResourceConfigSchema(_ context.Context, _ list.ListResourceSchemaRequest, response *list.ListResourceSchemaResponse) {
	response.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"vmware_cluster_id": schema.StringAttribute{
				Required: true,
				Validators: []validator.String{
					typehelpers.WrappedStringValidator{
						Func: validate.ClusterID,
					},
				},
			},
		},
	}
}

func (r NetappFileVolumeAttachmentListResource) List(ctx context.Context, request list.ListRequest, stream *list.ListResultsStream, metadata sdk.ResourceMetadata) {
	client := metadata.Client.Vmware.DataStoreClient

	var data NetappFileVolumeAttachmentListModel
	diags := request.Config.Get(ctx, &data)
	if diags.HasError() {
		stream.Results = list.ListResultsStreamDiagnostics(diags)
		return
	}

	clusterId, err := datastores.ParseClusterID(data.VmwareClusterId.ValueString())
	if err != nil {
		sdk.SetResponseErrorDiagnostic(stream, fmt.Sprintf("parsing `vmware_cluster_id` for `%s`", NetappFileVolumeAttachmentResource{}.ResourceType()), err)
		return
	}

	resp, err := client.ListComplete(ctx, *clusterId)
	if err != nil {
		sdk.SetResponseErrorDiagnostic(stream, fmt.Sprintf("listing `%s`", NetappFileVolumeAttachmentResource{}.ResourceType()), err)
		return
	}

	results := resp.Items

	resource := NetappFileVolumeAttachmentResource{}

	stream.Results = func(push func(list.ListResult) bool) {
		for _, datastore := range results {
			result := request.NewListResult(ctx)
			result.DisplayName = pointer.From(datastore.Name)

			id, err := datastores.ParseDataStoreID(pointer.From(datastore.Id))
			if err != nil {
				sdk.SetErrorDiagnosticAndPushListResult(result, push, "parsing VMware NetApp Volume Attachment ID", err)
				return
			}

			meta := sdk.NewResourceMetaData(metadata.Client, resource)
			meta.SetID(id)

			if err := resource.flatten(meta, id, &datastore); err != nil {
				sdk.SetErrorDiagnosticAndPushListResult(result, push, fmt.Sprintf("encoding `%s` resource data", pointer.From(datastore.Name)), err)
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
