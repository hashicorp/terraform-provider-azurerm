package parse

import (
	"fmt"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
)

type LinkedServiceId struct {
	ResourceGroup string
	FactoryName   string
	Name          string
}

func LinkedServiceID(input string) (*LinkedServiceId, error) {
	id, err := azure.ParseAzureResourceID(input)
	if err != nil {
		return nil, fmt.Errorf("[ERROR] Unable to parse Data Factory Linked Service ID %q: %+v", input, err)
	}

	dataFactoryIntegrationRuntime := LinkedServiceId{
		ResourceGroup: id.ResourceGroup,
	}

	if dataFactoryIntegrationRuntime.FactoryName, err = id.PopSegment("factories"); err != nil {
		return nil, err
	}

	if dataFactoryIntegrationRuntime.Name, err = id.PopSegment("linkedservices"); err != nil {
		return nil, err
	}

	return &dataFactoryIntegrationRuntime, nil
}
