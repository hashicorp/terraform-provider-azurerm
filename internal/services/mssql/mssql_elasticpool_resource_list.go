// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package mssql

import (
	"context"
	"fmt"

	"github.com/hashicorp/go-azure-helpers/framework/typehelpers"
	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/sql/2023-08-01-preview/elasticpools"
	"github.com/hashicorp/terraform-plugin-framework/list"
	"github.com/hashicorp/terraform-plugin-framework/list/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type MssqlElasticPoolListResource struct{}

type MssqlElasticPoolListModel struct {
	ServerId types.String `tfsdk:"server_id"`
}

var _ sdk.FrameworkListWrappedResource = new(MssqlElasticPoolListResource)

func (r MssqlElasticPoolListResource) ResourceFunc() *pluginsdk.Resource {
	return resourceMsSqlElasticPool()
}

func (r MssqlElasticPoolListResource) Metadata(_ context.Context, _ resource.MetadataRequest, response *resource.MetadataResponse) {
	response.TypeName = `azurerm_mssql_elasticpool`
}

func (r MssqlElasticPoolListResource) ListResourceConfigSchema(_ context.Context, _ list.ListResourceSchemaRequest, response *list.ListResourceSchemaResponse) {
	response.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"server_id": schema.StringAttribute{
				Required: true,
				Validators: []validator.String{
					typehelpers.WrappedStringValidator{
						Func: commonids.ValidateSqlServerID,
					},
				},
			},
		},
	}
}

func (r MssqlElasticPoolListResource) List(ctx context.Context, request list.ListRequest, stream *list.ListResultsStream, metadata sdk.ResourceMetadata) {
	client := metadata.Client.MSSQL.ElasticPoolsClient

	var data MssqlElasticPoolListModel
	diags := request.Config.Get(ctx, &data)
	if diags.HasError() {
		stream.Results = list.ListResultsStreamDiagnostics(diags)
		return
	}

	results := make([]elasticpools.ElasticPool, 0)
	if !data.ServerId.IsNull() {
		serverId, err := commonids.ParseSqlServerID(data.ServerId.ValueString())
		if err != nil {
			sdk.SetResponseErrorDiagnostic(stream, fmt.Sprintf("parsing Mssql Server ID for `%s`", "azurerm_mssql_elasticpool"), err)
			return
		}
		resp, err := client.ListByServerComplete(ctx, *serverId, elasticpools.DefaultListByServerOperationOptions())
		if err != nil {
			sdk.SetResponseErrorDiagnostic(stream, fmt.Sprintf("listing `%s`", `azurerm_mssql_elasticpool`), err)
			return
		}
		results = resp.Items
	}

	// Define the function that will push results into the stream
	stream.Results = func(push func(list.ListResult) bool) {
		for _, elasticPool := range results {

			// Initialize a new result object for each resource in the list
			result := request.NewListResult(ctx)

			// Set the display name of the item as the resource name
			result.DisplayName = pointer.From(elasticPool.Name)

			// Create a new ResourceData object to hold the state of the resource
			rd := resourceMsSqlElasticPool().Data(&terraform.InstanceState{})

			// Set the ID of the resource for the ResourceData object
			id, err := commonids.ParseSqlElasticPoolID(pointer.From(elasticPool.Id))
			if err != nil {
				sdk.SetListIteratorErrorDiagnostic(result, push, "parsing Mssql ElasticPool ID", err)
				return
			}
			rd.SetId(id.ID())

			// Use the resource flatten function to set the attributes into the resource state
			if err := resourceMssqlElasticPoolSetFlatten(rd, id, &elasticPool); err != nil {
				sdk.SetListIteratorErrorDiagnostic(result, push, fmt.Sprintf("encoding `%s` resource data", `azurerm_mssql_elasticpool`), err)
				return
			}

			sdk.EncodeListResult(ctx, rd, result, push)
		}
		return
	}
}
