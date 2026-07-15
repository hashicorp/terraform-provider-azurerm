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
	recaser.RegisterResourceId(&DatabaseId{})
}

var _ resourceids.ResourceId = &DatabaseId{}

// DatabaseId is a struct representing the Resource ID for a Database
type DatabaseId struct {
	SubscriptionId     string
	ResourceGroupName  string
	FlexibleServerName string
	DatabaseName       string
}

// NewDatabaseID returns a new DatabaseId struct
func NewDatabaseID(subscriptionId string, resourceGroupName string, flexibleServerName string, databaseName string) DatabaseId {
	return DatabaseId{
		SubscriptionId:     subscriptionId,
		ResourceGroupName:  resourceGroupName,
		FlexibleServerName: flexibleServerName,
		DatabaseName:       databaseName,
	}
}

// ParseDatabaseID parses 'input' into a DatabaseId
func ParseDatabaseID(input string) (*DatabaseId, error) {
	parser := resourceids.NewParserFromResourceIdType(&DatabaseId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := DatabaseId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParseDatabaseIDInsensitively parses 'input' case-insensitively into a DatabaseId
// note: this method should only be used for API response data and not user input
func ParseDatabaseIDInsensitively(input string) (*DatabaseId, error) {
	parser := resourceids.NewParserFromResourceIdType(&DatabaseId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := DatabaseId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *DatabaseId) FromParseResult(input resourceids.ParseResult) error {
	var ok bool

	if id.SubscriptionId, ok = input.Parsed["subscriptionId"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", input)
	}

	if id.ResourceGroupName, ok = input.Parsed["resourceGroupName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", input)
	}

	if id.FlexibleServerName, ok = input.Parsed["flexibleServerName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "flexibleServerName", input)
	}

	if id.DatabaseName, ok = input.Parsed["databaseName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "databaseName", input)
	}

	return nil
}

// ValidateDatabaseID checks that 'input' can be parsed as a Database ID
func ValidateDatabaseID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseDatabaseID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Database ID
func (id DatabaseId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.DBforPostgreSQL/flexibleServers/%s/databases/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.FlexibleServerName, id.DatabaseName)
}

// Segments returns a slice of Resource ID Segments which comprise this Database ID
func (id DatabaseId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftDBforPostgreSQL", "Microsoft.DBforPostgreSQL", "Microsoft.DBforPostgreSQL"),
		resourceids.StaticSegment("staticFlexibleServers", "flexibleServers", "flexibleServers"),
		resourceids.UserSpecifiedSegment("flexibleServerName", "flexibleServerName"),
		resourceids.StaticSegment("staticDatabases", "databases", "databases"),
		resourceids.UserSpecifiedSegment("databaseName", "databaseName"),
	}
}

// String returns a human-readable description of this Database ID
func (id DatabaseId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Flexible Server Name: %q", id.FlexibleServerName),
		fmt.Sprintf("Database Name: %q", id.DatabaseName),
	}
	return fmt.Sprintf("Database (%s)", strings.Join(components, "\n"))
}
