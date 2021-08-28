package signalr

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

type SubscriptionId struct {
	SubscriptionId string
}

func NewSubscriptionID(subscriptionId string) SubscriptionId {
	return SubscriptionId{
		SubscriptionId: subscriptionId,
	}
}

func (id SubscriptionId) String() string {
	segments := []string{}
	segmentsStr := strings.Join(segments, " / ")
	return fmt.Sprintf("%s: (%s)", "Subscription", segmentsStr)
}

func (id SubscriptionId) ID() string {
	fmtString := "/subscriptions/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId)
}

// ParseSubscriptionID parses a Subscription ID into an SubscriptionId struct
func ParseSubscriptionID(input string) (*SubscriptionId, error) {
	id, err := resourceids.ParseAzureResourceID(input)
	if err != nil {
		return nil, err
	}

	resourceId := SubscriptionId{
		SubscriptionId: id.SubscriptionID,
	}

	if resourceId.SubscriptionId == "" {
		return nil, fmt.Errorf("ID was missing the 'subscriptions' element")
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &resourceId, nil
}

// ParseSubscriptionIDInsensitively parses an Subscription ID into an SubscriptionId struct, insensitively
// This should only be used to parse an ID for rewriting to a consistent casing,
// the ParseSubscriptionID method should be used instead for validation etc.
func ParseSubscriptionIDInsensitively(input string) (*SubscriptionId, error) {
	id, err := resourceids.ParseAzureResourceID(input)
	if err != nil {
		return nil, err
	}

	resourceId := SubscriptionId{
		SubscriptionId: id.SubscriptionID,
	}

	if resourceId.SubscriptionId == "" {
		return nil, fmt.Errorf("ID was missing the 'subscriptions' element")
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &resourceId, nil
}
