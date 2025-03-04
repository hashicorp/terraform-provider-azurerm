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

type LoadBalancerInboundNatRuleId struct {
	SubscriptionId     string
	ResourceGroup      string
	LoadBalancerName   string
	InboundNatRuleName string
}

func NewLoadBalancerInboundNatRuleID(subscriptionId, resourceGroup, loadBalancerName, inboundNatRuleName string) LoadBalancerInboundNatRuleId {
	return LoadBalancerInboundNatRuleId{
		SubscriptionId:     subscriptionId,
		ResourceGroup:      resourceGroup,
		LoadBalancerName:   loadBalancerName,
		InboundNatRuleName: inboundNatRuleName,
	}
}

func (id LoadBalancerInboundNatRuleId) String() string {
	segments := []string{
		fmt.Sprintf("Inbound Nat Rule Name %q", id.InboundNatRuleName),
		fmt.Sprintf("Load Balancer Name %q", id.LoadBalancerName),
		fmt.Sprintf("Resource Group %q", id.ResourceGroup),
	}
	segmentsStr := strings.Join(segments, " / ")
	return fmt.Sprintf("%s: (%s)", "Load Balancer Inbound Nat Rule", segmentsStr)
}

func (id LoadBalancerInboundNatRuleId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Network/loadBalancers/%s/inboundNatRules/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroup, id.LoadBalancerName, id.InboundNatRuleName)
}

// LoadBalancerInboundNatRuleID parses a LoadBalancerInboundNatRule ID into an LoadBalancerInboundNatRuleId struct
func LoadBalancerInboundNatRuleID(input string) (*LoadBalancerInboundNatRuleId, error) {
	id, err := resourceids.ParseAzureResourceID(input)
	if err != nil {
		return nil, fmt.Errorf("parsing %q as an LoadBalancerInboundNatRule ID: %+v", input, err)
	}

	resourceId := LoadBalancerInboundNatRuleId{
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
	if resourceId.InboundNatRuleName, err = id.PopSegment("inboundNatRules"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &resourceId, nil
}
