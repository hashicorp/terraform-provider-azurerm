package parse

import (
	"fmt"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
)

type DataboxEdgeOrderId struct {
	ResourceGroup string
	DeviceName    string
}

func NewDataboxEdgeOrderID(resourcegroup string, deviceName string) DataboxEdgeOrderId {
	return DataboxEdgeOrderId{
		ResourceGroup: resourcegroup,
		DeviceName:    deviceName,
	}
}

func (id DataboxEdgeOrderId) ID(subscriptionId string) string {
	return fmt.Sprintf("/subscriptions/%s/resourceGroups/%s/providers/Microsoft.DataBoxEdge/dataBoxEdgeDevices/%s/orders/default", subscriptionId, id.ResourceGroup, id.DeviceName)
}

func DataboxEdgeOrderID(input string) (*DataboxEdgeOrderId, error) {
	id, err := azure.ParseAzureResourceID(input)
	if err != nil {
		return nil, fmt.Errorf("parsing Databox Edge Order ID %q: %+v", input, err)
	}

	databoxedgeOrder := DataboxEdgeOrderId{
		ResourceGroup: id.ResourceGroup,
	}
	if databoxedgeOrder.DeviceName, err = id.PopSegment("dataBoxEdgeDevices"); err != nil {
		return nil, err
	}
	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &databoxedgeOrder, nil
}
