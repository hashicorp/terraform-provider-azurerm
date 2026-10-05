package assignment

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/recaser"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

func init() {
	recaser.RegisterResourceId(&ScopedBlueprintAssignmentId{})
}

var _ resourceids.ResourceId = &ScopedBlueprintAssignmentId{}

// ScopedBlueprintAssignmentId is a struct representing the Resource ID for a Scoped Blueprint Assignment
type ScopedBlueprintAssignmentId struct {
	ResourceScope           string
	BlueprintAssignmentName string
}

// NewScopedBlueprintAssignmentID returns a new ScopedBlueprintAssignmentId struct
func NewScopedBlueprintAssignmentID(resourceScope string, blueprintAssignmentName string) ScopedBlueprintAssignmentId {
	return ScopedBlueprintAssignmentId{
		ResourceScope:           resourceScope,
		BlueprintAssignmentName: blueprintAssignmentName,
	}
}

// ParseScopedBlueprintAssignmentID parses 'input' into a ScopedBlueprintAssignmentId
func ParseScopedBlueprintAssignmentID(input string) (*ScopedBlueprintAssignmentId, error) {
	parser := resourceids.NewParserFromResourceIdType(&ScopedBlueprintAssignmentId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := ScopedBlueprintAssignmentId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParseScopedBlueprintAssignmentIDInsensitively parses 'input' case-insensitively into a ScopedBlueprintAssignmentId
// note: this method should only be used for API response data and not user input
func ParseScopedBlueprintAssignmentIDInsensitively(input string) (*ScopedBlueprintAssignmentId, error) {
	parser := resourceids.NewParserFromResourceIdType(&ScopedBlueprintAssignmentId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := ScopedBlueprintAssignmentId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *ScopedBlueprintAssignmentId) FromParseResult(input resourceids.ParseResult) error {
	var ok bool

	if id.ResourceScope, ok = input.Parsed["resourceScope"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "resourceScope", input)
	}

	if id.BlueprintAssignmentName, ok = input.Parsed["blueprintAssignmentName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "blueprintAssignmentName", input)
	}

	return nil
}

// ValidateScopedBlueprintAssignmentID checks that 'input' can be parsed as a Scoped Blueprint Assignment ID
func ValidateScopedBlueprintAssignmentID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseScopedBlueprintAssignmentID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Scoped Blueprint Assignment ID
func (id ScopedBlueprintAssignmentId) ID() string {
	fmtString := "/%s/providers/Microsoft.Blueprint/blueprintAssignments/%s"
	return fmt.Sprintf(fmtString, strings.TrimPrefix(id.ResourceScope, "/"), id.BlueprintAssignmentName)
}

// Segments returns a slice of Resource ID Segments which comprise this Scoped Blueprint Assignment ID
func (id ScopedBlueprintAssignmentId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.ScopeSegment("resourceScope", "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/some-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftBlueprint", "Microsoft.Blueprint", "Microsoft.Blueprint"),
		resourceids.StaticSegment("staticBlueprintAssignments", "blueprintAssignments", "blueprintAssignments"),
		resourceids.UserSpecifiedSegment("blueprintAssignmentName", "blueprintAssignmentName"),
	}
}

// String returns a human-readable description of this Scoped Blueprint Assignment ID
func (id ScopedBlueprintAssignmentId) String() string {
	components := []string{
		fmt.Sprintf("Resource Scope: %q", id.ResourceScope),
		fmt.Sprintf("Blueprint Assignment Name: %q", id.BlueprintAssignmentName),
	}
	return fmt.Sprintf("Scoped Blueprint Assignment (%s)", strings.Join(components, "\n"))
}
