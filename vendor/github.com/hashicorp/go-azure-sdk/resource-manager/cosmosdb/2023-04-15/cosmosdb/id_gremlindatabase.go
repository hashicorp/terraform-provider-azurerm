package cosmosdb

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ resourceids.ResourceId = GremlinDatabaseId{}

// GremlinDatabaseId is a struct representing the Resource ID for a Gremlin Database
type GremlinDatabaseId struct {
	SubscriptionId      string
	ResourceGroupName   string
	DatabaseAccountName string
	GremlinDatabaseName string
}

// NewGremlinDatabaseID returns a new GremlinDatabaseId struct
func NewGremlinDatabaseID(subscriptionId string, resourceGroupName string, databaseAccountName string, gremlinDatabaseName string) GremlinDatabaseId {
	return GremlinDatabaseId{
		SubscriptionId:      subscriptionId,
		ResourceGroupName:   resourceGroupName,
		DatabaseAccountName: databaseAccountName,
		GremlinDatabaseName: gremlinDatabaseName,
	}
}

// ParseGremlinDatabaseID parses 'input' into a GremlinDatabaseId
func ParseGremlinDatabaseID(input string) (*GremlinDatabaseId, error) {
	parser := resourceids.NewParserFromResourceIdType(GremlinDatabaseId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := GremlinDatabaseId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", *parsed)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", *parsed)
	}

	if id.DatabaseAccountName, ok = parsed.Parsed["databaseAccountName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "databaseAccountName", *parsed)
	}

	if id.GremlinDatabaseName, ok = parsed.Parsed["gremlinDatabaseName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "gremlinDatabaseName", *parsed)
	}

	return &id, nil
}

// ParseGremlinDatabaseIDInsensitively parses 'input' case-insensitively into a GremlinDatabaseId
// note: this method should only be used for API response data and not user input
func ParseGremlinDatabaseIDInsensitively(input string) (*GremlinDatabaseId, error) {
	parser := resourceids.NewParserFromResourceIdType(GremlinDatabaseId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := GremlinDatabaseId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", *parsed)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", *parsed)
	}

	if id.DatabaseAccountName, ok = parsed.Parsed["databaseAccountName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "databaseAccountName", *parsed)
	}

	if id.GremlinDatabaseName, ok = parsed.Parsed["gremlinDatabaseName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "gremlinDatabaseName", *parsed)
	}

	return &id, nil
}

// ValidateGremlinDatabaseID checks that 'input' can be parsed as a Gremlin Database ID
func ValidateGremlinDatabaseID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseGremlinDatabaseID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Gremlin Database ID
func (id GremlinDatabaseId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.DocumentDB/databaseAccounts/%s/gremlinDatabases/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.DatabaseAccountName, id.GremlinDatabaseName)
}

// Segments returns a slice of Resource ID Segments which comprise this Gremlin Database ID
func (id GremlinDatabaseId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftDocumentDB", "Microsoft.DocumentDB", "Microsoft.DocumentDB"),
		resourceids.StaticSegment("staticDatabaseAccounts", "databaseAccounts", "databaseAccounts"),
		resourceids.UserSpecifiedSegment("databaseAccountName", "databaseAccountValue"),
		resourceids.StaticSegment("staticGremlinDatabases", "gremlinDatabases", "gremlinDatabases"),
		resourceids.UserSpecifiedSegment("gremlinDatabaseName", "gremlinDatabaseValue"),
	}
}

// String returns a human-readable description of this Gremlin Database ID
func (id GremlinDatabaseId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Database Account Name: %q", id.DatabaseAccountName),
		fmt.Sprintf("Gremlin Database Name: %q", id.GremlinDatabaseName),
	}
	return fmt.Sprintf("Gremlin Database (%s)", strings.Join(components, "\n"))
}
