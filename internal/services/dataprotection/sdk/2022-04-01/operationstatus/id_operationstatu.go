package operationstatus

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

var _ resourceids.ResourceId = OperationStatuId{}

// OperationStatuId is a struct representing the Resource ID for a Operation Statu
type OperationStatuId struct {
	SubscriptionId string
	Location       string
	OperationId    string
}

// NewOperationStatuID returns a new OperationStatuId struct
func NewOperationStatuID(subscriptionId string, location string, operationId string) OperationStatuId {
	return OperationStatuId{
		SubscriptionId: subscriptionId,
		Location:       location,
		OperationId:    operationId,
	}
}

// ParseOperationStatuID parses 'input' into a OperationStatuId
func ParseOperationStatuID(input string) (*OperationStatuId, error) {
	parser := resourceids.NewParserFromResourceIdType(OperationStatuId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := OperationStatuId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, fmt.Errorf("the segment 'subscriptionId' was not found in the resource id %q", input)
	}

	if id.Location, ok = parsed.Parsed["location"]; !ok {
		return nil, fmt.Errorf("the segment 'location' was not found in the resource id %q", input)
	}

	if id.OperationId, ok = parsed.Parsed["operationId"]; !ok {
		return nil, fmt.Errorf("the segment 'operationId' was not found in the resource id %q", input)
	}

	return &id, nil
}

// ParseOperationStatuIDInsensitively parses 'input' case-insensitively into a OperationStatuId
// note: this method should only be used for API response data and not user input
func ParseOperationStatuIDInsensitively(input string) (*OperationStatuId, error) {
	parser := resourceids.NewParserFromResourceIdType(OperationStatuId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := OperationStatuId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, fmt.Errorf("the segment 'subscriptionId' was not found in the resource id %q", input)
	}

	if id.Location, ok = parsed.Parsed["location"]; !ok {
		return nil, fmt.Errorf("the segment 'location' was not found in the resource id %q", input)
	}

	if id.OperationId, ok = parsed.Parsed["operationId"]; !ok {
		return nil, fmt.Errorf("the segment 'operationId' was not found in the resource id %q", input)
	}

	return &id, nil
}

// ValidateOperationStatuID checks that 'input' can be parsed as a Operation Statu ID
func ValidateOperationStatuID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseOperationStatuID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Operation Statu ID
func (id OperationStatuId) ID() string {
	fmtString := "/subscriptions/%s/providers/Microsoft.DataProtection/locations/%s/operationStatus/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.Location, id.OperationId)
}

// Segments returns a slice of Resource ID Segments which comprise this Operation Statu ID
func (id OperationStatuId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftDataProtection", "Microsoft.DataProtection", "Microsoft.DataProtection"),
		resourceids.StaticSegment("staticLocations", "locations", "locations"),
		resourceids.UserSpecifiedSegment("location", "locationValue"),
		resourceids.StaticSegment("staticOperationStatus", "operationStatus", "operationStatus"),
		resourceids.UserSpecifiedSegment("operationId", "operationIdValue"),
	}
}

// String returns a human-readable description of this Operation Statu ID
func (id OperationStatuId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Location: %q", id.Location),
		fmt.Sprintf("Operation: %q", id.OperationId),
	}
	return fmt.Sprintf("Operation Statu (%s)", strings.Join(components, "\n"))
}
