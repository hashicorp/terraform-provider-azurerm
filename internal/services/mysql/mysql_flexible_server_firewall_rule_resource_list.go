package mysql

import (
	"context"
	"fmt"

	"github.com/hashicorp/go-azure-helpers/framework/typehelpers"
	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/mysql/2023-12-30/firewallrules"
	"github.com/hashicorp/terraform-plugin-framework/list"
	"github.com/hashicorp/terraform-plugin-framework/list/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type MysqlFlexibleServerFirewallRuleListResource struct{}

type MysqlFlexibleServerFirewallRuleListModel struct {
	FlexibleServerId types.String `tfsdk:"flexible_server_id"`
}

var _ sdk.FrameworkListWrappedResource = new(MysqlFlexibleServerFirewallRuleListResource)

func (r MysqlFlexibleServerFirewallRuleListResource) ResourceFunc() *pluginsdk.Resource {
	return resourceMySqlFlexibleServerFirewallRule()
}

func (r MysqlFlexibleServerFirewallRuleListResource) Metadata(_ context.Context, _ resource.MetadataRequest, response *resource.MetadataResponse) {
	response.TypeName = mysqlFlexibleServerFirewallResourceName
}

func (r MysqlFlexibleServerFirewallRuleListResource) ListResourceConfigSchema(_ context.Context, _ list.ListResourceSchemaRequest, response *list.ListResourceSchemaResponse) {
	response.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"flexible_server_id": schema.StringAttribute{
				Required: true,
				Validators: []validator.String{
					typehelpers.WrappedStringValidator{
						Func: firewallrules.ValidateFlexibleServerID,
					},
				},
			},
		},
	}
}

func (r MysqlFlexibleServerFirewallRuleListResource) List(ctx context.Context, request list.ListRequest, stream *list.ListResultsStream, metadata sdk.ResourceMetadata) {
	client := metadata.Client.MySQL.FlexibleServers.FirewallRules

	var data MysqlFlexibleServerFirewallRuleListModel
	diags := request.Config.Get(ctx, &data)
	if diags.HasError() {
		stream.Results = list.ListResultsStreamDiagnostics(diags)
		return
	}

	results := make([]firewallrules.FirewallRule, 0)

	serverId, err := firewallrules.ParseFlexibleServerID(data.FlexibleServerId.ValueString())
	if err != nil {
		sdk.SetResponseErrorDiagnostic(stream, fmt.Sprintf("parsing Mysql Server ID for `%s`", mysqlFlexibleServerFirewallResourceName), err)
		return
	}

	resp, err := client.ListByServerComplete(ctx, pointer.From(serverId))
	if err != nil {
		sdk.SetResponseErrorDiagnostic(stream, fmt.Sprintf("listing `%s`", mysqlFlexibleServerFirewallResourceName), err)
		return
	}
	results = resp.Items

	stream.Results = func(push func(list.ListResult) bool) {
		for _, rule := range results {
			result := request.NewListResult(ctx)
			result.DisplayName = pointer.From(rule.Name)

			id, err := firewallrules.ParseFirewallRuleID(pointer.From(rule.Id))
			if err != nil {
				sdk.SetErrorDiagnosticAndPushListResult(result, push, "parsing Mysql Firewall Rule ID", err)
				return
			}

			rd := resourceMySqlFlexibleServerFirewallRule().Data(&terraform.InstanceState{})
			rd.SetId(id.ID())

			if err := resourceMySqlFlexibleServerFirewallRuleFlatten(rd, id, &rule); err != nil {
				sdk.SetErrorDiagnosticAndPushListResult(result, push, fmt.Sprintf("encoding `%s` resource data", mysqlFlexibleServerFirewallResourceName), err)
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
