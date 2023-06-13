package workflowrunactions

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ resourceids.ResourceId = RepetitionRequestHistoryId{}

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
	parser := resourceids.NewParserFromResourceIdType(RepetitionRequestHistoryId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := RepetitionRequestHistoryId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", *parsed)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", *parsed)
	}

	if id.WorkflowName, ok = parsed.Parsed["workflowName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "workflowName", *parsed)
	}

	if id.RunName, ok = parsed.Parsed["runName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "runName", *parsed)
	}

	if id.ActionName, ok = parsed.Parsed["actionName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "actionName", *parsed)
	}

	if id.RepetitionName, ok = parsed.Parsed["repetitionName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "repetitionName", *parsed)
	}

	if id.RequestHistoryName, ok = parsed.Parsed["requestHistoryName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "requestHistoryName", *parsed)
	}

	return &id, nil
}

// ParseRepetitionRequestHistoryIDInsensitively parses 'input' case-insensitively into a RepetitionRequestHistoryId
// note: this method should only be used for API response data and not user input
func ParseRepetitionRequestHistoryIDInsensitively(input string) (*RepetitionRequestHistoryId, error) {
	parser := resourceids.NewParserFromResourceIdType(RepetitionRequestHistoryId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := RepetitionRequestHistoryId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", *parsed)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", *parsed)
	}

	if id.WorkflowName, ok = parsed.Parsed["workflowName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "workflowName", *parsed)
	}

	if id.RunName, ok = parsed.Parsed["runName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "runName", *parsed)
	}

	if id.ActionName, ok = parsed.Parsed["actionName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "actionName", *parsed)
	}

	if id.RepetitionName, ok = parsed.Parsed["repetitionName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "repetitionName", *parsed)
	}

	if id.RequestHistoryName, ok = parsed.Parsed["requestHistoryName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "requestHistoryName", *parsed)
	}

	return &id, nil
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
		resourceids.UserSpecifiedSegment("workflowName", "workflowValue"),
		resourceids.StaticSegment("staticRuns", "runs", "runs"),
		resourceids.UserSpecifiedSegment("runName", "runValue"),
		resourceids.StaticSegment("staticActions", "actions", "actions"),
		resourceids.UserSpecifiedSegment("actionName", "actionValue"),
		resourceids.StaticSegment("staticRepetitions", "repetitions", "repetitions"),
		resourceids.UserSpecifiedSegment("repetitionName", "repetitionValue"),
		resourceids.StaticSegment("staticRequestHistories", "requestHistories", "requestHistories"),
		resourceids.UserSpecifiedSegment("requestHistoryName", "requestHistoryValue"),
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
