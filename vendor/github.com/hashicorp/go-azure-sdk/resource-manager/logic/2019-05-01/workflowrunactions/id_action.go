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
	recaser.RegisterResourceId(&ActionId{})
}

var _ resourceids.ResourceId = &ActionId{}

// ActionId is a struct representing the Resource ID for a Action
type ActionId struct {
	SubscriptionId    string
	ResourceGroupName string
	WorkflowName      string
	RunName           string
	ActionName        string
}

// NewActionID returns a new ActionId struct
func NewActionID(subscriptionId string, resourceGroupName string, workflowName string, runName string, actionName string) ActionId {
	return ActionId{
		SubscriptionId:    subscriptionId,
		ResourceGroupName: resourceGroupName,
		WorkflowName:      workflowName,
		RunName:           runName,
		ActionName:        actionName,
	}
}

// ParseActionID parses 'input' into a ActionId
func ParseActionID(input string) (*ActionId, error) {
	parser := resourceids.NewParserFromResourceIdType(&ActionId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := ActionId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParseActionIDInsensitively parses 'input' case-insensitively into a ActionId
// note: this method should only be used for API response data and not user input
func ParseActionIDInsensitively(input string) (*ActionId, error) {
	parser := resourceids.NewParserFromResourceIdType(&ActionId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := ActionId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *ActionId) FromParseResult(input resourceids.ParseResult) error {
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

	return nil
}

// ValidateActionID checks that 'input' can be parsed as a Action ID
func ValidateActionID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseActionID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Action ID
func (id ActionId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Logic/workflows/%s/runs/%s/actions/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.WorkflowName, id.RunName, id.ActionName)
}

// Segments returns a slice of Resource ID Segments which comprise this Action ID
func (id ActionId) Segments() []resourceids.Segment {
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
	}
}

// String returns a human-readable description of this Action ID
func (id ActionId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Workflow Name: %q", id.WorkflowName),
		fmt.Sprintf("Run Name: %q", id.RunName),
		fmt.Sprintf("Action Name: %q", id.ActionName),
	}
	return fmt.Sprintf("Action (%s)", strings.Join(components, "\n"))
}
