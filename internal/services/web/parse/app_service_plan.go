// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package parse

// NOTE: this file is generated via 'go:generate' - manual changes will be overwritten

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

type AppServicePlanId struct {
	SubscriptionId string
	ResourceGroup  string
	ServerFarmName string
}

func NewAppServicePlanID(subscriptionId, resourceGroup, serverFarmName string) AppServicePlanId {
	return AppServicePlanId{
		SubscriptionId: subscriptionId,
		ResourceGroup:  resourceGroup,
		ServerFarmName: serverFarmName,
	}
}

func (id AppServicePlanId) String() string {
	segments := []string{
		fmt.Sprintf("Server Farm Name %q", id.ServerFarmName),
		fmt.Sprintf("Resource Group %q", id.ResourceGroup),
	}
	segmentsStr := strings.Join(segments, " / ")
	return fmt.Sprintf("%s: (%s)", "App Service Plan", segmentsStr)
}

func (id AppServicePlanId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Web/serverFarms/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroup, id.ServerFarmName)
}

// AppServicePlanID parses a AppServicePlan ID into an AppServicePlanId struct
func AppServicePlanID(input string) (*AppServicePlanId, error) {
	id, err := resourceids.ParseAzureResourceID(input)
	if err != nil {
		return nil, fmt.Errorf("parsing %q as an AppServicePlan ID: %+v", input, err)
	}

	resourceId := AppServicePlanId{
		SubscriptionId: id.SubscriptionID,
		ResourceGroup:  id.ResourceGroup,
	}

	if resourceId.SubscriptionId == "" {
		return nil, fmt.Errorf("ID was missing the 'subscriptions' element")
	}

	if resourceId.ResourceGroup == "" {
		return nil, fmt.Errorf("ID was missing the 'resourceGroups' element")
	}

	if resourceId.ServerFarmName, err = id.PopSegment("serverFarms"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &resourceId, nil
}
