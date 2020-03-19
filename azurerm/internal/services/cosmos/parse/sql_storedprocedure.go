package parse

import (
	"fmt"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
)

type StoredProcedureId struct {
	ResourceGroup string
	Account       string
	Database      string
	Container     string
	Name          string
}

func StoredProcedureID(input string) (*StoredProcedureId, error) {
	id, err := azure.ParseAzureResourceID(input)
	if err != nil {
		return nil, fmt.Errorf("[ERROR] Unable to parse Stored Procedure ID %q: %+v", input, err)
	}

	storedProcedure := StoredProcedureId{
		ResourceGroup: id.ResourceGroup,
	}

	if storedProcedure.Account, err = id.PopSegment("databaseAccounts"); err != nil {
		return nil, err
	}

	if storedProcedure.Database, err = id.PopSegment("sqlDatabases"); err != nil {
		return nil, err
	}

	if storedProcedure.Container, err = id.PopSegment("containers"); err != nil {
		return nil, err
	}

	if storedProcedure.Name, err = id.PopSegment("storedProcedures"); err != nil {
		return nil, err
	}

	return &storedProcedure, nil
}
