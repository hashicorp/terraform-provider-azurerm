package parse

// NOTE: this file is generated via 'go:generate' - manual changes will be overwritten

import (
	"fmt"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
)

type SqlContainerId struct {
	SubscriptionId      string
	ResourceGroup       string
	DatabaseAccountName string
	SqlDatabaseName     string
	ContainerName       string
}

func NewSqlContainerID(subscriptionId, resourceGroup, databaseAccountName, sqlDatabaseName, containerName string) SqlContainerId {
	return SqlContainerId{
		SubscriptionId:      subscriptionId,
		ResourceGroup:       resourceGroup,
		DatabaseAccountName: databaseAccountName,
		SqlDatabaseName:     sqlDatabaseName,
		ContainerName:       containerName,
	}
}

func (id SqlContainerId) ID(_ string) string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.DocumentDB/databaseAccounts/%s/sqlDatabases/%s/containers/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroup, id.DatabaseAccountName, id.SqlDatabaseName, id.ContainerName)
}

// SqlContainerID parses a SqlContainer ID into an SqlContainerId struct
func SqlContainerID(input string) (*SqlContainerId, error) {
	id, err := azure.ParseAzureResourceID(input)
	if err != nil {
		return nil, err
	}

	resourceId := SqlContainerId{
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

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &resourceId, nil
}
