package parse

// NOTE: this file is generated via 'go:generate' - manual changes will be overwritten

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

type AnalyticsFirewallRuleId struct {
	SubscriptionId   string
	ResourceGroup    string
	AccountName      string
	FirewallRuleName string
}

func NewAnalyticsFirewallRuleID(subscriptionId, resourceGroup, accountName, firewallRuleName string) AnalyticsFirewallRuleId {
	return AnalyticsFirewallRuleId{
		SubscriptionId:   subscriptionId,
		ResourceGroup:    resourceGroup,
		AccountName:      accountName,
		FirewallRuleName: firewallRuleName,
	}
}

func (id AnalyticsFirewallRuleId) String() string {
	segments := []string{
		fmt.Sprintf("Firewall Rule Name %q", id.FirewallRuleName),
		fmt.Sprintf("Account Name %q", id.AccountName),
		fmt.Sprintf("Resource Group %q", id.ResourceGroup),
	}
	segmentsStr := strings.Join(segments, " / ")
	return fmt.Sprintf("%s: (%s)", "Analytics Firewall Rule", segmentsStr)
}

func (id AnalyticsFirewallRuleId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.DataLakeAnalytics/accounts/%s/firewallRules/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroup, id.AccountName, id.FirewallRuleName)
}

// AnalyticsFirewallRuleID parses a AnalyticsFirewallRule ID into an AnalyticsFirewallRuleId struct
func AnalyticsFirewallRuleID(input string) (*AnalyticsFirewallRuleId, error) {
	id, err := resourceids.ParseAzureResourceID(input)
	if err != nil {
		return nil, err
	}

	resourceId := AnalyticsFirewallRuleId{
		SubscriptionId: id.SubscriptionID,
		ResourceGroup:  id.ResourceGroup,
	}

	if resourceId.SubscriptionId == "" {
		return nil, fmt.Errorf("ID was missing the 'subscriptions' element")
	}

	if resourceId.ResourceGroup == "" {
		return nil, fmt.Errorf("ID was missing the 'resourceGroups' element")
	}

	if resourceId.AccountName, err = id.PopSegment("accounts"); err != nil {
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
