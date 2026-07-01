// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package trafficmanager

import (
	"context"
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/trafficmanager/2022-04-01/profiles"
	"github.com/hashicorp/terraform-plugin-framework/list"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type TrafficManagerProfileListResource struct{}

var _ sdk.FrameworkListWrappedResource = new(TrafficManagerProfileListResource)

func (r TrafficManagerProfileListResource) ResourceFunc() *pluginsdk.Resource {
	return resourceArmTrafficManagerProfile()
}

func (r TrafficManagerProfileListResource) Metadata(_ context.Context, _ resource.MetadataRequest, response *resource.MetadataResponse) {
	response.TypeName = azureTrafficManagerProfileResourceName
}

func (r TrafficManagerProfileListResource) List(ctx context.Context, request list.ListRequest, stream *list.ListResultsStream, metadata sdk.ResourceMetadata) {
	client := metadata.Client.TrafficManager.ProfilesClient

	var data sdk.DefaultListModel
	diags := request.Config.Get(ctx, &data)
	if diags.HasError() {
		stream.Results = list.ListResultsStreamDiagnostics(diags)
		return
	}

	results := make([]profiles.Profile, 0)

	subscriptionID := metadata.SubscriptionId
	if !data.SubscriptionId.IsNull() {
		subscriptionID = data.SubscriptionId.ValueString()
	}

	switch {
	case !data.ResourceGroupName.IsNull():
		resp, err := client.ListByResourceGroupComplete(ctx, commonids.NewResourceGroupID(subscriptionID, data.ResourceGroupName.ValueString()))
		if err != nil {
			sdk.SetResponseErrorDiagnostic(stream, fmt.Sprintf("listing `%s`", azureTrafficManagerProfileResourceName), err)
			return
		}

		results = resp.Items
	default:
		resp, err := client.ListBySubscriptionComplete(ctx, commonids.NewSubscriptionID(subscriptionID))
		if err != nil {
			sdk.SetResponseErrorDiagnostic(stream, fmt.Sprintf("listing `%s`", azureTrafficManagerProfileResourceName), err)
			return
		}

		results = resp.Items
	}

	stream.Results = func(push func(list.ListResult) bool) {
		for _, profile := range results {
			result := request.NewListResult(ctx)
			result.DisplayName = pointer.From(profile.Name)

			id, err := profiles.ParseTrafficManagerProfileID(pointer.From(profile.Id))
			if err != nil {
				sdk.SetErrorDiagnosticAndPushListResult(result, push, "parsing Traffic Manager Profile ID", err)
				return
			}

			rd := resourceArmTrafficManagerProfile().Data(&terraform.InstanceState{})
			rd.SetId(id.ID())

			if err := resourceArmTrafficManagerProfileFlatten(rd, id, &profile); err != nil {
				sdk.SetErrorDiagnosticAndPushListResult(result, push, fmt.Sprintf("encoding `%s` resource data", azureTrafficManagerProfileResourceName), err)
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
