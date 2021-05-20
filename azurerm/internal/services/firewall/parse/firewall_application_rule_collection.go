package parse

// NOTE: this file is generated via 'go:generate' - manual changes will be overwritten

import (
	"fmt"
	"strings"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
)

type FirewallApplicationRuleCollectionId struct {
	SubscriptionId                string
	ResourceGroup                 string
	AzureFirewallName             string
	ApplicationRuleCollectionName string
}

func NewFirewallApplicationRuleCollectionID(subscriptionId, resourceGroup, azureFirewallName, applicationRuleCollectionName string) FirewallApplicationRuleCollectionId {
	return FirewallApplicationRuleCollectionId{
		SubscriptionId:                subscriptionId,
		ResourceGroup:                 resourceGroup,
		AzureFirewallName:             azureFirewallName,
		ApplicationRuleCollectionName: applicationRuleCollectionName,
	}
}

func (id FirewallApplicationRuleCollectionId) String() string {
	segments := []string{
		fmt.Sprintf("Application Rule Collection Name %q", id.ApplicationRuleCollectionName),
		fmt.Sprintf("Azure Firewall Name %q", id.AzureFirewallName),
		fmt.Sprintf("Resource Group %q", id.ResourceGroup),
	}
	segmentsStr := strings.Join(segments, " / ")
	return fmt.Sprintf("%s: (%s)", "Firewall Application Rule Collection", segmentsStr)
}

func (id FirewallApplicationRuleCollectionId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Network/azureFirewalls/%s/applicationRuleCollections/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroup, id.AzureFirewallName, id.ApplicationRuleCollectionName)
}

// FirewallApplicationRuleCollectionID parses a FirewallApplicationRuleCollection ID into an FirewallApplicationRuleCollectionId struct
func FirewallApplicationRuleCollectionID(input string) (*FirewallApplicationRuleCollectionId, error) {
	id, err := azure.ParseAzureResourceID(input)
	if err != nil {
		return nil, err
	}

	resourceId := FirewallApplicationRuleCollectionId{
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
	if resourceId.ApplicationRuleCollectionName, err = id.PopSegment("applicationRuleCollections"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &resourceId, nil
}
