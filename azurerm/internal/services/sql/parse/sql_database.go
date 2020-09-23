package parse

import (
	"fmt"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
)

type SqlDatabaseId struct {
	ResourceGroup string
	ServerName    string
	Name          string
}

func NewSqlDatabaseID(resourceGroup, serverName, name string) SqlDatabaseId {
	return SqlDatabaseId{
		ResourceGroup: resourceGroup,
		ServerName:    serverName,
		Name:          name,
	}
}

func (id SqlDatabaseId) ID(subscriptionId string) string {
	return fmt.Sprintf("/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Sql/servers/%s/databases/%s", subscriptionId, id.ResourceGroup, id.ServerName, id.Name)
}

func SqlDatabaseID(input string) (*SqlDatabaseId, error) {
	id, err := azure.ParseAzureResourceID(input)
	if err != nil {
		return nil, fmt.Errorf("parsing Synapse Sql Pool ID %q: %+v", input, err)
	}

	sqlDatabaseId := SqlDatabaseId{
		ResourceGroup: id.ResourceGroup,
	}
	if sqlDatabaseId.ServerName, err = id.PopSegment("servers"); err != nil {
		return nil, err
	}
	if sqlDatabaseId.Name, err = id.PopSegment("databases"); err != nil {
		return nil, err
	}
	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &sqlDatabaseId, nil
}
