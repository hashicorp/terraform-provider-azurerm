package parse

// NOTE: this file is generated via 'go:generate' - manual changes will be overwritten

import (
	"fmt"
	"strings"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
)

type AaaaRecordId struct {
	SubscriptionId     string
	ResourceGroup      string
	PrivateDnsZoneName string
	AAAAName           string
}

func NewAaaaRecordID(subscriptionId, resourceGroup, privateDnsZoneName, aAAAName string) AaaaRecordId {
	return AaaaRecordId{
		SubscriptionId:     subscriptionId,
		ResourceGroup:      resourceGroup,
		PrivateDnsZoneName: privateDnsZoneName,
		AAAAName:           aAAAName,
	}
}

func (id AaaaRecordId) String() string {
	segments := []string{
		fmt.Sprintf("A A A A Name %q", id.AAAAName),
		fmt.Sprintf("Private Dns Zone Name %q", id.PrivateDnsZoneName),
		fmt.Sprintf("Resource Group %q", id.ResourceGroup),
	}
	segmentsStr := strings.Join(segments, " / ")
	return fmt.Sprintf("%s: (%s)", "Aaaa Record", segmentsStr)
}

func (id AaaaRecordId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Network/privateDnsZones/%s/AAAA/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroup, id.PrivateDnsZoneName, id.AAAAName)
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

	if resourceId.PrivateDnsZoneName, err = id.PopSegment("privateDnsZones"); err != nil {
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
