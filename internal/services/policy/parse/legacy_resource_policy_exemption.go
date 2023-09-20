// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package parse

import (
	"fmt"
	"regexp"
	"strings"
)

type ResourcePolicyExemptionId struct {
	Name       string
	ResourceId string
}

func (id ResourcePolicyExemptionId) String() string {
	segments := []string{
		fmt.Sprintf("Resource Policy Exemption Name %q", id.Name),
		fmt.Sprintf("Resource ID %q", id.ResourceId),
	}
	segmentsStr := strings.Join(segments, " / ")
	return fmt.Sprintf("%s: (%s)", "Resource Policy Exemption ID", segmentsStr)
}

func (id ResourcePolicyExemptionId) ID() string {
	fmtString := "%s/providers/Microsoft.Authorization/policyExemptions/%s"
	return fmt.Sprintf(fmtString, id.ResourceId, id.Name)
}

func NewResourcePolicyExemptionId(resourceID, name string) ResourcePolicyExemptionId {
	return ResourcePolicyExemptionId{
		Name:       name,
		ResourceId: resourceID,
	}
}

// TODO: This paring function is currently suppressing every case difference due to github issue: https://github.com/Azure/azure-rest-api-specs/issues/8353
func ResourcePolicyExemptionID(input string) (*ResourcePolicyExemptionId, error) {
	regex := regexp.MustCompile(`/providers/Microsoft.Authorization/policyExemptions/`)
	if !regex.MatchString(input) {
		return nil, fmt.Errorf("unable to parse Resource Policy Exemption ID %q", input)
	}

	segments := regex.Split(input, -1)

	if len(segments) != 2 {
		return nil, fmt.Errorf("unable to parse Resource Policy Exemption ID %q: Expected 2 segments after split", input)
	}

	resourceId := segments[0]
	name := segments[1]
	if name == "" {
		return nil, fmt.Errorf("unable to parse Resource Policy Exemption ID %q: assignment name is empty", input)
	}

	return &ResourcePolicyExemptionId{
		Name:       name,
		ResourceId: resourceId,
	}, nil
}
