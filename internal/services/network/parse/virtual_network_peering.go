// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package parse

// NOTE: this file is generated via 'go:generate' - manual changes will be overwritten

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

type VirtualNetworkPeeringId struct {
	SubscriptionId     string
	ResourceGroup      string
	VirtualNetworkName string
	Name               string
}

func NewVirtualNetworkPeeringID(subscriptionId, resourceGroup, virtualNetworkName, name string) VirtualNetworkPeeringId {
	return VirtualNetworkPeeringId{
		SubscriptionId:     subscriptionId,
		ResourceGroup:      resourceGroup,
		VirtualNetworkName: virtualNetworkName,
		Name:               name,
	}
}

func (id VirtualNetworkPeeringId) String() string {
	segments := []string{
		fmt.Sprintf("Name %q", id.Name),
		fmt.Sprintf("Virtual Network Name %q", id.VirtualNetworkName),
		fmt.Sprintf("Resource Group %q", id.ResourceGroup),
	}
	segmentsStr := strings.Join(segments, " / ")
	return fmt.Sprintf("%s: (%s)", "Virtual Network Peering", segmentsStr)
}

func (id VirtualNetworkPeeringId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Network/virtualNetworks/%s/virtualNetworkPeerings/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroup, id.VirtualNetworkName, id.Name)
}

// VirtualNetworkPeeringID parses a VirtualNetworkPeering ID into an VirtualNetworkPeeringId struct
func VirtualNetworkPeeringID(input string) (*VirtualNetworkPeeringId, error) {
	id, err := resourceids.ParseAzureResourceID(input)
	if err != nil {
		return nil, fmt.Errorf("parsing %q as an VirtualNetworkPeering ID: %+v", input, err)
	}

	resourceId := VirtualNetworkPeeringId{
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
	if resourceId.Name, err = id.PopSegment("virtualNetworkPeerings"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &resourceId, nil
}
