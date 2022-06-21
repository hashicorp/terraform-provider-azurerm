package querykeys

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

var _ resourceids.ResourceId = DeleteQueryKeyId{}

// DeleteQueryKeyId is a struct representing the Resource ID for a Delete Query Key
type DeleteQueryKeyId struct {
	SubscriptionId    string
	ResourceGroupName string
	SearchServiceName string
	Key               string
}

// NewDeleteQueryKeyID returns a new DeleteQueryKeyId struct
func NewDeleteQueryKeyID(subscriptionId string, resourceGroupName string, searchServiceName string, key string) DeleteQueryKeyId {
	return DeleteQueryKeyId{
		SubscriptionId:    subscriptionId,
		ResourceGroupName: resourceGroupName,
		SearchServiceName: searchServiceName,
		Key:               key,
	}
}

// ParseDeleteQueryKeyID parses 'input' into a DeleteQueryKeyId
func ParseDeleteQueryKeyID(input string) (*DeleteQueryKeyId, error) {
	parser := resourceids.NewParserFromResourceIdType(DeleteQueryKeyId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := DeleteQueryKeyId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, fmt.Errorf("the segment 'subscriptionId' was not found in the resource id %q", input)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, fmt.Errorf("the segment 'resourceGroupName' was not found in the resource id %q", input)
	}

	if id.SearchServiceName, ok = parsed.Parsed["searchServiceName"]; !ok {
		return nil, fmt.Errorf("the segment 'searchServiceName' was not found in the resource id %q", input)
	}

	if id.Key, ok = parsed.Parsed["key"]; !ok {
		return nil, fmt.Errorf("the segment 'key' was not found in the resource id %q", input)
	}

	return &id, nil
}

// ParseDeleteQueryKeyIDInsensitively parses 'input' case-insensitively into a DeleteQueryKeyId
// note: this method should only be used for API response data and not user input
func ParseDeleteQueryKeyIDInsensitively(input string) (*DeleteQueryKeyId, error) {
	parser := resourceids.NewParserFromResourceIdType(DeleteQueryKeyId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := DeleteQueryKeyId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, fmt.Errorf("the segment 'subscriptionId' was not found in the resource id %q", input)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, fmt.Errorf("the segment 'resourceGroupName' was not found in the resource id %q", input)
	}

	if id.SearchServiceName, ok = parsed.Parsed["searchServiceName"]; !ok {
		return nil, fmt.Errorf("the segment 'searchServiceName' was not found in the resource id %q", input)
	}

	if id.Key, ok = parsed.Parsed["key"]; !ok {
		return nil, fmt.Errorf("the segment 'key' was not found in the resource id %q", input)
	}

	return &id, nil
}

// ValidateDeleteQueryKeyID checks that 'input' can be parsed as a Delete Query Key ID
func ValidateDeleteQueryKeyID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseDeleteQueryKeyID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Delete Query Key ID
func (id DeleteQueryKeyId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Search/searchServices/%s/deleteQueryKey/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.SearchServiceName, id.Key)
}

// Segments returns a slice of Resource ID Segments which comprise this Delete Query Key ID
func (id DeleteQueryKeyId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftSearch", "Microsoft.Search", "Microsoft.Search"),
		resourceids.StaticSegment("staticSearchServices", "searchServices", "searchServices"),
		resourceids.UserSpecifiedSegment("searchServiceName", "searchServiceValue"),
		resourceids.StaticSegment("staticDeleteQueryKey", "deleteQueryKey", "deleteQueryKey"),
		resourceids.UserSpecifiedSegment("key", "keyValue"),
	}
}

// String returns a human-readable description of this Delete Query Key ID
func (id DeleteQueryKeyId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Search Service Name: %q", id.SearchServiceName),
		fmt.Sprintf("Key: %q", id.Key),
	}
	return fmt.Sprintf("Delete Query Key (%s)", strings.Join(components, "\n"))
}
