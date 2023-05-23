package cosmosdb

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ resourceids.ResourceId = SqlDatabaseId{}

// SqlDatabaseId is a struct representing the Resource ID for a Sql Database
type SqlDatabaseId struct {
	SubscriptionId      string
	ResourceGroupName   string
	DatabaseAccountName string
	SqlDatabaseName     string
}

// NewSqlDatabaseID returns a new SqlDatabaseId struct
func NewSqlDatabaseID(subscriptionId string, resourceGroupName string, databaseAccountName string, sqlDatabaseName string) SqlDatabaseId {
	return SqlDatabaseId{
		SubscriptionId:      subscriptionId,
		ResourceGroupName:   resourceGroupName,
		DatabaseAccountName: databaseAccountName,
		SqlDatabaseName:     sqlDatabaseName,
	}
}

// ParseSqlDatabaseID parses 'input' into a SqlDatabaseId
func ParseSqlDatabaseID(input string) (*SqlDatabaseId, error) {
	parser := resourceids.NewParserFromResourceIdType(SqlDatabaseId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := SqlDatabaseId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", *parsed)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", *parsed)
	}

	if id.DatabaseAccountName, ok = parsed.Parsed["databaseAccountName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "databaseAccountName", *parsed)
	}

	if id.SqlDatabaseName, ok = parsed.Parsed["sqlDatabaseName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "sqlDatabaseName", *parsed)
	}

	return &id, nil
}

// ParseSqlDatabaseIDInsensitively parses 'input' case-insensitively into a SqlDatabaseId
// note: this method should only be used for API response data and not user input
func ParseSqlDatabaseIDInsensitively(input string) (*SqlDatabaseId, error) {
	parser := resourceids.NewParserFromResourceIdType(SqlDatabaseId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := SqlDatabaseId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", *parsed)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", *parsed)
	}

	if id.DatabaseAccountName, ok = parsed.Parsed["databaseAccountName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "databaseAccountName", *parsed)
	}

	if id.SqlDatabaseName, ok = parsed.Parsed["sqlDatabaseName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "sqlDatabaseName", *parsed)
	}

	return &id, nil
}

// ValidateSqlDatabaseID checks that 'input' can be parsed as a Sql Database ID
func ValidateSqlDatabaseID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseSqlDatabaseID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Sql Database ID
func (id SqlDatabaseId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.DocumentDB/databaseAccounts/%s/sqlDatabases/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.DatabaseAccountName, id.SqlDatabaseName)
}

// Segments returns a slice of Resource ID Segments which comprise this Sql Database ID
func (id SqlDatabaseId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftDocumentDB", "Microsoft.DocumentDB", "Microsoft.DocumentDB"),
		resourceids.StaticSegment("staticDatabaseAccounts", "databaseAccounts", "databaseAccounts"),
		resourceids.UserSpecifiedSegment("databaseAccountName", "databaseAccountValue"),
		resourceids.StaticSegment("staticSqlDatabases", "sqlDatabases", "sqlDatabases"),
		resourceids.UserSpecifiedSegment("sqlDatabaseName", "sqlDatabaseValue"),
	}
}

// String returns a human-readable description of this Sql Database ID
func (id SqlDatabaseId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Database Account Name: %q", id.DatabaseAccountName),
		fmt.Sprintf("Sql Database Name: %q", id.SqlDatabaseName),
	}
	return fmt.Sprintf("Sql Database (%s)", strings.Join(components, "\n"))
}
