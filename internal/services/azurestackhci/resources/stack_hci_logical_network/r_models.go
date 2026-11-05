// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package stack_hci_logical_network

func (Resource) ModelObject() any {
	return &StackHCILogicalNetworkResourceModel{}
}

type StackHCILogicalNetworkResourceModel struct {
	Name              string                `tfschema:"name"`
	ResourceGroupName string                `tfschema:"resource_group_name"`
	Location          string                `tfschema:"location"`
	CustomLocationId  string                `tfschema:"custom_location_id"`
	DNSServers        []string              `tfschema:"dns_servers"`
	Subnet            []StackHCISubnetModel `tfschema:"subnet"`
	VirtualSwitchName string                `tfschema:"virtual_switch_name"`
	Tags              map[string]any        `tfschema:"tags"`
}

type StackHCISubnetModel struct {
	AddressPrefix      string                `tfschema:"address_prefix"`
	IpAllocationMethod string                `tfschema:"ip_allocation_method"`
	IpPool             []StackHCIIPPoolModel `tfschema:"ip_pool"`
	Route              []StackHCIRouteModel  `tfschema:"route"`
	VlanId             int64                 `tfschema:"vlan_id"`
}

type StackHCIIPPoolModel struct {
	Start string `tfschema:"start"`
	End   string `tfschema:"end"`
}

type StackHCIRouteModel struct {
	Name             string `tfschema:"name"`
	AddressPrefix    string `tfschema:"address_prefix"`
	NextHopIpAddress string `tfschema:"next_hop_ip_address"`
}
