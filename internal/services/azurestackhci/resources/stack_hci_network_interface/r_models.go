// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package stack_hci_network_interface

func (Resource) ModelObject() any {
	return &StackHCINetworkInterfaceResourceModel{}
}

type StackHCINetworkInterfaceResourceModel struct {
	Name              string                         `tfschema:"name"`
	ResourceGroupName string                         `tfschema:"resource_group_name"`
	Location          string                         `tfschema:"location"`
	CustomLocationId  string                         `tfschema:"custom_location_id"`
	DNSServers        []string                       `tfschema:"dns_servers"`
	IPConfiguration   []StackHCIIPConfigurationModel `tfschema:"ip_configuration"`
	MACAddress        string                         `tfschema:"mac_address"`
	Tags              map[string]any                 `tfschema:"tags"`
}

type StackHCIIPConfigurationModel struct {
	Gateway          string `tfschema:"gateway"`
	PrefixLength     string `tfschema:"prefix_length"`
	PrivateIPAddress string `tfschema:"private_ip_address"`
	SubnetID         string `tfschema:"subnet_id"`
}
