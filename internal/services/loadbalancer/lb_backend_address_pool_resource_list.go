// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package loadbalancer

import (
	"context"

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

type ArmLoadBalancerBackendAddressPoolListResource struct{}

type ArmLoadBalancerBackendAddressPoolListModel struct {
	LoadBalancerId types.String `tfsdk:"loadbalancer_id"`
}

var _ sdk.FrameworkListWrappedResource = new(ArmLoadBalancerBackendAddressPoolListResource)

func (ArmLoadBalancerBackendAddressPoolListResource) Metadata(_ context.Context, _ resource.MetadataRequest, response *resource.MetadataResponse) {
	response.TypeName = backendAddressPoolResourceName
}

func (ArmLoadBalancerBackendAddressPoolListResource) ResourceFunc() *pluginsdk.Resource {
	return resourceArmLoadBalancerBackendAddressPool()
}

func (ArmLoadBalancerBackendAddressPoolListResource) ListResourceConfigSchema(_ context.Context, _ list.ListResourceSchemaRequest, response *list.ListResourceSchemaResponse) {
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

func (ArmLoadBalancerBackendAddressPoolListResource) List(ctx context.Context, request list.ListRequest, stream *list.ListResultsStream, metadata sdk.ResourceMetadata) {
	client := metadata.Client.LoadBalancers.LoadBalancersClient

	var data ArmLoadBalancerBackendAddressPoolListModel
	if diags := request.Config.Get(ctx, &data); diags.HasError() {
		stream.Results = list.ListResultsStreamDiagnostics(diags)
		return
	}

	parentID, err := loadbalancers.ParseLoadBalancerID(data.LoadBalancerId.ValueString())
	if err != nil {
		sdk.SetResponseErrorDiagnostic(stream, "parsing parent ID for "+backendAddressPoolResourceName, err)
		return
	}

	plbId := loadbalancers.NewProviderLoadBalancerID(parentID.SubscriptionId, parentID.ResourceGroupName, parentID.LoadBalancerName)

	resp, err := client.LoadBalancerBackendAddressPoolsListComplete(ctx, plbId)
	if err != nil {
		sdk.SetResponseErrorDiagnostic(stream, "listing "+backendAddressPoolResourceName, err)
		return
	}

	stream.Results = func(push func(list.ListResult) bool) {
		for _, item := range resp.Items {
			result := request.NewListResult(ctx)
			result.DisplayName = pointer.From(item.Name)

			rd := resourceArmLoadBalancerBackendAddressPool().Data(&terraform.InstanceState{})

			id, err := loadbalancers.ParseLoadBalancerBackendAddressPoolIDInsensitively(pointer.From(item.Id))
			if err != nil {
				sdk.SetErrorDiagnosticAndPushListResult(result, push, "parsing "+backendAddressPoolResourceName+" ID", err)
				return
			}
			rd.SetId(id.ID())

			if err := resourceArmLoadBalancerBackendAddressPoolFlatten(rd, id, &item); err != nil {
				sdk.SetErrorDiagnosticAndPushListResult(result, push, "encoding "+backendAddressPoolResourceName+" resource data", err)
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
