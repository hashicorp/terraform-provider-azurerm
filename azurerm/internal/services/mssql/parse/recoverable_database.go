package parse

import (
	"fmt"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
)

type RecoverableDBId struct {
	Name          string
	ServerName    string
	ResourceGroup string
}

func RecoverableDBID(input string) (*RecoverableDBId, error) {
	id, err := azure.ParseAzureResourceID(input)
	if err != nil {
		return nil, fmt.Errorf("[ERROR] Unable to parse Microsoft Sql Recoverable DB ID %q: %+v", input, err)
	}

	recoverableDBId := RecoverableDBId{
		ResourceGroup: id.ResourceGroup,
	}

	if recoverableDBId.ServerName, err = id.PopSegment("servers"); err != nil {
		return nil, err
	}

	if recoverableDBId.Name, err = id.PopSegment("recoverabledatabases"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &recoverableDBId, nil
}
