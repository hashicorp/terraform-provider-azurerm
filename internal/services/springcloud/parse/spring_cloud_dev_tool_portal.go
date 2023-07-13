// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package parse

// NOTE: this file is generated via 'go:generate' - manual changes will be overwritten

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

type SpringCloudDevToolPortalId struct {
	SubscriptionId    string
	ResourceGroup     string
	SpringName        string
	DevToolPortalName string
}

func NewSpringCloudDevToolPortalID(subscriptionId, resourceGroup, springName, devToolPortalName string) SpringCloudDevToolPortalId {
	return SpringCloudDevToolPortalId{
		SubscriptionId:    subscriptionId,
		ResourceGroup:     resourceGroup,
		SpringName:        springName,
		DevToolPortalName: devToolPortalName,
	}
}

func (id SpringCloudDevToolPortalId) String() string {
	segments := []string{
		fmt.Sprintf("Dev Tool Portal Name %q", id.DevToolPortalName),
		fmt.Sprintf("Spring Name %q", id.SpringName),
		fmt.Sprintf("Resource Group %q", id.ResourceGroup),
	}
	segmentsStr := strings.Join(segments, " / ")
	return fmt.Sprintf("%s: (%s)", "Spring Cloud Dev Tool Portal", segmentsStr)
}

func (id SpringCloudDevToolPortalId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.AppPlatform/Spring/%s/DevToolPortals/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroup, id.SpringName, id.DevToolPortalName)
}

// SpringCloudDevToolPortalID parses a SpringCloudDevToolPortal ID into an SpringCloudDevToolPortalId struct
func SpringCloudDevToolPortalID(input string) (*SpringCloudDevToolPortalId, error) {
	id, err := resourceids.ParseAzureResourceID(input)
	if err != nil {
		return nil, fmt.Errorf("parsing %q as an SpringCloudDevToolPortal ID: %+v", input, err)
	}

	resourceId := SpringCloudDevToolPortalId{
		SubscriptionId: id.SubscriptionID,
		ResourceGroup:  id.ResourceGroup,
	}

	if resourceId.SubscriptionId == "" {
		return nil, fmt.Errorf("ID was missing the 'subscriptions' element")
	}

	if resourceId.ResourceGroup == "" {
		return nil, fmt.Errorf("ID was missing the 'resourceGroups' element")
	}

	if resourceId.SpringName, err = id.PopSegment("Spring"); err != nil {
		return nil, err
	}
	if resourceId.DevToolPortalName, err = id.PopSegment("DevToolPortals"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &resourceId, nil
}
