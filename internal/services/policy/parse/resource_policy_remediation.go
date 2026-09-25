// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package parse

import (
	"fmt"
	"regexp"
)

type ResourcePolicyRemediationId struct {
	Name       string
	ResourceId string
}

// TODO: This paring function is currently suppressing every case difference due to github issue: https://github.com/Azure/azure-rest-api-specs/issues/8353
func ResourcePolicyRemediationID(input string) (*ResourcePolicyRemediationId, error) {
	// in general, the id of a policy remediation should be:
	// {scope}/providers/Microsoft.PolicyInsights/remediations/{name}
	regex := regexp.MustCompile(`/providers/[Mm]icrosoft\.[Pp]olicy[Ii]nsights/remediations/`)
	if !regex.MatchString(input) {
		return nil, fmt.Errorf("unable to parse Resource Policy Remediation ID %q", input)
	}

	segments := regex.Split(input, -1)

	if len(segments) != 2 {
		return nil, fmt.Errorf("unable to parse Resource Policy Remediation ID %q: Expected 2 segments after split", input)
	}

	resourceId := segments[0]
	name := segments[1]
	if name == "" {
		return nil, fmt.Errorf("unable to parse Resource Policy Remediation ID %q: assignment name is empty", input)
	}

	return &ResourcePolicyRemediationId{
		Name:       name,
		ResourceId: resourceId,
	}, nil
}
