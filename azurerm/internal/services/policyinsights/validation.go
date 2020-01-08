package policyinsights

import (
	"fmt"
	"strings"
)

func validateRemediationName(i interface{}, k string) (warnings []string, errors []error) {
	v, ok := i.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected type of %q to be string", k))
		return
	}

	// The service returns error when name of remediation is too long
	// error: The remediation name cannot be empty and must not exceed '260' characters.
	// By my additional test, the name of remediation cannot contain the following characters: %^#/\&?.
	if len(v) == 0 || len(v) > 260 {
		errors = append(errors, fmt.Errorf("%s cannot be empty and must not exceed '260' characters", k))
		return
	}
	invalidCharacters := `%^#/\&?`
	if strings.ContainsAny(v, invalidCharacters) {
		errors = append(errors, fmt.Errorf("%s cannot contain the following characters: %s", k, invalidCharacters))
	}

	return warnings, errors
}

func validatePolicyAssignmentID(i interface{}, k string) (warnings []string, errors []error) {
	v, ok := i.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected type of %q to be string", k))
		return
	}
	vLower := strings.ToLower(v)

	const provider = "/providers/microsoft.authorization/policyassignments"
	index := strings.LastIndex(vLower, provider)
	if index < 0 {
		errors = append(errors, fmt.Errorf("cannot recognize %s `%s` as PolicyAssignment ID", k, v))
		return
	}

	scope := v[0:index]
	assignmentPath := v[index+1:]
	// scope should be a resource ID, resource group ID, subscription ID, or Management Group ID
	_, err := ParseScope(scope)
	if err != nil {
		errors = append(errors, err)
	}
	// assignment should have a name
	segments := strings.Split(assignmentPath, "/")
	if len(segments) != 4 {
		errors = append(errors, fmt.Errorf("expect the following part of policy assignment id to have 4 segment, but got %d", len(segments)))
	}
	if segments[3] == "" {
		errors = append(errors, fmt.Errorf("the policy assignment name should not be empty"))
	}

	return warnings, errors
}

func validatePolicyDefinitionID(i interface{}, k string) (warnings []string, errors []error) {
	v, ok := i.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected type of %q to be string", k))
		return
	}
	vLower := strings.ToLower(v)

	const provider = "/providers/microsoft.authorization/policydefinitions/"
	index := strings.LastIndex(vLower, provider)
	if index < 0 {
		errors = append(errors, fmt.Errorf("cannot recognize %s `%s` as PolicyAssignment ID", k, v))
		return
	}

	scope := v[0:index]
	definitionPath := v[index+1:]
	// scope should be a Subscription ID or ManagementGroup ID
	_, err := ParseScope(scope)
	if err != nil {
		errors = append(errors, err)
	}
	segments := strings.Split(definitionPath, "/")
	if len(segments) != 4 {
		errors = append(errors, fmt.Errorf("expect the following part of policy definition id to have 4 segment, but got %d", len(segments)))
	}
	if segments[3] == "" {
		errors = append(errors, fmt.Errorf("the policy definition name should not be empty"))
	}

	return warnings, errors
}
