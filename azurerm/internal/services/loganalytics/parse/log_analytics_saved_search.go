package parse

import (
	"fmt"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
)

type LogAnalyticsSavedSearchId struct {
	ResourceGroup string
	WorkspaceName string
	Name          string
}

func LogAnalyticsSavedSearchID(input string) (*LogAnalyticsSavedSearchId, error) {
	id, err := azure.ParseAzureResourceID(input)
	if err != nil {
		return nil, fmt.Errorf("parsing Log Analytics Saved Search ID %q: %+v", input, err)
	}

	search := LogAnalyticsSavedSearchId{
		ResourceGroup: id.ResourceGroup,
	}

	if search.WorkspaceName, err = id.PopSegment("workspaces"); err != nil {
		return nil, err
	}

	if search.Name, err = id.PopSegment("savedSearches"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &search, nil
}
