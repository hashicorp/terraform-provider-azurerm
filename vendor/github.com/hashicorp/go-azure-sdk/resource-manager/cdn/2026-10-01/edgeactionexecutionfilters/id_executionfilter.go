package edgeactionexecutionfilters

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/recaser"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

func init() {
	recaser.RegisterResourceId(&ExecutionFilterId{})
}

var _ resourceids.ResourceId = &ExecutionFilterId{}

// ExecutionFilterId is a struct representing the Resource ID for a Execution Filter
type ExecutionFilterId struct {
	SubscriptionId      string
	ResourceGroupName   string
	EdgeActionName      string
	ExecutionFilterName string
}

// NewExecutionFilterID returns a new ExecutionFilterId struct
func NewExecutionFilterID(subscriptionId string, resourceGroupName string, edgeActionName string, executionFilterName string) ExecutionFilterId {
	return ExecutionFilterId{
		SubscriptionId:      subscriptionId,
		ResourceGroupName:   resourceGroupName,
		EdgeActionName:      edgeActionName,
		ExecutionFilterName: executionFilterName,
	}
}

// ParseExecutionFilterID parses 'input' into a ExecutionFilterId
func ParseExecutionFilterID(input string) (*ExecutionFilterId, error) {
	parser := resourceids.NewParserFromResourceIdType(&ExecutionFilterId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := ExecutionFilterId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParseExecutionFilterIDInsensitively parses 'input' case-insensitively into a ExecutionFilterId
// note: this method should only be used for API response data and not user input
func ParseExecutionFilterIDInsensitively(input string) (*ExecutionFilterId, error) {
	parser := resourceids.NewParserFromResourceIdType(&ExecutionFilterId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := ExecutionFilterId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *ExecutionFilterId) FromParseResult(input resourceids.ParseResult) error {
	var ok bool

	if id.SubscriptionId, ok = input.Parsed["subscriptionId"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", input)
	}

	if id.ResourceGroupName, ok = input.Parsed["resourceGroupName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", input)
	}

	if id.EdgeActionName, ok = input.Parsed["edgeActionName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "edgeActionName", input)
	}

	if id.ExecutionFilterName, ok = input.Parsed["executionFilterName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "executionFilterName", input)
	}

	return nil
}

// ValidateExecutionFilterID checks that 'input' can be parsed as a Execution Filter ID
func ValidateExecutionFilterID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseExecutionFilterID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Execution Filter ID
func (id ExecutionFilterId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Cdn/edgeActions/%s/executionFilters/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.EdgeActionName, id.ExecutionFilterName)
}

// Segments returns a slice of Resource ID Segments which comprise this Execution Filter ID
func (id ExecutionFilterId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftCdn", "Microsoft.Cdn", "Microsoft.Cdn"),
		resourceids.StaticSegment("staticEdgeActions", "edgeActions", "edgeActions"),
		resourceids.UserSpecifiedSegment("edgeActionName", "edgeActionName"),
		resourceids.StaticSegment("staticExecutionFilters", "executionFilters", "executionFilters"),
		resourceids.UserSpecifiedSegment("executionFilterName", "executionFilterName"),
	}
}

// String returns a human-readable description of this Execution Filter ID
func (id ExecutionFilterId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Edge Action Name: %q", id.EdgeActionName),
		fmt.Sprintf("Execution Filter Name: %q", id.ExecutionFilterName),
	}
	return fmt.Sprintf("Execution Filter (%s)", strings.Join(components, "\n"))
}
