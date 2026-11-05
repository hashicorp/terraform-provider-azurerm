package failovergroups

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/recaser"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

func init() {
	recaser.RegisterResourceId(&FailoverGroupId{})
}

var _ resourceids.ResourceId = &FailoverGroupId{}

// FailoverGroupId is a struct representing the Resource ID for a Failover Group
type FailoverGroupId struct {
	SubscriptionId    string
	ResourceGroupName string
	ServerName        string
	FailoverGroupName string
}

// NewFailoverGroupID returns a new FailoverGroupId struct
func NewFailoverGroupID(subscriptionId string, resourceGroupName string, serverName string, failoverGroupName string) FailoverGroupId {
	return FailoverGroupId{
		SubscriptionId:    subscriptionId,
		ResourceGroupName: resourceGroupName,
		ServerName:        serverName,
		FailoverGroupName: failoverGroupName,
	}
}

// ParseFailoverGroupID parses 'input' into a FailoverGroupId
func ParseFailoverGroupID(input string) (*FailoverGroupId, error) {
	parser := resourceids.NewParserFromResourceIdType(&FailoverGroupId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := FailoverGroupId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParseFailoverGroupIDInsensitively parses 'input' case-insensitively into a FailoverGroupId
// note: this method should only be used for API response data and not user input
func ParseFailoverGroupIDInsensitively(input string) (*FailoverGroupId, error) {
	parser := resourceids.NewParserFromResourceIdType(&FailoverGroupId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := FailoverGroupId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *FailoverGroupId) FromParseResult(input resourceids.ParseResult) error {
	var ok bool

	if id.SubscriptionId, ok = input.Parsed["subscriptionId"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", input)
	}

	if id.ResourceGroupName, ok = input.Parsed["resourceGroupName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", input)
	}

	if id.ServerName, ok = input.Parsed["serverName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "serverName", input)
	}

	if id.FailoverGroupName, ok = input.Parsed["failoverGroupName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "failoverGroupName", input)
	}

	return nil
}

// ValidateFailoverGroupID checks that 'input' can be parsed as a Failover Group ID
func ValidateFailoverGroupID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseFailoverGroupID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Failover Group ID
func (id FailoverGroupId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Sql/servers/%s/failoverGroups/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.ServerName, id.FailoverGroupName)
}

// Segments returns a slice of Resource ID Segments which comprise this Failover Group ID
func (id FailoverGroupId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftSql", "Microsoft.Sql", "Microsoft.Sql"),
		resourceids.StaticSegment("staticServers", "servers", "servers"),
		resourceids.UserSpecifiedSegment("serverName", "serverName"),
		resourceids.StaticSegment("staticFailoverGroups", "failoverGroups", "failoverGroups"),
		resourceids.UserSpecifiedSegment("failoverGroupName", "failoverGroupName"),
	}
}

// String returns a human-readable description of this Failover Group ID
func (id FailoverGroupId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Server Name: %q", id.ServerName),
		fmt.Sprintf("Failover Group Name: %q", id.FailoverGroupName),
	}
	return fmt.Sprintf("Failover Group (%s)", strings.Join(components, "\n"))
}
