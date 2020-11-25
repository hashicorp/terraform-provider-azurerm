package parse

import (
	"fmt"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
)

type ProjectId struct {
	ResourceGroup string
	ServiceName   string
	Name          string
}

func ProjectID(input string) (*ProjectId, error) {
	id, err := azure.ParseAzureResourceID(input)
	if err != nil {
		return nil, fmt.Errorf("[ERROR] Unable to parse Database Migration Project ID %q: %+v", input, err)
	}

	project := ProjectId{
		ResourceGroup: id.ResourceGroup,
	}

	if project.ServiceName, err = id.PopSegment("services"); err != nil {
		return nil, err
	}

	if project.Name, err = id.PopSegment("projects"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &project, nil
}
