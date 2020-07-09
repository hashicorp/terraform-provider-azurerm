package parse

import (
	"fmt"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
)

type KustoDatabasePrincipalAssignmentId struct {
	ResourceGroup string
	Cluster       string
	Database      string
	Name          string
}

func KustoDatabasePrincipalAssignmentID(input string) (*KustoDatabasePrincipalAssignmentId, error) {
	id, err := azure.ParseAzureResourceID(input)
	if err != nil {
		return nil, fmt.Errorf("[ERROR] Unable to parse Kusto Database Principal ID %q: %+v", input, err)
	}

	principal := KustoDatabasePrincipalAssignmentId{
		ResourceGroup: id.ResourceGroup,
	}

	if principal.Cluster, err = id.PopSegment("Clusters"); err != nil {
		return nil, err
	}

	if principal.Database, err = id.PopSegment("Databases"); err != nil {
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
