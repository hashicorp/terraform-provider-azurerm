// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package parse

// NOTE: this file is generated via 'go:generate' - manual changes will be overwritten

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

type VirtualNetworkDnsServersId struct {
	SubscriptionId     string
	ResourceGroup      string
	VirtualNetworkName string
	DnsServerName      string
}

func NewVirtualNetworkDnsServersID(subscriptionId, resourceGroup, virtualNetworkName, dnsServerName string) VirtualNetworkDnsServersId {
	return VirtualNetworkDnsServersId{
		SubscriptionId:     subscriptionId,
		ResourceGroup:      resourceGroup,
		VirtualNetworkName: virtualNetworkName,
		DnsServerName:      dnsServerName,
	}
}

func (id VirtualNetworkDnsServersId) String() string {
	segments := []string{
		fmt.Sprintf("Dns Server Name %q", id.DnsServerName),
		fmt.Sprintf("Virtual Network Name %q", id.VirtualNetworkName),
		fmt.Sprintf("Resource Group %q", id.ResourceGroup),
	}
	segmentsStr := strings.Join(segments, " / ")
	return fmt.Sprintf("%s: (%s)", "Virtual Network Dns Servers", segmentsStr)
}

func (id VirtualNetworkDnsServersId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Network/virtualNetworks/%s/dnsServers/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroup, id.VirtualNetworkName, id.DnsServerName)
}

// VirtualNetworkDnsServersID parses a VirtualNetworkDnsServers ID into an VirtualNetworkDnsServersId struct
func VirtualNetworkDnsServersID(input string) (*VirtualNetworkDnsServersId, error) {
	id, err := resourceids.ParseAzureResourceID(input)
	if err != nil {
		return nil, fmt.Errorf("parsing %q as an VirtualNetworkDnsServers ID: %+v", input, err)
	}

	resourceId := VirtualNetworkDnsServersId{
		SubscriptionId: id.SubscriptionID,
		ResourceGroup:  id.ResourceGroup,
	}

	if resourceId.SubscriptionId == "" {
		return nil, fmt.Errorf("ID was missing the 'subscriptions' element")
	}

	if resourceId.ResourceGroup == "" {
		return nil, fmt.Errorf("ID was missing the 'resourceGroups' element")
	}

	if resourceId.VirtualNetworkName, err = id.PopSegment("virtualNetworks"); err != nil {
		return nil, err
	}
	if resourceId.DnsServerName, err = id.PopSegment("dnsServers"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &resourceId, nil
}

// VirtualNetworkDnsServersIDInsensitively parses an VirtualNetworkDnsServers ID into an VirtualNetworkDnsServersId struct, insensitively
// This should only be used to parse an ID for rewriting, the VirtualNetworkDnsServersID
// method should be used instead for validation etc.
//
// Whilst this may seem strange, this enables Terraform have consistent casing
// which works around issues in Core, whilst handling broken API responses.
func VirtualNetworkDnsServersIDInsensitively(input string) (*VirtualNetworkDnsServersId, error) {
	id, err := resourceids.ParseAzureResourceID(input)
	if err != nil {
		return nil, err
	}

	resourceId := VirtualNetworkDnsServersId{
		SubscriptionId: id.SubscriptionID,
		ResourceGroup:  id.ResourceGroup,
	}

	if resourceId.SubscriptionId == "" {
		return nil, fmt.Errorf("ID was missing the 'subscriptions' element")
	}

	if resourceId.ResourceGroup == "" {
		return nil, fmt.Errorf("ID was missing the 'resourceGroups' element")
	}

	// find the correct casing for the 'virtualNetworks' segment
	virtualNetworksKey := "virtualNetworks"
	for key := range id.Path {
		if strings.EqualFold(key, virtualNetworksKey) {
			virtualNetworksKey = key
			break
		}
	}
	if resourceId.VirtualNetworkName, err = id.PopSegment(virtualNetworksKey); err != nil {
		return nil, err
	}

	// find the correct casing for the 'dnsServers' segment
	dnsServersKey := "dnsServers"
	for key := range id.Path {
		if strings.EqualFold(key, dnsServersKey) {
			dnsServersKey = key
			break
		}
	}
	if resourceId.DnsServerName, err = id.PopSegment(dnsServersKey); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &resourceId, nil
}
