package parse

import (
	"fmt"
	"strings"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
)

type DefinitionId struct {
	// "/subscriptions/00000000-0000-0000-0000-000000000000/providers/Microsoft.Blueprint/blueprints/simpleBlueprint"
	// "/providers/Microsoft.Management/managementGroups/ContosoOnlineGroup/providers/Microsoft.Blueprint/blueprints/simpleBlueprint"

	Name  string
	Scope string
}

func DefinitionID(input string) (*DefinitionId, error) {
	if len(input) == 0 {
		return nil, fmt.Errorf("Bad: Blueprint ID cannot be an empty string")
	}

	definitionId := DefinitionId{}

	idParts := strings.Split(strings.Trim(input, "/"), "/")
	if len(idParts) != 6 && len(idParts) != 8 {
		return nil, fmt.Errorf("Bad: Blueprint Version ID invalid: %q", input)
	}

	switch idParts[0] {
	case "providers":
		// check casing on segments
		if idParts[1] != "Microsoft.Management" || idParts[2] != "managementGroups" {
			return nil, fmt.Errorf("ID has invalid provider scope segment (should be `/providers/Microsoft.Management/managementGroups` case sensitive): %q", input)
		}
		if idParts[4] != "providers" || idParts[5] != "Microsoft.Blueprint" || idParts[6] != "blueprints" {
			return nil, fmt.Errorf("Bad: ID has invalid resource provider segment(s), shoud be `/providers/Microsoft.Blueprint/blueprints/`, case sensitive: %q", input)
		}

		definitionId = DefinitionId{
			Scope: fmt.Sprintf("providers/Microsoft.Management/managementGroups/%s", idParts[3]),
			Name:  idParts[6],
		}

	case "subscriptions":
		id, err := azure.ParseAzureResourceID(input)
		if err != nil {
			return nil, fmt.Errorf("[ERROR] Unable to parse Blueprint Definition ID %q: %+v", input, err)
		}

		definitionId.Scope = fmt.Sprintf("subscriptions/%s", id.SubscriptionID)
		definitionId.Name = idParts[5]

	default:
		return nil, fmt.Errorf("Bad: Invalid ID, should start with one of `/provider` or `/subscriptions`: %q", input)
	}

	return &definitionId, nil
}
