// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package parse

// NOTE: this file is generated via 'go:generate' - manual changes will be overwritten

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

type SpringCloudGatewayCustomDomainId struct {
	SubscriptionId string
	ResourceGroup  string
	SpringName     string
	GatewayName    string
	DomainName     string
}

func NewSpringCloudGatewayCustomDomainID(subscriptionId, resourceGroup, springName, gatewayName, domainName string) SpringCloudGatewayCustomDomainId {
	return SpringCloudGatewayCustomDomainId{
		SubscriptionId: subscriptionId,
		ResourceGroup:  resourceGroup,
		SpringName:     springName,
		GatewayName:    gatewayName,
		DomainName:     domainName,
	}
}

func (id SpringCloudGatewayCustomDomainId) String() string {
	segments := []string{
		fmt.Sprintf("Domain Name %q", id.DomainName),
		fmt.Sprintf("Gateway Name %q", id.GatewayName),
		fmt.Sprintf("Spring Name %q", id.SpringName),
		fmt.Sprintf("Resource Group %q", id.ResourceGroup),
	}
	segmentsStr := strings.Join(segments, " / ")
	return fmt.Sprintf("%s: (%s)", "Spring Cloud Gateway Custom Domain", segmentsStr)
}

func (id SpringCloudGatewayCustomDomainId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.AppPlatform/spring/%s/gateways/%s/domains/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroup, id.SpringName, id.GatewayName, id.DomainName)
}

// SpringCloudGatewayCustomDomainID parses a SpringCloudGatewayCustomDomain ID into an SpringCloudGatewayCustomDomainId struct
func SpringCloudGatewayCustomDomainID(input string) (*SpringCloudGatewayCustomDomainId, error) {
	id, err := resourceids.ParseAzureResourceID(input)
	if err != nil {
		return nil, fmt.Errorf("parsing %q as an SpringCloudGatewayCustomDomain ID: %+v", input, err)
	}

	resourceId := SpringCloudGatewayCustomDomainId{
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
	if resourceId.GatewayName, err = id.PopSegment("gateways"); err != nil {
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

// SpringCloudGatewayCustomDomainIDInsensitively parses an SpringCloudGatewayCustomDomain ID into an SpringCloudGatewayCustomDomainId struct, insensitively
// This should only be used to parse an ID for rewriting, the SpringCloudGatewayCustomDomainID
// method should be used instead for validation etc.
//
// Whilst this may seem strange, this enables Terraform have consistent casing
// which works around issues in Core, whilst handling broken API responses.
func SpringCloudGatewayCustomDomainIDInsensitively(input string) (*SpringCloudGatewayCustomDomainId, error) {
	id, err := resourceids.ParseAzureResourceID(input)
	if err != nil {
		return nil, err
	}

	resourceId := SpringCloudGatewayCustomDomainId{
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

	// find the correct casing for the 'gateways' segment
	gatewaysKey := "gateways"
	for key := range id.Path {
		if strings.EqualFold(key, gatewaysKey) {
			gatewaysKey = key
			break
		}
	}
	if resourceId.GatewayName, err = id.PopSegment(gatewaysKey); err != nil {
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
