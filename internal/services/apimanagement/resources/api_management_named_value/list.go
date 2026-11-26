// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package api_management_named_value

import (
	"context"
	"fmt"

	"github.com/hashicorp/go-azure-helpers/framework/typehelpers"
	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2022-08-01/namedvalue"
	"github.com/hashicorp/terraform-plugin-framework/list"
	"github.com/hashicorp/terraform-plugin-framework/list/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type ListResource struct{}

type ListModel struct {
	ApiManagementId types.String `tfsdk:"api_management_id"`
}

var _ sdk.FrameworkListWrappedResource = new(ListResource)

func (ListResource) Metadata(_ context.Context, _ resource.MetadataRequest, response *resource.MetadataResponse) {
	response.TypeName = azureApiManagementNamedValueResourceName
}

func (ListResource) ResourceFunc() *pluginsdk.Resource {
	return RegisterResource()
}

func (ListResource) ListResourceConfigSchema(_ context.Context, _ list.ListResourceSchemaRequest, response *list.ListResourceSchemaResponse) {
	response.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"api_management_id": schema.StringAttribute{
				Required: true,
				Validators: []validator.String{
					typehelpers.WrappedStringValidator{Func: namedvalue.ValidateServiceID},
				},
			},
		},
	}
}

func (ListResource) List(ctx context.Context, request list.ListRequest, stream *list.ListResultsStream, metadata sdk.ResourceMetadata) {
	client := metadata.Client.ApiManagement.NamedValueClient

	var data ListModel
	diags := request.Config.Get(ctx, &data)
	if diags.HasError() {
		stream.Results = list.ListResultsStreamDiagnostics(diags)
		return
	}

	parentID, err := namedvalue.ParseServiceID(data.ApiManagementId.ValueString())
	if err != nil {
		sdk.SetResponseErrorDiagnostic(stream, fmt.Sprintf("parsing Api Management Service ID for `%s`", azureApiManagementNamedValueResourceName), err)
		return
	}

	resp, err := client.ListByServiceComplete(ctx, *parentID, namedvalue.DefaultListByServiceOperationOptions())
	if err != nil {
		sdk.SetResponseErrorDiagnostic(stream, fmt.Sprintf("listing `%s`", azureApiManagementNamedValueResourceName), err)
		return
	}

	stream.Results = func(push func(list.ListResult) bool) {
		for _, item := range resp.Items {
			result := request.NewListResult(ctx)
			result.DisplayName = pointer.From(item.Name)

			id, err := namedvalue.ParseNamedValueID(pointer.From(item.Id))
			if err != nil {
				sdk.SetErrorDiagnosticAndPushListResult(result, push, fmt.Sprintf("parsing ID for `%s`", azureApiManagementNamedValueResourceName), err)
				return
			}

			rd := RegisterResource().Data(&terraform.InstanceState{})
			rd.SetId(id.ID())

			if err := resourceApiManagementNamedValueFlatten(rd, id, &item); err != nil {
				sdk.SetErrorDiagnosticAndPushListResult(result, push, fmt.Sprintf("encoding `%s` resource data", azureApiManagementNamedValueResourceName), err)
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
