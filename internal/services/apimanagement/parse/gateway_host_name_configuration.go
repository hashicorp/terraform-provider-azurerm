// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package parse

// NOTE: this file is generated via 'go:generate' - manual changes will be overwritten

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

type GatewayHostNameConfigurationId struct {
	SubscriptionId            string
	ResourceGroup             string
	ServiceName               string
	GatewayName               string
	HostnameConfigurationName string
}

func NewGatewayHostNameConfigurationID(subscriptionId, resourceGroup, serviceName, gatewayName, hostnameConfigurationName string) GatewayHostNameConfigurationId {
	return GatewayHostNameConfigurationId{
		SubscriptionId:            subscriptionId,
		ResourceGroup:             resourceGroup,
		ServiceName:               serviceName,
		GatewayName:               gatewayName,
		HostnameConfigurationName: hostnameConfigurationName,
	}
}

func (id GatewayHostNameConfigurationId) String() string {
	segments := []string{
		fmt.Sprintf("Hostname Configuration Name %q", id.HostnameConfigurationName),
		fmt.Sprintf("Gateway Name %q", id.GatewayName),
		fmt.Sprintf("Service Name %q", id.ServiceName),
		fmt.Sprintf("Resource Group %q", id.ResourceGroup),
	}
	segmentsStr := strings.Join(segments, " / ")
	return fmt.Sprintf("%s: (%s)", "Gateway Host Name Configuration", segmentsStr)
}

func (id GatewayHostNameConfigurationId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.ApiManagement/service/%s/gateways/%s/hostnameConfigurations/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroup, id.ServiceName, id.GatewayName, id.HostnameConfigurationName)
}

// GatewayHostNameConfigurationID parses a GatewayHostNameConfiguration ID into an GatewayHostNameConfigurationId struct
func GatewayHostNameConfigurationID(input string) (*GatewayHostNameConfigurationId, error) {
	id, err := resourceids.ParseAzureResourceID(input)
	if err != nil {
		return nil, fmt.Errorf("parsing %q as an GatewayHostNameConfiguration ID: %+v", input, err)
	}

	resourceId := GatewayHostNameConfigurationId{
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
	if resourceId.HostnameConfigurationName, err = id.PopSegment("hostnameConfigurations"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &resourceId, nil
}
