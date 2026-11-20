// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package palo_alto_next_generation_firewall_metrics

import (
	"context"
	"fmt"

	"github.com/hashicorp/go-azure-helpers/framework/typehelpers"
	"github.com/hashicorp/go-azure-sdk/resource-manager/paloaltonetworks/2025-10-08/metricsobjectfirewallresources"
	"github.com/hashicorp/terraform-plugin-framework/list"
	"github.com/hashicorp/terraform-plugin-framework/list/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type ListResource struct{}

type ListModel struct {
	FirewallId types.String `tfsdk:"firewall_id"`
}

var _ sdk.FrameworkListWrappedResourceWithConfig = new(ListResource)

func (ListResource) ResourceFunc() *pluginsdk.Resource {
	return sdk.WrappedResource(Resource{})
}

func (ListResource) Metadata(_ context.Context, _ resource.MetadataRequest, response *resource.MetadataResponse) {
	response.TypeName = Resource{}.ResourceType()
}

func (ListResource) ListResourceConfigSchema(_ context.Context, _ list.ListResourceSchemaRequest, response *list.ListResourceSchemaResponse) {
	response.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"firewall_id": schema.StringAttribute{
				Required: true,
				Validators: []validator.String{
					typehelpers.WrappedStringValidator{
						Func: metricsobjectfirewallresources.ValidateFirewallID,
					},
				},
			},
		},
	}
}

func (ListResource) List(ctx context.Context, request list.ListRequest, stream *list.ListResultsStream, metadata sdk.ResourceMetadata) {
	client := metadata.Client.PaloAlto.MetricsObjectFirewallResources

	var data ListModel
	diags := request.Config.Get(ctx, &data)
	if diags.HasError() {
		stream.Results = list.ListResultsStreamDiagnostics(diags)
		return
	}

	r := Resource{}

	firewallId, err := metricsobjectfirewallresources.ParseFirewallID(data.FirewallId.ValueString())
	if err != nil {
		sdk.SetResponseErrorDiagnostic(stream, "parsing Palo Alto Next Generation Firewall ID", err)
		return
	}

	resp, err := client.MetricsObjectFirewallListByFirewallsComplete(ctx, *firewallId)
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
