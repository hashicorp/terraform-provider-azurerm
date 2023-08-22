package groupuser

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ resourceids.ResourceId = GroupUserId{}

// GroupUserId is a struct representing the Resource ID for a Group User
type GroupUserId struct {
	SubscriptionId    string
	ResourceGroupName string
	ServiceName       string
	GroupId           string
	UserId            string
}

// NewGroupUserID returns a new GroupUserId struct
func NewGroupUserID(subscriptionId string, resourceGroupName string, serviceName string, groupId string, userId string) GroupUserId {
	return GroupUserId{
		SubscriptionId:    subscriptionId,
		ResourceGroupName: resourceGroupName,
		ServiceName:       serviceName,
		GroupId:           groupId,
		UserId:            userId,
	}
}

// ParseGroupUserID parses 'input' into a GroupUserId
func ParseGroupUserID(input string) (*GroupUserId, error) {
	parser := resourceids.NewParserFromResourceIdType(GroupUserId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := GroupUserId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", *parsed)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", *parsed)
	}

	if id.ServiceName, ok = parsed.Parsed["serviceName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "serviceName", *parsed)
	}

	if id.GroupId, ok = parsed.Parsed["groupId"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "groupId", *parsed)
	}

	if id.UserId, ok = parsed.Parsed["userId"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "userId", *parsed)
	}

	return &id, nil
}

// ParseGroupUserIDInsensitively parses 'input' case-insensitively into a GroupUserId
// note: this method should only be used for API response data and not user input
func ParseGroupUserIDInsensitively(input string) (*GroupUserId, error) {
	parser := resourceids.NewParserFromResourceIdType(GroupUserId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := GroupUserId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", *parsed)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", *parsed)
	}

	if id.ServiceName, ok = parsed.Parsed["serviceName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "serviceName", *parsed)
	}

	if id.GroupId, ok = parsed.Parsed["groupId"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "groupId", *parsed)
	}

	if id.UserId, ok = parsed.Parsed["userId"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "userId", *parsed)
	}

	return &id, nil
}

// ValidateGroupUserID checks that 'input' can be parsed as a Group User ID
func ValidateGroupUserID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseGroupUserID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Group User ID
func (id GroupUserId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.ApiManagement/service/%s/groups/%s/users/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.ServiceName, id.GroupId, id.UserId)
}

// Segments returns a slice of Resource ID Segments which comprise this Group User ID
func (id GroupUserId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftApiManagement", "Microsoft.ApiManagement", "Microsoft.ApiManagement"),
		resourceids.StaticSegment("staticService", "service", "service"),
		resourceids.UserSpecifiedSegment("serviceName", "serviceValue"),
		resourceids.StaticSegment("staticGroups", "groups", "groups"),
		resourceids.UserSpecifiedSegment("groupId", "groupIdValue"),
		resourceids.StaticSegment("staticUsers", "users", "users"),
		resourceids.UserSpecifiedSegment("userId", "userIdValue"),
	}
}

// String returns a human-readable description of this Group User ID
func (id GroupUserId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Service Name: %q", id.ServiceName),
		fmt.Sprintf("Group: %q", id.GroupId),
		fmt.Sprintf("User: %q", id.UserId),
	}
	return fmt.Sprintf("Group User (%s)", strings.Join(components, "\n"))
}
