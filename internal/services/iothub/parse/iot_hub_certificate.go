// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package parse

// NOTE: this file is generated via 'go:generate' - manual changes will be overwritten

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

type IotHubCertificateId struct {
	SubscriptionId  string
	ResourceGroup   string
	IotHubName      string
	CertificateName string
}

func NewIotHubCertificateID(subscriptionId, resourceGroup, iotHubName, certificateName string) IotHubCertificateId {
	return IotHubCertificateId{
		SubscriptionId:  subscriptionId,
		ResourceGroup:   resourceGroup,
		IotHubName:      iotHubName,
		CertificateName: certificateName,
	}
}

func (id IotHubCertificateId) String() string {
	segments := []string{
		fmt.Sprintf("Certificate Name %q", id.CertificateName),
		fmt.Sprintf("Iot Hub Name %q", id.IotHubName),
		fmt.Sprintf("Resource Group %q", id.ResourceGroup),
	}
	segmentsStr := strings.Join(segments, " / ")
	return fmt.Sprintf("%s: (%s)", "Iot Hub Certificate", segmentsStr)
}

func (id IotHubCertificateId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Devices/iotHubs/%s/certificates/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroup, id.IotHubName, id.CertificateName)
}

// IotHubCertificateID parses a IotHubCertificate ID into an IotHubCertificateId struct
func IotHubCertificateID(input string) (*IotHubCertificateId, error) {
	id, err := resourceids.ParseAzureResourceID(input)
	if err != nil {
		return nil, fmt.Errorf("parsing %q as an IotHubCertificate ID: %+v", input, err)
	}

	resourceId := IotHubCertificateId{
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
	if resourceId.CertificateName, err = id.PopSegment("certificates"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &resourceId, nil
}

// IotHubCertificateIDInsensitively parses an IotHubCertificate ID into an IotHubCertificateId struct, insensitively
// This should only be used to parse an ID for rewriting, the IotHubCertificateID
// method should be used instead for validation etc.
//
// Whilst this may seem strange, this enables Terraform have consistent casing
// which works around issues in Core, whilst handling broken API responses.
func IotHubCertificateIDInsensitively(input string) (*IotHubCertificateId, error) {
	id, err := resourceids.ParseAzureResourceID(input)
	if err != nil {
		return nil, err
	}

	resourceId := IotHubCertificateId{
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

	// find the correct casing for the 'certificates' segment
	certificatesKey := "certificates"
	for key := range id.Path {
		if strings.EqualFold(key, certificatesKey) {
			certificatesKey = key
			break
		}
	}
	if resourceId.CertificateName, err = id.PopSegment(certificatesKey); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &resourceId, nil
}
