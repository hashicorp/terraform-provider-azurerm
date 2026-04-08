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
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type (
	CustomCertWebPubsubListResource struct{}
	CustomCertWebPubsubListModel    struct {
		WebPubsubId types.String `tfsdk:"web_pubsub_id"`
	}
)

var _ sdk.FrameworkListWrappedResource = new(CustomCertWebPubsubListResource)

func (r CustomCertWebPubsubListResource) ResourceFunc() *pluginsdk.Resource {
	wrapper := sdk.NewResourceWrapper(CustomCertWebPubsubResource{})
	resource, err := wrapper.Resource()
	if err != nil {
		panic(fmt.Sprintf("building resource schema for `%s`: %+v", "azurerm_web_pubsub_custom_certificate", err))
	}

	return resource
}

func (r CustomCertWebPubsubListResource) Metadata(_ context.Context, _ resource.MetadataRequest, response *resource.MetadataResponse) {
	response.TypeName = "azurerm_web_pubsub_custom_certificate"
}

func (r CustomCertWebPubsubListResource) ListResourceConfigSchema(_ context.Context, _ list.ListResourceSchemaRequest, response *list.ListResourceSchemaResponse) {
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

func (r CustomCertWebPubsubListResource) List(ctx context.Context, request list.ListRequest, stream *list.ListResultsStream, metadata sdk.ResourceMetadata) {
	client := metadata.Client.SignalR.WebPubSubClient.WebPubSub

	var data CustomCertWebPubsubListModel
	diags := request.Config.Get(ctx, &data)
	if diags.HasError() {
		stream.Results = list.ListResultsStreamDiagnostics(diags)
		return
	}

	webPubsubId, err := webpubsub.ParseWebPubSubID(data.WebPubsubId.ValueString())
	if err != nil {
		sdk.SetResponseErrorDiagnostic(stream, fmt.Sprintf("parsing Web PubSub ID for `%s`", "azurerm_web_pubsub_custom_certificate"), err)
		return
	}

	resp, err := client.CustomCertificatesListComplete(ctx, *webPubsubId)
	if err != nil {
		sdk.SetResponseErrorDiagnostic(stream, fmt.Sprintf("listing `%s`", "azurerm_web_pubsub_custom_certificate"), err)
		return
	}

	stream.Results = func(push func(list.ListResult) bool) {
		for _, cert := range resp.Items {
			result := request.NewListResult(ctx)
			result.DisplayName = pointer.From(cert.Name)

			id, err := webpubsub.ParseCustomCertificateIDInsensitively(pointer.From(cert.Id))
			if err != nil {
				sdk.SetErrorDiagnosticAndPushListResult(result, push, "parsing Web PubSub Custom Certificate ID", err)
				return
			}

			state, err := flattenCustomCertWebPubsubModel(*id, &cert)
			if err != nil {
				sdk.SetErrorDiagnosticAndPushListResult(result, push, "flattening Web PubSub Custom Certificate", err)
				return
			}

			rd := r.ResourceFunc().Data(&terraform.InstanceState{})
			rd.SetId(id.ID())
			if err := pluginsdk.SetResourceIdentityData(rd, id); err != nil {
				sdk.SetErrorDiagnosticAndPushListResult(result, push, "setting resource identity data", err)
				return
			}

			if err := rd.Set("name", state.Name); err != nil {
				sdk.SetErrorDiagnosticAndPushListResult(result, push, "setting name", err)
				return
			}
			if err := rd.Set("web_pubsub_id", state.WebPubsubId); err != nil {
				sdk.SetErrorDiagnosticAndPushListResult(result, push, "setting web_pubsub_id", err)
				return
			}
			if err := rd.Set("custom_certificate_id", state.CustomCertId); err != nil {
				sdk.SetErrorDiagnosticAndPushListResult(result, push, "setting custom_certificate_id", err)
				return
			}
			if err := rd.Set("certificate_version", state.CertificateVersion); err != nil {
				sdk.SetErrorDiagnosticAndPushListResult(result, push, "setting certificate_version", err)
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
