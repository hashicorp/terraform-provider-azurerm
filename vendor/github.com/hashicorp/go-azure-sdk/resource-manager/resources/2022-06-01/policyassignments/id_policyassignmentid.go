package policyassignments

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/recaser"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

func init() {
	recaser.RegisterResourceId(&PolicyAssignmentIdId{})
}

var _ resourceids.ResourceId = &PolicyAssignmentIdId{}

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
	parser := resourceids.NewParserFromResourceIdType(&PolicyAssignmentIdId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := PolicyAssignmentIdId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParsePolicyAssignmentIdIDInsensitively parses 'input' case-insensitively into a PolicyAssignmentIdId
// note: this method should only be used for API response data and not user input
func ParsePolicyAssignmentIdIDInsensitively(input string) (*PolicyAssignmentIdId, error) {
	parser := resourceids.NewParserFromResourceIdType(&PolicyAssignmentIdId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := PolicyAssignmentIdId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *PolicyAssignmentIdId) FromParseResult(input resourceids.ParseResult) error {
	var ok bool

	if id.PolicyAssignmentId, ok = input.Parsed["policyAssignmentId"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "policyAssignmentId", input)
	}

	return nil
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
		resourceids.UserSpecifiedSegment("policyAssignmentId", "policyAssignmentId"),
	}
}

// String returns a human-readable description of this Policy Assignment Id ID
func (id PolicyAssignmentIdId) String() string {
	components := []string{
		fmt.Sprintf("Policy Assignment: %q", id.PolicyAssignmentId),
	}
	return fmt.Sprintf("Policy Assignment Id (%s)", strings.Join(components, "\n"))
}
