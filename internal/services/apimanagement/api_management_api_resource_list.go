package apimanagement

import (
	"context"
	"fmt"

	"github.com/hashicorp/go-azure-helpers/framework/typehelpers"
	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2022-08-01/api"
	"github.com/hashicorp/terraform-plugin-framework/list"
	"github.com/hashicorp/terraform-plugin-framework/list/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type ApiManagementApiListResource struct{}

type ApiManagementApiListModel struct {
	ServiceId types.String `tfsdk:"service_id"`
}

var _ sdk.FrameworkListWrappedResource = new(ApiManagementApiListResource)

func (ApiManagementApiListResource) Metadata(_ context.Context, _ resource.MetadataRequest, response *resource.MetadataResponse) {
	response.TypeName = azureApiManagementApiResourceName
}

func (ApiManagementApiListResource) ResourceFunc() *pluginsdk.Resource {
	return resourceApiManagementApi()
}

func (ApiManagementApiListResource) ListResourceConfigSchema(_ context.Context, _ list.ListResourceSchemaRequest, response *list.ListResourceSchemaResponse) {
	response.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"service_id": schema.StringAttribute{
				Required: true,
				Validators: []validator.String{
					typehelpers.WrappedStringValidator{Func: api.ValidateServiceID},
				},
			},
		},
	}
}

func (ApiManagementApiListResource) List(ctx context.Context, request list.ListRequest, stream *list.ListResultsStream, metadata sdk.ResourceMetadata) {
	client := metadata.Client.ApiManagement.ApiClient

	var data ApiManagementApiListModel
	diags := request.Config.Get(ctx, &data)
	if diags.HasError() {
		stream.Results = list.ListResultsStreamDiagnostics(diags)
		return
	}

	parentID, err := api.ParseServiceID(data.ServiceId.ValueString())
	if err != nil {
		sdk.SetResponseErrorDiagnostic(stream, fmt.Sprintf("parsing API Management Service ID for `%s`", azureApiManagementApiResourceName), err)
		return
	}

	resp, err := client.ListByServiceComplete(ctx, *parentID, api.DefaultListByServiceOperationOptions())
	if err != nil {
		sdk.SetResponseErrorDiagnostic(stream, fmt.Sprintf("listing `%s`", azureApiManagementApiResourceName), err)
		return
	}

	stream.Results = func(push func(list.ListResult) bool) {
		for _, item := range resp.Items {
			result := request.NewListResult(ctx)
			result.DisplayName = pointer.From(item.Name)

			rd := resourceApiManagementApi().Data(&terraform.InstanceState{})
			id, err := api.ParseApiIDInsensitively(pointer.From(item.Id))
			if err != nil {
				sdk.SetErrorDiagnosticAndPushListResult(result, push, fmt.Sprintf("parsing ID for `%s`", azureApiManagementApiResourceName), err)
				return
			}
			rd.SetId(id.ID())

			if err := resourceApiManagementApiFlatten(rd, id, &item); err != nil {
				sdk.SetErrorDiagnosticAndPushListResult(result, push, fmt.Sprintf("encoding `%s` resource data", azureApiManagementApiResourceName), err)
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
