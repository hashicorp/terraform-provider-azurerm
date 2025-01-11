package workflows

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/recaser"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

func init() {
	recaser.RegisterResourceId(&LocationWorkflowId{})
}

var _ resourceids.ResourceId = &LocationWorkflowId{}

// LocationWorkflowId is a struct representing the Resource ID for a Location Workflow
type LocationWorkflowId struct {
	SubscriptionId    string
	ResourceGroupName string
	LocationName      string
	WorkflowName      string
}

// NewLocationWorkflowID returns a new LocationWorkflowId struct
func NewLocationWorkflowID(subscriptionId string, resourceGroupName string, locationName string, workflowName string) LocationWorkflowId {
	return LocationWorkflowId{
		SubscriptionId:    subscriptionId,
		ResourceGroupName: resourceGroupName,
		LocationName:      locationName,
		WorkflowName:      workflowName,
	}
}

// ParseLocationWorkflowID parses 'input' into a LocationWorkflowId
func ParseLocationWorkflowID(input string) (*LocationWorkflowId, error) {
	parser := resourceids.NewParserFromResourceIdType(&LocationWorkflowId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := LocationWorkflowId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParseLocationWorkflowIDInsensitively parses 'input' case-insensitively into a LocationWorkflowId
// note: this method should only be used for API response data and not user input
func ParseLocationWorkflowIDInsensitively(input string) (*LocationWorkflowId, error) {
	parser := resourceids.NewParserFromResourceIdType(&LocationWorkflowId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := LocationWorkflowId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *LocationWorkflowId) FromParseResult(input resourceids.ParseResult) error {
	var ok bool

	if id.SubscriptionId, ok = input.Parsed["subscriptionId"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", input)
	}

	if id.ResourceGroupName, ok = input.Parsed["resourceGroupName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", input)
	}

	if id.LocationName, ok = input.Parsed["locationName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "locationName", input)
	}

	if id.WorkflowName, ok = input.Parsed["workflowName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "workflowName", input)
	}

	return nil
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
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.LocationName, id.WorkflowName)
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
		resourceids.UserSpecifiedSegment("locationName", "locationName"),
		resourceids.StaticSegment("staticWorkflows", "workflows", "workflows"),
		resourceids.UserSpecifiedSegment("workflowName", "workflowName"),
	}
}

// String returns a human-readable description of this Location Workflow ID
func (id LocationWorkflowId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Location Name: %q", id.LocationName),
		fmt.Sprintf("Workflow Name: %q", id.WorkflowName),
	}
	return fmt.Sprintf("Location Workflow (%s)", strings.Join(components, "\n"))
}
