package parse

import (
	"fmt"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
)

type DataboxedgeDeviceId struct {
	ResourceGroup string
	Name          string
}

func NewDataboxEdgeDeviceID(resourcegroup string, name string) DataboxedgeDeviceId {
	return DataboxedgeDeviceId{
		ResourceGroup: resourcegroup,
		Name:          name,
	}
}

func (id DataboxedgeDeviceId) ID(subscriptionId string) string {
	return fmt.Sprintf("/subscriptions/%s/resourceGroups/%s/providers/Microsoft.DataBoxEdge/dataBoxEdgeDevices/%s", subscriptionId, id.ResourceGroup, id.Name)
}

func DataboxEdgeDeviceID(input string) (*DataboxedgeDeviceId, error) {
	id, err := azure.ParseAzureResourceID(input)
	if err != nil {
		return nil, fmt.Errorf("parsing Databox Edge Device ID %q: %+v", input, err)
	}

	databoxedgeDevice := DataboxedgeDeviceId{
		ResourceGroup: id.ResourceGroup,
	}
	if databoxedgeDevice.Name, err = id.PopSegment("dataBoxEdgeDevices"); err != nil {
		return nil, err
	}
	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &databoxedgeDevice, nil
}
