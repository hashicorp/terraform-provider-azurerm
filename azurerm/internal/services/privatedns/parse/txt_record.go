package parse

// NOTE: this file is generated via 'go:generate' - manual changes will be overwritten

import (
	"fmt"
	"strings"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
)

type TxtRecordId struct {
	SubscriptionId     string
	ResourceGroup      string
	PrivateDnsZoneName string
	TXTName            string
}

func NewTxtRecordID(subscriptionId, resourceGroup, privateDnsZoneName, tXTName string) TxtRecordId {
	return TxtRecordId{
		SubscriptionId:     subscriptionId,
		ResourceGroup:      resourceGroup,
		PrivateDnsZoneName: privateDnsZoneName,
		TXTName:            tXTName,
	}
}

func (id TxtRecordId) String() string {
	segments := []string{
		fmt.Sprintf("T X T Name %q", id.TXTName),
		fmt.Sprintf("Private Dns Zone Name %q", id.PrivateDnsZoneName),
		fmt.Sprintf("Resource Group %q", id.ResourceGroup),
	}
	segmentsStr := strings.Join(segments, " / ")
	return fmt.Sprintf("%s: (%s)", "Txt Record", segmentsStr)
}

func (id TxtRecordId) ID(_ string) string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Network/privateDnsZones/%s/TXT/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroup, id.PrivateDnsZoneName, id.TXTName)
}

// TxtRecordID parses a TxtRecord ID into an TxtRecordId struct
func TxtRecordID(input string) (*TxtRecordId, error) {
	id, err := azure.ParseAzureResourceID(input)
	if err != nil {
		return nil, err
	}

	resourceId := TxtRecordId{
		SubscriptionId: id.SubscriptionID,
		ResourceGroup:  id.ResourceGroup,
	}

	if resourceId.SubscriptionId == "" {
		return nil, fmt.Errorf("ID was missing the 'subscriptions' element")
	}

	if resourceId.ResourceGroup == "" {
		return nil, fmt.Errorf("ID was missing the 'resourceGroups' element")
	}

	if resourceId.PrivateDnsZoneName, err = id.PopSegment("privateDnsZones"); err != nil {
		return nil, err
	}
	if resourceId.TXTName, err = id.PopSegment("TXT"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &resourceId, nil
}
