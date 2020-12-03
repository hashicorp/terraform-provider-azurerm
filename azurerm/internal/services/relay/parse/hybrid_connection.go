package parse

// NOTE: this file is generated via 'go:generate' - manual changes will be overwritten

import (
	"fmt"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
)

type HybridConnectionId struct {
	SubscriptionId string
	ResourceGroup  string
	NamespaceName  string
	Name           string
}

func NewHybridConnectionID(subscriptionId, resourceGroup, namespaceName, name string) HybridConnectionId {
	return HybridConnectionId{
		SubscriptionId: subscriptionId,
		ResourceGroup:  resourceGroup,
		NamespaceName:  namespaceName,
		Name:           name,
	}
}

func (id HybridConnectionId) ID(_ string) string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Relay/namespaces/%s/hybridConnections/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroup, id.NamespaceName, id.Name)
}

// HybridConnectionID parses a HybridConnection ID into an HybridConnectionId struct
func HybridConnectionID(input string) (*HybridConnectionId, error) {
	id, err := azure.ParseAzureResourceID(input)
	if err != nil {
		return nil, err
	}

	resourceId := HybridConnectionId{
		SubscriptionId: id.SubscriptionID,
		ResourceGroup:  id.ResourceGroup,
	}

	if resourceId.SubscriptionId == "" {
		return nil, fmt.Errorf("ID was missing the 'subscriptions' element")
	}

	if resourceId.ResourceGroup == "" {
		return nil, fmt.Errorf("ID was missing the 'resourceGroups' element")
	}

	if resourceId.NamespaceName, err = id.PopSegment("namespaces"); err != nil {
		return nil, err
	}
	if resourceId.Name, err = id.PopSegment("hybridConnections"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &resourceId, nil
}
