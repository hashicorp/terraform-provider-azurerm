package taskruns

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/recaser"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

func init() {
	recaser.RegisterResourceId(&TaskRunId{})
}

var _ resourceids.ResourceId = &TaskRunId{}

// TaskRunId is a struct representing the Resource ID for a Task Run
type TaskRunId struct {
	SubscriptionId    string
	ResourceGroupName string
	RegistryName      string
	TaskRunName       string
}

// NewTaskRunID returns a new TaskRunId struct
func NewTaskRunID(subscriptionId string, resourceGroupName string, registryName string, taskRunName string) TaskRunId {
	return TaskRunId{
		SubscriptionId:    subscriptionId,
		ResourceGroupName: resourceGroupName,
		RegistryName:      registryName,
		TaskRunName:       taskRunName,
	}
}

// ParseTaskRunID parses 'input' into a TaskRunId
func ParseTaskRunID(input string) (*TaskRunId, error) {
	parser := resourceids.NewParserFromResourceIdType(&TaskRunId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := TaskRunId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParseTaskRunIDInsensitively parses 'input' case-insensitively into a TaskRunId
// note: this method should only be used for API response data and not user input
func ParseTaskRunIDInsensitively(input string) (*TaskRunId, error) {
	parser := resourceids.NewParserFromResourceIdType(&TaskRunId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := TaskRunId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *TaskRunId) FromParseResult(input resourceids.ParseResult) error {
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

	if id.TaskRunName, ok = input.Parsed["taskRunName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "taskRunName", input)
	}

	return nil
}

// ValidateTaskRunID checks that 'input' can be parsed as a Task Run ID
func ValidateTaskRunID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseTaskRunID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Task Run ID
func (id TaskRunId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.ContainerRegistry/registries/%s/taskRuns/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.RegistryName, id.TaskRunName)
}

// Segments returns a slice of Resource ID Segments which comprise this Task Run ID
func (id TaskRunId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftContainerRegistry", "Microsoft.ContainerRegistry", "Microsoft.ContainerRegistry"),
		resourceids.StaticSegment("staticRegistries", "registries", "registries"),
		resourceids.UserSpecifiedSegment("registryName", "registryName"),
		resourceids.StaticSegment("staticTaskRuns", "taskRuns", "taskRuns"),
		resourceids.UserSpecifiedSegment("taskRunName", "taskRunName"),
	}
}

// String returns a human-readable description of this Task Run ID
func (id TaskRunId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Registry Name: %q", id.RegistryName),
		fmt.Sprintf("Task Run Name: %q", id.TaskRunName),
	}
	return fmt.Sprintf("Task Run (%s)", strings.Join(components, "\n"))
}
