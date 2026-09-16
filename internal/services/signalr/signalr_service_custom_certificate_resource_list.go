// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package signalr

import (
	"context"
	"fmt"

	"github.com/hashicorp/go-azure-helpers/framework/typehelpers"
	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/signalr/2024-03-01/signalr"
	"github.com/hashicorp/terraform-plugin-framework/list"
	"github.com/hashicorp/terraform-plugin-framework/list/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type (
	CustomCertSignalrServiceListResource struct{}
	CustomCertSignalrServiceListModel    struct {
		SignalRServiceId types.String `tfsdk:"signalr_service_id"`
	}
)

var _ sdk.FrameworkListWrappedResource = new(CustomCertSignalrServiceListResource)

func (r CustomCertSignalrServiceListResource) ResourceFunc() *pluginsdk.Resource {
	return sdk.WrappedResource(CustomCertSignalrServiceResource{})
}

func (r CustomCertSignalrServiceListResource) Metadata(_ context.Context, _ resource.MetadataRequest, response *resource.MetadataResponse) {
	response.TypeName = CustomCertSignalrServiceResource{}.ResourceType()
}

func (r CustomCertSignalrServiceListResource) ListResourceConfigSchema(_ context.Context, _ list.ListResourceSchemaRequest, response *list.ListResourceSchemaResponse) {
	response.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"signalr_service_id": schema.StringAttribute{
				Required: true,
				Validators: []validator.String{
					typehelpers.WrappedStringValidator{
						Func: signalr.ValidateSignalRID,
					},
				},
			},
		},
	}
}

func (r CustomCertSignalrServiceListResource) List(ctx context.Context, request list.ListRequest, stream *list.ListResultsStream, metadata sdk.ResourceMetadata) {
	client := metadata.Client.SignalR.SignalRClient

	var data CustomCertSignalrServiceListModel
	diags := request.Config.Get(ctx, &data)
	if diags.HasError() {
		stream.Results = list.ListResultsStreamDiagnostics(diags)
		return
	}

	signalRServiceId, err := signalr.ParseSignalRID(data.SignalRServiceId.ValueString())
	if err != nil {
		sdk.SetResponseErrorDiagnostic(stream, fmt.Sprintf("parsing SignalR Service ID for `%s`", "azurerm_signalr_service_custom_certificate"), err)
		return
	}

	resp, err := client.CustomCertificatesListComplete(ctx, *signalRServiceId)
	if err != nil {
		sdk.SetResponseErrorDiagnostic(stream, fmt.Sprintf("listing `%s`", "azurerm_signalr_service_custom_certificate"), err)
		return
	}

	stream.Results = func(push func(list.ListResult) bool) {
		for _, cert := range resp.Items {
			result := request.NewListResult(ctx)
			result.DisplayName = pointer.From(cert.Name)

			id, err := signalr.ParseCustomCertificateIDInsensitively(pointer.From(cert.Id))
			if err != nil {
				sdk.SetErrorDiagnosticAndPushListResult(result, push, "parsing SignalR Service Custom Certificate ID", err)
				return
			}

			r := CustomCertSignalrServiceResource{}
			rmd := sdk.NewResourceMetaData(metadata.Client, r)
			rmd.ResourceData.SetId(id.ID())

			if err := r.flatten(rmd, id, &cert); err != nil {
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
