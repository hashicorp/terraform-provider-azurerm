package parse

import (
	"fmt"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
)

type MsSqlDatabaseId struct {
	Name          string
	MsSqlServer   string
	ResourceGroup string
}

func NewMsSqlDatabaseID(resourceGroup, msSqlServer, name string) MsSqlDatabaseId {
	return MsSqlDatabaseId{
		ResourceGroup: resourceGroup,
		MsSqlServer:   msSqlServer,
		Name:          name,
	}
}

func (id MsSqlDatabaseId) ID(subscriptionId string) string {
	return fmt.Sprintf("/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Sql/servers/%s/databases/%s", subscriptionId, id.ResourceGroup, id.MsSqlServer, id.Name)
}

func MsSqlDatabaseID(input string) (*MsSqlDatabaseId, error) {
	id, err := azure.ParseAzureResourceID(input)
	if err != nil {
		return nil, fmt.Errorf("Unable to parse MsSql Database ID %q: %+v", input, err)
	}

	database := MsSqlDatabaseId{
		ResourceGroup: id.ResourceGroup,
	}

	if database.MsSqlServer, err = id.PopSegment("servers"); err != nil {
		return nil, err
	}

	if database.Name, err = id.PopSegment("databases"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &database, nil
}
