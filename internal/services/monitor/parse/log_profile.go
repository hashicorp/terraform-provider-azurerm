// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package parse

// NOTE: this file is generated via 'go:generate' - manual changes will be overwritten

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

type LogProfileId struct {
	SubscriptionId string
	Name           string
}

func NewLogProfileID(subscriptionId, name string) LogProfileId {
	return LogProfileId{
		SubscriptionId: subscriptionId,
		Name:           name,
	}
}

func (id LogProfileId) String() string {
	segments := []string{
		fmt.Sprintf("Name %q", id.Name),
	}
	segmentsStr := strings.Join(segments, " / ")
	return fmt.Sprintf("%s: (%s)", "Log Profile", segmentsStr)
}

func (id LogProfileId) ID() string {
	fmtString := "/subscriptions/%s/providers/Microsoft.Insights/logProfiles/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.Name)
}

// LogProfileID parses a LogProfile ID into an LogProfileId struct
func LogProfileID(input string) (*LogProfileId, error) {
	id, err := resourceids.ParseAzureResourceID(input)
	if err != nil {
		return nil, err
	}

	resourceId := LogProfileId{
		SubscriptionId: id.SubscriptionID,
	}

	if resourceId.SubscriptionId == "" {
		return nil, fmt.Errorf("ID was missing the 'subscriptions' element")
	}

	if resourceId.Name, err = id.PopSegment("logProfiles"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &resourceId, nil
}
