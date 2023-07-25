package serverkeys

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ resourceids.ResourceId = ServerId{}

// ServerId is a struct representing the Resource ID for a Server
type ServerId struct {
	SubscriptionId    string
	ResourceGroupName string
	ServerName        string
}

// NewServerID returns a new ServerId struct
func NewServerID(subscriptionId string, resourceGroupName string, serverName string) ServerId {
	return ServerId{
		SubscriptionId:    subscriptionId,
		ResourceGroupName: resourceGroupName,
		ServerName:        serverName,
	}
}

// ParseServerID parses 'input' into a ServerId
func ParseServerID(input string) (*ServerId, error) {
	parser := resourceids.NewParserFromResourceIdType(ServerId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := ServerId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", *parsed)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", *parsed)
	}

	if id.ServerName, ok = parsed.Parsed["serverName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "serverName", *parsed)
	}

	return &id, nil
}

// ParseServerIDInsensitively parses 'input' case-insensitively into a ServerId
// note: this method should only be used for API response data and not user input
func ParseServerIDInsensitively(input string) (*ServerId, error) {
	parser := resourceids.NewParserFromResourceIdType(ServerId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := ServerId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", *parsed)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", *parsed)
	}

	if id.ServerName, ok = parsed.Parsed["serverName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "serverName", *parsed)
	}

	return &id, nil
}

// ValidateServerID checks that 'input' can be parsed as a Server ID
func ValidateServerID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseServerID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Server ID
func (id ServerId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.DBforPostgreSQL/servers/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.ServerName)
}

// Segments returns a slice of Resource ID Segments which comprise this Server ID
func (id ServerId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftDBforPostgreSQL", "Microsoft.DBforPostgreSQL", "Microsoft.DBforPostgreSQL"),
		resourceids.StaticSegment("staticServers", "servers", "servers"),
		resourceids.UserSpecifiedSegment("serverName", "serverValue"),
	}
}

// String returns a human-readable description of this Server ID
func (id ServerId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Server Name: %q", id.ServerName),
	}
	return fmt.Sprintf("Server (%s)", strings.Join(components, "\n"))
}
