package privateendpointconnections

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/recaser"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

func init() {
	recaser.RegisterResourceId(&ScopedPrivateEndpointConnectionId{})
}

var _ resourceids.ResourceId = &ScopedPrivateEndpointConnectionId{}

// ScopedPrivateEndpointConnectionId is a struct representing the Resource ID for a Scoped Private Endpoint Connection
type ScopedPrivateEndpointConnectionId struct {
	Scope                         string
	PrivateEndpointConnectionName string
}

// NewScopedPrivateEndpointConnectionID returns a new ScopedPrivateEndpointConnectionId struct
func NewScopedPrivateEndpointConnectionID(scope string, privateEndpointConnectionName string) ScopedPrivateEndpointConnectionId {
	return ScopedPrivateEndpointConnectionId{
		Scope:                         scope,
		PrivateEndpointConnectionName: privateEndpointConnectionName,
	}
}

// ParseScopedPrivateEndpointConnectionID parses 'input' into a ScopedPrivateEndpointConnectionId
func ParseScopedPrivateEndpointConnectionID(input string) (*ScopedPrivateEndpointConnectionId, error) {
	parser := resourceids.NewParserFromResourceIdType(&ScopedPrivateEndpointConnectionId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := ScopedPrivateEndpointConnectionId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParseScopedPrivateEndpointConnectionIDInsensitively parses 'input' case-insensitively into a ScopedPrivateEndpointConnectionId
// note: this method should only be used for API response data and not user input
func ParseScopedPrivateEndpointConnectionIDInsensitively(input string) (*ScopedPrivateEndpointConnectionId, error) {
	parser := resourceids.NewParserFromResourceIdType(&ScopedPrivateEndpointConnectionId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := ScopedPrivateEndpointConnectionId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *ScopedPrivateEndpointConnectionId) FromParseResult(input resourceids.ParseResult) error {
	var ok bool

	if id.Scope, ok = input.Parsed["scope"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "scope", input)
	}

	if id.PrivateEndpointConnectionName, ok = input.Parsed["privateEndpointConnectionName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "privateEndpointConnectionName", input)
	}

	return nil
}

// ValidateScopedPrivateEndpointConnectionID checks that 'input' can be parsed as a Scoped Private Endpoint Connection ID
func ValidateScopedPrivateEndpointConnectionID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseScopedPrivateEndpointConnectionID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Scoped Private Endpoint Connection ID
func (id ScopedPrivateEndpointConnectionId) ID() string {
	fmtString := "/%s/privateEndpointConnections/%s"
	return fmt.Sprintf(fmtString, strings.TrimPrefix(id.Scope, "/"), id.PrivateEndpointConnectionName)
}

// Segments returns a slice of Resource ID Segments which comprise this Scoped Private Endpoint Connection ID
func (id ScopedPrivateEndpointConnectionId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.ScopeSegment("scope", "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/some-resource-group"),
		resourceids.StaticSegment("staticPrivateEndpointConnections", "privateEndpointConnections", "privateEndpointConnections"),
		resourceids.UserSpecifiedSegment("privateEndpointConnectionName", "privateEndpointConnectionName"),
	}
}

// String returns a human-readable description of this Scoped Private Endpoint Connection ID
func (id ScopedPrivateEndpointConnectionId) String() string {
	components := []string{
		fmt.Sprintf("Scope: %q", id.Scope),
		fmt.Sprintf("Private Endpoint Connection Name: %q", id.PrivateEndpointConnectionName),
	}
	return fmt.Sprintf("Scoped Private Endpoint Connection (%s)", strings.Join(components, "\n"))
}
