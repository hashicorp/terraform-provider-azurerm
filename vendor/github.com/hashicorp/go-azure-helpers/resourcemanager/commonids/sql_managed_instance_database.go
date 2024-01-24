// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package commonids

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

var _ resourceids.ResourceId = &SqlManagedInstanceDatabaseId{}

// SqlManagedInstanceDatabaseId is a struct representing the Resource ID for a Sql Managed Instance Database
type SqlManagedInstanceDatabaseId struct {
	SubscriptionId      string
	ResourceGroupName   string
	ManagedInstanceName string
	DatabaseName        string
}

// NewSqlManagedInstanceDatabaseID returns a new SqlManagedInstanceDatabaseId struct
func NewSqlManagedInstanceDatabaseID(subscriptionId string, resourceGroupName string, managedInstanceName string, databaseName string) SqlManagedInstanceDatabaseId {
	return SqlManagedInstanceDatabaseId{
		SubscriptionId:      subscriptionId,
		ResourceGroupName:   resourceGroupName,
		ManagedInstanceName: managedInstanceName,
		DatabaseName:        databaseName,
	}
}

// ParseManagedInstanceDatabaseID parses 'input' into a SqlManagedInstanceDatabaseId
func ParseManagedInstanceDatabaseID(input string) (*SqlManagedInstanceDatabaseId, error) {
	parser := resourceids.NewParserFromResourceIdType(&SqlManagedInstanceDatabaseId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := SqlManagedInstanceDatabaseId{}
	if err := id.FromParseResult(*parsed); err != nil {
		return nil, err
	}
	return &id, nil
}

// ParseSqlManagedInstanceDatabaseIDInsensitively parses 'input' case-insensitively into a SqlManagedInstanceDatabaseId
// note: this method should only be used for API response data and not user input
func ParseSqlManagedInstanceDatabaseIDInsensitively(input string) (*SqlManagedInstanceDatabaseId, error) {
	parser := resourceids.NewParserFromResourceIdType(&SqlManagedInstanceDatabaseId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := SqlManagedInstanceDatabaseId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *SqlManagedInstanceDatabaseId) FromParseResult(input resourceids.ParseResult) error {
	var ok bool

	if id.SubscriptionId, ok = input.Parsed["subscriptionId"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", input)
	}

	if id.ResourceGroupName, ok = input.Parsed["resourceGroupName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", input)
	}

	if id.ManagedInstanceName, ok = input.Parsed["managedInstanceName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "managedInstanceName", input)
	}

	if id.DatabaseName, ok = input.Parsed["databaseName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "databaseName", input)
	}

	return nil
}

// ValidateSqlManagedInstanceDatabaseID checks that 'input' can be parsed as a Sql Managed Instance Database ID
func ValidateSqlManagedInstanceDatabaseID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseManagedInstanceDatabaseID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Sql Managed Instance Database ID
func (id SqlManagedInstanceDatabaseId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Sql/managedInstances/%s/databases/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.ManagedInstanceName, id.DatabaseName)
}

// Segments returns a slice of Resource ID Segments which comprise this Sql Managed Instance Database ID
func (id SqlManagedInstanceDatabaseId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftSql", "Microsoft.Sql", "Microsoft.Sql"),
		resourceids.StaticSegment("staticManagedInstances", "managedInstances", "managedInstances"),
		resourceids.UserSpecifiedSegment("managedInstanceName", "managedInstanceValue"),
		resourceids.StaticSegment("staticDatabases", "databases", "databases"),
		resourceids.UserSpecifiedSegment("databaseName", "databaseValue"),
	}
}

// String returns a human-readable description of this Sql Managed Instance Database ID
func (id SqlManagedInstanceDatabaseId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Managed Instance Name: %q", id.ManagedInstanceName),
		fmt.Sprintf("Database Name: %q", id.DatabaseName),
	}
	return fmt.Sprintf("Managed Instance Database (%s)", strings.Join(components, "\n"))
}
