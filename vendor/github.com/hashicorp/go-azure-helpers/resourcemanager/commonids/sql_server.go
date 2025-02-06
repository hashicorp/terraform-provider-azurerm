// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package commonids

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

var _ resourceids.ResourceId = &SqlServerId{}

// SqlServerId is a struct representing the Resource ID for a Sql Server
type SqlServerId struct {
	SubscriptionId    string
	ResourceGroupName string
	ServerName        string
}

// NewSqlServerID returns a new SqlServerId struct
func NewSqlServerID(subscriptionId string, resourceGroupName string, serverName string) SqlServerId {
	return SqlServerId{
		SubscriptionId:    subscriptionId,
		ResourceGroupName: resourceGroupName,
		ServerName:        serverName,
	}
}

// ParseSqlServerID parses 'input' into a SqlServerId
func ParseSqlServerID(input string) (*SqlServerId, error) {
	parser := resourceids.NewParserFromResourceIdType(&SqlServerId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := SqlServerId{}
	if err := id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParseSqlServerIDInsensitively parses 'input' case-insensitively into a SqlServerId
// note: this method should only be used for API response data and not user input
func ParseSqlServerIDInsensitively(input string) (*SqlServerId, error) {
	parser := resourceids.NewParserFromResourceIdType(&SqlServerId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := SqlServerId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *SqlServerId) FromParseResult(input resourceids.ParseResult) error {
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

	return nil
}

// ValidateSqlServerID checks that 'input' can be parsed as a Sql Server ID
func ValidateSqlServerID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseSqlServerID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Sql Server ID
func (id SqlServerId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Sql/servers/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.ServerName)
}

// Segments returns a slice of Resource ID Segments which comprise this Sql Server ID
func (id SqlServerId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftSql", "Microsoft.Sql", "Microsoft.Sql"),
		resourceids.StaticSegment("staticServers", "servers", "servers"),
		resourceids.UserSpecifiedSegment("serverName", "serverValue"),
	}
}

// String returns a human-readable description of this Sql Server ID
func (id SqlServerId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Server Name: %q", id.ServerName),
	}
	return fmt.Sprintf("Server (%s)", strings.Join(components, "\n"))
}
