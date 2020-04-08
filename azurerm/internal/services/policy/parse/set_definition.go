package parse

import (
	"fmt"
	"regexp"
)

type PolicySetDefinitionId struct {
	Name string
	PolicyScopeId
}

// TODO: This parsing function is currently suppressing case difference due to github issue: https://github.com/Azure/azure-rest-api-specs/issues/8353
func PolicySetDefinitionID(input string) (*PolicySetDefinitionId, error) {
	// in general, the id of a set definition should be:
	// {scope}/providers/Microsoft.Authorization/policySetDefinitions/set1
	regex := regexp.MustCompile(`/providers/[Mm]icrosoft\.[Aa]uthorization/policy[Ss]et[Dd]efinitions/`)
	if !regex.MatchString(input) {
		return nil, fmt.Errorf("unable to parse Policy Set Definition ID %q", input)
	}

	segments := regex.Split(input, -1)

	if len(segments) != 2 {
		return nil, fmt.Errorf("unable to parse Policy Set Definition ID %q: Expected 2 segments after split", input)
	}

	scope := segments[0]
	name := segments[1]
	if name == "" {
		return nil, fmt.Errorf("unable to parse Policy Set Definition ID %q: set definition name is empty", input)
	}

	scopeId, err := PolicyScopeID(scope)
	if err != nil {
		return nil, fmt.Errorf("unable to parse Policy Set Definition ID %q: %+v", input, err)
	}

	return &PolicySetDefinitionId{
		Name:          name,
		PolicyScopeId: scopeId,
	}, nil
}
