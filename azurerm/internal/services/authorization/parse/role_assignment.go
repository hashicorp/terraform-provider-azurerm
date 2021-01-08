package parse

import (
	"fmt"
	"strings"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
)

// /subscriptions/00000000-0000-0000-0000-000000000000/providers/Microsoft.Authorization/roleAssignments/00000000-0000-0000-0000-000000000000
// /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.Authorization/roleAssignments/00000000-0000-0000-0000-000000000000
// /providers/Microsoft.Management/managementGroups/74253c84-2d7f-4522-a73f-8897fd715d21/providers/Microsoft.Authorization/roleAssignments/4d6124f3-1888-5851-0ffa-3b86bd56a2e4

type RoleAssignmentID struct {
	SubscriptionID  string
	ResourceGroup   string
	ManagementGroup string
	Name            string
}

func RoleAssignmentId(input string) (*RoleAssignmentID, error) {
	if len(input) == 0 {
		return nil, fmt.Errorf("Role Assignment ID is empty string")
	}

	roleAssignmentId := RoleAssignmentID{}

	if strings.HasPrefix(input, "/subscriptions/") {
		id, err := azure.ParseAzureResourceID(input)
		if err != nil {
			return nil, fmt.Errorf("could not parse %q as Azure resource ID", input)
		}
		roleAssignmentId.SubscriptionID = id.SubscriptionID
		roleAssignmentId.ResourceGroup = id.ResourceGroup
		if roleAssignmentId.Name, err = id.PopSegment("roleAssignments"); err != nil {
			return nil, err
		}
	} else if strings.HasPrefix(input, "/providers/Microsoft.Management/") {
		idParts := strings.Split(input, "/providers/Microsoft.Authorization/roleAssignments/")
		if len(idParts) != 2 {
			return nil, fmt.Errorf("could not parse Role Assignment ID %q for Management Group", input)
		}
		roleAssignmentId.Name = idParts[1]

		roleAssignmentId.ManagementGroup = strings.Trim(idParts[0], "/providers/Microsoft.Management/managementGroups/")

	} else {
		return nil, fmt.Errorf("could not parse Role Assignment ID %q", input)
	}

	return &roleAssignmentId, nil

}
