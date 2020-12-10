package parse

// NOTE: this file is generated via 'go:generate' - manual changes will be overwritten

import (
	"fmt"
	"strings"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
)

type NotificationHubId struct {
	SubscriptionId string
	ResourceGroup  string
	NamespaceName  string
	Name           string
}

func NewNotificationHubID(subscriptionId, resourceGroup, namespaceName, name string) NotificationHubId {
	return NotificationHubId{
		SubscriptionId: subscriptionId,
		ResourceGroup:  resourceGroup,
		NamespaceName:  namespaceName,
		Name:           name,
	}
}

func (id NotificationHubId) String() string {
	segments := []string{
		fmt.Sprintf("Resource Group %q", id.ResourceGroup),
		fmt.Sprintf("Namespace Name %q", id.NamespaceName),
		fmt.Sprintf("Name %q", id.Name),
	}
	return strings.Join(segments, " / ")
}

func (id NotificationHubId) ID(_ string) string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.NotificationHubs/namespaces/%s/notificationHubs/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroup, id.NamespaceName, id.Name)
}

// NotificationHubID parses a NotificationHub ID into an NotificationHubId struct
func NotificationHubID(input string) (*NotificationHubId, error) {
	id, err := azure.ParseAzureResourceID(input)
	if err != nil {
		return nil, err
	}

	resourceId := NotificationHubId{
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
	if resourceId.Name, err = id.PopSegment("notificationHubs"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &resourceId, nil
}
