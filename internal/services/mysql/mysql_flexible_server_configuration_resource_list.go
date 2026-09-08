package mysql

import (
	"context"
	"fmt"

	"github.com/hashicorp/go-azure-helpers/framework/typehelpers"
	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/mysql/2023-12-30/configurations"
	"github.com/hashicorp/terraform-plugin-framework/list"
	"github.com/hashicorp/terraform-plugin-framework/list/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type MysqlFlexibleServerConfigurationListResource struct{}

type MysqlFlexibleServerConfigurationListModel struct {
	FlexibleServerId types.String `tfsdk:"flexible_server_id"`
}

var _ sdk.FrameworkListWrappedResource = new(MysqlFlexibleServerConfigurationListResource)

func (r MysqlFlexibleServerConfigurationListResource) ResourceFunc() *pluginsdk.Resource {
	return resourceMySQLFlexibleServerConfiguration()
}

func (r MysqlFlexibleServerConfigurationListResource) Metadata(_ context.Context, _ resource.MetadataRequest, response *resource.MetadataResponse) {
	response.TypeName = mysqlFlexibleServerConfigurationResourceName
}

func (r MysqlFlexibleServerConfigurationListResource) ListResourceConfigSchema(_ context.Context, _ list.ListResourceSchemaRequest, response *list.ListResourceSchemaResponse) {
	response.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"flexible_server_id": schema.StringAttribute{
				Required: true,
				Validators: []validator.String{
					typehelpers.WrappedStringValidator{
						Func: configurations.ValidateFlexibleServerID,
					},
				},
			},
		},
	}
}

func (r MysqlFlexibleServerConfigurationListResource) List(ctx context.Context, request list.ListRequest, stream *list.ListResultsStream, metadata sdk.ResourceMetadata) {
	client := metadata.Client.MySQL.FlexibleServers.Configurations

	var data MysqlFlexibleServerConfigurationListModel
	diags := request.Config.Get(ctx, &data)
	if diags.HasError() {
		stream.Results = list.ListResultsStreamDiagnostics(diags)
		return
	}

	results := make([]configurations.Configuration, 0)

	serverId, err := configurations.ParseFlexibleServerID(data.FlexibleServerId.ValueString())
	if err != nil {
		sdk.SetResponseErrorDiagnostic(stream, fmt.Sprintf("parsing Mysql Server ID for `%s`", mysqlFlexibleServerConfigurationResourceName), err)
		return
	}

	resp, err := client.ListByServerComplete(ctx, pointer.From(serverId), configurations.DefaultListByServerOperationOptions())
	if err != nil {
		sdk.SetResponseErrorDiagnostic(stream, fmt.Sprintf("listing `%s`", mysqlFlexibleServerConfigurationResourceName), err)
		return
	}
	results = resp.Items

	stream.Results = func(push func(list.ListResult) bool) {
		for _, dbConfig := range results {
			result := request.NewListResult(ctx)
			result.DisplayName = pointer.From(dbConfig.Name)

			id, err := configurations.ParseConfigurationID(pointer.From(dbConfig.Id))
			if err != nil {
				sdk.SetErrorDiagnosticAndPushListResult(result, push, "parsing Mysql Configuration ID", err)
				return
			}

			rd := resourceMySQLFlexibleServerConfiguration().Data(&terraform.InstanceState{})
			rd.SetId(id.ID())

			if err := resourceMySQLFlexibleServerConfigurationFlatten(rd, id, &dbConfig); err != nil {
				sdk.SetErrorDiagnosticAndPushListResult(result, push, fmt.Sprintf("encoding `%s` resource data", mysqlFlexibleServerConfigurationResourceName), err)
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
