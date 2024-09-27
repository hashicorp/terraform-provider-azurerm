package privatelinkassociation

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/recaser"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

func init() {
	recaser.RegisterResourceId(&PrivateLinkAssociationId{})
}

var _ resourceids.ResourceId = &PrivateLinkAssociationId{}

// PrivateLinkAssociationId is a struct representing the Resource ID for a Private Link Association
type PrivateLinkAssociationId struct {
	GroupId string
	PlaId   string
}

// NewPrivateLinkAssociationID returns a new PrivateLinkAssociationId struct
func NewPrivateLinkAssociationID(groupId string, plaId string) PrivateLinkAssociationId {
	return PrivateLinkAssociationId{
		GroupId: groupId,
		PlaId:   plaId,
	}
}

// ParsePrivateLinkAssociationID parses 'input' into a PrivateLinkAssociationId
func ParsePrivateLinkAssociationID(input string) (*PrivateLinkAssociationId, error) {
	parser := resourceids.NewParserFromResourceIdType(&PrivateLinkAssociationId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := PrivateLinkAssociationId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParsePrivateLinkAssociationIDInsensitively parses 'input' case-insensitively into a PrivateLinkAssociationId
// note: this method should only be used for API response data and not user input
func ParsePrivateLinkAssociationIDInsensitively(input string) (*PrivateLinkAssociationId, error) {
	parser := resourceids.NewParserFromResourceIdType(&PrivateLinkAssociationId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := PrivateLinkAssociationId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *PrivateLinkAssociationId) FromParseResult(input resourceids.ParseResult) error {
	var ok bool

	if id.GroupId, ok = input.Parsed["groupId"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "groupId", input)
	}

	if id.PlaId, ok = input.Parsed["plaId"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "plaId", input)
	}

	return nil
}

// ValidatePrivateLinkAssociationID checks that 'input' can be parsed as a Private Link Association ID
func ValidatePrivateLinkAssociationID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParsePrivateLinkAssociationID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Private Link Association ID
func (id PrivateLinkAssociationId) ID() string {
	fmtString := "/providers/Microsoft.Management/managementGroups/%s/providers/Microsoft.Authorization/privateLinkAssociations/%s"
	return fmt.Sprintf(fmtString, id.GroupId, id.PlaId)
}

// Segments returns a slice of Resource ID Segments which comprise this Private Link Association ID
func (id PrivateLinkAssociationId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftManagement", "Microsoft.Management", "Microsoft.Management"),
		resourceids.StaticSegment("staticManagementGroups", "managementGroups", "managementGroups"),
		resourceids.UserSpecifiedSegment("groupId", "groupId"),
		resourceids.StaticSegment("staticProviders2", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftAuthorization", "Microsoft.Authorization", "Microsoft.Authorization"),
		resourceids.StaticSegment("staticPrivateLinkAssociations", "privateLinkAssociations", "privateLinkAssociations"),
		resourceids.UserSpecifiedSegment("plaId", "plaId"),
	}
}

// String returns a human-readable description of this Private Link Association ID
func (id PrivateLinkAssociationId) String() string {
	components := []string{
		fmt.Sprintf("Group: %q", id.GroupId),
		fmt.Sprintf("Pla: %q", id.PlaId),
	}
	return fmt.Sprintf("Private Link Association (%s)", strings.Join(components, "\n"))
}
