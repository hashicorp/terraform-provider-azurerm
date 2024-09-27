// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package parse

// NOTE: this file is generated via 'go:generate' - manual changes will be overwritten

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

type VirtualNetworkGatewayPolicyGroupId struct {
	SubscriptionId            string
	ResourceGroup             string
	VirtualNetworkGatewayName string
	Name                      string
}

func NewVirtualNetworkGatewayPolicyGroupID(subscriptionId, resourceGroup, virtualNetworkGatewayName, name string) VirtualNetworkGatewayPolicyGroupId {
	return VirtualNetworkGatewayPolicyGroupId{
		SubscriptionId:            subscriptionId,
		ResourceGroup:             resourceGroup,
		VirtualNetworkGatewayName: virtualNetworkGatewayName,
		Name:                      name,
	}
}

func (id VirtualNetworkGatewayPolicyGroupId) String() string {
	segments := []string{
		fmt.Sprintf("Name %q", id.Name),
		fmt.Sprintf("Virtual Network Gateway Name %q", id.VirtualNetworkGatewayName),
		fmt.Sprintf("Resource Group %q", id.ResourceGroup),
	}
	segmentsStr := strings.Join(segments, " / ")
	return fmt.Sprintf("%s: (%s)", "Virtual Network Gateway Policy Group", segmentsStr)
}

func (id VirtualNetworkGatewayPolicyGroupId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Network/virtualNetworkGateways/%s/virtualNetworkGatewayPolicyGroups/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroup, id.VirtualNetworkGatewayName, id.Name)
}

// VirtualNetworkGatewayPolicyGroupID parses a VirtualNetworkGatewayPolicyGroup ID into an VirtualNetworkGatewayPolicyGroupId struct
func VirtualNetworkGatewayPolicyGroupID(input string) (*VirtualNetworkGatewayPolicyGroupId, error) {
	id, err := resourceids.ParseAzureResourceID(input)
	if err != nil {
		return nil, fmt.Errorf("parsing %q as an VirtualNetworkGatewayPolicyGroup ID: %+v", input, err)
	}

	resourceId := VirtualNetworkGatewayPolicyGroupId{
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
	if resourceId.Name, err = id.PopSegment("virtualNetworkGatewayPolicyGroups"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &resourceId, nil
}
