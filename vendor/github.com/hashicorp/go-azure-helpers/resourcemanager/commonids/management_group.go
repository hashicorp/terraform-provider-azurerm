// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package commonids

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

var _ resourceids.ResourceId = ManagementGroupId{}

// ManagementGroupId is a struct representing the Resource ID for a Management Group
type ManagementGroupId struct {
	GroupId string
}

// NewManagementGroupID returns a new ManagementGroupId struct
func NewManagementGroupID(groupId string) ManagementGroupId {
	return ManagementGroupId{
		GroupId: groupId,
	}
}

// ParseManagementGroupID parses 'input' into a ManagementGroupId
func ParseManagementGroupID(input string) (*ManagementGroupId, error) {
	parser := resourceids.NewParserFromResourceIdType(ManagementGroupId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := ManagementGroupId{}

	if id.GroupId, ok = parsed.Parsed["groupId"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "groupId", *parsed)
	}

	return &id, nil
}

// ParseManagementGroupIDInsensitively parses 'input' case-insensitively into a ManagementGroupId
// note: this method should only be used for API response data and not user input
func ParseManagementGroupIDInsensitively(input string) (*ManagementGroupId, error) {
	parser := resourceids.NewParserFromResourceIdType(ManagementGroupId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := ManagementGroupId{}

	if id.GroupId, ok = parsed.Parsed["groupId"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "groupId", *parsed)
	}

	return &id, nil
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
	return fmt.Sprintf(fmtString, id.GroupId)
}

// Segments returns a slice of Resource ID Segments which comprise this Management Group ID
func (id ManagementGroupId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("providers", "providers", "providers"),
		resourceids.ResourceProviderSegment("resourceProvider", "Microsoft.Management", "Microsoft.Management"),
		resourceids.StaticSegment("managementGroups", "managementGroups", "managementGroups"),
		resourceids.UserSpecifiedSegment("groupId", "groupIdValue"),
	}
}

// String returns a human-readable description of this Management Group ID
func (id ManagementGroupId) String() string {
	components := []string{
		fmt.Sprintf("Group: %q", id.GroupId),
	}
	return fmt.Sprintf("Management Group (%s)", strings.Join(components, "\n"))
}
