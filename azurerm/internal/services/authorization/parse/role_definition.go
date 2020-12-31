package parse

import (
	"fmt"
	"strings"
)

type RoleDefinitionID struct {
	ResourceID string
	Scope      string
	RoleID     string
}

// RoleDefinitionId is a pseudo ID for storing Scope parameter as this it not retrievable from API
// It is formed of the Azure Resource ID for the Role and the Scope it is created against
func RoleDefinitionId(input string) (*RoleDefinitionID, error) {
	parts := strings.Split(input, "|")
	if len(parts) != 2 {
		return nil, fmt.Errorf("could not parse Role Definition ID, invalid format %q", input)
	}

	idParts := strings.Split(parts[0], "roleDefinitions/")

	if !strings.HasPrefix(parts[1], "/subscriptions/") && !strings.HasPrefix(parts[1], "/providers/Microsoft.Management/managementGroups/") {
		return nil, fmt.Errorf("failed to parse scope from Role Definition ID %q", input)
	}

	roleDefinitionID := RoleDefinitionID{
		ResourceID: parts[0],
		Scope:      parts[1],
	}

	if len(idParts) < 1 {
		return nil, fmt.Errorf("failed to parse Role Definition ID from resource ID %q", input)
	} else {
		roleDefinitionID.RoleID = idParts[1]
	}

	return &roleDefinitionID, nil
}
