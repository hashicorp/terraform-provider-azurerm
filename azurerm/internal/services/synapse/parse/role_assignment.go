package parse

import (
	"fmt"
	"strings"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/resourceid"
)

var _ resourceid.Formatter = RoleAssignmentId{}

type RoleAssignmentId struct {
	Workspace             WorkspaceId
	DataPlaneAssignmentId string
}

func NewRoleAssignmentId(workspace WorkspaceId, dataPlaneAssignmentId string) RoleAssignmentId {
	return RoleAssignmentId{
		Workspace:             workspace,
		DataPlaneAssignmentId: dataPlaneAssignmentId,
	}
}

func (id RoleAssignmentId) ID() string {
	workspaceId := id.Workspace.ID("")
	return fmt.Sprintf("%s|%s", workspaceId, id.DataPlaneAssignmentId)
}

func RoleAssignmentID(input string) (*RoleAssignmentId, error) {
	segments := strings.Split(input, "|")
	if len(segments) != 2 {
		return nil, fmt.Errorf("expected an ID in the format `{workspaceId}|{id} but got %q", input)
	}

	workspaceId, err := WorkspaceID(segments[0])
	if err != nil {
		return nil, err
	}

	return &RoleAssignmentId{
		Workspace:             *workspaceId,
		DataPlaneAssignmentId: segments[1],
	}, nil
}
