package parse

// NOTE: this file is generated via 'go:generate' - manual changes will be overwritten

import (
	"fmt"
	"strings"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
)

type MxRecordId struct {
	SubscriptionId     string
	ResourceGroup      string
	PrivateDnsZoneName string
	MXName             string
}

func NewMxRecordID(subscriptionId, resourceGroup, privateDnsZoneName, mXName string) MxRecordId {
	return MxRecordId{
		SubscriptionId:     subscriptionId,
		ResourceGroup:      resourceGroup,
		PrivateDnsZoneName: privateDnsZoneName,
		MXName:             mXName,
	}
}

func (id MxRecordId) String() string {
	segments := []string{
		fmt.Sprintf("M X Name %q", id.MXName),
		fmt.Sprintf("Private Dns Zone Name %q", id.PrivateDnsZoneName),
		fmt.Sprintf("Resource Group %q", id.ResourceGroup),
	}
	segmentsStr := strings.Join(segments, " / ")
	return fmt.Sprintf("%s: (%s)", "Mx Record", segmentsStr)
}

func (id MxRecordId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Network/privateDnsZones/%s/MX/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroup, id.PrivateDnsZoneName, id.MXName)
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

	if resourceId.SubscriptionId == "" {
		return nil, fmt.Errorf("ID was missing the 'subscriptions' element")
	}

	if resourceId.ResourceGroup == "" {
		return nil, fmt.Errorf("ID was missing the 'resourceGroups' element")
	}

	if resourceId.PrivateDnsZoneName, err = id.PopSegment("privateDnsZones"); err != nil {
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
