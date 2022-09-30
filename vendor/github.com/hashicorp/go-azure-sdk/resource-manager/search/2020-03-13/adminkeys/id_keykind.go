package adminkeys

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

var _ resourceids.ResourceId = KeyKindId{}

// KeyKindId is a struct representing the Resource ID for a Key Kind
type KeyKindId struct {
	SubscriptionId    string
	ResourceGroupName string
	SearchServiceName string
	KeyKind           AdminKeyKind
}

// NewKeyKindID returns a new KeyKindId struct
func NewKeyKindID(subscriptionId string, resourceGroupName string, searchServiceName string, keyKind AdminKeyKind) KeyKindId {
	return KeyKindId{
		SubscriptionId:    subscriptionId,
		ResourceGroupName: resourceGroupName,
		SearchServiceName: searchServiceName,
		KeyKind:           keyKind,
	}
}

// ParseKeyKindID parses 'input' into a KeyKindId
func ParseKeyKindID(input string) (*KeyKindId, error) {
	parser := resourceids.NewParserFromResourceIdType(KeyKindId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := KeyKindId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, fmt.Errorf("the segment 'subscriptionId' was not found in the resource id %q", input)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, fmt.Errorf("the segment 'resourceGroupName' was not found in the resource id %q", input)
	}

	if id.SearchServiceName, ok = parsed.Parsed["searchServiceName"]; !ok {
		return nil, fmt.Errorf("the segment 'searchServiceName' was not found in the resource id %q", input)
	}

	if v, ok := parsed.Parsed["keyKind"]; true {
		if !ok {
			return nil, fmt.Errorf("the segment 'keyKind' was not found in the resource id %q", input)
		}

		keyKind, err := parseAdminKeyKind(v)
		if err != nil {
			return nil, fmt.Errorf("parsing %q: %+v", v, err)
		}
		id.KeyKind = *keyKind
	}

	return &id, nil
}

// ParseKeyKindIDInsensitively parses 'input' case-insensitively into a KeyKindId
// note: this method should only be used for API response data and not user input
func ParseKeyKindIDInsensitively(input string) (*KeyKindId, error) {
	parser := resourceids.NewParserFromResourceIdType(KeyKindId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := KeyKindId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, fmt.Errorf("the segment 'subscriptionId' was not found in the resource id %q", input)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, fmt.Errorf("the segment 'resourceGroupName' was not found in the resource id %q", input)
	}

	if id.SearchServiceName, ok = parsed.Parsed["searchServiceName"]; !ok {
		return nil, fmt.Errorf("the segment 'searchServiceName' was not found in the resource id %q", input)
	}

	if v, ok := parsed.Parsed["keyKind"]; true {
		if !ok {
			return nil, fmt.Errorf("the segment 'keyKind' was not found in the resource id %q", input)
		}

		keyKind, err := parseAdminKeyKind(v)
		if err != nil {
			return nil, fmt.Errorf("parsing %q: %+v", v, err)
		}
		id.KeyKind = *keyKind
	}

	return &id, nil
}

// ValidateKeyKindID checks that 'input' can be parsed as a Key Kind ID
func ValidateKeyKindID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseKeyKindID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Key Kind ID
func (id KeyKindId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Search/searchServices/%s/regenerateAdminKey/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.SearchServiceName, string(id.KeyKind))
}

// Segments returns a slice of Resource ID Segments which comprise this Key Kind ID
func (id KeyKindId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftSearch", "Microsoft.Search", "Microsoft.Search"),
		resourceids.StaticSegment("staticSearchServices", "searchServices", "searchServices"),
		resourceids.UserSpecifiedSegment("searchServiceName", "searchServiceValue"),
		resourceids.StaticSegment("staticRegenerateAdminKey", "regenerateAdminKey", "regenerateAdminKey"),
		resourceids.ConstantSegment("keyKind", PossibleValuesForAdminKeyKind(), "primary"),
	}
}

// String returns a human-readable description of this Key Kind ID
func (id KeyKindId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Search Service Name: %q", id.SearchServiceName),
		fmt.Sprintf("Key Kind: %q", string(id.KeyKind)),
	}
	return fmt.Sprintf("Key Kind (%s)", strings.Join(components, "\n"))
}
