// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package cognitive

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/cognitive/2025-06-01/cognitiveservicesaccounts"
	"github.com/hashicorp/terraform-plugin-framework/list"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type CognitiveAccountListResource struct{}

var _ sdk.FrameworkListWrappedResource = new(CognitiveAccountListResource)

func (r CognitiveAccountListResource) ResourceFunc() *pluginsdk.Resource {
	return resourceCognitiveAccount()
}

func (r CognitiveAccountListResource) Metadata(_ context.Context, _ resource.MetadataRequest, response *resource.MetadataResponse) {
	response.TypeName = azureCognitiveAccountResourceName
}

func (r CognitiveAccountListResource) List(ctx context.Context, request list.ListRequest, stream *list.ListResultsStream, metadata sdk.ResourceMetadata) {
	client := metadata.Client.Cognitive.AccountsClient

	var data sdk.DefaultListModel
	diags := request.Config.Get(ctx, &data)
	if diags.HasError() {
		stream.Results = list.ListResultsStreamDiagnostics(diags)
		return
	}

	results := make([]cognitiveservicesaccounts.Account, 0)

	subscriptionID := metadata.SubscriptionId
	if !data.SubscriptionId.IsNull() {
		subscriptionID = data.SubscriptionId.ValueString()
	}

	switch {
	case !data.ResourceGroupName.IsNull():
		resp, err := client.AccountsListByResourceGroupComplete(ctx, commonids.NewResourceGroupID(subscriptionID, data.ResourceGroupName.ValueString()))
		if err != nil {
			sdk.SetResponseErrorDiagnostic(stream, fmt.Sprintf("listing `%s`", azureCognitiveAccountResourceName), err)
			return
		}

		results = resp.Items
	default:
		resp, err := client.AccountsListComplete(ctx, commonids.NewSubscriptionID(subscriptionID))
		if err != nil {
			sdk.SetResponseErrorDiagnostic(stream, fmt.Sprintf("listing `%s`", azureCognitiveAccountResourceName), err)
			return
		}

		results = resp.Items
	}
	ctxDeadline, ok := ctx.Deadline()
	if !ok {
		sdk.SetResponseErrorDiagnostic(stream, "obtaining ctxDeadline", "no deadline set on the context")
		return
	}
	stream.Results = func(push func(list.ListResult) bool) {
		ctx, cancel := context.WithTimeout(context.Background(), time.Until(ctxDeadline))
		defer cancel()
		for _, account := range results {

			result := request.NewListResult(ctx)
			result.DisplayName = pointer.From(account.Name)

			rd := resourceCognitiveAccount().Data(&terraform.InstanceState{})

			id, err := cognitiveservicesaccounts.ParseAccountID(pointer.From(account.Id))
			if err != nil {
				sdk.SetListIteratorErrorDiagnostic(result, push, "parsing Cognitive Account ID", err)
				return
			}
			rd.SetId(id.ID())

			if err := resourceCognitiveAccountFlatten(ctx, client, rd, id, &account, request.IncludeResource); err != nil {
				sdk.SetListIteratorErrorDiagnostic(result, push, fmt.Sprintf("encoding `%s` resource data", azureCognitiveAccountResourceName), err)
				return
			}

			tfTypeIdentity, err := rd.TfTypeIdentityState()
			if err != nil {
				sdk.SetListIteratorErrorDiagnostic(result, push, "converting Identity State", err)
				return
			}

			if err := result.Identity.Set(ctx, *tfTypeIdentity); err != nil {
				sdk.SetListIteratorErrorDiagnostic(result, push, "setting Identity Data", err)
				return
			}

			tfTypeResourceState, err := rd.TfTypeResourceState()
			if err != nil {
				sdk.SetListIteratorErrorDiagnostic(result, push, "converting Resource State", err)
				return
			}

			if err := result.Resource.Set(ctx, *tfTypeResourceState); err != nil {
				sdk.SetListIteratorErrorDiagnostic(result, push, "setting Resource Data", err)
				return
			}

			if !push(result) {
				return
			}
		}
	}
}
