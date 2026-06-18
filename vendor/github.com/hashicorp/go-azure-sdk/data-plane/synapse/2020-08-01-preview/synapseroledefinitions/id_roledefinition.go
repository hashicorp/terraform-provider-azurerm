package synapseroledefinitions

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ resourceids.ResourceId = &RoleDefinitionId{}

// RoleDefinitionId is a struct representing the Resource ID for a Role Definition
type RoleDefinitionId struct {
	BaseURI          string
	RoleDefinitionId string
}

// NewRoleDefinitionID returns a new RoleDefinitionId struct
func NewRoleDefinitionID(baseURI string, roleDefinitionId string) RoleDefinitionId {
	return RoleDefinitionId{
		BaseURI:          strings.TrimSuffix(baseURI, "/"),
		RoleDefinitionId: roleDefinitionId,
	}
}

// ParseRoleDefinitionID parses 'input' into a RoleDefinitionId
func ParseRoleDefinitionID(input string) (*RoleDefinitionId, error) {
	parser := resourceids.NewParserFromResourceIdType(&RoleDefinitionId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := RoleDefinitionId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParseRoleDefinitionIDInsensitively parses 'input' case-insensitively into a RoleDefinitionId
// note: this method should only be used for API response data and not user input
func ParseRoleDefinitionIDInsensitively(input string) (*RoleDefinitionId, error) {
	parser := resourceids.NewParserFromResourceIdType(&RoleDefinitionId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := RoleDefinitionId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *RoleDefinitionId) FromParseResult(input resourceids.ParseResult) error {
	var ok bool

	if id.BaseURI, ok = input.Parsed["baseURI"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "baseURI", input)
	}

	if id.RoleDefinitionId, ok = input.Parsed["roleDefinitionId"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "roleDefinitionId", input)
	}

	return nil
}

// ValidateRoleDefinitionID checks that 'input' can be parsed as a Role Definition ID
func ValidateRoleDefinitionID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseRoleDefinitionID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Role Definition ID
func (id RoleDefinitionId) ID() string {
	fmtString := "%s/roleDefinitions/%s"
	return fmt.Sprintf(fmtString, strings.TrimSuffix(id.BaseURI, "/"), id.RoleDefinitionId)
}

// Path returns the formatted Role Definition ID without the BaseURI
func (id RoleDefinitionId) Path() string {
	fmtString := "/roleDefinitions/%s"
	return fmt.Sprintf(fmtString, id.RoleDefinitionId)
}

// PathElements returns the values of Role Definition ID Segments without the BaseURI
func (id RoleDefinitionId) PathElements() []any {
	return []any{id.RoleDefinitionId}
}

// Segments returns a slice of Resource ID Segments which comprise this Role Definition ID
func (id RoleDefinitionId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.DataPlaneBaseURISegment("baseURI", "https://endpoint-url.example.com"),
		resourceids.StaticSegment("staticRoleDefinitions", "roleDefinitions", "roleDefinitions"),
		resourceids.UserSpecifiedSegment("roleDefinitionId", "roleDefinitionId"),
	}
}

// String returns a human-readable description of this Role Definition ID
func (id RoleDefinitionId) String() string {
	components := []string{
		fmt.Sprintf("Base URI: %q", id.BaseURI),
		fmt.Sprintf("Role Definition: %q", id.RoleDefinitionId),
	}
	return fmt.Sprintf("Role Definition (%s)", strings.Join(components, "\n"))
}
