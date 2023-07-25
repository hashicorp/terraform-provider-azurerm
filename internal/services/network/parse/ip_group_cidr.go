// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package parse

// NOTE: this file is generated via 'go:generate' - manual changes will be overwritten

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

type IpGroupCidrId struct {
	SubscriptionId string
	ResourceGroup  string
	IpGroupName    string
	CidrName       string
}

func NewIpGroupCidrID(subscriptionId, resourceGroup, ipGroupName, cidrName string) IpGroupCidrId {
	return IpGroupCidrId{
		SubscriptionId: subscriptionId,
		ResourceGroup:  resourceGroup,
		IpGroupName:    ipGroupName,
		CidrName:       cidrName,
	}
}

func (id IpGroupCidrId) String() string {
	segments := []string{
		fmt.Sprintf("Cidr Name %q", id.CidrName),
		fmt.Sprintf("Ip Group Name %q", id.IpGroupName),
		fmt.Sprintf("Resource Group %q", id.ResourceGroup),
	}
	segmentsStr := strings.Join(segments, " / ")
	return fmt.Sprintf("%s: (%s)", "Ip Group Cidr", segmentsStr)
}

func (id IpGroupCidrId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Network/ipGroups/%s/cidrs/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroup, id.IpGroupName, id.CidrName)
}

// IpGroupCidrID parses a IpGroupCidr ID into an IpGroupCidrId struct
func IpGroupCidrID(input string) (*IpGroupCidrId, error) {
	id, err := resourceids.ParseAzureResourceID(input)
	if err != nil {
		return nil, err
	}

	resourceId := IpGroupCidrId{
		SubscriptionId: id.SubscriptionID,
		ResourceGroup:  id.ResourceGroup,
	}

	if resourceId.SubscriptionId == "" {
		return nil, fmt.Errorf("ID was missing the 'subscriptions' element")
	}

	if resourceId.ResourceGroup == "" {
		return nil, fmt.Errorf("ID was missing the 'resourceGroups' element")
	}

	if resourceId.IpGroupName, err = id.PopSegment("ipGroups"); err != nil {
		return nil, err
	}
	if resourceId.CidrName, err = id.PopSegment("cidrs"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &resourceId, nil
}
