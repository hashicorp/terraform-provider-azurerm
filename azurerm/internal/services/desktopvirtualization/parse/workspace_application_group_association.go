package parse

import (
	"fmt"
	"strings"
)

type WorkspaceApplicationGroupAssociationId struct {
	Workspace        WorkspaceId
	ApplicationGroup ApplicationGroupId
}

func (id WorkspaceApplicationGroupAssociationId) ID(subscriptionId string) string {
	workspaceId := id.Workspace.ID(subscriptionId)
	applicationGroupId := id.ApplicationGroup.ID(subscriptionId)
	return fmt.Sprintf("%s|%s", workspaceId, applicationGroupId)
}

func NewWorkspaceApplicationGroupAssociationId(workspace WorkspaceId, applicationGroup ApplicationGroupId) WorkspaceApplicationGroupAssociationId {
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

	workspaceId, err := WorkspaceID(segments[0])
	if err != nil {
		return nil, fmt.Errorf("parsing Workspace ID for Workspace/Application Group Association %q: %+v", segments[0], err)
	}

	applicationGroupId, err := ApplicationGroupID(segments[1])
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

	workspaceId, err := WorkspaceID(segments[0])
	if err != nil {
		return nil, fmt.Errorf("parsing Workspace ID for Workspace/Application Group Association %q: %+v", segments[0], err)
	}

	applicationGroupId, err := ApplicationGroupIDInsensitively(segments[1])
	if err != nil {
		return nil, fmt.Errorf("parsing Application Group ID for Workspace/Application Group Association %q: %+v", segments[1], err)
	}

	return &WorkspaceApplicationGroupAssociationId{
		Workspace:        *workspaceId,
		ApplicationGroup: *applicationGroupId,
	}, nil
}
