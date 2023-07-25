// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package parse

// NOTE: this file is generated via 'go:generate' - manual changes will be overwritten

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

type ExpressRoutePortId struct {
	SubscriptionId string
	ResourceGroup  string
	Name           string
}

func NewExpressRoutePortID(subscriptionId, resourceGroup, name string) ExpressRoutePortId {
	return ExpressRoutePortId{
		SubscriptionId: subscriptionId,
		ResourceGroup:  resourceGroup,
		Name:           name,
	}
}

func (id ExpressRoutePortId) String() string {
	segments := []string{
		fmt.Sprintf("Name %q", id.Name),
		fmt.Sprintf("Resource Group %q", id.ResourceGroup),
	}
	segmentsStr := strings.Join(segments, " / ")
	return fmt.Sprintf("%s: (%s)", "Express Route Port", segmentsStr)
}

func (id ExpressRoutePortId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Network/expressRoutePorts/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroup, id.Name)
}

// ExpressRoutePortID parses a ExpressRoutePort ID into an ExpressRoutePortId struct
func ExpressRoutePortID(input string) (*ExpressRoutePortId, error) {
	id, err := resourceids.ParseAzureResourceID(input)
	if err != nil {
		return nil, fmt.Errorf("parsing %q as an ExpressRoutePort ID: %+v", input, err)
	}

	resourceId := ExpressRoutePortId{
		SubscriptionId: id.SubscriptionID,
		ResourceGroup:  id.ResourceGroup,
	}

	if resourceId.SubscriptionId == "" {
		return nil, fmt.Errorf("ID was missing the 'subscriptions' element")
	}

	if resourceId.ResourceGroup == "" {
		return nil, fmt.Errorf("ID was missing the 'resourceGroups' element")
	}

	if resourceId.Name, err = id.PopSegment("expressRoutePorts"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &resourceId, nil
}

// ExpressRoutePortIDInsensitively parses an ExpressRoutePort ID into an ExpressRoutePortId struct, insensitively
// This should only be used to parse an ID for rewriting, the ExpressRoutePortID
// method should be used instead for validation etc.
//
// Whilst this may seem strange, this enables Terraform have consistent casing
// which works around issues in Core, whilst handling broken API responses.
func ExpressRoutePortIDInsensitively(input string) (*ExpressRoutePortId, error) {
	id, err := resourceids.ParseAzureResourceID(input)
	if err != nil {
		return nil, err
	}

	resourceId := ExpressRoutePortId{
		SubscriptionId: id.SubscriptionID,
		ResourceGroup:  id.ResourceGroup,
	}

	if resourceId.SubscriptionId == "" {
		return nil, fmt.Errorf("ID was missing the 'subscriptions' element")
	}

	if resourceId.ResourceGroup == "" {
		return nil, fmt.Errorf("ID was missing the 'resourceGroups' element")
	}

	// find the correct casing for the 'expressRoutePorts' segment
	expressRoutePortsKey := "expressRoutePorts"
	for key := range id.Path {
		if strings.EqualFold(key, expressRoutePortsKey) {
			expressRoutePortsKey = key
			break
		}
	}
	if resourceId.Name, err = id.PopSegment(expressRoutePortsKey); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &resourceId, nil
}
