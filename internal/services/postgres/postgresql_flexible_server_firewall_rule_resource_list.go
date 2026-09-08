package postgres

import (
	"context"
	"fmt"

	"github.com/hashicorp/go-azure-helpers/framework/typehelpers"
	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/postgresql/2025-08-01/firewallrules"
	"github.com/hashicorp/terraform-plugin-framework/list"
	"github.com/hashicorp/terraform-plugin-framework/list/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type PostgresqlFlexibleServerFirewallRuleListResource struct{}

type PostgresqlFlexibleServerFirewallRuleListModel struct {
	FlexibleServerId types.String `tfsdk:"server_id"`
}

var _ sdk.FrameworkListWrappedResource = new(PostgresqlFlexibleServerFirewallRuleListResource)

func (PostgresqlFlexibleServerFirewallRuleListResource) Metadata(_ context.Context, _ resource.MetadataRequest, response *resource.MetadataResponse) {
	response.TypeName = azurePostgresqlFlexibleServerFirewallRuleResourceName
}

func (PostgresqlFlexibleServerFirewallRuleListResource) ResourceFunc() *pluginsdk.Resource {
	return resourcePostgresqlFlexibleServerFirewallRule()
}

func (PostgresqlFlexibleServerFirewallRuleListResource) ListResourceConfigSchema(_ context.Context, _ list.ListResourceSchemaRequest, response *list.ListResourceSchemaResponse) {
	response.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"server_id": schema.StringAttribute{
				Required: true,
				Validators: []validator.String{
					typehelpers.WrappedStringValidator{Func: firewallrules.ValidateFlexibleServerID},
				},
			},
		},
	}
}

func (PostgresqlFlexibleServerFirewallRuleListResource) List(ctx context.Context, request list.ListRequest, stream *list.ListResultsStream, metadata sdk.ResourceMetadata) {
	client := metadata.Client.Postgres.FlexibleServerFirewallRuleClient

	var data PostgresqlFlexibleServerFirewallRuleListModel
	diags := request.Config.Get(ctx, &data)
	if diags.HasError() {
		stream.Results = list.ListResultsStreamDiagnostics(diags)
		return
	}

	parentID, err := firewallrules.ParseFlexibleServerID(data.FlexibleServerId.ValueString())
	if err != nil {
		sdk.SetResponseErrorDiagnostic(stream, fmt.Sprintf("parsing Flexible Server ID for `%s`", azurePostgresqlFlexibleServerFirewallRuleResourceName), err)
		return
	}

	resp, err := client.ListByServerComplete(ctx, *parentID)
	if err != nil {
		sdk.SetResponseErrorDiagnostic(stream, fmt.Sprintf("listing `%s`", azurePostgresqlFlexibleServerFirewallRuleResourceName), err)
		return
	}

	stream.Results = func(push func(list.ListResult) bool) {
		for _, item := range resp.Items {
			result := request.NewListResult(ctx)
			result.DisplayName = pointer.From(item.Name)

			rd := resourcePostgresqlFlexibleServerFirewallRule().Data(&terraform.InstanceState{})
			id, err := firewallrules.ParseFirewallRuleIDInsensitively(pointer.From(item.Id))
			if err != nil {
				sdk.SetErrorDiagnosticAndPushListResult(result, push, fmt.Sprintf("parsing ID for `%s`", azurePostgresqlFlexibleServerFirewallRuleResourceName), err)
				return
			}
			rd.SetId(id.ID())

			if err := resourcePostgresqlFlexibleServerFirewallRuleFlatten(rd, id, &item); err != nil {
				sdk.SetErrorDiagnosticAndPushListResult(result, push, fmt.Sprintf("encoding `%s` resource data", azurePostgresqlFlexibleServerFirewallRuleResourceName), err)
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
