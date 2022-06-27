package application

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

var _ resourceids.ResourceId = ApplicationId{}

// ApplicationId is a struct representing the Resource ID for a Application
type ApplicationId struct {
	SubscriptionId       string
	ResourceGroupName    string
	ApplicationGroupName string
	ApplicationName      string
}

// NewApplicationID returns a new ApplicationId struct
func NewApplicationID(subscriptionId string, resourceGroupName string, applicationGroupName string, applicationName string) ApplicationId {
	return ApplicationId{
		SubscriptionId:       subscriptionId,
		ResourceGroupName:    resourceGroupName,
		ApplicationGroupName: applicationGroupName,
		ApplicationName:      applicationName,
	}
}

// ParseApplicationID parses 'input' into a ApplicationId
func ParseApplicationID(input string) (*ApplicationId, error) {
	parser := resourceids.NewParserFromResourceIdType(ApplicationId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := ApplicationId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, fmt.Errorf("the segment 'subscriptionId' was not found in the resource id %q", input)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, fmt.Errorf("the segment 'resourceGroupName' was not found in the resource id %q", input)
	}

	if id.ApplicationGroupName, ok = parsed.Parsed["applicationGroupName"]; !ok {
		return nil, fmt.Errorf("the segment 'applicationGroupName' was not found in the resource id %q", input)
	}

	if id.ApplicationName, ok = parsed.Parsed["applicationName"]; !ok {
		return nil, fmt.Errorf("the segment 'applicationName' was not found in the resource id %q", input)
	}

	return &id, nil
}

// ParseApplicationIDInsensitively parses 'input' case-insensitively into a ApplicationId
// note: this method should only be used for API response data and not user input
func ParseApplicationIDInsensitively(input string) (*ApplicationId, error) {
	parser := resourceids.NewParserFromResourceIdType(ApplicationId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := ApplicationId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, fmt.Errorf("the segment 'subscriptionId' was not found in the resource id %q", input)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, fmt.Errorf("the segment 'resourceGroupName' was not found in the resource id %q", input)
	}

	if id.ApplicationGroupName, ok = parsed.Parsed["applicationGroupName"]; !ok {
		return nil, fmt.Errorf("the segment 'applicationGroupName' was not found in the resource id %q", input)
	}

	if id.ApplicationName, ok = parsed.Parsed["applicationName"]; !ok {
		return nil, fmt.Errorf("the segment 'applicationName' was not found in the resource id %q", input)
	}

	return &id, nil
}

// ValidateApplicationID checks that 'input' can be parsed as a Application ID
func ValidateApplicationID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseApplicationID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Application ID
func (id ApplicationId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.DesktopVirtualization/applicationGroups/%s/applications/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.ApplicationGroupName, id.ApplicationName)
}

// Segments returns a slice of Resource ID Segments which comprise this Application ID
func (id ApplicationId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftDesktopVirtualization", "Microsoft.DesktopVirtualization", "Microsoft.DesktopVirtualization"),
		resourceids.StaticSegment("staticApplicationGroups", "applicationGroups", "applicationGroups"),
		resourceids.UserSpecifiedSegment("applicationGroupName", "applicationGroupValue"),
		resourceids.StaticSegment("staticApplications", "applications", "applications"),
		resourceids.UserSpecifiedSegment("applicationName", "applicationValue"),
	}
}

// String returns a human-readable description of this Application ID
func (id ApplicationId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Application Group Name: %q", id.ApplicationGroupName),
		fmt.Sprintf("Application Name: %q", id.ApplicationName),
	}
	return fmt.Sprintf("Application (%s)", strings.Join(components, "\n"))
}
