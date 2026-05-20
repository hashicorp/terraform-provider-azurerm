package querypacks

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/recaser"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

func init() {
	recaser.RegisterResourceId(&QueryPackId{})
}

var _ resourceids.ResourceId = &QueryPackId{}

// QueryPackId is a struct representing the Resource ID for a Query Pack
type QueryPackId struct {
	SubscriptionId    string
	ResourceGroupName string
	QueryPackName     string
}

// NewQueryPackID returns a new QueryPackId struct
func NewQueryPackID(subscriptionId string, resourceGroupName string, queryPackName string) QueryPackId {
	return QueryPackId{
		SubscriptionId:    subscriptionId,
		ResourceGroupName: resourceGroupName,
		QueryPackName:     queryPackName,
	}
}

// ParseQueryPackID parses 'input' into a QueryPackId
func ParseQueryPackID(input string) (*QueryPackId, error) {
	parser := resourceids.NewParserFromResourceIdType(&QueryPackId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := QueryPackId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParseQueryPackIDInsensitively parses 'input' case-insensitively into a QueryPackId
// note: this method should only be used for API response data and not user input
func ParseQueryPackIDInsensitively(input string) (*QueryPackId, error) {
	parser := resourceids.NewParserFromResourceIdType(&QueryPackId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := QueryPackId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *QueryPackId) FromParseResult(input resourceids.ParseResult) error {
	var ok bool

	if id.SubscriptionId, ok = input.Parsed["subscriptionId"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", input)
	}

	if id.ResourceGroupName, ok = input.Parsed["resourceGroupName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", input)
	}

	if id.QueryPackName, ok = input.Parsed["queryPackName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "queryPackName", input)
	}

	return nil
}

// ValidateQueryPackID checks that 'input' can be parsed as a Query Pack ID
func ValidateQueryPackID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseQueryPackID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Query Pack ID
func (id QueryPackId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.OperationalInsights/queryPacks/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.QueryPackName)
}

// Segments returns a slice of Resource ID Segments which comprise this Query Pack ID
func (id QueryPackId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftOperationalInsights", "Microsoft.OperationalInsights", "Microsoft.OperationalInsights"),
		resourceids.StaticSegment("staticQueryPacks", "queryPacks", "queryPacks"),
		resourceids.UserSpecifiedSegment("queryPackName", "queryPackName"),
	}
}

// String returns a human-readable description of this Query Pack ID
func (id QueryPackId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Query Pack Name: %q", id.QueryPackName),
	}
	return fmt.Sprintf("Query Pack (%s)", strings.Join(components, "\n"))
}
