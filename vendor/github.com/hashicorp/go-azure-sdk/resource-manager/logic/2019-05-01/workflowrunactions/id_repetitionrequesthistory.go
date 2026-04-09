package workflowrunactions

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/recaser"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

func init() {
	recaser.RegisterResourceId(&RepetitionRequestHistoryId{})
}

var _ resourceids.ResourceId = &RepetitionRequestHistoryId{}

// RepetitionRequestHistoryId is a struct representing the Resource ID for a Repetition Request History
type RepetitionRequestHistoryId struct {
	SubscriptionId     string
	ResourceGroupName  string
	WorkflowName       string
	RunName            string
	ActionName         string
	RepetitionName     string
	RequestHistoryName string
}

// NewRepetitionRequestHistoryID returns a new RepetitionRequestHistoryId struct
func NewRepetitionRequestHistoryID(subscriptionId string, resourceGroupName string, workflowName string, runName string, actionName string, repetitionName string, requestHistoryName string) RepetitionRequestHistoryId {
	return RepetitionRequestHistoryId{
		SubscriptionId:     subscriptionId,
		ResourceGroupName:  resourceGroupName,
		WorkflowName:       workflowName,
		RunName:            runName,
		ActionName:         actionName,
		RepetitionName:     repetitionName,
		RequestHistoryName: requestHistoryName,
	}
}

// ParseRepetitionRequestHistoryID parses 'input' into a RepetitionRequestHistoryId
func ParseRepetitionRequestHistoryID(input string) (*RepetitionRequestHistoryId, error) {
	parser := resourceids.NewParserFromResourceIdType(&RepetitionRequestHistoryId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := RepetitionRequestHistoryId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParseRepetitionRequestHistoryIDInsensitively parses 'input' case-insensitively into a RepetitionRequestHistoryId
// note: this method should only be used for API response data and not user input
func ParseRepetitionRequestHistoryIDInsensitively(input string) (*RepetitionRequestHistoryId, error) {
	parser := resourceids.NewParserFromResourceIdType(&RepetitionRequestHistoryId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := RepetitionRequestHistoryId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *RepetitionRequestHistoryId) FromParseResult(input resourceids.ParseResult) error {
	var ok bool

	if id.SubscriptionId, ok = input.Parsed["subscriptionId"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", input)
	}

	if id.ResourceGroupName, ok = input.Parsed["resourceGroupName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", input)
	}

	if id.WorkflowName, ok = input.Parsed["workflowName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "workflowName", input)
	}

	if id.RunName, ok = input.Parsed["runName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "runName", input)
	}

	if id.ActionName, ok = input.Parsed["actionName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "actionName", input)
	}

	if id.RepetitionName, ok = input.Parsed["repetitionName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "repetitionName", input)
	}

	if id.RequestHistoryName, ok = input.Parsed["requestHistoryName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "requestHistoryName", input)
	}

	return nil
}

// ValidateRepetitionRequestHistoryID checks that 'input' can be parsed as a Repetition Request History ID
func ValidateRepetitionRequestHistoryID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseRepetitionRequestHistoryID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Repetition Request History ID
func (id RepetitionRequestHistoryId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Logic/workflows/%s/runs/%s/actions/%s/repetitions/%s/requestHistories/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.WorkflowName, id.RunName, id.ActionName, id.RepetitionName, id.RequestHistoryName)
}

// Segments returns a slice of Resource ID Segments which comprise this Repetition Request History ID
func (id RepetitionRequestHistoryId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftLogic", "Microsoft.Logic", "Microsoft.Logic"),
		resourceids.StaticSegment("staticWorkflows", "workflows", "workflows"),
		resourceids.UserSpecifiedSegment("workflowName", "workflowName"),
		resourceids.StaticSegment("staticRuns", "runs", "runs"),
		resourceids.UserSpecifiedSegment("runName", "runName"),
		resourceids.StaticSegment("staticActions", "actions", "actions"),
		resourceids.UserSpecifiedSegment("actionName", "actionName"),
		resourceids.StaticSegment("staticRepetitions", "repetitions", "repetitions"),
		resourceids.UserSpecifiedSegment("repetitionName", "repetitionName"),
		resourceids.StaticSegment("staticRequestHistories", "requestHistories", "requestHistories"),
		resourceids.UserSpecifiedSegment("requestHistoryName", "requestHistoryName"),
	}
}

// String returns a human-readable description of this Repetition Request History ID
func (id RepetitionRequestHistoryId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Workflow Name: %q", id.WorkflowName),
		fmt.Sprintf("Run Name: %q", id.RunName),
		fmt.Sprintf("Action Name: %q", id.ActionName),
		fmt.Sprintf("Repetition Name: %q", id.RepetitionName),
		fmt.Sprintf("Request History Name: %q", id.RequestHistoryName),
	}
	return fmt.Sprintf("Repetition Request History (%s)", strings.Join(components, "\n"))
}
