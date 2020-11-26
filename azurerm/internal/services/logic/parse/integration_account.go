package parse

import (
	"fmt"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
)

type IntegrationAccountId struct {
	ResourceGroup string
	Name          string
}

func IntegrationAccountID(input string) (*IntegrationAccountId, error) {
	id, err := azure.ParseAzureResourceID(input)
	if err != nil {
		return nil, fmt.Errorf("parsing Integration Account ID %q: %+v", input, err)
	}

	IntegrationAccount := IntegrationAccountId{
		ResourceGroup: id.ResourceGroup,
	}
	if IntegrationAccount.Name, err = id.PopSegment("integrationAccounts"); err != nil {
		return nil, err
	}
	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &IntegrationAccount, nil
}
