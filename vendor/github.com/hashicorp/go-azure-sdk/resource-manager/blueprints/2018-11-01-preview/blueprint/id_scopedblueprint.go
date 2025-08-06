package blueprint

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/recaser"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

func init() {
	recaser.RegisterResourceId(&ScopedBlueprintId{})
}

var _ resourceids.ResourceId = &ScopedBlueprintId{}

// ScopedBlueprintId is a struct representing the Resource ID for a Scoped Blueprint
type ScopedBlueprintId struct {
	ResourceScope string
	BlueprintName string
}

// NewScopedBlueprintID returns a new ScopedBlueprintId struct
func NewScopedBlueprintID(resourceScope string, blueprintName string) ScopedBlueprintId {
	return ScopedBlueprintId{
		ResourceScope: resourceScope,
		BlueprintName: blueprintName,
	}
}

// ParseScopedBlueprintID parses 'input' into a ScopedBlueprintId
func ParseScopedBlueprintID(input string) (*ScopedBlueprintId, error) {
	parser := resourceids.NewParserFromResourceIdType(&ScopedBlueprintId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := ScopedBlueprintId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParseScopedBlueprintIDInsensitively parses 'input' case-insensitively into a ScopedBlueprintId
// note: this method should only be used for API response data and not user input
func ParseScopedBlueprintIDInsensitively(input string) (*ScopedBlueprintId, error) {
	parser := resourceids.NewParserFromResourceIdType(&ScopedBlueprintId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := ScopedBlueprintId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *ScopedBlueprintId) FromParseResult(input resourceids.ParseResult) error {
	var ok bool

	if id.ResourceScope, ok = input.Parsed["resourceScope"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "resourceScope", input)
	}

	if id.BlueprintName, ok = input.Parsed["blueprintName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "blueprintName", input)
	}

	return nil
}

// ValidateScopedBlueprintID checks that 'input' can be parsed as a Scoped Blueprint ID
func ValidateScopedBlueprintID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseScopedBlueprintID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Scoped Blueprint ID
func (id ScopedBlueprintId) ID() string {
	fmtString := "/%s/providers/Microsoft.Blueprint/blueprints/%s"
	return fmt.Sprintf(fmtString, strings.TrimPrefix(id.ResourceScope, "/"), id.BlueprintName)
}

// Segments returns a slice of Resource ID Segments which comprise this Scoped Blueprint ID
func (id ScopedBlueprintId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.ScopeSegment("resourceScope", "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/some-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftBlueprint", "Microsoft.Blueprint", "Microsoft.Blueprint"),
		resourceids.StaticSegment("staticBlueprints", "blueprints", "blueprints"),
		resourceids.UserSpecifiedSegment("blueprintName", "blueprintName"),
	}
}

// String returns a human-readable description of this Scoped Blueprint ID
func (id ScopedBlueprintId) String() string {
	components := []string{
		fmt.Sprintf("Resource Scope: %q", id.ResourceScope),
		fmt.Sprintf("Blueprint Name: %q", id.BlueprintName),
	}
	return fmt.Sprintf("Scoped Blueprint (%s)", strings.Join(components, "\n"))
}
