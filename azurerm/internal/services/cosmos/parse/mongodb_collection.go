package parse

import (
	"fmt"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
)

type MongoDbCollectionId struct {
	ResourceGroup       string
	DatabaseAccountName string
	MongodbDatabaseName string
	CollectionName      string
}

func MongoDbCollectionID(input string) (*MongoDbCollectionId, error) {
	id, err := azure.ParseAzureResourceID(input)
	if err != nil {
		return nil, fmt.Errorf("[ERROR] Unable to parse MongoDb Collection ID %q: %+v", input, err)
	}

	mongodbCollection := MongoDbCollectionId{
		ResourceGroup: id.ResourceGroup,
	}

	if mongodbCollection.DatabaseAccountName, err = id.PopSegment("databaseAccounts"); err != nil {
		return nil, err
	}

	if mongodbCollection.MongodbDatabaseName, err = id.PopSegment("mongodbDatabases"); err != nil {
		return nil, err
	}

	if mongodbCollection.CollectionName, err = id.PopSegment("collections"); err != nil {
		return nil, err
	}

	return &mongodbCollection, nil
}
