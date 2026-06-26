// Copyright IBM Corp.
// SPDX-License-Identifier: MPL-2.0

package network

import (
	"context"
	"fmt"

	"github.com/hashicorp/go-azure-helpers/framework/typehelpers"
	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2025-01-01/securityrules"
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
	NetworkSecurityRuleListResource struct{}
	NetworkSecurityRuleListModel    struct {
		NetworkSecurityGroupId types.String `tfsdk:"network_security_group_id"`
	}
)

var _ sdk.FrameworkListWrappedResource = new(NetworkSecurityRuleListResource)

func (r NetworkSecurityRuleListResource) ResourceFunc() *pluginsdk.Resource {
	return resourceNetworkSecurityRule()
}

func (r NetworkSecurityRuleListResource) Metadata(_ context.Context, _ resource.MetadataRequest, response *resource.MetadataResponse) {
	response.TypeName = "azurerm_network_security_rule"
}

func (r NetworkSecurityRuleListResource) ListResourceConfigSchema(_ context.Context, _ list.ListResourceSchemaRequest, response *list.ListResourceSchemaResponse) {
	response.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"network_security_group_id": schema.StringAttribute{
				Required: true,
				Validators: []validator.String{
					typehelpers.WrappedStringValidator{
						Func: securityrules.ValidateNetworkSecurityGroupID,
					},
				},
			},
		},
	}
}

func (r NetworkSecurityRuleListResource) List(ctx context.Context, request list.ListRequest, stream *list.ListResultsStream, metadata sdk.ResourceMetadata) {
	client := metadata.Client.Network.SecurityRules
	var data NetworkSecurityRuleListModel
	diags := request.Config.Get(ctx, &data)
	if diags.HasError() {
		stream.Results = list.ListResultsStreamDiagnostics(diags)
		return
	}

	results := make([]securityrules.SecurityRule, 0)

	if !data.NetworkSecurityGroupId.IsNull() {
		networksecuritygroupId, err := securityrules.ParseNetworkSecurityGroupID(data.NetworkSecurityGroupId.ValueString())
		if err != nil {
			sdk.SetResponseErrorDiagnostic(stream, fmt.Sprintf("parsing Network Security Group ID for `%s`", "azurerm_network_security_rule"), err)
			return
		}

		resp, err := client.ListComplete(ctx, *networksecuritygroupId)
		if err != nil {
			sdk.SetResponseErrorDiagnostic(stream, fmt.Sprintf("listing `%s`", "azurerm_network_security_rule"), err)
			return
		}

		results = resp.Items
	}

	stream.Results = func(push func(list.ListResult) bool) {
		for _, securityrule := range results {
			result := request.NewListResult(ctx)

			result.DisplayName = pointer.From(securityrule.Name)

			rd := resourceNetworkSecurityRule().Data(&terraform.InstanceState{})

			id, err := securityrules.ParseSecurityRuleID(pointer.From(securityrule.Id))
			if err != nil {
				sdk.SetErrorDiagnosticAndPushListResult(result, push, "parsing Network Security Rule ID", err)
				return
			}

			rd.SetId(id.ID())

			if err := resourceNetworkSecurityRuleFlatten(rd, id, &securityrule); err != nil {
				sdk.SetErrorDiagnosticAndPushListResult(result, push, fmt.Sprintf("encoding `%s` resource data", "azurerm_network_security_rule"), err)
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
