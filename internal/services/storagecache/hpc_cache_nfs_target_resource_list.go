// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package storagecache

import (
	"context"
	"fmt"

	"github.com/hashicorp/go-azure-helpers/framework/typehelpers"
	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/storagecache/2023-05-01/storagetargets"
	"github.com/hashicorp/terraform-plugin-framework/list"
	"github.com/hashicorp/terraform-plugin-framework/list/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type (
	HPCCacheNFSTargetListResource struct{}
	HPCCacheNFSTargetListModel    struct {
		CacheId types.String `tfsdk:"cache_id"`
	}
)

var _ sdk.FrameworkListWrappedResource = new(HPCCacheNFSTargetListResource)

func (r HPCCacheNFSTargetListResource) ResourceFunc() *pluginsdk.Resource {
	return resourceHPCCacheNFSTarget()
}

func (r HPCCacheNFSTargetListResource) Metadata(_ context.Context, _ resource.MetadataRequest, response *resource.MetadataResponse) {
	response.TypeName = "azurerm_hpc_cache_nfs_target"
}

func (r HPCCacheNFSTargetListResource) ListResourceConfigSchema(_ context.Context, _ list.ListResourceSchemaRequest, response *list.ListResourceSchemaResponse) {
	response.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"cache_id": schema.StringAttribute{
				Required: true,
				Validators: []validator.String{
					typehelpers.WrappedStringValidator{
						Func: storagetargets.ValidateCacheID,
					},
				},
			},
		},
	}
}

func (r HPCCacheNFSTargetListResource) List(ctx context.Context, request list.ListRequest, stream *list.ListResultsStream, metadata sdk.ResourceMetadata) {
	client := metadata.Client.StorageCache_2023_05_01.StorageTargets

	var data HPCCacheNFSTargetListModel
	if diags := request.Config.Get(ctx, &data); diags.HasError() {
		stream.Results = list.ListResultsStreamDiagnostics(diags)
		return
	}

	cacheId, err := storagetargets.ParseCacheID(data.CacheId.ValueString())
	if err != nil {
		sdk.SetResponseErrorDiagnostic(stream, fmt.Sprintf("parsing cache ID for `%s`", "azurerm_hpc_cache_nfs_target"), err)
		return
	}

	resp, err := client.ListByCacheComplete(ctx, *cacheId)
	if err != nil {
		sdk.SetResponseErrorDiagnostic(stream, fmt.Sprintf("listing `%s`", "azurerm_hpc_cache_nfs_target"), err)
		return
	}

	stream.Results = func(push func(list.ListResult) bool) {
		for _, target := range resp.Items {
			if target.Properties == nil || target.Properties.TargetType != storagetargets.StorageTargetTypeNfsThree {
				continue
			}

			result := request.NewListResult(ctx)
			result.DisplayName = pointer.From(target.Name)

			id, err := storagetargets.ParseStorageTargetIDInsensitively(pointer.From(target.Id))
			if err != nil {
				sdk.SetErrorDiagnosticAndPushListResult(result, push, "parsing Storage Target ID", err)
				return
			}

			rd := resourceHPCCacheNFSTarget().Data(&terraform.InstanceState{})
			rd.SetId(id.ID())

			if err := resourceHPCCacheNFSTargetFlatten(rd, id, &target); err != nil {
				sdk.SetErrorDiagnosticAndPushListResult(result, push, fmt.Sprintf("encoding `%s` resource data", "azurerm_hpc_cache_nfs_target"), err)
				return
			}

			sdk.EncodeListResult(ctx, rd, &result)
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
