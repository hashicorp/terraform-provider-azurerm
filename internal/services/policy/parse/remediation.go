package parse

import (
	"fmt"
	"regexp"
)

type PolicyRemediationId struct {
	Name string
	PolicyScopeId
}

// TODO: This paring function is currently suppressing every case difference due to github issue: https://github.com/Azure/azure-rest-api-specs/issues/8353
// Currently the returned Remediation response from the service will have all the IDs converted into lower cases
func PolicyRemediationID(input string) (*PolicyRemediationId, error) {
	// in general, the id of a remediation should be:
	// {scope}/providers/Microsoft.PolicyInsights/remediations/{name}
	regex := regexp.MustCompile(`/providers/[Mm]icrosoft\.[Pp]olicy[Ii]nsights/remediations/`)
	if !regex.MatchString(input) {
		return nil, fmt.Errorf("unable to parse Policy Remediation ID %q", input)
	}

	segments := regex.Split(input, -1)

	if len(segments) != 2 {
		return nil, fmt.Errorf("unable to parse Policy Remediation ID %q: Expected 2 segments after split", input)
	}

	scope := segments[0]
	name := segments[1]
	if name == "" {
		return nil, fmt.Errorf("unable to parse Policy Remediation ID %q: remediation name is empty", input)
	}

	scopeId, err := PolicyScopeID(scope)
	if err != nil {
		return nil, fmt.Errorf("unable to parse Policy Remediation ID %q: %+v", input, err)
	}

	return &PolicyRemediationId{
		Name:          name,
		PolicyScopeId: scopeId,
	}, nil
}
