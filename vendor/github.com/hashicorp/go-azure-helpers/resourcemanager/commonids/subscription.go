package commonids

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

var _ resourceids.ResourceId = SubscriptionId{}

type SubscriptionId struct {
	SubscriptionId string
}

func NewSubscriptionID(subscriptionId string) SubscriptionId {
	return SubscriptionId{
		SubscriptionId: subscriptionId,
	}
}

func ParseSubscriptionID(input string) (*SubscriptionId, error) {
	parser := resourceids.NewParserFromResourceIdType(SubscriptionId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := SubscriptionId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, fmt.Errorf("the segment 'subscriptionId' was not found in the resource id %q", input)
	}

	return &id, nil
}

func ParseSubscriptionIDInsensitively(input string) (*SubscriptionId, error) {
	parser := resourceids.NewParserFromResourceIdType(SubscriptionId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := SubscriptionId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, fmt.Errorf("the segment 'subscriptionId' was not found in the resource id %q", input)
	}

	return &id, nil
}

func (id SubscriptionId) ID() string {
	fmtString := "/subscriptions/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId)
}

func (id SubscriptionId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("subscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
	}
}

func (id SubscriptionId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
	}
	return fmt.Sprintf("Subscription (%s)", strings.Join(components, "\n"))
}
