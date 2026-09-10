package apimanagement

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/framework/typehelpers"
	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2022-08-01/apioperationpolicy"
	"github.com/hashicorp/terraform-plugin-framework/list"
	"github.com/hashicorp/terraform-plugin-framework/list/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type ApiManagementApiOperationPolicyListResource struct{}

type ApiManagementApiOperationPolicyListModel struct {
	OperationId types.String `tfsdk:"operation_id"`
}

var _ sdk.FrameworkListWrappedResource = new(ApiManagementApiOperationPolicyListResource)

func (ApiManagementApiOperationPolicyListResource) Metadata(_ context.Context, _ resource.MetadataRequest, response *resource.MetadataResponse) {
	response.TypeName = azureApiManagementApiOperationPolicyResourceName
}

func (ApiManagementApiOperationPolicyListResource) ResourceFunc() *pluginsdk.Resource {
	return resourceApiManagementApiOperationPolicy()
}

func (ApiManagementApiOperationPolicyListResource) ListResourceConfigSchema(_ context.Context, _ list.ListResourceSchemaRequest, response *list.ListResourceSchemaResponse) {
	response.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"operation_id": schema.StringAttribute{
				Required: true,
				Validators: []validator.String{
					typehelpers.WrappedStringValidator{Func: apioperationpolicy.ValidateOperationID},
				},
			},
		},
	}
}

func (ApiManagementApiOperationPolicyListResource) List(ctx context.Context, request list.ListRequest, stream *list.ListResultsStream, metadata sdk.ResourceMetadata) {
	client := metadata.Client.ApiManagement.ApiOperationPoliciesClient

	var data ApiManagementApiOperationPolicyListModel
	diags := request.Config.Get(ctx, &data)
	if diags.HasError() {
		stream.Results = list.ListResultsStreamDiagnostics(diags)
		return
	}

	parentID, err := apioperationpolicy.ParseOperationID(data.OperationId.ValueString())
	if err != nil {
		sdk.SetResponseErrorDiagnostic(stream, fmt.Sprintf("parsing Operation ID for `%s`", azureApiManagementApiOperationPolicyResourceName), err)
		return
	}

	resp, err := client.ListByOperationComplete(ctx, *parentID)
	if err != nil {
		sdk.SetResponseErrorDiagnostic(stream, fmt.Sprintf("listing `%s`", azureApiManagementApiOperationPolicyResourceName), err)
		return
	}

	stream.Results = func(push func(list.ListResult) bool) {
		for _, item := range resp.Items {
			result := request.NewListResult(ctx)
			result.DisplayName = pointer.From(item.Name)

			rd := resourceApiManagementApiOperationPolicy().Data(&terraform.InstanceState{})
			id, err := apioperationpolicy.ParseOperationIDInsensitively(strings.TrimSuffix(pointer.From(item.Id), "/policies/policy"))
			if err != nil {
				sdk.SetErrorDiagnosticAndPushListResult(result, push, fmt.Sprintf("parsing ID for `%s`", azureApiManagementApiOperationPolicyResourceName), err)
				return
			}
			rd.SetId(id.ID())

			if err := resourceApiManagementAPIOperationPolicyFlatten(rd, id, &item); err != nil {
				sdk.SetErrorDiagnosticAndPushListResult(result, push, fmt.Sprintf("encoding `%s` resource data", azureApiManagementApiOperationPolicyResourceName), err)
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
