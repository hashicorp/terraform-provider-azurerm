package parse

import (
	"fmt"
	"strings"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
)

type AssignmentId struct {
	// "/subscriptions/00000000-0000-0000-0000-000000000000/providers/Microsoft.Blueprint/blueprintAssignments/assignSimpleBlueprint",
	// "/managementGroups/ContosoOnlineGroup/providers/Microsoft.Blueprint/blueprintAssignments/assignSimpleBlueprint",

	Scope           string
	Subscription    string
	ManagementGroup string
	Name            string
}

func AssignmentID(input string) (*AssignmentId, error) {
	if len(input) == 0 {
		return nil, fmt.Errorf("Bad: Assignment ID is empty string")
	}

	assignmentID := AssignmentId{}

	idParts := strings.Split(strings.Trim(input, "/"), "/")
	if len(idParts) != 6 {
		return nil, fmt.Errorf("Bad: Blueprint Assignment ID invalid: %q", input)
	}

	// check casing on segments
	if idParts[2] != "providers" || idParts[3] != "Microsoft.Blueprint" {
		return nil, fmt.Errorf("ID has invalid provider segment (should be `providers/Microsoft.Blueprint` case sensitive): %q", input)
	}

	if idParts[4] != "blueprintAssignments" {
		return nil, fmt.Errorf("ID missing `blueprintAssignments` segment (case sensitive): %q", input)
	}

	switch idParts[0] {
	case "managementGroups":
		assignmentID = AssignmentId{
			Scope:           fmt.Sprintf("%s/%s", idParts[0], idParts[1]),
			ManagementGroup: idParts[1],
		}

	case "subscriptions":
		id, err := azure.ParseAzureResourceID(input)
		if err != nil {
			return nil, fmt.Errorf("[ERROR] Unable to parse Image ID %q: %+v", input, err)
		}

		assignmentID.Scope = fmt.Sprintf("subscriptions/%s", id.SubscriptionID)
		assignmentID.Subscription = id.SubscriptionID

	default:
		return nil, fmt.Errorf("Bad: Invalid ID, should start with one of `/managementGroups` or `/subscriptions`: %q", input)
	}

	assignmentID.Name = idParts[5]

	return &assignmentID, nil
}
