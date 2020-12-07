package parse

// NOTE: this file is generated via 'go:generate' - manual changes will be overwritten

import (
	"fmt"
	"strings"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
)

type SrvRecordId struct {
	SubscriptionId string
	ResourceGroup  string
	DnszoneName    string
	SRVName        string
}

func NewSrvRecordID(subscriptionId, resourceGroup, dnszoneName, sRVName string) SrvRecordId {
	return SrvRecordId{
		SubscriptionId: subscriptionId,
		ResourceGroup:  resourceGroup,
		DnszoneName:    dnszoneName,
		SRVName:        sRVName,
	}
}

func (id SrvRecordId) String() string {
	segments := []string{
		fmt.Sprintf("Resource Group %q", id.ResourceGroup),
		fmt.Sprintf("Dnszone Name %q", id.DnszoneName),
		fmt.Sprintf("S R V Name %q", id.SRVName),
	}
	return strings.Join(segments, " / ")
}

func (id SrvRecordId) ID(_ string) string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Network/dnszones/%s/SRV/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroup, id.DnszoneName, id.SRVName)
}

// SrvRecordID parses a SrvRecord ID into an SrvRecordId struct
func SrvRecordID(input string) (*SrvRecordId, error) {
	id, err := azure.ParseAzureResourceID(input)
	if err != nil {
		return nil, err
	}

	resourceId := SrvRecordId{
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
	if resourceId.SRVName, err = id.PopSegment("SRV"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &resourceId, nil
}
