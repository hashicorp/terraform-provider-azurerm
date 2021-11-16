package parse

// NOTE: this file is generated via 'go:generate' - manual changes will be overwritten

import (
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
)

type MariaDBFirewallRuleId struct {
	SubscriptionId   string
	ResourceGroup    string
	ServerName       string
	FirewallRuleName string
}

func NewMariaDBFirewallRuleID(subscriptionId, resourceGroup, serverName, firewallRuleName string) MariaDBFirewallRuleId {
	return MariaDBFirewallRuleId{
		SubscriptionId:   subscriptionId,
		ResourceGroup:    resourceGroup,
		ServerName:       serverName,
		FirewallRuleName: firewallRuleName,
	}
}

func (id MariaDBFirewallRuleId) String() string {
	segments := []string{
		fmt.Sprintf("Firewall Rule Name %q", id.FirewallRuleName),
		fmt.Sprintf("Server Name %q", id.ServerName),
		fmt.Sprintf("Resource Group %q", id.ResourceGroup),
	}
	segmentsStr := strings.Join(segments, " / ")
	return fmt.Sprintf("%s: (%s)", "Maria D B Firewall Rule", segmentsStr)
}

func (id MariaDBFirewallRuleId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.DBforMariaDB/servers/%s/firewallRules/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroup, id.ServerName, id.FirewallRuleName)
}

// MariaDBFirewallRuleID parses a MariaDBFirewallRule ID into an MariaDBFirewallRuleId struct
func MariaDBFirewallRuleID(input string) (*MariaDBFirewallRuleId, error) {
	id, err := azure.ParseAzureResourceID(input)
	if err != nil {
		return nil, err
	}

	resourceId := MariaDBFirewallRuleId{
		SubscriptionId: id.SubscriptionID,
		ResourceGroup:  id.ResourceGroup,
	}

	if resourceId.SubscriptionId == "" {
		return nil, fmt.Errorf("ID was missing the 'subscriptions' element")
	}

	if resourceId.ResourceGroup == "" {
		return nil, fmt.Errorf("ID was missing the 'resourceGroups' element")
	}

	if resourceId.ServerName, err = id.PopSegment("servers"); err != nil {
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
