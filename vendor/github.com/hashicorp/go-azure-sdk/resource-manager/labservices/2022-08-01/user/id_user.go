package user

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ resourceids.ResourceId = UserId{}

// UserId is a struct representing the Resource ID for a User
type UserId struct {
	SubscriptionId    string
	ResourceGroupName string
	LabName           string
	UserName          string
}

// NewUserID returns a new UserId struct
func NewUserID(subscriptionId string, resourceGroupName string, labName string, userName string) UserId {
	return UserId{
		SubscriptionId:    subscriptionId,
		ResourceGroupName: resourceGroupName,
		LabName:           labName,
		UserName:          userName,
	}
}

// ParseUserID parses 'input' into a UserId
func ParseUserID(input string) (*UserId, error) {
	parser := resourceids.NewParserFromResourceIdType(UserId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := UserId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", *parsed)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", *parsed)
	}

	if id.LabName, ok = parsed.Parsed["labName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "labName", *parsed)
	}

	if id.UserName, ok = parsed.Parsed["userName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "userName", *parsed)
	}

	return &id, nil
}

// ParseUserIDInsensitively parses 'input' case-insensitively into a UserId
// note: this method should only be used for API response data and not user input
func ParseUserIDInsensitively(input string) (*UserId, error) {
	parser := resourceids.NewParserFromResourceIdType(UserId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := UserId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", *parsed)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", *parsed)
	}

	if id.LabName, ok = parsed.Parsed["labName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "labName", *parsed)
	}

	if id.UserName, ok = parsed.Parsed["userName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "userName", *parsed)
	}

	return &id, nil
}

// ValidateUserID checks that 'input' can be parsed as a User ID
func ValidateUserID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseUserID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted User ID
func (id UserId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.LabServices/labs/%s/users/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.LabName, id.UserName)
}

// Segments returns a slice of Resource ID Segments which comprise this User ID
func (id UserId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftLabServices", "Microsoft.LabServices", "Microsoft.LabServices"),
		resourceids.StaticSegment("staticLabs", "labs", "labs"),
		resourceids.UserSpecifiedSegment("labName", "labValue"),
		resourceids.StaticSegment("staticUsers", "users", "users"),
		resourceids.UserSpecifiedSegment("userName", "userValue"),
	}
}

// String returns a human-readable description of this User ID
func (id UserId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Lab Name: %q", id.LabName),
		fmt.Sprintf("User Name: %q", id.UserName),
	}
	return fmt.Sprintf("User (%s)", strings.Join(components, "\n"))
}
