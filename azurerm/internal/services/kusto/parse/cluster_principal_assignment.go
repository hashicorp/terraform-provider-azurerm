package parse

import (
	"fmt"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
)

type ClusterPrincipalAssignmentId struct {
	ResourceGroup           string
	ClusterName             string
	PrincipalAssignmentName string
}

func ClusterPrincipalAssignmentID(input string) (*ClusterPrincipalAssignmentId, error) {
	id, err := azure.ParseAzureResourceID(input)
	if err != nil {
		return nil, fmt.Errorf("[ERROR] Unable to parse Kusto Cluster Principal ID %q: %+v", input, err)
	}

	principal := ClusterPrincipalAssignmentId{
		ResourceGroup: id.ResourceGroup,
	}

	if principal.ClusterName, err = id.PopSegment("Clusters"); err != nil {
		return nil, err
	}

	if principal.PrincipalAssignmentName, err = id.PopSegment("PrincipalAssignments"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &principal, nil
}
