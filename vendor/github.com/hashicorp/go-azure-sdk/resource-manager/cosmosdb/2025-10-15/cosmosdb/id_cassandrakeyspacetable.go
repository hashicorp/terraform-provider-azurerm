package cosmosdb

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/recaser"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

func init() {
	recaser.RegisterResourceId(&CassandraKeyspaceTableId{})
}

var _ resourceids.ResourceId = &CassandraKeyspaceTableId{}

// CassandraKeyspaceTableId is a struct representing the Resource ID for a Cassandra Keyspace Table
type CassandraKeyspaceTableId struct {
	SubscriptionId        string
	ResourceGroupName     string
	DatabaseAccountName   string
	CassandraKeyspaceName string
	TableName             string
}

// NewCassandraKeyspaceTableID returns a new CassandraKeyspaceTableId struct
func NewCassandraKeyspaceTableID(subscriptionId string, resourceGroupName string, databaseAccountName string, cassandraKeyspaceName string, tableName string) CassandraKeyspaceTableId {
	return CassandraKeyspaceTableId{
		SubscriptionId:        subscriptionId,
		ResourceGroupName:     resourceGroupName,
		DatabaseAccountName:   databaseAccountName,
		CassandraKeyspaceName: cassandraKeyspaceName,
		TableName:             tableName,
	}
}

// ParseCassandraKeyspaceTableID parses 'input' into a CassandraKeyspaceTableId
func ParseCassandraKeyspaceTableID(input string) (*CassandraKeyspaceTableId, error) {
	parser := resourceids.NewParserFromResourceIdType(&CassandraKeyspaceTableId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := CassandraKeyspaceTableId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParseCassandraKeyspaceTableIDInsensitively parses 'input' case-insensitively into a CassandraKeyspaceTableId
// note: this method should only be used for API response data and not user input
func ParseCassandraKeyspaceTableIDInsensitively(input string) (*CassandraKeyspaceTableId, error) {
	parser := resourceids.NewParserFromResourceIdType(&CassandraKeyspaceTableId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := CassandraKeyspaceTableId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *CassandraKeyspaceTableId) FromParseResult(input resourceids.ParseResult) error {
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

	if id.CassandraKeyspaceName, ok = input.Parsed["cassandraKeyspaceName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "cassandraKeyspaceName", input)
	}

	if id.TableName, ok = input.Parsed["tableName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "tableName", input)
	}

	return nil
}

// ValidateCassandraKeyspaceTableID checks that 'input' can be parsed as a Cassandra Keyspace Table ID
func ValidateCassandraKeyspaceTableID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseCassandraKeyspaceTableID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Cassandra Keyspace Table ID
func (id CassandraKeyspaceTableId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.DocumentDB/databaseAccounts/%s/cassandraKeyspaces/%s/tables/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.DatabaseAccountName, id.CassandraKeyspaceName, id.TableName)
}

// Segments returns a slice of Resource ID Segments which comprise this Cassandra Keyspace Table ID
func (id CassandraKeyspaceTableId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftDocumentDB", "Microsoft.DocumentDB", "Microsoft.DocumentDB"),
		resourceids.StaticSegment("staticDatabaseAccounts", "databaseAccounts", "databaseAccounts"),
		resourceids.UserSpecifiedSegment("databaseAccountName", "databaseAccountName"),
		resourceids.StaticSegment("staticCassandraKeyspaces", "cassandraKeyspaces", "cassandraKeyspaces"),
		resourceids.UserSpecifiedSegment("cassandraKeyspaceName", "cassandraKeyspaceName"),
		resourceids.StaticSegment("staticTables", "tables", "tables"),
		resourceids.UserSpecifiedSegment("tableName", "tableName"),
	}
}

// String returns a human-readable description of this Cassandra Keyspace Table ID
func (id CassandraKeyspaceTableId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Database Account Name: %q", id.DatabaseAccountName),
		fmt.Sprintf("Cassandra Keyspace Name: %q", id.CassandraKeyspaceName),
		fmt.Sprintf("Table Name: %q", id.TableName),
	}
	return fmt.Sprintf("Cassandra Keyspace Table (%s)", strings.Join(components, "\n"))
}
