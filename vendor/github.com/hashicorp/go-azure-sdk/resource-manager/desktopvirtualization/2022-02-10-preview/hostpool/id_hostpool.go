package hostpool

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

var _ resourceids.ResourceId = HostPoolId{}

// HostPoolId is a struct representing the Resource ID for a Host Pool
type HostPoolId struct {
	SubscriptionId    string
	ResourceGroupName string
	HostPoolName      string
}

// NewHostPoolID returns a new HostPoolId struct
func NewHostPoolID(subscriptionId string, resourceGroupName string, hostPoolName string) HostPoolId {
	return HostPoolId{
		SubscriptionId:    subscriptionId,
		ResourceGroupName: resourceGroupName,
		HostPoolName:      hostPoolName,
	}
}

// ParseHostPoolID parses 'input' into a HostPoolId
func ParseHostPoolID(input string) (*HostPoolId, error) {
	parser := resourceids.NewParserFromResourceIdType(HostPoolId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := HostPoolId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, fmt.Errorf("the segment 'subscriptionId' was not found in the resource id %q", input)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, fmt.Errorf("the segment 'resourceGroupName' was not found in the resource id %q", input)
	}

	if id.HostPoolName, ok = parsed.Parsed["hostPoolName"]; !ok {
		return nil, fmt.Errorf("the segment 'hostPoolName' was not found in the resource id %q", input)
	}

	return &id, nil
}

// ParseHostPoolIDInsensitively parses 'input' case-insensitively into a HostPoolId
// note: this method should only be used for API response data and not user input
func ParseHostPoolIDInsensitively(input string) (*HostPoolId, error) {
	parser := resourceids.NewParserFromResourceIdType(HostPoolId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := HostPoolId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, fmt.Errorf("the segment 'subscriptionId' was not found in the resource id %q", input)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, fmt.Errorf("the segment 'resourceGroupName' was not found in the resource id %q", input)
	}

	if id.HostPoolName, ok = parsed.Parsed["hostPoolName"]; !ok {
		return nil, fmt.Errorf("the segment 'hostPoolName' was not found in the resource id %q", input)
	}

	return &id, nil
}

// ValidateHostPoolID checks that 'input' can be parsed as a Host Pool ID
func ValidateHostPoolID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseHostPoolID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Host Pool ID
func (id HostPoolId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.DesktopVirtualization/hostPools/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.HostPoolName)
}

// Segments returns a slice of Resource ID Segments which comprise this Host Pool ID
func (id HostPoolId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftDesktopVirtualization", "Microsoft.DesktopVirtualization", "Microsoft.DesktopVirtualization"),
		resourceids.StaticSegment("staticHostPools", "hostPools", "hostPools"),
		resourceids.UserSpecifiedSegment("hostPoolName", "hostPoolValue"),
	}
}

// String returns a human-readable description of this Host Pool ID
func (id HostPoolId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Host Pool Name: %q", id.HostPoolName),
	}
	return fmt.Sprintf("Host Pool (%s)", strings.Join(components, "\n"))
}
