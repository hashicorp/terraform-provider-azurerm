// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package parse

// NOTE: this file is generated via 'go:generate' - manual changes will be overwritten

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

type ResourceGroupCostManagementExportId struct {
	SubscriptionId string
	ResourceGroup  string
	ExportName     string
}

func NewResourceGroupCostManagementExportID(subscriptionId, resourceGroup, exportName string) ResourceGroupCostManagementExportId {
	return ResourceGroupCostManagementExportId{
		SubscriptionId: subscriptionId,
		ResourceGroup:  resourceGroup,
		ExportName:     exportName,
	}
}

func (id ResourceGroupCostManagementExportId) String() string {
	segments := []string{
		fmt.Sprintf("Export Name %q", id.ExportName),
		fmt.Sprintf("Resource Group %q", id.ResourceGroup),
	}
	segmentsStr := strings.Join(segments, " / ")
	return fmt.Sprintf("%s: (%s)", "Resource Group Cost Management Export", segmentsStr)
}

func (id ResourceGroupCostManagementExportId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.CostManagement/exports/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroup, id.ExportName)
}

// ResourceGroupCostManagementExportID parses a ResourceGroupCostManagementExport ID into an ResourceGroupCostManagementExportId struct
func ResourceGroupCostManagementExportID(input string) (*ResourceGroupCostManagementExportId, error) {
	id, err := resourceids.ParseAzureResourceID(input)
	if err != nil {
		return nil, fmt.Errorf("parsing %q as an ResourceGroupCostManagementExport ID: %+v", input, err)
	}

	resourceId := ResourceGroupCostManagementExportId{
		SubscriptionId: id.SubscriptionID,
		ResourceGroup:  id.ResourceGroup,
	}

	if resourceId.SubscriptionId == "" {
		return nil, fmt.Errorf("ID was missing the 'subscriptions' element")
	}

	if resourceId.ResourceGroup == "" {
		return nil, fmt.Errorf("ID was missing the 'resourceGroups' element")
	}

	if resourceId.ExportName, err = id.PopSegment("exports"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &resourceId, nil
}
