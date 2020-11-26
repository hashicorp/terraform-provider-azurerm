package parse

import (
	"fmt"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
)

type DatabasePrincipalAssignmentId struct {
	ResourceGroup           string
	ClusterName             string
	DatabaseName            string
	PrincipalAssignmentName string
}

func DatabasePrincipalAssignmentID(input string) (*DatabasePrincipalAssignmentId, error) {
	id, err := azure.ParseAzureResourceID(input)
	if err != nil {
		return nil, fmt.Errorf("[ERROR] Unable to parse Kusto Database Principal ID %q: %+v", input, err)
	}

	principal := DatabasePrincipalAssignmentId{
		ResourceGroup: id.ResourceGroup,
	}

	if principal.ClusterName, err = id.PopSegment("Clusters"); err != nil {
		return nil, err
	}

	if principal.DatabaseName, err = id.PopSegment("Databases"); err != nil {
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
