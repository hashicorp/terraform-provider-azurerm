// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package parse

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

var _ resourceids.Id = RoleAssignmentId{}

type RoleAssignmentId struct {
	Scope                 string
	DataPlaneAssignmentId string
}

func NewRoleAssignmentId(scope string, dataPlaneAssignmentId string) RoleAssignmentId {
	return RoleAssignmentId{
		Scope:                 scope,
		DataPlaneAssignmentId: dataPlaneAssignmentId,
	}
}

func (id RoleAssignmentId) String() string {
	components := []string{
		fmt.Sprintf("Scope %q", id.Scope),
		fmt.Sprintf("Data Plane Assignment ID %q", id.DataPlaneAssignmentId),
	}
	return fmt.Sprintf("Role Assignment %s", strings.Join(components, " / "))
}

func (id RoleAssignmentId) ID() string {
	return fmt.Sprintf("%s|%s", id.Scope, id.DataPlaneAssignmentId)
}

func RoleAssignmentID(input string) (*RoleAssignmentId, error) {
	segments := strings.Split(input, "|")
	if len(segments) != 2 {
		return nil, fmt.Errorf("expected an ID in the format `{workspaceId}|{id} but got %q", input)
	}

	return &RoleAssignmentId{
		Scope:                 segments[0],
		DataPlaneAssignmentId: segments[1],
	}, nil
}

func SynapseScope(synapseScope string) (string, string, error) {
	workspaceId, err := WorkspaceID(synapseScope)
	if err == nil {
		return workspaceId.Name, fmt.Sprintf("workspaces/%s", workspaceId.Name), nil
	}

	sparkPoolID, err := SparkPoolID(synapseScope)
	if err == nil {
		return sparkPoolID.WorkspaceName, fmt.Sprintf("workspaces/%s/bigDataPools/%s", sparkPoolID.WorkspaceName, sparkPoolID.BigDataPoolName), nil
	}

	return "", "", fmt.Errorf("synapseScope format error")
}
