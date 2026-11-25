// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package palo_alto_next_generation_firewall_virtual_hub_local_rulestack

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/paloaltonetworks/2025-10-08/firewallresources"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/paloalto/schema"
)

func (r Resource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.PaloAlto.FirewallResources

			id, err := firewallresources.ParseFirewallID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			var state NextGenerationFirewallVHubLocalRuleStackModel

			existing, err := client.FirewallsGet(ctx, *id)
			if err != nil {
				if response.WasNotFound(existing.HttpResponse) {
					return metadata.MarkAsGone(id)
				}
				return fmt.Errorf("reading %s: %+v", *id, err)
			}

			state.Name = id.FirewallName
			state.ResourceGroupName = id.ResourceGroupName

			if model := existing.Model; model != nil {
				props := model.Properties

				state.DNSSettings = schema.FlattenDNSSettings(props.DnsSettings)

				netProfile, err := schema.FlattenNetworkProfileVHub(props.NetworkProfile)
				if err != nil {
					return fmt.Errorf("parsing Network Profile for %s: %+v", *id, err)
				}

				state.NetworkProfile = []schema.NetworkProfileVHub{*netProfile}

				state.FrontEnd = schema.FlattenDestinationNAT(props.FrontEndSettings)

				state.RuleStackId = pointer.From(props.AssociatedRulestack.ResourceId)

				state.MarketplaceOfferId = props.MarketplaceDetails.OfferId

				state.PlanId = props.PlanData.PlanId

				state.Tags = tags.Flatten(model.Tags)
			}

			return metadata.Encode(&state)
		},
	}
}
