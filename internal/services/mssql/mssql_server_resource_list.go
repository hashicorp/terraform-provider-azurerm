// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package mssql

import (
	"context"
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/sql/2023-08-01-preview/servers"
	"github.com/hashicorp/terraform-plugin-framework/list"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type MssqlServerListResource struct{}

var _ sdk.FrameworkListWrappedResource = new(MssqlServerListResource)

func (r MssqlServerListResource) ResourceFunc() *pluginsdk.Resource {
	return resourceMsSqlServer()
}

func (r MssqlServerListResource) Metadata(_ context.Context, _ resource.MetadataRequest, response *resource.MetadataResponse) {
	response.TypeName = `azurerm_mssql_server`
}

func (r MssqlServerListResource) List(ctx context.Context, request list.ListRequest, stream *list.ListResultsStream, metadata sdk.ResourceMetadata) {
	client := metadata.Client.MSSQL.ServersClient

	var data sdk.DefaultListModel
	diags := request.Config.Get(ctx, &data)
	if diags.HasError() {
		stream.Results = list.ListResultsStreamDiagnostics(diags)
		return
	}

	results := make([]servers.Server, 0)

	subscriptionID := metadata.SubscriptionId
	if !data.SubscriptionId.IsNull() {
		subscriptionID = data.SubscriptionId.ValueString()
	}

	switch {
	case !data.ResourceGroupName.IsNull():
		resp, err := client.ListByResourceGroupComplete(ctx, commonids.NewResourceGroupID(subscriptionID, data.ResourceGroupName.ValueString()), servers.DefaultListByResourceGroupOperationOptions())
		if err != nil {
			sdk.SetResponseErrorDiagnostic(stream, fmt.Sprintf("listing `%s`", `azurerm_mssql_server`), err)
			return
		}

		results = resp.Items
	default:
		resp, err := client.ListComplete(ctx, commonids.NewSubscriptionID(subscriptionID), servers.DefaultListOperationOptions())
		if err != nil {
			sdk.SetResponseErrorDiagnostic(stream, fmt.Sprintf("listing `%s`", `azurerm_mssql_server`), err)
			return
		}

		results = resp.Items
	}

	stream.Results = func(push func(list.ListResult) bool) {
		for _, server := range results {

			result := request.NewListResult(ctx)

			result.DisplayName = pointer.From(server.Name)

			rd := resourceMsSqlServer().Data(&terraform.InstanceState{})

			id, err := commonids.ParseSqlServerID(pointer.From(server.Id))
			if err != nil {
				sdk.SetErrorDiagnosticAndPushListResult(result, push, "parsing Mssql Server ID", err)
				return
			}
			rd.SetId(id.ID())

			if err := resourceMssqlServerSetFlatten(rd, id, &server, metadata.Client); err != nil {
				sdk.SetErrorDiagnosticAndPushListResult(result, push, fmt.Sprintf("encoding `%s` resource data", `azurerm_mssql_server`), err)
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
