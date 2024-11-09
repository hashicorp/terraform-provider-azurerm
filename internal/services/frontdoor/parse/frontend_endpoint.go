// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package parse

// NOTE: this file is generated via 'go:generate' - manual changes will be overwritten

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

type FrontendEndpointId struct {
	SubscriptionId string
	ResourceGroup  string
	FrontDoorName  string
	Name           string
}

func NewFrontendEndpointID(subscriptionId, resourceGroup, frontDoorName, name string) FrontendEndpointId {
	return FrontendEndpointId{
		SubscriptionId: subscriptionId,
		ResourceGroup:  resourceGroup,
		FrontDoorName:  frontDoorName,
		Name:           name,
	}
}

func (id FrontendEndpointId) String() string {
	segments := []string{
		fmt.Sprintf("Name %q", id.Name),
		fmt.Sprintf("Front Door Name %q", id.FrontDoorName),
		fmt.Sprintf("Resource Group %q", id.ResourceGroup),
	}
	segmentsStr := strings.Join(segments, " / ")
	return fmt.Sprintf("%s: (%s)", "Frontend Endpoint", segmentsStr)
}

func (id FrontendEndpointId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Network/frontDoors/%s/frontendEndpoints/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroup, id.FrontDoorName, id.Name)
}

// FrontendEndpointID parses a FrontendEndpoint ID into an FrontendEndpointId struct
func FrontendEndpointID(input string) (*FrontendEndpointId, error) {
	id, err := resourceids.ParseAzureResourceID(input)
	if err != nil {
		return nil, fmt.Errorf("parsing %q as an FrontendEndpoint ID: %+v", input, err)
	}

	resourceId := FrontendEndpointId{
		SubscriptionId: id.SubscriptionID,
		ResourceGroup:  id.ResourceGroup,
	}

	if resourceId.SubscriptionId == "" {
		return nil, fmt.Errorf("ID was missing the 'subscriptions' element")
	}

	if resourceId.ResourceGroup == "" {
		return nil, fmt.Errorf("ID was missing the 'resourceGroups' element")
	}

	if resourceId.FrontDoorName, err = id.PopSegment("frontDoors"); err != nil {
		return nil, err
	}
	if resourceId.Name, err = id.PopSegment("frontendEndpoints"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &resourceId, nil
}

// FrontendEndpointIDInsensitively parses an FrontendEndpoint ID into an FrontendEndpointId struct, insensitively
// This should only be used to parse an ID for rewriting, the FrontendEndpointID
// method should be used instead for validation etc.
//
// Whilst this may seem strange, this enables Terraform have consistent casing
// which works around issues in Core, whilst handling broken API responses.
func FrontendEndpointIDInsensitively(input string) (*FrontendEndpointId, error) {
	id, err := resourceids.ParseAzureResourceID(input)
	if err != nil {
		return nil, err
	}

	resourceId := FrontendEndpointId{
		SubscriptionId: id.SubscriptionID,
		ResourceGroup:  id.ResourceGroup,
	}

	if resourceId.SubscriptionId == "" {
		return nil, fmt.Errorf("ID was missing the 'subscriptions' element")
	}

	if resourceId.ResourceGroup == "" {
		return nil, fmt.Errorf("ID was missing the 'resourceGroups' element")
	}

	// find the correct casing for the 'frontDoors' segment
	frontDoorsKey := "frontDoors"
	for key := range id.Path {
		if strings.EqualFold(key, frontDoorsKey) {
			frontDoorsKey = key
			break
		}
	}
	if resourceId.FrontDoorName, err = id.PopSegment(frontDoorsKey); err != nil {
		return nil, err
	}

	// find the correct casing for the 'frontendEndpoints' segment
	frontendEndpointsKey := "frontendEndpoints"
	for key := range id.Path {
		if strings.EqualFold(key, frontendEndpointsKey) {
			frontendEndpointsKey = key
			break
		}
	}
	if resourceId.Name, err = id.PopSegment(frontendEndpointsKey); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &resourceId, nil
}
