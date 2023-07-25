// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package parse

// NOTE: this file is generated via 'go:generate' - manual changes will be overwritten

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

type ExpressRoutePortAuthorizationId struct {
	SubscriptionId       string
	ResourceGroup        string
	ExpressRoutePortName string
	AuthorizationName    string
}

func NewExpressRoutePortAuthorizationID(subscriptionId, resourceGroup, expressRoutePortName, authorizationName string) ExpressRoutePortAuthorizationId {
	return ExpressRoutePortAuthorizationId{
		SubscriptionId:       subscriptionId,
		ResourceGroup:        resourceGroup,
		ExpressRoutePortName: expressRoutePortName,
		AuthorizationName:    authorizationName,
	}
}

func (id ExpressRoutePortAuthorizationId) String() string {
	segments := []string{
		fmt.Sprintf("Authorization Name %q", id.AuthorizationName),
		fmt.Sprintf("Express Route Port Name %q", id.ExpressRoutePortName),
		fmt.Sprintf("Resource Group %q", id.ResourceGroup),
	}
	segmentsStr := strings.Join(segments, " / ")
	return fmt.Sprintf("%s: (%s)", "Express Route Port Authorization", segmentsStr)
}

func (id ExpressRoutePortAuthorizationId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Network/expressRoutePorts/%s/authorizations/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroup, id.ExpressRoutePortName, id.AuthorizationName)
}

// ExpressRoutePortAuthorizationID parses a ExpressRoutePortAuthorization ID into an ExpressRoutePortAuthorizationId struct
func ExpressRoutePortAuthorizationID(input string) (*ExpressRoutePortAuthorizationId, error) {
	id, err := resourceids.ParseAzureResourceID(input)
	if err != nil {
		return nil, fmt.Errorf("parsing %q as an ExpressRoutePortAuthorization ID: %+v", input, err)
	}

	resourceId := ExpressRoutePortAuthorizationId{
		SubscriptionId: id.SubscriptionID,
		ResourceGroup:  id.ResourceGroup,
	}

	if resourceId.SubscriptionId == "" {
		return nil, fmt.Errorf("ID was missing the 'subscriptions' element")
	}

	if resourceId.ResourceGroup == "" {
		return nil, fmt.Errorf("ID was missing the 'resourceGroups' element")
	}

	if resourceId.ExpressRoutePortName, err = id.PopSegment("expressRoutePorts"); err != nil {
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
