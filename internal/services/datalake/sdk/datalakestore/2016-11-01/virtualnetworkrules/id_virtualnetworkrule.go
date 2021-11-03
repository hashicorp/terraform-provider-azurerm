package virtualnetworkrules

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

type VirtualNetworkRuleId struct {
	SubscriptionId string
	ResourceGroup  string
	AccountName    string
	Name           string
}

func NewVirtualNetworkRuleID(subscriptionId, resourceGroup, accountName, name string) VirtualNetworkRuleId {
	return VirtualNetworkRuleId{
		SubscriptionId: subscriptionId,
		ResourceGroup:  resourceGroup,
		AccountName:    accountName,
		Name:           name,
	}
}

func (id VirtualNetworkRuleId) String() string {
	segments := []string{
		fmt.Sprintf("Name %q", id.Name),
		fmt.Sprintf("Account Name %q", id.AccountName),
		fmt.Sprintf("Resource Group %q", id.ResourceGroup),
	}
	segmentsStr := strings.Join(segments, " / ")
	return fmt.Sprintf("%s: (%s)", "Virtual Network Rule", segmentsStr)
}

func (id VirtualNetworkRuleId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.DataLakeStore/accounts/%s/virtualNetworkRules/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroup, id.AccountName, id.Name)
}

// ParseVirtualNetworkRuleID parses a VirtualNetworkRule ID into an VirtualNetworkRuleId struct
func ParseVirtualNetworkRuleID(input string) (*VirtualNetworkRuleId, error) {
	id, err := resourceids.ParseAzureResourceID(input)
	if err != nil {
		return nil, err
	}

	resourceId := VirtualNetworkRuleId{
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
	if resourceId.Name, err = id.PopSegment("virtualNetworkRules"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &resourceId, nil
}

// ParseVirtualNetworkRuleIDInsensitively parses an VirtualNetworkRule ID into an VirtualNetworkRuleId struct, insensitively
// This should only be used to parse an ID for rewriting to a consistent casing,
// the ParseVirtualNetworkRuleID method should be used instead for validation etc.
func ParseVirtualNetworkRuleIDInsensitively(input string) (*VirtualNetworkRuleId, error) {
	id, err := resourceids.ParseAzureResourceID(input)
	if err != nil {
		return nil, err
	}

	resourceId := VirtualNetworkRuleId{
		SubscriptionId: id.SubscriptionID,
		ResourceGroup:  id.ResourceGroup,
	}

	if resourceId.SubscriptionId == "" {
		return nil, fmt.Errorf("ID was missing the 'subscriptions' element")
	}

	if resourceId.ResourceGroup == "" {
		return nil, fmt.Errorf("ID was missing the 'resourceGroups' element")
	}

	// find the correct casing for the 'accounts' segment
	accountsKey := "accounts"
	for key := range id.Path {
		if strings.EqualFold(key, accountsKey) {
			accountsKey = key
			break
		}
	}
	if resourceId.AccountName, err = id.PopSegment(accountsKey); err != nil {
		return nil, err
	}

	// find the correct casing for the 'virtualNetworkRules' segment
	virtualNetworkRulesKey := "virtualNetworkRules"
	for key := range id.Path {
		if strings.EqualFold(key, virtualNetworkRulesKey) {
			virtualNetworkRulesKey = key
			break
		}
	}
	if resourceId.Name, err = id.PopSegment(virtualNetworkRulesKey); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &resourceId, nil
}
