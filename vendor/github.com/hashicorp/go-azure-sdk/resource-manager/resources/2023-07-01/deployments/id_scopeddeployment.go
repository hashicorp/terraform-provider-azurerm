package deployments

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/recaser"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

func init() {
	recaser.RegisterResourceId(&ScopedDeploymentId{})
}

var _ resourceids.ResourceId = &ScopedDeploymentId{}

// ScopedDeploymentId is a struct representing the Resource ID for a Scoped Deployment
type ScopedDeploymentId struct {
	Scope          string
	DeploymentName string
}

// NewScopedDeploymentID returns a new ScopedDeploymentId struct
func NewScopedDeploymentID(scope string, deploymentName string) ScopedDeploymentId {
	return ScopedDeploymentId{
		Scope:          scope,
		DeploymentName: deploymentName,
	}
}

// ParseScopedDeploymentID parses 'input' into a ScopedDeploymentId
func ParseScopedDeploymentID(input string) (*ScopedDeploymentId, error) {
	parser := resourceids.NewParserFromResourceIdType(&ScopedDeploymentId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := ScopedDeploymentId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParseScopedDeploymentIDInsensitively parses 'input' case-insensitively into a ScopedDeploymentId
// note: this method should only be used for API response data and not user input
func ParseScopedDeploymentIDInsensitively(input string) (*ScopedDeploymentId, error) {
	parser := resourceids.NewParserFromResourceIdType(&ScopedDeploymentId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := ScopedDeploymentId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *ScopedDeploymentId) FromParseResult(input resourceids.ParseResult) error {
	var ok bool

	if id.Scope, ok = input.Parsed["scope"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "scope", input)
	}

	if id.DeploymentName, ok = input.Parsed["deploymentName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "deploymentName", input)
	}

	return nil
}

// ValidateScopedDeploymentID checks that 'input' can be parsed as a Scoped Deployment ID
func ValidateScopedDeploymentID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseScopedDeploymentID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Scoped Deployment ID
func (id ScopedDeploymentId) ID() string {
	fmtString := "/%s/providers/Microsoft.Resources/deployments/%s"
	return fmt.Sprintf(fmtString, strings.TrimPrefix(id.Scope, "/"), id.DeploymentName)
}

// Segments returns a slice of Resource ID Segments which comprise this Scoped Deployment ID
func (id ScopedDeploymentId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.ScopeSegment("scope", "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/some-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftResources", "Microsoft.Resources", "Microsoft.Resources"),
		resourceids.StaticSegment("staticDeployments", "deployments", "deployments"),
		resourceids.UserSpecifiedSegment("deploymentName", "deploymentName"),
	}
}

// String returns a human-readable description of this Scoped Deployment ID
func (id ScopedDeploymentId) String() string {
	components := []string{
		fmt.Sprintf("Scope: %q", id.Scope),
		fmt.Sprintf("Deployment Name: %q", id.DeploymentName),
	}
	return fmt.Sprintf("Scoped Deployment (%s)", strings.Join(components, "\n"))
}
