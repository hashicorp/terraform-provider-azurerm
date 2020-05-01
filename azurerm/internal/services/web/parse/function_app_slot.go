package parse

import (
	"fmt"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
)

type FunctionAppSlotResourceID struct {
	ResourceGroup   string
	FunctionAppName string
	Name            string
}

func FunctionAppSlotID(input string) (*FunctionAppSlotResourceID, error) {
	id, err := azure.ParseAzureResourceID(input)
	if err != nil {
		return nil, fmt.Errorf("[ERROR] Unable to parse App Service Slot ID %q: %+v", input, err)
	}

	slot := FunctionAppSlotResourceID{
		ResourceGroup:   id.ResourceGroup,
		FunctionAppName: id.Path["sites"],
		Name:            id.Path["slots"],
	}

	if slot.FunctionAppName, err = id.PopSegment("sites"); err != nil {
		return nil, err
	}

	if slot.Name, err = id.PopSegment("slots"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &slot, nil
}
