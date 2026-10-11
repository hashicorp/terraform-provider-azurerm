// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package schema

import (
	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/paloaltonetworks/2025-10-08/firewallresources"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

type DNSSettings struct {
	DnsServers      []string `tfschema:"dns_servers"`
	AzureDNS        bool     `tfschema:"use_azure_dns"`
	AzureDNSServers []string `tfschema:"azure_dns_servers"`
}

func DNSSettingsSchema() *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Type:     pluginsdk.TypeList,
		Optional: true,
		MaxItems: 1,
		Elem: &schema.Resource{
			Schema: map[string]*pluginsdk.Schema{
				"dns_servers": {
					Type:     pluginsdk.TypeList,
					Optional: true,
					MaxItems: 2,
					Elem: &pluginsdk.Schema{
						Type:         pluginsdk.TypeString,
						ValidateFunc: validation.IsIPv4Address,
					},
					ConflictsWith: []string{
						"dns_settings.0.use_azure_dns",
					},
				},

				"use_azure_dns": {
					Type:     pluginsdk.TypeBool,
					Optional: true,
					Default:  false,
					ConflictsWith: []string{
						"dns_settings.0.dns_servers",
					},
				},

				"azure_dns_servers": {
					Type:     pluginsdk.TypeList,
					Computed: true,
					Elem: &pluginsdk.Schema{
						Type: pluginsdk.TypeString,
					},
				},
			},
		},
	}
}

func ExpandDNSSettings(input []DNSSettings) firewallresources.DNSSettings {
	result := firewallresources.DNSSettings{
		EnableDnsProxy: pointer.To(firewallresources.DNSProxyDISABLED),
		EnabledDnsType: pointer.To(firewallresources.EnabledDNSTypeCUSTOM),
	}

	if len(input) == 1 {
		result.EnableDnsProxy = pointer.To(firewallresources.DNSProxyENABLED)
		dns := input[0]
		if len(dns.DnsServers) > 0 {
			dnsServers := make([]firewallresources.IPAddress, 0)
			for _, v := range dns.DnsServers {
				dnsServers = append(dnsServers, firewallresources.IPAddress{
					Address: pointer.To(v),
				})
			}
			result.DnsServers = pointer.To(dnsServers)
		}

		if dns.AzureDNS {
			result.EnabledDnsType = pointer.To(firewallresources.EnabledDNSTypeAZURE)
		}
	}

	return result
}

func FlattenDNSSettings(input firewallresources.DNSSettings) []DNSSettings {
	result := DNSSettings{}
	if pointer.From(input.EnableDnsProxy) == firewallresources.DNSProxyDISABLED {
		return []DNSSettings{}
	}

	useAzureDNS := pointer.From(input.EnabledDnsType) == firewallresources.EnabledDNSTypeAZURE

	if !useAzureDNS {
		dnsServers := make([]string, 0)
		for _, v := range pointer.From(input.DnsServers) {
			dnsServers = append(dnsServers, pointer.From(v.Address))
		}
		result.DnsServers = dnsServers
	} else {
		dnsServers := make([]string, 0)
		for _, v := range pointer.From(input.DnsServers) {
			dnsServers = append(dnsServers, pointer.From(v.Address))
		}
		result.AzureDNSServers = dnsServers
	}

	result.AzureDNS = useAzureDNS

	return []DNSSettings{result}
}
