package redis

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ resourceids.ResourceId = AccessPolicyAssignmentId{}

// AccessPolicyAssignmentId is a struct representing the Resource ID for a Access Policy Assignment
type AccessPolicyAssignmentId struct {
	SubscriptionId             string
	ResourceGroupName          string
	RedisName                  string
	AccessPolicyAssignmentName string
}

// NewAccessPolicyAssignmentID returns a new AccessPolicyAssignmentId struct
func NewAccessPolicyAssignmentID(subscriptionId string, resourceGroupName string, redisName string, accessPolicyAssignmentName string) AccessPolicyAssignmentId {
	return AccessPolicyAssignmentId{
		SubscriptionId:             subscriptionId,
		ResourceGroupName:          resourceGroupName,
		RedisName:                  redisName,
		AccessPolicyAssignmentName: accessPolicyAssignmentName,
	}
}

// ParseAccessPolicyAssignmentID parses 'input' into a AccessPolicyAssignmentId
func ParseAccessPolicyAssignmentID(input string) (*AccessPolicyAssignmentId, error) {
	parser := resourceids.NewParserFromResourceIdType(AccessPolicyAssignmentId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := AccessPolicyAssignmentId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", *parsed)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", *parsed)
	}

	if id.RedisName, ok = parsed.Parsed["redisName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "redisName", *parsed)
	}

	if id.AccessPolicyAssignmentName, ok = parsed.Parsed["accessPolicyAssignmentName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "accessPolicyAssignmentName", *parsed)
	}

	return &id, nil
}

// ParseAccessPolicyAssignmentIDInsensitively parses 'input' case-insensitively into a AccessPolicyAssignmentId
// note: this method should only be used for API response data and not user input
func ParseAccessPolicyAssignmentIDInsensitively(input string) (*AccessPolicyAssignmentId, error) {
	parser := resourceids.NewParserFromResourceIdType(AccessPolicyAssignmentId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := AccessPolicyAssignmentId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", *parsed)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", *parsed)
	}

	if id.RedisName, ok = parsed.Parsed["redisName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "redisName", *parsed)
	}

	if id.AccessPolicyAssignmentName, ok = parsed.Parsed["accessPolicyAssignmentName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "accessPolicyAssignmentName", *parsed)
	}

	return &id, nil
}

// ValidateAccessPolicyAssignmentID checks that 'input' can be parsed as a Access Policy Assignment ID
func ValidateAccessPolicyAssignmentID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseAccessPolicyAssignmentID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Access Policy Assignment ID
func (id AccessPolicyAssignmentId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Cache/redis/%s/accessPolicyAssignments/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.RedisName, id.AccessPolicyAssignmentName)
}

// Segments returns a slice of Resource ID Segments which comprise this Access Policy Assignment ID
func (id AccessPolicyAssignmentId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftCache", "Microsoft.Cache", "Microsoft.Cache"),
		resourceids.StaticSegment("staticRedis", "redis", "redis"),
		resourceids.UserSpecifiedSegment("redisName", "redisValue"),
		resourceids.StaticSegment("staticAccessPolicyAssignments", "accessPolicyAssignments", "accessPolicyAssignments"),
		resourceids.UserSpecifiedSegment("accessPolicyAssignmentName", "accessPolicyAssignmentValue"),
	}
}

// String returns a human-readable description of this Access Policy Assignment ID
func (id AccessPolicyAssignmentId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Redis Name: %q", id.RedisName),
		fmt.Sprintf("Access Policy Assignment Name: %q", id.AccessPolicyAssignmentName),
	}
	return fmt.Sprintf("Access Policy Assignment (%s)", strings.Join(components, "\n"))
}
