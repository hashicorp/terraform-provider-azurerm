package parse

import (
	"fmt"
	"strings"
)

type VirtualDesktopWorkspaceApplicationGroupAssociationId struct {
	Workspace        VirtualDesktopWorkspaceId
	ApplicationGroup VirtualDesktopApplicationGroupId
}

func (id VirtualDesktopWorkspaceApplicationGroupAssociationId) ID(subscriptionId string) string {
	workspaceId := id.Workspace.ID(subscriptionId)
	applicationGroupId := id.ApplicationGroup.ID(subscriptionId)
	return fmt.Sprintf("%s|%s", workspaceId, applicationGroupId)
}

func NewVirtualDesktopWorkspaceApplicationGroupAssociationId(workspace VirtualDesktopWorkspaceId, applicationGroup VirtualDesktopApplicationGroupId) VirtualDesktopWorkspaceApplicationGroupAssociationId {
	return VirtualDesktopWorkspaceApplicationGroupAssociationId{
		Workspace:        workspace,
		ApplicationGroup: applicationGroup,
	}
}

func VirtualDesktopWorkspaceApplicationGroupAssociationID(input string) (*VirtualDesktopWorkspaceApplicationGroupAssociationId, error) {
	segments := strings.Split(input, "|")
	if len(segments) != 2 {
		return nil, fmt.Errorf("expected an ID in the format {workspaceID}|{applicationGroupID} but got %q", input)
	}

	workspaceId, err := VirtualDesktopWorkspaceID(segments[0])
	if err != nil {
		return nil, fmt.Errorf("parsing Workspace ID for Workspace/Application Group Association %q: %+v", segments[0], err)
	}

	applicationGroupId, err := VirtualDesktopApplicationGroupID(segments[1])
	if err != nil {
		return nil, fmt.Errorf("parsing Application Group ID for Workspace/Application Group Association %q: %+v", segments[1], err)
	}

	return &VirtualDesktopWorkspaceApplicationGroupAssociationId{
		Workspace:        *workspaceId,
		ApplicationGroup: *applicationGroupId,
	}, nil
}
