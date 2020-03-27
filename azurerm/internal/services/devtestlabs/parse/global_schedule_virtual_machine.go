package parse

import (
	"fmt"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
)

type GlobalScheduleVirtualMachineId struct {
	ResourceGroup string
	Name          string
}

func GlobalScheduleVirtualMachineID(input string) (*GlobalScheduleVirtualMachineId, error) {
	id, err := azure.ParseAzureResourceID(input)
	if err != nil {
		return nil, fmt.Errorf("[ERROR] Unable to parse Virtual Machine ID %q: %+v", input, err)
	}

	service := GlobalScheduleVirtualMachineId{
		ResourceGroup: id.ResourceGroup,
	}

	if service.Name, err = id.PopSegment("virtualMachines"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &service, nil
}
