package parse

import (
	"fmt"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
)

type ApplicationId struct {
	ResourceGroup    string
	BatchAccountName string
	ApplicationName  string
}

func ApplicationID(input string) (*ApplicationId, error) {
	id, err := azure.ParseAzureResourceID(input)
	if err != nil {
		return nil, fmt.Errorf("[ERROR] Unable to parse Batch Application ID %q: %+v", input, err)
	}

	application := ApplicationId{
		ResourceGroup: id.ResourceGroup,
	}

	if application.BatchAccountName, err = id.PopSegment("batchAccounts"); err != nil {
		return nil, err
	}

	if application.ApplicationName, err = id.PopSegment("applications"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &application, nil
}
