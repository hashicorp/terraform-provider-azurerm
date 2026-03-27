package taskhubs

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/recaser"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

func init() {
	recaser.RegisterResourceId(&TaskHubId{})
}

var _ resourceids.ResourceId = &TaskHubId{}

// TaskHubId is a struct representing the Resource ID for a Task Hub
type TaskHubId struct {
	SubscriptionId    string
	ResourceGroupName string
	SchedulerName     string
	TaskHubName       string
}

// NewTaskHubID returns a new TaskHubId struct
func NewTaskHubID(subscriptionId string, resourceGroupName string, schedulerName string, taskHubName string) TaskHubId {
	return TaskHubId{
		SubscriptionId:    subscriptionId,
		ResourceGroupName: resourceGroupName,
		SchedulerName:     schedulerName,
		TaskHubName:       taskHubName,
	}
}

// ParseTaskHubID parses 'input' into a TaskHubId
func ParseTaskHubID(input string) (*TaskHubId, error) {
	parser := resourceids.NewParserFromResourceIdType(&TaskHubId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := TaskHubId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParseTaskHubIDInsensitively parses 'input' case-insensitively into a TaskHubId
// note: this method should only be used for API response data and not user input
func ParseTaskHubIDInsensitively(input string) (*TaskHubId, error) {
	parser := resourceids.NewParserFromResourceIdType(&TaskHubId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := TaskHubId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *TaskHubId) FromParseResult(input resourceids.ParseResult) error {
	var ok bool

	if id.SubscriptionId, ok = input.Parsed["subscriptionId"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", input)
	}

	if id.ResourceGroupName, ok = input.Parsed["resourceGroupName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", input)
	}

	if id.SchedulerName, ok = input.Parsed["schedulerName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "schedulerName", input)
	}

	if id.TaskHubName, ok = input.Parsed["taskHubName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "taskHubName", input)
	}

	return nil
}

// ValidateTaskHubID checks that 'input' can be parsed as a Task Hub ID
func ValidateTaskHubID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseTaskHubID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Task Hub ID
func (id TaskHubId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.DurableTask/schedulers/%s/taskHubs/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.SchedulerName, id.TaskHubName)
}

// Segments returns a slice of Resource ID Segments which comprise this Task Hub ID
func (id TaskHubId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftDurableTask", "Microsoft.DurableTask", "Microsoft.DurableTask"),
		resourceids.StaticSegment("staticSchedulers", "schedulers", "schedulers"),
		resourceids.UserSpecifiedSegment("schedulerName", "schedulerName"),
		resourceids.StaticSegment("staticTaskHubs", "taskHubs", "taskHubs"),
		resourceids.UserSpecifiedSegment("taskHubName", "taskHubName"),
	}
}

// String returns a human-readable description of this Task Hub ID
func (id TaskHubId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Scheduler Name: %q", id.SchedulerName),
		fmt.Sprintf("Task Hub Name: %q", id.TaskHubName),
	}
	return fmt.Sprintf("Task Hub (%s)", strings.Join(components, "\n"))
}
