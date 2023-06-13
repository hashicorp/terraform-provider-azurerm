package roles

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ resourceids.ResourceId = RoleId{}

// RoleId is a struct representing the Resource ID for a Role
type RoleId struct {
	SubscriptionId     string
	ResourceGroupName  string
	ServerGroupsv2Name string
	RoleName           string
}

// NewRoleID returns a new RoleId struct
func NewRoleID(subscriptionId string, resourceGroupName string, serverGroupsv2Name string, roleName string) RoleId {
	return RoleId{
		SubscriptionId:     subscriptionId,
		ResourceGroupName:  resourceGroupName,
		ServerGroupsv2Name: serverGroupsv2Name,
		RoleName:           roleName,
	}
}

// ParseRoleID parses 'input' into a RoleId
func ParseRoleID(input string) (*RoleId, error) {
	parser := resourceids.NewParserFromResourceIdType(RoleId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := RoleId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", *parsed)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", *parsed)
	}

	if id.ServerGroupsv2Name, ok = parsed.Parsed["serverGroupsv2Name"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "serverGroupsv2Name", *parsed)
	}

	if id.RoleName, ok = parsed.Parsed["roleName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "roleName", *parsed)
	}

	return &id, nil
}

// ParseRoleIDInsensitively parses 'input' case-insensitively into a RoleId
// note: this method should only be used for API response data and not user input
func ParseRoleIDInsensitively(input string) (*RoleId, error) {
	parser := resourceids.NewParserFromResourceIdType(RoleId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := RoleId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", *parsed)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", *parsed)
	}

	if id.ServerGroupsv2Name, ok = parsed.Parsed["serverGroupsv2Name"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "serverGroupsv2Name", *parsed)
	}

	if id.RoleName, ok = parsed.Parsed["roleName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "roleName", *parsed)
	}

	return &id, nil
}

// ValidateRoleID checks that 'input' can be parsed as a Role ID
func ValidateRoleID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseRoleID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Role ID
func (id RoleId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.DBforPostgreSQL/serverGroupsv2/%s/roles/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.ServerGroupsv2Name, id.RoleName)
}

// Segments returns a slice of Resource ID Segments which comprise this Role ID
func (id RoleId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftDBforPostgreSQL", "Microsoft.DBforPostgreSQL", "Microsoft.DBforPostgreSQL"),
		resourceids.StaticSegment("staticServerGroupsv2", "serverGroupsv2", "serverGroupsv2"),
		resourceids.UserSpecifiedSegment("serverGroupsv2Name", "serverGroupsv2Value"),
		resourceids.StaticSegment("staticRoles", "roles", "roles"),
		resourceids.UserSpecifiedSegment("roleName", "roleValue"),
	}
}

// String returns a human-readable description of this Role ID
func (id RoleId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Server Groupsv 2 Name: %q", id.ServerGroupsv2Name),
		fmt.Sprintf("Role Name: %q", id.RoleName),
	}
	return fmt.Sprintf("Role (%s)", strings.Join(components, "\n"))
}
