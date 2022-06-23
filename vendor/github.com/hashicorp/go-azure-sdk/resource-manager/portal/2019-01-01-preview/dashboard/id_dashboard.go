package dashboard

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

var _ resourceids.ResourceId = DashboardId{}

// DashboardId is a struct representing the Resource ID for a Dashboard
type DashboardId struct {
	SubscriptionId    string
	ResourceGroupName string
	DashboardName     string
}

// NewDashboardID returns a new DashboardId struct
func NewDashboardID(subscriptionId string, resourceGroupName string, dashboardName string) DashboardId {
	return DashboardId{
		SubscriptionId:    subscriptionId,
		ResourceGroupName: resourceGroupName,
		DashboardName:     dashboardName,
	}
}

// ParseDashboardID parses 'input' into a DashboardId
func ParseDashboardID(input string) (*DashboardId, error) {
	parser := resourceids.NewParserFromResourceIdType(DashboardId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := DashboardId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, fmt.Errorf("the segment 'subscriptionId' was not found in the resource id %q", input)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, fmt.Errorf("the segment 'resourceGroupName' was not found in the resource id %q", input)
	}

	if id.DashboardName, ok = parsed.Parsed["dashboardName"]; !ok {
		return nil, fmt.Errorf("the segment 'dashboardName' was not found in the resource id %q", input)
	}

	return &id, nil
}

// ParseDashboardIDInsensitively parses 'input' case-insensitively into a DashboardId
// note: this method should only be used for API response data and not user input
func ParseDashboardIDInsensitively(input string) (*DashboardId, error) {
	parser := resourceids.NewParserFromResourceIdType(DashboardId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := DashboardId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, fmt.Errorf("the segment 'subscriptionId' was not found in the resource id %q", input)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, fmt.Errorf("the segment 'resourceGroupName' was not found in the resource id %q", input)
	}

	if id.DashboardName, ok = parsed.Parsed["dashboardName"]; !ok {
		return nil, fmt.Errorf("the segment 'dashboardName' was not found in the resource id %q", input)
	}

	return &id, nil
}

// ValidateDashboardID checks that 'input' can be parsed as a Dashboard ID
func ValidateDashboardID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseDashboardID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Dashboard ID
func (id DashboardId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Portal/dashboards/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.DashboardName)
}

// Segments returns a slice of Resource ID Segments which comprise this Dashboard ID
func (id DashboardId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftPortal", "Microsoft.Portal", "Microsoft.Portal"),
		resourceids.StaticSegment("staticDashboards", "dashboards", "dashboards"),
		resourceids.UserSpecifiedSegment("dashboardName", "dashboardValue"),
	}
}

// String returns a human-readable description of this Dashboard ID
func (id DashboardId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Dashboard Name: %q", id.DashboardName),
	}
	return fmt.Sprintf("Dashboard (%s)", strings.Join(components, "\n"))
}
