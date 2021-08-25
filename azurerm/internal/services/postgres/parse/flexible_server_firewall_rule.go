package parse

// NOTE: this file is generated via 'go:generate' - manual changes will be overwritten

import (
	"fmt"
	"strings"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
)

type FlexibleServerFirewallRuleId struct {
	SubscriptionId     string
	ResourceGroup      string
	FlexibleServerName string
	FirewallRuleName   string
}

func NewFlexibleServerFirewallRuleID(subscriptionId, resourceGroup, flexibleServerName, firewallRuleName string) FlexibleServerFirewallRuleId {
	return FlexibleServerFirewallRuleId{
		SubscriptionId:     subscriptionId,
		ResourceGroup:      resourceGroup,
		FlexibleServerName: flexibleServerName,
		FirewallRuleName:   firewallRuleName,
	}
}

func (id FlexibleServerFirewallRuleId) String() string {
	segments := []string{
		fmt.Sprintf("Firewall Rule Name %q", id.FirewallRuleName),
		fmt.Sprintf("Flexible Server Name %q", id.FlexibleServerName),
		fmt.Sprintf("Resource Group %q", id.ResourceGroup),
	}
	segmentsStr := strings.Join(segments, " / ")
	return fmt.Sprintf("%s: (%s)", "Flexible Server Firewall Rule", segmentsStr)
}

func (id FlexibleServerFirewallRuleId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.DBforPostgreSQL/flexibleServers/%s/firewallRules/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroup, id.FlexibleServerName, id.FirewallRuleName)
}

// FlexibleServerFirewallRuleID parses a FlexibleServerFirewallRule ID into an FlexibleServerFirewallRuleId struct
func FlexibleServerFirewallRuleID(input string) (*FlexibleServerFirewallRuleId, error) {
	id, err := azure.ParseAzureResourceID(input)
	if err != nil {
		return nil, err
	}

	resourceId := FlexibleServerFirewallRuleId{
		SubscriptionId: id.SubscriptionID,
		ResourceGroup:  id.ResourceGroup,
	}

	if resourceId.SubscriptionId == "" {
		return nil, fmt.Errorf("ID was missing the 'subscriptions' element")
	}

	if resourceId.ResourceGroup == "" {
		return nil, fmt.Errorf("ID was missing the 'resourceGroups' element")
	}

	if resourceId.FlexibleServerName, err = id.PopSegment("flexibleServers"); err != nil {
		return nil, err
	}
	if resourceId.FirewallRuleName, err = id.PopSegment("firewallRules"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &resourceId, nil
}
