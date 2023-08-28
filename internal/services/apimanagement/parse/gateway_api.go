// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package parse

// NOTE: this file is generated via 'go:generate' - manual changes will be overwritten

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

type GatewayApiId struct {
	SubscriptionId string
	ResourceGroup  string
	ServiceName    string
	GatewayName    string
	ApiName        string
}

func NewGatewayApiID(subscriptionId, resourceGroup, serviceName, gatewayName, apiName string) GatewayApiId {
	return GatewayApiId{
		SubscriptionId: subscriptionId,
		ResourceGroup:  resourceGroup,
		ServiceName:    serviceName,
		GatewayName:    gatewayName,
		ApiName:        apiName,
	}
}

func (id GatewayApiId) String() string {
	segments := []string{
		fmt.Sprintf("Api Name %q", id.ApiName),
		fmt.Sprintf("Gateway Name %q", id.GatewayName),
		fmt.Sprintf("Service Name %q", id.ServiceName),
		fmt.Sprintf("Resource Group %q", id.ResourceGroup),
	}
	segmentsStr := strings.Join(segments, " / ")
	return fmt.Sprintf("%s: (%s)", "Gateway Api", segmentsStr)
}

func (id GatewayApiId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.ApiManagement/service/%s/gateways/%s/apis/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroup, id.ServiceName, id.GatewayName, id.ApiName)
}

// GatewayApiID parses a GatewayApi ID into an GatewayApiId struct
func GatewayApiID(input string) (*GatewayApiId, error) {
	id, err := resourceids.ParseAzureResourceID(input)
	if err != nil {
		return nil, fmt.Errorf("parsing %q as an GatewayApi ID: %+v", input, err)
	}

	resourceId := GatewayApiId{
		SubscriptionId: id.SubscriptionID,
		ResourceGroup:  id.ResourceGroup,
	}

	if resourceId.SubscriptionId == "" {
		return nil, fmt.Errorf("ID was missing the 'subscriptions' element")
	}

	if resourceId.ResourceGroup == "" {
		return nil, fmt.Errorf("ID was missing the 'resourceGroups' element")
	}

	if resourceId.ServiceName, err = id.PopSegment("service"); err != nil {
		return nil, err
	}
	if resourceId.GatewayName, err = id.PopSegment("gateways"); err != nil {
		return nil, err
	}
	if resourceId.ApiName, err = id.PopSegment("apis"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &resourceId, nil
}
