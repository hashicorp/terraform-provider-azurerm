// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package privatednsresolver

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-sdk/resource-manager/dnsresolver/2022-07-01/dnsforwardingrulesets"
	"github.com/hashicorp/go-azure-sdk/resource-manager/dnsresolver/2022-07-01/outboundendpoints"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

type PrivateDNSResolverDnsForwardingRulesetModel struct {
	Name                         string            `tfschema:"name"`
	ResourceGroupName            string            `tfschema:"resource_group_name"`
	DnsResolverOutboundEndpoints []string          `tfschema:"private_dns_resolver_outbound_endpoint_ids"`
	Location                     string            `tfschema:"location"`
	Tags                         map[string]string `tfschema:"tags"`
}

type PrivateDNSResolverDnsForwardingRulesetResource struct{}

var _ sdk.ResourceWithUpdate = PrivateDNSResolverDnsForwardingRulesetResource{}

func (r PrivateDNSResolverDnsForwardingRulesetResource) ResourceType() string {
	return "azurerm_private_dns_resolver_dns_forwarding_ruleset"
}

func (r PrivateDNSResolverDnsForwardingRulesetResource) ModelObject() interface{} {
	return &PrivateDNSResolverDnsForwardingRulesetModel{}
}

func (r PrivateDNSResolverDnsForwardingRulesetResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return dnsforwardingrulesets.ValidateDnsForwardingRulesetID
}

func (r PrivateDNSResolverDnsForwardingRulesetResource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"resource_group_name": commonschema.ResourceGroupName(),

		"private_dns_resolver_outbound_endpoint_ids": {
			Type:     pluginsdk.TypeList,
			Required: true,
			Elem: &pluginsdk.Schema{
				Type:         pluginsdk.TypeString,
				ValidateFunc: outboundendpoints.ValidateOutboundEndpointID,
			},
		},

		"location": commonschema.Location(),

		"tags": commonschema.Tags(),
	}
}

func (r PrivateDNSResolverDnsForwardingRulesetResource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{}
}

func (r PrivateDNSResolverDnsForwardingRulesetResource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			var model PrivateDNSResolverDnsForwardingRulesetModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			client := metadata.Client.PrivateDnsResolver.DnsForwardingRulesetsClient
			subscriptionId := metadata.Client.Account.SubscriptionId
			id := dnsforwardingrulesets.NewDnsForwardingRulesetID(subscriptionId, model.ResourceGroupName, model.Name)
			existing, err := client.Get(ctx, id)
			if err != nil && !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for existing %s: %+v", id, err)
			}

			if !response.WasNotFound(existing.HttpResponse) {
				return metadata.ResourceRequiresImport(r.ResourceType(), id)
			}

			properties := &dnsforwardingrulesets.DnsForwardingRuleset{
				Location:   location.Normalize(model.Location),
				Properties: dnsforwardingrulesets.DnsForwardingRulesetProperties{},
				Tags:       &model.Tags,
			}

			dnsResolverOutboundEndpointsValue := expandDnsResolverOutboundEndpoints(model.DnsResolverOutboundEndpoints)

			if dnsResolverOutboundEndpointsValue != nil {
				properties.Properties.DnsResolverOutboundEndpoints = *dnsResolverOutboundEndpointsValue
			}

			if err := client.CreateOrUpdateThenPoll(ctx, id, *properties, dnsforwardingrulesets.CreateOrUpdateOperationOptions{}); err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}

			metadata.SetID(id)
			return nil
		},
	}
}

func (r PrivateDNSResolverDnsForwardingRulesetResource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.PrivateDnsResolver.DnsForwardingRulesetsClient

			id, err := dnsforwardingrulesets.ParseDnsForwardingRulesetID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			var model PrivateDNSResolverDnsForwardingRulesetModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			resp, err := client.Get(ctx, *id)
			if err != nil {
				return fmt.Errorf("retrieving %s: %+v", *id, err)
			}

			properties := resp.Model
			if properties == nil {
				return fmt.Errorf("retrieving %s: properties was nil", id)
			}

			if metadata.ResourceData.HasChange("dns_resolver_outbound_endpoints") {
				dnsResolverOutboundEndpointsValue := expandDnsResolverOutboundEndpoints(model.DnsResolverOutboundEndpoints)

				if dnsResolverOutboundEndpointsValue != nil {
					properties.Properties.DnsResolverOutboundEndpoints = *dnsResolverOutboundEndpointsValue
				}
			}

			if metadata.ResourceData.HasChange("tags") {
				properties.Tags = &model.Tags
			}

			if err := client.CreateOrUpdateThenPoll(ctx, *id, *properties, dnsforwardingrulesets.CreateOrUpdateOperationOptions{}); err != nil {
				return fmt.Errorf("updating %s: %+v", *id, err)
			}

			return nil
		},
	}
}

func (r PrivateDNSResolverDnsForwardingRulesetResource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.PrivateDnsResolver.DnsForwardingRulesetsClient

			id, err := dnsforwardingrulesets.ParseDnsForwardingRulesetID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			resp, err := client.Get(ctx, *id)
			if err != nil {
				if response.WasNotFound(resp.HttpResponse) {
					return metadata.MarkAsGone(id)
				}

				return fmt.Errorf("retrieving %s: %+v", *id, err)
			}

			model := resp.Model
			if model == nil {
				return fmt.Errorf("retrieving %s: model was nil", id)
			}

			state := PrivateDNSResolverDnsForwardingRulesetModel{
				Name:              id.DnsForwardingRulesetName,
				ResourceGroupName: id.ResourceGroupName,
				Location:          location.Normalize(model.Location),
			}

			properties := &model.Properties

			state.DnsResolverOutboundEndpoints = flattenDnsResolverOutboundEndpoints(&properties.DnsResolverOutboundEndpoints)

			if model.Tags != nil {
				state.Tags = *model.Tags
			}

			return metadata.Encode(&state)
		},
	}
}

func (r PrivateDNSResolverDnsForwardingRulesetResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.PrivateDnsResolver.DnsForwardingRulesetsClient

			id, err := dnsforwardingrulesets.ParseDnsForwardingRulesetID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			if err := client.DeleteThenPoll(ctx, *id, dnsforwardingrulesets.DeleteOperationOptions{}); err != nil {
				return fmt.Errorf("deleting %s: %+v", id, err)
			}

			return nil
		},
	}
}

func expandDnsResolverOutboundEndpoints(inputList []string) *[]dnsforwardingrulesets.SubResource {
	var outputList []dnsforwardingrulesets.SubResource
	for _, v := range inputList {
		output := dnsforwardingrulesets.SubResource{
			Id: v,
		}
		outputList = append(outputList, output)
	}

	return &outputList
}

func flattenDnsResolverOutboundEndpoints(inputList *[]dnsforwardingrulesets.SubResource) []string {
	var outputList []string
	if inputList == nil {
		return outputList
	}

	for _, input := range *inputList {
		output := input.Id
		outputList = append(outputList, output)
	}

	return outputList
}
