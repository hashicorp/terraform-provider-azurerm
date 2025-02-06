// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package parse

import (
	"fmt"
	"regexp"
	"strings"
)

type ManagementGroupId struct {
	Name     string
	TenantID string
}

func ManagementGroupID(input string) (*ManagementGroupId, error) {
	regex := regexp.MustCompile(`^/providers/[Mm]icrosoft\.[Mm]anagement/[Mm]anagement[Gg]roups/`)
	if !regex.MatchString(input) {
		return nil, fmt.Errorf("Unable to parse Management Group ID %q", input)
	}

	// Split the input ID by the regex
	segments := regex.Split(input, -1)
	if len(segments) != 2 {
		return nil, fmt.Errorf("Unable to parse Management Group ID %q: expected id to have two segments after splitting", input)
	}

	groupID := segments[1]
	if groupID == "" {
		return nil, fmt.Errorf("unable to parse Management Group ID %q: management group name is empty", input)
	}
	if segments := strings.Split(groupID, "/"); len(segments) != 1 {
		return nil, fmt.Errorf("unable to parse Management Group ID %q: ID has extra segments", input)
	}

	id := ManagementGroupId{
		Name: groupID,
	}

	return &id, nil
}

func TenantScopedManagementGroupID(input string) (*ManagementGroupId, error) {
	regex := regexp.MustCompile(`^/tenants/.*-.*-.*-.*-.*/providers/Microsoft\.Management/managementGroups/`)
	if !regex.MatchString(input) {
		return nil, fmt.Errorf("Unable to parse Management Group ID for System Topic %q, format should look like '/tenants/<tenantID>/providers/Microsoft.Management/managementGroups/<management_group_name>'", input)
	}

	segments := strings.Split(input, "/")
	if len(segments) != 7 {
		return nil, fmt.Errorf("Unable to parse Management Group ID %q: expected id to have seven segments after splitting", input)
	}

	groupID := segments[len(segments)-1]
	if groupID == "" {
		return nil, fmt.Errorf("unable to parse Management Group ID %q: management group name is empty", input)
	}

	tenantID := segments[2]
	if tenantID == "" {
		return nil, fmt.Errorf("unable to parse Management Group ID %q: tenant id is empty", input)
	}

	id := ManagementGroupId{
		Name:     groupID,
		TenantID: tenantID,
	}

	return &id, nil
}

func NewManagementGroupId(managementGroupName string) ManagementGroupId {
	return ManagementGroupId{
		Name: managementGroupName,
	}
}

func NewTenantScopedManagementGroupID(tenantID, groupName string) ManagementGroupId {
	return ManagementGroupId{
		Name:     groupName,
		TenantID: tenantID,
	}
}

func (r ManagementGroupId) ID() string {
	managementGroupIdFmt := "/providers/Microsoft.Management/managementGroups/%s"
	return fmt.Sprintf(managementGroupIdFmt, r.Name)
}

func (r ManagementGroupId) TenantScopedID() string {
	managementGroupIDForSystemTopicFormat := "/tenants/%s/providers/Microsoft.Management/managementGroups/%s"

	return fmt.Sprintf(managementGroupIDForSystemTopicFormat, r.TenantID, r.Name)
}
