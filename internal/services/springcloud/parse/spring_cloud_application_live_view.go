// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package parse

// NOTE: this file is generated via 'go:generate' - manual changes will be overwritten

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

type SpringCloudApplicationLiveViewId struct {
	SubscriptionId          string
	ResourceGroup           string
	SpringName              string
	ApplicationLiveViewName string
}

func NewSpringCloudApplicationLiveViewID(subscriptionId, resourceGroup, springName, applicationLiveViewName string) SpringCloudApplicationLiveViewId {
	return SpringCloudApplicationLiveViewId{
		SubscriptionId:          subscriptionId,
		ResourceGroup:           resourceGroup,
		SpringName:              springName,
		ApplicationLiveViewName: applicationLiveViewName,
	}
}

func (id SpringCloudApplicationLiveViewId) String() string {
	segments := []string{
		fmt.Sprintf("Application Live View Name %q", id.ApplicationLiveViewName),
		fmt.Sprintf("Spring Name %q", id.SpringName),
		fmt.Sprintf("Resource Group %q", id.ResourceGroup),
	}
	segmentsStr := strings.Join(segments, " / ")
	return fmt.Sprintf("%s: (%s)", "Spring Cloud Application Live View", segmentsStr)
}

func (id SpringCloudApplicationLiveViewId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.AppPlatform/spring/%s/applicationLiveViews/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroup, id.SpringName, id.ApplicationLiveViewName)
}

// SpringCloudApplicationLiveViewID parses a SpringCloudApplicationLiveView ID into an SpringCloudApplicationLiveViewId struct
func SpringCloudApplicationLiveViewID(input string) (*SpringCloudApplicationLiveViewId, error) {
	id, err := resourceids.ParseAzureResourceID(input)
	if err != nil {
		return nil, fmt.Errorf("parsing %q as an SpringCloudApplicationLiveView ID: %+v", input, err)
	}

	resourceId := SpringCloudApplicationLiveViewId{
		SubscriptionId: id.SubscriptionID,
		ResourceGroup:  id.ResourceGroup,
	}

	if resourceId.SubscriptionId == "" {
		return nil, fmt.Errorf("ID was missing the 'subscriptions' element")
	}

	if resourceId.ResourceGroup == "" {
		return nil, fmt.Errorf("ID was missing the 'resourceGroups' element")
	}

	if resourceId.SpringName, err = id.PopSegment("spring"); err != nil {
		return nil, err
	}
	if resourceId.ApplicationLiveViewName, err = id.PopSegment("applicationLiveViews"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &resourceId, nil
}
