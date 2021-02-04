package parse

import (
	"fmt"
	"strings"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
)

type IotSecurityDeviceGroupId struct {
	IotHubID string
	Name     string
}

func NewIotSecurityDeviceGroupId(iotHubID string, name string) IotSecurityDeviceGroupId {
	return IotSecurityDeviceGroupId{
		IotHubID: iotHubID,
		Name:     name,
	}
}

func (id IotSecurityDeviceGroupId) String() string {
	segments := []string{
		fmt.Sprintf("Name %q", id.Name),
		fmt.Sprintf("IotHub Id %q", id.IotHubID),
	}
	segmentsStr := strings.Join(segments, " / ")
	return fmt.Sprintf("%s: (%s)", "Iot Security Device Group", segmentsStr)
}

func (id IotSecurityDeviceGroupId) ID() string {
	fmtString := "%s/providers/Microsoft.Security/deviceSecurityGroups/%s"
	return fmt.Sprintf(fmtString, id.IotHubID, id.Name)
}

func IotSecurityDeviceGroupID(input string) (*IotSecurityDeviceGroupId, error) {
	id, err := azure.ParseAzureResourceID(input)
	if err != nil {
		return nil, fmt.Errorf("parsing Iot Security Device Group ID %q: %+v", input, err)
	}

	parts := strings.Split(input, "/providers/Microsoft.Security/deviceSecurityGroups/")
	if len(parts) != 2 {
		return nil, fmt.Errorf("parsing Iot Security Device Group ID: %q", id)
	}

	return &IotSecurityDeviceGroupId{
		IotHubID: parts[0],
		Name:     parts[1],
	}, nil
}
