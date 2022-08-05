package queuesauthorizationrule

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

var _ resourceids.ResourceId = QueueAuthorizationRuleId{}

// QueueAuthorizationRuleId is a struct representing the Resource ID for a Queue Authorization Rule
type QueueAuthorizationRuleId struct {
	SubscriptionId        string
	ResourceGroupName     string
	NamespaceName         string
	QueueName             string
	AuthorizationRuleName string
}

// NewQueueAuthorizationRuleID returns a new QueueAuthorizationRuleId struct
func NewQueueAuthorizationRuleID(subscriptionId string, resourceGroupName string, namespaceName string, queueName string, authorizationRuleName string) QueueAuthorizationRuleId {
	return QueueAuthorizationRuleId{
		SubscriptionId:        subscriptionId,
		ResourceGroupName:     resourceGroupName,
		NamespaceName:         namespaceName,
		QueueName:             queueName,
		AuthorizationRuleName: authorizationRuleName,
	}
}

// ParseQueueAuthorizationRuleID parses 'input' into a QueueAuthorizationRuleId
func ParseQueueAuthorizationRuleID(input string) (*QueueAuthorizationRuleId, error) {
	parser := resourceids.NewParserFromResourceIdType(QueueAuthorizationRuleId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := QueueAuthorizationRuleId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, fmt.Errorf("the segment 'subscriptionId' was not found in the resource id %q", input)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, fmt.Errorf("the segment 'resourceGroupName' was not found in the resource id %q", input)
	}

	if id.NamespaceName, ok = parsed.Parsed["namespaceName"]; !ok {
		return nil, fmt.Errorf("the segment 'namespaceName' was not found in the resource id %q", input)
	}

	if id.QueueName, ok = parsed.Parsed["queueName"]; !ok {
		return nil, fmt.Errorf("the segment 'queueName' was not found in the resource id %q", input)
	}

	if id.AuthorizationRuleName, ok = parsed.Parsed["authorizationRuleName"]; !ok {
		return nil, fmt.Errorf("the segment 'authorizationRuleName' was not found in the resource id %q", input)
	}

	return &id, nil
}

// ParseQueueAuthorizationRuleIDInsensitively parses 'input' case-insensitively into a QueueAuthorizationRuleId
// note: this method should only be used for API response data and not user input
func ParseQueueAuthorizationRuleIDInsensitively(input string) (*QueueAuthorizationRuleId, error) {
	parser := resourceids.NewParserFromResourceIdType(QueueAuthorizationRuleId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := QueueAuthorizationRuleId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, fmt.Errorf("the segment 'subscriptionId' was not found in the resource id %q", input)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, fmt.Errorf("the segment 'resourceGroupName' was not found in the resource id %q", input)
	}

	if id.NamespaceName, ok = parsed.Parsed["namespaceName"]; !ok {
		return nil, fmt.Errorf("the segment 'namespaceName' was not found in the resource id %q", input)
	}

	if id.QueueName, ok = parsed.Parsed["queueName"]; !ok {
		return nil, fmt.Errorf("the segment 'queueName' was not found in the resource id %q", input)
	}

	if id.AuthorizationRuleName, ok = parsed.Parsed["authorizationRuleName"]; !ok {
		return nil, fmt.Errorf("the segment 'authorizationRuleName' was not found in the resource id %q", input)
	}

	return &id, nil
}

// ValidateQueueAuthorizationRuleID checks that 'input' can be parsed as a Queue Authorization Rule ID
func ValidateQueueAuthorizationRuleID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseQueueAuthorizationRuleID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Queue Authorization Rule ID
func (id QueueAuthorizationRuleId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.ServiceBus/namespaces/%s/queues/%s/authorizationRules/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.NamespaceName, id.QueueName, id.AuthorizationRuleName)
}

// Segments returns a slice of Resource ID Segments which comprise this Queue Authorization Rule ID
func (id QueueAuthorizationRuleId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftServiceBus", "Microsoft.ServiceBus", "Microsoft.ServiceBus"),
		resourceids.StaticSegment("staticNamespaces", "namespaces", "namespaces"),
		resourceids.UserSpecifiedSegment("namespaceName", "namespaceValue"),
		resourceids.StaticSegment("staticQueues", "queues", "queues"),
		resourceids.UserSpecifiedSegment("queueName", "queueValue"),
		resourceids.StaticSegment("staticAuthorizationRules", "authorizationRules", "authorizationRules"),
		resourceids.UserSpecifiedSegment("authorizationRuleName", "authorizationRuleValue"),
	}
}

// String returns a human-readable description of this Queue Authorization Rule ID
func (id QueueAuthorizationRuleId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Namespace Name: %q", id.NamespaceName),
		fmt.Sprintf("Queue Name: %q", id.QueueName),
		fmt.Sprintf("Authorization Rule Name: %q", id.AuthorizationRuleName),
	}
	return fmt.Sprintf("Queue Authorization Rule (%s)", strings.Join(components, "\n"))
}
