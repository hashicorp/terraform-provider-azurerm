package parse

import (
	"fmt"
	"strings"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
)

type VersionId struct {
	// "/{resourceScope}/providers/Microsoft.Blueprint/blueprints/{blueprintName}/versions/{versionId}"
	// "/subscriptions/00000000-0000-0000-0000-000000000000/providers/Microsoft.Blueprint/blueprints/simpleBlueprint/versions/v1"
	// "/providers/Microsoft.Management/managementGroups/ContosoOnlineGroup/providers/Microsoft.Blueprint/blueprints/simpleBlueprint/versions/v1"

	Scope     string
	Blueprint string
	Name      string
}

func VersionID(input string) (*VersionId, error) {
	if len(input) == 0 {
		return nil, fmt.Errorf("Bad: Blueprint version ID cannot be an empty string")
	}

	versionId := VersionId{}

	idParts := strings.Split(strings.Trim(input, "/"), "/")
	if len(idParts) != 8 && len(idParts) != 10 {
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

		versionId = VersionId{
			Scope:     fmt.Sprintf("providers/Microsoft.Management/managementGroups/%s", idParts[3]),
			Blueprint: idParts[6],
			Name:      idParts[9],
		}

	case "subscriptions":
		id, err := azure.ParseAzureResourceID(input)
		if err != nil {
			return nil, fmt.Errorf("unable to parse Resource ID %q: %+v", input, err)
		}

		versionId.Scope = fmt.Sprintf("subscriptions/%s", id.SubscriptionID)
		versionId.Blueprint = idParts[5]
		versionId.Name = idParts[7]

	default:
		return nil, fmt.Errorf("Bad: Invalid ID, should start with one of `/provider` or `/subscriptions`: %q", input)
	}

	return &versionId, nil
}
