package parse

import (
	"fmt"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
)

type IntegrationRuntimeId struct {
	ResourceGroup string
	FactoryName   string
	Name          string
}

func IntegrationRuntimeID(input string) (*IntegrationRuntimeId, error) {
	id, err := azure.ParseAzureResourceID(input)
	if err != nil {
		return nil, fmt.Errorf("[ERROR] Unable to parse Data Factory Integration Runtime ID %q: %+v", input, err)
	}

	dataFactoryIntegrationRuntime := IntegrationRuntimeId{
		ResourceGroup: id.ResourceGroup,
	}

	if dataFactoryIntegrationRuntime.FactoryName, err = id.PopSegment("factories"); err != nil {
		return nil, err
	}

	if dataFactoryIntegrationRuntime.Name, err = id.PopSegment("integrationruntimes"); err != nil {
		return nil, err
	}

	return &dataFactoryIntegrationRuntime, nil
}
