// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package discovery

import (
	"context"
	"fmt"

	"github.com/hashicorp/go-azure-helpers/framework/typehelpers"
	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/discovery/2026-06-01/bookshelfprivateendpointconnections"
	"github.com/hashicorp/terraform-plugin-framework/list"
	"github.com/hashicorp/terraform-plugin-framework/list/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type DiscoveryBookshelfPrivateEndpointConnectionListResource struct{}

type DiscoveryBookshelfPrivateEndpointConnectionListModel struct {
	BookshelfId types.String `tfsdk:"bookshelf_id"`
}

var _ sdk.FrameworkListWrappedResource = new(DiscoveryBookshelfPrivateEndpointConnectionListResource)

func (DiscoveryBookshelfPrivateEndpointConnectionListResource) Metadata(_ context.Context, _ resource.MetadataRequest, response *resource.MetadataResponse) {
	response.TypeName = DiscoveryBookshelfPrivateEndpointConnectionResource{}.ResourceType()
}

func (DiscoveryBookshelfPrivateEndpointConnectionListResource) ResourceFunc() *pluginsdk.Resource {
	return sdk.WrappedResource(DiscoveryBookshelfPrivateEndpointConnectionResource{})
}

func (DiscoveryBookshelfPrivateEndpointConnectionListResource) ListResourceConfigSchema(_ context.Context, _ list.ListResourceSchemaRequest, response *list.ListResourceSchemaResponse) {
	response.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"bookshelf_id": schema.StringAttribute{
				Required: true,
				Validators: []validator.String{
					typehelpers.WrappedStringValidator{Func: bookshelfprivateendpointconnections.ValidateBookshelfID},
				},
			},
		},
	}
}

func (DiscoveryBookshelfPrivateEndpointConnectionListResource) List(ctx context.Context, request list.ListRequest, stream *list.ListResultsStream, metadata sdk.ResourceMetadata) {
	client := metadata.Client.Discovery.BookshelfPrivateEndpointConnections

	var data DiscoveryBookshelfPrivateEndpointConnectionListModel
	diags := request.Config.Get(ctx, &data)
	if diags.HasError() {
		stream.Results = list.ListResultsStreamDiagnostics(diags)
		return
	}

	r := DiscoveryBookshelfPrivateEndpointConnectionResource{}

	bookshelfId, err := bookshelfprivateendpointconnections.ParseBookshelfID(data.BookshelfId.ValueString())
	if err != nil {
		sdk.SetResponseErrorDiagnostic(stream, fmt.Sprintf("parsing Discovery Bookshelf ID for `%s`", r.ResourceType()), err)
		return
	}

	resp, err := client.ListByBookshelfComplete(ctx, *bookshelfId)
	if err != nil {
		sdk.SetResponseErrorDiagnostic(stream, fmt.Sprintf("listing `%s`", r.ResourceType()), err)
		return
	}

	stream.Results = func(push func(list.ListResult) bool) {
		for _, item := range resp.Items {
			result := request.NewListResult(ctx)
			result.DisplayName = pointer.From(item.Name)

			id, err := bookshelfprivateendpointconnections.ParseBookshelfPrivateEndpointConnectionIDInsensitively(pointer.From(item.Id))
			if err != nil {
				sdk.SetErrorDiagnosticAndPushListResult(result, push, "parsing Discovery Bookshelf Private Endpoint Connection ID", err)
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
