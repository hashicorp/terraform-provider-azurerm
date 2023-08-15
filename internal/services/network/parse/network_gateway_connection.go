// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package parse

// NOTE: this file is generated via 'go:generate' - manual changes will be overwritten

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

type NetworkGatewayConnectionId struct {
	SubscriptionId string
	ResourceGroup  string
	ConnectionName string
}

func NewNetworkGatewayConnectionID(subscriptionId, resourceGroup, connectionName string) NetworkGatewayConnectionId {
	return NetworkGatewayConnectionId{
		SubscriptionId: subscriptionId,
		ResourceGroup:  resourceGroup,
		ConnectionName: connectionName,
	}
}

func (id NetworkGatewayConnectionId) String() string {
	segments := []string{
		fmt.Sprintf("Connection Name %q", id.ConnectionName),
		fmt.Sprintf("Resource Group %q", id.ResourceGroup),
	}
	segmentsStr := strings.Join(segments, " / ")
	return fmt.Sprintf("%s: (%s)", "Network Gateway Connection", segmentsStr)
}

func (id NetworkGatewayConnectionId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Network/connections/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroup, id.ConnectionName)
}

// NetworkGatewayConnectionID parses a NetworkGatewayConnection ID into an NetworkGatewayConnectionId struct
func NetworkGatewayConnectionID(input string) (*NetworkGatewayConnectionId, error) {
	id, err := resourceids.ParseAzureResourceID(input)
	if err != nil {
		return nil, fmt.Errorf("parsing %q as an NetworkGatewayConnection ID: %+v", input, err)
	}

	resourceId := NetworkGatewayConnectionId{
		SubscriptionId: id.SubscriptionID,
		ResourceGroup:  id.ResourceGroup,
	}

	if resourceId.SubscriptionId == "" {
		return nil, fmt.Errorf("ID was missing the 'subscriptions' element")
	}

	if resourceId.ResourceGroup == "" {
		return nil, fmt.Errorf("ID was missing the 'resourceGroups' element")
	}

	if resourceId.ConnectionName, err = id.PopSegment("connections"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &resourceId, nil
}
