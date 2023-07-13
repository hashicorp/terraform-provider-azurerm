// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package parse

// NOTE: this file is generated via 'go:generate' - manual changes will be overwritten

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

type ResourceGroupCostManagementViewId struct {
	SubscriptionId string
	ResourceGroup  string
	ViewName       string
}

func NewResourceGroupCostManagementViewID(subscriptionId, resourceGroup, viewName string) ResourceGroupCostManagementViewId {
	return ResourceGroupCostManagementViewId{
		SubscriptionId: subscriptionId,
		ResourceGroup:  resourceGroup,
		ViewName:       viewName,
	}
}

func (id ResourceGroupCostManagementViewId) String() string {
	segments := []string{
		fmt.Sprintf("View Name %q", id.ViewName),
		fmt.Sprintf("Resource Group %q", id.ResourceGroup),
	}
	segmentsStr := strings.Join(segments, " / ")
	return fmt.Sprintf("%s: (%s)", "Resource Group Cost Management View", segmentsStr)
}

func (id ResourceGroupCostManagementViewId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.CostManagement/views/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroup, id.ViewName)
}

// ResourceGroupCostManagementViewID parses a ResourceGroupCostManagementView ID into an ResourceGroupCostManagementViewId struct
func ResourceGroupCostManagementViewID(input string) (*ResourceGroupCostManagementViewId, error) {
	id, err := resourceids.ParseAzureResourceID(input)
	if err != nil {
		return nil, fmt.Errorf("parsing %q as an ResourceGroupCostManagementView ID: %+v", input, err)
	}

	resourceId := ResourceGroupCostManagementViewId{
		SubscriptionId: id.SubscriptionID,
		ResourceGroup:  id.ResourceGroup,
	}

	if resourceId.SubscriptionId == "" {
		return nil, fmt.Errorf("ID was missing the 'subscriptions' element")
	}

	if resourceId.ResourceGroup == "" {
		return nil, fmt.Errorf("ID was missing the 'resourceGroups' element")
	}

	if resourceId.ViewName, err = id.PopSegment("views"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &resourceId, nil
}
