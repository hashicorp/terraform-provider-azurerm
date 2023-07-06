// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package parse

// NOTE: this file is generated via 'go:generate' - manual changes will be overwritten

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

type SecurityRuleId struct {
	SubscriptionId           string
	ResourceGroup            string
	NetworkSecurityGroupName string
	Name                     string
}

func NewSecurityRuleID(subscriptionId, resourceGroup, networkSecurityGroupName, name string) SecurityRuleId {
	return SecurityRuleId{
		SubscriptionId:           subscriptionId,
		ResourceGroup:            resourceGroup,
		NetworkSecurityGroupName: networkSecurityGroupName,
		Name:                     name,
	}
}

func (id SecurityRuleId) String() string {
	segments := []string{
		fmt.Sprintf("Name %q", id.Name),
		fmt.Sprintf("Network Security Group Name %q", id.NetworkSecurityGroupName),
		fmt.Sprintf("Resource Group %q", id.ResourceGroup),
	}
	segmentsStr := strings.Join(segments, " / ")
	return fmt.Sprintf("%s: (%s)", "Security Rule", segmentsStr)
}

func (id SecurityRuleId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Network/networkSecurityGroups/%s/securityRules/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroup, id.NetworkSecurityGroupName, id.Name)
}

// SecurityRuleID parses a SecurityRule ID into an SecurityRuleId struct
func SecurityRuleID(input string) (*SecurityRuleId, error) {
	id, err := resourceids.ParseAzureResourceID(input)
	if err != nil {
		return nil, fmt.Errorf("parsing %q as an SecurityRule ID: %+v", input, err)
	}

	resourceId := SecurityRuleId{
		SubscriptionId: id.SubscriptionID,
		ResourceGroup:  id.ResourceGroup,
	}

	if resourceId.SubscriptionId == "" {
		return nil, fmt.Errorf("ID was missing the 'subscriptions' element")
	}

	if resourceId.ResourceGroup == "" {
		return nil, fmt.Errorf("ID was missing the 'resourceGroups' element")
	}

	if resourceId.NetworkSecurityGroupName, err = id.PopSegment("networkSecurityGroups"); err != nil {
		return nil, err
	}
	if resourceId.Name, err = id.PopSegment("securityRules"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &resourceId, nil
}
