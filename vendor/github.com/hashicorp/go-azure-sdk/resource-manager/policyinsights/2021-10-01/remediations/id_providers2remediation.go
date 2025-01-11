package remediations

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/recaser"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

func init() {
	recaser.RegisterResourceId(&Providers2RemediationId{})
}

var _ resourceids.ResourceId = &Providers2RemediationId{}

// Providers2RemediationId is a struct representing the Resource ID for a Providers 2 Remediation
type Providers2RemediationId struct {
	ManagementGroupId string
	RemediationName   string
}

// NewProviders2RemediationID returns a new Providers2RemediationId struct
func NewProviders2RemediationID(managementGroupId string, remediationName string) Providers2RemediationId {
	return Providers2RemediationId{
		ManagementGroupId: managementGroupId,
		RemediationName:   remediationName,
	}
}

// ParseProviders2RemediationID parses 'input' into a Providers2RemediationId
func ParseProviders2RemediationID(input string) (*Providers2RemediationId, error) {
	parser := resourceids.NewParserFromResourceIdType(&Providers2RemediationId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := Providers2RemediationId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParseProviders2RemediationIDInsensitively parses 'input' case-insensitively into a Providers2RemediationId
// note: this method should only be used for API response data and not user input
func ParseProviders2RemediationIDInsensitively(input string) (*Providers2RemediationId, error) {
	parser := resourceids.NewParserFromResourceIdType(&Providers2RemediationId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := Providers2RemediationId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *Providers2RemediationId) FromParseResult(input resourceids.ParseResult) error {
	var ok bool

	if id.ManagementGroupId, ok = input.Parsed["managementGroupId"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "managementGroupId", input)
	}

	if id.RemediationName, ok = input.Parsed["remediationName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "remediationName", input)
	}

	return nil
}

// ValidateProviders2RemediationID checks that 'input' can be parsed as a Providers 2 Remediation ID
func ValidateProviders2RemediationID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseProviders2RemediationID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Providers 2 Remediation ID
func (id Providers2RemediationId) ID() string {
	fmtString := "/providers/Microsoft.Management/managementGroups/%s/providers/Microsoft.PolicyInsights/remediations/%s"
	return fmt.Sprintf(fmtString, id.ManagementGroupId, id.RemediationName)
}

// Segments returns a slice of Resource ID Segments which comprise this Providers 2 Remediation ID
func (id Providers2RemediationId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.StaticSegment("managementGroupsNamespace", "Microsoft.Management", "Microsoft.Management"),
		resourceids.StaticSegment("staticManagementGroups", "managementGroups", "managementGroups"),
		resourceids.UserSpecifiedSegment("managementGroupId", "managementGroupId"),
		resourceids.StaticSegment("staticProviders2", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftPolicyInsights", "Microsoft.PolicyInsights", "Microsoft.PolicyInsights"),
		resourceids.StaticSegment("staticRemediations", "remediations", "remediations"),
		resourceids.UserSpecifiedSegment("remediationName", "remediationName"),
	}
}

// String returns a human-readable description of this Providers 2 Remediation ID
func (id Providers2RemediationId) String() string {
	components := []string{
		fmt.Sprintf("Management Group: %q", id.ManagementGroupId),
		fmt.Sprintf("Remediation Name: %q", id.RemediationName),
	}
	return fmt.Sprintf("Providers 2 Remediation (%s)", strings.Join(components, "\n"))
}
