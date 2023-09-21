package jobs

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ resourceids.ResourceId = ExecutionId{}

// ExecutionId is a struct representing the Resource ID for a Execution
type ExecutionId struct {
	SubscriptionId    string
	ResourceGroupName string
	JobName           string
	ExecutionName     string
}

// NewExecutionID returns a new ExecutionId struct
func NewExecutionID(subscriptionId string, resourceGroupName string, jobName string, executionName string) ExecutionId {
	return ExecutionId{
		SubscriptionId:    subscriptionId,
		ResourceGroupName: resourceGroupName,
		JobName:           jobName,
		ExecutionName:     executionName,
	}
}

// ParseExecutionID parses 'input' into a ExecutionId
func ParseExecutionID(input string) (*ExecutionId, error) {
	parser := resourceids.NewParserFromResourceIdType(ExecutionId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := ExecutionId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", *parsed)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", *parsed)
	}

	if id.JobName, ok = parsed.Parsed["jobName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "jobName", *parsed)
	}

	if id.ExecutionName, ok = parsed.Parsed["executionName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "executionName", *parsed)
	}

	return &id, nil
}

// ParseExecutionIDInsensitively parses 'input' case-insensitively into a ExecutionId
// note: this method should only be used for API response data and not user input
func ParseExecutionIDInsensitively(input string) (*ExecutionId, error) {
	parser := resourceids.NewParserFromResourceIdType(ExecutionId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := ExecutionId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", *parsed)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", *parsed)
	}

	if id.JobName, ok = parsed.Parsed["jobName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "jobName", *parsed)
	}

	if id.ExecutionName, ok = parsed.Parsed["executionName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "executionName", *parsed)
	}

	return &id, nil
}

// ValidateExecutionID checks that 'input' can be parsed as a Execution ID
func ValidateExecutionID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseExecutionID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Execution ID
func (id ExecutionId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.App/jobs/%s/executions/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.JobName, id.ExecutionName)
}

// Segments returns a slice of Resource ID Segments which comprise this Execution ID
func (id ExecutionId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftApp", "Microsoft.App", "Microsoft.App"),
		resourceids.StaticSegment("staticJobs", "jobs", "jobs"),
		resourceids.UserSpecifiedSegment("jobName", "jobValue"),
		resourceids.StaticSegment("staticExecutions", "executions", "executions"),
		resourceids.UserSpecifiedSegment("executionName", "executionValue"),
	}
}

// String returns a human-readable description of this Execution ID
func (id ExecutionId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Job Name: %q", id.JobName),
		fmt.Sprintf("Execution Name: %q", id.ExecutionName),
	}
	return fmt.Sprintf("Execution (%s)", strings.Join(components, "\n"))
}
