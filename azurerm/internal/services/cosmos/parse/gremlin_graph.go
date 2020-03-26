package parse

import (
	"fmt"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
)

type GremlinGraphId struct {
	ResourceGroup string
	Account       string
	Database      string
	Name          string
}

func GremlinGraphID(input string) (*GremlinGraphId, error) {
	id, err := azure.ParseAzureResourceID(input)
	if err != nil {
		return nil, fmt.Errorf("[ERROR] Unable to parse Gremlin Graph ID %q: %+v", input, err)
	}

	gremlinGraph := GremlinGraphId{
		ResourceGroup: id.ResourceGroup,
	}

	if gremlinGraph.Account, err = id.PopSegment("databaseAccounts"); err != nil {
		return nil, err
	}

	if gremlinGraph.Database, err = id.PopSegment("gremlinDatabases"); err != nil {
		return nil, err
	}

	if gremlinGraph.Name, err = id.PopSegment("graphs"); err != nil {
		return nil, err
	}

	return &gremlinGraph, nil
}
