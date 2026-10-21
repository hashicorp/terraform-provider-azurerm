// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package network_manager_ipam_pool_static_cidr

func (Resource) ModelObject() any {
	return &ManagerIpamPoolStaticCidrResourceModel{}
}

type ManagerIpamPoolStaticCidrResourceModel struct {
	AddressPrefixes               []string `tfschema:"address_prefixes"`
	IpamPoolId                    string   `tfschema:"ipam_pool_id"`
	Name                          string   `tfschema:"name"`
	NumberOfIPAddressesToAllocate string   `tfschema:"number_of_ip_addresses_to_allocate"`
}
