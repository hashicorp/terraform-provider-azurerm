// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package parse

// NOTE: this file is generated via 'go:generate' - manual changes will be overwritten

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

type VirtualNetworkGatewayNatRuleId struct {
	SubscriptionId            string
	ResourceGroup             string
	VirtualNetworkGatewayName string
	NatRuleName               string
}

func NewVirtualNetworkGatewayNatRuleID(subscriptionId, resourceGroup, virtualNetworkGatewayName, natRuleName string) VirtualNetworkGatewayNatRuleId {
	return VirtualNetworkGatewayNatRuleId{
		SubscriptionId:            subscriptionId,
		ResourceGroup:             resourceGroup,
		VirtualNetworkGatewayName: virtualNetworkGatewayName,
		NatRuleName:               natRuleName,
	}
}

func (id VirtualNetworkGatewayNatRuleId) String() string {
	segments := []string{
		fmt.Sprintf("Nat Rule Name %q", id.NatRuleName),
		fmt.Sprintf("Virtual Network Gateway Name %q", id.VirtualNetworkGatewayName),
		fmt.Sprintf("Resource Group %q", id.ResourceGroup),
	}
	segmentsStr := strings.Join(segments, " / ")
	return fmt.Sprintf("%s: (%s)", "Virtual Network Gateway Nat Rule", segmentsStr)
}

func (id VirtualNetworkGatewayNatRuleId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Network/virtualNetworkGateways/%s/natRules/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroup, id.VirtualNetworkGatewayName, id.NatRuleName)
}

// VirtualNetworkGatewayNatRuleID parses a VirtualNetworkGatewayNatRule ID into an VirtualNetworkGatewayNatRuleId struct
func VirtualNetworkGatewayNatRuleID(input string) (*VirtualNetworkGatewayNatRuleId, error) {
	id, err := resourceids.ParseAzureResourceID(input)
	if err != nil {
		return nil, fmt.Errorf("parsing %q as an VirtualNetworkGatewayNatRule ID: %+v", input, err)
	}

	resourceId := VirtualNetworkGatewayNatRuleId{
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
	if resourceId.NatRuleName, err = id.PopSegment("natRules"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &resourceId, nil
}
