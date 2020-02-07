package parse

import (
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	accountParser "github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/storage/parsers"
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

func AccountIDCaseDiffSuppress(input string) (*accountParser.AccountID, error) {
	accountId, err := accountParser.ParseAccountID(input)
	if err == nil {
		return accountId, nil
	}
	id, err := azure.ParseAzureResourceID(input)
	if err != nil {
		return nil, err
	}

	account := accountParser.AccountID{
		ResourceGroup: id.ResourceGroup,
	}

	if account.Name, err = id.PopSegment("storageaccounts"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &account, nil
}
