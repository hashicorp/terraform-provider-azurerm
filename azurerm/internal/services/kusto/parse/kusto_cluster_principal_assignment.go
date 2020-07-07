package parse

import (
	"fmt"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
)

type KustoClusterPrincipalAssignmentId struct {
	ResourceGroup string
	Cluster       string
	Name          string
}

func KustoClusterPrincipalAssignmentID(input string) (*KustoClusterPrincipalAssignmentId, error) {
	id, err := azure.ParseAzureResourceID(input)
	if err != nil {
		return nil, fmt.Errorf("[ERROR] Unable to parse Kusto Cluster Principal ID %q: %+v", input, err)
	}

	principal := KustoClusterPrincipalAssignmentId{
		ResourceGroup: id.ResourceGroup,
	}

	if principal.Cluster, err = id.PopSegment("Clusters"); err != nil {
		return nil, err
	}

	if principal.Name, err = id.PopSegment("PrincipalAssignments"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &principal, nil
}
