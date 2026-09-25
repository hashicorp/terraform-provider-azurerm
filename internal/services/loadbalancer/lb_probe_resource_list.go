// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package loadbalancer

import (
	"context"
	"fmt"

	"github.com/hashicorp/go-azure-helpers/framework/typehelpers"
	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2025-07-01/loadbalancers"
	"github.com/hashicorp/terraform-plugin-framework/list"
	"github.com/hashicorp/terraform-plugin-framework/list/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type LoadBalancerProbeListResource struct{}

type ArmLoadBalancerProbeListModel struct {
	ProviderLoadBalancerId types.String `tfsdk:"loadbalancer_id"`
}

var _ sdk.FrameworkListWrappedResource = new(LoadBalancerProbeListResource)

func (LoadBalancerProbeListResource) Metadata(_ context.Context, _ resource.MetadataRequest, response *resource.MetadataResponse) {
	response.TypeName = azureLoadBalancerProbeResourceName
}

func (LoadBalancerProbeListResource) ResourceFunc() *pluginsdk.Resource {
	return resourceArmLoadBalancerProbe()
}

func (LoadBalancerProbeListResource) ListResourceConfigSchema(_ context.Context, _ list.ListResourceSchemaRequest, response *list.ListResourceSchemaResponse) {
	response.Schema = schema.Schema{Attributes: map[string]schema.Attribute{
		"loadbalancer_id": schema.StringAttribute{
			Required: true,
			Validators: []validator.String{
				typehelpers.WrappedStringValidator{Func: loadbalancers.ValidateProviderLoadBalancerID},
			},
		},
	}}
}

func (LoadBalancerProbeListResource) List(ctx context.Context, request list.ListRequest, stream *list.ListResultsStream, metadata sdk.ResourceMetadata) {
	client := metadata.Client.LoadBalancers.LoadBalancersClient

	var data ArmLoadBalancerProbeListModel
	if diags := request.Config.Get(ctx, &data); diags.HasError() {
		stream.Results = list.ListResultsStreamDiagnostics(diags)
		return
	}

	parentID, err := loadbalancers.ParseProviderLoadBalancerID(data.ProviderLoadBalancerId.ValueString())
	if err != nil {
		sdk.SetResponseErrorDiagnostic(stream, fmt.Sprintf("parsing Load Balancer ID for `%s`", azureLoadBalancerProbeResourceName), err)
		return
	}

	resp, err := client.LoadBalancerProbesListComplete(ctx, *parentID)
	if err != nil {
		sdk.SetResponseErrorDiagnostic(stream, fmt.Sprintf("listing `%s`", azureLoadBalancerProbeResourceName), err)
		return
	}

	stream.Results = func(push func(list.ListResult) bool) {
		for _, item := range resp.Items {
			result := request.NewListResult(ctx)
			result.DisplayName = pointer.From(item.Name)

			rd := resourceArmLoadBalancerProbe().Data(&terraform.InstanceState{})

			id, err := loadbalancers.ParseProbeIDInsensitively(pointer.From(item.Id))
			if err != nil {
				sdk.SetErrorDiagnosticAndPushListResult(result, push, fmt.Sprintf("parsing ID for `%s`", azureLoadBalancerProbeResourceName), err)
				return
			}
			rd.SetId(id.ID())

			if err := resourceArmLoadBalancerProbeFlatten(rd, id, &item); err != nil {
				sdk.SetErrorDiagnosticAndPushListResult(result, push, fmt.Sprintf("encoding `%s` resource data", azureLoadBalancerProbeResourceName), err)
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
