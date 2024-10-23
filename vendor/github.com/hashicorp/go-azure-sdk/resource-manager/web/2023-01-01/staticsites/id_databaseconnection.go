package staticsites

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/recaser"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

func init() {
	recaser.RegisterResourceId(&DatabaseConnectionId{})
}

var _ resourceids.ResourceId = &DatabaseConnectionId{}

// DatabaseConnectionId is a struct representing the Resource ID for a Database Connection
type DatabaseConnectionId struct {
	SubscriptionId         string
	ResourceGroupName      string
	StaticSiteName         string
	DatabaseConnectionName string
}

// NewDatabaseConnectionID returns a new DatabaseConnectionId struct
func NewDatabaseConnectionID(subscriptionId string, resourceGroupName string, staticSiteName string, databaseConnectionName string) DatabaseConnectionId {
	return DatabaseConnectionId{
		SubscriptionId:         subscriptionId,
		ResourceGroupName:      resourceGroupName,
		StaticSiteName:         staticSiteName,
		DatabaseConnectionName: databaseConnectionName,
	}
}

// ParseDatabaseConnectionID parses 'input' into a DatabaseConnectionId
func ParseDatabaseConnectionID(input string) (*DatabaseConnectionId, error) {
	parser := resourceids.NewParserFromResourceIdType(&DatabaseConnectionId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := DatabaseConnectionId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParseDatabaseConnectionIDInsensitively parses 'input' case-insensitively into a DatabaseConnectionId
// note: this method should only be used for API response data and not user input
func ParseDatabaseConnectionIDInsensitively(input string) (*DatabaseConnectionId, error) {
	parser := resourceids.NewParserFromResourceIdType(&DatabaseConnectionId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := DatabaseConnectionId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *DatabaseConnectionId) FromParseResult(input resourceids.ParseResult) error {
	var ok bool

	if id.SubscriptionId, ok = input.Parsed["subscriptionId"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", input)
	}

	if id.ResourceGroupName, ok = input.Parsed["resourceGroupName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", input)
	}

	if id.StaticSiteName, ok = input.Parsed["staticSiteName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "staticSiteName", input)
	}

	if id.DatabaseConnectionName, ok = input.Parsed["databaseConnectionName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "databaseConnectionName", input)
	}

	return nil
}

// ValidateDatabaseConnectionID checks that 'input' can be parsed as a Database Connection ID
func ValidateDatabaseConnectionID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseDatabaseConnectionID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Database Connection ID
func (id DatabaseConnectionId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Web/staticSites/%s/databaseConnections/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.StaticSiteName, id.DatabaseConnectionName)
}

// Segments returns a slice of Resource ID Segments which comprise this Database Connection ID
func (id DatabaseConnectionId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftWeb", "Microsoft.Web", "Microsoft.Web"),
		resourceids.StaticSegment("staticStaticSites", "staticSites", "staticSites"),
		resourceids.UserSpecifiedSegment("staticSiteName", "staticSiteName"),
		resourceids.StaticSegment("staticDatabaseConnections", "databaseConnections", "databaseConnections"),
		resourceids.UserSpecifiedSegment("databaseConnectionName", "databaseConnectionName"),
	}
}

// String returns a human-readable description of this Database Connection ID
func (id DatabaseConnectionId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Static Site Name: %q", id.StaticSiteName),
		fmt.Sprintf("Database Connection Name: %q", id.DatabaseConnectionName),
	}
	return fmt.Sprintf("Database Connection (%s)", strings.Join(components, "\n"))
}
