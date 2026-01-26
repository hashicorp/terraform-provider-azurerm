package rbacs

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ resourceids.ResourceId = &ScopeId{}

// ScopeId is a struct representing the Resource ID for a Scope
type ScopeId struct {
	BaseURI string
	Scope   string
}

// NewScopeID returns a new ScopeId struct
func NewScopeID(baseURI string, scope string) ScopeId {
	return ScopeId{
		BaseURI: strings.TrimSuffix(baseURI, "/"),
		Scope:   scope,
	}
}

// ParseScopeID parses 'input' into a ScopeId
func ParseScopeID(input string) (*ScopeId, error) {
	parser := resourceids.NewParserFromResourceIdType(&ScopeId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := ScopeId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParseScopeIDInsensitively parses 'input' case-insensitively into a ScopeId
// note: this method should only be used for API response data and not user input
func ParseScopeIDInsensitively(input string) (*ScopeId, error) {
	parser := resourceids.NewParserFromResourceIdType(&ScopeId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := ScopeId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *ScopeId) FromParseResult(input resourceids.ParseResult) error {
	var ok bool

	if id.BaseURI, ok = input.Parsed["baseURI"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "baseURI", input)
	}

	if id.Scope, ok = input.Parsed["scope"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "scope", input)
	}

	return nil
}

// ValidateScopeID checks that 'input' can be parsed as a Scope ID
func ValidateScopeID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseScopeID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Scope ID
func (id ScopeId) ID() string {
	fmtString := "%s/%s"
	return fmt.Sprintf(fmtString, strings.TrimSuffix(id.BaseURI, "/"), strings.TrimPrefix(id.Scope, "/"))
}

// Path returns the formatted Scope ID without the BaseURI
func (id ScopeId) Path() string {
	fmtString := "/%s"
	return fmt.Sprintf(fmtString, strings.TrimPrefix(id.Scope, "/"))
}

// PathElements returns the values of Scope ID Segments without the BaseURI
func (id ScopeId) PathElements() []any {
	return []any{strings.TrimPrefix(id.Scope, "/")}
}

// Segments returns a slice of Resource ID Segments which comprise this Scope ID
func (id ScopeId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.DataPlaneBaseURISegment("baseURI", "https://endpoint-url.example.com"),
		resourceids.ScopeSegment("scope", "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/some-resource-group"),
	}
}

// String returns a human-readable description of this Scope ID
func (id ScopeId) String() string {
	components := []string{
		fmt.Sprintf("Base URI: %q", id.BaseURI),
		fmt.Sprintf("Scope: %q", id.Scope),
	}
	return fmt.Sprintf("Scope (%s)", strings.Join(components, "\n"))
}
