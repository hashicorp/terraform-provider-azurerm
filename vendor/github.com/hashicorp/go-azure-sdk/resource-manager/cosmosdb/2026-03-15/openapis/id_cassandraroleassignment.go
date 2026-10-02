package openapis

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/recaser"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

func init() {
	recaser.RegisterResourceId(&CassandraRoleAssignmentId{})
}

var _ resourceids.ResourceId = &CassandraRoleAssignmentId{}

// CassandraRoleAssignmentId is a struct representing the Resource ID for a Cassandra Role Assignment
type CassandraRoleAssignmentId struct {
	SubscriptionId      string
	ResourceGroupName   string
	DatabaseAccountName string
	RoleAssignmentId    string
}

// NewCassandraRoleAssignmentID returns a new CassandraRoleAssignmentId struct
func NewCassandraRoleAssignmentID(subscriptionId string, resourceGroupName string, databaseAccountName string, roleAssignmentId string) CassandraRoleAssignmentId {
	return CassandraRoleAssignmentId{
		SubscriptionId:      subscriptionId,
		ResourceGroupName:   resourceGroupName,
		DatabaseAccountName: databaseAccountName,
		RoleAssignmentId:    roleAssignmentId,
	}
}

// ParseCassandraRoleAssignmentID parses 'input' into a CassandraRoleAssignmentId
func ParseCassandraRoleAssignmentID(input string) (*CassandraRoleAssignmentId, error) {
	parser := resourceids.NewParserFromResourceIdType(&CassandraRoleAssignmentId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := CassandraRoleAssignmentId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParseCassandraRoleAssignmentIDInsensitively parses 'input' case-insensitively into a CassandraRoleAssignmentId
// note: this method should only be used for API response data and not user input
func ParseCassandraRoleAssignmentIDInsensitively(input string) (*CassandraRoleAssignmentId, error) {
	parser := resourceids.NewParserFromResourceIdType(&CassandraRoleAssignmentId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := CassandraRoleAssignmentId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *CassandraRoleAssignmentId) FromParseResult(input resourceids.ParseResult) error {
	var ok bool

	if id.SubscriptionId, ok = input.Parsed["subscriptionId"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", input)
	}

	if id.ResourceGroupName, ok = input.Parsed["resourceGroupName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", input)
	}

	if id.DatabaseAccountName, ok = input.Parsed["databaseAccountName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "databaseAccountName", input)
	}

	if id.RoleAssignmentId, ok = input.Parsed["roleAssignmentId"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "roleAssignmentId", input)
	}

	return nil
}

// ValidateCassandraRoleAssignmentID checks that 'input' can be parsed as a Cassandra Role Assignment ID
func ValidateCassandraRoleAssignmentID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseCassandraRoleAssignmentID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Cassandra Role Assignment ID
func (id CassandraRoleAssignmentId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.DocumentDB/databaseAccounts/%s/cassandraRoleAssignments/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.DatabaseAccountName, id.RoleAssignmentId)
}

// Segments returns a slice of Resource ID Segments which comprise this Cassandra Role Assignment ID
func (id CassandraRoleAssignmentId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftDocumentDB", "Microsoft.DocumentDB", "Microsoft.DocumentDB"),
		resourceids.StaticSegment("staticDatabaseAccounts", "databaseAccounts", "databaseAccounts"),
		resourceids.UserSpecifiedSegment("databaseAccountName", "databaseAccountName"),
		resourceids.StaticSegment("staticCassandraRoleAssignments", "cassandraRoleAssignments", "cassandraRoleAssignments"),
		resourceids.UserSpecifiedSegment("roleAssignmentId", "roleAssignmentId"),
	}
}

// String returns a human-readable description of this Cassandra Role Assignment ID
func (id CassandraRoleAssignmentId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Database Account Name: %q", id.DatabaseAccountName),
		fmt.Sprintf("Role Assignment: %q", id.RoleAssignmentId),
	}
	return fmt.Sprintf("Cassandra Role Assignment (%s)", strings.Join(components, "\n"))
}
