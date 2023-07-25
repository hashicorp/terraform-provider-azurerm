// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package parse

// NOTE: this file is generated via 'go:generate' - manual changes will be overwritten

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

type SpringCloudCertificateId struct {
	SubscriptionId  string
	ResourceGroup   string
	SpringName      string
	CertificateName string
}

func NewSpringCloudCertificateID(subscriptionId, resourceGroup, springName, certificateName string) SpringCloudCertificateId {
	return SpringCloudCertificateId{
		SubscriptionId:  subscriptionId,
		ResourceGroup:   resourceGroup,
		SpringName:      springName,
		CertificateName: certificateName,
	}
}

func (id SpringCloudCertificateId) String() string {
	segments := []string{
		fmt.Sprintf("Certificate Name %q", id.CertificateName),
		fmt.Sprintf("Spring Name %q", id.SpringName),
		fmt.Sprintf("Resource Group %q", id.ResourceGroup),
	}
	segmentsStr := strings.Join(segments, " / ")
	return fmt.Sprintf("%s: (%s)", "Spring Cloud Certificate", segmentsStr)
}

func (id SpringCloudCertificateId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.AppPlatform/spring/%s/certificates/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroup, id.SpringName, id.CertificateName)
}

// SpringCloudCertificateID parses a SpringCloudCertificate ID into an SpringCloudCertificateId struct
func SpringCloudCertificateID(input string) (*SpringCloudCertificateId, error) {
	id, err := resourceids.ParseAzureResourceID(input)
	if err != nil {
		return nil, fmt.Errorf("parsing %q as an SpringCloudCertificate ID: %+v", input, err)
	}

	resourceId := SpringCloudCertificateId{
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
	if resourceId.CertificateName, err = id.PopSegment("certificates"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &resourceId, nil
}

// SpringCloudCertificateIDInsensitively parses an SpringCloudCertificate ID into an SpringCloudCertificateId struct, insensitively
// This should only be used to parse an ID for rewriting, the SpringCloudCertificateID
// method should be used instead for validation etc.
//
// Whilst this may seem strange, this enables Terraform have consistent casing
// which works around issues in Core, whilst handling broken API responses.
func SpringCloudCertificateIDInsensitively(input string) (*SpringCloudCertificateId, error) {
	id, err := resourceids.ParseAzureResourceID(input)
	if err != nil {
		return nil, err
	}

	resourceId := SpringCloudCertificateId{
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
