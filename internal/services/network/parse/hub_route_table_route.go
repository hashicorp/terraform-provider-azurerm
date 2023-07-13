// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package parse

// NOTE: this file is generated via 'go:generate' - manual changes will be overwritten

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

type HubRouteTableRouteId struct {
	SubscriptionId    string
	ResourceGroup     string
	VirtualHubName    string
	HubRouteTableName string
	RouteName         string
}

func NewHubRouteTableRouteID(subscriptionId, resourceGroup, virtualHubName, hubRouteTableName, routeName string) HubRouteTableRouteId {
	return HubRouteTableRouteId{
		SubscriptionId:    subscriptionId,
		ResourceGroup:     resourceGroup,
		VirtualHubName:    virtualHubName,
		HubRouteTableName: hubRouteTableName,
		RouteName:         routeName,
	}
}

func (id HubRouteTableRouteId) String() string {
	segments := []string{
		fmt.Sprintf("Route Name %q", id.RouteName),
		fmt.Sprintf("Hub Route Table Name %q", id.HubRouteTableName),
		fmt.Sprintf("Virtual Hub Name %q", id.VirtualHubName),
		fmt.Sprintf("Resource Group %q", id.ResourceGroup),
	}
	segmentsStr := strings.Join(segments, " / ")
	return fmt.Sprintf("%s: (%s)", "Hub Route Table Route", segmentsStr)
}

func (id HubRouteTableRouteId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Network/virtualHubs/%s/hubRouteTables/%s/routes/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroup, id.VirtualHubName, id.HubRouteTableName, id.RouteName)
}

// HubRouteTableRouteID parses a HubRouteTableRoute ID into an HubRouteTableRouteId struct
func HubRouteTableRouteID(input string) (*HubRouteTableRouteId, error) {
	id, err := resourceids.ParseAzureResourceID(input)
	if err != nil {
		return nil, fmt.Errorf("parsing %q as an HubRouteTableRoute ID: %+v", input, err)
	}

	resourceId := HubRouteTableRouteId{
		SubscriptionId: id.SubscriptionID,
		ResourceGroup:  id.ResourceGroup,
	}

	if resourceId.SubscriptionId == "" {
		return nil, fmt.Errorf("ID was missing the 'subscriptions' element")
	}

	if resourceId.ResourceGroup == "" {
		return nil, fmt.Errorf("ID was missing the 'resourceGroups' element")
	}

	if resourceId.VirtualHubName, err = id.PopSegment("virtualHubs"); err != nil {
		return nil, err
	}
	if resourceId.HubRouteTableName, err = id.PopSegment("hubRouteTables"); err != nil {
		return nil, err
	}
	if resourceId.RouteName, err = id.PopSegment("routes"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &resourceId, nil
}
