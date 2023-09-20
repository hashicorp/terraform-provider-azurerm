// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package parse

import (
	"fmt"
	"regexp"
	"strings"
)

type ResourcePolicyRemediationId struct {
	Name       string
	ResourceId string
}

func (id ResourcePolicyRemediationId) String() string {
	segments := []string{
		fmt.Sprintf("Resource Policy Remediation Name %q", id.Name),
		fmt.Sprintf("Resource ID %q", id.ResourceId),
	}
	segmentsStr := strings.Join(segments, " / ")
	return fmt.Sprintf("%s: (%s)", "Resource Policy Remediation ID", segmentsStr)
}

func (id ResourcePolicyRemediationId) ID() string {
	fmtString := "%s/providers/Microsoft.PolicyInsights/remediations/%s"
	return fmt.Sprintf(fmtString, id.ResourceId, id.Name)
}

func NewResourcePolicyRemediationId(resourceID, name string) ResourcePolicyRemediationId {
	return ResourcePolicyRemediationId{
		Name:       name,
		ResourceId: resourceID,
	}
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
