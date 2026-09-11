package network

import (
	"context"
	"fmt"

	"github.com/hashicorp/go-azure-helpers/framework/typehelpers"
	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2025-01-01/flowlogs"
	"github.com/hashicorp/terraform-plugin-framework/list"
	"github.com/hashicorp/terraform-plugin-framework/list/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type NetworkWatcherFlowLogListResource struct{}

type NetworkWatcherFlowLogListModel struct {
	NetworkWatcherId types.String `tfsdk:"network_watcher_id"`
}

var _ sdk.FrameworkListWrappedResource = new(NetworkWatcherFlowLogListResource)

func (NetworkWatcherFlowLogListResource) Metadata(_ context.Context, _ resource.MetadataRequest, response *resource.MetadataResponse) {
	response.TypeName = azureNetworkWatcherFlowLogResourceName
}

func (NetworkWatcherFlowLogListResource) ResourceFunc() *pluginsdk.Resource {
	return resourceNetworkWatcherFlowLog()
}

func (NetworkWatcherFlowLogListResource) ListResourceConfigSchema(_ context.Context, _ list.ListResourceSchemaRequest, response *list.ListResourceSchemaResponse) {
	response.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"network_watcher_id": schema.StringAttribute{
				Required: true,
				Validators: []validator.String{
					typehelpers.WrappedStringValidator{Func: flowlogs.ValidateNetworkWatcherID},
				},
			},
		},
	}
}

func (NetworkWatcherFlowLogListResource) List(ctx context.Context, request list.ListRequest, stream *list.ListResultsStream, metadata sdk.ResourceMetadata) {
	client := metadata.Client.Network.FlowLogs

	var data NetworkWatcherFlowLogListModel
	diags := request.Config.Get(ctx, &data)
	if diags.HasError() {
		stream.Results = list.ListResultsStreamDiagnostics(diags)
		return
	}

	parentID, err := flowlogs.ParseNetworkWatcherID(data.NetworkWatcherId.ValueString())
	if err != nil {
		sdk.SetResponseErrorDiagnostic(stream, fmt.Sprintf("parsing parent ID for `%s`", azureNetworkWatcherFlowLogResourceName), err)
		return
	}

	resp, err := client.ListComplete(ctx, *parentID)
	if err != nil {
		sdk.SetResponseErrorDiagnostic(stream, fmt.Sprintf("listing `%s`", azureNetworkWatcherFlowLogResourceName), err)
		return
	}

	stream.Results = func(push func(list.ListResult) bool) {
		for _, item := range resp.Items {
			result := request.NewListResult(ctx)
			result.DisplayName = pointer.From(item.Name)

			rd := resourceNetworkWatcherFlowLog().Data(&terraform.InstanceState{})

			id, err := flowlogs.ParseFlowLogIDInsensitively(pointer.From(item.Id))
			if err != nil {
				sdk.SetErrorDiagnosticAndPushListResult(result, push, fmt.Sprintf("parsing `%s` ID", azureNetworkWatcherFlowLogResourceName), err)
				return
			}
			rd.SetId(id.ID())

			if err := resourceNetworkWatcherFlowLogFlatten(rd, id, &item); err != nil {
				sdk.SetErrorDiagnosticAndPushListResult(result, push, fmt.Sprintf("encoding `%s` resource data", azureNetworkWatcherFlowLogResourceName), err)
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
