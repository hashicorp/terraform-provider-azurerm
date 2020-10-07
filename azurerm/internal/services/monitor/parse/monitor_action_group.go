package parse

import (
	"fmt"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
)

type ActionGroupId struct {
	ResourceGroup string
	Name          string
}

func ActionGroupID(input string) (*ActionGroupId, error) {
	id, err := azure.ParseAzureResourceID(input)
	if err != nil {
		return nil, fmt.Errorf("parsing Action Group ID %q: %+v", input, err)
	}

	actionGroup := ActionGroupId{
		ResourceGroup: id.ResourceGroup,
	}

	if actionGroup.Name, err = id.PopSegment("actionGroups"); err != nil {
		return nil, fmt.Errorf("parsing Action Group ID %q: %+v", input, err)
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, fmt.Errorf("parsing Action Group ID %q: %+v", input, err)
	}

	return &actionGroup, nil
}
