package parse

import (
	"fmt"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
)

type DatabasePrincipalId struct {
	ResourceGroup string
	ClusterName   string
	DatabaseName  string
	RoleName      string
	FQNName       string
}

func DatabasePrincipalID(input string) (*DatabasePrincipalId, error) {
	id, err := azure.ParseAzureResourceID(input)
	if err != nil {
		return nil, fmt.Errorf("[ERROR] Unable to parse Kusto Database Principal ID %q: %+v", input, err)
	}

	principal := DatabasePrincipalId{
		ResourceGroup: id.ResourceGroup,
	}

	if principal.ClusterName, err = id.PopSegment("Clusters"); err != nil {
		return nil, err
	}

	if principal.DatabaseName, err = id.PopSegment("Databases"); err != nil {
		return nil, err
	}

	if principal.RoleName, err = id.PopSegment("Role"); err != nil {
		return nil, err
	}

	if principal.FQNName, err = id.PopSegment("FQN"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &principal, nil
}
