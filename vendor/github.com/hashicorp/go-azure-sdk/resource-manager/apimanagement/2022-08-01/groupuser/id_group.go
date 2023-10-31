package groupuser

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ resourceids.ResourceId = GroupId{}

// GroupId is a struct representing the Resource ID for a Group
type GroupId struct {
	SubscriptionId    string
	ResourceGroupName string
	ServiceName       string
	GroupId           string
}

// NewGroupID returns a new GroupId struct
func NewGroupID(subscriptionId string, resourceGroupName string, serviceName string, groupId string) GroupId {
	return GroupId{
		SubscriptionId:    subscriptionId,
		ResourceGroupName: resourceGroupName,
		ServiceName:       serviceName,
		GroupId:           groupId,
	}
}

// ParseGroupID parses 'input' into a GroupId
func ParseGroupID(input string) (*GroupId, error) {
	parser := resourceids.NewParserFromResourceIdType(GroupId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := GroupId{}

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

	return &id, nil
}

// ParseGroupIDInsensitively parses 'input' case-insensitively into a GroupId
// note: this method should only be used for API response data and not user input
func ParseGroupIDInsensitively(input string) (*GroupId, error) {
	parser := resourceids.NewParserFromResourceIdType(GroupId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := GroupId{}

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

	return &id, nil
}

// ValidateGroupID checks that 'input' can be parsed as a Group ID
func ValidateGroupID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseGroupID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Group ID
func (id GroupId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.ApiManagement/service/%s/groups/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.ServiceName, id.GroupId)
}

// Segments returns a slice of Resource ID Segments which comprise this Group ID
func (id GroupId) Segments() []resourceids.Segment {
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
	}
}

// String returns a human-readable description of this Group ID
func (id GroupId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Service Name: %q", id.ServiceName),
		fmt.Sprintf("Group: %q", id.GroupId),
	}
	return fmt.Sprintf("Group (%s)", strings.Join(components, "\n"))
}
