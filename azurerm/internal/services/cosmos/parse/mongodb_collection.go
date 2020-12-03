package parse

// NOTE: this file is generated via 'go:generate' - manual changes will be overwritten

import (
	"fmt"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
)

type MongodbCollectionId struct {
	SubscriptionId      string
	ResourceGroup       string
	DatabaseAccountName string
	MongodbDatabaseName string
	CollectionName      string
}

func NewMongodbCollectionID(subscriptionId, resourceGroup, databaseAccountName, mongodbDatabaseName, collectionName string) MongodbCollectionId {
	return MongodbCollectionId{
		SubscriptionId:      subscriptionId,
		ResourceGroup:       resourceGroup,
		DatabaseAccountName: databaseAccountName,
		MongodbDatabaseName: mongodbDatabaseName,
		CollectionName:      collectionName,
	}
}

func (id MongodbCollectionId) ID(_ string) string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.DocumentDB/databaseAccounts/%s/mongodbDatabases/%s/collections/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroup, id.DatabaseAccountName, id.MongodbDatabaseName, id.CollectionName)
}

// MongodbCollectionID parses a MongodbCollection ID into an MongodbCollectionId struct
func MongodbCollectionID(input string) (*MongodbCollectionId, error) {
	id, err := azure.ParseAzureResourceID(input)
	if err != nil {
		return nil, err
	}

	resourceId := MongodbCollectionId{
		SubscriptionId: id.SubscriptionID,
		ResourceGroup:  id.ResourceGroup,
	}

	if resourceId.SubscriptionId == "" {
		return nil, fmt.Errorf("ID was missing the 'subscriptions' element")
	}

	if resourceId.ResourceGroup == "" {
		return nil, fmt.Errorf("ID was missing the 'resourceGroups' element")
	}

	if resourceId.DatabaseAccountName, err = id.PopSegment("databaseAccounts"); err != nil {
		return nil, err
	}
	if resourceId.MongodbDatabaseName, err = id.PopSegment("mongodbDatabases"); err != nil {
		return nil, err
	}
	if resourceId.CollectionName, err = id.PopSegment("collections"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &resourceId, nil
}
