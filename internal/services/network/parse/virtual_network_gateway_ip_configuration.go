// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package parse

// NOTE: this file is generated via 'go:generate' - manual changes will be overwritten

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

type VirtualNetworkGatewayIpConfigurationId struct {
	SubscriptionId            string
	ResourceGroup             string
	VirtualNetworkGatewayName string
	IpConfigurationName       string
}

func NewVirtualNetworkGatewayIpConfigurationID(subscriptionId, resourceGroup, virtualNetworkGatewayName, ipConfigurationName string) VirtualNetworkGatewayIpConfigurationId {
	return VirtualNetworkGatewayIpConfigurationId{
		SubscriptionId:            subscriptionId,
		ResourceGroup:             resourceGroup,
		VirtualNetworkGatewayName: virtualNetworkGatewayName,
		IpConfigurationName:       ipConfigurationName,
	}
}

func (id VirtualNetworkGatewayIpConfigurationId) String() string {
	segments := []string{
		fmt.Sprintf("Ip Configuration Name %q", id.IpConfigurationName),
		fmt.Sprintf("Virtual Network Gateway Name %q", id.VirtualNetworkGatewayName),
		fmt.Sprintf("Resource Group %q", id.ResourceGroup),
	}
	segmentsStr := strings.Join(segments, " / ")
	return fmt.Sprintf("%s: (%s)", "Virtual Network Gateway Ip Configuration", segmentsStr)
}

func (id VirtualNetworkGatewayIpConfigurationId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Network/virtualNetworkGateways/%s/ipConfigurations/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroup, id.VirtualNetworkGatewayName, id.IpConfigurationName)
}

// VirtualNetworkGatewayIpConfigurationID parses a VirtualNetworkGatewayIpConfiguration ID into an VirtualNetworkGatewayIpConfigurationId struct
func VirtualNetworkGatewayIpConfigurationID(input string) (*VirtualNetworkGatewayIpConfigurationId, error) {
	id, err := resourceids.ParseAzureResourceID(input)
	if err != nil {
		return nil, fmt.Errorf("parsing %q as an VirtualNetworkGatewayIpConfiguration ID: %+v", input, err)
	}

	resourceId := VirtualNetworkGatewayIpConfigurationId{
		SubscriptionId: id.SubscriptionID,
		ResourceGroup:  id.ResourceGroup,
	}

	if resourceId.SubscriptionId == "" {
		return nil, fmt.Errorf("ID was missing the 'subscriptions' element")
	}

	if resourceId.ResourceGroup == "" {
		return nil, fmt.Errorf("ID was missing the 'resourceGroups' element")
	}

	if resourceId.VirtualNetworkGatewayName, err = id.PopSegment("virtualNetworkGateways"); err != nil {
		return nil, err
	}
	if resourceId.IpConfigurationName, err = id.PopSegment("ipConfigurations"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &resourceId, nil
}

// VirtualNetworkGatewayIpConfigurationIDInsensitively parses an VirtualNetworkGatewayIpConfiguration ID into an VirtualNetworkGatewayIpConfigurationId struct, insensitively
// This should only be used to parse an ID for rewriting, the VirtualNetworkGatewayIpConfigurationID
// method should be used instead for validation etc.
//
// Whilst this may seem strange, this enables Terraform have consistent casing
// which works around issues in Core, whilst handling broken API responses.
func VirtualNetworkGatewayIpConfigurationIDInsensitively(input string) (*VirtualNetworkGatewayIpConfigurationId, error) {
	id, err := resourceids.ParseAzureResourceID(input)
	if err != nil {
		return nil, err
	}

	resourceId := VirtualNetworkGatewayIpConfigurationId{
		SubscriptionId: id.SubscriptionID,
		ResourceGroup:  id.ResourceGroup,
	}

	if resourceId.SubscriptionId == "" {
		return nil, fmt.Errorf("ID was missing the 'subscriptions' element")
	}

	if resourceId.ResourceGroup == "" {
		return nil, fmt.Errorf("ID was missing the 'resourceGroups' element")
	}

	// find the correct casing for the 'virtualNetworkGateways' segment
	virtualNetworkGatewaysKey := "virtualNetworkGateways"
	for key := range id.Path {
		if strings.EqualFold(key, virtualNetworkGatewaysKey) {
			virtualNetworkGatewaysKey = key
			break
		}
	}
	if resourceId.VirtualNetworkGatewayName, err = id.PopSegment(virtualNetworkGatewaysKey); err != nil {
		return nil, err
	}

	// find the correct casing for the 'ipConfigurations' segment
	ipConfigurationsKey := "ipConfigurations"
	for key := range id.Path {
		if strings.EqualFold(key, ipConfigurationsKey) {
			ipConfigurationsKey = key
			break
		}
	}
	if resourceId.IpConfigurationName, err = id.PopSegment(ipConfigurationsKey); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &resourceId, nil
}
