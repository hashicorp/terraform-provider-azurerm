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

type MsSqlServerId struct {
	Name          string
	ResourceGroup string
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

func MsSqlServerID(input string) (*MsSqlServerId, error) {
	id, err := azure.ParseAzureResourceID(input)
	if err != nil {
		return nil, fmt.Errorf("Unable to parse MsSql Server ID %q: %+v", input, err)
	}

	server := MsSqlServerId{
		ResourceGroup: id.ResourceGroup,
	}

	if server.Name, err = id.PopSegment("servers"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &server, nil
}

func ParseArmMSSqlElasticPoolId(sqlElasticPoolId string) (string, string, string, error) {
	id, err := azure.ParseAzureResourceID(sqlElasticPoolId)
	if err != nil {
		return "", "", "", fmt.Errorf("Unable to parse SQL ElasticPool ID %q: %+v", sqlElasticPoolId, err)
	}

	return id.ResourceGroup, id.Path["servers"], id.Path["elasticPools"], nil
}
