package querykeys

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

var _ resourceids.ResourceId = CreateQueryKeyId{}

// CreateQueryKeyId is a struct representing the Resource ID for a Create Query Key
type CreateQueryKeyId struct {
	SubscriptionId    string
	ResourceGroupName string
	SearchServiceName string
	Name              string
}

// NewCreateQueryKeyID returns a new CreateQueryKeyId struct
func NewCreateQueryKeyID(subscriptionId string, resourceGroupName string, searchServiceName string, name string) CreateQueryKeyId {
	return CreateQueryKeyId{
		SubscriptionId:    subscriptionId,
		ResourceGroupName: resourceGroupName,
		SearchServiceName: searchServiceName,
		Name:              name,
	}
}

// ParseCreateQueryKeyID parses 'input' into a CreateQueryKeyId
func ParseCreateQueryKeyID(input string) (*CreateQueryKeyId, error) {
	parser := resourceids.NewParserFromResourceIdType(CreateQueryKeyId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := CreateQueryKeyId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, fmt.Errorf("the segment 'subscriptionId' was not found in the resource id %q", input)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, fmt.Errorf("the segment 'resourceGroupName' was not found in the resource id %q", input)
	}

	if id.SearchServiceName, ok = parsed.Parsed["searchServiceName"]; !ok {
		return nil, fmt.Errorf("the segment 'searchServiceName' was not found in the resource id %q", input)
	}

	if id.Name, ok = parsed.Parsed["name"]; !ok {
		return nil, fmt.Errorf("the segment 'name' was not found in the resource id %q", input)
	}

	return &id, nil
}

// ParseCreateQueryKeyIDInsensitively parses 'input' case-insensitively into a CreateQueryKeyId
// note: this method should only be used for API response data and not user input
func ParseCreateQueryKeyIDInsensitively(input string) (*CreateQueryKeyId, error) {
	parser := resourceids.NewParserFromResourceIdType(CreateQueryKeyId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := CreateQueryKeyId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, fmt.Errorf("the segment 'subscriptionId' was not found in the resource id %q", input)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, fmt.Errorf("the segment 'resourceGroupName' was not found in the resource id %q", input)
	}

	if id.SearchServiceName, ok = parsed.Parsed["searchServiceName"]; !ok {
		return nil, fmt.Errorf("the segment 'searchServiceName' was not found in the resource id %q", input)
	}

	if id.Name, ok = parsed.Parsed["name"]; !ok {
		return nil, fmt.Errorf("the segment 'name' was not found in the resource id %q", input)
	}

	return &id, nil
}

// ValidateCreateQueryKeyID checks that 'input' can be parsed as a Create Query Key ID
func ValidateCreateQueryKeyID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseCreateQueryKeyID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Create Query Key ID
func (id CreateQueryKeyId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Search/searchServices/%s/createQueryKey/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.SearchServiceName, id.Name)
}

// Segments returns a slice of Resource ID Segments which comprise this Create Query Key ID
func (id CreateQueryKeyId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftSearch", "Microsoft.Search", "Microsoft.Search"),
		resourceids.StaticSegment("staticSearchServices", "searchServices", "searchServices"),
		resourceids.UserSpecifiedSegment("searchServiceName", "searchServiceValue"),
		resourceids.StaticSegment("staticCreateQueryKey", "createQueryKey", "createQueryKey"),
		resourceids.UserSpecifiedSegment("name", "nameValue"),
	}
}

// String returns a human-readable description of this Create Query Key ID
func (id CreateQueryKeyId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Search Service Name: %q", id.SearchServiceName),
		fmt.Sprintf("Name: %q", id.Name),
	}
	return fmt.Sprintf("Create Query Key (%s)", strings.Join(components, "\n"))
}
