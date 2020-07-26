package parse

import (
	"fmt"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
)

type DataFactoryIntegrationRuntimeId struct {
	ResourceGroup string
	Name          string
	DataFactory   string
}

func DataFactoryIntegrationRuntimeID(input string) (*DataFactoryIntegrationRuntimeId, error) {
	id, err := azure.ParseAzureResourceID(input)
	if err != nil {
		return nil, fmt.Errorf("[ERROR] Unable to parse Data Factory Integration Runtime ID %q: %+v", input, err)
	}

	dataFactoryIntegrationRuntime := DataFactoryIntegrationRuntimeId{
		ResourceGroup: id.ResourceGroup,
	}

	if dataFactoryIntegrationRuntime.DataFactory, err = id.PopSegment("factories"); err != nil {
		return nil, err
	}

	if dataFactoryIntegrationRuntime.Name, err = id.PopSegment("integrationruntimes"); err != nil {
		return nil, err
	}

	return &dataFactoryIntegrationRuntime, nil
}
