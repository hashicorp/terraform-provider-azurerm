package parse

// NOTE: this file is generated via 'go:generate' - manual changes will be overwritten

import (
	"fmt"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
)

type SqlStoredProcedureId struct {
	SubscriptionId      string
	ResourceGroup       string
	DatabaseAccountName string
	SqlDatabaseName     string
	ContainerName       string
	StoredProcedureName string
}

func NewSqlStoredProcedureID(subscriptionId, resourceGroup, databaseAccountName, sqlDatabaseName, containerName, storedProcedureName string) SqlStoredProcedureId {
	return SqlStoredProcedureId{
		SubscriptionId:      subscriptionId,
		ResourceGroup:       resourceGroup,
		DatabaseAccountName: databaseAccountName,
		SqlDatabaseName:     sqlDatabaseName,
		ContainerName:       containerName,
		StoredProcedureName: storedProcedureName,
	}
}

func (id SqlStoredProcedureId) ID(_ string) string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.DocumentDB/databaseAccounts/%s/sqlDatabases/%s/containers/%s/storedProcedures/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroup, id.DatabaseAccountName, id.SqlDatabaseName, id.ContainerName, id.StoredProcedureName)
}

func SqlStoredProcedureID(input string) (*SqlStoredProcedureId, error) {
	id, err := azure.ParseAzureResourceID(input)
	if err != nil {
		return nil, err
	}

	resourceId := SqlStoredProcedureId{
		SubscriptionId: id.SubscriptionID,
		ResourceGroup:  id.ResourceGroup,
	}

	if resourceId.DatabaseAccountName, err = id.PopSegment("databaseAccounts"); err != nil {
		return nil, err
	}
	if resourceId.SqlDatabaseName, err = id.PopSegment("sqlDatabases"); err != nil {
		return nil, err
	}
	if resourceId.ContainerName, err = id.PopSegment("containers"); err != nil {
		return nil, err
	}
	if resourceId.StoredProcedureName, err = id.PopSegment("storedProcedures"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &resourceId, nil
}
