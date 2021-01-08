package parse

// NOTE: this file is generated via 'go:generate' - manual changes will be overwritten

import (
	"fmt"
	"strings"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
)

type FirewallNatRuleCollectionId struct {
	SubscriptionId            string
	ResourceGroup             string
	AzureFirewallName         string
	NetworkRuleCollectionName string
}

func NewFirewallNatRuleCollectionID(subscriptionId, resourceGroup, azureFirewallName, networkRuleCollectionName string) FirewallNatRuleCollectionId {
	return FirewallNatRuleCollectionId{
		SubscriptionId:            subscriptionId,
		ResourceGroup:             resourceGroup,
		AzureFirewallName:         azureFirewallName,
		NetworkRuleCollectionName: networkRuleCollectionName,
	}
}

func (id FirewallNatRuleCollectionId) String() string {
	segments := []string{
		fmt.Sprintf("Network Rule Collection Name %q", id.NetworkRuleCollectionName),
		fmt.Sprintf("Azure Firewall Name %q", id.AzureFirewallName),
		fmt.Sprintf("Resource Group %q", id.ResourceGroup),
	}
	segmentsStr := strings.Join(segments, " / ")
	return fmt.Sprintf("%s: (%s)", "Firewall Nat Rule Collection", segmentsStr)
}

func (id FirewallNatRuleCollectionId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Network/azureFirewalls/%s/networkRuleCollections/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroup, id.AzureFirewallName, id.NetworkRuleCollectionName)
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
	if resourceId.NetworkRuleCollectionName, err = id.PopSegment("networkRuleCollections"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &resourceId, nil
}
