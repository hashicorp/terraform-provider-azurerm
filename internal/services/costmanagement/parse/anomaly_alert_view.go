// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package parse

// NOTE: this file is generated via 'go:generate' - manual changes will be overwritten

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

type AnomalyAlertViewId struct {
	SubscriptionId string
	ViewName       string
}

func NewAnomalyAlertViewID(subscriptionId, viewName string) AnomalyAlertViewId {
	return AnomalyAlertViewId{
		SubscriptionId: subscriptionId,
		ViewName:       viewName,
	}
}

func (id AnomalyAlertViewId) String() string {
	segments := []string{
		fmt.Sprintf("View Name %q", id.ViewName),
	}
	segmentsStr := strings.Join(segments, " / ")
	return fmt.Sprintf("%s: (%s)", "Anomaly Alert View", segmentsStr)
}

func (id AnomalyAlertViewId) ID() string {
	fmtString := "/subscriptions/%s/providers/Microsoft.CostManagement/views/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ViewName)
}

// AnomalyAlertViewID parses a AnomalyAlertView ID into an AnomalyAlertViewId struct
func AnomalyAlertViewID(input string) (*AnomalyAlertViewId, error) {
	id, err := resourceids.ParseAzureResourceID(input)
	if err != nil {
		return nil, fmt.Errorf("parsing %q as an AnomalyAlertView ID: %+v", input, err)
	}

	resourceId := AnomalyAlertViewId{
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
