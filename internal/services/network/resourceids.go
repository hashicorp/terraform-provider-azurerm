// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package network

// Core bits and pieces
//go:generate go run ../../tools/generator-resource-id/main.go -path=./ -name=Subnet -id=/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.Network/virtualNetworks/network1/subnets/subnet1 -rewrite=true
//go:generate go run ../../tools/generator-resource-id/main.go -path=./ -name=VirtualNetworkDnsServers -id=/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.Network/virtualNetworks/network1/dnsServers/default -rewrite=true
//go:generate go run ../../tools/generator-resource-id/main.go -path=./ -name=IpGroupCidr -id=/subscriptions/4fe7b06f-0eb5-489a-b053-f2afbdd8a4f3/resourceGroups/rg-core-hub-networking-ip-groups/providers/Microsoft.Network/ipGroups/my-ips/cidrs/127.0.0.1_32

// Private Link
//go:generate go run ../../tools/generator-resource-id/main.go -path=./ -name=PrivateDnsZoneConfig -id=/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.Network/privateEndpoints/endpoint1/privateDnsZoneGroups/privateDnsZoneGroup1/privateDnsZoneConfigs/privateDnsZoneConfig1
// ^ this one looks like it should be in Private DNS - alas no, it's actually nested and entirely managed within the Private Endpoint
