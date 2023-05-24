package publishedblueprint

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ resourceids.ResourceId = ScopedVersionId{}

// ScopedVersionId is a struct representing the Resource ID for a Scoped Version
type ScopedVersionId struct {
	ResourceScope string
	BlueprintName string
	VersionId     string
}

// NewScopedVersionID returns a new ScopedVersionId struct
func NewScopedVersionID(resourceScope string, blueprintName string, versionId string) ScopedVersionId {
	return ScopedVersionId{
		ResourceScope: resourceScope,
		BlueprintName: blueprintName,
		VersionId:     versionId,
	}
}

// ParseScopedVersionID parses 'input' into a ScopedVersionId
func ParseScopedVersionID(input string) (*ScopedVersionId, error) {
	parser := resourceids.NewParserFromResourceIdType(ScopedVersionId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := ScopedVersionId{}

	if id.ResourceScope, ok = parsed.Parsed["resourceScope"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "resourceScope", *parsed)
	}

	if id.BlueprintName, ok = parsed.Parsed["blueprintName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "blueprintName", *parsed)
	}

	if id.VersionId, ok = parsed.Parsed["versionId"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "versionId", *parsed)
	}

	return &id, nil
}

// ParseScopedVersionIDInsensitively parses 'input' case-insensitively into a ScopedVersionId
// note: this method should only be used for API response data and not user input
func ParseScopedVersionIDInsensitively(input string) (*ScopedVersionId, error) {
	parser := resourceids.NewParserFromResourceIdType(ScopedVersionId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := ScopedVersionId{}

	if id.ResourceScope, ok = parsed.Parsed["resourceScope"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "resourceScope", *parsed)
	}

	if id.BlueprintName, ok = parsed.Parsed["blueprintName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "blueprintName", *parsed)
	}

	if id.VersionId, ok = parsed.Parsed["versionId"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "versionId", *parsed)
	}

	return &id, nil
}

// ValidateScopedVersionID checks that 'input' can be parsed as a Scoped Version ID
func ValidateScopedVersionID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseScopedVersionID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Scoped Version ID
func (id ScopedVersionId) ID() string {
	fmtString := "/%s/providers/Microsoft.Blueprint/blueprints/%s/versions/%s"
	return fmt.Sprintf(fmtString, strings.TrimPrefix(id.ResourceScope, "/"), id.BlueprintName, id.VersionId)
}

// Segments returns a slice of Resource ID Segments which comprise this Scoped Version ID
func (id ScopedVersionId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.ScopeSegment("resourceScope", "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/some-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftBlueprint", "Microsoft.Blueprint", "Microsoft.Blueprint"),
		resourceids.StaticSegment("staticBlueprints", "blueprints", "blueprints"),
		resourceids.UserSpecifiedSegment("blueprintName", "blueprintValue"),
		resourceids.StaticSegment("staticVersions", "versions", "versions"),
		resourceids.UserSpecifiedSegment("versionId", "versionIdValue"),
	}
}

// String returns a human-readable description of this Scoped Version ID
func (id ScopedVersionId) String() string {
	components := []string{
		fmt.Sprintf("Resource Scope: %q", id.ResourceScope),
		fmt.Sprintf("Blueprint Name: %q", id.BlueprintName),
		fmt.Sprintf("Version: %q", id.VersionId),
	}
	return fmt.Sprintf("Scoped Version (%s)", strings.Join(components, "\n"))
}
