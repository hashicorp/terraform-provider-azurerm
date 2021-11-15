package parse

// NOTE: this file is generated via 'go:generate' - manual changes will be overwritten

import (
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
)

type VpnNatRuleId struct {
	SubscriptionId string
	ResourceGroup  string
	VpnGatewayName string
	NatRuleName    string
}

func NewVpnNatRuleID(subscriptionId, resourceGroup, vpnGatewayName, natRuleName string) VpnNatRuleId {
	return VpnNatRuleId{
		SubscriptionId: subscriptionId,
		ResourceGroup:  resourceGroup,
		VpnGatewayName: vpnGatewayName,
		NatRuleName:    natRuleName,
	}
}

func (id VpnNatRuleId) String() string {
	segments := []string{
		fmt.Sprintf("Nat Rule Name %q", id.NatRuleName),
		fmt.Sprintf("Vpn Gateway Name %q", id.VpnGatewayName),
		fmt.Sprintf("Resource Group %q", id.ResourceGroup),
	}
	segmentsStr := strings.Join(segments, " / ")
	return fmt.Sprintf("%s: (%s)", "Vpn Nat Rule", segmentsStr)
}

func (id VpnNatRuleId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Network/vpnGateways/%s/natRules/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroup, id.VpnGatewayName, id.NatRuleName)
}

// VpnNatRuleID parses a VpnNatRule ID into an VpnNatRuleId struct
func VpnNatRuleID(input string) (*VpnNatRuleId, error) {
	id, err := azure.ParseAzureResourceID(input)
	if err != nil {
		return nil, err
	}

	resourceId := VpnNatRuleId{
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
