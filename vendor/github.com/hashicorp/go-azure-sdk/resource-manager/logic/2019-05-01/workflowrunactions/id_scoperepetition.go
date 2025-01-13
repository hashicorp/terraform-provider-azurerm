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
	recaser.RegisterResourceId(&ScopeRepetitionId{})
}

var _ resourceids.ResourceId = &ScopeRepetitionId{}

// ScopeRepetitionId is a struct representing the Resource ID for a Scope Repetition
type ScopeRepetitionId struct {
	SubscriptionId      string
	ResourceGroupName   string
	WorkflowName        string
	RunName             string
	ActionName          string
	ScopeRepetitionName string
}

// NewScopeRepetitionID returns a new ScopeRepetitionId struct
func NewScopeRepetitionID(subscriptionId string, resourceGroupName string, workflowName string, runName string, actionName string, scopeRepetitionName string) ScopeRepetitionId {
	return ScopeRepetitionId{
		SubscriptionId:      subscriptionId,
		ResourceGroupName:   resourceGroupName,
		WorkflowName:        workflowName,
		RunName:             runName,
		ActionName:          actionName,
		ScopeRepetitionName: scopeRepetitionName,
	}
}

// ParseScopeRepetitionID parses 'input' into a ScopeRepetitionId
func ParseScopeRepetitionID(input string) (*ScopeRepetitionId, error) {
	parser := resourceids.NewParserFromResourceIdType(&ScopeRepetitionId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := ScopeRepetitionId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParseScopeRepetitionIDInsensitively parses 'input' case-insensitively into a ScopeRepetitionId
// note: this method should only be used for API response data and not user input
func ParseScopeRepetitionIDInsensitively(input string) (*ScopeRepetitionId, error) {
	parser := resourceids.NewParserFromResourceIdType(&ScopeRepetitionId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := ScopeRepetitionId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *ScopeRepetitionId) FromParseResult(input resourceids.ParseResult) error {
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

	if id.ScopeRepetitionName, ok = input.Parsed["scopeRepetitionName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "scopeRepetitionName", input)
	}

	return nil
}

// ValidateScopeRepetitionID checks that 'input' can be parsed as a Scope Repetition ID
func ValidateScopeRepetitionID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseScopeRepetitionID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Scope Repetition ID
func (id ScopeRepetitionId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Logic/workflows/%s/runs/%s/actions/%s/scopeRepetitions/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.WorkflowName, id.RunName, id.ActionName, id.ScopeRepetitionName)
}

// Segments returns a slice of Resource ID Segments which comprise this Scope Repetition ID
func (id ScopeRepetitionId) Segments() []resourceids.Segment {
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
		resourceids.StaticSegment("staticScopeRepetitions", "scopeRepetitions", "scopeRepetitions"),
		resourceids.UserSpecifiedSegment("scopeRepetitionName", "scopeRepetitionName"),
	}
}

// String returns a human-readable description of this Scope Repetition ID
func (id ScopeRepetitionId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Workflow Name: %q", id.WorkflowName),
		fmt.Sprintf("Run Name: %q", id.RunName),
		fmt.Sprintf("Action Name: %q", id.ActionName),
		fmt.Sprintf("Scope Repetition Name: %q", id.ScopeRepetitionName),
	}
	return fmt.Sprintf("Scope Repetition (%s)", strings.Join(components, "\n"))
}
