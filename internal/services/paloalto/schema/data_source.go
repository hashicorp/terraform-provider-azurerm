// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package schema

import "github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"

func DNSSettingsDataSourceSchema() *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Type:     pluginsdk.TypeList,
		Computed: true,
		Elem: &pluginsdk.Resource{
			Schema: map[string]*pluginsdk.Schema{
				"azure_dns_servers": computedStringListSchema(),
				"dns_servers":       computedStringListSchema(),
				"use_azure_dns": {
					Type:     pluginsdk.TypeBool,
					Computed: true,
				},
			},
		},
	}
}

func DestinationNATDataSourceSchema() *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Type:     pluginsdk.TypeList,
		Computed: true,
		Elem: &pluginsdk.Resource{
			Schema: map[string]*pluginsdk.Schema{
				"backend_config":  endpointConfigurationDataSourceSchema(false),
				"frontend_config": endpointConfigurationDataSourceSchema(true),
				"name": {
					Type:     pluginsdk.TypeString,
					Computed: true,
				},
				"protocol": {
					Type:     pluginsdk.TypeString,
					Computed: true,
				},
			},
		},
	}
}

func VHubNetworkProfileDataSourceSchema() *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Type:     pluginsdk.TypeList,
		Computed: true,
		Elem: &pluginsdk.Resource{
			Schema: map[string]*pluginsdk.Schema{
				"egress_nat_ip_address_ids": computedStringListSchema(),
				"egress_nat_ip_addresses":   computedStringListSchema(),
				"ip_of_trust_for_user_defined_routes": {
					Type:     pluginsdk.TypeString,
					Computed: true,
				},
				"network_virtual_appliance_id": {
					Type:     pluginsdk.TypeString,
					Computed: true,
				},
				"public_ip_address_ids":  computedStringListSchema(),
				"public_ip_addresses":    computedStringListSchema(),
				"trusted_address_ranges": computedStringListSchema(),
				"trusted_subnet_id": {
					Type:     pluginsdk.TypeString,
					Computed: true,
				},
				"untrusted_subnet_id": {
					Type:     pluginsdk.TypeString,
					Computed: true,
				},
				"virtual_hub_id": {
					Type:     pluginsdk.TypeString,
					Computed: true,
				},
			},
		},
	}
}

func VnetNetworkProfileDataSourceSchema() *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Type:     pluginsdk.TypeList,
		Computed: true,
		Elem: &pluginsdk.Resource{
			Schema: map[string]*pluginsdk.Schema{
				"egress_nat_ip_address_ids": computedStringListSchema(),
				"egress_nat_ip_addresses":   computedStringListSchema(),
				"public_ip_address_ids":     computedStringListSchema(),
				"public_ip_addresses":       computedStringListSchema(),
				"trusted_address_ranges":    computedStringListSchema(),
				"vnet_configuration": {
					Type:     pluginsdk.TypeList,
					Computed: true,
					Elem: &pluginsdk.Resource{
						Schema: map[string]*pluginsdk.Schema{
							"ip_of_trust_for_user_defined_routes": {
								Type:     pluginsdk.TypeString,
								Computed: true,
							},
							"trusted_subnet_id": {
								Type:     pluginsdk.TypeString,
								Computed: true,
							},
							"untrusted_subnet_id": {
								Type:     pluginsdk.TypeString,
								Computed: true,
							},
							"virtual_network_id": {
								Type:     pluginsdk.TypeString,
								Computed: true,
							},
						},
					},
				},
			},
		},
	}
}

func computedStringListSchema() *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Type:     pluginsdk.TypeList,
		Computed: true,
		Elem: &pluginsdk.Schema{
			Type: pluginsdk.TypeString,
		},
	}
}

func endpointConfigurationDataSourceSchema(frontend bool) *pluginsdk.Schema {
	addressField := "public_ip_address"
	if frontend {
		addressField = "public_ip_address_id"
	}

	return &pluginsdk.Schema{
		Type:     pluginsdk.TypeList,
		Computed: true,
		Elem: &pluginsdk.Resource{
			Schema: map[string]*pluginsdk.Schema{
				addressField: {
					Type:     pluginsdk.TypeString,
					Computed: true,
				},
				"port": {
					Type:     pluginsdk.TypeInt,
					Computed: true,
				},
			},
		},
	}
}
