package parse

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/hashicorp/go-uuid"
)

type LocationId struct {
	SubscriptionId string
}

func NewLocationID(subscriptionId string) LocationId {
	return LocationId{
		SubscriptionId: subscriptionId,
	}
}

func (id LocationId) String() string {
	segments := []string{
		fmt.Sprintf("SubscriptionId %q", id.SubscriptionId),
	}
	segmentsStr := strings.Join(segments, " / ")
	return fmt.Sprintf("%s: (%s)", "Location", segmentsStr)
}

func (id LocationId) ID() string {
	fmtString := "/subscriptions/%s/locations"
	return fmt.Sprintf(fmtString, id.SubscriptionId)
}

func LocationID(input string) (*LocationId, error) {
	idURL, err := url.ParseRequestURI(input)
	if err != nil {
		return nil, fmt.Errorf("cannot parse Location ID: %s", err)
	}

	path := idURL.Path

	path = strings.TrimPrefix(path, "/")
	path = strings.TrimSuffix(path, "/")

	components := strings.Split(path, "/")

	resourceId := LocationId{}

	if len(components) == 3 {
		if components[0] != "subscriptions" {
			return nil, fmt.Errorf("ID should start with 'subscriptions'")
		}

		if components[2] != "locations" {
			return nil, fmt.Errorf("ID should start with 'locations'")
		}

		if components[1] == "" {
			return nil, fmt.Errorf("ID was missing 'SubscriptionID'")
		}

		if _, err := uuid.ParseUUID(components[1]); err != nil {
			return nil, fmt.Errorf("'SubscriptionID' should be valid UUID")
		}

		resourceId.SubscriptionId = components[1]
	} else {
		return nil, fmt.Errorf("ID should include 'SubscriptionID' and start with '/subscriptions' and end with '/locations'")
	}

	return &resourceId, nil
}
