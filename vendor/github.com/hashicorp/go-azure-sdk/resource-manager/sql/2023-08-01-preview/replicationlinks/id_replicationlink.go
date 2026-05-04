package replicationlinks

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/recaser"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

func init() {
	recaser.RegisterResourceId(&ReplicationLinkId{})
}

var _ resourceids.ResourceId = &ReplicationLinkId{}

// ReplicationLinkId is a struct representing the Resource ID for a Replication Link
type ReplicationLinkId struct {
	SubscriptionId    string
	ResourceGroupName string
	ServerName        string
	DatabaseName      string
	LinkId            string
}

// NewReplicationLinkID returns a new ReplicationLinkId struct
func NewReplicationLinkID(subscriptionId string, resourceGroupName string, serverName string, databaseName string, linkId string) ReplicationLinkId {
	return ReplicationLinkId{
		SubscriptionId:    subscriptionId,
		ResourceGroupName: resourceGroupName,
		ServerName:        serverName,
		DatabaseName:      databaseName,
		LinkId:            linkId,
	}
}

// ParseReplicationLinkID parses 'input' into a ReplicationLinkId
func ParseReplicationLinkID(input string) (*ReplicationLinkId, error) {
	parser := resourceids.NewParserFromResourceIdType(&ReplicationLinkId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := ReplicationLinkId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParseReplicationLinkIDInsensitively parses 'input' case-insensitively into a ReplicationLinkId
// note: this method should only be used for API response data and not user input
func ParseReplicationLinkIDInsensitively(input string) (*ReplicationLinkId, error) {
	parser := resourceids.NewParserFromResourceIdType(&ReplicationLinkId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := ReplicationLinkId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *ReplicationLinkId) FromParseResult(input resourceids.ParseResult) error {
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

	if id.DatabaseName, ok = input.Parsed["databaseName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "databaseName", input)
	}

	if id.LinkId, ok = input.Parsed["linkId"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "linkId", input)
	}

	return nil
}

// ValidateReplicationLinkID checks that 'input' can be parsed as a Replication Link ID
func ValidateReplicationLinkID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseReplicationLinkID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Replication Link ID
func (id ReplicationLinkId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Sql/servers/%s/databases/%s/replicationLinks/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.ServerName, id.DatabaseName, id.LinkId)
}

// Segments returns a slice of Resource ID Segments which comprise this Replication Link ID
func (id ReplicationLinkId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftSql", "Microsoft.Sql", "Microsoft.Sql"),
		resourceids.StaticSegment("staticServers", "servers", "servers"),
		resourceids.UserSpecifiedSegment("serverName", "serverName"),
		resourceids.StaticSegment("staticDatabases", "databases", "databases"),
		resourceids.UserSpecifiedSegment("databaseName", "databaseName"),
		resourceids.StaticSegment("staticReplicationLinks", "replicationLinks", "replicationLinks"),
		resourceids.UserSpecifiedSegment("linkId", "linkId"),
	}
}

// String returns a human-readable description of this Replication Link ID
func (id ReplicationLinkId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Server Name: %q", id.ServerName),
		fmt.Sprintf("Database Name: %q", id.DatabaseName),
		fmt.Sprintf("Link: %q", id.LinkId),
	}
	return fmt.Sprintf("Replication Link (%s)", strings.Join(components, "\n"))
}
