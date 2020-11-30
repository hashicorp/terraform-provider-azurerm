package parse

// NOTE: this file is generated via 'go:generate' - manual changes will be overwritten

import (
	"fmt"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
)

type NamespaceNetworkRuleSetId struct {
	SubscriptionId     string
	ResourceGroup      string
	NamespaceName      string
	NetworkrulesetName string
}

func NewNamespaceNetworkRuleSetID(subscriptionId, resourceGroup, namespaceName, networkrulesetName string) NamespaceNetworkRuleSetId {
	return NamespaceNetworkRuleSetId{
		SubscriptionId:     subscriptionId,
		ResourceGroup:      resourceGroup,
		NamespaceName:      namespaceName,
		NetworkrulesetName: networkrulesetName,
	}
}

func (id NamespaceNetworkRuleSetId) ID(_ string) string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.ServiceBus/namespaces/%s/networkrulesets/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroup, id.NamespaceName, id.NetworkrulesetName)
}

// NamespaceNetworkRuleSetID parses a NamespaceNetworkRuleSet ID into an NamespaceNetworkRuleSetId struct
func NamespaceNetworkRuleSetID(input string) (*NamespaceNetworkRuleSetId, error) {
	id, err := azure.ParseAzureResourceID(input)
	if err != nil {
		return nil, err
	}

	resourceId := NamespaceNetworkRuleSetId{
		SubscriptionId: id.SubscriptionID,
		ResourceGroup:  id.ResourceGroup,
	}

	if resourceId.NamespaceName, err = id.PopSegment("namespaces"); err != nil {
		return nil, err
	}
	if resourceId.NetworkrulesetName, err = id.PopSegment("networkrulesets"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &resourceId, nil
}
