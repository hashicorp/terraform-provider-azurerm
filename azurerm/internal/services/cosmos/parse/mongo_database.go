package parse

import (
	"fmt"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
)

type MongoDbDatabaseId struct {
	ResourceGroup string
	Account       string
	Name          string
}

func MongoDbDatabaseID(input string) (*MongoDbDatabaseId, error) {
	id, err := azure.ParseAzureResourceID(input)
	if err != nil {
		return nil, fmt.Errorf("[ERROR] Unable to parse MongoDb Database ID %q: %+v", input, err)
	}

	mongodbDatabase := MongoDbDatabaseId{
		ResourceGroup: id.ResourceGroup,
	}

	if mongodbDatabase.Account, err = id.PopSegment("databaseAccounts"); err != nil {
		return nil, err
	}

	if mongodbDatabase.Name, err = id.PopSegment("mongodbDatabases"); err != nil {
		return nil, err
	}

	return &mongodbDatabase, nil
}
