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
	recaser.RegisterResourceId(&ScheduleId{})
}

var _ resourceids.ResourceId = &ScheduleId{}

// ScheduleId is a struct representing the Resource ID for a Schedule
type ScheduleId struct {
	SubscriptionId    string
	ResourceGroupName string
	ProjectName       string
	PoolName          string
	ScheduleName      string
}

// NewScheduleID returns a new ScheduleId struct
func NewScheduleID(subscriptionId string, resourceGroupName string, projectName string, poolName string, scheduleName string) ScheduleId {
	return ScheduleId{
		SubscriptionId:    subscriptionId,
		ResourceGroupName: resourceGroupName,
		ProjectName:       projectName,
		PoolName:          poolName,
		ScheduleName:      scheduleName,
	}
}

// ParseScheduleID parses 'input' into a ScheduleId
func ParseScheduleID(input string) (*ScheduleId, error) {
	parser := resourceids.NewParserFromResourceIdType(&ScheduleId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := ScheduleId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParseScheduleIDInsensitively parses 'input' case-insensitively into a ScheduleId
// note: this method should only be used for API response data and not user input
func ParseScheduleIDInsensitively(input string) (*ScheduleId, error) {
	parser := resourceids.NewParserFromResourceIdType(&ScheduleId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := ScheduleId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *ScheduleId) FromParseResult(input resourceids.ParseResult) error {
	var ok bool

	if id.SubscriptionId, ok = input.Parsed["subscriptionId"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", input)
	}

	if id.ResourceGroupName, ok = input.Parsed["resourceGroupName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", input)
	}

	if id.ProjectName, ok = input.Parsed["projectName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "projectName", input)
	}

	if id.PoolName, ok = input.Parsed["poolName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "poolName", input)
	}

	if id.ScheduleName, ok = input.Parsed["scheduleName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "scheduleName", input)
	}

	return nil
}

// ValidateScheduleID checks that 'input' can be parsed as a Schedule ID
func ValidateScheduleID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseScheduleID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Schedule ID
func (id ScheduleId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.DevCenter/projects/%s/pools/%s/schedules/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.ProjectName, id.PoolName, id.ScheduleName)
}

// Segments returns a slice of Resource ID Segments which comprise this Schedule ID
func (id ScheduleId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftDevCenter", "Microsoft.DevCenter", "Microsoft.DevCenter"),
		resourceids.StaticSegment("staticProjects", "projects", "projects"),
		resourceids.UserSpecifiedSegment("projectName", "projectName"),
		resourceids.StaticSegment("staticPools", "pools", "pools"),
		resourceids.UserSpecifiedSegment("poolName", "poolName"),
		resourceids.StaticSegment("staticSchedules", "schedules", "schedules"),
		resourceids.UserSpecifiedSegment("scheduleName", "scheduleName"),
	}
}

// String returns a human-readable description of this Schedule ID
func (id ScheduleId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Project Name: %q", id.ProjectName),
		fmt.Sprintf("Pool Name: %q", id.PoolName),
		fmt.Sprintf("Schedule Name: %q", id.ScheduleName),
	}
	return fmt.Sprintf("Schedule (%s)", strings.Join(components, "\n"))
}
