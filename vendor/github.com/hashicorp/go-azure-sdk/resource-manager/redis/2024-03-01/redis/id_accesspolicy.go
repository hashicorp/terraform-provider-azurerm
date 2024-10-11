package redis

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/recaser"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

func init() {
	recaser.RegisterResourceId(&AccessPolicyId{})
}

var _ resourceids.ResourceId = &AccessPolicyId{}

// AccessPolicyId is a struct representing the Resource ID for a Access Policy
type AccessPolicyId struct {
	SubscriptionId    string
	ResourceGroupName string
	RedisName         string
	AccessPolicyName  string
}

// NewAccessPolicyID returns a new AccessPolicyId struct
func NewAccessPolicyID(subscriptionId string, resourceGroupName string, redisName string, accessPolicyName string) AccessPolicyId {
	return AccessPolicyId{
		SubscriptionId:    subscriptionId,
		ResourceGroupName: resourceGroupName,
		RedisName:         redisName,
		AccessPolicyName:  accessPolicyName,
	}
}

// ParseAccessPolicyID parses 'input' into a AccessPolicyId
func ParseAccessPolicyID(input string) (*AccessPolicyId, error) {
	parser := resourceids.NewParserFromResourceIdType(&AccessPolicyId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := AccessPolicyId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParseAccessPolicyIDInsensitively parses 'input' case-insensitively into a AccessPolicyId
// note: this method should only be used for API response data and not user input
func ParseAccessPolicyIDInsensitively(input string) (*AccessPolicyId, error) {
	parser := resourceids.NewParserFromResourceIdType(&AccessPolicyId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := AccessPolicyId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *AccessPolicyId) FromParseResult(input resourceids.ParseResult) error {
	var ok bool

	if id.SubscriptionId, ok = input.Parsed["subscriptionId"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", input)
	}

	if id.ResourceGroupName, ok = input.Parsed["resourceGroupName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", input)
	}

	if id.RedisName, ok = input.Parsed["redisName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "redisName", input)
	}

	if id.AccessPolicyName, ok = input.Parsed["accessPolicyName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "accessPolicyName", input)
	}

	return nil
}

// ValidateAccessPolicyID checks that 'input' can be parsed as a Access Policy ID
func ValidateAccessPolicyID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseAccessPolicyID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Access Policy ID
func (id AccessPolicyId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Cache/redis/%s/accessPolicies/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.RedisName, id.AccessPolicyName)
}

// Segments returns a slice of Resource ID Segments which comprise this Access Policy ID
func (id AccessPolicyId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftCache", "Microsoft.Cache", "Microsoft.Cache"),
		resourceids.StaticSegment("staticRedis", "redis", "redis"),
		resourceids.UserSpecifiedSegment("redisName", "redisName"),
		resourceids.StaticSegment("staticAccessPolicies", "accessPolicies", "accessPolicies"),
		resourceids.UserSpecifiedSegment("accessPolicyName", "accessPolicyName"),
	}
}

// String returns a human-readable description of this Access Policy ID
func (id AccessPolicyId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Redis Name: %q", id.RedisName),
		fmt.Sprintf("Access Policy Name: %q", id.AccessPolicyName),
	}
	return fmt.Sprintf("Access Policy (%s)", strings.Join(components, "\n"))
}
