package firewallrules

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ resourceids.ResourceId = RediId{}

// RediId is a struct representing the Resource ID for a Redi
type RediId struct {
	SubscriptionId    string
	ResourceGroupName string
	RedisName         string
}

// NewRediID returns a new RediId struct
func NewRediID(subscriptionId string, resourceGroupName string, redisName string) RediId {
	return RediId{
		SubscriptionId:    subscriptionId,
		ResourceGroupName: resourceGroupName,
		RedisName:         redisName,
	}
}

// ParseRediID parses 'input' into a RediId
func ParseRediID(input string) (*RediId, error) {
	parser := resourceids.NewParserFromResourceIdType(RediId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := RediId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", *parsed)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", *parsed)
	}

	if id.RedisName, ok = parsed.Parsed["redisName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "redisName", *parsed)
	}

	return &id, nil
}

// ParseRediIDInsensitively parses 'input' case-insensitively into a RediId
// note: this method should only be used for API response data and not user input
func ParseRediIDInsensitively(input string) (*RediId, error) {
	parser := resourceids.NewParserFromResourceIdType(RediId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := RediId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", *parsed)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", *parsed)
	}

	if id.RedisName, ok = parsed.Parsed["redisName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "redisName", *parsed)
	}

	return &id, nil
}

// ValidateRediID checks that 'input' can be parsed as a Redi ID
func ValidateRediID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseRediID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Redi ID
func (id RediId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Cache/redis/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.RedisName)
}

// Segments returns a slice of Resource ID Segments which comprise this Redi ID
func (id RediId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftCache", "Microsoft.Cache", "Microsoft.Cache"),
		resourceids.StaticSegment("staticRedis", "redis", "redis"),
		resourceids.UserSpecifiedSegment("redisName", "redisValue"),
	}
}

// String returns a human-readable description of this Redi ID
func (id RediId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Redis Name: %q", id.RedisName),
	}
	return fmt.Sprintf("Redi (%s)", strings.Join(components, "\n"))
}
