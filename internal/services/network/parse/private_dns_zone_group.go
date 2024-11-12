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

type PrivateDnsZoneGroupId struct {
	SubscriptionId      string
	ResourceGroup       string
	PrivateEndpointName string
	Name                string
}

func NewPrivateDnsZoneGroupID(subscriptionId, resourceGroup, privateEndpointName, name string) PrivateDnsZoneGroupId {
	return PrivateDnsZoneGroupId{
		SubscriptionId:      subscriptionId,
		ResourceGroup:       resourceGroup,
		PrivateEndpointName: privateEndpointName,
		Name:                name,
	}
}

func (id PrivateDnsZoneGroupId) String() string {
	segments := []string{
		fmt.Sprintf("Name %q", id.Name),
		fmt.Sprintf("Private Endpoint Name %q", id.PrivateEndpointName),
		fmt.Sprintf("Resource Group %q", id.ResourceGroup),
	}
	segmentsStr := strings.Join(segments, " / ")
	return fmt.Sprintf("%s: (%s)", "Private Dns Zone Group", segmentsStr)
}

func (id PrivateDnsZoneGroupId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Network/privateEndpoints/%s/privateDnsZoneGroups/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroup, id.PrivateEndpointName, id.Name)
}

// PrivateDnsZoneGroupID parses a PrivateDnsZoneGroup ID into an PrivateDnsZoneGroupId struct
func PrivateDnsZoneGroupID(input string) (*PrivateDnsZoneGroupId, error) {
	id, err := resourceids.ParseAzureResourceID(input)
	if err != nil {
		return nil, fmt.Errorf("parsing %q as an PrivateDnsZoneGroup ID: %+v", input, err)
	}

	resourceId := PrivateDnsZoneGroupId{
		SubscriptionId: id.SubscriptionID,
		ResourceGroup:  id.ResourceGroup,
	}

	if resourceId.SubscriptionId == "" {
		return nil, errors.New("ID was missing the 'subscriptions' element")
	}

	if resourceId.ResourceGroup == "" {
		return nil, errors.New("ID was missing the 'resourceGroups' element")
	}

	if resourceId.PrivateEndpointName, err = id.PopSegment("privateEndpoints"); err != nil {
		return nil, err
	}
	if resourceId.Name, err = id.PopSegment("privateDnsZoneGroups"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &resourceId, nil
}
