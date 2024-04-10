package remediations

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ resourceids.ResourceId = &ManagementGroupId{}

// ManagementGroupId is a struct representing the Resource ID for a Management Group
type ManagementGroupId struct {
	ManagementGroupId string
}

// NewManagementGroupID returns a new ManagementGroupId struct
func NewManagementGroupID(managementGroupId string) ManagementGroupId {
	return ManagementGroupId{
		ManagementGroupId: managementGroupId,
	}
}

// ParseManagementGroupID parses 'input' into a ManagementGroupId
func ParseManagementGroupID(input string) (*ManagementGroupId, error) {
	parser := resourceids.NewParserFromResourceIdType(&ManagementGroupId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := ManagementGroupId{}
	if err := id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParseManagementGroupIDInsensitively parses 'input' case-insensitively into a ManagementGroupId
// note: this method should only be used for API response data and not user input
func ParseManagementGroupIDInsensitively(input string) (*ManagementGroupId, error) {
	parser := resourceids.NewParserFromResourceIdType(&ManagementGroupId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := ManagementGroupId{}
	if err := id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *ManagementGroupId) FromParseResult(input resourceids.ParseResult) error {
	var ok bool

	if id.ManagementGroupId, ok = input.Parsed["managementGroupId"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "managementGroupId", input)
	}

	return nil
}

// ValidateManagementGroupID checks that 'input' can be parsed as a Management Group ID
func ValidateManagementGroupID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseManagementGroupID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Management Group ID
func (id ManagementGroupId) ID() string {
	fmtString := "/providers/Microsoft.Management/managementGroups/%s"
	return fmt.Sprintf(fmtString, id.ManagementGroupId)
}

// Segments returns a slice of Resource ID Segments which comprise this Management Group ID
func (id ManagementGroupId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.StaticSegment("managementGroupsNamespace", "Microsoft.Management", "Microsoft.Management"),
		resourceids.StaticSegment("staticManagementGroups", "managementGroups", "managementGroups"),
		resourceids.UserSpecifiedSegment("managementGroupId", "managementGroupIdValue"),
	}
}

// String returns a human-readable description of this Management Group ID
func (id ManagementGroupId) String() string {
	components := []string{
		fmt.Sprintf("Management Group: %q", id.ManagementGroupId),
	}
	return fmt.Sprintf("Management Group (%s)", strings.Join(components, "\n"))
}
