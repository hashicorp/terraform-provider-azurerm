package accounts

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

var _ resourceids.ResourceId = LocationId{}

// LocationId is a struct representing the Resource ID for a Location
type LocationId struct {
	SubscriptionId string
	Location       string
}

// NewLocationID returns a new LocationId struct
func NewLocationID(subscriptionId string, location string) LocationId {
	return LocationId{
		SubscriptionId: subscriptionId,
		Location:       location,
	}
}

// ParseLocationID parses 'input' into a LocationId
func ParseLocationID(input string) (*LocationId, error) {
	parser := resourceids.NewParserFromResourceIdType(LocationId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := LocationId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, fmt.Errorf("the segment 'subscriptionId' was not found in the resource id %q", input)
	}

	if id.Location, ok = parsed.Parsed["location"]; !ok {
		return nil, fmt.Errorf("the segment 'location' was not found in the resource id %q", input)
	}

	return &id, nil
}

// ParseLocationIDInsensitively parses 'input' case-insensitively into a LocationId
// note: this method should only be used for API response data and not user input
func ParseLocationIDInsensitively(input string) (*LocationId, error) {
	parser := resourceids.NewParserFromResourceIdType(LocationId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := LocationId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, fmt.Errorf("the segment 'subscriptionId' was not found in the resource id %q", input)
	}

	if id.Location, ok = parsed.Parsed["location"]; !ok {
		return nil, fmt.Errorf("the segment 'location' was not found in the resource id %q", input)
	}

	return &id, nil
}

// ValidateLocationID checks that 'input' can be parsed as a Location ID
func ValidateLocationID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseLocationID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Location ID
func (id LocationId) ID() string {
	fmtString := "/subscriptions/%s/providers/Microsoft.DataLakeAnalytics/locations/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.Location)
}

// Segments returns a slice of Resource ID Segments which comprise this Location ID
func (id LocationId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftDataLakeAnalytics", "Microsoft.DataLakeAnalytics", "Microsoft.DataLakeAnalytics"),
		resourceids.StaticSegment("staticLocations", "locations", "locations"),
		resourceids.UserSpecifiedSegment("location", "locationValue"),
	}
}

// String returns a human-readable description of this Location ID
func (id LocationId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Location: %q", id.Location),
	}
	return fmt.Sprintf("Location (%s)", strings.Join(components, "\n"))
}
