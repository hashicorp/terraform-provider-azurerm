package loganalytics

import (
	"context"
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/operationalinsights/2025-07-01/workspaces"
	"github.com/hashicorp/terraform-plugin-framework/list"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type LogAnalyticsWorkspaceListResource struct{}

var _ sdk.FrameworkListWrappedResource = new(LogAnalyticsWorkspaceListResource)

func (LogAnalyticsWorkspaceListResource) Metadata(_ context.Context, _ resource.MetadataRequest, response *resource.MetadataResponse) {
	response.TypeName = azureLogAnalyticsWorkspaceName
}

func (LogAnalyticsWorkspaceListResource) ResourceFunc() *pluginsdk.Resource {
	return resourceLogAnalyticsWorkspace()
}

func (LogAnalyticsWorkspaceListResource) List(ctx context.Context, request list.ListRequest, stream *list.ListResultsStream, metadata sdk.ResourceMetadata) {
	client := metadata.Client.LogAnalytics.WorkspaceClient
	sharedKeyClient := metadata.Client.LogAnalytics.SharedKeyWorkspacesClient

	var data sdk.DefaultListModel
	diags := request.Config.Get(ctx, &data)
	if diags.HasError() {
		stream.Results = list.ListResultsStreamDiagnostics(diags)
		return
	}

	var results []workspaces.Workspace
	subscriptionID := metadata.SubscriptionId
	if !data.SubscriptionId.IsNull() {
		subscriptionID = data.SubscriptionId.ValueString()
	}

	switch {
	case !data.ResourceGroupName.IsNull():
		resp, err := client.ListByResourceGroupComplete(ctx, commonids.NewResourceGroupID(subscriptionID, data.ResourceGroupName.ValueString()))
		if err != nil {
			sdk.SetResponseErrorDiagnostic(stream, fmt.Sprintf("listing `%s`", azureLogAnalyticsWorkspaceName), err)
			return
		}
		results = resp.Items
	default:
		resp, err := client.ListComplete(ctx, commonids.NewSubscriptionID(subscriptionID))
		if err != nil {
			sdk.SetResponseErrorDiagnostic(stream, fmt.Sprintf("listing `%s`", azureLogAnalyticsWorkspaceName), err)
			return
		}
		results = resp.Items
	}

	deadline, ok := ctx.Deadline()
	if !ok {
		sdk.SetResponseErrorDiagnostic(stream, "internal-error", fmt.Errorf("context had no deadline"))
		return
	}

	stream.Results = func(push func(list.ListResult) bool) {
		ctx, cancel := context.WithDeadline(context.Background(), deadline)
		defer cancel()

		for _, item := range results {
			result := request.NewListResult(ctx)
			result.DisplayName = pointer.From(item.Name)

			rd := resourceLogAnalyticsWorkspace().Data(&terraform.InstanceState{})
			id, err := workspaces.ParseWorkspaceIDInsensitively(pointer.From(item.Id))
			if err != nil {
				sdk.SetErrorDiagnosticAndPushListResult(result, push, fmt.Sprintf("parsing `%s` ID", azureLogAnalyticsWorkspaceName), err)
				return
			}
			rd.SetId(id.ID())

			if err := resourceLogAnalyticsWorkspaceFlatten(ctx, sharedKeyClient, rd, id, &item, request.IncludeResource); err != nil {
				sdk.SetErrorDiagnosticAndPushListResult(result, push, fmt.Sprintf("encoding `%s` resource data", azureLogAnalyticsWorkspaceName), err)
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
