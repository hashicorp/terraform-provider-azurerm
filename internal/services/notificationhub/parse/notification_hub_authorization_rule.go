package parse

// NOTE: this file is generated via 'go:generate' - manual changes will be overwritten

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
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

func (id NotificationHubAuthorizationRuleId) String() string {
	segments := []string{
		fmt.Sprintf("Authorization Rule Name %q", id.AuthorizationRuleName),
		fmt.Sprintf("Notification Hub Name %q", id.NotificationHubName),
		fmt.Sprintf("Namespace Name %q", id.NamespaceName),
		fmt.Sprintf("Resource Group %q", id.ResourceGroup),
	}
	segmentsStr := strings.Join(segments, " / ")
	return fmt.Sprintf("%s: (%s)", "Notification Hub Authorization Rule", segmentsStr)
}

func (id NotificationHubAuthorizationRuleId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.NotificationHubs/namespaces/%s/notificationHubs/%s/authorizationRules/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroup, id.NamespaceName, id.NotificationHubName, id.AuthorizationRuleName)
}

// NotificationHubAuthorizationRuleID parses a NotificationHubAuthorizationRule ID into an NotificationHubAuthorizationRuleId struct
func NotificationHubAuthorizationRuleID(input string) (*NotificationHubAuthorizationRuleId, error) {
	id, err := resourceids.ParseAzureResourceID(input)
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
	if resourceId.AuthorizationRuleName, err = id.PopSegment("authorizationRules"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &resourceId, nil
}

// NotificationHubAuthorizationRuleIDInsensitively parses an NotificationHubAuthorizationRule ID into an NotificationHubAuthorizationRuleId struct, insensitively
// This should only be used to parse an ID for rewriting, the NotificationHubAuthorizationRuleID
// method should be used instead for validation etc.
//
// Whilst this may seem strange, this enables Terraform have consistent casing
// which works around issues in Core, whilst handling broken API responses.
func NotificationHubAuthorizationRuleIDInsensitively(input string) (*NotificationHubAuthorizationRuleId, error) {
	id, err := resourceids.ParseAzureResourceID(input)
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

	// find the correct casing for the 'namespaces' segment
	namespacesKey := "namespaces"
	for key := range id.Path {
		if strings.EqualFold(key, namespacesKey) {
			namespacesKey = key
			break
		}
	}
	if resourceId.NamespaceName, err = id.PopSegment(namespacesKey); err != nil {
		return nil, err
	}

	// find the correct casing for the 'notificationHubs' segment
	notificationHubsKey := "notificationHubs"
	for key := range id.Path {
		if strings.EqualFold(key, notificationHubsKey) {
			notificationHubsKey = key
			break
		}
	}
	if resourceId.NotificationHubName, err = id.PopSegment(notificationHubsKey); err != nil {
		return nil, err
	}

	// find the correct casing for the 'authorizationRules' segment
	authorizationRulesKey := "authorizationRules"
	for key := range id.Path {
		if strings.EqualFold(key, authorizationRulesKey) {
			authorizationRulesKey = key
			break
		}
	}
	if resourceId.AuthorizationRuleName, err = id.PopSegment(authorizationRulesKey); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &resourceId, nil
}
