// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package parse

// NOTE: this file is generated via 'go:generate' - manual changes will be overwritten

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

type FrontDoorFirewallPolicyId struct {
	SubscriptionId                            string
	ResourceGroup                             string
	FrontDoorWebApplicationFirewallPolicyName string
}

func NewFrontDoorFirewallPolicyID(subscriptionId, resourceGroup, frontDoorWebApplicationFirewallPolicyName string) FrontDoorFirewallPolicyId {
	return FrontDoorFirewallPolicyId{
		SubscriptionId: subscriptionId,
		ResourceGroup:  resourceGroup,
		FrontDoorWebApplicationFirewallPolicyName: frontDoorWebApplicationFirewallPolicyName,
	}
}

func (id FrontDoorFirewallPolicyId) String() string {
	segments := []string{
		fmt.Sprintf("Front Door Web Application Firewall Policy Name %q", id.FrontDoorWebApplicationFirewallPolicyName),
		fmt.Sprintf("Resource Group %q", id.ResourceGroup),
	}
	segmentsStr := strings.Join(segments, " / ")
	return fmt.Sprintf("%s: (%s)", "Front Door Firewall Policy", segmentsStr)
}

func (id FrontDoorFirewallPolicyId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Network/frontDoorWebApplicationFirewallPolicies/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroup, id.FrontDoorWebApplicationFirewallPolicyName)
}

// FrontDoorFirewallPolicyID parses a FrontDoorFirewallPolicy ID into an FrontDoorFirewallPolicyId struct
func FrontDoorFirewallPolicyID(input string) (*FrontDoorFirewallPolicyId, error) {
	id, err := resourceids.ParseAzureResourceID(input)
	if err != nil {
		return nil, fmt.Errorf("parsing %q as an FrontDoorFirewallPolicy ID: %+v", input, err)
	}

	resourceId := FrontDoorFirewallPolicyId{
		SubscriptionId: id.SubscriptionID,
		ResourceGroup:  id.ResourceGroup,
	}

	if resourceId.SubscriptionId == "" {
		return nil, fmt.Errorf("ID was missing the 'subscriptions' element")
	}

	if resourceId.ResourceGroup == "" {
		return nil, fmt.Errorf("ID was missing the 'resourceGroups' element")
	}

	if resourceId.FrontDoorWebApplicationFirewallPolicyName, err = id.PopSegment("frontDoorWebApplicationFirewallPolicies"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &resourceId, nil
}

// FrontDoorFirewallPolicyIDInsensitively parses an FrontDoorFirewallPolicy ID into an FrontDoorFirewallPolicyId struct, insensitively
// This should only be used to parse an ID for rewriting, the FrontDoorFirewallPolicyID
// method should be used instead for validation etc.
//
// Whilst this may seem strange, this enables Terraform have consistent casing
// which works around issues in Core, whilst handling broken API responses.
func FrontDoorFirewallPolicyIDInsensitively(input string) (*FrontDoorFirewallPolicyId, error) {
	id, err := resourceids.ParseAzureResourceID(input)
	if err != nil {
		return nil, err
	}

	resourceId := FrontDoorFirewallPolicyId{
		SubscriptionId: id.SubscriptionID,
		ResourceGroup:  id.ResourceGroup,
	}

	if resourceId.SubscriptionId == "" {
		return nil, fmt.Errorf("ID was missing the 'subscriptions' element")
	}

	if resourceId.ResourceGroup == "" {
		return nil, fmt.Errorf("ID was missing the 'resourceGroups' element")
	}

	// find the correct casing for the 'frontDoorWebApplicationFirewallPolicies' segment
	frontDoorWebApplicationFirewallPoliciesKey := "frontDoorWebApplicationFirewallPolicies"
	for key := range id.Path {
		if strings.EqualFold(key, frontDoorWebApplicationFirewallPoliciesKey) {
			frontDoorWebApplicationFirewallPoliciesKey = key
			break
		}
	}
	if resourceId.FrontDoorWebApplicationFirewallPolicyName, err = id.PopSegment(frontDoorWebApplicationFirewallPoliciesKey); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &resourceId, nil
}
