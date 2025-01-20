// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package parse

// NOTE: this file is generated via 'go:generate' - manual changes will be overwritten

import (
	"errors"
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

type SkusId struct {
	SubscriptionId string
}

func NewSkusID(subscriptionId string) SkusId {
	return SkusId{
		SubscriptionId: subscriptionId,
	}
}

func (id SkusId) String() string {
	segments := []string{}
	segmentsStr := strings.Join(segments, " / ")
	return fmt.Sprintf("%s: (%s)", "Skus", segmentsStr)
}

func (id SkusId) ID() string {
	fmtString := "/subscriptions/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId)
}

// SkusID parses a Skus ID into an SkusId struct
func SkusID(input string) (*SkusId, error) {
	id, err := resourceids.ParseAzureResourceID(input)
	if err != nil {
		return nil, fmt.Errorf("parsing %q as an Skus ID: %+v", input, err)
	}

	resourceId := SkusId{
		SubscriptionId: id.SubscriptionID,
	}

	if resourceId.SubscriptionId == "" {
		return nil, errors.New("ID was missing the 'subscriptions' element")
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &resourceId, nil
}
