package parse

import (
	"fmt"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
)

type ConnectionId struct {
	ResourceGroup         string
	AutomationAccountName string
	ConnectionName        string
}

func ConnectionID(input string) (*ConnectionId, error) {
	id, err := azure.ParseAzureResourceID(input)
	if err != nil {
		return nil, fmt.Errorf("parsing Automation Connection ID %q: %+v", input, err)
	}

	connection := ConnectionId{
		ResourceGroup: id.ResourceGroup,
	}

	if connection.AutomationAccountName, err = id.PopSegment("automationAccounts"); err != nil {
		return nil, err
	}

	if connection.ConnectionName, err = id.PopSegment("connections"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &connection, nil
}
