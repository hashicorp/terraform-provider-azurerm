package groupquotassubscriptions

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/recaser"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

func init() {
	recaser.RegisterResourceId(&GroupQuotaId{})
}

var _ resourceids.ResourceId = &GroupQuotaId{}

// GroupQuotaId is a struct representing the Resource ID for a Group Quota
type GroupQuotaId struct {
	ManagementGroupId string
	GroupQuotaName    string
}

// NewGroupQuotaID returns a new GroupQuotaId struct
func NewGroupQuotaID(managementGroupId string, groupQuotaName string) GroupQuotaId {
	return GroupQuotaId{
		ManagementGroupId: managementGroupId,
		GroupQuotaName:    groupQuotaName,
	}
}

// ParseGroupQuotaID parses 'input' into a GroupQuotaId
func ParseGroupQuotaID(input string) (*GroupQuotaId, error) {
	parser := resourceids.NewParserFromResourceIdType(&GroupQuotaId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := GroupQuotaId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParseGroupQuotaIDInsensitively parses 'input' case-insensitively into a GroupQuotaId
// note: this method should only be used for API response data and not user input
func ParseGroupQuotaIDInsensitively(input string) (*GroupQuotaId, error) {
	parser := resourceids.NewParserFromResourceIdType(&GroupQuotaId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := GroupQuotaId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *GroupQuotaId) FromParseResult(input resourceids.ParseResult) error {
	var ok bool

	if id.ManagementGroupId, ok = input.Parsed["managementGroupId"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "managementGroupId", input)
	}

	if id.GroupQuotaName, ok = input.Parsed["groupQuotaName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "groupQuotaName", input)
	}

	return nil
}

// ValidateGroupQuotaID checks that 'input' can be parsed as a Group Quota ID
func ValidateGroupQuotaID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseGroupQuotaID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Group Quota ID
func (id GroupQuotaId) ID() string {
	fmtString := "/providers/Microsoft.Management/managementGroups/%s/providers/Microsoft.Quota/groupQuotas/%s"
	return fmt.Sprintf(fmtString, id.ManagementGroupId, id.GroupQuotaName)
}

// Segments returns a slice of Resource ID Segments which comprise this Group Quota ID
func (id GroupQuotaId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftManagement", "Microsoft.Management", "Microsoft.Management"),
		resourceids.StaticSegment("staticManagementGroups", "managementGroups", "managementGroups"),
		resourceids.UserSpecifiedSegment("managementGroupId", "managementGroupId"),
		resourceids.StaticSegment("staticProviders2", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftQuota", "Microsoft.Quota", "Microsoft.Quota"),
		resourceids.StaticSegment("staticGroupQuotas", "groupQuotas", "groupQuotas"),
		resourceids.UserSpecifiedSegment("groupQuotaName", "groupQuotaName"),
	}
}

// String returns a human-readable description of this Group Quota ID
func (id GroupQuotaId) String() string {
	components := []string{
		fmt.Sprintf("Management Group: %q", id.ManagementGroupId),
		fmt.Sprintf("Group Quota Name: %q", id.GroupQuotaName),
	}
	return fmt.Sprintf("Group Quota (%s)", strings.Join(components, "\n"))
}
