// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package privatednsresolver

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/dnsresolver/2022-07-01/dnsforwardingrulesets"
	"github.com/hashicorp/go-azure-sdk/resource-manager/dnsresolver/2022-07-01/virtualnetworklinks"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

type PrivateDNSResolverVirtualNetworkLinkModel struct {
	Name                   string            `tfschema:"name"`
	DnsForwardingRulesetId string            `tfschema:"dns_forwarding_ruleset_id"`
	Metadata               map[string]string `tfschema:"metadata"`
	VirtualNetworkId       string            `tfschema:"virtual_network_id"`
}

type PrivateDNSResolverVirtualNetworkLinkResource struct{}

var _ sdk.ResourceWithUpdate = PrivateDNSResolverVirtualNetworkLinkResource{}

func (r PrivateDNSResolverVirtualNetworkLinkResource) ResourceType() string {
	return "azurerm_private_dns_resolver_virtual_network_link"
}

func (r PrivateDNSResolverVirtualNetworkLinkResource) ModelObject() interface{} {
	return &PrivateDNSResolverVirtualNetworkLinkModel{}
}

func (r PrivateDNSResolverVirtualNetworkLinkResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return virtualnetworklinks.ValidateVirtualNetworkLinkID
}

func (r PrivateDNSResolverVirtualNetworkLinkResource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"dns_forwarding_ruleset_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: dnsforwardingrulesets.ValidateDnsForwardingRulesetID,
		},

		"virtual_network_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: commonids.ValidateVirtualNetworkID,
		},

		"metadata": {
			Type:     pluginsdk.TypeMap,
			Optional: true,
			Elem: &pluginsdk.Schema{
				Type: pluginsdk.TypeString,
			},
		},
	}
}

func (r PrivateDNSResolverVirtualNetworkLinkResource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{}
}

func (r PrivateDNSResolverVirtualNetworkLinkResource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			var model PrivateDNSResolverVirtualNetworkLinkModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			client := metadata.Client.PrivateDnsResolver.VirtualNetworkLinksClient
			dnsForwardingRulesetId, err := dnsforwardingrulesets.ParseDnsForwardingRulesetID(model.DnsForwardingRulesetId)
			if err != nil {
				return err
			}

			id := virtualnetworklinks.NewVirtualNetworkLinkID(dnsForwardingRulesetId.SubscriptionId, dnsForwardingRulesetId.ResourceGroupName, dnsForwardingRulesetId.DnsForwardingRulesetName, model.Name)
			existing, err := client.Get(ctx, id)
			if err != nil && !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for existing %s: %+v", id, err)
			}

			if !response.WasNotFound(existing.HttpResponse) {
				return metadata.ResourceRequiresImport(r.ResourceType(), id)
			}

			properties := &virtualnetworklinks.VirtualNetworkLink{
				Properties: virtualnetworklinks.VirtualNetworkLinkProperties{
					Metadata: &model.Metadata,
					VirtualNetwork: virtualnetworklinks.SubResource{
						Id: model.VirtualNetworkId,
					},
				},
			}

			if err := client.CreateOrUpdateThenPoll(ctx, id, *properties, virtualnetworklinks.CreateOrUpdateOperationOptions{}); err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}

			metadata.SetID(id)
			return nil
		},
	}
}

func (r PrivateDNSResolverVirtualNetworkLinkResource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.PrivateDnsResolver.VirtualNetworkLinksClient

			id, err := virtualnetworklinks.ParseVirtualNetworkLinkID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			var model PrivateDNSResolverVirtualNetworkLinkModel
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

			if metadata.ResourceData.HasChange("metadata") {
				properties.Properties.Metadata = &model.Metadata
			}

			if metadata.ResourceData.HasChange("virtual_network_id") {
				properties.Properties.VirtualNetwork = virtualnetworklinks.SubResource{
					Id: model.VirtualNetworkId,
				}
			}

			if err := client.CreateOrUpdateThenPoll(ctx, *id, *properties, virtualnetworklinks.CreateOrUpdateOperationOptions{}); err != nil {
				return fmt.Errorf("updating %s: %+v", *id, err)
			}

			return nil
		},
	}
}

func (r PrivateDNSResolverVirtualNetworkLinkResource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.PrivateDnsResolver.VirtualNetworkLinksClient

			id, err := virtualnetworklinks.ParseVirtualNetworkLinkID(metadata.ResourceData.Id())
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

			state := PrivateDNSResolverVirtualNetworkLinkModel{
				Name:                   id.VirtualNetworkLinkName,
				DnsForwardingRulesetId: dnsforwardingrulesets.NewDnsForwardingRulesetID(id.SubscriptionId, id.ResourceGroupName, id.DnsForwardingRulesetName).ID(),
			}

			properties := &model.Properties
			if properties.Metadata != nil {
				state.Metadata = *properties.Metadata
			}

			state.VirtualNetworkId = properties.VirtualNetwork.Id

			return metadata.Encode(&state)
		},
	}
}

func (r PrivateDNSResolverVirtualNetworkLinkResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.PrivateDnsResolver.VirtualNetworkLinksClient

			id, err := virtualnetworklinks.ParseVirtualNetworkLinkID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			if err := client.DeleteThenPoll(ctx, *id, virtualnetworklinks.DeleteOperationOptions{}); err != nil {
				return fmt.Errorf("deleting %s: %+v", id, err)
			}

			return nil
		},
	}
}
