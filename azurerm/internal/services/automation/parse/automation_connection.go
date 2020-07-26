package parse

import (
	"fmt"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
)

type AutomationConnectionId struct {
	ResourceGroup string
	AccountName   string
	Name          string
}

func AutomationConnectionID(input string) (*AutomationConnectionId, error) {
	id, err := azure.ParseAzureResourceID(input)
	if err != nil {
		return nil, fmt.Errorf("parsing Automation Connection ID %q: %+v", input, err)
	}

	connection := AutomationConnectionId{
		ResourceGroup: id.ResourceGroup,
	}

	if connection.AccountName, err = id.PopSegment("automationAccounts"); err != nil {
		return nil, err
	}

	if connection.Name, err = id.PopSegment("connections"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &connection, nil
}
