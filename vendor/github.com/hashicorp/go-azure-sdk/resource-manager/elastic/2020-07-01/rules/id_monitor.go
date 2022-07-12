package rules

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

var _ resourceids.ResourceId = MonitorId{}

// MonitorId is a struct representing the Resource ID for a Monitor
type MonitorId struct {
	SubscriptionId    string
	ResourceGroupName string
	MonitorName       string
}

// NewMonitorID returns a new MonitorId struct
func NewMonitorID(subscriptionId string, resourceGroupName string, monitorName string) MonitorId {
	return MonitorId{
		SubscriptionId:    subscriptionId,
		ResourceGroupName: resourceGroupName,
		MonitorName:       monitorName,
	}
}

// ParseMonitorID parses 'input' into a MonitorId
func ParseMonitorID(input string) (*MonitorId, error) {
	parser := resourceids.NewParserFromResourceIdType(MonitorId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := MonitorId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, fmt.Errorf("the segment 'subscriptionId' was not found in the resource id %q", input)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, fmt.Errorf("the segment 'resourceGroupName' was not found in the resource id %q", input)
	}

	if id.MonitorName, ok = parsed.Parsed["monitorName"]; !ok {
		return nil, fmt.Errorf("the segment 'monitorName' was not found in the resource id %q", input)
	}

	return &id, nil
}

// ParseMonitorIDInsensitively parses 'input' case-insensitively into a MonitorId
// note: this method should only be used for API response data and not user input
func ParseMonitorIDInsensitively(input string) (*MonitorId, error) {
	parser := resourceids.NewParserFromResourceIdType(MonitorId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := MonitorId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, fmt.Errorf("the segment 'subscriptionId' was not found in the resource id %q", input)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, fmt.Errorf("the segment 'resourceGroupName' was not found in the resource id %q", input)
	}

	if id.MonitorName, ok = parsed.Parsed["monitorName"]; !ok {
		return nil, fmt.Errorf("the segment 'monitorName' was not found in the resource id %q", input)
	}

	return &id, nil
}

// ValidateMonitorID checks that 'input' can be parsed as a Monitor ID
func ValidateMonitorID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseMonitorID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Monitor ID
func (id MonitorId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Elastic/monitors/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.MonitorName)
}

// Segments returns a slice of Resource ID Segments which comprise this Monitor ID
func (id MonitorId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftElastic", "Microsoft.Elastic", "Microsoft.Elastic"),
		resourceids.StaticSegment("staticMonitors", "monitors", "monitors"),
		resourceids.UserSpecifiedSegment("monitorName", "monitorValue"),
	}
}

// String returns a human-readable description of this Monitor ID
func (id MonitorId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Monitor Name: %q", id.MonitorName),
	}
	return fmt.Sprintf("Monitor (%s)", strings.Join(components, "\n"))
}
