package parse

// NOTE: this file is generated via 'go:generate' - manual changes will be overwritten

import (
	"fmt"
	"strings"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
)

type CnameRecordId struct {
	SubscriptionId string
	ResourceGroup  string
	DnszoneName    string
	CNAMEName      string
}

func NewCnameRecordID(subscriptionId, resourceGroup, dnszoneName, cNAMEName string) CnameRecordId {
	return CnameRecordId{
		SubscriptionId: subscriptionId,
		ResourceGroup:  resourceGroup,
		DnszoneName:    dnszoneName,
		CNAMEName:      cNAMEName,
	}
}

func (id CnameRecordId) String() string {
	segments := []string{
		fmt.Sprintf("C N A M E Name %q", id.CNAMEName),
		fmt.Sprintf("Dnszone Name %q", id.DnszoneName),
		fmt.Sprintf("Resource Group %q", id.ResourceGroup),
	}
	segmentsStr := strings.Join(segments, " / ")
	return fmt.Sprintf("%s: (%s)", "Cname Record", segmentsStr)
}

func (id CnameRecordId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Network/dnszones/%s/CNAME/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroup, id.DnszoneName, id.CNAMEName)
}

// CnameRecordID parses a CnameRecord ID into an CnameRecordId struct
func CnameRecordID(input string) (*CnameRecordId, error) {
	id, err := azure.ParseAzureResourceID(input)
	if err != nil {
		return nil, err
	}

	resourceId := CnameRecordId{
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
	if resourceId.CNAMEName, err = id.PopSegment("CNAME"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &resourceId, nil
}
