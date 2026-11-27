// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package horizondb

import (
	"context"
	"fmt"

	"github.com/hashicorp/go-azure-helpers/framework/typehelpers"
	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/horizondb/2026-10-01-preview/firewallrules"
	"github.com/hashicorp/terraform-plugin-framework/list"
	"github.com/hashicorp/terraform-plugin-framework/list/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/horizondb/client"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type HorizondbFirewallRuleListResource struct{}

type HorizondbFirewallRuleListModel struct {
	HorizondbPoolId types.String `tfsdk:"horizondb_pool_id"`
}

var _ sdk.FrameworkListWrappedResource = HorizondbFirewallRuleListResource{}

func (HorizondbFirewallRuleListResource) ResourceFunc() *pluginsdk.Resource {
	return sdk.WrappedResource(HorizondbFirewallRuleResource{})
}

func (HorizondbFirewallRuleListResource) Metadata(_ context.Context, _ resource.MetadataRequest, response *resource.MetadataResponse) {
	response.TypeName = HorizondbFirewallRuleResource{}.ResourceType()
}

func (HorizondbFirewallRuleListResource) ListResourceConfigSchema(_ context.Context, _ list.ListResourceSchemaRequest, response *list.ListResourceSchemaResponse) {
	response.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"horizondb_pool_id": schema.StringAttribute{
				Required: true,
				Validators: []validator.String{
					typehelpers.WrappedStringValidator{Func: firewallrules.ValidatePoolID},
				},
			},
		},
	}
}

func (r HorizondbFirewallRuleListResource) List(ctx context.Context, request list.ListRequest, stream *list.ListResultsStream, metadata sdk.ResourceMetadata) {
	var config HorizondbFirewallRuleListModel

	if diags := request.Config.Get(ctx, &config); diags.HasError() {
		stream.Results = list.ListResultsStreamDiagnostics(diags)
		return
	}

	poolId, err := firewallrules.ParsePoolID(config.HorizondbPoolId.ValueString())
	if err != nil {
		sdk.SetResponseErrorDiagnostic(stream, "parsing HorizonDB Pool ID", err)
		return
	}

	resp, err := client.From(metadata.Client).FirewallRulesClient.ListComplete(ctx, *poolId)
	if err != nil {
		sdk.SetResponseErrorDiagnostic(stream, "listing HorizonDB Firewall Rules", err)
		return
	}

	stream.Results = func(push func(list.ListResult) bool) {
		for _, item := range resp.Items {
			result := request.NewListResult(ctx)

			id, err := firewallrules.ParseFirewallRuleIDInsensitively(pointer.From(item.Id))
			if err != nil {
				sdk.SetErrorDiagnosticAndPushListResult(result, push, "parsing HorizonDB Firewall Rule ID", err)
				return
			}

			result.DisplayName = id.FirewallRuleName
			resourceMetadata := sdk.NewResourceMetaData(metadata.Client, HorizondbFirewallRuleResource{})
			resourceMetadata.ResourceData.SetId(id.ID())

			if err := (HorizondbFirewallRuleResource{}).flatten(resourceMetadata, id, &item); err != nil {
				sdk.SetErrorDiagnosticAndPushListResult(result, push, fmt.Sprintf("encoding `%s` resource data", HorizondbFirewallRuleResource{}.ResourceType()), err)
				return
			}

			sdk.EncodeListResult(ctx, resourceMetadata.ResourceData, &result)

			if !push(result) || result.Diagnostics.HasError() {
				return
			}
		}
	}
}
