package compute

import (
	"fmt"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
)

type VirtualMachineExtensionID struct {
	ResourceGroup  string
	Name           string
	VirtualMachine string
}

func ParseVirtualMachineExtensionID(input string) (*VirtualMachineExtensionID, error) {
	id, err := azure.ParseAzureResourceID(input)
	if err != nil {
		return nil, fmt.Errorf("[ERROR] Unable to parse App Service ID %q: %+v", input, err)
	}

	virtualMachineExtension := VirtualMachineExtensionID{
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

func ValidateVirtualMachineExtensionID(i interface{}, k string) (warnings []string, errors []error) {
	v, ok := i.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected type of %q to be string", k))
		return
	}

	if _, err := ParseVirtualMachineExtensionID(v); err != nil {
		errors = append(errors, fmt.Errorf("Can not parse %q as a resource id: %v", k, err))
		return
	}

	return warnings, errors
}
