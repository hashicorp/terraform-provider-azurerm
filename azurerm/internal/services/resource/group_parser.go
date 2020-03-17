package resource

import (
	"fmt"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
)

type ResourceGroupResourceID struct {
	Name string
}

func ParseResourceGroupID(input string) (*ResourceGroupResourceID, error) {
	id, err := azure.ParseAzureResourceID(input)
	if err != nil {
		return nil, fmt.Errorf("[ERROR] Unable to parse Resource Group ID %q: %+v", input, err)
	}

	group := ResourceGroupResourceID{
		Name: id.ResourceGroup,
	}

	if group.Name == "" {
		return nil, fmt.Errorf("ID contained no `resourceGroups` segment!")
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &group, nil
}
