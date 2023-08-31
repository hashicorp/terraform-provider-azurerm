// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package parse

// NOTE: this file is generated via 'go:generate' - manual changes will be overwritten

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

type CertificateVersionlessId struct {
	SubscriptionId  string
	ResourceGroup   string
	VaultName       string
	CertificateName string
}

func NewCertificateVersionlessID(subscriptionId, resourceGroup, vaultName, certificateName string) CertificateVersionlessId {
	return CertificateVersionlessId{
		SubscriptionId:  subscriptionId,
		ResourceGroup:   resourceGroup,
		VaultName:       vaultName,
		CertificateName: certificateName,
	}
}

func (id CertificateVersionlessId) String() string {
	segments := []string{
		fmt.Sprintf("Certificate Name %q", id.CertificateName),
		fmt.Sprintf("Vault Name %q", id.VaultName),
		fmt.Sprintf("Resource Group %q", id.ResourceGroup),
	}
	segmentsStr := strings.Join(segments, " / ")
	return fmt.Sprintf("%s: (%s)", "Certificate Versionless", segmentsStr)
}

func (id CertificateVersionlessId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.KeyVault/vaults/%s/certificates/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroup, id.VaultName, id.CertificateName)
}

// CertificateVersionlessID parses a CertificateVersionless ID into an CertificateVersionlessId struct
func CertificateVersionlessID(input string) (*CertificateVersionlessId, error) {
	id, err := resourceids.ParseAzureResourceID(input)
	if err != nil {
		return nil, fmt.Errorf("parsing %q as an CertificateVersionless ID: %+v", input, err)
	}

	resourceId := CertificateVersionlessId{
		SubscriptionId: id.SubscriptionID,
		ResourceGroup:  id.ResourceGroup,
	}

	if resourceId.SubscriptionId == "" {
		return nil, fmt.Errorf("ID was missing the 'subscriptions' element")
	}

	if resourceId.ResourceGroup == "" {
		return nil, fmt.Errorf("ID was missing the 'resourceGroups' element")
	}

	if resourceId.VaultName, err = id.PopSegment("vaults"); err != nil {
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
