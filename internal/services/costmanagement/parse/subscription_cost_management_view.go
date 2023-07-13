// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package parse

// NOTE: this file is generated via 'go:generate' - manual changes will be overwritten

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

type SubscriptionCostManagementViewId struct {
	SubscriptionId string
	ViewName       string
}

func NewSubscriptionCostManagementViewID(subscriptionId, viewName string) SubscriptionCostManagementViewId {
	return SubscriptionCostManagementViewId{
		SubscriptionId: subscriptionId,
		ViewName:       viewName,
	}
}

func (id SubscriptionCostManagementViewId) String() string {
	segments := []string{
		fmt.Sprintf("View Name %q", id.ViewName),
	}
	segmentsStr := strings.Join(segments, " / ")
	return fmt.Sprintf("%s: (%s)", "Subscription Cost Management View", segmentsStr)
}

func (id SubscriptionCostManagementViewId) ID() string {
	fmtString := "/subscriptions/%s/providers/Microsoft.CostManagement/views/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ViewName)
}

// SubscriptionCostManagementViewID parses a SubscriptionCostManagementView ID into an SubscriptionCostManagementViewId struct
func SubscriptionCostManagementViewID(input string) (*SubscriptionCostManagementViewId, error) {
	id, err := resourceids.ParseAzureResourceID(input)
	if err != nil {
		return nil, fmt.Errorf("parsing %q as an SubscriptionCostManagementView ID: %+v", input, err)
	}

	resourceId := SubscriptionCostManagementViewId{
		SubscriptionId: id.SubscriptionID,
	}

	if resourceId.SubscriptionId == "" {
		return nil, fmt.Errorf("ID was missing the 'subscriptions' element")
	}

	if resourceId.ViewName, err = id.PopSegment("views"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &resourceId, nil
}
