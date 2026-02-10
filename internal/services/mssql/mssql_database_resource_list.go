// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package mssql

import (
	"context"
	"fmt"

	"github.com/hashicorp/go-azure-helpers/framework/typehelpers"
	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/sql/2023-08-01-preview/databases"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/list"
	"github.com/hashicorp/terraform-plugin-framework/list/schema"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

var _ sdk.FrameworkServiceRegistration = Registration{}

type MssqlDatabaseListResource struct{}

type MssqlDatabaseListModel struct {
	ServerId      types.String `tfsdk:"server_id"`
	ElasticPoolId types.String `tfsdk:"elastic_pool_id"`
}

var _ sdk.FrameworkListWrappedResource = new(MssqlDatabaseListResource)

func (r MssqlDatabaseListResource) ResourceFunc() *pluginsdk.Resource {
	return resourceMsSqlDatabase()
}

func (r MssqlDatabaseListResource) Metadata(_ context.Context, _ resource.MetadataRequest, response *resource.MetadataResponse) {
	response.TypeName = `azurerm_mssql_database`
}

func (r MssqlDatabaseListResource) ListResourceConfigSchema(_ context.Context, _ list.ListResourceSchemaRequest, response *list.ListResourceSchemaResponse) {
	response.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"server_id": schema.StringAttribute{
				Optional: true,
				Validators: []validator.String{
					typehelpers.WrappedStringValidator{
						Func: commonids.ValidateSqlServerID,
					},
					stringvalidator.ExactlyOneOf(path.MatchRoot("server_id"), path.MatchRoot("elastic_pool_id")),
				},
			},
			"elastic_pool_id": schema.StringAttribute{
				Optional: true,
				Validators: []validator.String{
					typehelpers.WrappedStringValidator{
						Func: commonids.ValidateSqlElasticPoolID,
					},
					stringvalidator.ExactlyOneOf(path.MatchRoot("server_id"), path.MatchRoot("elastic_pool_id")),
				},
			},
		},
	}
}

func (r MssqlDatabaseListResource) List(ctx context.Context, request list.ListRequest, stream *list.ListResultsStream, metadata sdk.ResourceMetadata) {
	client := metadata.Client.MSSQL.DatabasesClient

	var data MssqlDatabaseListModel
	diags := request.Config.Get(ctx, &data)
	if diags.HasError() {
		stream.Results = list.ListResultsStreamDiagnostics(diags)
		return
	}

	results := make([]databases.Database, 0)

	switch {
	case !data.ServerId.IsNull():
		serverId, err := commonids.ParseSqlServerID(data.ServerId.ValueString())
		if err != nil {
			sdk.SetResponseErrorDiagnostic(stream, fmt.Sprintf("parsing Mssql Server ID for `%s`", "azurerm_mssql_database"), err)
			return
		}
		resp, err := client.ListByServerComplete(ctx, *serverId)
		if err != nil {
			sdk.SetResponseErrorDiagnostic(stream, fmt.Sprintf("listing `%s`", `azurerm_mssql_database`), err)
			return
		}
		results = resp.Items

	case !data.ElasticPoolId.IsNull():
		elasticPoolId, err := commonids.ParseSqlElasticPoolID(data.ElasticPoolId.ValueString())
		if err != nil {
			sdk.SetResponseErrorDiagnostic(stream, fmt.Sprintf("parsing Mssql Elastic Pool ID for `%s`", "azurerm_mssql_database"), err)
			return
		}
		resp, err := client.ListByElasticPoolComplete(ctx, *elasticPoolId)
		if err != nil {
			sdk.SetResponseErrorDiagnostic(stream, fmt.Sprintf("listing `%s`", `azurerm_mssql_database`), err)
			return
		}
		results = resp.Items
	}

	stream.Results = func(push func(list.ListResult) bool) {
		for _, database := range results {
			// the default master database is special and will fail in the list function
			if *database.Name == "master" {
				continue
			}
			result := request.NewListResult(ctx)

			result.DisplayName = pointer.From(database.Name)

			rd := resourceMsSqlDatabase().Data(&terraform.InstanceState{})

			id, err := commonids.ParseSqlDatabaseID(pointer.From(database.Id))
			if err != nil {
				sdk.SetErrorDiagnosticAndPushListResult(result, push, "parsing Mssql Database ID", err)
				return
			}
			rd.SetId(id.ID())

			if err := resourceMssqlDatabaseSetFlatten(rd, id, &database, metadata.Client); err != nil {
				sdk.SetErrorDiagnosticAndPushListResult(result, push, fmt.Sprintf("encoding `%s` resource data", `azurerm_mssql_database`), err)
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
