// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package parse

// NOTE: this file is generated via 'go:generate' - manual changes will be overwritten

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

type SubscriptionCostManagementExportId struct {
	SubscriptionId string
	ExportName     string
}

func NewSubscriptionCostManagementExportID(subscriptionId, exportName string) SubscriptionCostManagementExportId {
	return SubscriptionCostManagementExportId{
		SubscriptionId: subscriptionId,
		ExportName:     exportName,
	}
}

func (id SubscriptionCostManagementExportId) String() string {
	segments := []string{
		fmt.Sprintf("Export Name %q", id.ExportName),
	}
	segmentsStr := strings.Join(segments, " / ")
	return fmt.Sprintf("%s: (%s)", "Subscription Cost Management Export", segmentsStr)
}

func (id SubscriptionCostManagementExportId) ID() string {
	fmtString := "/subscriptions/%s/providers/Microsoft.CostManagement/exports/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ExportName)
}

// SubscriptionCostManagementExportID parses a SubscriptionCostManagementExport ID into an SubscriptionCostManagementExportId struct
func SubscriptionCostManagementExportID(input string) (*SubscriptionCostManagementExportId, error) {
	id, err := resourceids.ParseAzureResourceID(input)
	if err != nil {
		return nil, fmt.Errorf("parsing %q as an SubscriptionCostManagementExport ID: %+v", input, err)
	}

	resourceId := SubscriptionCostManagementExportId{
		SubscriptionId: id.SubscriptionID,
	}

	if resourceId.SubscriptionId == "" {
		return nil, fmt.Errorf("ID was missing the 'subscriptions' element")
	}

	if resourceId.ExportName, err = id.PopSegment("exports"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &resourceId, nil
}
