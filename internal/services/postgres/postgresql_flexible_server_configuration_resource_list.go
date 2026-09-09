package postgres

import (
	"context"
	"fmt"

	"github.com/hashicorp/go-azure-helpers/framework/typehelpers"
	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/postgresql/2025-08-01/configurations"
	"github.com/hashicorp/go-azure-sdk/resource-manager/postgresql/2025-08-01/servers"
	"github.com/hashicorp/terraform-plugin-framework/list"
	"github.com/hashicorp/terraform-plugin-framework/list/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type PostgresqlFlexibleServerConfigurationListResource struct{}

type PostgresqlFlexibleServerConfigurationListModel struct {
	FlexibleServerId types.String `tfsdk:"server_id"`
}

var _ sdk.FrameworkListWrappedResource = new(PostgresqlFlexibleServerConfigurationListResource)

func (PostgresqlFlexibleServerConfigurationListResource) Metadata(_ context.Context, _ resource.MetadataRequest, response *resource.MetadataResponse) {
	response.TypeName = azurePostgresqlFlexibleServerConfigurationResourceName
}

func (PostgresqlFlexibleServerConfigurationListResource) ResourceFunc() *pluginsdk.Resource {
	return resourcePostgresqlFlexibleServerConfiguration()
}

func (PostgresqlFlexibleServerConfigurationListResource) ListResourceConfigSchema(_ context.Context, _ list.ListResourceSchemaRequest, response *list.ListResourceSchemaResponse) {
	response.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"server_id": schema.StringAttribute{
				Required: true,
				Validators: []validator.String{
					typehelpers.WrappedStringValidator{Func: servers.ValidateFlexibleServerID},
				},
			},
		},
	}
}

func (PostgresqlFlexibleServerConfigurationListResource) List(ctx context.Context, request list.ListRequest, stream *list.ListResultsStream, metadata sdk.ResourceMetadata) {
	client := metadata.Client.Postgres.FlexibleServersConfigurationsClient

	var data PostgresqlFlexibleServerConfigurationListModel
	diags := request.Config.Get(ctx, &data)
	if diags.HasError() {
		stream.Results = list.ListResultsStreamDiagnostics(diags)
		return
	}

	parentID, err := servers.ParseFlexibleServerID(data.FlexibleServerId.ValueString())
	if err != nil {
		sdk.SetResponseErrorDiagnostic(stream, fmt.Sprintf("parsing Flexible Server ID for `%s`", azurePostgresqlFlexibleServerConfigurationResourceName), err)
		return
	}

	resp, err := client.ListByServerComplete(ctx, configurations.NewFlexibleServerID(parentID.SubscriptionId, parentID.ResourceGroupName, parentID.FlexibleServerName))
	if err != nil {
		sdk.SetResponseErrorDiagnostic(stream, fmt.Sprintf("listing `%s`", azurePostgresqlFlexibleServerConfigurationResourceName), err)
		return
	}

	stream.Results = func(push func(list.ListResult) bool) {
		for _, item := range resp.Items {
			result := request.NewListResult(ctx)
			result.DisplayName = pointer.From(item.Name)

			rd := resourcePostgresqlFlexibleServerConfiguration().Data(&terraform.InstanceState{})
			id, err := configurations.ParseConfigurationIDInsensitively(pointer.From(item.Id))
			if err != nil {
				sdk.SetErrorDiagnosticAndPushListResult(result, push, fmt.Sprintf("parsing ID for `%s`", azurePostgresqlFlexibleServerConfigurationResourceName), err)
				return
			}
			rd.SetId(id.ID())

			if err := resourcePostgresqlFlexibleServerConfigurationFlatten(rd, id, &item); err != nil {
				sdk.SetErrorDiagnosticAndPushListResult(result, push, fmt.Sprintf("encoding `%s` resource data", azurePostgresqlFlexibleServerConfigurationResourceName), err)
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
