// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package parse

// NOTE: this file is generated via 'go:generate' - manual changes will be overwritten

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

type RestorableDatabaseAccountId struct {
	SubscriptionId string
	LocationName   string
	Name           string
}

func NewRestorableDatabaseAccountID(subscriptionId, locationName, name string) RestorableDatabaseAccountId {
	return RestorableDatabaseAccountId{
		SubscriptionId: subscriptionId,
		LocationName:   locationName,
		Name:           name,
	}
}

func (id RestorableDatabaseAccountId) String() string {
	segments := []string{
		fmt.Sprintf("Name %q", id.Name),
		fmt.Sprintf("Location Name %q", id.LocationName),
	}
	segmentsStr := strings.Join(segments, " / ")
	return fmt.Sprintf("%s: (%s)", "Restorable Database Account", segmentsStr)
}

func (id RestorableDatabaseAccountId) ID() string {
	fmtString := "/subscriptions/%s/providers/Microsoft.DocumentDB/locations/%s/restorableDatabaseAccounts/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.LocationName, id.Name)
}

// RestorableDatabaseAccountID parses a RestorableDatabaseAccount ID into an RestorableDatabaseAccountId struct
func RestorableDatabaseAccountID(input string) (*RestorableDatabaseAccountId, error) {
	id, err := resourceids.ParseAzureResourceID(input)
	if err != nil {
		return nil, fmt.Errorf("parsing %q as an RestorableDatabaseAccount ID: %+v", input, err)
	}

	resourceId := RestorableDatabaseAccountId{
		SubscriptionId: id.SubscriptionID,
	}

	if resourceId.SubscriptionId == "" {
		return nil, fmt.Errorf("ID was missing the 'subscriptions' element")
	}

	if resourceId.LocationName, err = id.PopSegment("locations"); err != nil {
		return nil, err
	}
	if resourceId.Name, err = id.PopSegment("restorableDatabaseAccounts"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &resourceId, nil
}
