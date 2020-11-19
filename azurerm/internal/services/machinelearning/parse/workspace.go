package parse

import (
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	accountParser "github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/storage/parse"
)

type WorkspaceId struct {
	Name          string
	ResourceGroup string
}

func WorkspaceID(input string) (*WorkspaceId, error) {
	id, err := azure.ParseAzureResourceID(input)
	if err != nil {
		return nil, err
	}

	workspace := WorkspaceId{
		ResourceGroup: id.ResourceGroup,
	}

	if workspace.Name, err = id.PopSegment("workspaces"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &workspace, nil
}

// TODO -- use parse function "github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/storage/parsers".ParseAccountID
// when issue https://github.com/Azure/azure-rest-api-specs/issues/8323 is addressed
func AccountIDCaseDiffSuppress(input string) (*accountParser.AccountId, error) {
	id, err := azure.ParseAzureResourceID(input)
	if err != nil {
		return nil, err
	}

	account := accountParser.AccountId{
		ResourceGroup: id.ResourceGroup,
	}

	if account.Name, err = id.PopSegment("storageAccounts"); err != nil {
		if account.Name, err = id.PopSegment("storageaccounts"); err != nil {
			return nil, err
		}
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &account, nil
}
