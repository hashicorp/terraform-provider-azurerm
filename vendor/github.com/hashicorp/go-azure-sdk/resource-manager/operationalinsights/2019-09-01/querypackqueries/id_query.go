package querypackqueries

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ resourceids.ResourceId = QueryId{}

// QueryId is a struct representing the Resource ID for a Query
type QueryId struct {
	SubscriptionId    string
	ResourceGroupName string
	QueryPackName     string
	QueryName         string
}

// NewQueryID returns a new QueryId struct
func NewQueryID(subscriptionId string, resourceGroupName string, queryPackName string, queryName string) QueryId {
	return QueryId{
		SubscriptionId:    subscriptionId,
		ResourceGroupName: resourceGroupName,
		QueryPackName:     queryPackName,
		QueryName:         queryName,
	}
}

// ParseQueryID parses 'input' into a QueryId
func ParseQueryID(input string) (*QueryId, error) {
	parser := resourceids.NewParserFromResourceIdType(QueryId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := QueryId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", *parsed)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", *parsed)
	}

	if id.QueryPackName, ok = parsed.Parsed["queryPackName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "queryPackName", *parsed)
	}

	if id.QueryName, ok = parsed.Parsed["queryName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "queryName", *parsed)
	}

	return &id, nil
}

// ParseQueryIDInsensitively parses 'input' case-insensitively into a QueryId
// note: this method should only be used for API response data and not user input
func ParseQueryIDInsensitively(input string) (*QueryId, error) {
	parser := resourceids.NewParserFromResourceIdType(QueryId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := QueryId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", *parsed)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", *parsed)
	}

	if id.QueryPackName, ok = parsed.Parsed["queryPackName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "queryPackName", *parsed)
	}

	if id.QueryName, ok = parsed.Parsed["queryName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "queryName", *parsed)
	}

	return &id, nil
}

// ValidateQueryID checks that 'input' can be parsed as a Query ID
func ValidateQueryID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseQueryID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Query ID
func (id QueryId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.OperationalInsights/queryPacks/%s/queries/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.QueryPackName, id.QueryName)
}

// Segments returns a slice of Resource ID Segments which comprise this Query ID
func (id QueryId) Segments() []resourceids.Segment {
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
		resourceids.UserSpecifiedSegment("queryName", "queryValue"),
	}
}

// String returns a human-readable description of this Query ID
func (id QueryId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Query Pack Name: %q", id.QueryPackName),
		fmt.Sprintf("Query Name: %q", id.QueryName),
	}
	return fmt.Sprintf("Query (%s)", strings.Join(components, "\n"))
}
