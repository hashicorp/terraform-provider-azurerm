// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package parse

// NOTE: this file is generated via 'go:generate' - manual changes will be overwritten

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

type AnalyticsUserItemId struct {
	SubscriptionId      string
	ResourceGroup       string
	ComponentName       string
	MyAnalyticsItemName string
}

func NewAnalyticsUserItemID(subscriptionId, resourceGroup, componentName, myAnalyticsItemName string) AnalyticsUserItemId {
	return AnalyticsUserItemId{
		SubscriptionId:      subscriptionId,
		ResourceGroup:       resourceGroup,
		ComponentName:       componentName,
		MyAnalyticsItemName: myAnalyticsItemName,
	}
}

func (id AnalyticsUserItemId) String() string {
	segments := []string{
		fmt.Sprintf("My Analytics Item Name %q", id.MyAnalyticsItemName),
		fmt.Sprintf("Component Name %q", id.ComponentName),
		fmt.Sprintf("Resource Group %q", id.ResourceGroup),
	}
	segmentsStr := strings.Join(segments, " / ")
	return fmt.Sprintf("%s: (%s)", "Analytics User Item", segmentsStr)
}

func (id AnalyticsUserItemId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Insights/components/%s/myAnalyticsItems/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroup, id.ComponentName, id.MyAnalyticsItemName)
}

// AnalyticsUserItemID parses a AnalyticsUserItem ID into an AnalyticsUserItemId struct
func AnalyticsUserItemID(input string) (*AnalyticsUserItemId, error) {
	id, err := resourceids.ParseAzureResourceID(input)
	if err != nil {
		return nil, fmt.Errorf("parsing %q as an AnalyticsUserItem ID: %+v", input, err)
	}

	resourceId := AnalyticsUserItemId{
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
	if resourceId.MyAnalyticsItemName, err = id.PopSegment("myAnalyticsItems"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &resourceId, nil
}
