// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package schema

import (
	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/paloaltonetworks/2023-09-01/firewalls"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
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
						ValidateFunc: validate.IPv4Address,
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

func ExpandDNSSettings(input []DNSSettings) firewalls.DNSSettings {
	result := firewalls.DNSSettings{
		EnableDnsProxy: pointer.To(firewalls.DNSProxyDISABLED),
		EnabledDnsType: pointer.To(firewalls.EnabledDNSTypeCUSTOM),
	}

	if len(input) == 1 {
		result.EnableDnsProxy = pointer.To(firewalls.DNSProxyENABLED)
		dns := input[0]
		if len(dns.DnsServers) > 0 {
			dnsServers := make([]firewalls.IPAddress, 0)
			for _, v := range dns.DnsServers {
				dnsServers = append(dnsServers, firewalls.IPAddress{
					Address: pointer.To(v),
				})
			}
			result.DnsServers = pointer.To(dnsServers)
		}

		if dns.AzureDNS {
			result.EnabledDnsType = pointer.To(firewalls.EnabledDNSTypeAZURE)
		}
	}

	return result
}

func FlattenDNSSettings(input firewalls.DNSSettings) []DNSSettings {
	result := DNSSettings{}
	if pointer.From(input.EnableDnsProxy) == firewalls.DNSProxyDISABLED {
		return []DNSSettings{}
	}

	useAzureDNS := pointer.From(input.EnabledDnsType) == firewalls.EnabledDNSTypeAZURE

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
