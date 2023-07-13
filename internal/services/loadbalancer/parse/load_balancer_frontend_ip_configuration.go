// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package parse

// NOTE: this file is generated via 'go:generate' - manual changes will be overwritten

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

type LoadBalancerFrontendIpConfigurationId struct {
	SubscriptionId              string
	ResourceGroup               string
	LoadBalancerName            string
	FrontendIPConfigurationName string
}

func NewLoadBalancerFrontendIpConfigurationID(subscriptionId, resourceGroup, loadBalancerName, frontendIPConfigurationName string) LoadBalancerFrontendIpConfigurationId {
	return LoadBalancerFrontendIpConfigurationId{
		SubscriptionId:              subscriptionId,
		ResourceGroup:               resourceGroup,
		LoadBalancerName:            loadBalancerName,
		FrontendIPConfigurationName: frontendIPConfigurationName,
	}
}

func (id LoadBalancerFrontendIpConfigurationId) String() string {
	segments := []string{
		fmt.Sprintf("Frontend I P Configuration Name %q", id.FrontendIPConfigurationName),
		fmt.Sprintf("Load Balancer Name %q", id.LoadBalancerName),
		fmt.Sprintf("Resource Group %q", id.ResourceGroup),
	}
	segmentsStr := strings.Join(segments, " / ")
	return fmt.Sprintf("%s: (%s)", "Load Balancer Frontend Ip Configuration", segmentsStr)
}

func (id LoadBalancerFrontendIpConfigurationId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Network/loadBalancers/%s/frontendIPConfigurations/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroup, id.LoadBalancerName, id.FrontendIPConfigurationName)
}

// LoadBalancerFrontendIpConfigurationID parses a LoadBalancerFrontendIpConfiguration ID into an LoadBalancerFrontendIpConfigurationId struct
func LoadBalancerFrontendIpConfigurationID(input string) (*LoadBalancerFrontendIpConfigurationId, error) {
	id, err := resourceids.ParseAzureResourceID(input)
	if err != nil {
		return nil, fmt.Errorf("parsing %q as an LoadBalancerFrontendIpConfiguration ID: %+v", input, err)
	}

	resourceId := LoadBalancerFrontendIpConfigurationId{
		SubscriptionId: id.SubscriptionID,
		ResourceGroup:  id.ResourceGroup,
	}

	if resourceId.SubscriptionId == "" {
		return nil, fmt.Errorf("ID was missing the 'subscriptions' element")
	}

	if resourceId.ResourceGroup == "" {
		return nil, fmt.Errorf("ID was missing the 'resourceGroups' element")
	}

	if resourceId.LoadBalancerName, err = id.PopSegment("loadBalancers"); err != nil {
		return nil, err
	}
	if resourceId.FrontendIPConfigurationName, err = id.PopSegment("frontendIPConfigurations"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &resourceId, nil
}

// LoadBalancerFrontendIpConfigurationIDInsensitively parses an LoadBalancerFrontendIpConfiguration ID into an LoadBalancerFrontendIpConfigurationId struct, insensitively
// This should only be used to parse an ID for rewriting, the LoadBalancerFrontendIpConfigurationID
// method should be used instead for validation etc.
//
// Whilst this may seem strange, this enables Terraform have consistent casing
// which works around issues in Core, whilst handling broken API responses.
func LoadBalancerFrontendIpConfigurationIDInsensitively(input string) (*LoadBalancerFrontendIpConfigurationId, error) {
	id, err := resourceids.ParseAzureResourceID(input)
	if err != nil {
		return nil, err
	}

	resourceId := LoadBalancerFrontendIpConfigurationId{
		SubscriptionId: id.SubscriptionID,
		ResourceGroup:  id.ResourceGroup,
	}

	if resourceId.SubscriptionId == "" {
		return nil, fmt.Errorf("ID was missing the 'subscriptions' element")
	}

	if resourceId.ResourceGroup == "" {
		return nil, fmt.Errorf("ID was missing the 'resourceGroups' element")
	}

	// find the correct casing for the 'loadBalancers' segment
	loadBalancersKey := "loadBalancers"
	for key := range id.Path {
		if strings.EqualFold(key, loadBalancersKey) {
			loadBalancersKey = key
			break
		}
	}
	if resourceId.LoadBalancerName, err = id.PopSegment(loadBalancersKey); err != nil {
		return nil, err
	}

	// find the correct casing for the 'frontendIPConfigurations' segment
	frontendIPConfigurationsKey := "frontendIPConfigurations"
	for key := range id.Path {
		if strings.EqualFold(key, frontendIPConfigurationsKey) {
			frontendIPConfigurationsKey = key
			break
		}
	}
	if resourceId.FrontendIPConfigurationName, err = id.PopSegment(frontendIPConfigurationsKey); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &resourceId, nil
}
