package parse

import (
	"fmt"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
)

type SqlDatabaseId struct {
	SubscriptionID string
	ResourceGroup  string
	ServerName     string
	Name           string
}

func SqlDatabaseID(input string) (*SqlDatabaseId, error) {
	id, err := azure.ParseAzureResourceID(input)
	if err != nil {
		return nil, fmt.Errorf("parsing Synapse Sql Pool ID %q: %+v", input, err)
	}

	sqlDatabaseId := SqlDatabaseId{
		SubscriptionID: id.SubscriptionID,
		ResourceGroup:  id.ResourceGroup,
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

func (id *SqlDatabaseId) String() string {
	return fmt.Sprintf("/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Sql/servers/%s/databases/%s", id.SubscriptionID, id.ResourceGroup, id.ServerName, id.Name)
}
