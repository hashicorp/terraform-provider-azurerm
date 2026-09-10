// Copyright IBM Corp.
// SPDX-License-Identifier: MPL-2.0

package network

import (
	"context"
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2025-01-01/ddosprotectionplans"
	"github.com/hashicorp/terraform-plugin-framework/list"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type NetworkDDoSProtectionPlanListResource struct{}

var _ sdk.FrameworkListWrappedResource = new(NetworkDDoSProtectionPlanListResource)

func (r NetworkDDoSProtectionPlanListResource) ResourceFunc() *pluginsdk.Resource {
	return resourceNetworkDDoSProtectionPlan()
}

func (r NetworkDDoSProtectionPlanListResource) Metadata(_ context.Context, _ resource.MetadataRequest, response *resource.MetadataResponse) {
	response.TypeName = "azurerm_network_ddos_protection_plan"
}

func (r NetworkDDoSProtectionPlanListResource) List(ctx context.Context, request list.ListRequest, stream *list.ListResultsStream, metadata sdk.ResourceMetadata) {
	client := metadata.Client.Network.DdosProtectionPlans
	var data sdk.DefaultListModel
	diags := request.Config.Get(ctx, &data)
	if diags.HasError() {
		stream.Results = list.ListResultsStreamDiagnostics(diags)
		return
	}

	results := make([]ddosprotectionplans.DdosProtectionPlan, 0)
	subscriptionID := metadata.SubscriptionId
	if !data.SubscriptionId.IsNull() {
		subscriptionID = data.SubscriptionId.ValueString()
	}

	switch {
	case !data.ResourceGroupName.IsNull():
		resp, err := client.ListByResourceGroupComplete(ctx, commonids.NewResourceGroupID(subscriptionID, data.ResourceGroupName.ValueString()))
		if err != nil {
			sdk.SetResponseErrorDiagnostic(stream, fmt.Sprintf("listing `%s`", `azurerm_network_ddos_protection_plan`), err)
			return
		}

		results = resp.Items
	default:
		resp, err := client.ListComplete(ctx, commonids.NewSubscriptionID(subscriptionID))
		if err != nil {
			sdk.SetResponseErrorDiagnostic(stream, fmt.Sprintf("listing `%s`", `azurerm_network_ddos_protection_plan`), err)
			return
		}

		results = resp.Items
	}

	stream.Results = func(push func(list.ListResult) bool) {
		for _, ddosprotectionplan := range results {
			result := request.NewListResult(ctx)

			result.DisplayName = pointer.From(ddosprotectionplan.Name)

			rd := resourceNetworkDDoSProtectionPlan().Data(&terraform.InstanceState{})

			id, err := ddosprotectionplans.ParseDdosProtectionPlanID(pointer.From(ddosprotectionplan.Id))
			if err != nil {
				sdk.SetErrorDiagnosticAndPushListResult(result, push, "parsing Network DDoS Protection Plan ID", err)
				return
			}

			rd.SetId(id.ID())

			if err := resourceNetworkDDoSProtectionPlanFlatten(rd, id, &ddosprotectionplan); err != nil {
				sdk.SetErrorDiagnosticAndPushListResult(result, push, fmt.Sprintf("encoding `%s` resource data", "azurerm_network_ddos_protection_plan"), err)
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
