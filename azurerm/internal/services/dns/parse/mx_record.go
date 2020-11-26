package parse

// NOTE: this file is generated via 'go:generate' - manual changes will be overwritten

import (
	"fmt"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
)

type MxRecordId struct {
	SubscriptionId string
	ResourceGroup  string
	DnszoneName    string
	MXName         string
}

func NewMxRecordID(subscriptionId, resourceGroup, dnszoneName, mXName string) MxRecordId {
	return MxRecordId{
		SubscriptionId: subscriptionId,
		ResourceGroup:  resourceGroup,
		DnszoneName:    dnszoneName,
		MXName:         mXName,
	}
}

func (id MxRecordId) ID(_ string) string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Network/dnszones/%s/MX/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroup, id.DnszoneName, id.MXName)
}

// MxRecordID parses a MxRecord ID into an MxRecordId struct
func MxRecordID(input string) (*MxRecordId, error) {
	id, err := azure.ParseAzureResourceID(input)
	if err != nil {
		return nil, err
	}

	resourceId := MxRecordId{
		SubscriptionId: id.SubscriptionID,
		ResourceGroup:  id.ResourceGroup,
	}

	if resourceId.DnszoneName, err = id.PopSegment("dnszones"); err != nil {
		return nil, err
	}
	if resourceId.MXName, err = id.PopSegment("MX"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &resourceId, nil
}
