// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package parse

// NOTE: this file is generated via 'go:generate' - manual changes will be overwritten

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

type CustomIpv4PrefixId struct {
	SubscriptionId      string
	ResourceGroup       string
	CustomIpPrefixeName string
}

func NewCustomIpv4PrefixID(subscriptionId, resourceGroup, customIpPrefixeName string) CustomIpv4PrefixId {
	return CustomIpv4PrefixId{
		SubscriptionId:      subscriptionId,
		ResourceGroup:       resourceGroup,
		CustomIpPrefixeName: customIpPrefixeName,
	}
}

func (id CustomIpv4PrefixId) String() string {
	segments := []string{
		fmt.Sprintf("Custom Ip Prefixe Name %q", id.CustomIpPrefixeName),
		fmt.Sprintf("Resource Group %q", id.ResourceGroup),
	}
	segmentsStr := strings.Join(segments, " / ")
	return fmt.Sprintf("%s: (%s)", "Custom Ipv4 Prefix", segmentsStr)
}

func (id CustomIpv4PrefixId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Network/customIpPrefixes/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroup, id.CustomIpPrefixeName)
}

// CustomIpv4PrefixID parses a CustomIpv4Prefix ID into an CustomIpv4PrefixId struct
func CustomIpv4PrefixID(input string) (*CustomIpv4PrefixId, error) {
	id, err := resourceids.ParseAzureResourceID(input)
	if err != nil {
		return nil, err
	}

	resourceId := CustomIpv4PrefixId{
		SubscriptionId: id.SubscriptionID,
		ResourceGroup:  id.ResourceGroup,
	}

	if resourceId.SubscriptionId == "" {
		return nil, fmt.Errorf("ID was missing the 'subscriptions' element")
	}

	if resourceId.ResourceGroup == "" {
		return nil, fmt.Errorf("ID was missing the 'resourceGroups' element")
	}

	if resourceId.CustomIpPrefixeName, err = id.PopSegment("customIpPrefixes"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &resourceId, nil
}
