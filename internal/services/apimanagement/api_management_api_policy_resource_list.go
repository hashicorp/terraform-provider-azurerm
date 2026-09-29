// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package apimanagement

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/framework/typehelpers"
	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2022-08-01/apipolicy"
	"github.com/hashicorp/terraform-plugin-framework/list"
	"github.com/hashicorp/terraform-plugin-framework/list/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type ApiManagementApiPolicyListResource struct{}

type ApiManagementApiPolicyListModel struct {
	ApiId types.String `tfsdk:"api_management_api_id"`
}

var _ sdk.FrameworkListWrappedResource = new(ApiManagementApiPolicyListResource)

func (ApiManagementApiPolicyListResource) Metadata(_ context.Context, _ resource.MetadataRequest, response *resource.MetadataResponse) {
	response.TypeName = "azurerm_api_management_api_policy"
}

func (ApiManagementApiPolicyListResource) ResourceFunc() *pluginsdk.Resource {
	return resourceApiManagementApiPolicy()
}

func (ApiManagementApiPolicyListResource) ListResourceConfigSchema(_ context.Context, _ list.ListResourceSchemaRequest, response *list.ListResourceSchemaResponse) {
	response.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"api_management_api_id": schema.StringAttribute{
				Required: true,
				Validators: []validator.String{
					typehelpers.WrappedStringValidator{Func: apipolicy.ValidateApiID},
				},
			},
		},
	}
}

func (ApiManagementApiPolicyListResource) List(ctx context.Context, request list.ListRequest, stream *list.ListResultsStream, metadata sdk.ResourceMetadata) {
	client := metadata.Client.ApiManagement.ApiPoliciesClient

	var data ApiManagementApiPolicyListModel
	diags := request.Config.Get(ctx, &data)
	if diags.HasError() {
		stream.Results = list.ListResultsStreamDiagnostics(diags)
		return
	}

	parentID, err := apipolicy.ParseApiID(data.ApiId.ValueString())
	if err != nil {
		sdk.SetResponseErrorDiagnostic(stream, fmt.Sprintf("parsing parent ID for %s", azurermApiManagementApiPolicy), err)
		return
	}

	resp, err := client.ListByApiComplete(ctx, *parentID)
	if err != nil {
		sdk.SetResponseErrorDiagnostic(stream, fmt.Sprintf("listing %s", azurermApiManagementApiPolicy), err)
		return
	}

	var selected *apipolicy.PolicyContract
	var policy *apipolicy.PolicyContract
	for _, item := range resp.Items {
		if strings.EqualFold(pointer.From(item.Id), parentID.ID()+"/policies/policy") {
			selected = &item
			break
		}
	}
	if selected != nil {
		readResp, err := client.Get(ctx, *parentID, apipolicy.GetOperationOptions{Format: pointer.To(apipolicy.PolicyExportFormatXml)})
		if err != nil {
			sdk.SetResponseErrorDiagnostic(stream, fmt.Sprintf("reading %s", azurermApiManagementApiPolicy), err)
			return
		}
		policy = readResp.Model
	}

	// Note: there can only be 1 policy at the api_manangement_api level, so we only return 1 result if it exists
	stream.Results = func(push func(list.ListResult) bool) {
		result := request.NewListResult(ctx)
		result.DisplayName = pointer.From(policy.Name)
		rd := resourceApiManagementApiPolicy().Data(&terraform.InstanceState{})

		rd.SetId(parentID.ID())

		if err := resourceApiManagementApiPolicyFlatten(rd, parentID, policy); err != nil {
			sdk.SetErrorDiagnosticAndPushListResult(result, push, fmt.Sprintf("encoding `%s` resource data", azurermApiManagementApiPolicy), err)
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
