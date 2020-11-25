package parse

// NOTE: this file is generated via 'go:generate' - manual changes will be overwritten

import (
	"fmt"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
)

type PtrRecordId struct {
	SubscriptionId string
	ResourceGroup  string
	DnszoneName    string
	PTRName        string
}

func NewPtrRecordID(subscriptionId, resourceGroup, dnszoneName, pTRName string) PtrRecordId {
	return PtrRecordId{
		SubscriptionId: subscriptionId,
		ResourceGroup:  resourceGroup,
		DnszoneName:    dnszoneName,
		PTRName:        pTRName,
	}
}

func (id PtrRecordId) ID(_ string) string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Network/dnszones/%s/PTR/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroup, id.DnszoneName, id.PTRName)
}

func PtrRecordID(input string) (*PtrRecordId, error) {
	id, err := azure.ParseAzureResourceID(input)
	if err != nil {
		return nil, err
	}

	resourceId := PtrRecordId{
		SubscriptionId: id.SubscriptionID,
		ResourceGroup:  id.ResourceGroup,
	}

	if resourceId.DnszoneName, err = id.PopSegment("dnszones"); err != nil {
		return nil, err
	}
	if resourceId.PTRName, err = id.PopSegment("PTR"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &resourceId, nil
}
