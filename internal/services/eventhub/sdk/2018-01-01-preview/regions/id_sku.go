package regions

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

var _ resourceids.ResourceId = SkuId{}

// SkuId is a struct representing the Resource ID for a Sku
type SkuId struct {
	SubscriptionId string
	Sku            string
}

// NewSkuID returns a new SkuId struct
func NewSkuID(subscriptionId string, sku string) SkuId {
	return SkuId{
		SubscriptionId: subscriptionId,
		Sku:            sku,
	}
}

// ParseSkuID parses 'input' into a SkuId
func ParseSkuID(input string) (*SkuId, error) {
	parser := resourceids.NewParserFromResourceIdType(SkuId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := SkuId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, fmt.Errorf("the segment 'subscriptionId' was not found in the resource id %q", input)
	}

	if id.Sku, ok = parsed.Parsed["sku"]; !ok {
		return nil, fmt.Errorf("the segment 'sku' was not found in the resource id %q", input)
	}

	return &id, nil
}

// ParseSkuIDInsensitively parses 'input' case-insensitively into a SkuId
// note: this method should only be used for API response data and not user input
func ParseSkuIDInsensitively(input string) (*SkuId, error) {
	parser := resourceids.NewParserFromResourceIdType(SkuId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := SkuId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, fmt.Errorf("the segment 'subscriptionId' was not found in the resource id %q", input)
	}

	if id.Sku, ok = parsed.Parsed["sku"]; !ok {
		return nil, fmt.Errorf("the segment 'sku' was not found in the resource id %q", input)
	}

	return &id, nil
}

// ValidateSkuID checks that 'input' can be parsed as a Sku ID
func ValidateSkuID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseSkuID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Sku ID
func (id SkuId) ID() string {
	fmtString := "/subscriptions/%s/providers/Microsoft.EventHub/sku/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.Sku)
}

// Segments returns a slice of Resource ID Segments which comprise this Sku ID
func (id SkuId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("subscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("providers", "providers", "providers"),
		resourceids.ResourceProviderSegment("microsoftEventHub", "Microsoft.EventHub", "Microsoft.EventHub"),
		resourceids.StaticSegment("sku", "sku", "sku"),
		resourceids.UserSpecifiedSegment("sku", "skuValue"),
	}
}

// String returns a human-readable description of this Sku ID
func (id SkuId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Sku: %q", id.Sku),
	}
	return fmt.Sprintf("Sku (%s)", strings.Join(components, "\n"))
}
