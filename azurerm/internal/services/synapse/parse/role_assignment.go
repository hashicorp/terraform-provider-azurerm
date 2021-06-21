package parse

import (
	"fmt"
	"strings"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/resourceid"
)

var _ resourceid.Formatter = RoleAssignmentId{}

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
