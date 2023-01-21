package workflows

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

var _ resourceids.ResourceId = LocationWorkflowId{}

// LocationWorkflowId is a struct representing the Resource ID for a Location Workflow
type LocationWorkflowId struct {
	SubscriptionId    string
	ResourceGroupName string
	Location          string
	WorkflowName      string
}

// NewLocationWorkflowID returns a new LocationWorkflowId struct
func NewLocationWorkflowID(subscriptionId string, resourceGroupName string, location string, workflowName string) LocationWorkflowId {
	return LocationWorkflowId{
		SubscriptionId:    subscriptionId,
		ResourceGroupName: resourceGroupName,
		Location:          location,
		WorkflowName:      workflowName,
	}
}

// ParseLocationWorkflowID parses 'input' into a LocationWorkflowId
func ParseLocationWorkflowID(input string) (*LocationWorkflowId, error) {
	parser := resourceids.NewParserFromResourceIdType(LocationWorkflowId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := LocationWorkflowId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, fmt.Errorf("the segment 'subscriptionId' was not found in the resource id %q", input)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, fmt.Errorf("the segment 'resourceGroupName' was not found in the resource id %q", input)
	}

	if id.Location, ok = parsed.Parsed["location"]; !ok {
		return nil, fmt.Errorf("the segment 'location' was not found in the resource id %q", input)
	}

	if id.WorkflowName, ok = parsed.Parsed["workflowName"]; !ok {
		return nil, fmt.Errorf("the segment 'workflowName' was not found in the resource id %q", input)
	}

	return &id, nil
}

// ParseLocationWorkflowIDInsensitively parses 'input' case-insensitively into a LocationWorkflowId
// note: this method should only be used for API response data and not user input
func ParseLocationWorkflowIDInsensitively(input string) (*LocationWorkflowId, error) {
	parser := resourceids.NewParserFromResourceIdType(LocationWorkflowId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := LocationWorkflowId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, fmt.Errorf("the segment 'subscriptionId' was not found in the resource id %q", input)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, fmt.Errorf("the segment 'resourceGroupName' was not found in the resource id %q", input)
	}

	if id.Location, ok = parsed.Parsed["location"]; !ok {
		return nil, fmt.Errorf("the segment 'location' was not found in the resource id %q", input)
	}

	if id.WorkflowName, ok = parsed.Parsed["workflowName"]; !ok {
		return nil, fmt.Errorf("the segment 'workflowName' was not found in the resource id %q", input)
	}

	return &id, nil
}

// ValidateLocationWorkflowID checks that 'input' can be parsed as a Location Workflow ID
func ValidateLocationWorkflowID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseLocationWorkflowID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Location Workflow ID
func (id LocationWorkflowId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Logic/locations/%s/workflows/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.Location, id.WorkflowName)
}

// Segments returns a slice of Resource ID Segments which comprise this Location Workflow ID
func (id LocationWorkflowId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftLogic", "Microsoft.Logic", "Microsoft.Logic"),
		resourceids.StaticSegment("staticLocations", "locations", "locations"),
		resourceids.UserSpecifiedSegment("location", "locationValue"),
		resourceids.StaticSegment("staticWorkflows", "workflows", "workflows"),
		resourceids.UserSpecifiedSegment("workflowName", "workflowValue"),
	}
}

// String returns a human-readable description of this Location Workflow ID
func (id LocationWorkflowId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Location: %q", id.Location),
		fmt.Sprintf("Workflow Name: %q", id.WorkflowName),
	}
	return fmt.Sprintf("Location Workflow (%s)", strings.Join(components, "\n"))
}
