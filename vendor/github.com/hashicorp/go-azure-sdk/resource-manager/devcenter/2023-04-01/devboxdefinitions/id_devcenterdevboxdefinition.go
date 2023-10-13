package devboxdefinitions

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ resourceids.ResourceId = DevCenterDevBoxDefinitionId{}

// DevCenterDevBoxDefinitionId is a struct representing the Resource ID for a Dev Center Dev Box Definition
type DevCenterDevBoxDefinitionId struct {
	SubscriptionId       string
	ResourceGroupName    string
	DevCenterName        string
	DevBoxDefinitionName string
}

// NewDevCenterDevBoxDefinitionID returns a new DevCenterDevBoxDefinitionId struct
func NewDevCenterDevBoxDefinitionID(subscriptionId string, resourceGroupName string, devCenterName string, devBoxDefinitionName string) DevCenterDevBoxDefinitionId {
	return DevCenterDevBoxDefinitionId{
		SubscriptionId:       subscriptionId,
		ResourceGroupName:    resourceGroupName,
		DevCenterName:        devCenterName,
		DevBoxDefinitionName: devBoxDefinitionName,
	}
}

// ParseDevCenterDevBoxDefinitionID parses 'input' into a DevCenterDevBoxDefinitionId
func ParseDevCenterDevBoxDefinitionID(input string) (*DevCenterDevBoxDefinitionId, error) {
	parser := resourceids.NewParserFromResourceIdType(DevCenterDevBoxDefinitionId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := DevCenterDevBoxDefinitionId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", *parsed)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", *parsed)
	}

	if id.DevCenterName, ok = parsed.Parsed["devCenterName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "devCenterName", *parsed)
	}

	if id.DevBoxDefinitionName, ok = parsed.Parsed["devBoxDefinitionName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "devBoxDefinitionName", *parsed)
	}

	return &id, nil
}

// ParseDevCenterDevBoxDefinitionIDInsensitively parses 'input' case-insensitively into a DevCenterDevBoxDefinitionId
// note: this method should only be used for API response data and not user input
func ParseDevCenterDevBoxDefinitionIDInsensitively(input string) (*DevCenterDevBoxDefinitionId, error) {
	parser := resourceids.NewParserFromResourceIdType(DevCenterDevBoxDefinitionId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := DevCenterDevBoxDefinitionId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", *parsed)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", *parsed)
	}

	if id.DevCenterName, ok = parsed.Parsed["devCenterName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "devCenterName", *parsed)
	}

	if id.DevBoxDefinitionName, ok = parsed.Parsed["devBoxDefinitionName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "devBoxDefinitionName", *parsed)
	}

	return &id, nil
}

// ValidateDevCenterDevBoxDefinitionID checks that 'input' can be parsed as a Dev Center Dev Box Definition ID
func ValidateDevCenterDevBoxDefinitionID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseDevCenterDevBoxDefinitionID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Dev Center Dev Box Definition ID
func (id DevCenterDevBoxDefinitionId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.DevCenter/devCenters/%s/devBoxDefinitions/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.DevCenterName, id.DevBoxDefinitionName)
}

// Segments returns a slice of Resource ID Segments which comprise this Dev Center Dev Box Definition ID
func (id DevCenterDevBoxDefinitionId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftDevCenter", "Microsoft.DevCenter", "Microsoft.DevCenter"),
		resourceids.StaticSegment("staticDevCenters", "devCenters", "devCenters"),
		resourceids.UserSpecifiedSegment("devCenterName", "devCenterValue"),
		resourceids.StaticSegment("staticDevBoxDefinitions", "devBoxDefinitions", "devBoxDefinitions"),
		resourceids.UserSpecifiedSegment("devBoxDefinitionName", "devBoxDefinitionValue"),
	}
}

// String returns a human-readable description of this Dev Center Dev Box Definition ID
func (id DevCenterDevBoxDefinitionId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Dev Center Name: %q", id.DevCenterName),
		fmt.Sprintf("Dev Box Definition Name: %q", id.DevBoxDefinitionName),
	}
	return fmt.Sprintf("Dev Center Dev Box Definition (%s)", strings.Join(components, "\n"))
}
