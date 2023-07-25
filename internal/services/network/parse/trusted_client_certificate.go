// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package parse

// NOTE: this file is generated via 'go:generate' - manual changes will be overwritten

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

type TrustedClientCertificateId struct {
	SubscriptionId         string
	ResourceGroup          string
	ApplicationGatewayName string
	Name                   string
}

func NewTrustedClientCertificateID(subscriptionId, resourceGroup, applicationGatewayName, name string) TrustedClientCertificateId {
	return TrustedClientCertificateId{
		SubscriptionId:         subscriptionId,
		ResourceGroup:          resourceGroup,
		ApplicationGatewayName: applicationGatewayName,
		Name:                   name,
	}
}

func (id TrustedClientCertificateId) String() string {
	segments := []string{
		fmt.Sprintf("Name %q", id.Name),
		fmt.Sprintf("Application Gateway Name %q", id.ApplicationGatewayName),
		fmt.Sprintf("Resource Group %q", id.ResourceGroup),
	}
	segmentsStr := strings.Join(segments, " / ")
	return fmt.Sprintf("%s: (%s)", "Trusted Client Certificate", segmentsStr)
}

func (id TrustedClientCertificateId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Network/applicationGateways/%s/trustedClientCertificates/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroup, id.ApplicationGatewayName, id.Name)
}

// TrustedClientCertificateID parses a TrustedClientCertificate ID into an TrustedClientCertificateId struct
func TrustedClientCertificateID(input string) (*TrustedClientCertificateId, error) {
	id, err := resourceids.ParseAzureResourceID(input)
	if err != nil {
		return nil, fmt.Errorf("parsing %q as an TrustedClientCertificate ID: %+v", input, err)
	}

	resourceId := TrustedClientCertificateId{
		SubscriptionId: id.SubscriptionID,
		ResourceGroup:  id.ResourceGroup,
	}

	if resourceId.SubscriptionId == "" {
		return nil, fmt.Errorf("ID was missing the 'subscriptions' element")
	}

	if resourceId.ResourceGroup == "" {
		return nil, fmt.Errorf("ID was missing the 'resourceGroups' element")
	}

	if resourceId.ApplicationGatewayName, err = id.PopSegment("applicationGateways"); err != nil {
		return nil, err
	}
	if resourceId.Name, err = id.PopSegment("trustedClientCertificates"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &resourceId, nil
}

// TrustedClientCertificateIDInsensitively parses an TrustedClientCertificate ID into an TrustedClientCertificateId struct, insensitively
// This should only be used to parse an ID for rewriting, the TrustedClientCertificateID
// method should be used instead for validation etc.
//
// Whilst this may seem strange, this enables Terraform have consistent casing
// which works around issues in Core, whilst handling broken API responses.
func TrustedClientCertificateIDInsensitively(input string) (*TrustedClientCertificateId, error) {
	id, err := resourceids.ParseAzureResourceID(input)
	if err != nil {
		return nil, err
	}

	resourceId := TrustedClientCertificateId{
		SubscriptionId: id.SubscriptionID,
		ResourceGroup:  id.ResourceGroup,
	}

	if resourceId.SubscriptionId == "" {
		return nil, fmt.Errorf("ID was missing the 'subscriptions' element")
	}

	if resourceId.ResourceGroup == "" {
		return nil, fmt.Errorf("ID was missing the 'resourceGroups' element")
	}

	// find the correct casing for the 'applicationGateways' segment
	applicationGatewaysKey := "applicationGateways"
	for key := range id.Path {
		if strings.EqualFold(key, applicationGatewaysKey) {
			applicationGatewaysKey = key
			break
		}
	}
	if resourceId.ApplicationGatewayName, err = id.PopSegment(applicationGatewaysKey); err != nil {
		return nil, err
	}

	// find the correct casing for the 'trustedClientCertificates' segment
	trustedClientCertificatesKey := "trustedClientCertificates"
	for key := range id.Path {
		if strings.EqualFold(key, trustedClientCertificatesKey) {
			trustedClientCertificatesKey = key
			break
		}
	}
	if resourceId.Name, err = id.PopSegment(trustedClientCertificatesKey); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &resourceId, nil
}
