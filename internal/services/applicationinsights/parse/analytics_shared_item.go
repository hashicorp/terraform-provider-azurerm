// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package parse

// NOTE: this file is generated via 'go:generate' - manual changes will be overwritten

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

type AnalyticsSharedItemId struct {
	SubscriptionId    string
	ResourceGroup     string
	ComponentName     string
	AnalyticsItemName string
}

func NewAnalyticsSharedItemID(subscriptionId, resourceGroup, componentName, analyticsItemName string) AnalyticsSharedItemId {
	return AnalyticsSharedItemId{
		SubscriptionId:    subscriptionId,
		ResourceGroup:     resourceGroup,
		ComponentName:     componentName,
		AnalyticsItemName: analyticsItemName,
	}
}

func (id AnalyticsSharedItemId) String() string {
	segments := []string{
		fmt.Sprintf("Analytics Item Name %q", id.AnalyticsItemName),
		fmt.Sprintf("Component Name %q", id.ComponentName),
		fmt.Sprintf("Resource Group %q", id.ResourceGroup),
	}
	segmentsStr := strings.Join(segments, " / ")
	return fmt.Sprintf("%s: (%s)", "Analytics Shared Item", segmentsStr)
}

func (id AnalyticsSharedItemId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Insights/components/%s/analyticsItems/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroup, id.ComponentName, id.AnalyticsItemName)
}

// AnalyticsSharedItemID parses a AnalyticsSharedItem ID into an AnalyticsSharedItemId struct
func AnalyticsSharedItemID(input string) (*AnalyticsSharedItemId, error) {
	id, err := resourceids.ParseAzureResourceID(input)
	if err != nil {
		return nil, fmt.Errorf("parsing %q as an AnalyticsSharedItem ID: %+v", input, err)
	}

	resourceId := AnalyticsSharedItemId{
		SubscriptionId: id.SubscriptionID,
		ResourceGroup:  id.ResourceGroup,
	}

	if resourceId.SubscriptionId == "" {
		return nil, fmt.Errorf("ID was missing the 'subscriptions' element")
	}

	if resourceId.ResourceGroup == "" {
		return nil, fmt.Errorf("ID was missing the 'resourceGroups' element")
	}

	if resourceId.ComponentName, err = id.PopSegment("components"); err != nil {
		return nil, err
	}
	if resourceId.AnalyticsItemName, err = id.PopSegment("analyticsItems"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &resourceId, nil
}
