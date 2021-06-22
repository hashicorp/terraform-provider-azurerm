package parse

import (
	"fmt"
	"regexp"
	"strings"
)

type PolicyAssignmentId struct {
	Name  string
	Scope string
	PolicyScopeId
}

func (id PolicyAssignmentId) String() string {
	segments := []string{
		fmt.Sprintf("Assignment Name %q", id.Name),
		fmt.Sprintf("Scope %q", id.Scope),
	}
	segmentsStr := strings.Join(segments, " / ")
	return fmt.Sprintf("%s: (%s)", "Policy Assignment ID", segmentsStr)
}

func (id PolicyAssignmentId) ID() string {
	fmtString := "%s/providers/Microsoft.Authorization/policyAssignments/%s"
	return fmt.Sprintf(fmtString, id.Scope, id.Name)
}

func NewPolicyAssignmentId(scope, name string) (*PolicyAssignmentId, error) {
	scopeId, err := PolicyScopeID(scope)
	if err != nil {
		return nil, fmt.Errorf("parsing Policy Scope ID %q: %+v", scope, err)
	}

	return &PolicyAssignmentId{
		Name:          name,
		PolicyScopeId: scopeId,
		Scope:         scope,
	}, nil
}

// TODO: This paring function is currently suppressing every case difference due to github issue: https://github.com/Azure/azure-rest-api-specs/issues/8353
func PolicyAssignmentID(input string) (*PolicyAssignmentId, error) {
	// in general, the id of a assignment should be:
	// {scope}/providers/Microsoft.Authorization/policyAssignment/{name}
	regex := regexp.MustCompile(`/providers/[Mm]icrosoft\.[Aa]uthorization/policy[Aa]ssignments/`)
	if !regex.MatchString(input) {
		return nil, fmt.Errorf("unable to parse Policy Assignment ID %q", input)
	}

	segments := regex.Split(input, -1)

	if len(segments) != 2 {
		return nil, fmt.Errorf("unable to parse Policy Assignment ID %q: Expected 2 segments after split", input)
	}

	scope := segments[0]
	name := segments[1]
	if name == "" {
		return nil, fmt.Errorf("unable to parse Policy Assignment ID %q: assignment name is empty", input)
	}

	scopeId, err := PolicyScopeID(scope)
	if err != nil {
		return nil, fmt.Errorf("unable to parse Policy Assignment ID %q: %+v", input, err)
	}

	return &PolicyAssignmentId{
		Name:          name,
		Scope:         scope,
		PolicyScopeId: scopeId,
	}, nil
}
