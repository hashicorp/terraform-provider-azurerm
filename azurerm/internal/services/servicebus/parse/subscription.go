package parse

// NOTE: this file is generated via 'go:generate' - manual changes will be overwritten

import (
	"fmt"
	"strings"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
)

type SubscriptionId struct {
	SubscriptionId string
	ResourceGroup  string
	NamespaceName  string
	TopicName      string
	Name           string
}

func NewSubscriptionID(subscriptionId, resourceGroup, namespaceName, topicName, name string) SubscriptionId {
	return SubscriptionId{
		SubscriptionId: subscriptionId,
		ResourceGroup:  resourceGroup,
		NamespaceName:  namespaceName,
		TopicName:      topicName,
		Name:           name,
	}
}

func (id SubscriptionId) String() string {
	segments := []string{
		fmt.Sprintf("Resource Group %q", id.ResourceGroup),
		fmt.Sprintf("Namespace Name %q", id.NamespaceName),
		fmt.Sprintf("Topic Name %q", id.TopicName),
		fmt.Sprintf("Name %q", id.Name),
	}
	return strings.Join(segments, " / ")
}

func (id SubscriptionId) ID(_ string) string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.ServiceBus/namespaces/%s/topics/%s/subscriptions/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroup, id.NamespaceName, id.TopicName, id.Name)
}

// SubscriptionID parses a Subscription ID into an SubscriptionId struct
func SubscriptionID(input string) (*SubscriptionId, error) {
	id, err := azure.ParseAzureResourceID(input)
	if err != nil {
		return nil, err
	}

	resourceId := SubscriptionId{
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
	if resourceId.TopicName, err = id.PopSegment("topics"); err != nil {
		return nil, err
	}

	if resourceId.Name == "" {
		return nil, fmt.Errorf("ID was missing the 'subscriptions' element")
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &resourceId, nil
}
