package querykeys

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

var _ resourceids.ResourceId = SearchServiceId{}

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
	parser := resourceids.NewParserFromResourceIdType(SearchServiceId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := SearchServiceId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, fmt.Errorf("the segment 'subscriptionId' was not found in the resource id %q", input)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, fmt.Errorf("the segment 'resourceGroupName' was not found in the resource id %q", input)
	}

	if id.SearchServiceName, ok = parsed.Parsed["searchServiceName"]; !ok {
		return nil, fmt.Errorf("the segment 'searchServiceName' was not found in the resource id %q", input)
	}

	return &id, nil
}

// ParseSearchServiceIDInsensitively parses 'input' case-insensitively into a SearchServiceId
// note: this method should only be used for API response data and not user input
func ParseSearchServiceIDInsensitively(input string) (*SearchServiceId, error) {
	parser := resourceids.NewParserFromResourceIdType(SearchServiceId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := SearchServiceId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, fmt.Errorf("the segment 'subscriptionId' was not found in the resource id %q", input)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, fmt.Errorf("the segment 'resourceGroupName' was not found in the resource id %q", input)
	}

	if id.SearchServiceName, ok = parsed.Parsed["searchServiceName"]; !ok {
		return nil, fmt.Errorf("the segment 'searchServiceName' was not found in the resource id %q", input)
	}

	return &id, nil
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
		resourceids.UserSpecifiedSegment("searchServiceName", "searchServiceValue"),
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
