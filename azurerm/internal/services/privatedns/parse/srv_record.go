package parse

// NOTE: this file is generated via 'go:generate' - manual changes will be overwritten

import (
	"fmt"
	"strings"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
)

type SrvRecordId struct {
	SubscriptionId     string
	ResourceGroup      string
	PrivateDnsZoneName string
	SRVName            string
}

func NewSrvRecordID(subscriptionId, resourceGroup, privateDnsZoneName, sRVName string) SrvRecordId {
	return SrvRecordId{
		SubscriptionId:     subscriptionId,
		ResourceGroup:      resourceGroup,
		PrivateDnsZoneName: privateDnsZoneName,
		SRVName:            sRVName,
	}
}

func (id SrvRecordId) String() string {
	segments := []string{
		fmt.Sprintf("S R V Name %q", id.SRVName),
		fmt.Sprintf("Private Dns Zone Name %q", id.PrivateDnsZoneName),
		fmt.Sprintf("Resource Group %q", id.ResourceGroup),
	}
	segmentsStr := strings.Join(segments, " / ")
	return fmt.Sprintf("%s: (%s)", "Srv Record", segmentsStr)
}

func (id SrvRecordId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Network/privateDnsZones/%s/SRV/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroup, id.PrivateDnsZoneName, id.SRVName)
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

	if resourceId.PrivateDnsZoneName, err = id.PopSegment("privateDnsZones"); err != nil {
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
