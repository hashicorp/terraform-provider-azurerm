package parse

import (
	"fmt"
	"strings"
)

type RoleAssignmentMarketplaceId struct {
	Name     string
	TenantId string
}

func (id RoleAssignmentMarketplaceId) AzureResourceID() string {
	return fmt.Sprintf("/providers/Microsoft.Marketplace/providers/Microsoft.Authorization/roleAssignments/%s", id.Name)
}

func RoleAssignmentMarketplaceID(input string) (*RoleAssignmentMarketplaceId, error) {
	if len(input) == 0 {
		return nil, fmt.Errorf("Role Assignment Marketplace ID is empty string")
	}

	roleAssignmentId := RoleAssignmentMarketplaceId{}

	parts := strings.Split(input, "|")
	if len(parts) == 2 {
		input = parts[0]
		roleAssignmentId.TenantId = parts[1]
	}

	idParts := strings.Split(input, "/providers/Microsoft.Authorization/roleAssignments/")
	if len(idParts) != 2 {
		return nil, fmt.Errorf("could not parse Role Assignment Marketplace ID %q", input)
	}

	if idParts[0] != "/providers/Microsoft.Marketplace" {
		return nil, fmt.Errorf("resource provider %s is invalid", idParts[0])
	}

	if idParts[1] == "" {
		return nil, fmt.Errorf("ID was missing a value for the roleAssignments element")
	}

	roleAssignmentId.Name = idParts[1]

	return &roleAssignmentId, nil
}
