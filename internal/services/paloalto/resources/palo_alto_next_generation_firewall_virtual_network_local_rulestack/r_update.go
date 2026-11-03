// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package palo_alto_next_generation_firewall_virtual_network_local_rulestack

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/paloaltonetworks/2025-10-08/firewallresources"
	"github.com/hashicorp/go-azure-sdk/resource-manager/paloaltonetworks/2025-10-08/localrulestackresources"
	"github.com/hashicorp/terraform-provider-azurerm/internal/locks"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/paloalto/schema"
)

func (r Resource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 3 * time.Hour,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.PaloAlto.FirewallResources

			id, err := firewallresources.ParseFirewallID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			model := NextGenerationFirewallVnetLocalRulestackModel{}

			if err = metadata.Decode(&model); err != nil {
				return err
			}

			existing, err := client.FirewallsGet(ctx, *id)
			if err != nil {
				if response.WasNotFound(existing.HttpResponse) {
					return metadata.MarkAsGone(id)
				}
				return fmt.Errorf("reading %s: %+v", *id, err)
			}

			firewall := *existing.Model
			props := firewall.Properties

			if metadata.ResourceData.HasChange("rulestack_id") {
				ruleStackID, err := localrulestackresources.ParseLocalRulestackID(model.RuleStackId)
				if err != nil {
					return err
				}

				props.AssociatedRulestack = &firewallresources.RulestackDetails{
					Location:    props.AssociatedRulestack.Location,
					RulestackId: pointer.To(ruleStackID.ID()),
				}
				locks.ByID(ruleStackID.ID())
				defer locks.UnlockByID(ruleStackID.ID())
			}

			if metadata.ResourceData.HasChange("network_profile") {
				props.NetworkProfile = schema.ExpandNetworkProfileVnet(model.NetworkProfile)
			}

			if metadata.ResourceData.HasChange("dns_settings") {
				props.DnsSettings = schema.ExpandDNSSettings(model.DNSSettings)
			}

			if metadata.ResourceData.HasChange("destination_nat") {
				props.FrontEndSettings = schema.ExpandDestinationNAT(model.FrontEnd)
			}

			if metadata.ResourceData.HasChange("plan_id") {
				props.PlanData.PlanId = model.PlanId
			}

			firewall.Properties = props

			if metadata.ResourceData.HasChange("tags") {
				firewall.Tags = tags.Expand(model.Tags)
			}

			if err = client.FirewallsCreateOrUpdateThenPoll(ctx, *id, firewall); err != nil {
				return err
			}

			return nil
		},
	}
}
