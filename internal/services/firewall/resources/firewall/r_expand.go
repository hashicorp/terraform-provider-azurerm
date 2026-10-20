// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package firewall

import (
	"slices"
	"strings"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2025-07-01/azurefirewalls"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func expandFirewallIPConfigurations(configs []any) (*[]azurefirewalls.AzureFirewallIPConfiguration, *[]string, *[]string, error) {
	ipConfigs := make([]azurefirewalls.AzureFirewallIPConfiguration, 0)
	subnetNamesToLock := make([]string, 0)
	virtualNetworkNamesToLock := make([]string, 0)

	for _, configRaw := range configs {
		data := configRaw.(map[string]any)
		name := data["name"].(string)
		subnetId := data["subnet_id"].(string)
		pubID := data["public_ip_address_id"].(string)

		ipConfig := azurefirewalls.AzureFirewallIPConfiguration{
			Name:       pointer.To(name),
			Properties: &azurefirewalls.AzureFirewallIPConfigurationPropertiesFormat{},
		}

		if pubID != "" {
			ipConfig.Properties.PublicIPAddress = &azurefirewalls.SubResource{
				Id: pointer.To(pubID),
			}
		}

		if subnetId != "" {
			subnetID, err := commonids.ParseSubnetID(subnetId)
			if err != nil {
				return nil, nil, nil, err
			}

			if !slices.Contains(subnetNamesToLock, subnetID.SubnetName) {
				subnetNamesToLock = append(subnetNamesToLock, subnetID.SubnetName)
			}

			if !slices.Contains(virtualNetworkNamesToLock, subnetID.VirtualNetworkName) {
				virtualNetworkNamesToLock = append(virtualNetworkNamesToLock, subnetID.VirtualNetworkName)
			}

			ipConfig.Properties.Subnet = &azurefirewalls.SubResource{
				Id: pointer.To(subnetId),
			}
		}
		ipConfigs = append(ipConfigs, ipConfig)
	}
	return &ipConfigs, &subnetNamesToLock, &virtualNetworkNamesToLock, nil
}

func expandFirewallAdditionalProperty(d *pluginsdk.ResourceData) map[string]string {
	// Swagger issue asking finalize these properties: https://github.com/Azure/azure-rest-api-specs/issues/11278
	res := map[string]string{}
	if servers := d.Get("dns_servers").([]any); len(servers) > 0 {
		var servs []string
		for _, server := range servers {
			servs = append(servs, server.(string))
		}
		res["Network.DNS.EnableProxy"] = "true"
		res["Network.DNS.Servers"] = strings.Join(servs, ",")
	}
	if enabled := d.Get("dns_proxy_enabled").(bool); enabled {
		res["Network.DNS.EnableProxy"] = "true"
	}
	return res
}

func expandFirewallPrivateIpRange(input []any) map[string]string {
	if len(input) == 0 {
		return nil
	}

	rangeSlice := *pluginsdk.ExpandStringSlice(input)
	if len(rangeSlice) == 0 {
		return nil
	}

	// Swagger issue asking finalize these properties: https://github.com/Azure/azure-rest-api-specs/issues/10015
	return map[string]string{
		"Network.SNAT.PrivateRanges": strings.Join(rangeSlice, ","),
	}
}

func expandFirewallVirtualHubSetting(existing *azurefirewalls.AzureFirewall, input []any) (vhub *azurefirewalls.SubResource, ipAddresses *azurefirewalls.HubIPAddresses, ok bool) {
	if len(input) == 0 {
		return nil, nil, false
	}

	b := input[0].(map[string]any)

	// The API requires both "Count" and "Addresses" for the "PublicIPs" setting.
	// The "Count" means how many PIP to provision.
	// The "Addresses" means differently in different cases:
	// - Create: only "Count" is needed, "Addresses" is not necessary
	// - Update: both "Count" and "Addresses" are needed:
	//   Scale up: "Addresses" should remain same as before scaling up
	//   Scale down: "Addresses" should indicate the addresses to be retained (in this case we retain the first new "Count" ones)
	newCount := b["public_ip_count"].(int)
	var addresses *[]azurefirewalls.AzureFirewallPublicIPAddress
	if existing != nil {
		if prop := existing.Properties; prop != nil {
			if ipaddress := prop.HubIPAddresses; ipaddress != nil {
				if pips := ipaddress.PublicIPs; pips != nil {
					if count := pips.Count; count != nil {
						oldCount := int(*count)
						addresses = pips.Addresses

						// In case of scale down, keep the first new "Count" addresses.
						if oldCount > newCount {
							keptAddresses := make([]azurefirewalls.AzureFirewallPublicIPAddress, newCount)
							for i := range newCount {
								keptAddresses[i] = (*addresses)[i]
							}
							addresses = &keptAddresses
						}
					}
				}
			}
		}
	}

	vhub = &azurefirewalls.SubResource{Id: pointer.To(b["virtual_hub_id"].(string))}
	ipAddresses = &azurefirewalls.HubIPAddresses{
		PublicIPs: &azurefirewalls.HubPublicIPAddresses{
			Count:     pointer.To(int64(b["public_ip_count"].(int))),
			Addresses: addresses,
		},
	}

	return vhub, ipAddresses, true
}
