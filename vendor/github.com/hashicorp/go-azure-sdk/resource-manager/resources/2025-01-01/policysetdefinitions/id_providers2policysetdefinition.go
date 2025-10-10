package policysetdefinitions

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/recaser"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

func init() {
	recaser.RegisterResourceId(&Providers2PolicySetDefinitionId{})
}

var _ resourceids.ResourceId = &Providers2PolicySetDefinitionId{}

// Providers2PolicySetDefinitionId is a struct representing the Resource ID for a Providers 2 Policy Set Definition
type Providers2PolicySetDefinitionId struct {
	ManagementGroupName     string
	PolicySetDefinitionName string
}

// NewProviders2PolicySetDefinitionID returns a new Providers2PolicySetDefinitionId struct
func NewProviders2PolicySetDefinitionID(managementGroupName string, policySetDefinitionName string) Providers2PolicySetDefinitionId {
	return Providers2PolicySetDefinitionId{
		ManagementGroupName:     managementGroupName,
		PolicySetDefinitionName: policySetDefinitionName,
	}
}

// ParseProviders2PolicySetDefinitionID parses 'input' into a Providers2PolicySetDefinitionId
func ParseProviders2PolicySetDefinitionID(input string) (*Providers2PolicySetDefinitionId, error) {
	parser := resourceids.NewParserFromResourceIdType(&Providers2PolicySetDefinitionId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := Providers2PolicySetDefinitionId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParseProviders2PolicySetDefinitionIDInsensitively parses 'input' case-insensitively into a Providers2PolicySetDefinitionId
// note: this method should only be used for API response data and not user input
func ParseProviders2PolicySetDefinitionIDInsensitively(input string) (*Providers2PolicySetDefinitionId, error) {
	parser := resourceids.NewParserFromResourceIdType(&Providers2PolicySetDefinitionId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := Providers2PolicySetDefinitionId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *Providers2PolicySetDefinitionId) FromParseResult(input resourceids.ParseResult) error {
	var ok bool

	if id.ManagementGroupName, ok = input.Parsed["managementGroupName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "managementGroupName", input)
	}

	if id.PolicySetDefinitionName, ok = input.Parsed["policySetDefinitionName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "policySetDefinitionName", input)
	}

	return nil
}

// ValidateProviders2PolicySetDefinitionID checks that 'input' can be parsed as a Providers 2 Policy Set Definition ID
func ValidateProviders2PolicySetDefinitionID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseProviders2PolicySetDefinitionID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Providers 2 Policy Set Definition ID
func (id Providers2PolicySetDefinitionId) ID() string {
	fmtString := "/providers/Microsoft.Management/managementGroups/%s/providers/Microsoft.Authorization/policySetDefinitions/%s"
	return fmt.Sprintf(fmtString, id.ManagementGroupName, id.PolicySetDefinitionName)
}

// Segments returns a slice of Resource ID Segments which comprise this Providers 2 Policy Set Definition ID
func (id Providers2PolicySetDefinitionId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftManagement", "Microsoft.Management", "Microsoft.Management"),
		resourceids.StaticSegment("staticManagementGroups", "managementGroups", "managementGroups"),
		resourceids.UserSpecifiedSegment("managementGroupName", "managementGroupName"),
		resourceids.StaticSegment("staticProviders2", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftAuthorization", "Microsoft.Authorization", "Microsoft.Authorization"),
		resourceids.StaticSegment("staticPolicySetDefinitions", "policySetDefinitions", "policySetDefinitions"),
		resourceids.UserSpecifiedSegment("policySetDefinitionName", "policySetDefinitionName"),
	}
}

// String returns a human-readable description of this Providers 2 Policy Set Definition ID
func (id Providers2PolicySetDefinitionId) String() string {
	components := []string{
		fmt.Sprintf("Management Group Name: %q", id.ManagementGroupName),
		fmt.Sprintf("Policy Set Definition Name: %q", id.PolicySetDefinitionName),
	}
	return fmt.Sprintf("Providers 2 Policy Set Definition (%s)", strings.Join(components, "\n"))
}
