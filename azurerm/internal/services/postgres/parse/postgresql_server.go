package parse

import (
	"fmt"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
)

type PostgresqlServerId struct {
	ResourceGroup string
	Name          string
}

func PostgresqlServerID(input string) (*PostgresqlServerId, error) {
	id, err := azure.ParseAzureResourceID(input)
	if err != nil {
		return nil, fmt.Errorf("[ERROR] Unable to parse Postgresql Server ID %q: %+v", input, err)
	}

	set := PostgresqlServerId{
		ResourceGroup: id.ResourceGroup,
	}

	if set.Name, err = id.PopSegment("servers"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &set, nil
}
