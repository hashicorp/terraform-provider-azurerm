// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package eventgrid_partner_namespace_channel_topic

import (
	"context"
	"fmt"

	"github.com/hashicorp/go-azure-helpers/framework/typehelpers"
	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/eventgrid/2025-02-15/channels"
	"github.com/hashicorp/terraform-plugin-framework/list"
	"github.com/hashicorp/terraform-plugin-framework/list/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type ListResource struct{}

type EventGridPartnerNamespaceChannelTopicListModel struct {
	PartnerNamespaceId types.String `tfsdk:"partner_namespace_id"`
}

var _ sdk.FrameworkListWrappedResource = ListResource{}

func (r ListResource) ResourceFunc() *pluginsdk.Resource {
	return sdk.WrappedResource(Resource{})
}

func (r ListResource) Metadata(_ context.Context, _ resource.MetadataRequest, response *resource.MetadataResponse) {
	response.TypeName = Resource{}.ResourceType()
}

func (r ListResource) ListResourceConfigSchema(_ context.Context, _ list.ListResourceSchemaRequest, response *list.ListResourceSchemaResponse) {
	response.Schema = schema.Schema{Attributes: map[string]schema.Attribute{
		"partner_namespace_id": schema.StringAttribute{
			Required: true,
			Validators: []validator.String{
				typehelpers.WrappedStringValidator{Func: channels.ValidatePartnerNamespaceID},
			},
		},
	}}
}

func (r ListResource) List(ctx context.Context, request list.ListRequest, stream *list.ListResultsStream, metadata sdk.ResourceMetadata) {
	var config EventGridPartnerNamespaceChannelTopicListModel
	if diags := request.Config.Get(ctx, &config); diags.HasError() {
		stream.Results = list.ListResultsStreamDiagnostics(diags)
		return
	}
	parent, err := channels.ParsePartnerNamespaceID(config.PartnerNamespaceId.ValueString())
	if err != nil {
		sdk.SetResponseErrorDiagnostic(stream, "parsing Partner Namespace ID", err)
		return
	}
	resp, err := metadata.Client.EventGrid.Channels.ListByPartnerNamespaceComplete(ctx, *parent, channels.DefaultListByPartnerNamespaceOperationOptions())
	if err != nil {
		sdk.SetResponseErrorDiagnostic(stream, "listing Partner Topic channels", err)
		return
	}
	rc := Resource{}
	stream.Results = func(push func(list.ListResult) bool) {
		for _, channel := range resp.Items {
			if channel.Properties == nil || channel.Properties.ChannelType == nil || *channel.Properties.ChannelType != channels.ChannelTypePartnerTopic {
				continue
			}
			result := request.NewListResult(ctx)
			result.DisplayName = pointer.From(channel.Name)
			id, err := channels.ParseChannelID(pointer.From(channel.Id))
			if err != nil {
				sdk.SetErrorDiagnosticAndPushListResult(result, push, "parsing Channel ID", err)
				return
			}
			rmd := sdk.NewResourceMetaData(metadata.Client, rc)
			rmd.SetID(id)
			if err := rc.flatten(rmd, id, &channel); err != nil {
				sdk.SetErrorDiagnosticAndPushListResult(result, push, fmt.Sprintf("encoding `%s` resource data", rc.ResourceType()), err)
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
