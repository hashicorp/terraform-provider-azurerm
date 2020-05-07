package parse

import (
	"fmt"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
)

type VirtualMachineId struct {
	ResourceGroup string
	Name          string
}

func VirtualMachineID(input string) (*VirtualMachineId, error) {
	id, err := azure.ParseAzureResourceID(input)
	if err != nil {
		return nil, fmt.Errorf("[ERROR] Unable to parse VM ID %q: %+v", input, err)
	}

	vm := VirtualMachineId{
		ResourceGroup: id.ResourceGroup,
	}

	if vm.Name, err = id.PopSegment("virtualMachines"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &vm, nil
}
