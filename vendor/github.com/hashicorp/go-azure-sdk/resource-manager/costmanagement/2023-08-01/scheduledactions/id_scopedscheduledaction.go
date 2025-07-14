package scheduledactions

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/recaser"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

func init() {
	recaser.RegisterResourceId(&ScopedScheduledActionId{})
}

var _ resourceids.ResourceId = &ScopedScheduledActionId{}

// ScopedScheduledActionId is a struct representing the Resource ID for a Scoped Scheduled Action
type ScopedScheduledActionId struct {
	Scope               string
	ScheduledActionName string
}

// NewScopedScheduledActionID returns a new ScopedScheduledActionId struct
func NewScopedScheduledActionID(scope string, scheduledActionName string) ScopedScheduledActionId {
	return ScopedScheduledActionId{
		Scope:               scope,
		ScheduledActionName: scheduledActionName,
	}
}

// ParseScopedScheduledActionID parses 'input' into a ScopedScheduledActionId
func ParseScopedScheduledActionID(input string) (*ScopedScheduledActionId, error) {
	parser := resourceids.NewParserFromResourceIdType(&ScopedScheduledActionId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := ScopedScheduledActionId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParseScopedScheduledActionIDInsensitively parses 'input' case-insensitively into a ScopedScheduledActionId
// note: this method should only be used for API response data and not user input
func ParseScopedScheduledActionIDInsensitively(input string) (*ScopedScheduledActionId, error) {
	parser := resourceids.NewParserFromResourceIdType(&ScopedScheduledActionId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := ScopedScheduledActionId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *ScopedScheduledActionId) FromParseResult(input resourceids.ParseResult) error {
	var ok bool

	if id.Scope, ok = input.Parsed["scope"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "scope", input)
	}

	if id.ScheduledActionName, ok = input.Parsed["scheduledActionName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "scheduledActionName", input)
	}

	return nil
}

// ValidateScopedScheduledActionID checks that 'input' can be parsed as a Scoped Scheduled Action ID
func ValidateScopedScheduledActionID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseScopedScheduledActionID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Scoped Scheduled Action ID
func (id ScopedScheduledActionId) ID() string {
	fmtString := "/%s/providers/Microsoft.CostManagement/scheduledActions/%s"
	return fmt.Sprintf(fmtString, strings.TrimPrefix(id.Scope, "/"), id.ScheduledActionName)
}

// Segments returns a slice of Resource ID Segments which comprise this Scoped Scheduled Action ID
func (id ScopedScheduledActionId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.ScopeSegment("scope", "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/some-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftCostManagement", "Microsoft.CostManagement", "Microsoft.CostManagement"),
		resourceids.StaticSegment("staticScheduledActions", "scheduledActions", "scheduledActions"),
		resourceids.UserSpecifiedSegment("scheduledActionName", "scheduledActionName"),
	}
}

// String returns a human-readable description of this Scoped Scheduled Action ID
func (id ScopedScheduledActionId) String() string {
	components := []string{
		fmt.Sprintf("Scope: %q", id.Scope),
		fmt.Sprintf("Scheduled Action Name: %q", id.ScheduledActionName),
	}
	return fmt.Sprintf("Scoped Scheduled Action (%s)", strings.Join(components, "\n"))
}
