package parse

import (
	"fmt"
	"regexp"
)

type PolicyDefinitionId struct {
	Name string
	PolicyScopeId
}

// TODO: This paring function is currently suppressing every case difference due to github issue: https://github.com/Azure/azure-rest-api-specs/issues/8353
func PolicyDefinitionID(input string) (*PolicyDefinitionId, error) {
	// in general, the id of a definition should be:
	// {scope}/providers/Microsoft.PolicyInsights/remediations/{name}
	regex := regexp.MustCompile(`/providers/[Mm]icrosoft\.[Aa]uthorization/policy[Dd]efinitions/`)
	if !regex.MatchString(input) {
		return nil, fmt.Errorf("unable to parse Policy Definition ID %q", input)
	}

	segments := regex.Split(input, -1)

	if len(segments) != 2 {
		return nil, fmt.Errorf("unable to parse Policy Definition ID %q: Expected 2 segments after split", input)
	}

	scope := segments[0]
	name := segments[1]
	if name == "" {
		return nil, fmt.Errorf("unable to parse Policy Definition ID %q: assignment name is empty", input)
	}

	scopeId, err := PolicyScopeID(scope)
	if err != nil {
		return nil, fmt.Errorf("unable to parse Policy Definition ID %q: %+v", input, err)
	}

	return &PolicyDefinitionId{
		Name:          name,
		PolicyScopeId: scopeId,
	}, nil
}
