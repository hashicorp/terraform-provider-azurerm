package environmenttypes

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ resourceids.ResourceId = DevCenterEnvironmentTypeId{}

// DevCenterEnvironmentTypeId is a struct representing the Resource ID for a Dev Center Environment Type
type DevCenterEnvironmentTypeId struct {
	SubscriptionId      string
	ResourceGroupName   string
	DevCenterName       string
	EnvironmentTypeName string
}

// NewDevCenterEnvironmentTypeID returns a new DevCenterEnvironmentTypeId struct
func NewDevCenterEnvironmentTypeID(subscriptionId string, resourceGroupName string, devCenterName string, environmentTypeName string) DevCenterEnvironmentTypeId {
	return DevCenterEnvironmentTypeId{
		SubscriptionId:      subscriptionId,
		ResourceGroupName:   resourceGroupName,
		DevCenterName:       devCenterName,
		EnvironmentTypeName: environmentTypeName,
	}
}

// ParseDevCenterEnvironmentTypeID parses 'input' into a DevCenterEnvironmentTypeId
func ParseDevCenterEnvironmentTypeID(input string) (*DevCenterEnvironmentTypeId, error) {
	parser := resourceids.NewParserFromResourceIdType(DevCenterEnvironmentTypeId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := DevCenterEnvironmentTypeId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", *parsed)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", *parsed)
	}

	if id.DevCenterName, ok = parsed.Parsed["devCenterName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "devCenterName", *parsed)
	}

	if id.EnvironmentTypeName, ok = parsed.Parsed["environmentTypeName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "environmentTypeName", *parsed)
	}

	return &id, nil
}

// ParseDevCenterEnvironmentTypeIDInsensitively parses 'input' case-insensitively into a DevCenterEnvironmentTypeId
// note: this method should only be used for API response data and not user input
func ParseDevCenterEnvironmentTypeIDInsensitively(input string) (*DevCenterEnvironmentTypeId, error) {
	parser := resourceids.NewParserFromResourceIdType(DevCenterEnvironmentTypeId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := DevCenterEnvironmentTypeId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", *parsed)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", *parsed)
	}

	if id.DevCenterName, ok = parsed.Parsed["devCenterName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "devCenterName", *parsed)
	}

	if id.EnvironmentTypeName, ok = parsed.Parsed["environmentTypeName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "environmentTypeName", *parsed)
	}

	return &id, nil
}

// ValidateDevCenterEnvironmentTypeID checks that 'input' can be parsed as a Dev Center Environment Type ID
func ValidateDevCenterEnvironmentTypeID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseDevCenterEnvironmentTypeID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Dev Center Environment Type ID
func (id DevCenterEnvironmentTypeId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.DevCenter/devCenters/%s/environmentTypes/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.DevCenterName, id.EnvironmentTypeName)
}

// Segments returns a slice of Resource ID Segments which comprise this Dev Center Environment Type ID
func (id DevCenterEnvironmentTypeId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftDevCenter", "Microsoft.DevCenter", "Microsoft.DevCenter"),
		resourceids.StaticSegment("staticDevCenters", "devCenters", "devCenters"),
		resourceids.UserSpecifiedSegment("devCenterName", "devCenterValue"),
		resourceids.StaticSegment("staticEnvironmentTypes", "environmentTypes", "environmentTypes"),
		resourceids.UserSpecifiedSegment("environmentTypeName", "environmentTypeValue"),
	}
}

// String returns a human-readable description of this Dev Center Environment Type ID
func (id DevCenterEnvironmentTypeId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Dev Center Name: %q", id.DevCenterName),
		fmt.Sprintf("Environment Type Name: %q", id.EnvironmentTypeName),
	}
	return fmt.Sprintf("Dev Center Environment Type (%s)", strings.Join(components, "\n"))
}
