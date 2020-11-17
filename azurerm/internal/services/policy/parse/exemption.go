package parse

import (
	"fmt"
	"regexp"
)

type PolicyExemptionId struct {
	Name string
	PolicyScopeId
}

func PolicyExemptionID(input string) (*PolicyExemptionId, error) {
	// in general, the id of an exemption should be:
	// {scope}/providers/Microsoft.Authorization/policyExemptions/{name}
	regex := regexp.MustCompile(`/providers/Microsoft\.Authorization/policyExemptions/`)
	if !regex.MatchString(input) {
		return nil, fmt.Errorf("unable to parse Policy Exemption ID %q", input)
	}

	segments := regex.Split(input, -1)

	if len(segments) != 2 {
		return nil, fmt.Errorf("unable to parse Policy Exemption ID %q: Expected 2 segments after split", input)
	}

	scope, name := segments[0], segments[1]
	if name == "" {
		return nil, fmt.Errorf("unable to parse Policy Exemption ID %q: exemption name is empty", input)
	}

	scopeId, err := PolicyScopeID(scope)
	if err != nil {
		return nil, fmt.Errorf("unable to parse Policy Exemption ID %q: %+v", input, err)
	}

	return &PolicyExemptionId{
		Name:          name,
		PolicyScopeId: scopeId,
	}, nil
}
