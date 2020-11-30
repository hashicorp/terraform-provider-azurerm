package parse

import (
	"fmt"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
)

type DatabaseId struct {
	ResourceGroup string
	ServerName    string
	Name          string
}

func NewDatabaseID(resourceGroup, serverName, name string) DatabaseId {
	return DatabaseId{
		ResourceGroup: resourceGroup,
		ServerName:    serverName,
		Name:          name,
	}
}

func (id DatabaseId) ID(subscriptionId string) string {
	return fmt.Sprintf("/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Sql/servers/%s/databases/%s", subscriptionId, id.ResourceGroup, id.ServerName, id.Name)
}

func DatabaseID(input string) (*DatabaseId, error) {
	id, err := azure.ParseAzureResourceID(input)
	if err != nil {
		return nil, fmt.Errorf("parsing Synapse Sql Pool ID %q: %+v", input, err)
	}

	sqlDatabaseId := DatabaseId{
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
