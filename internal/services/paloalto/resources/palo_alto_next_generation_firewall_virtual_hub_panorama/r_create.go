// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package palo_alto_next_generation_firewall_virtual_hub_panorama

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/paloaltonetworks/2025-10-08/firewallresources"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/paloalto/schema"
)

func (r Resource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 2 * time.Hour,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.PaloAlto.FirewallResources

			var model NextGenerationFirewallVHubPanoramaResourceModel

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

			firewall := firewallresources.FirewallResource{
				Location: location.Normalize(model.Location),
				Properties: firewallresources.FirewallDeploymentProperties{
					PanoramaConfig: &firewallresources.PanoramaConfig{
						ConfigString: model.PanoramaBase64Config,
					},
					IsPanoramaManaged: pointer.To(firewallresources.BooleanEnumTRUE),
					DnsSettings:       schema.ExpandDNSSettings(model.DNSSettings),
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

			if err := client.FirewallsCreateOrUpdateCallbackThenPoll(ctx, id, firewall, metadata.SetIDCallback(&id)); err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}
			metadata.SetID(id)

			return nil
		},
	}
}
