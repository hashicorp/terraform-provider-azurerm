package groupquotalimits

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/recaser"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

func init() {
	recaser.RegisterResourceId(&GroupQuotaLimitId{})
}

var _ resourceids.ResourceId = &GroupQuotaLimitId{}

// GroupQuotaLimitId is a struct representing the Resource ID for a Group Quota Limit
type GroupQuotaLimitId struct {
	ManagementGroupId    string
	GroupQuotaName       string
	ResourceProviderName string
	GroupQuotaLimitName  string
}

// NewGroupQuotaLimitID returns a new GroupQuotaLimitId struct
func NewGroupQuotaLimitID(managementGroupId string, groupQuotaName string, resourceProviderName string, groupQuotaLimitName string) GroupQuotaLimitId {
	return GroupQuotaLimitId{
		ManagementGroupId:    managementGroupId,
		GroupQuotaName:       groupQuotaName,
		ResourceProviderName: resourceProviderName,
		GroupQuotaLimitName:  groupQuotaLimitName,
	}
}

// ParseGroupQuotaLimitID parses 'input' into a GroupQuotaLimitId
func ParseGroupQuotaLimitID(input string) (*GroupQuotaLimitId, error) {
	parser := resourceids.NewParserFromResourceIdType(&GroupQuotaLimitId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := GroupQuotaLimitId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParseGroupQuotaLimitIDInsensitively parses 'input' case-insensitively into a GroupQuotaLimitId
// note: this method should only be used for API response data and not user input
func ParseGroupQuotaLimitIDInsensitively(input string) (*GroupQuotaLimitId, error) {
	parser := resourceids.NewParserFromResourceIdType(&GroupQuotaLimitId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := GroupQuotaLimitId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *GroupQuotaLimitId) FromParseResult(input resourceids.ParseResult) error {
	var ok bool

	if id.ManagementGroupId, ok = input.Parsed["managementGroupId"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "managementGroupId", input)
	}

	if id.GroupQuotaName, ok = input.Parsed["groupQuotaName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "groupQuotaName", input)
	}

	if id.ResourceProviderName, ok = input.Parsed["resourceProviderName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "resourceProviderName", input)
	}

	if id.GroupQuotaLimitName, ok = input.Parsed["groupQuotaLimitName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "groupQuotaLimitName", input)
	}

	return nil
}

// ValidateGroupQuotaLimitID checks that 'input' can be parsed as a Group Quota Limit ID
func ValidateGroupQuotaLimitID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseGroupQuotaLimitID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Group Quota Limit ID
func (id GroupQuotaLimitId) ID() string {
	fmtString := "/providers/Microsoft.Management/managementGroups/%s/providers/Microsoft.Quota/groupQuotas/%s/resourceProviders/%s/groupQuotaLimits/%s"
	return fmt.Sprintf(fmtString, id.ManagementGroupId, id.GroupQuotaName, id.ResourceProviderName, id.GroupQuotaLimitName)
}

// Segments returns a slice of Resource ID Segments which comprise this Group Quota Limit ID
func (id GroupQuotaLimitId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftManagement", "Microsoft.Management", "Microsoft.Management"),
		resourceids.StaticSegment("staticManagementGroups", "managementGroups", "managementGroups"),
		resourceids.UserSpecifiedSegment("managementGroupId", "managementGroupId"),
		resourceids.StaticSegment("staticProviders2", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftQuota", "Microsoft.Quota", "Microsoft.Quota"),
		resourceids.StaticSegment("staticGroupQuotas", "groupQuotas", "groupQuotas"),
		resourceids.UserSpecifiedSegment("groupQuotaName", "groupQuotaName"),
		resourceids.StaticSegment("staticResourceProviders", "resourceProviders", "resourceProviders"),
		resourceids.UserSpecifiedSegment("resourceProviderName", "resourceProviderName"),
		resourceids.StaticSegment("staticGroupQuotaLimits", "groupQuotaLimits", "groupQuotaLimits"),
		resourceids.UserSpecifiedSegment("groupQuotaLimitName", "groupQuotaLimitName"),
	}
}

// String returns a human-readable description of this Group Quota Limit ID
func (id GroupQuotaLimitId) String() string {
	components := []string{
		fmt.Sprintf("Management Group: %q", id.ManagementGroupId),
		fmt.Sprintf("Group Quota Name: %q", id.GroupQuotaName),
		fmt.Sprintf("Resource Provider Name: %q", id.ResourceProviderName),
		fmt.Sprintf("Group Quota Limit Name: %q", id.GroupQuotaLimitName),
	}
	return fmt.Sprintf("Group Quota Limit (%s)", strings.Join(components, "\n"))
}
