package web

import (
	"context"
	"fmt"

	"github.com/hashicorp/go-azure-helpers/framework/typehelpers"
	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/web/2023-12-01/webapps"
	"github.com/hashicorp/terraform-plugin-framework/list"
	"github.com/hashicorp/terraform-plugin-framework/list/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type AppServiceCustomHostnameBindingListResource struct{}

type AppServiceCustomHostnameBindingListModel struct {
	AppServiceId types.String `tfsdk:"app_service_id"`
}

var _ sdk.FrameworkListWrappedResource = new(AppServiceCustomHostnameBindingListResource)

func (AppServiceCustomHostnameBindingListResource) Metadata(_ context.Context, _ resource.MetadataRequest, response *resource.MetadataResponse) {
	response.TypeName = appServiceCustomHostnameBindingResourceName
}

func (AppServiceCustomHostnameBindingListResource) ResourceFunc() *pluginsdk.Resource {
	return resourceAppServiceCustomHostnameBinding()
}

func (AppServiceCustomHostnameBindingListResource) ListResourceConfigSchema(_ context.Context, _ list.ListResourceSchemaRequest, response *list.ListResourceSchemaResponse) {
	response.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"app_service_id": schema.StringAttribute{
				Required: true,
				Validators: []validator.String{
					typehelpers.WrappedStringValidator{Func: commonids.ValidateAppServiceID},
				},
			},
		},
	}
}

func (AppServiceCustomHostnameBindingListResource) List(ctx context.Context, request list.ListRequest, stream *list.ListResultsStream, metadata sdk.ResourceMetadata) {
	client := metadata.Client.Web.WebAppsClient

	var data AppServiceCustomHostnameBindingListModel
	diags := request.Config.Get(ctx, &data)
	if diags.HasError() {
		stream.Results = list.ListResultsStreamDiagnostics(diags)
		return
	}

	parentID, err := commonids.ParseAppServiceID(data.AppServiceId.ValueString())
	if err != nil {
		sdk.SetResponseErrorDiagnostic(stream, fmt.Sprintf("parsing App Service ID for `%s`", appServiceCustomHostnameBindingResourceName), err)
		return
	}

	resp, err := client.ListHostNameBindingsComplete(ctx, *parentID)
	if err != nil {
		sdk.SetResponseErrorDiagnostic(stream, fmt.Sprintf("listing `%s`", appServiceCustomHostnameBindingResourceName), err)
		return
	}

	stream.Results = func(push func(list.ListResult) bool) {
		for _, item := range resp.Items {
			result := request.NewListResult(ctx)
			result.DisplayName = pointer.From(item.Name)

			rd := resourceAppServiceCustomHostnameBinding().Data(&terraform.InstanceState{})
			id, err := webapps.ParseHostNameBindingIDInsensitively(pointer.From(item.Id))
			if err != nil {
				sdk.SetErrorDiagnosticAndPushListResult(result, push, fmt.Sprintf("parsing ID for `%s`", appServiceCustomHostnameBindingResourceName), err)
				return
			}
			rd.SetId(id.ID())

			if err := resourceAppServiceCustomHostnameBindingFlatten(rd, id, &item); err != nil {
				sdk.SetErrorDiagnosticAndPushListResult(result, push, fmt.Sprintf("encoding `%s` resource data", appServiceCustomHostnameBindingResourceName), err)
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
