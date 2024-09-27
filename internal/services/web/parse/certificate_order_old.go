// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package parse

// NOTE: this file is generated via 'go:generate' - manual changes will be overwritten

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

type CertificateOrderOldId struct {
	SubscriptionId       string
	ResourceGroup        string
	CertificateOrderName string
}

func NewCertificateOrderOldID(subscriptionId, resourceGroup, certificateOrderName string) CertificateOrderOldId {
	return CertificateOrderOldId{
		SubscriptionId:       subscriptionId,
		ResourceGroup:        resourceGroup,
		CertificateOrderName: certificateOrderName,
	}
}

func (id CertificateOrderOldId) String() string {
	segments := []string{
		fmt.Sprintf("Certificate Order Name %q", id.CertificateOrderName),
		fmt.Sprintf("Resource Group %q", id.ResourceGroup),
	}
	segmentsStr := strings.Join(segments, " / ")
	return fmt.Sprintf("%s: (%s)", "Certificate Order Old", segmentsStr)
}

func (id CertificateOrderOldId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Web/certificateOrders/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroup, id.CertificateOrderName)
}

// CertificateOrderOldID parses a CertificateOrderOld ID into an CertificateOrderOldId struct
func CertificateOrderOldID(input string) (*CertificateOrderOldId, error) {
	id, err := resourceids.ParseAzureResourceID(input)
	if err != nil {
		return nil, fmt.Errorf("parsing %q as an CertificateOrderOld ID: %+v", input, err)
	}

	resourceId := CertificateOrderOldId{
		SubscriptionId: id.SubscriptionID,
		ResourceGroup:  id.ResourceGroup,
	}

	if resourceId.SubscriptionId == "" {
		return nil, fmt.Errorf("ID was missing the 'subscriptions' element")
	}

	if resourceId.ResourceGroup == "" {
		return nil, fmt.Errorf("ID was missing the 'resourceGroups' element")
	}

	if resourceId.CertificateOrderName, err = id.PopSegment("certificateOrders"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &resourceId, nil
}
