// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package parse

// NOTE: this file is generated via 'go:generate' - manual changes will be overwritten

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

type VpnGatewayNatRuleId struct {
	SubscriptionId string
	ResourceGroup  string
	VpnGatewayName string
	NatRuleName    string
}

func NewVpnGatewayNatRuleID(subscriptionId, resourceGroup, vpnGatewayName, natRuleName string) VpnGatewayNatRuleId {
	return VpnGatewayNatRuleId{
		SubscriptionId: subscriptionId,
		ResourceGroup:  resourceGroup,
		VpnGatewayName: vpnGatewayName,
		NatRuleName:    natRuleName,
	}
}

func (id VpnGatewayNatRuleId) String() string {
	segments := []string{
		fmt.Sprintf("Nat Rule Name %q", id.NatRuleName),
		fmt.Sprintf("Vpn Gateway Name %q", id.VpnGatewayName),
		fmt.Sprintf("Resource Group %q", id.ResourceGroup),
	}
	segmentsStr := strings.Join(segments, " / ")
	return fmt.Sprintf("%s: (%s)", "Vpn Gateway Nat Rule", segmentsStr)
}

func (id VpnGatewayNatRuleId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Network/vpnGateways/%s/natRules/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroup, id.VpnGatewayName, id.NatRuleName)
}

// VpnGatewayNatRuleID parses a VpnGatewayNatRule ID into an VpnGatewayNatRuleId struct
func VpnGatewayNatRuleID(input string) (*VpnGatewayNatRuleId, error) {
	id, err := resourceids.ParseAzureResourceID(input)
	if err != nil {
		return nil, fmt.Errorf("parsing %q as an VpnGatewayNatRule ID: %+v", input, err)
	}

	resourceId := VpnGatewayNatRuleId{
		SubscriptionId: id.SubscriptionID,
		ResourceGroup:  id.ResourceGroup,
	}

	if resourceId.SubscriptionId == "" {
		return nil, fmt.Errorf("ID was missing the 'subscriptions' element")
	}

	if resourceId.ResourceGroup == "" {
		return nil, fmt.Errorf("ID was missing the 'resourceGroups' element")
	}

	if resourceId.VpnGatewayName, err = id.PopSegment("vpnGateways"); err != nil {
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
