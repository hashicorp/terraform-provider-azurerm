// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package parse

// NOTE: this file is generated via 'go:generate' - manual changes will be overwritten

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

type ServiceVersionId struct {
	SubscriptionId string
	LocationName   string
}

func NewServiceVersionID(subscriptionId, locationName string) ServiceVersionId {
	return ServiceVersionId{
		SubscriptionId: subscriptionId,
		LocationName:   locationName,
	}
}

func (id ServiceVersionId) String() string {
	segments := []string{
		fmt.Sprintf("Location Name %q", id.LocationName),
	}
	segmentsStr := strings.Join(segments, " / ")
	return fmt.Sprintf("%s: (%s)", "Service Version", segmentsStr)
}

func (id ServiceVersionId) ID() string {
	fmtString := "/subscriptions/%s/providers/Microsoft.ContainerService/locations/%s/orchestrators"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.LocationName)
}

// ServiceVersionID parses a ServiceVersion ID into an ServiceVersionId struct
func ServiceVersionID(input string) (*ServiceVersionId, error) {
	id, err := resourceids.ParseAzureResourceID(input)
	if err != nil {
		return nil, err
	}

	resourceId := ServiceVersionId{
		SubscriptionId: id.SubscriptionID,
	}

	if resourceId.SubscriptionId == "" {
		return nil, fmt.Errorf("ID was missing the 'subscriptions' element")
	}

	if resourceId.LocationName, err = id.PopSegment("locations"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &resourceId, nil
}
