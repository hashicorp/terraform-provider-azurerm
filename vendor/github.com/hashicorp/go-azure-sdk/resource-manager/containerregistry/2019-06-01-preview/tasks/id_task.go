package tasks

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/recaser"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

func init() {
	recaser.RegisterResourceId(&TaskId{})
}

var _ resourceids.ResourceId = &TaskId{}

// TaskId is a struct representing the Resource ID for a Task
type TaskId struct {
	SubscriptionId    string
	ResourceGroupName string
	RegistryName      string
	TaskName          string
}

// NewTaskID returns a new TaskId struct
func NewTaskID(subscriptionId string, resourceGroupName string, registryName string, taskName string) TaskId {
	return TaskId{
		SubscriptionId:    subscriptionId,
		ResourceGroupName: resourceGroupName,
		RegistryName:      registryName,
		TaskName:          taskName,
	}
}

// ParseTaskID parses 'input' into a TaskId
func ParseTaskID(input string) (*TaskId, error) {
	parser := resourceids.NewParserFromResourceIdType(&TaskId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := TaskId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParseTaskIDInsensitively parses 'input' case-insensitively into a TaskId
// note: this method should only be used for API response data and not user input
func ParseTaskIDInsensitively(input string) (*TaskId, error) {
	parser := resourceids.NewParserFromResourceIdType(&TaskId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := TaskId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *TaskId) FromParseResult(input resourceids.ParseResult) error {
	var ok bool

	if id.SubscriptionId, ok = input.Parsed["subscriptionId"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", input)
	}

	if id.ResourceGroupName, ok = input.Parsed["resourceGroupName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", input)
	}

	if id.RegistryName, ok = input.Parsed["registryName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "registryName", input)
	}

	if id.TaskName, ok = input.Parsed["taskName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "taskName", input)
	}

	return nil
}

// ValidateTaskID checks that 'input' can be parsed as a Task ID
func ValidateTaskID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseTaskID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Task ID
func (id TaskId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.ContainerRegistry/registries/%s/tasks/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.RegistryName, id.TaskName)
}

// Segments returns a slice of Resource ID Segments which comprise this Task ID
func (id TaskId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftContainerRegistry", "Microsoft.ContainerRegistry", "Microsoft.ContainerRegistry"),
		resourceids.StaticSegment("staticRegistries", "registries", "registries"),
		resourceids.UserSpecifiedSegment("registryName", "registryName"),
		resourceids.StaticSegment("staticTasks", "tasks", "tasks"),
		resourceids.UserSpecifiedSegment("taskName", "taskName"),
	}
}

// String returns a human-readable description of this Task ID
func (id TaskId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Registry Name: %q", id.RegistryName),
		fmt.Sprintf("Task Name: %q", id.TaskName),
	}
	return fmt.Sprintf("Task (%s)", strings.Join(components, "\n"))
}
