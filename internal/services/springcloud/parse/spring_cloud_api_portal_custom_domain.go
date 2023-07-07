// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package parse

// NOTE: this file is generated via 'go:generate' - manual changes will be overwritten

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

type SpringCloudAPIPortalCustomDomainId struct {
	SubscriptionId string
	ResourceGroup  string
	SpringName     string
	ApiPortalName  string
	DomainName     string
}

func NewSpringCloudAPIPortalCustomDomainID(subscriptionId, resourceGroup, springName, apiPortalName, domainName string) SpringCloudAPIPortalCustomDomainId {
	return SpringCloudAPIPortalCustomDomainId{
		SubscriptionId: subscriptionId,
		ResourceGroup:  resourceGroup,
		SpringName:     springName,
		ApiPortalName:  apiPortalName,
		DomainName:     domainName,
	}
}

func (id SpringCloudAPIPortalCustomDomainId) String() string {
	segments := []string{
		fmt.Sprintf("Domain Name %q", id.DomainName),
		fmt.Sprintf("Api Portal Name %q", id.ApiPortalName),
		fmt.Sprintf("Spring Name %q", id.SpringName),
		fmt.Sprintf("Resource Group %q", id.ResourceGroup),
	}
	segmentsStr := strings.Join(segments, " / ")
	return fmt.Sprintf("%s: (%s)", "Spring Cloud A P I Portal Custom Domain", segmentsStr)
}

func (id SpringCloudAPIPortalCustomDomainId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.AppPlatform/spring/%s/apiPortals/%s/domains/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroup, id.SpringName, id.ApiPortalName, id.DomainName)
}

// SpringCloudAPIPortalCustomDomainID parses a SpringCloudAPIPortalCustomDomain ID into an SpringCloudAPIPortalCustomDomainId struct
func SpringCloudAPIPortalCustomDomainID(input string) (*SpringCloudAPIPortalCustomDomainId, error) {
	id, err := resourceids.ParseAzureResourceID(input)
	if err != nil {
		return nil, fmt.Errorf("parsing %q as an SpringCloudAPIPortalCustomDomain ID: %+v", input, err)
	}

	resourceId := SpringCloudAPIPortalCustomDomainId{
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
	if resourceId.ApiPortalName, err = id.PopSegment("apiPortals"); err != nil {
		return nil, err
	}
	if resourceId.DomainName, err = id.PopSegment("domains"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &resourceId, nil
}

// SpringCloudAPIPortalCustomDomainIDInsensitively parses an SpringCloudAPIPortalCustomDomain ID into an SpringCloudAPIPortalCustomDomainId struct, insensitively
// This should only be used to parse an ID for rewriting, the SpringCloudAPIPortalCustomDomainID
// method should be used instead for validation etc.
//
// Whilst this may seem strange, this enables Terraform have consistent casing
// which works around issues in Core, whilst handling broken API responses.
func SpringCloudAPIPortalCustomDomainIDInsensitively(input string) (*SpringCloudAPIPortalCustomDomainId, error) {
	id, err := resourceids.ParseAzureResourceID(input)
	if err != nil {
		return nil, err
	}

	resourceId := SpringCloudAPIPortalCustomDomainId{
		SubscriptionId: id.SubscriptionID,
		ResourceGroup:  id.ResourceGroup,
	}

	if resourceId.SubscriptionId == "" {
		return nil, fmt.Errorf("ID was missing the 'subscriptions' element")
	}

	if resourceId.ResourceGroup == "" {
		return nil, fmt.Errorf("ID was missing the 'resourceGroups' element")
	}

	// find the correct casing for the 'spring' segment
	springKey := "spring"
	for key := range id.Path {
		if strings.EqualFold(key, springKey) {
			springKey = key
			break
		}
	}
	if resourceId.SpringName, err = id.PopSegment(springKey); err != nil {
		return nil, err
	}

	// find the correct casing for the 'apiPortals' segment
	apiPortalsKey := "apiPortals"
	for key := range id.Path {
		if strings.EqualFold(key, apiPortalsKey) {
			apiPortalsKey = key
			break
		}
	}
	if resourceId.ApiPortalName, err = id.PopSegment(apiPortalsKey); err != nil {
		return nil, err
	}

	// find the correct casing for the 'domains' segment
	domainsKey := "domains"
	for key := range id.Path {
		if strings.EqualFold(key, domainsKey) {
			domainsKey = key
			break
		}
	}
	if resourceId.DomainName, err = id.PopSegment(domainsKey); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &resourceId, nil
}
