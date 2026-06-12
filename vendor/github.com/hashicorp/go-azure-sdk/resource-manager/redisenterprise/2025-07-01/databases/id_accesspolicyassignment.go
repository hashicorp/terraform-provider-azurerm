package databases

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/recaser"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

func init() {
	recaser.RegisterResourceId(&AccessPolicyAssignmentId{})
}

var _ resourceids.ResourceId = &AccessPolicyAssignmentId{}

// AccessPolicyAssignmentId is a struct representing the Resource ID for a Access Policy Assignment
type AccessPolicyAssignmentId struct {
	SubscriptionId             string
	ResourceGroupName          string
	RedisEnterpriseName        string
	DatabaseName               string
	AccessPolicyAssignmentName string
}

// NewAccessPolicyAssignmentID returns a new AccessPolicyAssignmentId struct
func NewAccessPolicyAssignmentID(subscriptionId string, resourceGroupName string, redisEnterpriseName string, databaseName string, accessPolicyAssignmentName string) AccessPolicyAssignmentId {
	return AccessPolicyAssignmentId{
		SubscriptionId:             subscriptionId,
		ResourceGroupName:          resourceGroupName,
		RedisEnterpriseName:        redisEnterpriseName,
		DatabaseName:               databaseName,
		AccessPolicyAssignmentName: accessPolicyAssignmentName,
	}
}

// ParseAccessPolicyAssignmentID parses 'input' into a AccessPolicyAssignmentId
func ParseAccessPolicyAssignmentID(input string) (*AccessPolicyAssignmentId, error) {
	parser := resourceids.NewParserFromResourceIdType(&AccessPolicyAssignmentId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := AccessPolicyAssignmentId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParseAccessPolicyAssignmentIDInsensitively parses 'input' case-insensitively into a AccessPolicyAssignmentId
// note: this method should only be used for API response data and not user input
func ParseAccessPolicyAssignmentIDInsensitively(input string) (*AccessPolicyAssignmentId, error) {
	parser := resourceids.NewParserFromResourceIdType(&AccessPolicyAssignmentId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := AccessPolicyAssignmentId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *AccessPolicyAssignmentId) FromParseResult(input resourceids.ParseResult) error {
	var ok bool

	if id.SubscriptionId, ok = input.Parsed["subscriptionId"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", input)
	}

	if id.ResourceGroupName, ok = input.Parsed["resourceGroupName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", input)
	}

	if id.RedisEnterpriseName, ok = input.Parsed["redisEnterpriseName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "redisEnterpriseName", input)
	}

	if id.DatabaseName, ok = input.Parsed["databaseName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "databaseName", input)
	}

	if id.AccessPolicyAssignmentName, ok = input.Parsed["accessPolicyAssignmentName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "accessPolicyAssignmentName", input)
	}

	return nil
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
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Cache/redisEnterprise/%s/databases/%s/accessPolicyAssignments/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.RedisEnterpriseName, id.DatabaseName, id.AccessPolicyAssignmentName)
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
		resourceids.StaticSegment("staticRedisEnterprise", "redisEnterprise", "redisEnterprise"),
		resourceids.UserSpecifiedSegment("redisEnterpriseName", "redisEnterpriseName"),
		resourceids.StaticSegment("staticDatabases", "databases", "databases"),
		resourceids.UserSpecifiedSegment("databaseName", "databaseName"),
		resourceids.StaticSegment("staticAccessPolicyAssignments", "accessPolicyAssignments", "accessPolicyAssignments"),
		resourceids.UserSpecifiedSegment("accessPolicyAssignmentName", "accessPolicyAssignmentName"),
	}
}

// String returns a human-readable description of this Access Policy Assignment ID
func (id AccessPolicyAssignmentId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Redis Enterprise Name: %q", id.RedisEnterpriseName),
		fmt.Sprintf("Database Name: %q", id.DatabaseName),
		fmt.Sprintf("Access Policy Assignment Name: %q", id.AccessPolicyAssignmentName),
	}
	return fmt.Sprintf("Access Policy Assignment (%s)", strings.Join(components, "\n"))
}
