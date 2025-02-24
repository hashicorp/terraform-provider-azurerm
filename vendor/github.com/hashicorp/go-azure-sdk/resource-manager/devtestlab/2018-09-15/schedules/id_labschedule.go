package schedules

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/recaser"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

func init() {
	recaser.RegisterResourceId(&LabScheduleId{})
}

var _ resourceids.ResourceId = &LabScheduleId{}

// LabScheduleId is a struct representing the Resource ID for a Lab Schedule
type LabScheduleId struct {
	SubscriptionId    string
	ResourceGroupName string
	LabName           string
	ScheduleName      string
}

// NewLabScheduleID returns a new LabScheduleId struct
func NewLabScheduleID(subscriptionId string, resourceGroupName string, labName string, scheduleName string) LabScheduleId {
	return LabScheduleId{
		SubscriptionId:    subscriptionId,
		ResourceGroupName: resourceGroupName,
		LabName:           labName,
		ScheduleName:      scheduleName,
	}
}

// ParseLabScheduleID parses 'input' into a LabScheduleId
func ParseLabScheduleID(input string) (*LabScheduleId, error) {
	parser := resourceids.NewParserFromResourceIdType(&LabScheduleId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := LabScheduleId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParseLabScheduleIDInsensitively parses 'input' case-insensitively into a LabScheduleId
// note: this method should only be used for API response data and not user input
func ParseLabScheduleIDInsensitively(input string) (*LabScheduleId, error) {
	parser := resourceids.NewParserFromResourceIdType(&LabScheduleId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := LabScheduleId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *LabScheduleId) FromParseResult(input resourceids.ParseResult) error {
	var ok bool

	if id.SubscriptionId, ok = input.Parsed["subscriptionId"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", input)
	}

	if id.ResourceGroupName, ok = input.Parsed["resourceGroupName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", input)
	}

	if id.LabName, ok = input.Parsed["labName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "labName", input)
	}

	if id.ScheduleName, ok = input.Parsed["scheduleName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "scheduleName", input)
	}

	return nil
}

// ValidateLabScheduleID checks that 'input' can be parsed as a Lab Schedule ID
func ValidateLabScheduleID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseLabScheduleID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Lab Schedule ID
func (id LabScheduleId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.DevTestLab/labs/%s/schedules/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.LabName, id.ScheduleName)
}

// Segments returns a slice of Resource ID Segments which comprise this Lab Schedule ID
func (id LabScheduleId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftDevTestLab", "Microsoft.DevTestLab", "Microsoft.DevTestLab"),
		resourceids.StaticSegment("staticLabs", "labs", "labs"),
		resourceids.UserSpecifiedSegment("labName", "labName"),
		resourceids.StaticSegment("staticSchedules", "schedules", "schedules"),
		resourceids.UserSpecifiedSegment("scheduleName", "scheduleName"),
	}
}

// String returns a human-readable description of this Lab Schedule ID
func (id LabScheduleId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Lab Name: %q", id.LabName),
		fmt.Sprintf("Schedule Name: %q", id.ScheduleName),
	}
	return fmt.Sprintf("Lab Schedule (%s)", strings.Join(components, "\n"))
}
