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
	recaser.RegisterResourceId(&ScopedRemediationId{})
}

var _ resourceids.ResourceId = &ScopedRemediationId{}

// ScopedRemediationId is a struct representing the Resource ID for a Scoped Remediation
type ScopedRemediationId struct {
	ResourceId      string
	RemediationName string
}

// NewScopedRemediationID returns a new ScopedRemediationId struct
func NewScopedRemediationID(resourceId string, remediationName string) ScopedRemediationId {
	return ScopedRemediationId{
		ResourceId:      resourceId,
		RemediationName: remediationName,
	}
}

// ParseScopedRemediationID parses 'input' into a ScopedRemediationId
func ParseScopedRemediationID(input string) (*ScopedRemediationId, error) {
	parser := resourceids.NewParserFromResourceIdType(&ScopedRemediationId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := ScopedRemediationId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParseScopedRemediationIDInsensitively parses 'input' case-insensitively into a ScopedRemediationId
// note: this method should only be used for API response data and not user input
func ParseScopedRemediationIDInsensitively(input string) (*ScopedRemediationId, error) {
	parser := resourceids.NewParserFromResourceIdType(&ScopedRemediationId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := ScopedRemediationId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *ScopedRemediationId) FromParseResult(input resourceids.ParseResult) error {
	var ok bool

	if id.ResourceId, ok = input.Parsed["resourceId"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "resourceId", input)
	}

	if id.RemediationName, ok = input.Parsed["remediationName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "remediationName", input)
	}

	return nil
}

// ValidateScopedRemediationID checks that 'input' can be parsed as a Scoped Remediation ID
func ValidateScopedRemediationID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseScopedRemediationID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Scoped Remediation ID
func (id ScopedRemediationId) ID() string {
	fmtString := "/%s/providers/Microsoft.PolicyInsights/remediations/%s"
	return fmt.Sprintf(fmtString, strings.TrimPrefix(id.ResourceId, "/"), id.RemediationName)
}

// Segments returns a slice of Resource ID Segments which comprise this Scoped Remediation ID
func (id ScopedRemediationId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.ScopeSegment("resourceId", "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/some-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftPolicyInsights", "Microsoft.PolicyInsights", "Microsoft.PolicyInsights"),
		resourceids.StaticSegment("staticRemediations", "remediations", "remediations"),
		resourceids.UserSpecifiedSegment("remediationName", "remediationName"),
	}
}

// String returns a human-readable description of this Scoped Remediation ID
func (id ScopedRemediationId) String() string {
	components := []string{
		fmt.Sprintf("Resource: %q", id.ResourceId),
		fmt.Sprintf("Remediation Name: %q", id.RemediationName),
	}
	return fmt.Sprintf("Scoped Remediation (%s)", strings.Join(components, "\n"))
}
