package parse

import (
	"fmt"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
)

type DatabricksWorkspaceId struct {
	ResourceGroup string
	Name          string
}

func DatabricksWorkspaceID(input string) (*DatabricksWorkspaceId, error) {
	id, err := azure.ParseAzureResourceID(input)
	if err != nil {
		return nil, fmt.Errorf("[ERROR] Unable to parse Databricks Workspace ID %q: %+v", input, err)
	}

	workspace := DatabricksWorkspaceId{
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
