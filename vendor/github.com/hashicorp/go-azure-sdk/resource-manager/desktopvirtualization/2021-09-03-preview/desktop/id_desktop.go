package desktop

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

var _ resourceids.ResourceId = DesktopId{}

// DesktopId is a struct representing the Resource ID for a Desktop
type DesktopId struct {
	SubscriptionId       string
	ResourceGroupName    string
	ApplicationGroupName string
	DesktopName          string
}

// NewDesktopID returns a new DesktopId struct
func NewDesktopID(subscriptionId string, resourceGroupName string, applicationGroupName string, desktopName string) DesktopId {
	return DesktopId{
		SubscriptionId:       subscriptionId,
		ResourceGroupName:    resourceGroupName,
		ApplicationGroupName: applicationGroupName,
		DesktopName:          desktopName,
	}
}

// ParseDesktopID parses 'input' into a DesktopId
func ParseDesktopID(input string) (*DesktopId, error) {
	parser := resourceids.NewParserFromResourceIdType(DesktopId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := DesktopId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, fmt.Errorf("the segment 'subscriptionId' was not found in the resource id %q", input)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, fmt.Errorf("the segment 'resourceGroupName' was not found in the resource id %q", input)
	}

	if id.ApplicationGroupName, ok = parsed.Parsed["applicationGroupName"]; !ok {
		return nil, fmt.Errorf("the segment 'applicationGroupName' was not found in the resource id %q", input)
	}

	if id.DesktopName, ok = parsed.Parsed["desktopName"]; !ok {
		return nil, fmt.Errorf("the segment 'desktopName' was not found in the resource id %q", input)
	}

	return &id, nil
}

// ParseDesktopIDInsensitively parses 'input' case-insensitively into a DesktopId
// note: this method should only be used for API response data and not user input
func ParseDesktopIDInsensitively(input string) (*DesktopId, error) {
	parser := resourceids.NewParserFromResourceIdType(DesktopId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := DesktopId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, fmt.Errorf("the segment 'subscriptionId' was not found in the resource id %q", input)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, fmt.Errorf("the segment 'resourceGroupName' was not found in the resource id %q", input)
	}

	if id.ApplicationGroupName, ok = parsed.Parsed["applicationGroupName"]; !ok {
		return nil, fmt.Errorf("the segment 'applicationGroupName' was not found in the resource id %q", input)
	}

	if id.DesktopName, ok = parsed.Parsed["desktopName"]; !ok {
		return nil, fmt.Errorf("the segment 'desktopName' was not found in the resource id %q", input)
	}

	return &id, nil
}

// ValidateDesktopID checks that 'input' can be parsed as a Desktop ID
func ValidateDesktopID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseDesktopID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Desktop ID
func (id DesktopId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.DesktopVirtualization/applicationGroups/%s/desktops/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.ApplicationGroupName, id.DesktopName)
}

// Segments returns a slice of Resource ID Segments which comprise this Desktop ID
func (id DesktopId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftDesktopVirtualization", "Microsoft.DesktopVirtualization", "Microsoft.DesktopVirtualization"),
		resourceids.StaticSegment("staticApplicationGroups", "applicationGroups", "applicationGroups"),
		resourceids.UserSpecifiedSegment("applicationGroupName", "applicationGroupValue"),
		resourceids.StaticSegment("staticDesktops", "desktops", "desktops"),
		resourceids.UserSpecifiedSegment("desktopName", "desktopValue"),
	}
}

// String returns a human-readable description of this Desktop ID
func (id DesktopId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Application Group Name: %q", id.ApplicationGroupName),
		fmt.Sprintf("Desktop Name: %q", id.DesktopName),
	}
	return fmt.Sprintf("Desktop (%s)", strings.Join(components, "\n"))
}
