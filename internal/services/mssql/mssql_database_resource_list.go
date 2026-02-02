// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package mssql

import (
	"context"
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/sql/2023-08-01-preview/databases"
	"github.com/hashicorp/go-azure-sdk/resource-manager/sql/2023-08-01-preview/servers"
	"github.com/hashicorp/terraform-plugin-framework/list"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type MssqlDatabaseListResource struct{}

var _ sdk.FrameworkListWrappedResource = new(MssqlDatabaseListResource)

func (r MssqlDatabaseListResource) ResourceFunc() *pluginsdk.Resource {
	return resourceMsSqlDatabase()
}

func (r MssqlDatabaseListResource) Metadata(_ context.Context, _ resource.MetadataRequest, response *resource.MetadataResponse) {
	response.TypeName = `azurerm_mssql_D\database`
}

func (r MssqlDatabaseListResource) List(ctx context.Context, request list.ListRequest, stream *list.ListResultsStream, metadata sdk.ResourceMetadata) {
	client := metadata.Client.MSSQL.DatabasesClient

	var data sdk.DefaultListModel
	diags := request.Config.Get(ctx, &data)
	if diags.HasError() {
		stream.Results = list.ListResultsStreamDiagnostics(diags)
		return
	}

	results := make([]databases.Database, 0)

	subscriptionID := metadata.SubscriptionId
	if !data.SubscriptionId.IsNull() {
		subscriptionID = data.SubscriptionId.ValueString()
	}

	// Make the request based on which list parameters have been set in the config
	switch {
	//case !data.ResourceGroupName.IsNull():
	//	resp, err := client.ListByServerComplete(ctx, commonids.NewResourceGroupID(subscriptionID, data.ResourceGroupName.ValueString()), databases.DefaultListByResourceGroupOperationOptions())
	//	if err != nil {
	//		sdk.SetResponseErrorDiagnostic(stream, fmt.Sprintf("listing `%s`", `azurerm_mssql_database`), err)
	//		return
	//	}
	//
	//	results = resp.Items
	case !data.ResourceGroupName.IsNull():
		resp, err := client.ListByServerComplete(ctx, commonids.NewSqlServerID(subscriptionID, data.ResourceGroupName.ValueString(), "servername"))
		if err != nil {
			sdk.SetResponseErrorDiagnostic(stream, fmt.Sprintf("listing `%s`", `azurerm_mssql_database`), err)
			return
		}

		results = resp.Items
		//default:
		//	resp, err := client.ListComplete(ctx, commonids.NewSubscriptionID(subscriptionID), databases.DefaultListOperationOptions())
		//	if err != nil {
		//		sdk.SetResponseErrorDiagnostic(stream, fmt.Sprintf("listing `%s`", `azurerm_mssql_database`), err)
		//		return
		//	}
		//
		//	results = resp.Items
	}

	// Define the function that will push results into the stream
	stream.Results = func(push func(list.ListResult) bool) {
		for _, database := range results {

			// Initialize a new result object for each resource in the list
			result := request.NewListResult(ctx)

			// Set the display name of the item as the resource name
			result.DisplayName = pointer.From(database.Name)

			// Create a new ResourceData object to hold the state of the resource
			rd := resourceMsSqlDatabase().Data(&terraform.InstanceState{})

			// Set the ID of the resource for the ResourceData object
			id, err := commonids.ParseSqlDatabaseID(pointer.From(database.Id))
			if err != nil {
				sdk.SetListIteratorErrorDiagnostic(result, push, "parsing Mssql Database ID", err)
				return
			}
			rd.SetId(id.ID())

			// Use the resource flatten function to set the attributes into the resource state
			if err := resourceMssqlDatabaseSetFlatten(rd, id, &database, metadata.Client); err != nil {
				sdk.SetListIteratorErrorDiagnostic(result, push, fmt.Sprintf("encoding `%s` resource data", `azurerm_mssql_database`), err)
				return
			}

			sdk.EncodeListResult(ctx, rd, result, push)
		}
		return
	}
}
