package parse

import (
	"fmt"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
)

type FunctionID struct {
	ResourceGroup   string
	FunctionAppName string
	Name            string
}

func ParseFunctionID(input string) (*FunctionID, error) {
	id, err := azure.ParseAzureResourceID(input)
	if err != nil {
		return nil, fmt.Errorf("[ERROR] Unable to parse Function ID %q: %+v", input, err)
	}

	function := FunctionID{
		ResourceGroup:   id.ResourceGroup,
		FunctionAppName: id.Path["sites"],
		Name:            id.Path["functions"],
	}

	if function.FunctionAppName, err = id.PopSegment("sites"); err != nil {
		return nil, err
	}

	if function.Name, err = id.PopSegment("functions"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &function, nil
}
