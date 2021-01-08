package parse

import (
	"fmt"
	"strings"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
)

type DeviceSecurityGroupId struct {
	TargetResourceID string
	Name             string
}

func NewDeviceSecurityGroupId(targetResourceId string, name string) DeviceSecurityGroupId {
	return DeviceSecurityGroupId{
		TargetResourceID: targetResourceId,
		Name:             name,
	}
}

func (id DeviceSecurityGroupId) ID() string {
	fmtString := "%s/providers/Microsoft.Security/deviceSecurityGroups/%s"
	return fmt.Sprintf(fmtString, id.TargetResourceID, id.Name)
}

func DeviceSecurityGroupID(input string) (*DeviceSecurityGroupId, error) {
	id, err := azure.ParseAzureResourceID(input)
	if err != nil {
		return nil, fmt.Errorf("parsing Device Security Group ID %q: %+v", input, err)
	}

	parts := strings.Split(input, "/providers/Microsoft.Security/deviceSecurityGroups/")
	if len(parts) != 2 {
		return nil, fmt.Errorf("parsing Device Security Group ID: %q", id)
	}

	return &DeviceSecurityGroupId{
		TargetResourceID: parts[0],
		Name:             parts[1],
	}, nil
}
