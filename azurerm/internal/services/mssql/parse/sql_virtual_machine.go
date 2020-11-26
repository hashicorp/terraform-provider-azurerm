package parse

import (
	"fmt"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
)

type SqlVirtualMachineId struct {
	ResourceGroup string
	Name          string
}

func SqlVirtualMachineID(input string) (*SqlVirtualMachineId, error) {
	id, err := azure.ParseAzureResourceID(input)
	if err != nil {
		return nil, fmt.Errorf("[ERROR] Unable to parse Microsoft Sql VM ID %q: %+v", input, err)
	}

	sqlvm := SqlVirtualMachineId{
		ResourceGroup: id.ResourceGroup,
	}

	if sqlvm.Name, err = id.PopSegment("sqlVirtualMachines"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &sqlvm, nil
}
