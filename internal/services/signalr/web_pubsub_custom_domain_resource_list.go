// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package signalr

import (
	"context"
	"fmt"

	"github.com/hashicorp/go-azure-helpers/framework/typehelpers"
	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/webpubsub/2024-03-01/webpubsub"
	"github.com/hashicorp/terraform-plugin-framework/list"
	"github.com/hashicorp/terraform-plugin-framework/list/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type (
	CustomDomainWebPubsubListResource struct{}
	CustomDomainWebPubsubListModel    struct {
		WebPubsubId types.String `tfsdk:"web_pubsub_id"`
	}
)

var _ sdk.FrameworkListWrappedResource = new(CustomDomainWebPubsubListResource)

func (r CustomDomainWebPubsubListResource) ResourceFunc() *pluginsdk.Resource {
	return sdk.WrappedResource(CustomDomainWebPubsubResource{})
}

func (r CustomDomainWebPubsubListResource) Metadata(_ context.Context, _ resource.MetadataRequest, response *resource.MetadataResponse) {
	response.TypeName = CustomDomainWebPubsubResource{}.ResourceType()
}

func (r CustomDomainWebPubsubListResource) ListResourceConfigSchema(_ context.Context, _ list.ListResourceSchemaRequest, response *list.ListResourceSchemaResponse) {
	response.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"web_pubsub_id": schema.StringAttribute{
				Required: true,
				Validators: []validator.String{
					typehelpers.WrappedStringValidator{Func: webpubsub.ValidateWebPubSubID},
				},
			},
		},
	}
}

func (r CustomDomainWebPubsubListResource) List(ctx context.Context, request list.ListRequest, stream *list.ListResultsStream, metadata sdk.ResourceMetadata) {
	client := metadata.Client.SignalR.WebPubSubClient.WebPubSub

	var data CustomDomainWebPubsubListModel
	if diags := request.Config.Get(ctx, &data); diags.HasError() {
		stream.Results = list.ListResultsStreamDiagnostics(diags)
		return
	}

	webPubsubId, err := webpubsub.ParseWebPubSubID(data.WebPubsubId.ValueString())
	if err != nil {
		sdk.SetResponseErrorDiagnostic(stream, fmt.Sprintf("parsing Web PubSub ID for `%s`", webPubsubCustomDomainResourceType), err)
		return
	}

	resp, err := client.CustomDomainsListComplete(ctx, *webPubsubId)
	if err != nil {
		sdk.SetResponseErrorDiagnostic(stream, fmt.Sprintf("listing `%s`", webPubsubCustomDomainResourceType), err)
		return
	}

	stream.Results = func(push func(list.ListResult) bool) {
		for _, domain := range resp.Items {
			result := request.NewListResult(ctx)
			result.DisplayName = pointer.From(domain.Name)

			id, err := webpubsub.ParseCustomDomainIDInsensitively(pointer.From(domain.Id))
			if err != nil {
				sdk.SetErrorDiagnosticAndPushListResult(result, push, "parsing Web PubSub Custom Domain ID", err)
				return
			}

			r := CustomDomainWebPubsubResource{}
			rmd := sdk.NewResourceMetaData(metadata.Client, r)
			rmd.ResourceData.SetId(id.ID())

			if err := r.flatten(rmd, id, &domain); err != nil {
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
