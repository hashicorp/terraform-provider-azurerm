// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package parse

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/desktopvirtualization/2022-02-10-preview/applicationgroup"
	"github.com/hashicorp/go-azure-sdk/resource-manager/desktopvirtualization/2022-02-10-preview/workspace"
)

var _ resourceids.Id = WorkspaceApplicationGroupAssociationId{}

type WorkspaceApplicationGroupAssociationId struct {
	Workspace        workspace.WorkspaceId
	ApplicationGroup applicationgroup.ApplicationGroupId
}

func (id WorkspaceApplicationGroupAssociationId) String() string {
	components := []string{
		fmt.Sprintf("Workspace %s", id.Workspace.String()),
		fmt.Sprintf("Application Group %s", id.ApplicationGroup.String()),
	}
	return fmt.Sprintf("Workspace Application Group Association %s", strings.Join(components, " / "))
}

func (id WorkspaceApplicationGroupAssociationId) ID() string {
	workspaceId := id.Workspace.ID()
	applicationGroupId := id.ApplicationGroup.ID()
	return fmt.Sprintf("%s|%s", workspaceId, applicationGroupId)
}

func NewWorkspaceApplicationGroupAssociationId(workspace workspace.WorkspaceId, applicationGroup applicationgroup.ApplicationGroupId) WorkspaceApplicationGroupAssociationId {
	return WorkspaceApplicationGroupAssociationId{
		Workspace:        workspace,
		ApplicationGroup: applicationGroup,
	}
}

func WorkspaceApplicationGroupAssociationID(input string) (*WorkspaceApplicationGroupAssociationId, error) {
	segments := strings.Split(input, "|")
	if len(segments) != 2 {
		return nil, fmt.Errorf("expected an ID in the format {workspaceID}|{applicationGroupID} but got %q", input)
	}

	workspaceId, err := workspace.ParseWorkspaceID(segments[0])
	if err != nil {
		return nil, fmt.Errorf("parsing Workspace ID for Workspace/Application Group Association %q: %+v", segments[0], err)
	}

	applicationGroupId, err := applicationgroup.ParseApplicationGroupID(segments[1])
	if err != nil {
		return nil, fmt.Errorf("parsing Application Group ID for Workspace/Application Group Association %q: %+v", segments[1], err)
	}

	return &WorkspaceApplicationGroupAssociationId{
		Workspace:        *workspaceId,
		ApplicationGroup: *applicationGroupId,
	}, nil
}

func WorkspaceApplicationGroupAssociationIDInsensitively(input string) (*WorkspaceApplicationGroupAssociationId, error) {
	segments := strings.Split(input, "|")
	if len(segments) != 2 {
		return nil, fmt.Errorf("expected an ID in the format {workspaceID}|{applicationGroupID} but got %q", input)
	}

	workspaceId, err := workspace.ParseWorkspaceIDInsensitively(segments[0])
	if err != nil {
		return nil, fmt.Errorf("parsing Workspace ID for Workspace/Application Group Association %q: %+v", segments[0], err)
	}

	applicationGroupId, err := applicationgroup.ParseApplicationGroupIDInsensitively(segments[1])
	if err != nil {
		return nil, fmt.Errorf("parsing Application Group ID for Workspace/Application Group Association %q: %+v", segments[1], err)
	}

	return &WorkspaceApplicationGroupAssociationId{
		Workspace:        *workspaceId,
		ApplicationGroup: *applicationGroupId,
	}, nil
}
