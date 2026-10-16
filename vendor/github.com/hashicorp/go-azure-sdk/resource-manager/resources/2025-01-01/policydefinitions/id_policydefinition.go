package policydefinitions

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/recaser"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

func init() {
	recaser.RegisterResourceId(&PolicyDefinitionId{})
}

var _ resourceids.ResourceId = &PolicyDefinitionId{}

// PolicyDefinitionId is a struct representing the Resource ID for a Policy Definition
type PolicyDefinitionId struct {
	PolicyDefinitionName string
}

// NewPolicyDefinitionID returns a new PolicyDefinitionId struct
func NewPolicyDefinitionID(policyDefinitionName string) PolicyDefinitionId {
	return PolicyDefinitionId{
		PolicyDefinitionName: policyDefinitionName,
	}
}

// ParsePolicyDefinitionID parses 'input' into a PolicyDefinitionId
func ParsePolicyDefinitionID(input string) (*PolicyDefinitionId, error) {
	parser := resourceids.NewParserFromResourceIdType(&PolicyDefinitionId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := PolicyDefinitionId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParsePolicyDefinitionIDInsensitively parses 'input' case-insensitively into a PolicyDefinitionId
// note: this method should only be used for API response data and not user input
func ParsePolicyDefinitionIDInsensitively(input string) (*PolicyDefinitionId, error) {
	parser := resourceids.NewParserFromResourceIdType(&PolicyDefinitionId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := PolicyDefinitionId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *PolicyDefinitionId) FromParseResult(input resourceids.ParseResult) error {
	var ok bool

	if id.PolicyDefinitionName, ok = input.Parsed["policyDefinitionName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "policyDefinitionName", input)
	}

	return nil
}

// ValidatePolicyDefinitionID checks that 'input' can be parsed as a Policy Definition ID
func ValidatePolicyDefinitionID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParsePolicyDefinitionID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Policy Definition ID
func (id PolicyDefinitionId) ID() string {
	fmtString := "/providers/Microsoft.Authorization/policyDefinitions/%s"
	return fmt.Sprintf(fmtString, id.PolicyDefinitionName)
}

// Segments returns a slice of Resource ID Segments which comprise this Policy Definition ID
func (id PolicyDefinitionId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftAuthorization", "Microsoft.Authorization", "Microsoft.Authorization"),
		resourceids.StaticSegment("staticPolicyDefinitions", "policyDefinitions", "policyDefinitions"),
		resourceids.UserSpecifiedSegment("policyDefinitionName", "policyDefinitionName"),
	}
}

// String returns a human-readable description of this Policy Definition ID
func (id PolicyDefinitionId) String() string {
	components := []string{
		fmt.Sprintf("Policy Definition Name: %q", id.PolicyDefinitionName),
	}
	return fmt.Sprintf("Policy Definition (%s)", strings.Join(components, "\n"))
}
