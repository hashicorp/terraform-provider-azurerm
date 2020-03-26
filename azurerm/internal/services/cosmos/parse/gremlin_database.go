package parse

import (
	"fmt"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
)

type GremlinDatabaseId struct {
	ResourceGroup string
	Account       string
	Name          string
}

func GremlinDatabaseID(input string) (*GremlinDatabaseId, error) {
	id, err := azure.ParseAzureResourceID(input)
	if err != nil {
		return nil, fmt.Errorf("[ERROR] Unable to parse Gremlin Database ID %q: %+v", input, err)
	}

	gremlinDatabase := GremlinDatabaseId{
		ResourceGroup: id.ResourceGroup,
	}

	if gremlinDatabase.Account, err = id.PopSegment("databaseAccounts"); err != nil {
		return nil, err
	}

	if gremlinDatabase.Name, err = id.PopSegment("gremlinDatabases"); err != nil {
		return nil, err
	}

	return &gremlinDatabase, nil
}
