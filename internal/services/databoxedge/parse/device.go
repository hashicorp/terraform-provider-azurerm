package parse

// NOTE: this file is generated via 'go:generate' - manual changes will be overwritten

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

type DeviceId struct {
	SubscriptionId        string
	ResourceGroup         string
	DataBoxEdgeDeviceName string
}

func NewDeviceID(subscriptionId, resourceGroup, dataBoxEdgeDeviceName string) DeviceId {
	return DeviceId{
		SubscriptionId:        subscriptionId,
		ResourceGroup:         resourceGroup,
		DataBoxEdgeDeviceName: dataBoxEdgeDeviceName,
	}
}

func (id DeviceId) String() string {
	segments := []string{
		fmt.Sprintf("Data Box Edge Device Name %q", id.DataBoxEdgeDeviceName),
		fmt.Sprintf("Resource Group %q", id.ResourceGroup),
	}
	segmentsStr := strings.Join(segments, " / ")
	return fmt.Sprintf("%s: (%s)", "Device", segmentsStr)
}

func (id DeviceId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.DataBoxEdge/dataBoxEdgeDevices/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroup, id.DataBoxEdgeDeviceName)
}

// DeviceID parses a Device ID into an DeviceId struct
func DeviceID(input string) (*DeviceId, error) {
	id, err := resourceids.ParseAzureResourceID(input)
	if err != nil {
		return nil, err
	}

	resourceId := DeviceId{
		SubscriptionId: id.SubscriptionID,
		ResourceGroup:  id.ResourceGroup,
	}

	if resourceId.SubscriptionId == "" {
		return nil, fmt.Errorf("ID was missing the 'subscriptions' element")
	}

	if resourceId.ResourceGroup == "" {
		return nil, fmt.Errorf("ID was missing the 'resourceGroups' element")
	}

	if resourceId.DataBoxEdgeDeviceName, err = id.PopSegment("dataBoxEdgeDevices"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &resourceId, nil
}
