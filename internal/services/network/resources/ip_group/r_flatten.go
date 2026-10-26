// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package ip_group

import (
	"fmt"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2025-07-01/azurefirewalls"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2025-07-01/firewallpolicies"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2025-07-01/ipgroups"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func resourceIpGroupFlatten(d *pluginsdk.ResourceData, id *ipgroups.IPGroupId, model *ipgroups.IPGroup) error {
	d.Set("name", id.IpGroupName)
	d.Set("resource_group_name", id.ResourceGroupName)

	if model != nil {
		d.Set("location", location.NormalizeNilable(model.Location))

		if props := model.Properties; props != nil {
			if props.IPAddresses == nil {
				return fmt.Errorf("list of ipAddresses returned is nil")
			}
			if err := d.Set("cidrs", props.IPAddresses); err != nil {
				return fmt.Errorf("setting `cidrs`: %+v", err)
			}

			firewallIDs := make([]string, 0)
			for _, idStr := range getIds(props.Firewalls) {
				firewallID, err := azurefirewalls.ParseAzureFirewallIDInsensitively(idStr)
				if err != nil {
					return fmt.Errorf("parsing Azure Firewall ID %q: %+v", idStr, err)
				}
				firewallIDs = append(firewallIDs, firewallID.ID())
			}
			d.Set("firewall_ids", firewallIDs)

			firewallPolicyIDs := make([]string, 0)
			for _, idStr := range getIds(props.FirewallPolicies) {
				policyID, err := firewallpolicies.ParseFirewallPolicyIDInsensitively(idStr)
				if err != nil {
					return fmt.Errorf("parsing Azure Firewall Policy ID %q: %+v", idStr, err)
				}
				firewallPolicyIDs = append(firewallPolicyIDs, policyID.ID())
			}
			d.Set("firewall_policy_ids", firewallPolicyIDs)
		}
		if err := tags.FlattenAndSet(d, model.Tags); err != nil {
			return err
		}
	}

	return pluginsdk.SetResourceIdentityData(d, id)
}
