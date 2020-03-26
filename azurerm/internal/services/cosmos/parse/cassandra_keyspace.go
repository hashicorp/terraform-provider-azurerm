package parse

import (
	"fmt"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
)

type CassandraKeyspaceId struct {
	ResourceGroup string
	Account       string
	Name          string
}

func CassandraKeyspaceID(input string) (*CassandraKeyspaceId, error) {
	id, err := azure.ParseAzureResourceID(input)
	if err != nil {
		return nil, fmt.Errorf("[ERROR] Unable to parse Cassandra Keyspace ID %q: %+v", input, err)
	}

	cassandraKeyspace := CassandraKeyspaceId{
		ResourceGroup: id.ResourceGroup,
	}

	if cassandraKeyspace.Account, err = id.PopSegment("databaseAccounts"); err != nil {
		return nil, err
	}

	if cassandraKeyspace.Name, err = id.PopSegment("cassandraKeyspaces"); err != nil {
		return nil, err
	}

	return &cassandraKeyspace, nil
}
