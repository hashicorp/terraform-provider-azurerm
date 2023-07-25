// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package parse

// NOTE: this file is generated via 'go:generate' - manual changes will be overwritten

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

type LoadBalancerBackendAddressPoolId struct {
	SubscriptionId         string
	ResourceGroup          string
	LoadBalancerName       string
	BackendAddressPoolName string
}

func NewLoadBalancerBackendAddressPoolID(subscriptionId, resourceGroup, loadBalancerName, backendAddressPoolName string) LoadBalancerBackendAddressPoolId {
	return LoadBalancerBackendAddressPoolId{
		SubscriptionId:         subscriptionId,
		ResourceGroup:          resourceGroup,
		LoadBalancerName:       loadBalancerName,
		BackendAddressPoolName: backendAddressPoolName,
	}
}

func (id LoadBalancerBackendAddressPoolId) String() string {
	segments := []string{
		fmt.Sprintf("Backend Address Pool Name %q", id.BackendAddressPoolName),
		fmt.Sprintf("Load Balancer Name %q", id.LoadBalancerName),
		fmt.Sprintf("Resource Group %q", id.ResourceGroup),
	}
	segmentsStr := strings.Join(segments, " / ")
	return fmt.Sprintf("%s: (%s)", "Load Balancer Backend Address Pool", segmentsStr)
}

func (id LoadBalancerBackendAddressPoolId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Network/loadBalancers/%s/backendAddressPools/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroup, id.LoadBalancerName, id.BackendAddressPoolName)
}

// LoadBalancerBackendAddressPoolID parses a LoadBalancerBackendAddressPool ID into an LoadBalancerBackendAddressPoolId struct
func LoadBalancerBackendAddressPoolID(input string) (*LoadBalancerBackendAddressPoolId, error) {
	id, err := resourceids.ParseAzureResourceID(input)
	if err != nil {
		return nil, fmt.Errorf("parsing %q as an LoadBalancerBackendAddressPool ID: %+v", input, err)
	}

	resourceId := LoadBalancerBackendAddressPoolId{
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
	if resourceId.BackendAddressPoolName, err = id.PopSegment("backendAddressPools"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &resourceId, nil
}

// LoadBalancerBackendAddressPoolIDInsensitively parses an LoadBalancerBackendAddressPool ID into an LoadBalancerBackendAddressPoolId struct, insensitively
// This should only be used to parse an ID for rewriting, the LoadBalancerBackendAddressPoolID
// method should be used instead for validation etc.
//
// Whilst this may seem strange, this enables Terraform have consistent casing
// which works around issues in Core, whilst handling broken API responses.
func LoadBalancerBackendAddressPoolIDInsensitively(input string) (*LoadBalancerBackendAddressPoolId, error) {
	id, err := resourceids.ParseAzureResourceID(input)
	if err != nil {
		return nil, err
	}

	resourceId := LoadBalancerBackendAddressPoolId{
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

	// find the correct casing for the 'backendAddressPools' segment
	backendAddressPoolsKey := "backendAddressPools"
	for key := range id.Path {
		if strings.EqualFold(key, backendAddressPoolsKey) {
			backendAddressPoolsKey = key
			break
		}
	}
	if resourceId.BackendAddressPoolName, err = id.PopSegment(backendAddressPoolsKey); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &resourceId, nil
}
