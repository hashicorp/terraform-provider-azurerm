// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package parse

// NOTE: this file is generated via 'go:generate' - manual changes will be overwritten

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

type ExpressRouteCircuitPeeringId struct {
	SubscriptionId          string
	ResourceGroup           string
	ExpressRouteCircuitName string
	PeeringName             string
}

func NewExpressRouteCircuitPeeringID(subscriptionId, resourceGroup, expressRouteCircuitName, peeringName string) ExpressRouteCircuitPeeringId {
	return ExpressRouteCircuitPeeringId{
		SubscriptionId:          subscriptionId,
		ResourceGroup:           resourceGroup,
		ExpressRouteCircuitName: expressRouteCircuitName,
		PeeringName:             peeringName,
	}
}

func (id ExpressRouteCircuitPeeringId) String() string {
	segments := []string{
		fmt.Sprintf("Peering Name %q", id.PeeringName),
		fmt.Sprintf("Express Route Circuit Name %q", id.ExpressRouteCircuitName),
		fmt.Sprintf("Resource Group %q", id.ResourceGroup),
	}
	segmentsStr := strings.Join(segments, " / ")
	return fmt.Sprintf("%s: (%s)", "Express Route Circuit Peering", segmentsStr)
}

func (id ExpressRouteCircuitPeeringId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Network/expressRouteCircuits/%s/peerings/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroup, id.ExpressRouteCircuitName, id.PeeringName)
}

// ExpressRouteCircuitPeeringID parses a ExpressRouteCircuitPeering ID into an ExpressRouteCircuitPeeringId struct
func ExpressRouteCircuitPeeringID(input string) (*ExpressRouteCircuitPeeringId, error) {
	id, err := resourceids.ParseAzureResourceID(input)
	if err != nil {
		return nil, fmt.Errorf("parsing %q as an ExpressRouteCircuitPeering ID: %+v", input, err)
	}

	resourceId := ExpressRouteCircuitPeeringId{
		SubscriptionId: id.SubscriptionID,
		ResourceGroup:  id.ResourceGroup,
	}

	if resourceId.SubscriptionId == "" {
		return nil, fmt.Errorf("ID was missing the 'subscriptions' element")
	}

	if resourceId.ResourceGroup == "" {
		return nil, fmt.Errorf("ID was missing the 'resourceGroups' element")
	}

	if resourceId.ExpressRouteCircuitName, err = id.PopSegment("expressRouteCircuits"); err != nil {
		return nil, err
	}
	if resourceId.PeeringName, err = id.PopSegment("peerings"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &resourceId, nil
}

// ExpressRouteCircuitPeeringIDInsensitively parses an ExpressRouteCircuitPeering ID into an ExpressRouteCircuitPeeringId struct, insensitively
// This should only be used to parse an ID for rewriting, the ExpressRouteCircuitPeeringID
// method should be used instead for validation etc.
//
// Whilst this may seem strange, this enables Terraform have consistent casing
// which works around issues in Core, whilst handling broken API responses.
func ExpressRouteCircuitPeeringIDInsensitively(input string) (*ExpressRouteCircuitPeeringId, error) {
	id, err := resourceids.ParseAzureResourceID(input)
	if err != nil {
		return nil, err
	}

	resourceId := ExpressRouteCircuitPeeringId{
		SubscriptionId: id.SubscriptionID,
		ResourceGroup:  id.ResourceGroup,
	}

	if resourceId.SubscriptionId == "" {
		return nil, fmt.Errorf("ID was missing the 'subscriptions' element")
	}

	if resourceId.ResourceGroup == "" {
		return nil, fmt.Errorf("ID was missing the 'resourceGroups' element")
	}

	// find the correct casing for the 'expressRouteCircuits' segment
	expressRouteCircuitsKey := "expressRouteCircuits"
	for key := range id.Path {
		if strings.EqualFold(key, expressRouteCircuitsKey) {
			expressRouteCircuitsKey = key
			break
		}
	}
	if resourceId.ExpressRouteCircuitName, err = id.PopSegment(expressRouteCircuitsKey); err != nil {
		return nil, err
	}

	// find the correct casing for the 'peerings' segment
	peeringsKey := "peerings"
	for key := range id.Path {
		if strings.EqualFold(key, peeringsKey) {
			peeringsKey = key
			break
		}
	}
	if resourceId.PeeringName, err = id.PopSegment(peeringsKey); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &resourceId, nil
}
