package parse

import "github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"

type ApiVersionSetId struct {
	ResourceGroup string
	ServiceName   string
	Name          string
}

func APIVersionSetID(input string) (*ApiVersionSetId, error) {
	id, err := azure.ParseAzureResourceID(input)
	if err != nil {
		return nil, err
	}

	logger := ApiVersionSetId{
		ResourceGroup: id.ResourceGroup,
	}

	if logger.ServiceName, err = id.PopSegment("service"); err != nil {
		return nil, err
	}

	if logger.Name, err = id.PopSegment("apiVersionSets"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &logger, nil
}
