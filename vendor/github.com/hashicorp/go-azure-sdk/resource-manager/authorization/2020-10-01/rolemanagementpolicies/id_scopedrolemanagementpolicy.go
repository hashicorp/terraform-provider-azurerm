package rolemanagementpolicies

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/recaser"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

func init() {
	recaser.RegisterResourceId(&ScopedRoleManagementPolicyId{})
}

var _ resourceids.ResourceId = &ScopedRoleManagementPolicyId{}

// ScopedRoleManagementPolicyId is a struct representing the Resource ID for a Scoped Role Management Policy
type ScopedRoleManagementPolicyId struct {
	Scope                    string
	RoleManagementPolicyName string
}

// NewScopedRoleManagementPolicyID returns a new ScopedRoleManagementPolicyId struct
func NewScopedRoleManagementPolicyID(scope string, roleManagementPolicyName string) ScopedRoleManagementPolicyId {
	return ScopedRoleManagementPolicyId{
		Scope:                    scope,
		RoleManagementPolicyName: roleManagementPolicyName,
	}
}

// ParseScopedRoleManagementPolicyID parses 'input' into a ScopedRoleManagementPolicyId
func ParseScopedRoleManagementPolicyID(input string) (*ScopedRoleManagementPolicyId, error) {
	parser := resourceids.NewParserFromResourceIdType(&ScopedRoleManagementPolicyId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := ScopedRoleManagementPolicyId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParseScopedRoleManagementPolicyIDInsensitively parses 'input' case-insensitively into a ScopedRoleManagementPolicyId
// note: this method should only be used for API response data and not user input
func ParseScopedRoleManagementPolicyIDInsensitively(input string) (*ScopedRoleManagementPolicyId, error) {
	parser := resourceids.NewParserFromResourceIdType(&ScopedRoleManagementPolicyId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := ScopedRoleManagementPolicyId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *ScopedRoleManagementPolicyId) FromParseResult(input resourceids.ParseResult) error {
	var ok bool

	if id.Scope, ok = input.Parsed["scope"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "scope", input)
	}

	if id.RoleManagementPolicyName, ok = input.Parsed["roleManagementPolicyName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "roleManagementPolicyName", input)
	}

	return nil
}

// ValidateScopedRoleManagementPolicyID checks that 'input' can be parsed as a Scoped Role Management Policy ID
func ValidateScopedRoleManagementPolicyID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseScopedRoleManagementPolicyID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Scoped Role Management Policy ID
func (id ScopedRoleManagementPolicyId) ID() string {
	fmtString := "/%s/providers/Microsoft.Authorization/roleManagementPolicies/%s"
	return fmt.Sprintf(fmtString, strings.TrimPrefix(id.Scope, "/"), id.RoleManagementPolicyName)
}

// Segments returns a slice of Resource ID Segments which comprise this Scoped Role Management Policy ID
func (id ScopedRoleManagementPolicyId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.ScopeSegment("scope", "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/some-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftAuthorization", "Microsoft.Authorization", "Microsoft.Authorization"),
		resourceids.StaticSegment("staticRoleManagementPolicies", "roleManagementPolicies", "roleManagementPolicies"),
		resourceids.UserSpecifiedSegment("roleManagementPolicyName", "roleManagementPolicyName"),
	}
}

// String returns a human-readable description of this Scoped Role Management Policy ID
func (id ScopedRoleManagementPolicyId) String() string {
	components := []string{
		fmt.Sprintf("Scope: %q", id.Scope),
		fmt.Sprintf("Role Management Policy Name: %q", id.RoleManagementPolicyName),
	}
	return fmt.Sprintf("Scoped Role Management Policy (%s)", strings.Join(components, "\n"))
}
