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
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type (
	CustomDomainSignalrServiceListResource struct{}
	CustomDomainSignalrServiceListModel    struct {
		SignalRServiceId types.String `tfsdk:"signalr_service_id"`
	}
)

var _ sdk.FrameworkListWrappedResource = new(CustomDomainSignalrServiceListResource)

func (r CustomDomainSignalrServiceListResource) ResourceFunc() *pluginsdk.Resource {
	wrapper := sdk.NewResourceWrapper(CustomDomainSignalrServiceResource{})
	resource, err := wrapper.Resource()
	if err != nil {
		panic(fmt.Sprintf("building resource schema for `%s`: %+v", "azurerm_signalr_service_custom_domain", err))
	}

	return resource
}

func (r CustomDomainSignalrServiceListResource) Metadata(_ context.Context, _ resource.MetadataRequest, response *resource.MetadataResponse) {
	response.TypeName = "azurerm_signalr_service_custom_domain"
}

func (r CustomDomainSignalrServiceListResource) ListResourceConfigSchema(_ context.Context, _ list.ListResourceSchemaRequest, response *list.ListResourceSchemaResponse) {
	response.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"signalr_service_id": schema.StringAttribute{
				Required: true,
				Validators: []validator.String{
					typehelpers.WrappedStringValidator{Func: signalr.ValidateSignalRID},
				},
			},
		},
	}
}

func (r CustomDomainSignalrServiceListResource) List(ctx context.Context, request list.ListRequest, stream *list.ListResultsStream, metadata sdk.ResourceMetadata) {
	client := metadata.Client.SignalR.SignalRClient

	var data CustomDomainSignalrServiceListModel
	diags := request.Config.Get(ctx, &data)
	if diags.HasError() {
		stream.Results = list.ListResultsStreamDiagnostics(diags)
		return
	}

	signalRServiceId, err := signalr.ParseSignalRID(data.SignalRServiceId.ValueString())
	if err != nil {
		sdk.SetResponseErrorDiagnostic(stream, fmt.Sprintf("parsing SignalR Service ID for `%s`", "azurerm_signalr_service_custom_domain"), err)
		return
	}

	resp, err := client.CustomDomainsListComplete(ctx, *signalRServiceId)
	if err != nil {
		sdk.SetResponseErrorDiagnostic(stream, fmt.Sprintf("listing `%s`", "azurerm_signalr_service_custom_domain"), err)
		return
	}

	stream.Results = func(push func(list.ListResult) bool) {
		for _, domain := range resp.Items {
			result := request.NewListResult(ctx)
			result.DisplayName = pointer.From(domain.Name)

			id, err := signalr.ParseCustomDomainIDInsensitively(pointer.From(domain.Id))
			if err != nil {
				sdk.SetErrorDiagnosticAndPushListResult(result, push, "parsing SignalR Service Custom Domain ID", err)
				return
			}

			state, err := flattenCustomDomainSignalrServiceModel(*id, &domain)
			if err != nil {
				sdk.SetErrorDiagnosticAndPushListResult(result, push, "flattening SignalR Service Custom Domain", err)
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
			if err := rd.Set("signalr_service_id", state.SignalRServiceId); err != nil {
				sdk.SetErrorDiagnosticAndPushListResult(result, push, "setting signalr_service_id", err)
				return
			}
			if err := rd.Set("domain_name", state.DomainName); err != nil {
				sdk.SetErrorDiagnosticAndPushListResult(result, push, "setting domain_name", err)
				return
			}
			if err := rd.Set("signalr_custom_certificate_id", state.SignalrCustomCertificateId); err != nil {
				sdk.SetErrorDiagnosticAndPushListResult(result, push, "setting signalr_custom_certificate_id", err)
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
