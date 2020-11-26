package parse

// NOTE: this file is generated via 'go:generate' - manual changes will be overwritten

import (
	"fmt"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
)

type AaaaRecordId struct {
	SubscriptionId string
	ResourceGroup  string
	DnszoneName    string
	AAAAName       string
}

func NewAaaaRecordID(subscriptionId, resourceGroup, dnszoneName, aAAAName string) AaaaRecordId {
	return AaaaRecordId{
		SubscriptionId: subscriptionId,
		ResourceGroup:  resourceGroup,
		DnszoneName:    dnszoneName,
		AAAAName:       aAAAName,
	}
}

func (id AaaaRecordId) ID(_ string) string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Network/dnszones/%s/AAAA/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroup, id.DnszoneName, id.AAAAName)
}

// AaaaRecordID parses a AaaaRecord ID into an AaaaRecordId struct
func AaaaRecordID(input string) (*AaaaRecordId, error) {
	id, err := azure.ParseAzureResourceID(input)
	if err != nil {
		return nil, err
	}

	resourceId := AaaaRecordId{
		SubscriptionId: id.SubscriptionID,
		ResourceGroup:  id.ResourceGroup,
	}

	if resourceId.DnszoneName, err = id.PopSegment("dnszones"); err != nil {
		return nil, err
	}
	if resourceId.AAAAName, err = id.PopSegment("AAAA"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &resourceId, nil
}
