// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package palo_alto_next_generation_firewall_virtual_hub_strata_cloud_manager

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
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

			var state NextGenerationFirewallVHubStrataCloudManagerModel

			existing, err := client.FirewallsGet(ctx, *id)
			if err != nil {
				if response.WasNotFound(existing.HttpResponse) {
					return metadata.MarkAsGone(id)
				}
				return fmt.Errorf("retrieving %s: %+v", *id, err)
			}

			state.Name = id.FirewallName

			state.ResourceGroupName = id.ResourceGroupName

			if model := existing.Model; model != nil {
				props := model.Properties
				state.Location = location.Normalize(model.Location)
				state.DNSSettings = schema.FlattenDNSSettings(props.DnsSettings)
				state.FrontEnd = schema.FlattenDestinationNAT(props.FrontEndSettings)

				networkProfile, err := schema.FlattenNetworkProfileVHub(props.NetworkProfile)
				if err != nil {
					return fmt.Errorf("flattening Network Profile for %s: %+v", *id, err)
				}
				state.NetworkProfile = []schema.NetworkProfileVHub{*networkProfile}

				state.MarketplaceOfferId = props.MarketplaceDetails.OfferId
				state.PlanId = props.PlanData.PlanId

				expandedIdentity, err := identity.FlattenUserAssignedMapToModel(model.Identity)
				if err != nil {
					return fmt.Errorf("flattening `identity`: %+v", err)
				}
				state.Identity = pointer.From(expandedIdentity)

				state.Tags = tags.Flatten(model.Tags)

				if strataCloudManagerConfig := props.StrataCloudManagerConfig; strataCloudManagerConfig != nil {
					state.StrataCloudManagerTenantName = strataCloudManagerConfig.CloudManagerName
				}
			}

			return metadata.Encode(&state)
		},
	}
}
