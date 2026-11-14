// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package palo_alto_next_generation_firewall_virtual_hub_local_rulestack

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/paloaltonetworks/2025-10-08/firewallresources"
	"github.com/hashicorp/go-azure-sdk/resource-manager/paloaltonetworks/2025-10-08/localrulestackresources"
	"github.com/hashicorp/terraform-provider-azurerm/internal/locks"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/paloalto/schema"
)

func (r Resource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 3 * time.Hour,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.PaloAlto.FirewallResources
			localRuleStackClient := metadata.Client.PaloAlto.LocalRulestackResources
			var model NextGenerationFirewallVHubLocalRuleStackModel

			if err := metadata.Decode(&model); err != nil {
				return err
			}

			id := firewallresources.NewFirewallID(metadata.Client.Account.SubscriptionId, model.ResourceGroupName, model.Name)

			if !metadata.Client.Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
				existing, err := client.FirewallsGet(ctx, id)
				if err != nil {
					if !response.WasNotFound(existing.HttpResponse) {
						return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
					}
				}
				if !response.WasNotFound(existing.HttpResponse) {
					return metadata.ResourceRequiresImport(r.ResourceType(), id)
				}
			}

			ruleStackID, err := localrulestackresources.ParseLocalRulestackID(model.RuleStackId)
			if err != nil {
				return err
			}

			ruleStack, err := localRuleStackClient.LocalRulestacksGet(ctx, *ruleStackID)
			if err != nil {
				return fmt.Errorf("reading %s for %s: %+v", ruleStackID, id, err)
			}

			loc := location.Normalize(ruleStack.Model.Location)

			firewall := firewallresources.FirewallResource{
				Location: loc,
				Properties: firewallresources.FirewallDeploymentProperties{
					AssociatedRulestack: &firewallresources.RulestackDetails{
						ResourceId: pointer.To(ruleStackID.ID()),
						Location:   pointer.To(loc),
					},
					DnsSettings: schema.ExpandDNSSettings(model.DNSSettings),
					MarketplaceDetails: firewallresources.MarketplaceDetails{
						OfferId:     model.MarketplaceOfferId,
						PublisherId: "paloaltonetworks",
					},
					NetworkProfile: schema.ExpandNetworkProfileVHub(model.NetworkProfile),
					PlanData: firewallresources.PlanData{
						BillingCycle: firewallresources.BillingCycleMONTHLY,
						PlanId:       model.PlanId,
					},
					FrontEndSettings: schema.ExpandDestinationNAT(model.FrontEnd),
				},

				Tags: tags.Expand(model.Tags),
			}

			locks.ByID(ruleStackID.ID())
			defer locks.UnlockByID(ruleStackID.ID())

			if err := client.FirewallsCreateOrUpdateCallbackThenPoll(ctx, id, firewall, metadata.SetIDCallback(&id)); err != nil {
				return err
			}
			metadata.SetID(id)

			return nil
		},
	}
}
