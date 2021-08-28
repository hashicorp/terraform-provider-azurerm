package signalr

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

type LocationId struct {
	SubscriptionId string
	Name           string
}

func NewLocationID(subscriptionId, name string) LocationId {
	return LocationId{
		SubscriptionId: subscriptionId,
		Name:           name,
	}
}

func (id LocationId) String() string {
	segments := []string{
		fmt.Sprintf("Name %q", id.Name),
	}
	segmentsStr := strings.Join(segments, " / ")
	return fmt.Sprintf("%s: (%s)", "Location", segmentsStr)
}

func (id LocationId) ID() string {
	fmtString := "/subscriptions/%s/providers/Microsoft.SignalRService/locations/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.Name)
}

// ParseLocationID parses a Location ID into an LocationId struct
func ParseLocationID(input string) (*LocationId, error) {
	id, err := resourceids.ParseAzureResourceID(input)
	if err != nil {
		return nil, err
	}

	resourceId := LocationId{
		SubscriptionId: id.SubscriptionID,
	}

	if resourceId.SubscriptionId == "" {
		return nil, fmt.Errorf("ID was missing the 'subscriptions' element")
	}

	if resourceId.Name, err = id.PopSegment("locations"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &resourceId, nil
}

// ParseLocationIDInsensitively parses an Location ID into an LocationId struct, insensitively
// This should only be used to parse an ID for rewriting to a consistent casing,
// the ParseLocationID method should be used instead for validation etc.
func ParseLocationIDInsensitively(input string) (*LocationId, error) {
	id, err := resourceids.ParseAzureResourceID(input)
	if err != nil {
		return nil, err
	}

	resourceId := LocationId{
		SubscriptionId: id.SubscriptionID,
	}

	if resourceId.SubscriptionId == "" {
		return nil, fmt.Errorf("ID was missing the 'subscriptions' element")
	}

	// find the correct casing for the 'locations' segment
	locationsKey := "locations"
	for key := range id.Path {
		if strings.EqualFold(key, locationsKey) {
			locationsKey = key
			break
		}
	}
	if resourceId.Name, err = id.PopSegment(locationsKey); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &resourceId, nil
}
