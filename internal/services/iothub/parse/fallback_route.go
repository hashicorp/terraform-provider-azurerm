// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package parse

// NOTE: this file is generated via 'go:generate' - manual changes will be overwritten

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

type FallbackRouteId struct {
	SubscriptionId    string
	ResourceGroup     string
	IotHubName        string
	FallbackRouteName string
}

func NewFallbackRouteID(subscriptionId, resourceGroup, iotHubName, fallbackRouteName string) FallbackRouteId {
	return FallbackRouteId{
		SubscriptionId:    subscriptionId,
		ResourceGroup:     resourceGroup,
		IotHubName:        iotHubName,
		FallbackRouteName: fallbackRouteName,
	}
}

func (id FallbackRouteId) String() string {
	segments := []string{
		fmt.Sprintf("Fallback Route Name %q", id.FallbackRouteName),
		fmt.Sprintf("Iot Hub Name %q", id.IotHubName),
		fmt.Sprintf("Resource Group %q", id.ResourceGroup),
	}
	segmentsStr := strings.Join(segments, " / ")
	return fmt.Sprintf("%s: (%s)", "Fallback Route", segmentsStr)
}

func (id FallbackRouteId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Devices/iotHubs/%s/fallbackRoute/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroup, id.IotHubName, id.FallbackRouteName)
}

// FallbackRouteID parses a FallbackRoute ID into an FallbackRouteId struct
func FallbackRouteID(input string) (*FallbackRouteId, error) {
	id, err := resourceids.ParseAzureResourceID(input)
	if err != nil {
		return nil, fmt.Errorf("parsing %q as an FallbackRoute ID: %+v", input, err)
	}

	resourceId := FallbackRouteId{
		SubscriptionId: id.SubscriptionID,
		ResourceGroup:  id.ResourceGroup,
	}

	if resourceId.SubscriptionId == "" {
		return nil, fmt.Errorf("ID was missing the 'subscriptions' element")
	}

	if resourceId.ResourceGroup == "" {
		return nil, fmt.Errorf("ID was missing the 'resourceGroups' element")
	}

	if resourceId.IotHubName, err = id.PopSegment("iotHubs"); err != nil {
		return nil, err
	}
	if resourceId.FallbackRouteName, err = id.PopSegment("fallbackRoute"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &resourceId, nil
}

// FallbackRouteIDInsensitively parses an FallbackRoute ID into an FallbackRouteId struct, insensitively
// This should only be used to parse an ID for rewriting, the FallbackRouteID
// method should be used instead for validation etc.
//
// Whilst this may seem strange, this enables Terraform have consistent casing
// which works around issues in Core, whilst handling broken API responses.
func FallbackRouteIDInsensitively(input string) (*FallbackRouteId, error) {
	id, err := resourceids.ParseAzureResourceID(input)
	if err != nil {
		return nil, err
	}

	resourceId := FallbackRouteId{
		SubscriptionId: id.SubscriptionID,
		ResourceGroup:  id.ResourceGroup,
	}

	if resourceId.SubscriptionId == "" {
		return nil, fmt.Errorf("ID was missing the 'subscriptions' element")
	}

	if resourceId.ResourceGroup == "" {
		return nil, fmt.Errorf("ID was missing the 'resourceGroups' element")
	}

	// find the correct casing for the 'iotHubs' segment
	iotHubsKey := "iotHubs"
	for key := range id.Path {
		if strings.EqualFold(key, iotHubsKey) {
			iotHubsKey = key
			break
		}
	}
	if resourceId.IotHubName, err = id.PopSegment(iotHubsKey); err != nil {
		return nil, err
	}

	// find the correct casing for the 'fallbackRoute' segment
	fallbackRouteKey := "fallbackRoute"
	for key := range id.Path {
		if strings.EqualFold(key, fallbackRouteKey) {
			fallbackRouteKey = key
			break
		}
	}
	if resourceId.FallbackRouteName, err = id.PopSegment(fallbackRouteKey); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &resourceId, nil
}
