// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package compute

import (
	"context"
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/compute/2024-03-01/virtualmachines"
	"github.com/hashicorp/terraform-plugin-framework/list"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type WindowsVirtualMachineListResource struct{}

var _ sdk.FrameworkListWrappedResource = new(WindowsVirtualMachineListResource)

func (WindowsVirtualMachineListResource) ResourceFunc() *pluginsdk.Resource {
	return resourceWindowsVirtualMachine()
}

func (r WindowsVirtualMachineListResource) Metadata(_ context.Context, _ resource.MetadataRequest, response *resource.MetadataResponse) {
	response.TypeName = azureWindowsVirtualMachineResourceName
}

func (WindowsVirtualMachineListResource) List(ctx context.Context, request list.ListRequest, stream *list.ListResultsStream, metadata sdk.ResourceMetadata) {
	client := metadata.Client.Compute.VirtualMachinesClient

	// retrieve the deadline from the supplied context, since the List wrapper cancels it
	// before the iterator below runs and `flatten` makes additional API calls.
	deadline, ok := ctx.Deadline()
	if !ok {
		// This *should* never happen given the List Wrapper instantiates a context with a timeout
		sdk.SetResponseErrorDiagnostic(stream, "internal-error", "context had no deadline")
		return
	}

	var data sdk.DefaultListModel
	diags := request.Config.Get(ctx, &data)
	if diags.HasError() {
		stream.Results = list.ListResultsStreamDiagnostics(diags)
		return
	}

	subscriptionID := metadata.SubscriptionId
	if !data.SubscriptionId.IsNull() {
		subscriptionID = data.SubscriptionId.ValueString()
	}

	results := make([]virtualmachines.VirtualMachine, 0)

	switch {
	case !data.ResourceGroupName.IsNull():
		resp, err := client.ListComplete(ctx, commonids.NewResourceGroupID(subscriptionID, data.ResourceGroupName.ValueString()), virtualmachines.DefaultListOperationOptions())
		if err != nil {
			sdk.SetResponseErrorDiagnostic(stream, fmt.Sprintf("listing `%s`", azureWindowsVirtualMachineResourceName), err)
			return
		}

		results = resp.Items
	default:
		resp, err := client.ListAllComplete(ctx, commonids.NewSubscriptionID(subscriptionID), virtualmachines.DefaultListAllOperationOptions())
		if err != nil {
			sdk.SetResponseErrorDiagnostic(stream, fmt.Sprintf("listing `%s`", azureWindowsVirtualMachineResourceName), err)
			return
		}

		results = resp.Items
	}

	stream.Results = func(push func(list.ListResult) bool) {
		deadlineCtx, cancel := context.WithDeadline(context.Background(), deadline)
		defer cancel()

		for _, item := range results {
			// the API returns all Virtual Machines in scope, so filter out any that aren't Windows
			if props := item.Properties; props == nil || props.StorageProfile == nil || props.StorageProfile.OsDisk == nil || pointer.From(props.StorageProfile.OsDisk.OsType) != virtualmachines.OperatingSystemTypesWindows {
				continue
			}

			result := request.NewListResult(deadlineCtx)
			result.DisplayName = pointer.From(item.Name)

			rd := resourceWindowsVirtualMachine().Data(&terraform.InstanceState{})

			id, err := virtualmachines.ParseVirtualMachineIDInsensitively(pointer.From(item.Id))
			if err != nil {
				sdk.SetErrorDiagnosticAndPushListResult(result, push, fmt.Sprintf("parsing ID for `%s`", azureWindowsVirtualMachineResourceName), err)
				return
			}
			rd.SetId(id.ID())

			if err := resourceWindowsVirtualMachineFlatten(deadlineCtx, metadata.Client, rd, id, &item, request.IncludeResource); err != nil {
				sdk.SetErrorDiagnosticAndPushListResult(result, push, fmt.Sprintf("encoding `%s` resource data", azureWindowsVirtualMachineResourceName), err)
				return
			}

			sdk.EncodeListResult(deadlineCtx, rd, &result)
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
