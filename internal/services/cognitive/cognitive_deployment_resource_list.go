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

func cognitiveDeploymentListResourceConfigSchema() schema.Schema {
	return schema.Schema{
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

type CognitiveDeploymentListResource struct{}

var _ sdk.FrameworkListWrappedResource = new(CognitiveDeploymentListResource)

func (CognitiveDeploymentListResource) ResourceFunc() *pluginsdk.Resource {
	return sdk.WrappedResource(CognitiveDeploymentResource{})
}

func (CognitiveDeploymentListResource) Metadata(_ context.Context, _ resource.MetadataRequest, response *resource.MetadataResponse) {
	response.TypeName = CognitiveDeploymentResource{}.ResourceType()
}

func (CognitiveDeploymentListResource) ListResourceConfigSchema(_ context.Context, _ list.ListResourceSchemaRequest, response *list.ListResourceSchemaResponse) {
	response.Schema = cognitiveDeploymentListResourceConfigSchema()
}

func (CognitiveDeploymentListResource) List(ctx context.Context, request list.ListRequest, stream *list.ListResultsStream, metadata sdk.ResourceMetadata) {
	client := metadata.Client.Cognitive.DeploymentsClient

	var data cognitiveDeploymentListModel
	diags := request.Config.Get(ctx, &data)
	if diags.HasError() {
		stream.Results = list.ListResultsStreamDiagnostics(diags)
		return
	}

	deadline, ok := ctx.Deadline()
	if !ok {
		sdk.SetResponseErrorDiagnostic(stream, "internal-error", "context had no deadline")
		return
	}

	accountId, err := cognitiveservicesaccounts.ParseAccountID(data.CognitiveAccountId.ValueString())
	if err != nil {
		sdk.SetResponseErrorDiagnostic(stream, "parsing Cognitive Account ID", err)
		return
	}

	stream.Results = func(push func(list.ListResult) bool) {
		listCtx, cancel := context.WithDeadline(context.Background(), deadline)
		defer cancel()

		deploymentsResp, err := client.ListComplete(listCtx, deployments.NewAccountID(accountId.SubscriptionId, accountId.ResourceGroupName, accountId.AccountName))
		if err != nil {
			result := request.NewListResult(listCtx)
			sdk.SetErrorDiagnosticAndPushListResult(result, push, fmt.Sprintf("listing deployments for `%s`", accountId.AccountName), err)
			return
		}

		for _, deployment := range deploymentsResp.Items {
			deploymentId, err := deployments.ParseDeploymentID(pointer.From(deployment.Id))
			if err != nil {
				result := request.NewListResult(listCtx)
				sdk.SetErrorDiagnosticAndPushListResult(result, push, "parsing Cognitive Deployment ID", err)
				return
			}

			result := request.NewListResult(listCtx)
			result.DisplayName = pointer.From(deployment.Name)

			r := CognitiveDeploymentResource{}
			meta := sdk.NewResourceMetaData(metadata.Client, r)
			meta.SetID(deploymentId)

			if err := r.flatten(meta, deploymentId, &deployment); err != nil {
				sdk.SetErrorDiagnosticAndPushListResult(result, push, fmt.Sprintf("encoding `%s` resource data", CognitiveDeploymentResource{}.ResourceType()), err)
				return
			}

			sdk.EncodeListResult(listCtx, meta.ResourceData, &result)
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
