package parse

import (
	"fmt"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
)

type VirtualMachineExtensionId struct {
	ResourceGroup  string
	Name           string
	VirtualMachine string
}

func VirtualMachineExtensionID(input string) (*VirtualMachineExtensionId, error) {
	id, err := azure.ParseAzureResourceID(input)
	if err != nil {
		return nil, fmt.Errorf("[ERROR] Unable to parse App Service ID %q: %+v", input, err)
	}

	virtualMachineExtension := VirtualMachineExtensionId{
		ResourceGroup: id.ResourceGroup,
	}

	if virtualMachineExtension.VirtualMachine, err = id.PopSegment("virtualMachines"); err != nil {
		return nil, err
	}

	if virtualMachineExtension.Name, err = id.PopSegment("extensions"); err != nil {
		return nil, err
	}

	return &virtualMachineExtension, nil
}
