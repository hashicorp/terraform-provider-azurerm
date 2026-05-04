package links

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/recaser"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

func init() {
	recaser.RegisterResourceId(&ScopedLinkerId{})
}

var _ resourceids.ResourceId = &ScopedLinkerId{}

// ScopedLinkerId is a struct representing the Resource ID for a Scoped Linker
type ScopedLinkerId struct {
	ResourceUri string
	LinkerName  string
}

// NewScopedLinkerID returns a new ScopedLinkerId struct
func NewScopedLinkerID(resourceUri string, linkerName string) ScopedLinkerId {
	return ScopedLinkerId{
		ResourceUri: resourceUri,
		LinkerName:  linkerName,
	}
}

// ParseScopedLinkerID parses 'input' into a ScopedLinkerId
func ParseScopedLinkerID(input string) (*ScopedLinkerId, error) {
	parser := resourceids.NewParserFromResourceIdType(&ScopedLinkerId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := ScopedLinkerId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParseScopedLinkerIDInsensitively parses 'input' case-insensitively into a ScopedLinkerId
// note: this method should only be used for API response data and not user input
func ParseScopedLinkerIDInsensitively(input string) (*ScopedLinkerId, error) {
	parser := resourceids.NewParserFromResourceIdType(&ScopedLinkerId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := ScopedLinkerId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *ScopedLinkerId) FromParseResult(input resourceids.ParseResult) error {
	var ok bool

	if id.ResourceUri, ok = input.Parsed["resourceUri"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "resourceUri", input)
	}

	if id.LinkerName, ok = input.Parsed["linkerName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "linkerName", input)
	}

	return nil
}

// ValidateScopedLinkerID checks that 'input' can be parsed as a Scoped Linker ID
func ValidateScopedLinkerID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseScopedLinkerID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Scoped Linker ID
func (id ScopedLinkerId) ID() string {
	fmtString := "/%s/providers/Microsoft.ServiceLinker/linkers/%s"
	return fmt.Sprintf(fmtString, strings.TrimPrefix(id.ResourceUri, "/"), id.LinkerName)
}

// Segments returns a slice of Resource ID Segments which comprise this Scoped Linker ID
func (id ScopedLinkerId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.ScopeSegment("resourceUri", "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/some-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftServiceLinker", "Microsoft.ServiceLinker", "Microsoft.ServiceLinker"),
		resourceids.StaticSegment("staticLinkers", "linkers", "linkers"),
		resourceids.UserSpecifiedSegment("linkerName", "linkerName"),
	}
}

// String returns a human-readable description of this Scoped Linker ID
func (id ScopedLinkerId) String() string {
	components := []string{
		fmt.Sprintf("Resource Uri: %q", id.ResourceUri),
		fmt.Sprintf("Linker Name: %q", id.LinkerName),
	}
	return fmt.Sprintf("Scoped Linker (%s)", strings.Join(components, "\n"))
}
