package parse

import (
	"fmt"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
)

type DedicatedHardwareSecurityModuleId struct {
	ResourceGroup string
	Name          string
}

func DedicatedHardwareSecurityModuleID(input string) (*DedicatedHardwareSecurityModuleId, error) {
	id, err := azure.ParseAzureResourceID(input)
	if err != nil {
		return nil, fmt.Errorf("parsing DedicatedHardwareSecurityModule ID %q: %+v", input, err)
	}

	dedicatedHardwareSecurityModule := DedicatedHardwareSecurityModuleId{
		ResourceGroup: id.ResourceGroup,
	}

	if dedicatedHardwareSecurityModule.Name, err = id.PopSegment("dedicatedHSMs"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &dedicatedHardwareSecurityModule, nil
}
