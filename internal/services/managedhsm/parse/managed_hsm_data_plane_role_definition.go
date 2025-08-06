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
var _ resourceids.Id = ManagedHSMDataPlaneRoleDefinitionId{}

type ManagedHSMDataPlaneRoleDefinitionId struct {
	ManagedHSMName     string
	DomainSuffix       string
	Scope              string
	RoleDefinitionName string
}

func NewManagedHSMDataPlaneRoleDefinitionID(managedHSMName, domainSuffix, scope, roleDefinitionName string) ManagedHSMDataPlaneRoleDefinitionId {
	return ManagedHSMDataPlaneRoleDefinitionId{
		ManagedHSMName:     managedHSMName,
		DomainSuffix:       domainSuffix,
		Scope:              scope,
		RoleDefinitionName: roleDefinitionName,
	}
}

func ManagedHSMDataPlaneRoleDefinitionID(input string, domainSuffix *string) (*ManagedHSMDataPlaneRoleDefinitionId, error) {
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

	// and then the Scope and RoleDefinitionName from the URI
	if !strings.HasPrefix(uri.Path, "/") {
		// sanity-checking, but we're expecting at least a `//` on the front
		return nil, fmt.Errorf("expected the path to start with `//` but got %q", uri.Path)
	}
	pathRaw := strings.TrimPrefix(uri.Path, "/")
	path, err := parseManagedHSMRoleDefinitionFromPath(pathRaw)
	if err != nil {
		return nil, err
	}

	return &ManagedHSMDataPlaneRoleDefinitionId{
		ManagedHSMName:     endpoint.ManagedHSMName,
		DomainSuffix:       endpoint.DomainSuffix,
		Scope:              path.scope,
		RoleDefinitionName: path.roleDefinitionName,
	}, nil
}

func (id ManagedHSMDataPlaneRoleDefinitionId) BaseURI() string {
	return ManagedHSMDataPlaneEndpoint{
		ManagedHSMName: id.ManagedHSMName,
		DomainSuffix:   id.DomainSuffix,
	}.BaseURI()
}

func (id ManagedHSMDataPlaneRoleDefinitionId) ID() string {
	path := managedHSMRoleDefinitionPathId{
		scope:              id.Scope,
		roleDefinitionName: id.RoleDefinitionName,
	}
	return fmt.Sprintf("https://%s.%s%s", id.ManagedHSMName, id.DomainSuffix, path.ID())
}

func (id ManagedHSMDataPlaneRoleDefinitionId) String() string {
	components := []string{
		fmt.Sprintf("Managed HSM Name %q", id.ManagedHSMName),
		fmt.Sprintf("Domain Suffix Name %q", id.DomainSuffix),
		fmt.Sprintf("Scope %q", id.Scope),
		fmt.Sprintf("Role Definition Name %q", id.RoleDefinitionName),
	}
	return fmt.Sprintf("Managed HSM Data Plane Role Definition ID (%s)", strings.Join(components, " | "))
}

var _ resourceids.ResourceId = &managedHSMRoleDefinitionPathId{}

// managedHSMRoleDefinitionPathId parses the Path component
type managedHSMRoleDefinitionPathId struct {
	scope              string
	roleDefinitionName string
}

func parseManagedHSMRoleDefinitionFromPath(input string) (*managedHSMRoleDefinitionPathId, error) {
	id := managedHSMRoleDefinitionPathId{}
	parsed, err := resourceids.NewParserFromResourceIdType(&id).Parse(input, false)
	if err != nil {
		return nil, err
	}

	if err := id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *managedHSMRoleDefinitionPathId) FromParseResult(parsed resourceids.ParseResult) error {
	var ok bool
	if id.scope, ok = parsed.Parsed["scope"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "scope", parsed)
	}
	if id.roleDefinitionName, ok = parsed.Parsed["roleDefinitionName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "roleDefinitionName", parsed)
	}

	return nil
}

func (id *managedHSMRoleDefinitionPathId) ID() string {
	return fmt.Sprintf("/%s/providers/Microsoft.Authorization/roleDefinitions/%s", id.scope, id.roleDefinitionName)
}

func (id *managedHSMRoleDefinitionPathId) String() string {
	return fmt.Sprintf("Role Definition %q (Scope %q)", id.roleDefinitionName, id.scope)
}

func (id *managedHSMRoleDefinitionPathId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		// /{scope}/providers/Microsoft.Authorization/roleDefinitions/{roleDefinitionName}
		resourceids.ScopeSegment("scope", "/"),
		resourceids.StaticSegment("providers", "providers", "providers"),
		resourceids.ResourceProviderSegment("resourceProvider", "Microsoft.Authorization", "Microsoft.Authorization"),
		resourceids.StaticSegment("roleDefinitions", "roleDefinitions", "roleDefinitions"),
		resourceids.UserSpecifiedSegment("roleDefinitionName", "roleDefinitionValue"),
	}
}
