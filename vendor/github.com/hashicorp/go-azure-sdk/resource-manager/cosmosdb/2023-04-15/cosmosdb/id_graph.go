package cosmosdb

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ resourceids.ResourceId = GraphId{}

// GraphId is a struct representing the Resource ID for a Graph
type GraphId struct {
	SubscriptionId      string
	ResourceGroupName   string
	DatabaseAccountName string
	GremlinDatabaseName string
	GraphName           string
}

// NewGraphID returns a new GraphId struct
func NewGraphID(subscriptionId string, resourceGroupName string, databaseAccountName string, gremlinDatabaseName string, graphName string) GraphId {
	return GraphId{
		SubscriptionId:      subscriptionId,
		ResourceGroupName:   resourceGroupName,
		DatabaseAccountName: databaseAccountName,
		GremlinDatabaseName: gremlinDatabaseName,
		GraphName:           graphName,
	}
}

// ParseGraphID parses 'input' into a GraphId
func ParseGraphID(input string) (*GraphId, error) {
	parser := resourceids.NewParserFromResourceIdType(GraphId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := GraphId{}

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

	if id.GraphName, ok = parsed.Parsed["graphName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "graphName", *parsed)
	}

	return &id, nil
}

// ParseGraphIDInsensitively parses 'input' case-insensitively into a GraphId
// note: this method should only be used for API response data and not user input
func ParseGraphIDInsensitively(input string) (*GraphId, error) {
	parser := resourceids.NewParserFromResourceIdType(GraphId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := GraphId{}

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

	if id.GraphName, ok = parsed.Parsed["graphName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "graphName", *parsed)
	}

	return &id, nil
}

// ValidateGraphID checks that 'input' can be parsed as a Graph ID
func ValidateGraphID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseGraphID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Graph ID
func (id GraphId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.DocumentDB/databaseAccounts/%s/gremlinDatabases/%s/graphs/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.DatabaseAccountName, id.GremlinDatabaseName, id.GraphName)
}

// Segments returns a slice of Resource ID Segments which comprise this Graph ID
func (id GraphId) Segments() []resourceids.Segment {
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
		resourceids.StaticSegment("staticGraphs", "graphs", "graphs"),
		resourceids.UserSpecifiedSegment("graphName", "graphValue"),
	}
}

// String returns a human-readable description of this Graph ID
func (id GraphId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Database Account Name: %q", id.DatabaseAccountName),
		fmt.Sprintf("Gremlin Database Name: %q", id.GremlinDatabaseName),
		fmt.Sprintf("Graph Name: %q", id.GraphName),
	}
	return fmt.Sprintf("Graph (%s)", strings.Join(components, "\n"))
}
