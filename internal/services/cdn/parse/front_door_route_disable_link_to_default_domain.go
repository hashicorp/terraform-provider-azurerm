// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package parse

// NOTE: this file is generated via 'go:generate' - manual changes will be overwritten

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

type FrontDoorRouteDisableLinkToDefaultDomainId struct {
	SubscriptionId                 string
	ResourceGroup                  string
	ProfileName                    string
	AfdEndpointName                string
	RouteName                      string
	DisableLinkToDefaultDomainName string
}

func NewFrontDoorRouteDisableLinkToDefaultDomainID(subscriptionId, resourceGroup, profileName, afdEndpointName, routeName, disableLinkToDefaultDomainName string) FrontDoorRouteDisableLinkToDefaultDomainId {
	return FrontDoorRouteDisableLinkToDefaultDomainId{
		SubscriptionId:                 subscriptionId,
		ResourceGroup:                  resourceGroup,
		ProfileName:                    profileName,
		AfdEndpointName:                afdEndpointName,
		RouteName:                      routeName,
		DisableLinkToDefaultDomainName: disableLinkToDefaultDomainName,
	}
}

func (id FrontDoorRouteDisableLinkToDefaultDomainId) String() string {
	segments := []string{
		fmt.Sprintf("Disable Link To Default Domain Name %q", id.DisableLinkToDefaultDomainName),
		fmt.Sprintf("Route Name %q", id.RouteName),
		fmt.Sprintf("Afd Endpoint Name %q", id.AfdEndpointName),
		fmt.Sprintf("Profile Name %q", id.ProfileName),
		fmt.Sprintf("Resource Group %q", id.ResourceGroup),
	}
	segmentsStr := strings.Join(segments, " / ")
	return fmt.Sprintf("%s: (%s)", "Front Door Route Disable Link To Default Domain", segmentsStr)
}

func (id FrontDoorRouteDisableLinkToDefaultDomainId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Cdn/profiles/%s/afdEndpoints/%s/routes/%s/disableLinkToDefaultDomain/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroup, id.ProfileName, id.AfdEndpointName, id.RouteName, id.DisableLinkToDefaultDomainName)
}

// FrontDoorRouteDisableLinkToDefaultDomainID parses a FrontDoorRouteDisableLinkToDefaultDomain ID into an FrontDoorRouteDisableLinkToDefaultDomainId struct
func FrontDoorRouteDisableLinkToDefaultDomainID(input string) (*FrontDoorRouteDisableLinkToDefaultDomainId, error) {
	id, err := resourceids.ParseAzureResourceID(input)
	if err != nil {
		return nil, fmt.Errorf("parsing %q as an FrontDoorRouteDisableLinkToDefaultDomain ID: %+v", input, err)
	}

	resourceId := FrontDoorRouteDisableLinkToDefaultDomainId{
		SubscriptionId: id.SubscriptionID,
		ResourceGroup:  id.ResourceGroup,
	}

	if resourceId.SubscriptionId == "" {
		return nil, fmt.Errorf("ID was missing the 'subscriptions' element")
	}

	if resourceId.ResourceGroup == "" {
		return nil, fmt.Errorf("ID was missing the 'resourceGroups' element")
	}

	if resourceId.ProfileName, err = id.PopSegment("profiles"); err != nil {
		return nil, err
	}
	if resourceId.AfdEndpointName, err = id.PopSegment("afdEndpoints"); err != nil {
		return nil, err
	}
	if resourceId.RouteName, err = id.PopSegment("routes"); err != nil {
		return nil, err
	}
	if resourceId.DisableLinkToDefaultDomainName, err = id.PopSegment("disableLinkToDefaultDomain"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &resourceId, nil
}
