// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package mssql

import (
	"context"
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/sqlvirtualmachine/2023-10-01/sqlvirtualmachines"
	"github.com/hashicorp/terraform-plugin-framework/list"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

var _ sdk.FrameworkServiceRegistration = Registration{}

type MssqlVirtualMachineListResource struct{}

var _ sdk.FrameworkListWrappedResource = new(MssqlVirtualMachineListResource)

func (r MssqlVirtualMachineListResource) ResourceFunc() *pluginsdk.Resource {
	return resourceMsSqlVirtualMachine()
}

func (r MssqlVirtualMachineListResource) Metadata(_ context.Context, _ resource.MetadataRequest, response *resource.MetadataResponse) {
	response.TypeName = `azurerm_mssql_virtual_machine`
}

func (r MssqlVirtualMachineListResource) List(ctx context.Context, request list.ListRequest, stream *list.ListResultsStream, metadata sdk.ResourceMetadata) {
	client := metadata.Client.MSSQL.VirtualMachinesClient

	var data sdk.DefaultListModel
	diags := request.Config.Get(ctx, &data)
	if diags.HasError() {
		stream.Results = list.ListResultsStreamDiagnostics(diags)
		return
	}

	results := make([]sqlvirtualmachines.SqlVirtualMachine, 0)

	subscriptionID := metadata.SubscriptionId
	if !data.SubscriptionId.IsNull() {
		subscriptionID = data.SubscriptionId.ValueString()
	}

	switch {
	case !data.ResourceGroupName.IsNull():
		resp, err := client.ListByResourceGroupComplete(ctx, commonids.NewResourceGroupID(subscriptionID, data.ResourceGroupName.ValueString()))
		if err != nil {
			sdk.SetResponseErrorDiagnostic(stream, fmt.Sprintf("listing `%s`", `azurerm_mssql_virtual_machine`), err)
			return
		}

		results = resp.Items
	default:
		resp, err := client.ListComplete(ctx, commonids.NewSubscriptionID(subscriptionID))
		if err != nil {
			sdk.SetResponseErrorDiagnostic(stream, fmt.Sprintf("listing `%s`", `azurerm_mssql_virtual_machine`), err)
			return
		}

		results = resp.Items
	}

	stream.Results = func(push func(list.ListResult) bool) {
		for _, virtualMachine := range results {

			result := request.NewListResult(ctx)

			result.DisplayName = pointer.From(virtualMachine.Name)

			rd := resourceMsSqlVirtualMachine().Data(&terraform.InstanceState{})

			id, err := sqlvirtualmachines.ParseSqlVirtualMachineID(pointer.From(virtualMachine.Id))
			if err != nil {
				sdk.SetErrorDiagnosticAndPushListResult(result, push, "parsing Mssql virtual machine ID", err)
				return
			}
			rd.SetId(id.ID())

			if err := resourceMssqlVirtualMachineSetFlatten(rd, id, &virtualMachine); err != nil {
				sdk.SetErrorDiagnosticAndPushListResult(result, push, fmt.Sprintf("encoding `%s` resource data", `azurerm_mssql_virtual_machine`), err)
				return
			}

			sdk.EncodeListResult(ctx, rd, &result)
		}
		return
	}
}
