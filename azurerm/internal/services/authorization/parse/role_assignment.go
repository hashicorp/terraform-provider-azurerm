package parse

import (
	"fmt"
	"strings"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
)

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
