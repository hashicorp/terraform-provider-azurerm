// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package parse

// NOTE: this file is generated via 'go:generate' - manual changes will be overwritten

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

type SpringCloudCustomDomainId struct {
	SubscriptionId string
	ResourceGroup  string
	SpringName     string
	AppName        string
	DomainName     string
}

func NewSpringCloudCustomDomainID(subscriptionId, resourceGroup, springName, appName, domainName string) SpringCloudCustomDomainId {
	return SpringCloudCustomDomainId{
		SubscriptionId: subscriptionId,
		ResourceGroup:  resourceGroup,
		SpringName:     springName,
		AppName:        appName,
		DomainName:     domainName,
	}
}

func (id SpringCloudCustomDomainId) String() string {
	segments := []string{
		fmt.Sprintf("Domain Name %q", id.DomainName),
		fmt.Sprintf("App Name %q", id.AppName),
		fmt.Sprintf("Spring Name %q", id.SpringName),
		fmt.Sprintf("Resource Group %q", id.ResourceGroup),
	}
	segmentsStr := strings.Join(segments, " / ")
	return fmt.Sprintf("%s: (%s)", "Spring Cloud Custom Domain", segmentsStr)
}

func (id SpringCloudCustomDomainId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.AppPlatform/spring/%s/apps/%s/domains/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroup, id.SpringName, id.AppName, id.DomainName)
}

// SpringCloudCustomDomainID parses a SpringCloudCustomDomain ID into an SpringCloudCustomDomainId struct
func SpringCloudCustomDomainID(input string) (*SpringCloudCustomDomainId, error) {
	id, err := resourceids.ParseAzureResourceID(input)
	if err != nil {
		return nil, fmt.Errorf("parsing %q as an SpringCloudCustomDomain ID: %+v", input, err)
	}

	resourceId := SpringCloudCustomDomainId{
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
	if resourceId.AppName, err = id.PopSegment("apps"); err != nil {
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

// SpringCloudCustomDomainIDInsensitively parses an SpringCloudCustomDomain ID into an SpringCloudCustomDomainId struct, insensitively
// This should only be used to parse an ID for rewriting, the SpringCloudCustomDomainID
// method should be used instead for validation etc.
//
// Whilst this may seem strange, this enables Terraform have consistent casing
// which works around issues in Core, whilst handling broken API responses.
func SpringCloudCustomDomainIDInsensitively(input string) (*SpringCloudCustomDomainId, error) {
	id, err := resourceids.ParseAzureResourceID(input)
	if err != nil {
		return nil, err
	}

	resourceId := SpringCloudCustomDomainId{
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

	// find the correct casing for the 'apps' segment
	appsKey := "apps"
	for key := range id.Path {
		if strings.EqualFold(key, appsKey) {
			appsKey = key
			break
		}
	}
	if resourceId.AppName, err = id.PopSegment(appsKey); err != nil {
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
