package parse

import (
	"fmt"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
)

type ElasticPoolId struct {
	Name          string
	ServerName    string
	ResourceGroup string
}

func ElasticPoolID(input string) (*ElasticPoolId, error) {
	id, err := azure.ParseAzureResourceID(input)
	if err != nil {
		return nil, fmt.Errorf("Unable to parse MsSql Elastic Pool ID %q: %+v", input, err)
	}

	elasticPool := ElasticPoolId{
		ResourceGroup: id.ResourceGroup,
	}

	if elasticPool.ServerName, err = id.PopSegment("servers"); err != nil {
		return nil, err
	}

	if elasticPool.Name, err = id.PopSegment("elasticPools"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &elasticPool, nil
}
