// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package parse

// NOTE: this file is generated via 'go:generate' - manual changes will be overwritten

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

type ExpressRouteCircuitAuthorizationId struct {
	SubscriptionId          string
	ResourceGroup           string
	ExpressRouteCircuitName string
	AuthorizationName       string
}

func NewExpressRouteCircuitAuthorizationID(subscriptionId, resourceGroup, expressRouteCircuitName, authorizationName string) ExpressRouteCircuitAuthorizationId {
	return ExpressRouteCircuitAuthorizationId{
		SubscriptionId:          subscriptionId,
		ResourceGroup:           resourceGroup,
		ExpressRouteCircuitName: expressRouteCircuitName,
		AuthorizationName:       authorizationName,
	}
}

func (id ExpressRouteCircuitAuthorizationId) String() string {
	segments := []string{
		fmt.Sprintf("Authorization Name %q", id.AuthorizationName),
		fmt.Sprintf("Express Route Circuit Name %q", id.ExpressRouteCircuitName),
		fmt.Sprintf("Resource Group %q", id.ResourceGroup),
	}
	segmentsStr := strings.Join(segments, " / ")
	return fmt.Sprintf("%s: (%s)", "Express Route Circuit Authorization", segmentsStr)
}

func (id ExpressRouteCircuitAuthorizationId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Network/expressRouteCircuits/%s/authorizations/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroup, id.ExpressRouteCircuitName, id.AuthorizationName)
}

// ExpressRouteCircuitAuthorizationID parses a ExpressRouteCircuitAuthorization ID into an ExpressRouteCircuitAuthorizationId struct
func ExpressRouteCircuitAuthorizationID(input string) (*ExpressRouteCircuitAuthorizationId, error) {
	id, err := resourceids.ParseAzureResourceID(input)
	if err != nil {
		return nil, fmt.Errorf("parsing %q as an ExpressRouteCircuitAuthorization ID: %+v", input, err)
	}

	resourceId := ExpressRouteCircuitAuthorizationId{
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
	if resourceId.AuthorizationName, err = id.PopSegment("authorizations"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &resourceId, nil
}
