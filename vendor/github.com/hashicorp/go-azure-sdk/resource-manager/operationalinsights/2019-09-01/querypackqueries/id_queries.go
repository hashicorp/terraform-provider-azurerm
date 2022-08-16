package querypackqueries

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

var _ resourceids.ResourceId = QueriesId{}

// QueriesId is a struct representing the Resource ID for a Queries
type QueriesId struct {
	SubscriptionId    string
	ResourceGroupName string
	QueryPackName     string
	Id                string
}

// NewQueriesID returns a new QueriesId struct
func NewQueriesID(subscriptionId string, resourceGroupName string, queryPackName string, id string) QueriesId {
	return QueriesId{
		SubscriptionId:    subscriptionId,
		ResourceGroupName: resourceGroupName,
		QueryPackName:     queryPackName,
		Id:                id,
	}
}

// ParseQueriesID parses 'input' into a QueriesId
func ParseQueriesID(input string) (*QueriesId, error) {
	parser := resourceids.NewParserFromResourceIdType(QueriesId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := QueriesId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, fmt.Errorf("the segment 'subscriptionId' was not found in the resource id %q", input)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, fmt.Errorf("the segment 'resourceGroupName' was not found in the resource id %q", input)
	}

	if id.QueryPackName, ok = parsed.Parsed["queryPackName"]; !ok {
		return nil, fmt.Errorf("the segment 'queryPackName' was not found in the resource id %q", input)
	}

	if id.Id, ok = parsed.Parsed["id"]; !ok {
		return nil, fmt.Errorf("the segment 'id' was not found in the resource id %q", input)
	}

	return &id, nil
}

// ParseQueriesIDInsensitively parses 'input' case-insensitively into a QueriesId
// note: this method should only be used for API response data and not user input
func ParseQueriesIDInsensitively(input string) (*QueriesId, error) {
	parser := resourceids.NewParserFromResourceIdType(QueriesId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := QueriesId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, fmt.Errorf("the segment 'subscriptionId' was not found in the resource id %q", input)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, fmt.Errorf("the segment 'resourceGroupName' was not found in the resource id %q", input)
	}

	if id.QueryPackName, ok = parsed.Parsed["queryPackName"]; !ok {
		return nil, fmt.Errorf("the segment 'queryPackName' was not found in the resource id %q", input)
	}

	if id.Id, ok = parsed.Parsed["id"]; !ok {
		return nil, fmt.Errorf("the segment 'id' was not found in the resource id %q", input)
	}

	return &id, nil
}

// ValidateQueriesID checks that 'input' can be parsed as a Queries ID
func ValidateQueriesID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseQueriesID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Queries ID
func (id QueriesId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.OperationalInsights/queryPacks/%s/queries/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.QueryPackName, id.Id)
}

// Segments returns a slice of Resource ID Segments which comprise this Queries ID
func (id QueriesId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftOperationalInsights", "Microsoft.OperationalInsights", "Microsoft.OperationalInsights"),
		resourceids.StaticSegment("staticQueryPacks", "queryPacks", "queryPacks"),
		resourceids.UserSpecifiedSegment("queryPackName", "queryPackValue"),
		resourceids.StaticSegment("staticQueries", "queries", "queries"),
		resourceids.UserSpecifiedSegment("id", "idValue"),
	}
}

// String returns a human-readable description of this Queries ID
func (id QueriesId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Query Pack Name: %q", id.QueryPackName),
		fmt.Sprintf(": %q", id.Id),
	}
	return fmt.Sprintf("Queries (%s)", strings.Join(components, "\n"))
}
