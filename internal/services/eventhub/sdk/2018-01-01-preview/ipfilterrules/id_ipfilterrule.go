package ipfilterrules

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

var _ resourceids.ResourceId = IpfilterruleId{}

// IpfilterruleId is a struct representing the Resource ID for a Ipfilterrule
type IpfilterruleId struct {
	SubscriptionId    string
	ResourceGroupName string
	NamespaceName     string
	IpFilterRuleName  string
}

// NewIpfilterruleID returns a new IpfilterruleId struct
func NewIpfilterruleID(subscriptionId string, resourceGroupName string, namespaceName string, ipFilterRuleName string) IpfilterruleId {
	return IpfilterruleId{
		SubscriptionId:    subscriptionId,
		ResourceGroupName: resourceGroupName,
		NamespaceName:     namespaceName,
		IpFilterRuleName:  ipFilterRuleName,
	}
}

// ParseIpfilterruleID parses 'input' into a IpfilterruleId
func ParseIpfilterruleID(input string) (*IpfilterruleId, error) {
	parser := resourceids.NewParserFromResourceIdType(IpfilterruleId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := IpfilterruleId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, fmt.Errorf("the segment 'subscriptionId' was not found in the resource id %q", input)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, fmt.Errorf("the segment 'resourceGroupName' was not found in the resource id %q", input)
	}

	if id.NamespaceName, ok = parsed.Parsed["namespaceName"]; !ok {
		return nil, fmt.Errorf("the segment 'namespaceName' was not found in the resource id %q", input)
	}

	if id.IpFilterRuleName, ok = parsed.Parsed["ipFilterRuleName"]; !ok {
		return nil, fmt.Errorf("the segment 'ipFilterRuleName' was not found in the resource id %q", input)
	}

	return &id, nil
}

// ParseIpfilterruleIDInsensitively parses 'input' case-insensitively into a IpfilterruleId
// note: this method should only be used for API response data and not user input
func ParseIpfilterruleIDInsensitively(input string) (*IpfilterruleId, error) {
	parser := resourceids.NewParserFromResourceIdType(IpfilterruleId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := IpfilterruleId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, fmt.Errorf("the segment 'subscriptionId' was not found in the resource id %q", input)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, fmt.Errorf("the segment 'resourceGroupName' was not found in the resource id %q", input)
	}

	if id.NamespaceName, ok = parsed.Parsed["namespaceName"]; !ok {
		return nil, fmt.Errorf("the segment 'namespaceName' was not found in the resource id %q", input)
	}

	if id.IpFilterRuleName, ok = parsed.Parsed["ipFilterRuleName"]; !ok {
		return nil, fmt.Errorf("the segment 'ipFilterRuleName' was not found in the resource id %q", input)
	}

	return &id, nil
}

// ValidateIpfilterruleID checks that 'input' can be parsed as a Ipfilterrule ID
func ValidateIpfilterruleID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseIpfilterruleID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Ipfilterrule ID
func (id IpfilterruleId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.EventHub/namespaces/%s/ipfilterrules/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.NamespaceName, id.IpFilterRuleName)
}

// Segments returns a slice of Resource ID Segments which comprise this Ipfilterrule ID
func (id IpfilterruleId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftEventHub", "Microsoft.EventHub", "Microsoft.EventHub"),
		resourceids.StaticSegment("staticNamespaces", "namespaces", "namespaces"),
		resourceids.UserSpecifiedSegment("namespaceName", "namespaceValue"),
		resourceids.StaticSegment("staticIpfilterrules", "ipfilterrules", "ipfilterrules"),
		resourceids.UserSpecifiedSegment("ipFilterRuleName", "ipFilterRuleValue"),
	}
}

// String returns a human-readable description of this Ipfilterrule ID
func (id IpfilterruleId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Namespace Name: %q", id.NamespaceName),
		fmt.Sprintf("Ip Filter Rule Name: %q", id.IpFilterRuleName),
	}
	return fmt.Sprintf("Ipfilterrule (%s)", strings.Join(components, "\n"))
}
