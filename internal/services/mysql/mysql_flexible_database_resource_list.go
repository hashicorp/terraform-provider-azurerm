package mysql

import (
	"context"
	"fmt"

	"github.com/hashicorp/go-azure-helpers/framework/typehelpers"
	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/mysql/2023-12-30/databases"
	"github.com/hashicorp/terraform-plugin-framework/list"
	"github.com/hashicorp/terraform-plugin-framework/list/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type MysqlFlexibleDatabaseListResource struct{}

type MysqlFlexibleDatabaseListModel struct {
	FlexibleServerId types.String `tfsdk:"flexible_server_id"`
}

var _ sdk.FrameworkListWrappedResource = new(MysqlFlexibleDatabaseListResource)

func (r MysqlFlexibleDatabaseListResource) ResourceFunc() *pluginsdk.Resource {
	return resourceMySqlFlexibleDatabase()
}

func (r MysqlFlexibleDatabaseListResource) Metadata(_ context.Context, _ resource.MetadataRequest, response *resource.MetadataResponse) {
	response.TypeName = mysqlFlexibleDatabaseResourceName
}

func (r MysqlFlexibleDatabaseListResource) ListResourceConfigSchema(_ context.Context, _ list.ListResourceSchemaRequest, response *list.ListResourceSchemaResponse) {
	response.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"flexible_server_id": schema.StringAttribute{
				Required: true,
				Validators: []validator.String{
					typehelpers.WrappedStringValidator{
						Func: databases.ValidateFlexibleServerID,
					},
				},
			},
		},
	}
}

//goland:noinspection ALL
func (r MysqlFlexibleDatabaseListResource) List(ctx context.Context, request list.ListRequest, stream *list.ListResultsStream, metadata sdk.ResourceMetadata) {
	client := metadata.Client.MySQL.FlexibleServers.Databases

	var data MysqlFlexibleDatabaseListModel
	diags := request.Config.Get(ctx, &data)
	if diags.HasError() {
		stream.Results = list.ListResultsStreamDiagnostics(diags)
		return
	}

	results := make([]databases.Database, 0)

	serverId, err := databases.ParseFlexibleServerID(data.FlexibleServerId.ValueString())
	if err != nil {
		sdk.SetResponseErrorDiagnostic(stream, fmt.Sprintf("parsing Mysql Server ID for `%s`", mysqlFlexibleDatabaseResourceName), err)
		return
	}

	resp, err := client.ListByServerComplete(ctx, pointer.From(serverId))
	if err != nil {
		sdk.SetResponseErrorDiagnostic(stream, fmt.Sprintf("listing `%s`", mysqlFlexibleDatabaseResourceName), err)
		return
	}
	results = resp.Items

	stream.Results = func(push func(list.ListResult) bool) {
		for _, db := range results {
			result := request.NewListResult(ctx)
			result.DisplayName = pointer.From(db.Name)

			id, err := databases.ParseDatabaseID(pointer.From(db.Id))
			if err != nil {
				sdk.SetListIteratorErrorDiagnostic(result, push, "parsing Mysql Database ID", err)
				return
			}

			rd := resourceMySqlFlexibleDatabase().Data(&terraform.InstanceState{})
			rd.SetId(id.ID())

			if err := resourceMySqlFlexibleDatabaseSetResourceData(rd, id, &db); err != nil {
				sdk.SetListIteratorErrorDiagnostic(result, push, fmt.Sprintf("encoding `%s` resource data", mysqlFlexibleDatabaseResourceName), err)
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
