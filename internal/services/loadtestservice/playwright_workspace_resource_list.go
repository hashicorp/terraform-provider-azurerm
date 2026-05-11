// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package loadtestservice

import (
	"context"
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/loadtestservice/2025-09-01/playwrightworkspaces"
	"github.com/hashicorp/terraform-plugin-framework/list"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type PlaywrightWorkspaceListResource struct{}

var _ sdk.FrameworkListWrappedResource = new(PlaywrightWorkspaceListResource)

func (PlaywrightWorkspaceListResource) ResourceFunc() *pluginsdk.Resource {
	return sdk.WrappedResource(PlaywrightWorkspaceResource{})
}

func (PlaywrightWorkspaceListResource) Metadata(_ context.Context, _ resource.MetadataRequest, response *resource.MetadataResponse) {
	response.TypeName = PlaywrightWorkspaceResource{}.ResourceType()
}

func (PlaywrightWorkspaceListResource) List(ctx context.Context, request list.ListRequest, stream *list.ListResultsStream, metadata sdk.ResourceMetadata) {
	client := metadata.Client.LoadTestService.PlaywrightWorkspacesClient

	var data sdk.DefaultListModel
	diags := request.Config.Get(ctx, &data)
	if diags.HasError() {
		stream.Results = list.ListResultsStreamDiagnostics(diags)
		return
	}

	var results []playwrightworkspaces.PlaywrightWorkspace

	subscriptionID := metadata.Client.Account.SubscriptionId
	if !data.SubscriptionId.IsNull() {
		subscriptionID = data.SubscriptionId.ValueString()
	}

	r := PlaywrightWorkspaceResource{}

	switch {
	case !data.ResourceGroupName.IsNull():
		resp, err := client.ListByResourceGroupComplete(ctx, commonids.NewResourceGroupID(subscriptionID, data.ResourceGroupName.ValueString()))
		if err != nil {
			sdk.SetResponseErrorDiagnostic(stream, fmt.Sprintf("listing `%s`", r.ResourceType()), err)
			return
		}

		results = resp.Items
	default:
		resp, err := client.ListBySubscriptionComplete(ctx, commonids.NewSubscriptionID(subscriptionID))
		if err != nil {
			sdk.SetResponseErrorDiagnostic(stream, fmt.Sprintf("listing `%s`", r.ResourceType()), err)
			return
		}

		results = resp.Items
	}

	stream.Results = func(push func(list.ListResult) bool) {
		for _, playwrightWorkspaceResult := range results {
			result := request.NewListResult(ctx)
			result.DisplayName = pointer.From(playwrightWorkspaceResult.Name)

			id, err := playwrightworkspaces.ParsePlaywrightWorkspaceID(pointer.From(playwrightWorkspaceResult.Id))
			if err != nil {
				sdk.SetErrorDiagnosticAndPushListResult(result, push, "parsing Playwright Workspace ID", err)
				return
			}

			rmd := sdk.NewResourceMetaData(metadata.Client, r)
			rmd.SetID(id)

			if err := r.flatten(rmd, id, pointer.To(playwrightWorkspaceResult)); err != nil {
				sdk.SetErrorDiagnosticAndPushListResult(result, push, fmt.Sprintf("encoding `%s` resource data", r.ResourceType()), err)
				return
			}

			sdk.EncodeListResult(ctx, rmd.ResourceData, pointer.To(result))
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
