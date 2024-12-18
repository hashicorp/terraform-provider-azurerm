// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package parse

// NOTE: this file is generated via 'go:generate' - manual changes will be overwritten

import (
	"errors"
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

type LoadBalancerInboundNatPoolId struct {
	SubscriptionId     string
	ResourceGroup      string
	LoadBalancerName   string
	InboundNatPoolName string
}

func NewLoadBalancerInboundNatPoolID(subscriptionId, resourceGroup, loadBalancerName, inboundNatPoolName string) LoadBalancerInboundNatPoolId {
	return LoadBalancerInboundNatPoolId{
		SubscriptionId:     subscriptionId,
		ResourceGroup:      resourceGroup,
		LoadBalancerName:   loadBalancerName,
		InboundNatPoolName: inboundNatPoolName,
	}
}

func (id LoadBalancerInboundNatPoolId) String() string {
	segments := []string{
		fmt.Sprintf("Inbound Nat Pool Name %q", id.InboundNatPoolName),
		fmt.Sprintf("Load Balancer Name %q", id.LoadBalancerName),
		fmt.Sprintf("Resource Group %q", id.ResourceGroup),
	}
	segmentsStr := strings.Join(segments, " / ")
	return fmt.Sprintf("%s: (%s)", "Load Balancer Inbound Nat Pool", segmentsStr)
}

func (id LoadBalancerInboundNatPoolId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Network/loadBalancers/%s/inboundNatPools/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroup, id.LoadBalancerName, id.InboundNatPoolName)
}

// LoadBalancerInboundNatPoolID parses a LoadBalancerInboundNatPool ID into an LoadBalancerInboundNatPoolId struct
func LoadBalancerInboundNatPoolID(input string) (*LoadBalancerInboundNatPoolId, error) {
	id, err := resourceids.ParseAzureResourceID(input)
	if err != nil {
		return nil, fmt.Errorf("parsing %q as an LoadBalancerInboundNatPool ID: %+v", input, err)
	}

	resourceId := LoadBalancerInboundNatPoolId{
		SubscriptionId: id.SubscriptionID,
		ResourceGroup:  id.ResourceGroup,
	}

	if resourceId.SubscriptionId == "" {
		return nil, errors.New("ID was missing the 'subscriptions' element")
	}

	if resourceId.ResourceGroup == "" {
		return nil, errors.New("ID was missing the 'resourceGroups' element")
	}

	if resourceId.LoadBalancerName, err = id.PopSegment("loadBalancers"); err != nil {
		return nil, err
	}
	if resourceId.InboundNatPoolName, err = id.PopSegment("inboundNatPools"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &resourceId, nil
}
