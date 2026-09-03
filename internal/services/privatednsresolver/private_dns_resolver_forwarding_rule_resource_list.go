// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package privatednsresolver

import (
	"context"
	"fmt"

	"github.com/hashicorp/go-azure-helpers/framework/typehelpers"
	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/dnsresolver/2022-07-01/dnsforwardingrulesets"
	"github.com/hashicorp/go-azure-sdk/resource-manager/dnsresolver/2022-07-01/forwardingrules"
	"github.com/hashicorp/terraform-plugin-framework/list"
	"github.com/hashicorp/terraform-plugin-framework/list/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type PrivateDNSResolverForwardingRuleListResource struct{}

type PrivateDNSResolverForwardingRuleListModel struct {
	DnsForwardingRulesetId types.String `tfsdk:"dns_forwarding_ruleset_id"`
}

var _ sdk.FrameworkListWrappedResource = new(PrivateDNSResolverForwardingRuleListResource)

func (PrivateDNSResolverForwardingRuleListResource) Metadata(_ context.Context, _ resource.MetadataRequest, response *resource.MetadataResponse) {
	response.TypeName = PrivateDNSResolverForwardingRuleResource{}.ResourceType()
}

func (PrivateDNSResolverForwardingRuleListResource) ResourceFunc() *pluginsdk.Resource {
	return sdk.WrappedResource(PrivateDNSResolverForwardingRuleResource{})
}

func (PrivateDNSResolverForwardingRuleListResource) ListResourceConfigSchema(_ context.Context, _ list.ListResourceSchemaRequest, response *list.ListResourceSchemaResponse) {
	response.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"dns_forwarding_ruleset_id": schema.StringAttribute{
				Required: true,
				Validators: []validator.String{
					typehelpers.WrappedStringValidator{Func: dnsforwardingrulesets.ValidateDnsForwardingRulesetID},
				},
			},
		},
	}
}

func (PrivateDNSResolverForwardingRuleListResource) List(ctx context.Context, request list.ListRequest, stream *list.ListResultsStream, metadata sdk.ResourceMetadata) {
	client := metadata.Client.PrivateDnsResolver.ForwardingRulesClient

	var data PrivateDNSResolverForwardingRuleListModel
	diags := request.Config.Get(ctx, &data)
	if diags.HasError() {
		stream.Results = list.ListResultsStreamDiagnostics(diags)
		return
	}

	resource := PrivateDNSResolverForwardingRuleResource{}

	parentID, err := dnsforwardingrulesets.ParseDnsForwardingRulesetID(data.DnsForwardingRulesetId.ValueString())
	if err != nil {
		sdk.SetResponseErrorDiagnostic(stream, fmt.Sprintf("parsing parent ID for `%s`", resource.ResourceType()), err)
		return
	}

	rulesetId := forwardingrules.NewDnsForwardingRulesetID(parentID.SubscriptionId, parentID.ResourceGroupName, parentID.DnsForwardingRulesetName)

	resp, err := client.ListComplete(ctx, rulesetId, forwardingrules.DefaultListOperationOptions())
	if err != nil {
		sdk.SetResponseErrorDiagnostic(stream, fmt.Sprintf("listing `%s`", resource.ResourceType()), err)
		return
	}

	stream.Results = func(push func(list.ListResult) bool) {
		for _, item := range resp.Items {
			result := request.NewListResult(ctx)
			result.DisplayName = pointer.From(item.Name)

			id, err := forwardingrules.ParseForwardingRuleIDInsensitively(pointer.From(item.Id))
			if err != nil {
				sdk.SetErrorDiagnosticAndPushListResult(result, push, fmt.Sprintf("parsing %s ID", resource.ResourceType()), err)
				return
			}

			meta := sdk.NewResourceMetaData(metadata.Client, resource)
			meta.SetID(id)

			if err := resource.flatten(meta, id, &item); err != nil {
				sdk.SetErrorDiagnosticAndPushListResult(result, push, fmt.Sprintf("encoding `%s` resource data", resource.ResourceType()), err)
				return
			}

			sdk.EncodeListResult(ctx, meta.ResourceData, &result)
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
