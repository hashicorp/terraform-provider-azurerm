package user

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/recaser"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

func init() {
	recaser.RegisterResourceId(&UserId{})
}

var _ resourceids.ResourceId = &UserId{}

// UserId is a struct representing the Resource ID for a User
type UserId struct {
	SubscriptionId    string
	ResourceGroupName string
	ServiceName       string
	UserId            string
}

// NewUserID returns a new UserId struct
func NewUserID(subscriptionId string, resourceGroupName string, serviceName string, userId string) UserId {
	return UserId{
		SubscriptionId:    subscriptionId,
		ResourceGroupName: resourceGroupName,
		ServiceName:       serviceName,
		UserId:            userId,
	}
}

// ParseUserID parses 'input' into a UserId
func ParseUserID(input string) (*UserId, error) {
	parser := resourceids.NewParserFromResourceIdType(&UserId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := UserId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParseUserIDInsensitively parses 'input' case-insensitively into a UserId
// note: this method should only be used for API response data and not user input
func ParseUserIDInsensitively(input string) (*UserId, error) {
	parser := resourceids.NewParserFromResourceIdType(&UserId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := UserId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *UserId) FromParseResult(input resourceids.ParseResult) error {
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

	return nil
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
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.ApiManagement/service/%s/users/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.ServiceName, id.UserId)
}

// Segments returns a slice of Resource ID Segments which comprise this User ID
func (id UserId) Segments() []resourceids.Segment {
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
	}
}

// String returns a human-readable description of this User ID
func (id UserId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Service Name: %q", id.ServiceName),
		fmt.Sprintf("User: %q", id.UserId),
	}
	return fmt.Sprintf("User (%s)", strings.Join(components, "\n"))
}
