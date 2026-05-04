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
	recaser.RegisterResourceId(&LinkedServerId{})
}

var _ resourceids.ResourceId = &LinkedServerId{}

// LinkedServerId is a struct representing the Resource ID for a Linked Server
type LinkedServerId struct {
	SubscriptionId    string
	ResourceGroupName string
	RedisName         string
	LinkedServerName  string
}

// NewLinkedServerID returns a new LinkedServerId struct
func NewLinkedServerID(subscriptionId string, resourceGroupName string, redisName string, linkedServerName string) LinkedServerId {
	return LinkedServerId{
		SubscriptionId:    subscriptionId,
		ResourceGroupName: resourceGroupName,
		RedisName:         redisName,
		LinkedServerName:  linkedServerName,
	}
}

// ParseLinkedServerID parses 'input' into a LinkedServerId
func ParseLinkedServerID(input string) (*LinkedServerId, error) {
	parser := resourceids.NewParserFromResourceIdType(&LinkedServerId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := LinkedServerId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParseLinkedServerIDInsensitively parses 'input' case-insensitively into a LinkedServerId
// note: this method should only be used for API response data and not user input
func ParseLinkedServerIDInsensitively(input string) (*LinkedServerId, error) {
	parser := resourceids.NewParserFromResourceIdType(&LinkedServerId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := LinkedServerId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *LinkedServerId) FromParseResult(input resourceids.ParseResult) error {
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

	if id.LinkedServerName, ok = input.Parsed["linkedServerName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "linkedServerName", input)
	}

	return nil
}

// ValidateLinkedServerID checks that 'input' can be parsed as a Linked Server ID
func ValidateLinkedServerID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseLinkedServerID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Linked Server ID
func (id LinkedServerId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Cache/redis/%s/linkedServers/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.RedisName, id.LinkedServerName)
}

// Segments returns a slice of Resource ID Segments which comprise this Linked Server ID
func (id LinkedServerId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftCache", "Microsoft.Cache", "Microsoft.Cache"),
		resourceids.StaticSegment("staticRedis", "redis", "redis"),
		resourceids.UserSpecifiedSegment("redisName", "redisName"),
		resourceids.StaticSegment("staticLinkedServers", "linkedServers", "linkedServers"),
		resourceids.UserSpecifiedSegment("linkedServerName", "linkedServerName"),
	}
}

// String returns a human-readable description of this Linked Server ID
func (id LinkedServerId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Redis Name: %q", id.RedisName),
		fmt.Sprintf("Linked Server Name: %q", id.LinkedServerName),
	}
	return fmt.Sprintf("Linked Server (%s)", strings.Join(components, "\n"))
}
