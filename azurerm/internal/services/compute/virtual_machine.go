package compute

import (
	"fmt"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
)

type VirtualMachineID struct {
	ResourceGroup string
	Name          string
}

func ParseVirtualMachineID(input string) (*VirtualMachineID, error) {
	id, err := azure.ParseAzureResourceID(input)
	if err != nil {
		return nil, fmt.Errorf("[ERROR] Unable to parse Virtual Machine ID %q: %+v", input, err)
	}

	virtualMachine := VirtualMachineID{
		ResourceGroup: id.ResourceGroup,
	}

	if virtualMachine.Name, err = id.PopSegment("virtualMachines"); err != nil {
		return nil, err
	}

	return &virtualMachine, nil
}

func ValidateVirtualMachineID(i interface{}, k string) (warnings []string, errors []error) {
	v, ok := i.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected type of %q to be string", k))
		return
	}

	if _, err := ParseVirtualMachineID(v); err != nil {
		errors = append(errors, fmt.Errorf("Can not parse %q as a resource id: %v", k, err))
		return
	}

	return warnings, errors
}
