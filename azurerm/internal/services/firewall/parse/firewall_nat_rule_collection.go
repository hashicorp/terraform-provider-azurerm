package parse

// NOTE: this file is generated via 'go:generate' - manual changes will be overwritten

import (
	"fmt"
	"strings"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
)

type FirewallNatRuleCollectionId struct {
	SubscriptionId        string
	ResourceGroup         string
	AzureFirewallName     string
	NatRuleCollectionName string
}

func NewFirewallNatRuleCollectionID(subscriptionId, resourceGroup, azureFirewallName, natRuleCollectionName string) FirewallNatRuleCollectionId {
	return FirewallNatRuleCollectionId{
		SubscriptionId:        subscriptionId,
		ResourceGroup:         resourceGroup,
		AzureFirewallName:     azureFirewallName,
		NatRuleCollectionName: natRuleCollectionName,
	}
}

func (id FirewallNatRuleCollectionId) String() string {
	segments := []string{
		fmt.Sprintf("Nat Rule Collection Name %q", id.NatRuleCollectionName),
		fmt.Sprintf("Azure Firewall Name %q", id.AzureFirewallName),
		fmt.Sprintf("Resource Group %q", id.ResourceGroup),
	}
	segmentsStr := strings.Join(segments, " / ")
	return fmt.Sprintf("%s: (%s)", "Firewall Nat Rule Collection", segmentsStr)
}

func (id FirewallNatRuleCollectionId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Network/azureFirewalls/%s/natRuleCollections/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroup, id.AzureFirewallName, id.NatRuleCollectionName)
}

// FirewallNatRuleCollectionID parses a FirewallNatRuleCollection ID into an FirewallNatRuleCollectionId struct
func FirewallNatRuleCollectionID(input string) (*FirewallNatRuleCollectionId, error) {
	id, err := azure.ParseAzureResourceID(input)
	if err != nil {
		return nil, err
	}

	resourceId := FirewallNatRuleCollectionId{
		SubscriptionId: id.SubscriptionID,
		ResourceGroup:  id.ResourceGroup,
	}

	if resourceId.SubscriptionId == "" {
		return nil, fmt.Errorf("ID was missing the 'subscriptions' element")
	}

	if resourceId.ResourceGroup == "" {
		return nil, fmt.Errorf("ID was missing the 'resourceGroups' element")
	}

	if resourceId.AzureFirewallName, err = id.PopSegment("azureFirewalls"); err != nil {
		return nil, err
	}
	if resourceId.NatRuleCollectionName, err = id.PopSegment("natRuleCollections"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &resourceId, nil
}
