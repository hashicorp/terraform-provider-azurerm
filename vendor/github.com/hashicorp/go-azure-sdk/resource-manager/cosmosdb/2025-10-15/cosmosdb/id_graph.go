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
	recaser.RegisterResourceId(&GraphId{})
}

var _ resourceids.ResourceId = &GraphId{}

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
	parser := resourceids.NewParserFromResourceIdType(&GraphId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := GraphId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParseGraphIDInsensitively parses 'input' case-insensitively into a GraphId
// note: this method should only be used for API response data and not user input
func ParseGraphIDInsensitively(input string) (*GraphId, error) {
	parser := resourceids.NewParserFromResourceIdType(&GraphId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := GraphId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *GraphId) FromParseResult(input resourceids.ParseResult) error {
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

	if id.GremlinDatabaseName, ok = input.Parsed["gremlinDatabaseName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "gremlinDatabaseName", input)
	}

	if id.GraphName, ok = input.Parsed["graphName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "graphName", input)
	}

	return nil
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
		resourceids.UserSpecifiedSegment("databaseAccountName", "databaseAccountName"),
		resourceids.StaticSegment("staticGremlinDatabases", "gremlinDatabases", "gremlinDatabases"),
		resourceids.UserSpecifiedSegment("gremlinDatabaseName", "gremlinDatabaseName"),
		resourceids.StaticSegment("staticGraphs", "graphs", "graphs"),
		resourceids.UserSpecifiedSegment("graphName", "graphName"),
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
