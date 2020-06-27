package parse

import (
	"fmt"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
)

type KustoDatabasePrincipalId struct {
	ResourceGroup string
	Cluster       string
	Database      string
	Role          string
	Name          string
}

func KustoDatabasePrincipalID(input string) (*KustoDatabasePrincipalId, error) {
	id, err := azure.ParseAzureResourceID(input)
	if err != nil {
		return nil, fmt.Errorf("[ERROR] Unable to parse Kusto Database Principal ID %q: %+v", input, err)
	}

	principal := KustoDatabasePrincipalId{
		ResourceGroup: id.ResourceGroup,
	}

	if principal.Cluster, err = id.PopSegment("Clusters"); err != nil {
		return nil, err
	}

	if principal.Database, err = id.PopSegment("Databases"); err != nil {
		return nil, err
	}

	if principal.Role, err = id.PopSegment("Role"); err != nil {
		return nil, err
	}

	if principal.Name, err = id.PopSegment("FQN"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &principal, nil
}
