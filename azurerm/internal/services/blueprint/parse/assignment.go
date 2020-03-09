package parse

import (
	"fmt"
	"regexp"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
)

type BlueprintAssignmentId struct {
	Name string
	BlueprintAssignmentScopeId
}

func BlueprintAssignmentID(input string) (*BlueprintAssignmentId, error) {
	regex := regexp.MustCompile(`/providers/Microsoft\.Blueprint/blueprintAssignments/`)
	if !regex.MatchString(input) {
		return nil, fmt.Errorf("unable to parse Blueprint Assignment ID %q", input)
	}

	segments := regex.Split(input, -1)

	if len(segments) != 2 {
		return nil, fmt.Errorf("unable to parse Blueprint Assignment ID %q: expeceted 2 segmetns", input)
	}

	scope := segments[0]
	scopeId, err := BlueprintAssignmentScopeID(scope)
	if err != nil {
		return nil, fmt.Errorf("unable to parse Blueprint Assignment ID %q: %+v", input, err)
	}

	name := segments[1]
	if name == "" {
		return nil, fmt.Errorf("unable to parse Blueprint Assignment ID %q: asssignment name is empty", input)
	}

	id := BlueprintAssignmentId{
		Name:                       name,
		BlueprintAssignmentScopeId: *scopeId,
	}

	return &id, nil
}

type BlueprintAssignmentScopeId struct {
	ScopeId        string
	SubscriptionId string
}

func BlueprintAssignmentScopeID(input string) (*BlueprintAssignmentScopeId, error) {
	id, err := azure.ParseAzureResourceID(input)
	if err != nil {
		return nil, fmt.Errorf("unable to parse Blueprint Assignment Scope ID %q: %+v", input, err)
	}
	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, fmt.Errorf("unable to parse Blueprint Assignment Scope ID %q: %+v", input, err)
	}
	if id.ResourceGroup != "" {
		return nil, fmt.Errorf("unable to parse Blueprint Assignment Scope ID %q: scope cannot have resource groups", input)
	}
	scopeId := BlueprintAssignmentScopeId{
		ScopeId:        input,
		SubscriptionId: id.SubscriptionID,
	}

	return &scopeId, nil
}
