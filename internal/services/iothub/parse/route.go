// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package parse

// NOTE: this file is generated via 'go:generate' - manual changes will be overwritten

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

type RouteId struct {
	SubscriptionId string
	ResourceGroup  string
	IotHubName     string
	Name           string
}

func NewRouteID(subscriptionId, resourceGroup, iotHubName, name string) RouteId {
	return RouteId{
		SubscriptionId: subscriptionId,
		ResourceGroup:  resourceGroup,
		IotHubName:     iotHubName,
		Name:           name,
	}
}

func (id RouteId) String() string {
	segments := []string{
		fmt.Sprintf("Name %q", id.Name),
		fmt.Sprintf("Iot Hub Name %q", id.IotHubName),
		fmt.Sprintf("Resource Group %q", id.ResourceGroup),
	}
	segmentsStr := strings.Join(segments, " / ")
	return fmt.Sprintf("%s: (%s)", "Route", segmentsStr)
}

func (id RouteId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Devices/iotHubs/%s/routes/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroup, id.IotHubName, id.Name)
}

// RouteID parses a Route ID into an RouteId struct
func RouteID(input string) (*RouteId, error) {
	id, err := resourceids.ParseAzureResourceID(input)
	if err != nil {
		return nil, fmt.Errorf("parsing %q as an Route ID: %+v", input, err)
	}

	resourceId := RouteId{
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
	if resourceId.Name, err = id.PopSegment("routes"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &resourceId, nil
}

// RouteIDInsensitively parses an Route ID into an RouteId struct, insensitively
// This should only be used to parse an ID for rewriting, the RouteID
// method should be used instead for validation etc.
//
// Whilst this may seem strange, this enables Terraform have consistent casing
// which works around issues in Core, whilst handling broken API responses.
func RouteIDInsensitively(input string) (*RouteId, error) {
	id, err := resourceids.ParseAzureResourceID(input)
	if err != nil {
		return nil, err
	}

	resourceId := RouteId{
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

	// find the correct casing for the 'routes' segment
	routesKey := "routes"
	for key := range id.Path {
		if strings.EqualFold(key, routesKey) {
			routesKey = key
			break
		}
	}
	if resourceId.Name, err = id.PopSegment(routesKey); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &resourceId, nil
}
