package parse

// NOTE: this file is generated via 'go:generate' - manual changes will be overwritten

import (
	"fmt"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
)

type NotificationHubAuthorizationRuleId struct {
	SubscriptionId        string
	ResourceGroup         string
	NamespaceName         string
	NotificationHubName   string
	AuthorizationRuleName string
}

func NewNotificationHubAuthorizationRuleID(subscriptionId, resourceGroup, namespaceName, notificationHubName, authorizationRuleName string) NotificationHubAuthorizationRuleId {
	return NotificationHubAuthorizationRuleId{
		SubscriptionId:        subscriptionId,
		ResourceGroup:         resourceGroup,
		NamespaceName:         namespaceName,
		NotificationHubName:   notificationHubName,
		AuthorizationRuleName: authorizationRuleName,
	}
}

func (id NotificationHubAuthorizationRuleId) ID(_ string) string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.NotificationHubs/namespaces/%s/notificationHubs/%s/AuthorizationRules/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroup, id.NamespaceName, id.NotificationHubName, id.AuthorizationRuleName)
}

// NotificationHubAuthorizationRuleID parses a NotificationHubAuthorizationRule ID into an NotificationHubAuthorizationRuleId struct
func NotificationHubAuthorizationRuleID(input string) (*NotificationHubAuthorizationRuleId, error) {
	id, err := azure.ParseAzureResourceID(input)
	if err != nil {
		return nil, err
	}

	resourceId := NotificationHubAuthorizationRuleId{
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
	if resourceId.NotificationHubName, err = id.PopSegment("notificationHubs"); err != nil {
		return nil, err
	}
	if resourceId.AuthorizationRuleName, err = id.PopSegment("AuthorizationRules"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &resourceId, nil
}
