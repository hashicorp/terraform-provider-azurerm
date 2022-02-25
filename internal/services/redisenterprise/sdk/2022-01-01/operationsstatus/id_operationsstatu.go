package operationsstatus

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

var _ resourceids.ResourceId = OperationsStatuId{}

// OperationsStatuId is a struct representing the Resource ID for a Operations Statu
type OperationsStatuId struct {
	SubscriptionId string
	Location       string
	OperationId    string
}

// NewOperationsStatuID returns a new OperationsStatuId struct
func NewOperationsStatuID(subscriptionId string, location string, operationId string) OperationsStatuId {
	return OperationsStatuId{
		SubscriptionId: subscriptionId,
		Location:       location,
		OperationId:    operationId,
	}
}

// ParseOperationsStatuID parses 'input' into a OperationsStatuId
func ParseOperationsStatuID(input string) (*OperationsStatuId, error) {
	parser := resourceids.NewParserFromResourceIdType(OperationsStatuId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := OperationsStatuId{}

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

// ParseOperationsStatuIDInsensitively parses 'input' case-insensitively into a OperationsStatuId
// note: this method should only be used for API response data and not user input
func ParseOperationsStatuIDInsensitively(input string) (*OperationsStatuId, error) {
	parser := resourceids.NewParserFromResourceIdType(OperationsStatuId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := OperationsStatuId{}

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

// ValidateOperationsStatuID checks that 'input' can be parsed as a Operations Statu ID
func ValidateOperationsStatuID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseOperationsStatuID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Operations Statu ID
func (id OperationsStatuId) ID() string {
	fmtString := "/subscriptions/%s/providers/Microsoft.Cache/locations/%s/operationsStatus/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.Location, id.OperationId)
}

// Segments returns a slice of Resource ID Segments which comprise this Operations Statu ID
func (id OperationsStatuId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftCache", "Microsoft.Cache", "Microsoft.Cache"),
		resourceids.StaticSegment("staticLocations", "locations", "locations"),
		resourceids.UserSpecifiedSegment("location", "locationValue"),
		resourceids.StaticSegment("staticOperationsStatus", "operationsStatus", "operationsStatus"),
		resourceids.UserSpecifiedSegment("operationId", "operationIdValue"),
	}
}

// String returns a human-readable description of this Operations Statu ID
func (id OperationsStatuId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Location: %q", id.Location),
		fmt.Sprintf("Operation: %q", id.OperationId),
	}
	return fmt.Sprintf("Operations Statu (%s)", strings.Join(components, "\n"))
}
