package parse

import (
	"fmt"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
)

type FlexibleServerId struct {
	ResourceGroup string
	Name          string
}

func FlexibleServerID(input string) (*FlexibleServerId, error) {
	id, err := azure.ParseAzureResourceID(input)
	if err != nil {
		return nil, fmt.Errorf("parsing postgresqlflexible Server ID %q: %+v", input, err)
	}

	flexibleServer := FlexibleServerId{
		ResourceGroup: id.ResourceGroup,
	}
	if flexibleServer.Name, err = id.PopSegment("flexibleServers"); err != nil {
		return nil, err
	}
	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &flexibleServer, nil
}
