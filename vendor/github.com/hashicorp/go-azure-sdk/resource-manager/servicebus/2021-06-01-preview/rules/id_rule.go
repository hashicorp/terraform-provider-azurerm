package rules

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

var _ resourceids.ResourceId = RuleId{}

// RuleId is a struct representing the Resource ID for a Rule
type RuleId struct {
	SubscriptionId    string
	ResourceGroupName string
	NamespaceName     string
	TopicName         string
	SubscriptionName  string
	RuleName          string
}

// NewRuleID returns a new RuleId struct
func NewRuleID(subscriptionId string, resourceGroupName string, namespaceName string, topicName string, subscriptionName string, ruleName string) RuleId {
	return RuleId{
		SubscriptionId:    subscriptionId,
		ResourceGroupName: resourceGroupName,
		NamespaceName:     namespaceName,
		TopicName:         topicName,
		SubscriptionName:  subscriptionName,
		RuleName:          ruleName,
	}
}

// ParseRuleID parses 'input' into a RuleId
func ParseRuleID(input string) (*RuleId, error) {
	parser := resourceids.NewParserFromResourceIdType(RuleId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := RuleId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, fmt.Errorf("the segment 'subscriptionId' was not found in the resource id %q", input)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, fmt.Errorf("the segment 'resourceGroupName' was not found in the resource id %q", input)
	}

	if id.NamespaceName, ok = parsed.Parsed["namespaceName"]; !ok {
		return nil, fmt.Errorf("the segment 'namespaceName' was not found in the resource id %q", input)
	}

	if id.TopicName, ok = parsed.Parsed["topicName"]; !ok {
		return nil, fmt.Errorf("the segment 'topicName' was not found in the resource id %q", input)
	}

	if id.SubscriptionName, ok = parsed.Parsed["subscriptionName"]; !ok {
		return nil, fmt.Errorf("the segment 'subscriptionName' was not found in the resource id %q", input)
	}

	if id.RuleName, ok = parsed.Parsed["ruleName"]; !ok {
		return nil, fmt.Errorf("the segment 'ruleName' was not found in the resource id %q", input)
	}

	return &id, nil
}

// ParseRuleIDInsensitively parses 'input' case-insensitively into a RuleId
// note: this method should only be used for API response data and not user input
func ParseRuleIDInsensitively(input string) (*RuleId, error) {
	parser := resourceids.NewParserFromResourceIdType(RuleId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := RuleId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, fmt.Errorf("the segment 'subscriptionId' was not found in the resource id %q", input)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, fmt.Errorf("the segment 'resourceGroupName' was not found in the resource id %q", input)
	}

	if id.NamespaceName, ok = parsed.Parsed["namespaceName"]; !ok {
		return nil, fmt.Errorf("the segment 'namespaceName' was not found in the resource id %q", input)
	}

	if id.TopicName, ok = parsed.Parsed["topicName"]; !ok {
		return nil, fmt.Errorf("the segment 'topicName' was not found in the resource id %q", input)
	}

	if id.SubscriptionName, ok = parsed.Parsed["subscriptionName"]; !ok {
		return nil, fmt.Errorf("the segment 'subscriptionName' was not found in the resource id %q", input)
	}

	if id.RuleName, ok = parsed.Parsed["ruleName"]; !ok {
		return nil, fmt.Errorf("the segment 'ruleName' was not found in the resource id %q", input)
	}

	return &id, nil
}

// ValidateRuleID checks that 'input' can be parsed as a Rule ID
func ValidateRuleID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseRuleID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Rule ID
func (id RuleId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.ServiceBus/namespaces/%s/topics/%s/subscriptions/%s/rules/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.NamespaceName, id.TopicName, id.SubscriptionName, id.RuleName)
}

// Segments returns a slice of Resource ID Segments which comprise this Rule ID
func (id RuleId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftServiceBus", "Microsoft.ServiceBus", "Microsoft.ServiceBus"),
		resourceids.StaticSegment("staticNamespaces", "namespaces", "namespaces"),
		resourceids.UserSpecifiedSegment("namespaceName", "namespaceValue"),
		resourceids.StaticSegment("staticTopics", "topics", "topics"),
		resourceids.UserSpecifiedSegment("topicName", "topicValue"),
		resourceids.StaticSegment("staticSubscriptions2", "subscriptions", "subscriptions"),
		resourceids.UserSpecifiedSegment("subscriptionName", "subscriptionValue"),
		resourceids.StaticSegment("staticRules", "rules", "rules"),
		resourceids.UserSpecifiedSegment("ruleName", "ruleValue"),
	}
}

// String returns a human-readable description of this Rule ID
func (id RuleId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Namespace Name: %q", id.NamespaceName),
		fmt.Sprintf("Topic Name: %q", id.TopicName),
		fmt.Sprintf("Subscription Name: %q", id.SubscriptionName),
		fmt.Sprintf("Rule Name: %q", id.RuleName),
	}
	return fmt.Sprintf("Rule (%s)", strings.Join(components, "\n"))
}
