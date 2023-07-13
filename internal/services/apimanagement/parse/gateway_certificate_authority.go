// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package parse

// NOTE: this file is generated via 'go:generate' - manual changes will be overwritten

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

type GatewayCertificateAuthorityId struct {
	SubscriptionId           string
	ResourceGroup            string
	ServiceName              string
	GatewayName              string
	CertificateAuthorityName string
}

func NewGatewayCertificateAuthorityID(subscriptionId, resourceGroup, serviceName, gatewayName, certificateAuthorityName string) GatewayCertificateAuthorityId {
	return GatewayCertificateAuthorityId{
		SubscriptionId:           subscriptionId,
		ResourceGroup:            resourceGroup,
		ServiceName:              serviceName,
		GatewayName:              gatewayName,
		CertificateAuthorityName: certificateAuthorityName,
	}
}

func (id GatewayCertificateAuthorityId) String() string {
	segments := []string{
		fmt.Sprintf("Certificate Authority Name %q", id.CertificateAuthorityName),
		fmt.Sprintf("Gateway Name %q", id.GatewayName),
		fmt.Sprintf("Service Name %q", id.ServiceName),
		fmt.Sprintf("Resource Group %q", id.ResourceGroup),
	}
	segmentsStr := strings.Join(segments, " / ")
	return fmt.Sprintf("%s: (%s)", "Gateway Certificate Authority", segmentsStr)
}

func (id GatewayCertificateAuthorityId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.ApiManagement/service/%s/gateways/%s/certificateAuthorities/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroup, id.ServiceName, id.GatewayName, id.CertificateAuthorityName)
}

// GatewayCertificateAuthorityID parses a GatewayCertificateAuthority ID into an GatewayCertificateAuthorityId struct
func GatewayCertificateAuthorityID(input string) (*GatewayCertificateAuthorityId, error) {
	id, err := resourceids.ParseAzureResourceID(input)
	if err != nil {
		return nil, fmt.Errorf("parsing %q as an GatewayCertificateAuthority ID: %+v", input, err)
	}

	resourceId := GatewayCertificateAuthorityId{
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
	if resourceId.CertificateAuthorityName, err = id.PopSegment("certificateAuthorities"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &resourceId, nil
}
