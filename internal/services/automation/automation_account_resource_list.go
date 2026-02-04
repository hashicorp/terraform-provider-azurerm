package automation

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/automation/2019-06-01/agentregistrationinformation"
	"github.com/hashicorp/go-azure-sdk/resource-manager/automation/2024-10-23/automationaccount"
	"github.com/hashicorp/terraform-plugin-framework/list"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type AutomationAccountListResource struct{}

var _ sdk.FrameworkListWrappedResource = new(AutomationAccountListResource)

func (r AutomationAccountListResource) ResourceFunc() *pluginsdk.Resource {
	return resourceAutomationAccount()
}

func (r AutomationAccountListResource) Metadata(_ context.Context, _ resource.MetadataRequest, response *resource.MetadataResponse) {
	response.TypeName = "azurerm_automation_account"
}

func (r AutomationAccountListResource) List(ctx context.Context, request list.ListRequest, stream *list.ListResultsStream, metadata sdk.ResourceMetadata) {
	client := metadata.Client.Automation.AutomationAccount

	var data sdk.DefaultListModel
	diags := request.Config.Get(ctx, &data)
	if diags.HasError() {
		stream.Results = list.ListResultsStreamDiagnostics(diags)
		return
	}

	var results []automationaccount.AutomationAccount

	subscriptionID := metadata.SubscriptionId
	if !data.SubscriptionId.IsNull() {
		subscriptionID = data.SubscriptionId.ValueString()
	}

	switch {
	case !data.ResourceGroupName.IsNull():
		resp, err := client.ListByResourceGroupComplete(ctx, commonids.NewResourceGroupID(subscriptionID, data.ResourceGroupName.ValueString()))
		if err != nil {
			sdk.SetResponseErrorDiagnostic(stream, fmt.Sprintf("listing `%s`", "azurerm_automation_account"), err)
			return
		}

		results = resp.Items
	default:
		resp, err := client.ListComplete(ctx, commonids.NewSubscriptionID(subscriptionID))
		if err != nil {
			sdk.SetResponseErrorDiagnostic(stream, fmt.Sprintf("listing `%s`", "azurerm_automation_account"), err)
			return
		}

		results = resp.Items
	}

	stream.Results = func(push func(list.ListResult) bool) {
		ctx, cancel := context.WithTimeout(context.Background(), time.Minute*5)
		defer cancel()
		for _, account := range results {
			result := request.NewListResult(ctx)
			result.DisplayName = pointer.From(account.Name)

			id, err := automationaccount.ParseAutomationAccountID(pointer.From(account.Id))
			if err != nil {
				sdk.SetListIteratorErrorDiagnostic(result, push, "parsing Automation Account ID", err)
				return
			}

			rd := resourceAutomationAccount().Data(&terraform.InstanceState{})
			rd.SetId(id.ID())

			var registration *agentregistrationinformation.AgentRegistration

			if request.IncludeResource {
				infoId := agentregistrationinformation.NewAutomationAccountID(id.SubscriptionId, id.ResourceGroupName, id.AutomationAccountName)
				if keysResp, err := metadata.Client.Automation.AgentRegistrationInfoClient.Get(ctx, infoId); err != nil {
					sdk.SetListIteratorErrorDiagnostic(result, push, "retrieving Automation Account Agent Registration Information", err)
					return
				} else {
					registration = keysResp.Model
				}
			}

			if err := resourceAutomationAccountFlatten(rd, id, &account, registration); err != nil {
				sdk.SetListIteratorErrorDiagnostic(result, push, fmt.Sprintf("encoding `%s` resource data", "azurerm_automation_account"), err)
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
