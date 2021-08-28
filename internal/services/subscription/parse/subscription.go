package parse

import (
	"fmt"
	"strings"

	"github.com/google/uuid"
)

type SubscriptionId struct {
	SubscriptionID string
}

// NewSubscriptionId returns a SubscriptionId object for the given Subscription UUID
func NewSubscriptionId(uuid string) SubscriptionId {
	return SubscriptionId{
		SubscriptionID: uuid,
	}
}

func (id SubscriptionId) ID() string {
	return fmt.Sprintf("/subscriptions/%s", id.SubscriptionID)
}

// SubscriptionID returns a SubscriptionId pointer for a valid input string.
// This string can be either a UUID, or an Azure Resource ID in the form
// `/subscriptions/{SubscriptionID}`
func SubscriptionID(input string) (*SubscriptionId, error) {
	id, err := uuid.Parse(strings.TrimPrefix(input, "/subscriptions/"))
	if err != nil {
		return nil, err
	}

	return &SubscriptionId{
		SubscriptionID: id.String(),
	}, err
}
