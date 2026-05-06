// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package netapp

import (
	"context"
	"fmt"

	"github.com/hashicorp/go-azure-helpers/framework/typehelpers"
	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/netapp/2026-01-01/buckets"
	"github.com/hashicorp/go-azure-sdk/resource-manager/netapp/2026-01-01/volumes"
	"github.com/hashicorp/terraform-plugin-framework/list"
	"github.com/hashicorp/terraform-plugin-framework/list/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type NetAppVolumeBucketListResource struct{}

type NetAppVolumeBucketListModel struct {
	VolumeId types.String `tfsdk:"volume_id"`
}

var _ sdk.FrameworkListWrappedResource = new(NetAppVolumeBucketListResource)

func (NetAppVolumeBucketListResource) Metadata(_ context.Context, _ resource.MetadataRequest, response *resource.MetadataResponse) {
	response.TypeName = NetAppVolumeBucketResource{}.ResourceType()
}

func (NetAppVolumeBucketListResource) ResourceFunc() *pluginsdk.Resource {
	return sdk.WrappedResource(NetAppVolumeBucketResource{})
}

func (NetAppVolumeBucketListResource) ListResourceConfigSchema(_ context.Context, _ list.ListResourceSchemaRequest, response *list.ListResourceSchemaResponse) {
	response.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"volume_id": schema.StringAttribute{
				Required: true,
				Validators: []validator.String{
					typehelpers.WrappedStringValidator{Func: volumes.ValidateVolumeID},
				},
			},
		},
	}
}

func (NetAppVolumeBucketListResource) List(ctx context.Context, request list.ListRequest, stream *list.ListResultsStream, metadata sdk.ResourceMetadata) {
	client := metadata.Client.NetApp.BucketsClient

	var data NetAppVolumeBucketListModel
	diags := request.Config.Get(ctx, &data)
	if diags.HasError() {
		stream.Results = list.ListResultsStreamDiagnostics(diags)
		return
	}

	volumeID, err := volumes.ParseVolumeID(data.VolumeId.ValueString())
	if err != nil {
		sdk.SetResponseErrorDiagnostic(stream, fmt.Sprintf("parsing Volume ID for `%s`", NetAppVolumeBucketResource{}.ResourceType()), err)
		return
	}

	resp, err := client.ListComplete(ctx, buckets.NewVolumeID(volumeID.SubscriptionId, volumeID.ResourceGroupName, volumeID.NetAppAccountName, volumeID.CapacityPoolName, volumeID.VolumeName))
	if err != nil {
		sdk.SetResponseErrorDiagnostic(stream, fmt.Sprintf("listing `%s`", NetAppVolumeBucketResource{}.ResourceType()), err)
		return
	}

	r := NetAppVolumeBucketResource{}
	stream.Results = func(push func(list.ListResult) bool) {
		for _, item := range resp.Items {
			result := request.NewListResult(ctx)
			result.DisplayName = pointer.From(item.Name)

			id, err := buckets.ParseBucketIDInsensitively(pointer.From(item.Id))
			if err != nil {
				sdk.SetErrorDiagnosticAndPushListResult(result, push, "parsing NetApp Volume Bucket ID", err)
				return
			}

			meta := sdk.NewResourceMetaData(metadata.Client, r)
			meta.SetID(id)

			if err := r.flatten(meta, id, &item); err != nil {
				sdk.SetErrorDiagnosticAndPushListResult(result, push, fmt.Sprintf("encoding `%s` resource data", r.ResourceType()), err)
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
