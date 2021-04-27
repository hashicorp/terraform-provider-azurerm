package parse

// NOTE: this file is generated via 'go:generate' - manual changes will be overwritten

import (
	"fmt"
	"strings"

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

func (id AaaaRecordId) String() string {
	segments := []string{
		fmt.Sprintf("A A A A Name %q", id.AAAAName),
		fmt.Sprintf("Dnszone Name %q", id.DnszoneName),
		fmt.Sprintf("Resource Group %q", id.ResourceGroup),
	}
	segmentsStr := strings.Join(segments, " / ")
	return fmt.Sprintf("%s: (%s)", "Aaaa Record", segmentsStr)
}

func (id AaaaRecordId) ID() string {
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

	if resourceId.SubscriptionId == "" {
		return nil, fmt.Errorf("ID was missing the 'subscriptions' element")
	}

	if resourceId.ResourceGroup == "" {
		return nil, fmt.Errorf("ID was missing the 'resourceGroups' element")
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
