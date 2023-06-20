package databases

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ resourceids.ResourceId = RedisEnterpriseId{}

// RedisEnterpriseId is a struct representing the Resource ID for a Redis Enterprise
type RedisEnterpriseId struct {
	SubscriptionId      string
	ResourceGroupName   string
	RedisEnterpriseName string
}

// NewRedisEnterpriseID returns a new RedisEnterpriseId struct
func NewRedisEnterpriseID(subscriptionId string, resourceGroupName string, redisEnterpriseName string) RedisEnterpriseId {
	return RedisEnterpriseId{
		SubscriptionId:      subscriptionId,
		ResourceGroupName:   resourceGroupName,
		RedisEnterpriseName: redisEnterpriseName,
	}
}

// ParseRedisEnterpriseID parses 'input' into a RedisEnterpriseId
func ParseRedisEnterpriseID(input string) (*RedisEnterpriseId, error) {
	parser := resourceids.NewParserFromResourceIdType(RedisEnterpriseId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := RedisEnterpriseId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", *parsed)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", *parsed)
	}

	if id.RedisEnterpriseName, ok = parsed.Parsed["redisEnterpriseName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "redisEnterpriseName", *parsed)
	}

	return &id, nil
}

// ParseRedisEnterpriseIDInsensitively parses 'input' case-insensitively into a RedisEnterpriseId
// note: this method should only be used for API response data and not user input
func ParseRedisEnterpriseIDInsensitively(input string) (*RedisEnterpriseId, error) {
	parser := resourceids.NewParserFromResourceIdType(RedisEnterpriseId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := RedisEnterpriseId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", *parsed)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", *parsed)
	}

	if id.RedisEnterpriseName, ok = parsed.Parsed["redisEnterpriseName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "redisEnterpriseName", *parsed)
	}

	return &id, nil
}

// ValidateRedisEnterpriseID checks that 'input' can be parsed as a Redis Enterprise ID
func ValidateRedisEnterpriseID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseRedisEnterpriseID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Redis Enterprise ID
func (id RedisEnterpriseId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Cache/redisEnterprise/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.RedisEnterpriseName)
}

// Segments returns a slice of Resource ID Segments which comprise this Redis Enterprise ID
func (id RedisEnterpriseId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftCache", "Microsoft.Cache", "Microsoft.Cache"),
		resourceids.StaticSegment("staticRedisEnterprise", "redisEnterprise", "redisEnterprise"),
		resourceids.UserSpecifiedSegment("redisEnterpriseName", "redisEnterpriseValue"),
	}
}

// String returns a human-readable description of this Redis Enterprise ID
func (id RedisEnterpriseId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Redis Enterprise Name: %q", id.RedisEnterpriseName),
	}
	return fmt.Sprintf("Redis Enterprise (%s)", strings.Join(components, "\n"))
}
