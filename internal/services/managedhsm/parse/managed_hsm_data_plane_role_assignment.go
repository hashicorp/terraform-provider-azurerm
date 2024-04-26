// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package parse

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// @tombuildsstuff: note this intentionally implements `resourceids.Id` and not `resourceids.ResourceId` since we
// can't currently represent the full URI in the Segment types until https://github.com/hashicorp/go-azure-helpers/issues/187
// is implemented (including having Pandora be aware of the new segment types)
var _ resourceids.Id = ManagedHSMDataPlaneRoleAssignmentId{}

type ManagedHSMDataPlaneRoleAssignmentId struct {
	ManagedHSMName     string
	DomainSuffix       string
	Scope              string
	RoleAssignmentName string
}

func NewManagedHSMDataPlaneRoleAssignmentID(managedHSMName, domainSuffix, scope, roleAssignmentName string) ManagedHSMDataPlaneRoleAssignmentId {
	return ManagedHSMDataPlaneRoleAssignmentId{
		ManagedHSMName:     managedHSMName,
		DomainSuffix:       domainSuffix,
		Scope:              scope,
		RoleAssignmentName: roleAssignmentName,
	}
}

func ManagedHSMDataPlaneRoleAssignmentID(input string, domainSuffix *string) (*ManagedHSMDataPlaneRoleAssignmentId, error) {
	if input == "" {
		return nil, fmt.Errorf("`input` was empty")
	}
	if domainSuffix != nil && !strings.HasPrefix(strings.ToLower(*domainSuffix), "managedhsm.") {
		return nil, fmt.Errorf("internal-error: the domainSuffix for Managed HSM %q didn't contain `managedhsm.`", *domainSuffix)
	}

	uri, err := url.Parse(input)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	// we need the ManagedHSMName and DomainSuffix from the Host
	endpoint, err := parseDataPlaneEndpoint(uri, domainSuffix)
	if err != nil {
		// intentionally not wrapping this
		return nil, err
	}

	// and then the Scope and RoleAssignmentName from the URI
	if !strings.HasPrefix(uri.Path, "/") {
		// sanity-checking, but we're expecting at least a `//` on the front
		return nil, fmt.Errorf("expected the path to start with `//` but got %q", uri.Path)
	}
	pathRaw := strings.TrimPrefix(uri.Path, "/")
	path, err := parseManagedHSMRoleAssignmentFromPath(pathRaw)
	if err != nil {
		return nil, err
	}

	return &ManagedHSMDataPlaneRoleAssignmentId{
		ManagedHSMName:     endpoint.ManagedHSMName,
		DomainSuffix:       endpoint.DomainSuffix,
		Scope:              path.scope,
		RoleAssignmentName: path.roleAssignmentName,
	}, nil
}

func (id ManagedHSMDataPlaneRoleAssignmentId) BaseURI() string {
	return ManagedHSMDataPlaneEndpoint{
		ManagedHSMName: id.ManagedHSMName,
		DomainSuffix:   id.DomainSuffix,
	}.BaseURI()
}

func (id ManagedHSMDataPlaneRoleAssignmentId) ID() string {
	path := managedHSMRoleAssignmentPathId{
		scope:              id.Scope,
		roleAssignmentName: id.RoleAssignmentName,
	}
	return fmt.Sprintf("https://%s.%s%s", id.ManagedHSMName, id.DomainSuffix, path.ID())
}

func (id ManagedHSMDataPlaneRoleAssignmentId) String() string {
	components := []string{
		fmt.Sprintf("Managed HSM Name %q", id.ManagedHSMName),
		fmt.Sprintf("Domain Suffix Name %q", id.DomainSuffix),
		fmt.Sprintf("Scope %q", id.Scope),
		fmt.Sprintf("Role Assignment Name %q", id.RoleAssignmentName),
	}
	return fmt.Sprintf("Managed HSM Data Plane Role Assignment ID (%s)", strings.Join(components, " | "))
}

var _ resourceids.ResourceId = &managedHSMRoleAssignmentPathId{}

// managedHSMRoleAssignmentPathId parses the Path component
type managedHSMRoleAssignmentPathId struct {
	scope              string
	roleAssignmentName string
}

func parseManagedHSMRoleAssignmentFromPath(input string) (*managedHSMRoleAssignmentPathId, error) {
	id := managedHSMRoleAssignmentPathId{}
	parsed, err := resourceids.NewParserFromResourceIdType(&id).Parse(input, false)
	if err != nil {
		return nil, err
	}

	if err := id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *managedHSMRoleAssignmentPathId) FromParseResult(parsed resourceids.ParseResult) error {
	var ok bool
	if id.scope, ok = parsed.Parsed["scope"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "scope", parsed)
	}
	if id.roleAssignmentName, ok = parsed.Parsed["roleAssignmentName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "roleAssignmentName", parsed)
	}

	return nil
}

func (id *managedHSMRoleAssignmentPathId) ID() string {
	return fmt.Sprintf("/%s/providers/Microsoft.Authorization/roleAssignments/%s", id.scope, id.roleAssignmentName)
}

func (id *managedHSMRoleAssignmentPathId) String() string {
	return fmt.Sprintf("Role Assignment %q (Scope %q)", id.roleAssignmentName, id.scope)
}

func (id *managedHSMRoleAssignmentPathId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		// /{scope}/providers/Microsoft.Authorization/roleAssignments/{roleAssignmentName}
		resourceids.ScopeSegment("scope", "/"),
		resourceids.StaticSegment("providers", "providers", "providers"),
		resourceids.ResourceProviderSegment("resourceProvider", "Microsoft.Authorization", "Microsoft.Authorization"),
		resourceids.StaticSegment("roleAssignments", "roleAssignments", "roleAssignments"),
		resourceids.UserSpecifiedSegment("roleAssignmentName", "roleAssignmentValue"),
	}
}
