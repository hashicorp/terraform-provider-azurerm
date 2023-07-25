package policyassignments

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ resourceids.ResourceId = PolicyAssignmentIdId{}

// PolicyAssignmentIdId is a struct representing the Resource ID for a Policy Assignment Id
type PolicyAssignmentIdId struct {
	PolicyAssignmentId string
}

// NewPolicyAssignmentIdID returns a new PolicyAssignmentIdId struct
func NewPolicyAssignmentIdID(policyAssignmentId string) PolicyAssignmentIdId {
	return PolicyAssignmentIdId{
		PolicyAssignmentId: policyAssignmentId,
	}
}

// ParsePolicyAssignmentIdID parses 'input' into a PolicyAssignmentIdId
func ParsePolicyAssignmentIdID(input string) (*PolicyAssignmentIdId, error) {
	parser := resourceids.NewParserFromResourceIdType(PolicyAssignmentIdId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := PolicyAssignmentIdId{}

	if id.PolicyAssignmentId, ok = parsed.Parsed["policyAssignmentId"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "policyAssignmentId", *parsed)
	}

	return &id, nil
}

// ParsePolicyAssignmentIdIDInsensitively parses 'input' case-insensitively into a PolicyAssignmentIdId
// note: this method should only be used for API response data and not user input
func ParsePolicyAssignmentIdIDInsensitively(input string) (*PolicyAssignmentIdId, error) {
	parser := resourceids.NewParserFromResourceIdType(PolicyAssignmentIdId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := PolicyAssignmentIdId{}

	if id.PolicyAssignmentId, ok = parsed.Parsed["policyAssignmentId"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "policyAssignmentId", *parsed)
	}

	return &id, nil
}

// ValidatePolicyAssignmentIdID checks that 'input' can be parsed as a Policy Assignment Id ID
func ValidatePolicyAssignmentIdID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParsePolicyAssignmentIdID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Policy Assignment Id ID
func (id PolicyAssignmentIdId) ID() string {
	fmtString := "/%s"
	return fmt.Sprintf(fmtString, id.PolicyAssignmentId)
}

// Segments returns a slice of Resource ID Segments which comprise this Policy Assignment Id ID
func (id PolicyAssignmentIdId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.UserSpecifiedSegment("policyAssignmentId", "policyAssignmentIdValue"),
	}
}

// String returns a human-readable description of this Policy Assignment Id ID
func (id PolicyAssignmentIdId) String() string {
	components := []string{
		fmt.Sprintf("Policy Assignment: %q", id.PolicyAssignmentId),
	}
	return fmt.Sprintf("Policy Assignment Id (%s)", strings.Join(components, "\n"))
}
