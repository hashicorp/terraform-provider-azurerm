package querykeys

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/recaser"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

func init() {
	recaser.RegisterResourceId(&SearchServiceId{})
}

var _ resourceids.ResourceId = &SearchServiceId{}

// SearchServiceId is a struct representing the Resource ID for a Search Service
type SearchServiceId struct {
	SubscriptionId    string
	ResourceGroupName string
	SearchServiceName string
}

// NewSearchServiceID returns a new SearchServiceId struct
func NewSearchServiceID(subscriptionId string, resourceGroupName string, searchServiceName string) SearchServiceId {
	return SearchServiceId{
		SubscriptionId:    subscriptionId,
		ResourceGroupName: resourceGroupName,
		SearchServiceName: searchServiceName,
	}
}

// ParseSearchServiceID parses 'input' into a SearchServiceId
func ParseSearchServiceID(input string) (*SearchServiceId, error) {
	parser := resourceids.NewParserFromResourceIdType(&SearchServiceId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := SearchServiceId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParseSearchServiceIDInsensitively parses 'input' case-insensitively into a SearchServiceId
// note: this method should only be used for API response data and not user input
func ParseSearchServiceIDInsensitively(input string) (*SearchServiceId, error) {
	parser := resourceids.NewParserFromResourceIdType(&SearchServiceId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := SearchServiceId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *SearchServiceId) FromParseResult(input resourceids.ParseResult) error {
	var ok bool

	if id.SubscriptionId, ok = input.Parsed["subscriptionId"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", input)
	}

	if id.ResourceGroupName, ok = input.Parsed["resourceGroupName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", input)
	}

	if id.SearchServiceName, ok = input.Parsed["searchServiceName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "searchServiceName", input)
	}

	return nil
}

// ValidateSearchServiceID checks that 'input' can be parsed as a Search Service ID
func ValidateSearchServiceID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseSearchServiceID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Search Service ID
func (id SearchServiceId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Search/searchServices/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.SearchServiceName)
}

// Segments returns a slice of Resource ID Segments which comprise this Search Service ID
func (id SearchServiceId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftSearch", "Microsoft.Search", "Microsoft.Search"),
		resourceids.StaticSegment("staticSearchServices", "searchServices", "searchServices"),
		resourceids.UserSpecifiedSegment("searchServiceName", "searchServiceName"),
	}
}

// String returns a human-readable description of this Search Service ID
func (id SearchServiceId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Search Service Name: %q", id.SearchServiceName),
	}
	return fmt.Sprintf("Search Service (%s)", strings.Join(components, "\n"))
}
