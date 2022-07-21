package notificationhubs

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

var _ resourceids.ResourceId = NotificationHubAuthorizationRuleId{}

// NotificationHubAuthorizationRuleId is a struct representing the Resource ID for a Notification Hub Authorization Rule
type NotificationHubAuthorizationRuleId struct {
	SubscriptionId        string
	ResourceGroupName     string
	NamespaceName         string
	NotificationHubName   string
	AuthorizationRuleName string
}

// NewNotificationHubAuthorizationRuleID returns a new NotificationHubAuthorizationRuleId struct
func NewNotificationHubAuthorizationRuleID(subscriptionId string, resourceGroupName string, namespaceName string, notificationHubName string, authorizationRuleName string) NotificationHubAuthorizationRuleId {
	return NotificationHubAuthorizationRuleId{
		SubscriptionId:        subscriptionId,
		ResourceGroupName:     resourceGroupName,
		NamespaceName:         namespaceName,
		NotificationHubName:   notificationHubName,
		AuthorizationRuleName: authorizationRuleName,
	}
}

// ParseNotificationHubAuthorizationRuleID parses 'input' into a NotificationHubAuthorizationRuleId
func ParseNotificationHubAuthorizationRuleID(input string) (*NotificationHubAuthorizationRuleId, error) {
	parser := resourceids.NewParserFromResourceIdType(NotificationHubAuthorizationRuleId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := NotificationHubAuthorizationRuleId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, fmt.Errorf("the segment 'subscriptionId' was not found in the resource id %q", input)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, fmt.Errorf("the segment 'resourceGroupName' was not found in the resource id %q", input)
	}

	if id.NamespaceName, ok = parsed.Parsed["namespaceName"]; !ok {
		return nil, fmt.Errorf("the segment 'namespaceName' was not found in the resource id %q", input)
	}

	if id.NotificationHubName, ok = parsed.Parsed["notificationHubName"]; !ok {
		return nil, fmt.Errorf("the segment 'notificationHubName' was not found in the resource id %q", input)
	}

	if id.AuthorizationRuleName, ok = parsed.Parsed["authorizationRuleName"]; !ok {
		return nil, fmt.Errorf("the segment 'authorizationRuleName' was not found in the resource id %q", input)
	}

	return &id, nil
}

// ParseNotificationHubAuthorizationRuleIDInsensitively parses 'input' case-insensitively into a NotificationHubAuthorizationRuleId
// note: this method should only be used for API response data and not user input
func ParseNotificationHubAuthorizationRuleIDInsensitively(input string) (*NotificationHubAuthorizationRuleId, error) {
	parser := resourceids.NewParserFromResourceIdType(NotificationHubAuthorizationRuleId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := NotificationHubAuthorizationRuleId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, fmt.Errorf("the segment 'subscriptionId' was not found in the resource id %q", input)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, fmt.Errorf("the segment 'resourceGroupName' was not found in the resource id %q", input)
	}

	if id.NamespaceName, ok = parsed.Parsed["namespaceName"]; !ok {
		return nil, fmt.Errorf("the segment 'namespaceName' was not found in the resource id %q", input)
	}

	if id.NotificationHubName, ok = parsed.Parsed["notificationHubName"]; !ok {
		return nil, fmt.Errorf("the segment 'notificationHubName' was not found in the resource id %q", input)
	}

	if id.AuthorizationRuleName, ok = parsed.Parsed["authorizationRuleName"]; !ok {
		return nil, fmt.Errorf("the segment 'authorizationRuleName' was not found in the resource id %q", input)
	}

	return &id, nil
}

// ValidateNotificationHubAuthorizationRuleID checks that 'input' can be parsed as a Notification Hub Authorization Rule ID
func ValidateNotificationHubAuthorizationRuleID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseNotificationHubAuthorizationRuleID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Notification Hub Authorization Rule ID
func (id NotificationHubAuthorizationRuleId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.NotificationHubs/namespaces/%s/notificationHubs/%s/authorizationRules/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.NamespaceName, id.NotificationHubName, id.AuthorizationRuleName)
}

// Segments returns a slice of Resource ID Segments which comprise this Notification Hub Authorization Rule ID
func (id NotificationHubAuthorizationRuleId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftNotificationHubs", "Microsoft.NotificationHubs", "Microsoft.NotificationHubs"),
		resourceids.StaticSegment("staticNamespaces", "namespaces", "namespaces"),
		resourceids.UserSpecifiedSegment("namespaceName", "namespaceValue"),
		resourceids.StaticSegment("staticNotificationHubs", "notificationHubs", "notificationHubs"),
		resourceids.UserSpecifiedSegment("notificationHubName", "notificationHubValue"),
		resourceids.StaticSegment("staticAuthorizationRules", "authorizationRules", "authorizationRules"),
		resourceids.UserSpecifiedSegment("authorizationRuleName", "authorizationRuleValue"),
	}
}

// String returns a human-readable description of this Notification Hub Authorization Rule ID
func (id NotificationHubAuthorizationRuleId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Namespace Name: %q", id.NamespaceName),
		fmt.Sprintf("Notification Hub Name: %q", id.NotificationHubName),
		fmt.Sprintf("Authorization Rule Name: %q", id.AuthorizationRuleName),
	}
	return fmt.Sprintf("Notification Hub Authorization Rule (%s)", strings.Join(components, "\n"))
}
