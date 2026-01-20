// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package mysql

import (
	"context"
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/mysql/2023-12-30/servers"
	"github.com/hashicorp/terraform-plugin-framework/list"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type MysqlFlexibleServerListResource struct{}

var _ sdk.FrameworkListWrappedResource = new(MysqlFlexibleServerListResource)

func (r MysqlFlexibleServerListResource) ResourceFunc() *pluginsdk.Resource {
	return resourceMysqlFlexibleServer()
}

func (r MysqlFlexibleServerListResource) Metadata(_ context.Context, _ resource.MetadataRequest, response *resource.MetadataResponse) {
	response.TypeName = `azurerm_mysql_flexible_server`
}

func (r MysqlFlexibleServerListResource) List(ctx context.Context, request list.ListRequest, stream *list.ListResultsStream, metadata sdk.ResourceMetadata) {
	client := metadata.Client.MySQL.FlexibleServers.Servers

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

	// mysql only supports subscription list
	resp, err := client.ListComplete(ctx, commonids.NewSubscriptionID(subscriptionID))
	if err != nil {
		sdk.SetResponseErrorDiagnostic(stream, fmt.Sprintf("listing `%s`", mysqlFlexibleServerResourceName), err)
		return
	}

	results = resp.Items

	stream.Results = func(push func(list.ListResult) bool) {
		for _, server := range results {
			result := request.NewListResult(ctx)
			result.DisplayName = pointer.From(server.Name)

			id, err := servers.ParseFlexibleServerID(pointer.From(server.Id))
			if err != nil {
				sdk.SetListIteratorErrorDiagnostic(result, push, "parsing Mysql Server ID", err)
				return
			}

			// filter out based on resource_group_name if supplied
			if !data.ResourceGroupName.IsNull() && !data.ResourceGroupName.Equal(basetypes.NewStringValue(id.ResourceGroupName)) {
				continue
			}

			rd := resourceMysqlFlexibleServer().Data(&terraform.InstanceState{})
			rd.SetId(id.ID())

			if err := resourceMysqlFlexibleServerSetResourceData(rd, id, &server, metadata); err != nil {
				sdk.SetListIteratorErrorDiagnostic(result, push, fmt.Sprintf("encoding `%s` resource data", mysqlFlexibleServerResourceName), err)
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
