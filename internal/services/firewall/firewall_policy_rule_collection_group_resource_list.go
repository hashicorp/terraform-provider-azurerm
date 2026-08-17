// Copyright IBM Corp.
// SPDX-License-Identifier: MPL-2.0

package firewall

import (
	"context"
	"fmt"

	"github.com/hashicorp/go-azure-helpers/framework/typehelpers"
	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2025-01-01/firewallpolicyrulecollectiongroups"
	"github.com/hashicorp/terraform-plugin-framework/list"
	"github.com/hashicorp/terraform-plugin-framework/list/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type (
	FirewallPolicyRuleCollectionGroupListResource struct{}
	FirewallPolicyRuleCollectionGroupListModel    struct {
		FirewallPolicyId types.String `tfsdk:"firewall_policy_id"`
	}
)

var _ sdk.FrameworkListWrappedResource = new(FirewallPolicyRuleCollectionGroupListResource)

func (r FirewallPolicyRuleCollectionGroupListResource) ResourceFunc() *pluginsdk.Resource {
	return resourceFirewallPolicyRuleCollectionGroup()
}

func (r FirewallPolicyRuleCollectionGroupListResource) Metadata(_ context.Context, _ resource.MetadataRequest, response *resource.MetadataResponse) {
	response.TypeName = "azurerm_firewall_policy_rule_collection_group"
}

func (r FirewallPolicyRuleCollectionGroupListResource) ListResourceConfigSchema(_ context.Context, _ list.ListResourceSchemaRequest, response *list.ListResourceSchemaResponse) {
	response.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"firewall_policy_id": schema.StringAttribute{
				Required: true,
				Validators: []validator.String{
					typehelpers.WrappedStringValidator{
						Func: firewallpolicyrulecollectiongroups.ValidateFirewallPolicyID,
					},
				},
			},
		},
	}
}

func (r FirewallPolicyRuleCollectionGroupListResource) List(ctx context.Context, request list.ListRequest, stream *list.ListResultsStream, metadata sdk.ResourceMetadata) {
	client := metadata.Client.Network.FirewallPolicyRuleCollectionGroups
	var data FirewallPolicyRuleCollectionGroupListModel
	diags := request.Config.Get(ctx, &data)
	if diags.HasError() {
		stream.Results = list.ListResultsStreamDiagnostics(diags)
		return
	}

	results := make([]firewallpolicyrulecollectiongroups.FirewallPolicyRuleCollectionGroup, 0)

	if !data.FirewallPolicyId.IsNull() {
		firewallPolicyId, err := firewallpolicyrulecollectiongroups.ParseFirewallPolicyID(data.FirewallPolicyId.ValueString())
		if err != nil {
			sdk.SetResponseErrorDiagnostic(stream, fmt.Sprintf("parsing firewall Firewall ID for `%s`", "azurerm_firewall_policy_rule_collection_group"), err)
			return
		}

		resp, err := client.ListComplete(ctx, *firewallPolicyId)
		if err != nil {
			sdk.SetResponseErrorDiagnostic(stream, fmt.Sprintf("listing `%s`", "azurerm_firewall_policy_rule_collection_group"), err)
			return
		}

		results = resp.Items
	}

	stream.Results = func(push func(list.ListResult) bool) {
		for _, policyrulecollectiongroup := range results {
			result := request.NewListResult(ctx)

			result.DisplayName = pointer.From(policyrulecollectiongroup.Name)

			rd := resourceFirewallPolicyRuleCollectionGroup().Data(&terraform.InstanceState{})

			id, err := firewallpolicyrulecollectiongroups.ParseRuleCollectionGroupID(pointer.From(policyrulecollectiongroup.Id))
			if err != nil {
				sdk.SetErrorDiagnosticAndPushListResult(result, push, "parsing Firewall PolicyRuleCollectionGroup ID", err)
				return
			}

			rd.SetId(id.ID())

			if err := resourceFirewallPolicyRuleCollectionGroupSetFlatten(rd, id, &policyrulecollectiongroup); err != nil {
				sdk.SetErrorDiagnosticAndPushListResult(result, push, fmt.Sprintf("encoding `%s` resource data", "azurerm_firewall_policy_rule_collection_group"), err)
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
