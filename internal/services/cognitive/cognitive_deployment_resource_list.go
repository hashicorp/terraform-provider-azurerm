// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package cognitive

import (
	"context"
	"fmt"

	"github.com/hashicorp/go-azure-helpers/framework/typehelpers"
	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/cognitive/2026-03-01/cognitiveservicesaccounts"
	"github.com/hashicorp/go-azure-sdk/resource-manager/cognitive/2026-03-01/deployments"
	"github.com/hashicorp/terraform-plugin-framework/list"
	"github.com/hashicorp/terraform-plugin-framework/list/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type cognitiveDeploymentListModel struct {
	CognitiveAccountId types.String `tfsdk:"cognitive_account_id"`
}

type CognitiveDeploymentListResource struct{}

var _ sdk.FrameworkListWrappedResource = new(CognitiveDeploymentListResource)

func (CognitiveDeploymentListResource) ResourceFunc() *pluginsdk.Resource {
	return sdk.WrappedResource(CognitiveDeploymentResource{})
}

func (CognitiveDeploymentListResource) Metadata(_ context.Context, _ resource.MetadataRequest, response *resource.MetadataResponse) {
	response.TypeName = CognitiveDeploymentResource{}.ResourceType()
}

func (CognitiveDeploymentListResource) ListResourceConfigSchema(_ context.Context, _ list.ListResourceSchemaRequest, response *list.ListResourceSchemaResponse) {
	response.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"cognitive_account_id": schema.StringAttribute{
				Required: true,
				Validators: []validator.String{
					typehelpers.WrappedStringValidator{
						Func: cognitiveservicesaccounts.ValidateAccountID,
					},
				},
			},
		},
	}
}

func (CognitiveDeploymentListResource) List(ctx context.Context, request list.ListRequest, stream *list.ListResultsStream, metadata sdk.ResourceMetadata) {
	client := metadata.Client.Cognitive.DeploymentsClient

	var data cognitiveDeploymentListModel
	diags := request.Config.Get(ctx, &data)
	if diags.HasError() {
		stream.Results = list.ListResultsStreamDiagnostics(diags)
		return
	}

	accountId, err := deployments.ParseAccountID(data.CognitiveAccountId.ValueString())
	if err != nil {
		sdk.SetResponseErrorDiagnostic(stream, "parsing Cognitive Account ID", err)
		return
	}

	deploymentsResp, err := client.ListComplete(ctx, *accountId)
	if err != nil {
		sdk.SetResponseErrorDiagnostic(stream, fmt.Sprintf("listing deployments for `%s`", accountId.AccountName), err)
		return
	}

	r := CognitiveDeploymentResource{}
	stream.Results = func(push func(list.ListResult) bool) {
		for _, deployment := range deploymentsResp.Items {
			deploymentId, err := deployments.ParseDeploymentID(pointer.From(deployment.Id))
			if err != nil {
				result := request.NewListResult(ctx)
				sdk.SetErrorDiagnosticAndPushListResult(result, push, "parsing Cognitive Deployment ID", err)
				return
			}

			result := request.NewListResult(ctx)
			result.DisplayName = pointer.From(deployment.Name)

			meta := sdk.NewResourceMetaData(metadata.Client, r)
			meta.SetID(deploymentId)

			if err := r.flatten(meta, deploymentId, &deployment); err != nil {
				sdk.SetErrorDiagnosticAndPushListResult(result, push, fmt.Sprintf("encoding `%s` resource data", r.ResourceType()), err)
				return
			}

			sdk.EncodeListResult(ctx, meta.ResourceData, &result)
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
