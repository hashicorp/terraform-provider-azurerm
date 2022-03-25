package managedapis

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

var _ resourceids.ResourceId = ManagedApiId{}

// ManagedApiId is a struct representing the Resource ID for a Managed Api
type ManagedApiId struct {
	SubscriptionId string
	Location       string
	ApiName        string
}

// NewManagedApiID returns a new ManagedApiId struct
func NewManagedApiID(subscriptionId string, location string, apiName string) ManagedApiId {
	return ManagedApiId{
		SubscriptionId: subscriptionId,
		Location:       location,
		ApiName:        apiName,
	}
}

// ParseManagedApiID parses 'input' into a ManagedApiId
func ParseManagedApiID(input string) (*ManagedApiId, error) {
	parser := resourceids.NewParserFromResourceIdType(ManagedApiId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := ManagedApiId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, fmt.Errorf("the segment 'subscriptionId' was not found in the resource id %q", input)
	}

	if id.Location, ok = parsed.Parsed["location"]; !ok {
		return nil, fmt.Errorf("the segment 'location' was not found in the resource id %q", input)
	}

	if id.ApiName, ok = parsed.Parsed["apiName"]; !ok {
		return nil, fmt.Errorf("the segment 'apiName' was not found in the resource id %q", input)
	}

	return &id, nil
}

// ParseManagedApiIDInsensitively parses 'input' case-insensitively into a ManagedApiId
// note: this method should only be used for API response data and not user input
func ParseManagedApiIDInsensitively(input string) (*ManagedApiId, error) {
	parser := resourceids.NewParserFromResourceIdType(ManagedApiId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := ManagedApiId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, fmt.Errorf("the segment 'subscriptionId' was not found in the resource id %q", input)
	}

	if id.Location, ok = parsed.Parsed["location"]; !ok {
		return nil, fmt.Errorf("the segment 'location' was not found in the resource id %q", input)
	}

	if id.ApiName, ok = parsed.Parsed["apiName"]; !ok {
		return nil, fmt.Errorf("the segment 'apiName' was not found in the resource id %q", input)
	}

	return &id, nil
}

// ValidateManagedApiID checks that 'input' can be parsed as a Managed Api ID
func ValidateManagedApiID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseManagedApiID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Managed Api ID
func (id ManagedApiId) ID() string {
	fmtString := "/subscriptions/%s/providers/Microsoft.Web/locations/%s/managedApis/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.Location, id.ApiName)
}

// Segments returns a slice of Resource ID Segments which comprise this Managed Api ID
func (id ManagedApiId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftWeb", "Microsoft.Web", "Microsoft.Web"),
		resourceids.StaticSegment("staticLocations", "locations", "locations"),
		resourceids.UserSpecifiedSegment("location", "locationValue"),
		resourceids.StaticSegment("staticManagedApis", "managedApis", "managedApis"),
		resourceids.UserSpecifiedSegment("apiName", "apiValue"),
	}
}

// String returns a human-readable description of this Managed Api ID
func (id ManagedApiId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Location: %q", id.Location),
		fmt.Sprintf("Api Name: %q", id.ApiName),
	}
	return fmt.Sprintf("Managed Api (%s)", strings.Join(components, "\n"))
}
