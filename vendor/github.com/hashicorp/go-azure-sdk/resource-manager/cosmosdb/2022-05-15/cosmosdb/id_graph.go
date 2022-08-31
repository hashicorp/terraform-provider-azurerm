package cosmosdb

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

var _ resourceids.ResourceId = GraphId{}

// GraphId is a struct representing the Resource ID for a Graph
type GraphId struct {
	SubscriptionId    string
	ResourceGroupName string
	AccountName       string
	DatabaseName      string
	GraphName         string
}

// NewGraphID returns a new GraphId struct
func NewGraphID(subscriptionId string, resourceGroupName string, accountName string, databaseName string, graphName string) GraphId {
	return GraphId{
		SubscriptionId:    subscriptionId,
		ResourceGroupName: resourceGroupName,
		AccountName:       accountName,
		DatabaseName:      databaseName,
		GraphName:         graphName,
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
		return nil, fmt.Errorf("the segment 'subscriptionId' was not found in the resource id %q", input)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, fmt.Errorf("the segment 'resourceGroupName' was not found in the resource id %q", input)
	}

	if id.AccountName, ok = parsed.Parsed["accountName"]; !ok {
		return nil, fmt.Errorf("the segment 'accountName' was not found in the resource id %q", input)
	}

	if id.DatabaseName, ok = parsed.Parsed["databaseName"]; !ok {
		return nil, fmt.Errorf("the segment 'databaseName' was not found in the resource id %q", input)
	}

	if id.GraphName, ok = parsed.Parsed["graphName"]; !ok {
		return nil, fmt.Errorf("the segment 'graphName' was not found in the resource id %q", input)
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
		return nil, fmt.Errorf("the segment 'subscriptionId' was not found in the resource id %q", input)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, fmt.Errorf("the segment 'resourceGroupName' was not found in the resource id %q", input)
	}

	if id.AccountName, ok = parsed.Parsed["accountName"]; !ok {
		return nil, fmt.Errorf("the segment 'accountName' was not found in the resource id %q", input)
	}

	if id.DatabaseName, ok = parsed.Parsed["databaseName"]; !ok {
		return nil, fmt.Errorf("the segment 'databaseName' was not found in the resource id %q", input)
	}

	if id.GraphName, ok = parsed.Parsed["graphName"]; !ok {
		return nil, fmt.Errorf("the segment 'graphName' was not found in the resource id %q", input)
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
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.AccountName, id.DatabaseName, id.GraphName)
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
		resourceids.UserSpecifiedSegment("accountName", "accountValue"),
		resourceids.StaticSegment("staticGremlinDatabases", "gremlinDatabases", "gremlinDatabases"),
		resourceids.UserSpecifiedSegment("databaseName", "databaseValue"),
		resourceids.StaticSegment("staticGraphs", "graphs", "graphs"),
		resourceids.UserSpecifiedSegment("graphName", "graphValue"),
	}
}

// String returns a human-readable description of this Graph ID
func (id GraphId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Account Name: %q", id.AccountName),
		fmt.Sprintf("Database Name: %q", id.DatabaseName),
		fmt.Sprintf("Graph Name: %q", id.GraphName),
	}
	return fmt.Sprintf("Graph (%s)", strings.Join(components, "\n"))
}
