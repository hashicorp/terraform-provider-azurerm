// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package voiceservices

import (
	"context"
	"fmt"

	"github.com/hashicorp/go-azure-helpers/framework/typehelpers"
	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/voiceservices/2023-04-03/communicationsgateways"
	"github.com/hashicorp/go-azure-sdk/resource-manager/voiceservices/2023-04-03/testlines"
	"github.com/hashicorp/terraform-plugin-framework/list"
	"github.com/hashicorp/terraform-plugin-framework/list/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type CommunicationsGatewayTestLineListResource struct{}

var _ sdk.FrameworkListWrappedResource = new(CommunicationsGatewayTestLineListResource)

func (CommunicationsGatewayTestLineListResource) ResourceFunc() *pluginsdk.Resource {
	return sdk.WrappedResource(CommunicationsGatewayTestLineResource{})
}

func (CommunicationsGatewayTestLineListResource) Metadata(_ context.Context, _ resource.MetadataRequest, response *resource.MetadataResponse) {
	response.TypeName = CommunicationsGatewayTestLineResource{}.ResourceType()
}

type CommunicationsGatewayTestLineListModel struct {
	VoiceServicesCommunicationsGatewayId types.String `tfsdk:"voice_services_communications_gateway_id"`
}

func (CommunicationsGatewayTestLineListResource) ListResourceConfigSchema(_ context.Context, _ list.ListResourceSchemaRequest, response *list.ListResourceSchemaResponse) {
	response.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"voice_services_communications_gateway_id": schema.StringAttribute{
				Required: true,
				Validators: []validator.String{
					typehelpers.WrappedStringValidator{
						Func: communicationsgateways.ValidateCommunicationsGatewayID,
					},
				},
			},
		},
	}
}

func (CommunicationsGatewayTestLineListResource) List(ctx context.Context, request list.ListRequest, stream *list.ListResultsStream, metadata sdk.ResourceMetadata) {
	client := metadata.Client.VoiceServices.TestLinesClient

	var data CommunicationsGatewayTestLineListModel
	diags := request.Config.Get(ctx, &data)
	if diags.HasError() {
		stream.Results = list.ListResultsStreamDiagnostics(diags)
		return
	}

	gatewayId, err := testlines.ParseCommunicationsGatewayID(data.VoiceServicesCommunicationsGatewayId.ValueString())
	if err != nil {
		sdk.SetResponseErrorDiagnostic(stream, "parsing `voice_services_communications_gateway_id`", err)
		return
	}

	r := CommunicationsGatewayTestLineResource{}

	resp, err := client.ListByCommunicationsGatewayComplete(ctx, *gatewayId)
	if err != nil {
		sdk.SetResponseErrorDiagnostic(stream, fmt.Sprintf("listing `%s`", r.ResourceType()), err)
		return
	}

	results := resp.Items

	stream.Results = func(push func(list.ListResult) bool) {
		for _, testLine := range results {
			result := request.NewListResult(ctx)
			result.DisplayName = pointer.From(testLine.Name)

			id, err := testlines.ParseTestLineID(pointer.From(testLine.Id))
			if err != nil {
				sdk.SetErrorDiagnosticAndPushListResult(result, push, "parsing Test Line ID", err)
				return
			}

			rmd := sdk.NewResourceMetaData(metadata.Client, r)
			rmd.SetID(id)

			if err := r.flatten(rmd, id, &testLine); err != nil {
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
