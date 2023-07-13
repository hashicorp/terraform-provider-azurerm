// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package parse

// NOTE: this file is generated via 'go:generate' - manual changes will be overwritten

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

type StaticSiteCustomDomainId struct {
	SubscriptionId   string
	ResourceGroup    string
	StaticSiteName   string
	CustomDomainName string
}

func NewStaticSiteCustomDomainID(subscriptionId, resourceGroup, staticSiteName, customDomainName string) StaticSiteCustomDomainId {
	return StaticSiteCustomDomainId{
		SubscriptionId:   subscriptionId,
		ResourceGroup:    resourceGroup,
		StaticSiteName:   staticSiteName,
		CustomDomainName: customDomainName,
	}
}

func (id StaticSiteCustomDomainId) String() string {
	segments := []string{
		fmt.Sprintf("Custom Domain Name %q", id.CustomDomainName),
		fmt.Sprintf("Static Site Name %q", id.StaticSiteName),
		fmt.Sprintf("Resource Group %q", id.ResourceGroup),
	}
	segmentsStr := strings.Join(segments, " / ")
	return fmt.Sprintf("%s: (%s)", "Static Site Custom Domain", segmentsStr)
}

func (id StaticSiteCustomDomainId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Web/staticSites/%s/customDomains/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroup, id.StaticSiteName, id.CustomDomainName)
}

// StaticSiteCustomDomainID parses a StaticSiteCustomDomain ID into an StaticSiteCustomDomainId struct
func StaticSiteCustomDomainID(input string) (*StaticSiteCustomDomainId, error) {
	id, err := resourceids.ParseAzureResourceID(input)
	if err != nil {
		return nil, fmt.Errorf("parsing %q as an StaticSiteCustomDomain ID: %+v", input, err)
	}

	resourceId := StaticSiteCustomDomainId{
		SubscriptionId: id.SubscriptionID,
		ResourceGroup:  id.ResourceGroup,
	}

	if resourceId.SubscriptionId == "" {
		return nil, fmt.Errorf("ID was missing the 'subscriptions' element")
	}

	if resourceId.ResourceGroup == "" {
		return nil, fmt.Errorf("ID was missing the 'resourceGroups' element")
	}

	if resourceId.StaticSiteName, err = id.PopSegment("staticSites"); err != nil {
		return nil, err
	}
	if resourceId.CustomDomainName, err = id.PopSegment("customDomains"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &resourceId, nil
}
