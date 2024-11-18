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

type VirtualNetworkSwiftConnectionId struct {
	SubscriptionId string
	ResourceGroup  string
	SiteName       string
	ConfigName     string
}

func NewVirtualNetworkSwiftConnectionID(subscriptionId, resourceGroup, siteName, configName string) VirtualNetworkSwiftConnectionId {
	return VirtualNetworkSwiftConnectionId{
		SubscriptionId: subscriptionId,
		ResourceGroup:  resourceGroup,
		SiteName:       siteName,
		ConfigName:     configName,
	}
}

func (id VirtualNetworkSwiftConnectionId) String() string {
	segments := []string{
		fmt.Sprintf("Config Name %q", id.ConfigName),
		fmt.Sprintf("Site Name %q", id.SiteName),
		fmt.Sprintf("Resource Group %q", id.ResourceGroup),
	}
	segmentsStr := strings.Join(segments, " / ")
	return fmt.Sprintf("%s: (%s)", "Virtual Network Swift Connection", segmentsStr)
}

func (id VirtualNetworkSwiftConnectionId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Web/sites/%s/config/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroup, id.SiteName, id.ConfigName)
}

// VirtualNetworkSwiftConnectionID parses a VirtualNetworkSwiftConnection ID into an VirtualNetworkSwiftConnectionId struct
func VirtualNetworkSwiftConnectionID(input string) (*VirtualNetworkSwiftConnectionId, error) {
	id, err := resourceids.ParseAzureResourceID(input)
	if err != nil {
		return nil, fmt.Errorf("parsing %q as an VirtualNetworkSwiftConnection ID: %+v", input, err)
	}

	resourceId := VirtualNetworkSwiftConnectionId{
		SubscriptionId: id.SubscriptionID,
		ResourceGroup:  id.ResourceGroup,
	}

	if resourceId.SubscriptionId == "" {
		return nil, errors.New("ID was missing the 'subscriptions' element")
	}

	if resourceId.ResourceGroup == "" {
		return nil, errors.New("ID was missing the 'resourceGroups' element")
	}

	if resourceId.SiteName, err = id.PopSegment("sites"); err != nil {
		return nil, err
	}
	if resourceId.ConfigName, err = id.PopSegment("config"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &resourceId, nil
}
