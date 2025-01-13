package subscription

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/recaser"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

func init() {
	recaser.RegisterResourceId(&UserSubscriptions2Id{})
}

var _ resourceids.ResourceId = &UserSubscriptions2Id{}

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
	parser := resourceids.NewParserFromResourceIdType(&UserSubscriptions2Id{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := UserSubscriptions2Id{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParseUserSubscriptions2IDInsensitively parses 'input' case-insensitively into a UserSubscriptions2Id
// note: this method should only be used for API response data and not user input
func ParseUserSubscriptions2IDInsensitively(input string) (*UserSubscriptions2Id, error) {
	parser := resourceids.NewParserFromResourceIdType(&UserSubscriptions2Id{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := UserSubscriptions2Id{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *UserSubscriptions2Id) FromParseResult(input resourceids.ParseResult) error {
	var ok bool

	if id.SubscriptionId, ok = input.Parsed["subscriptionId"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", input)
	}

	if id.ResourceGroupName, ok = input.Parsed["resourceGroupName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", input)
	}

	if id.ServiceName, ok = input.Parsed["serviceName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "serviceName", input)
	}

	if id.UserId, ok = input.Parsed["userId"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "userId", input)
	}

	if id.SubscriptionName, ok = input.Parsed["subscriptionName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "subscriptionName", input)
	}

	return nil
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
		resourceids.UserSpecifiedSegment("serviceName", "serviceName"),
		resourceids.StaticSegment("staticUsers", "users", "users"),
		resourceids.UserSpecifiedSegment("userId", "userId"),
		resourceids.StaticSegment("staticSubscriptions2", "subscriptions", "subscriptions"),
		resourceids.UserSpecifiedSegment("subscriptionName", "subscriptionName"),
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
