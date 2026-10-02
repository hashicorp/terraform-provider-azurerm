// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package paloalto

import (
	"context"
	"fmt"

	"github.com/hashicorp/go-azure-helpers/framework/typehelpers"
	"github.com/hashicorp/go-azure-sdk/resource-manager/paloaltonetworks/2025-10-08/firewallresources"
	"github.com/hashicorp/go-azure-sdk/resource-manager/paloaltonetworks/2025-10-08/metricsobjectfirewallresources"
	"github.com/hashicorp/terraform-plugin-framework/list"
	"github.com/hashicorp/terraform-plugin-framework/list/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type NextGenerationFirewallMetricsListResource struct{}

type NextGenerationFirewallMetricsListModel struct {
	FirewallId types.String `tfsdk:"firewall_id"`
}

var _ sdk.FrameworkListWrappedResourceWithConfig = new(NextGenerationFirewallMetricsListResource)

func (NextGenerationFirewallMetricsListResource) ResourceFunc() *pluginsdk.Resource {
	return sdk.WrappedResource(NextGenerationFirewallMetricsResource{})
}

func (NextGenerationFirewallMetricsListResource) Metadata(_ context.Context, _ resource.MetadataRequest, response *resource.MetadataResponse) {
	response.TypeName = NextGenerationFirewallMetricsResource{}.ResourceType()
}

func (NextGenerationFirewallMetricsListResource) ListResourceConfigSchema(_ context.Context, _ list.ListResourceSchemaRequest, response *list.ListResourceSchemaResponse) {
	response.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"firewall_id": schema.StringAttribute{
				Required: true,
				Validators: []validator.String{
					typehelpers.WrappedStringValidator{
						Func: firewallresources.ValidateFirewallID,
					},
				},
			},
		},
	}
}

func (NextGenerationFirewallMetricsListResource) List(ctx context.Context, request list.ListRequest, stream *list.ListResultsStream, metadata sdk.ResourceMetadata) {
	client := metadata.Client.PaloAlto.MetricsObjectFirewallResources

	var data NextGenerationFirewallMetricsListModel
	diags := request.Config.Get(ctx, &data)
	if diags.HasError() {
		stream.Results = list.ListResultsStreamDiagnostics(diags)
		return
	}

	r := NextGenerationFirewallMetricsResource{}

	firewallId, err := firewallresources.ParseFirewallID(data.FirewallId.ValueString())
	if err != nil {
		sdk.SetResponseErrorDiagnostic(stream, "parsing Palo Alto Next Generation Firewall ID", err)
		return
	}

	metricsFirewallId := metricsobjectfirewallresources.NewFirewallID(firewallId.SubscriptionId, firewallId.ResourceGroupName, firewallId.FirewallName)

	resp, err := client.MetricsObjectFirewallListByFirewallsComplete(ctx, metricsFirewallId)
	if err != nil {
		sdk.SetResponseErrorDiagnostic(stream, fmt.Sprintf("listing `%s`", r.ResourceType()), err)
		return
	}

	stream.Results = func(push func(list.ListResult) bool) {
		for _, item := range resp.Items {
			result := request.NewListResult(ctx)
			// the Metrics Object is a singleton beneath the Firewall, so the Firewall's name and ID are used to identify it
			result.DisplayName = firewallId.FirewallName

			rmd := sdk.NewResourceMetaData(metadata.Client, r)
			rmd.SetID(firewallId)

			if err := r.flatten(rmd, firewallId, &item); err != nil {
				sdk.SetErrorDiagnosticAndPushListResult(result, push, fmt.Sprintf("encoding `%s` resource data", r.ResourceType()), err)
				return
			}

			sdk.EncodeListResult(ctx, rmd.ResourceData, &result)
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
