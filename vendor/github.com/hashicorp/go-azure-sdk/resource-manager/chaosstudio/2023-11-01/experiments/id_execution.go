package experiments

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/recaser"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

func init() {
	recaser.RegisterResourceId(&ExecutionId{})
}

var _ resourceids.ResourceId = &ExecutionId{}

// ExecutionId is a struct representing the Resource ID for a Execution
type ExecutionId struct {
	SubscriptionId    string
	ResourceGroupName string
	ExperimentName    string
	ExecutionId       string
}

// NewExecutionID returns a new ExecutionId struct
func NewExecutionID(subscriptionId string, resourceGroupName string, experimentName string, executionId string) ExecutionId {
	return ExecutionId{
		SubscriptionId:    subscriptionId,
		ResourceGroupName: resourceGroupName,
		ExperimentName:    experimentName,
		ExecutionId:       executionId,
	}
}

// ParseExecutionID parses 'input' into a ExecutionId
func ParseExecutionID(input string) (*ExecutionId, error) {
	parser := resourceids.NewParserFromResourceIdType(&ExecutionId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := ExecutionId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParseExecutionIDInsensitively parses 'input' case-insensitively into a ExecutionId
// note: this method should only be used for API response data and not user input
func ParseExecutionIDInsensitively(input string) (*ExecutionId, error) {
	parser := resourceids.NewParserFromResourceIdType(&ExecutionId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := ExecutionId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *ExecutionId) FromParseResult(input resourceids.ParseResult) error {
	var ok bool

	if id.SubscriptionId, ok = input.Parsed["subscriptionId"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", input)
	}

	if id.ResourceGroupName, ok = input.Parsed["resourceGroupName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", input)
	}

	if id.ExperimentName, ok = input.Parsed["experimentName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "experimentName", input)
	}

	if id.ExecutionId, ok = input.Parsed["executionId"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "executionId", input)
	}

	return nil
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
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Chaos/experiments/%s/executions/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.ExperimentName, id.ExecutionId)
}

// Segments returns a slice of Resource ID Segments which comprise this Execution ID
func (id ExecutionId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftChaos", "Microsoft.Chaos", "Microsoft.Chaos"),
		resourceids.StaticSegment("staticExperiments", "experiments", "experiments"),
		resourceids.UserSpecifiedSegment("experimentName", "experimentName"),
		resourceids.StaticSegment("staticExecutions", "executions", "executions"),
		resourceids.UserSpecifiedSegment("executionId", "executionId"),
	}
}

// String returns a human-readable description of this Execution ID
func (id ExecutionId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Experiment Name: %q", id.ExperimentName),
		fmt.Sprintf("Execution: %q", id.ExecutionId),
	}
	return fmt.Sprintf("Execution (%s)", strings.Join(components, "\n"))
}
