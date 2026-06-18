package roleassignments

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ resourceids.ResourceId = &RoleAssignmentIdId{}

// RoleAssignmentIdId is a struct representing the Resource ID for a Role Assignment Id
type RoleAssignmentIdId struct {
	BaseURI          string
	RoleAssignmentId string
}

// NewRoleAssignmentIdID returns a new RoleAssignmentIdId struct
func NewRoleAssignmentIdID(baseURI string, roleAssignmentId string) RoleAssignmentIdId {
	return RoleAssignmentIdId{
		BaseURI:          strings.TrimSuffix(baseURI, "/"),
		RoleAssignmentId: roleAssignmentId,
	}
}

// ParseRoleAssignmentIdID parses 'input' into a RoleAssignmentIdId
func ParseRoleAssignmentIdID(input string) (*RoleAssignmentIdId, error) {
	parser := resourceids.NewParserFromResourceIdType(&RoleAssignmentIdId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := RoleAssignmentIdId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParseRoleAssignmentIdIDInsensitively parses 'input' case-insensitively into a RoleAssignmentIdId
// note: this method should only be used for API response data and not user input
func ParseRoleAssignmentIdIDInsensitively(input string) (*RoleAssignmentIdId, error) {
	parser := resourceids.NewParserFromResourceIdType(&RoleAssignmentIdId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := RoleAssignmentIdId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *RoleAssignmentIdId) FromParseResult(input resourceids.ParseResult) error {
	var ok bool

	if id.BaseURI, ok = input.Parsed["baseURI"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "baseURI", input)
	}

	if id.RoleAssignmentId, ok = input.Parsed["roleAssignmentId"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "roleAssignmentId", input)
	}

	return nil
}

// ValidateRoleAssignmentIdID checks that 'input' can be parsed as a Role Assignment Id ID
func ValidateRoleAssignmentIdID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseRoleAssignmentIdID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Role Assignment Id ID
func (id RoleAssignmentIdId) ID() string {
	fmtString := "%s/roleAssignments/%s"
	return fmt.Sprintf(fmtString, strings.TrimSuffix(id.BaseURI, "/"), strings.TrimPrefix(id.RoleAssignmentId, "/"))
}

// Path returns the formatted Role Assignment Id ID without the BaseURI
func (id RoleAssignmentIdId) Path() string {
	fmtString := "/roleAssignments/%s"
	return fmt.Sprintf(fmtString, strings.TrimPrefix(id.RoleAssignmentId, "/"))
}

// PathElements returns the values of Role Assignment Id ID Segments without the BaseURI
func (id RoleAssignmentIdId) PathElements() []any {
	return []any{strings.TrimPrefix(id.RoleAssignmentId, "/")}
}

// Segments returns a slice of Resource ID Segments which comprise this Role Assignment Id ID
func (id RoleAssignmentIdId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.DataPlaneBaseURISegment("baseURI", "https://endpoint-url.example.com"),
		resourceids.StaticSegment("staticRoleAssignments", "roleAssignments", "roleAssignments"),
		resourceids.ScopeSegment("roleAssignmentId", "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/some-resource-group"),
	}
}

// String returns a human-readable description of this Role Assignment Id ID
func (id RoleAssignmentIdId) String() string {
	components := []string{
		fmt.Sprintf("Base URI: %q", id.BaseURI),
		fmt.Sprintf("Role Assignment: %q", id.RoleAssignmentId),
	}
	return fmt.Sprintf("Role Assignment Id (%s)", strings.Join(components, "\n"))
}
