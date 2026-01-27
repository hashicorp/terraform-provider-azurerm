package rbacs

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ resourceids.ResourceId = &ScopedRoleDefinitionId{}

// ScopedRoleDefinitionId is a struct representing the Resource ID for a Scoped Role Definition
type ScopedRoleDefinitionId struct {
	BaseURI            string
	Scope              string
	RoleDefinitionName string
}

// NewScopedRoleDefinitionID returns a new ScopedRoleDefinitionId struct
func NewScopedRoleDefinitionID(baseURI string, scope string, roleDefinitionName string) ScopedRoleDefinitionId {
	return ScopedRoleDefinitionId{
		BaseURI:            strings.TrimSuffix(baseURI, "/"),
		Scope:              scope,
		RoleDefinitionName: roleDefinitionName,
	}
}

// ParseScopedRoleDefinitionID parses 'input' into a ScopedRoleDefinitionId
func ParseScopedRoleDefinitionID(input string) (*ScopedRoleDefinitionId, error) {
	parser := resourceids.NewParserFromResourceIdType(&ScopedRoleDefinitionId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := ScopedRoleDefinitionId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParseScopedRoleDefinitionIDInsensitively parses 'input' case-insensitively into a ScopedRoleDefinitionId
// note: this method should only be used for API response data and not user input
func ParseScopedRoleDefinitionIDInsensitively(input string) (*ScopedRoleDefinitionId, error) {
	parser := resourceids.NewParserFromResourceIdType(&ScopedRoleDefinitionId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := ScopedRoleDefinitionId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *ScopedRoleDefinitionId) FromParseResult(input resourceids.ParseResult) error {
	var ok bool

	if id.BaseURI, ok = input.Parsed["baseURI"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "baseURI", input)
	}

	if id.Scope, ok = input.Parsed["scope"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "scope", input)
	}

	if id.RoleDefinitionName, ok = input.Parsed["roleDefinitionName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "roleDefinitionName", input)
	}

	return nil
}

// ValidateScopedRoleDefinitionID checks that 'input' can be parsed as a Scoped Role Definition ID
func ValidateScopedRoleDefinitionID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseScopedRoleDefinitionID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Scoped Role Definition ID
func (id ScopedRoleDefinitionId) ID() string {
	fmtString := "%s/%s/providers/Microsoft.Authorization/roleDefinitions/%s"
	return fmt.Sprintf(fmtString, id.BaseURI, strings.TrimPrefix(id.Scope, "/"), id.RoleDefinitionName)
}

// Path returns the formatted Scoped Role Definition ID without the BaseURI
func (id ScopedRoleDefinitionId) Path() string {
	fmtString := "/%s/providers/Microsoft.Authorization/roleDefinitions/%s"
	return fmt.Sprintf(fmtString, strings.TrimPrefix(id.Scope, "/"), id.RoleDefinitionName)
}

// PathElements returns the values of Scoped Role Definition ID Segments without the BaseURI
func (id ScopedRoleDefinitionId) PathElements() []any {
	return []any{strings.TrimPrefix(id.Scope, "/"), id.RoleDefinitionName}
}

// Segments returns a slice of Resource ID Segments which comprise this Scoped Role Definition ID
func (id ScopedRoleDefinitionId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.DataPlaneBaseURISegment("baseURI", "https://endpoint_url"),
		resourceids.ScopeSegment("scope", "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/some-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftAuthorization", "Microsoft.Authorization", "Microsoft.Authorization"),
		resourceids.StaticSegment("staticRoleDefinitions", "roleDefinitions", "roleDefinitions"),
		resourceids.UserSpecifiedSegment("roleDefinitionName", "roleDefinitionName"),
	}
}

// String returns a human-readable description of this Scoped Role Definition ID
func (id ScopedRoleDefinitionId) String() string {
	components := []string{
		fmt.Sprintf("Base URI: %q", id.BaseURI),
		fmt.Sprintf("Scope: %q", id.Scope),
		fmt.Sprintf("Role Definition Name: %q", id.RoleDefinitionName),
	}
	return fmt.Sprintf("Scoped Role Definition (%s)", strings.Join(components, "\n"))
}
