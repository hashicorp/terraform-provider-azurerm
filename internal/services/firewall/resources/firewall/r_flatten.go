// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package firewall

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/zones"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2025-07-01/azurefirewalls"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2025-07-01/firewallpolicies"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func resourceFirewallSetFlatten(d *pluginsdk.ResourceData, id *azurefirewalls.AzureFirewallId, model *azurefirewalls.AzureFirewall) error {
	d.Set("name", id.AzureFirewallName)
	d.Set("resource_group_name", id.ResourceGroupName)

	if model != nil {
		d.Set("location", location.NormalizeNilable(model.Location))
		d.Set("zones", zones.FlattenUntyped(model.Zones))

		if props := model.Properties; props != nil {
			if err := d.Set("ip_configuration", flattenFirewallIPConfigurations(props.IPConfigurations)); err != nil {
				return fmt.Errorf("setting `ip_configuration`: %+v", err)
			}
			managementIPConfigs := make([]any, 0)
			if props.ManagementIPConfiguration != nil {
				managementIPConfigs = flattenFirewallIPConfigurations(&[]azurefirewalls.AzureFirewallIPConfiguration{
					*props.ManagementIPConfiguration,
				})
			}
			if err := d.Set("management_ip_configuration", managementIPConfigs); err != nil {
				return fmt.Errorf("setting `management_ip_configuration`: %+v", err)
			}

			d.Set("threat_intel_mode", pointer.FromEnum(props.ThreatIntelMode))

			dnsProxyEnabled, dnsServers := flattenFirewallAdditionalProperty(props.AdditionalProperties)
			if err := d.Set("dns_proxy_enabled", dnsProxyEnabled); err != nil {
				return fmt.Errorf("setting `dns_proxy_enabled`: %+v", err)
			}
			if err := d.Set("dns_servers", dnsServers); err != nil {
				return fmt.Errorf("setting `dns_servers`: %+v", err)
			}

			if err := d.Set("private_ip_ranges", flattenFirewallPrivateIpRange(props.AdditionalProperties)); err != nil {
				return fmt.Errorf("setting `private_ip_ranges`: %+v", err)
			}

			firewallPolicyId := ""
			if props.FirewallPolicy != nil && props.FirewallPolicy.Id != nil {
				firewallPolicyId = *props.FirewallPolicy.Id
				if policyId, err := firewallpolicies.ParseFirewallPolicyIDInsensitively(firewallPolicyId); err == nil {
					firewallPolicyId = policyId.ID()
				}
			}
			d.Set("firewall_policy_id", firewallPolicyId)

			if sku := props.Sku; sku != nil {
				d.Set("sku_name", pointer.FromEnum(sku.Name))
				d.Set("sku_tier", pointer.FromEnum(sku.Tier))
			}

			if err := d.Set("virtual_hub", flattenFirewallVirtualHubSetting(props)); err != nil {
				return fmt.Errorf("setting `virtual_hub`: %+v", err)
			}
		}

		if err := tags.FlattenAndSet(d, model.Tags); err != nil {
			return fmt.Errorf("flattening `tags`: %+v", err)
		}
	}

	return pluginsdk.SetResourceIdentityData(d, id)
}

func flattenFirewallIPConfigurations(input *[]azurefirewalls.AzureFirewallIPConfiguration) []any {
	result := make([]any, 0)
	if input == nil {
		return result
	}

	for _, v := range *input {
		afIPConfig := make(map[string]any)
		props := v.Properties
		if props == nil {
			continue
		}

		if name := v.Name; name != nil {
			afIPConfig["name"] = *name
		}

		if subnet := props.Subnet; subnet != nil {
			if id := subnet.Id; id != nil {
				afIPConfig["subnet_id"] = *id
			}
		}

		if ipAddress := props.PrivateIPAddress; ipAddress != nil {
			afIPConfig["private_ip_address"] = *ipAddress
		}

		if pip := props.PublicIPAddress; pip != nil {
			if id := pip.Id; id != nil {
				afIPConfig["public_ip_address_id"] = *id
			}
		}
		result = append(result, afIPConfig)
	}

	return result
}

func flattenFirewallAdditionalProperty(input *map[string]string) (enabled any, servers []any) {
	if input == nil || len(*input) == 0 {
		return nil, []any{}
	}

	if enabledPtr, ok := (*input)["Network.DNS.EnableProxy"]; ok {
		enabled = enabledPtr == "true"
	}

	if serversPtr, ok := (*input)["Network.DNS.Servers"]; ok {
		for val := range strings.SplitSeq(serversPtr, ",") {
			servers = append(servers, val)
		}
	}
	return
}

func flattenFirewallPrivateIpRange(input *map[string]string) []any {
	if input == nil || len(*input) == 0 {
		return []any{}
	}

	attrs := *input
	rangeSlice := []string{}
	if privateIpRanges := attrs["Network.SNAT.PrivateRanges"]; privateIpRanges != "" {
		rangeSlice = strings.Split(attrs["Network.SNAT.PrivateRanges"], ",")
	}
	return pluginsdk.FlattenSlice(&rangeSlice)
}

func flattenFirewallVirtualHubSetting(props *azurefirewalls.AzureFirewallPropertiesFormat) []any {
	if props.VirtualHub == nil {
		return []any{}
	}

	var (
		publicIpCount int
		publicIps     []string
		privateIp     string
	)
	if hubIP := props.HubIPAddresses; hubIP != nil {
		if hubIP.PrivateIPAddress != nil {
			privateIp = *hubIP.PrivateIPAddress
		}
		if pubIPs := hubIP.PublicIPs; pubIPs != nil {
			if pubIPs.Count != nil {
				publicIpCount = int(*pubIPs.Count)
			}
			if pubIPs.Addresses != nil {
				for _, addr := range *pubIPs.Addresses {
					if addr.Address != nil {
						publicIps = append(publicIps, *addr.Address)
					}
				}
			}
		}
	}

	return []any{
		map[string]any{
			"virtual_hub_id":      pointer.From(props.VirtualHub.Id),
			"public_ip_count":     publicIpCount,
			"public_ip_addresses": publicIps,
			"private_ip_address":  privateIp,
		},
	}
}
