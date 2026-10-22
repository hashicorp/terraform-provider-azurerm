// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package palo_alto_next_generation_firewall_virtual_network_strata_cloud_manager

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/paloaltonetworks/2025-10-08/firewallresources"
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

			var model NextGenerationFirewallVNetStrataCloudManagerModel
			if err = metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding model: %+v", err)
			}

			existing, err := client.FirewallsGet(ctx, *id)
			if err != nil {
				if response.WasNotFound(existing.HttpResponse) {
					return metadata.MarkAsGone(id)
				}
				return fmt.Errorf("retrieving %s: %+v", *id, err)
			}

			if existing.Model == nil {
				return fmt.Errorf("retrieving %s: model was nil", *id)
			}

			firewall := *existing.Model
			props := firewall.Properties

			if metadata.ResourceData.HasChange("plan_id") {
				props.PlanData.PlanId = model.PlanId
			}

			if metadata.ResourceData.HasChange("strata_cloud_manager_tenant_name") {
				if props.StrataCloudManagerConfig == nil {
					props.StrataCloudManagerConfig = &firewallresources.StrataCloudManagerConfig{}
				}
				props.StrataCloudManagerConfig.CloudManagerName = model.StrataCloudManagerTenantName
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

			firewall.Properties = props

			if metadata.ResourceData.HasChange("identity") {
				identityValue, err := identity.ExpandUserAssignedMapFromModel(model.Identity)
				if err != nil {
					return fmt.Errorf("expanding `identity`: %+v", err)
				}
				firewall.Identity = identityValue
			}

			if metadata.ResourceData.HasChange("tags") {
				firewall.Tags = tags.Expand(model.Tags)
			}

			if err = client.FirewallsCreateOrUpdateThenPoll(ctx, *id, firewall); err != nil {
				return fmt.Errorf("updating %s: %+v", *id, err)
			}

			return nil
		},
	}
}
