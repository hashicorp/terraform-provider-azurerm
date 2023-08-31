// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package parse

// NOTE: this file is generated via 'go:generate' - manual changes will be overwritten

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

type SoaRecordId struct {
	SubscriptionId string
	ResourceGroup  string
	DnszoneName    string
	SOAName        string
}

func NewSoaRecordID(subscriptionId, resourceGroup, dnszoneName, sOAName string) SoaRecordId {
	return SoaRecordId{
		SubscriptionId: subscriptionId,
		ResourceGroup:  resourceGroup,
		DnszoneName:    dnszoneName,
		SOAName:        sOAName,
	}
}

func (id SoaRecordId) String() string {
	segments := []string{
		fmt.Sprintf("SOA Name %q", id.SOAName),
		fmt.Sprintf("Dnszone Name %q", id.DnszoneName),
		fmt.Sprintf("Resource Group %q", id.ResourceGroup),
	}
	segmentsStr := strings.Join(segments, " / ")
	return fmt.Sprintf("%s: (%s)", "SOA Record", segmentsStr)
}

func (id SoaRecordId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Network/dnszones/%s/SOA/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroup, id.DnszoneName, id.SOAName)
}

// SoaRecordID parses a SOARecord ID into an SoaRecordId struct
func SoaRecordID(input string) (*SoaRecordId, error) {
	id, err := resourceids.ParseAzureResourceID(input)
	if err != nil {
		return nil, err
	}

	resourceId := SoaRecordId{
		SubscriptionId: id.SubscriptionID,
		ResourceGroup:  id.ResourceGroup,
	}

	if resourceId.SubscriptionId == "" {
		return nil, fmt.Errorf("ID was missing the 'subscriptions' element")
	}

	if resourceId.ResourceGroup == "" {
		return nil, fmt.Errorf("ID was missing the 'resourceGroups' element")
	}

	if resourceId.DnszoneName, err = id.PopSegment("dnszones"); err != nil {
		return nil, err
	}
	if resourceId.SOAName, err = id.PopSegment("SOA"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &resourceId, nil
}
