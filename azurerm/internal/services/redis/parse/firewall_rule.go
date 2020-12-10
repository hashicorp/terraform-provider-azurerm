package parse

// NOTE: this file is generated via 'go:generate' - manual changes will be overwritten

import (
	"fmt"
	"strings"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
)

type FirewallRuleId struct {
	SubscriptionId string
	ResourceGroup  string
	RediName       string
	Name           string
}

func NewFirewallRuleID(subscriptionId, resourceGroup, rediName, name string) FirewallRuleId {
	return FirewallRuleId{
		SubscriptionId: subscriptionId,
		ResourceGroup:  resourceGroup,
		RediName:       rediName,
		Name:           name,
	}
}

func (id FirewallRuleId) String() string {
	segments := []string{
		fmt.Sprintf("Resource Group %q", id.ResourceGroup),
		fmt.Sprintf("Redi Name %q", id.RediName),
		fmt.Sprintf("Name %q", id.Name),
	}
	return strings.Join(segments, " / ")
}

func (id FirewallRuleId) ID(_ string) string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Cache/Redis/%s/firewallRules/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroup, id.RediName, id.Name)
}

// FirewallRuleID parses a FirewallRule ID into an FirewallRuleId struct
func FirewallRuleID(input string) (*FirewallRuleId, error) {
	id, err := azure.ParseAzureResourceID(input)
	if err != nil {
		return nil, err
	}

	resourceId := FirewallRuleId{
		SubscriptionId: id.SubscriptionID,
		ResourceGroup:  id.ResourceGroup,
	}

	if resourceId.SubscriptionId == "" {
		return nil, fmt.Errorf("ID was missing the 'subscriptions' element")
	}

	if resourceId.ResourceGroup == "" {
		return nil, fmt.Errorf("ID was missing the 'resourceGroups' element")
	}

	if resourceId.RediName, err = id.PopSegment("Redis"); err != nil {
		return nil, err
	}
	if resourceId.Name, err = id.PopSegment("firewallRules"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &resourceId, nil
}
