// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package parse

// NOTE: this file is generated via 'go:generate' - manual changes will be overwritten

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

type SpringCloudGatewayRouteConfigId struct {
	SubscriptionId  string
	ResourceGroup   string
	SpringName      string
	GatewayName     string
	RouteConfigName string
}

func NewSpringCloudGatewayRouteConfigID(subscriptionId, resourceGroup, springName, gatewayName, routeConfigName string) SpringCloudGatewayRouteConfigId {
	return SpringCloudGatewayRouteConfigId{
		SubscriptionId:  subscriptionId,
		ResourceGroup:   resourceGroup,
		SpringName:      springName,
		GatewayName:     gatewayName,
		RouteConfigName: routeConfigName,
	}
}

func (id SpringCloudGatewayRouteConfigId) String() string {
	segments := []string{
		fmt.Sprintf("Route Config Name %q", id.RouteConfigName),
		fmt.Sprintf("Gateway Name %q", id.GatewayName),
		fmt.Sprintf("Spring Name %q", id.SpringName),
		fmt.Sprintf("Resource Group %q", id.ResourceGroup),
	}
	segmentsStr := strings.Join(segments, " / ")
	return fmt.Sprintf("%s: (%s)", "Spring Cloud Gateway Route Config", segmentsStr)
}

func (id SpringCloudGatewayRouteConfigId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.AppPlatform/spring/%s/gateways/%s/routeConfigs/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroup, id.SpringName, id.GatewayName, id.RouteConfigName)
}

// SpringCloudGatewayRouteConfigID parses a SpringCloudGatewayRouteConfig ID into an SpringCloudGatewayRouteConfigId struct
func SpringCloudGatewayRouteConfigID(input string) (*SpringCloudGatewayRouteConfigId, error) {
	id, err := resourceids.ParseAzureResourceID(input)
	if err != nil {
		return nil, fmt.Errorf("parsing %q as an SpringCloudGatewayRouteConfig ID: %+v", input, err)
	}

	resourceId := SpringCloudGatewayRouteConfigId{
		SubscriptionId: id.SubscriptionID,
		ResourceGroup:  id.ResourceGroup,
	}

	if resourceId.SubscriptionId == "" {
		return nil, fmt.Errorf("ID was missing the 'subscriptions' element")
	}

	if resourceId.ResourceGroup == "" {
		return nil, fmt.Errorf("ID was missing the 'resourceGroups' element")
	}

	if resourceId.SpringName, err = id.PopSegment("spring"); err != nil {
		return nil, err
	}
	if resourceId.GatewayName, err = id.PopSegment("gateways"); err != nil {
		return nil, err
	}
	if resourceId.RouteConfigName, err = id.PopSegment("routeConfigs"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &resourceId, nil
}

// SpringCloudGatewayRouteConfigIDInsensitively parses an SpringCloudGatewayRouteConfig ID into an SpringCloudGatewayRouteConfigId struct, insensitively
// This should only be used to parse an ID for rewriting, the SpringCloudGatewayRouteConfigID
// method should be used instead for validation etc.
//
// Whilst this may seem strange, this enables Terraform have consistent casing
// which works around issues in Core, whilst handling broken API responses.
func SpringCloudGatewayRouteConfigIDInsensitively(input string) (*SpringCloudGatewayRouteConfigId, error) {
	id, err := resourceids.ParseAzureResourceID(input)
	if err != nil {
		return nil, err
	}

	resourceId := SpringCloudGatewayRouteConfigId{
		SubscriptionId: id.SubscriptionID,
		ResourceGroup:  id.ResourceGroup,
	}

	if resourceId.SubscriptionId == "" {
		return nil, fmt.Errorf("ID was missing the 'subscriptions' element")
	}

	if resourceId.ResourceGroup == "" {
		return nil, fmt.Errorf("ID was missing the 'resourceGroups' element")
	}

	// find the correct casing for the 'spring' segment
	springKey := "spring"
	for key := range id.Path {
		if strings.EqualFold(key, springKey) {
			springKey = key
			break
		}
	}
	if resourceId.SpringName, err = id.PopSegment(springKey); err != nil {
		return nil, err
	}

	// find the correct casing for the 'gateways' segment
	gatewaysKey := "gateways"
	for key := range id.Path {
		if strings.EqualFold(key, gatewaysKey) {
			gatewaysKey = key
			break
		}
	}
	if resourceId.GatewayName, err = id.PopSegment(gatewaysKey); err != nil {
		return nil, err
	}

	// find the correct casing for the 'routeConfigs' segment
	routeConfigsKey := "routeConfigs"
	for key := range id.Path {
		if strings.EqualFold(key, routeConfigsKey) {
			routeConfigsKey = key
			break
		}
	}
	if resourceId.RouteConfigName, err = id.PopSegment(routeConfigsKey); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &resourceId, nil
}
