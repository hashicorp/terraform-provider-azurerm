// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package loadbalancer

import (
	"context"
	"fmt"

	"github.com/hashicorp/go-azure-helpers/framework/typehelpers"
	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2025-01-01/loadbalancers"
	"github.com/hashicorp/terraform-plugin-framework/list"
	"github.com/hashicorp/terraform-plugin-framework/list/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

const azurermLbRuleResourceName = "azurerm_lb_rule"

type (
	LbRuleListResource struct{}
	LbRuleListModel    struct {
		LoadBalancerId types.String `tfsdk:"loadbalancer_id"`
	}
)

var _ sdk.FrameworkListWrappedResource = new(LbRuleListResource)

func (LbRuleListResource) Metadata(_ context.Context, _ resource.MetadataRequest, response *resource.MetadataResponse) {
	response.TypeName = azurermLbRuleResourceName
}

func (LbRuleListResource) ResourceFunc() *pluginsdk.Resource {
	return resourceArmLoadBalancerRule()
}

func (LbRuleListResource) ListResourceConfigSchema(_ context.Context, _ list.ListResourceSchemaRequest, response *list.ListResourceSchemaResponse) {
	response.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"loadbalancer_id": schema.StringAttribute{
				Required: true,
				Validators: []validator.String{
					typehelpers.WrappedStringValidator{Func: loadbalancers.ValidateLoadBalancerID},
				},
			},
		},
	}
}

func (LbRuleListResource) List(ctx context.Context, request list.ListRequest, stream *list.ListResultsStream, metadata sdk.ResourceMetadata) {
	client := metadata.Client.LoadBalancers.LoadBalancersClient

	var data LbRuleListModel
	if diags := request.Config.Get(ctx, &data); diags.HasError() {
		stream.Results = list.ListResultsStreamDiagnostics(diags)
		return
	}

	lbId, err := loadbalancers.ParseLoadBalancerID(data.LoadBalancerId.ValueString())
	if err != nil {
		sdk.SetResponseErrorDiagnostic(stream, fmt.Sprintf("parsing parent ID for `%s`", azurermLbRuleResourceName), err)
		return
	}

	plbId := loadbalancers.ProviderLoadBalancerId{
		SubscriptionId:    lbId.SubscriptionId,
		ResourceGroupName: lbId.ResourceGroupName,
		LoadBalancerName:  lbId.LoadBalancerName,
	}

	lbResp, err := client.Get(ctx, plbId, loadbalancers.GetOperationOptions{})
	if err != nil {
		sdk.SetResponseErrorDiagnostic(stream, fmt.Sprintf("retrieving Load Balancer for `%s`", azurermLbRuleResourceName), err)
		return
	}

	resp, err := client.LoadBalancerLoadBalancingRulesListComplete(ctx, plbId)
	if err != nil {
		sdk.SetResponseErrorDiagnostic(stream, fmt.Sprintf("listing `%s`", azurermLbRuleResourceName), err)
		return
	}

	stream.Results = func(push func(list.ListResult) bool) {
		for _, item := range resp.Items {
			result := request.NewListResult(ctx)
			result.DisplayName = pointer.From(item.Name)

			rd := resourceArmLoadBalancerRule().Data(&terraform.InstanceState{})

			id, err := loadbalancers.ParseLoadBalancingRuleIDInsensitively(pointer.From(item.Id))
			if err != nil {
				sdk.SetErrorDiagnosticAndPushListResult(result, push, fmt.Sprintf("parsing `%s` ID", azurermLbRuleResourceName), err)
				return
			}
			rd.SetId(id.ID())

			if err := resourceArmLoadBalancerRuleFlatten(rd, id, lbResp.Model); err != nil {
				sdk.SetErrorDiagnosticAndPushListResult(result, push, fmt.Sprintf("encoding `%s` resource data", azurermLbRuleResourceName), err)
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
