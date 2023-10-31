package subscription

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ resourceids.ResourceId = UserSubscriptions2Id{}

// UserSubscriptions2Id is a struct representing the Resource ID for a User Subscriptions 2
type UserSubscriptions2Id struct {
	SubscriptionId    string
	ResourceGroupName string
	ServiceName       string
	UserId            string
	SubscriptionName  string
}

// NewUserSubscriptions2ID returns a new UserSubscriptions2Id struct
func NewUserSubscriptions2ID(subscriptionId string, resourceGroupName string, serviceName string, userId string, subscriptionName string) UserSubscriptions2Id {
	return UserSubscriptions2Id{
		SubscriptionId:    subscriptionId,
		ResourceGroupName: resourceGroupName,
		ServiceName:       serviceName,
		UserId:            userId,
		SubscriptionName:  subscriptionName,
	}
}

// ParseUserSubscriptions2ID parses 'input' into a UserSubscriptions2Id
func ParseUserSubscriptions2ID(input string) (*UserSubscriptions2Id, error) {
	parser := resourceids.NewParserFromResourceIdType(UserSubscriptions2Id{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := UserSubscriptions2Id{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", *parsed)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", *parsed)
	}

	if id.ServiceName, ok = parsed.Parsed["serviceName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "serviceName", *parsed)
	}

	if id.UserId, ok = parsed.Parsed["userId"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "userId", *parsed)
	}

	if id.SubscriptionName, ok = parsed.Parsed["subscriptionName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "subscriptionName", *parsed)
	}

	return &id, nil
}

// ParseUserSubscriptions2IDInsensitively parses 'input' case-insensitively into a UserSubscriptions2Id
// note: this method should only be used for API response data and not user input
func ParseUserSubscriptions2IDInsensitively(input string) (*UserSubscriptions2Id, error) {
	parser := resourceids.NewParserFromResourceIdType(UserSubscriptions2Id{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := UserSubscriptions2Id{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", *parsed)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", *parsed)
	}

	if id.ServiceName, ok = parsed.Parsed["serviceName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "serviceName", *parsed)
	}

	if id.UserId, ok = parsed.Parsed["userId"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "userId", *parsed)
	}

	if id.SubscriptionName, ok = parsed.Parsed["subscriptionName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "subscriptionName", *parsed)
	}

	return &id, nil
}

// ValidateUserSubscriptions2ID checks that 'input' can be parsed as a User Subscriptions 2 ID
func ValidateUserSubscriptions2ID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseUserSubscriptions2ID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted User Subscriptions 2 ID
func (id UserSubscriptions2Id) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.ApiManagement/service/%s/users/%s/subscriptions/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.ServiceName, id.UserId, id.SubscriptionName)
}

// Segments returns a slice of Resource ID Segments which comprise this User Subscriptions 2 ID
func (id UserSubscriptions2Id) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftApiManagement", "Microsoft.ApiManagement", "Microsoft.ApiManagement"),
		resourceids.StaticSegment("staticService", "service", "service"),
		resourceids.UserSpecifiedSegment("serviceName", "serviceValue"),
		resourceids.StaticSegment("staticUsers", "users", "users"),
		resourceids.UserSpecifiedSegment("userId", "userIdValue"),
		resourceids.StaticSegment("staticSubscriptions2", "subscriptions", "subscriptions"),
		resourceids.UserSpecifiedSegment("subscriptionName", "subscriptionValue"),
	}
}

// String returns a human-readable description of this User Subscriptions 2 ID
func (id UserSubscriptions2Id) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Service Name: %q", id.ServiceName),
		fmt.Sprintf("User: %q", id.UserId),
		fmt.Sprintf("Subscription Name: %q", id.SubscriptionName),
	}
	return fmt.Sprintf("User Subscriptions 2 (%s)", strings.Join(components, "\n"))
}
