// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package cognitive

import (
	"context"
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/cognitive/2026-03-01/accountconnectionresource"
	"github.com/hashicorp/go-azure-sdk/resource-manager/cognitive/2026-03-01/cognitiveservicesaccounts"
	"github.com/hashicorp/terraform-plugin-framework/list"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type CognitiveAccountConnectionAccountKeyListResource struct{}

var _ sdk.FrameworkListWrappedResource = new(CognitiveAccountConnectionAccountKeyListResource)

func (CognitiveAccountConnectionAccountKeyListResource) ResourceFunc() *pluginsdk.Resource {
	return sdk.WrappedResource(CognitiveAccountConnectionAccountKeyResource{})
}

func (CognitiveAccountConnectionAccountKeyListResource) Metadata(_ context.Context, _ resource.MetadataRequest, response *resource.MetadataResponse) {
	response.TypeName = "azurerm_cognitive_account_connection_account_key"
}

func (CognitiveAccountConnectionAccountKeyListResource) List(ctx context.Context, request list.ListRequest, stream *list.ListResultsStream, metadata sdk.ResourceMetadata) {
	client := metadata.Client.Cognitive.AccountConnectionResourceClient

	var data sdk.DefaultListModel
	diags := request.Config.Get(ctx, &data)
	if diags.HasError() {
		stream.Results = list.ListResultsStreamDiagnostics(diags)
		return
	}

	if data.ResourceGroupName.IsNull() {
		sdk.SetResponseErrorDiagnostic(stream, "listing `azurerm_cognitive_account_connection_account_key`", "`resource_group_name` is required")
		return
	}

	deadline, ok := ctx.Deadline()
	if !ok {
		sdk.SetResponseErrorDiagnostic(stream, "internal-error", "context had no deadline")
		return
	}

	subscriptionID := metadata.SubscriptionId
	if !data.SubscriptionId.IsNull() {
		subscriptionID = data.SubscriptionId.ValueString()
	}

	accountsResp, err := metadata.Client.Cognitive.AccountsClient.AccountsListByResourceGroupComplete(ctx, commonids.NewResourceGroupID(subscriptionID, data.ResourceGroupName.ValueString()))
	if err != nil {
		sdk.SetResponseErrorDiagnostic(stream, "listing `azurerm_cognitive_account_connection_account_key`", err)
		return
	}

	stream.Results = func(push func(list.ListResult) bool) {
		listCtx, cancel := context.WithDeadline(context.Background(), deadline)
		defer cancel()

		for _, account := range accountsResp.Items {
			accountId, err := cognitiveservicesaccounts.ParseAccountID(pointer.From(account.Id))
			if err != nil {
				result := request.NewListResult(listCtx)
				sdk.SetErrorDiagnosticAndPushListResult(result, push, "parsing Cognitive Account ID", err)
				return
			}

			connectionsResp, err := client.AccountConnectionsListComplete(listCtx, accountconnectionresource.NewAccountID(accountId.SubscriptionId, accountId.ResourceGroupName, accountId.AccountName), accountconnectionresource.AccountConnectionsListOperationOptions{})
			if err != nil {
				result := request.NewListResult(listCtx)
				sdk.SetErrorDiagnosticAndPushListResult(result, push, fmt.Sprintf("listing connections for `%s`", accountId.AccountName), err)
				return
			}

			for _, connection := range connectionsResp.Items {
				base := connection.Properties.ConnectionPropertiesV2()
				if string(base.AuthType) != string(accountconnectionresource.ConnectionAuthTypeAccountKey) {
					continue
				}

				connectionId, err := accountconnectionresource.ParseConnectionID(pointer.From(connection.Id))
				if err != nil {
					result := request.NewListResult(listCtx)
					sdk.SetErrorDiagnosticAndPushListResult(result, push, "parsing Cognitive Account Connection ID", err)
					return
				}

				result := request.NewListResult(listCtx)
				result.DisplayName = pointer.From(connection.Name)

				rd := sdk.WrappedResource(CognitiveAccountConnectionAccountKeyResource{}).Data(&terraform.InstanceState{})
				rd.SetId(connectionId.ID())
				if err := pluginsdk.SetResourceIdentityData(rd, connectionId); err != nil {
					sdk.SetErrorDiagnosticAndPushListResult(result, push, "setting Cognitive Account Connection identity", err)
					return
				}
				_ = rd.Set("name", connectionId.ConnectionName)
				_ = rd.Set("cognitive_account_id", accountconnectionresource.NewAccountID(connectionId.SubscriptionId, connectionId.ResourceGroupName, connectionId.AccountName).ID())
				_ = rd.Set("category", pointer.FromEnum(base.Category))
				_ = rd.Set("target", pointer.From(base.Target))
				_ = rd.Set("metadata", pointer.From(base.Metadata))
				_ = rd.Set("account_key", "")

				sdk.EncodeListResult(listCtx, rd, &result)
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
}
