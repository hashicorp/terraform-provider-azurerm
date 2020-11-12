package parse

import (
	"fmt"
	"strings"
)

type SynapseRoleAssignmentId struct {
	Workspace *SynapseWorkspaceId
	Id        string
}

func SynapseRoleAssignmentID(input string) (*SynapseRoleAssignmentId, error) {
	segments := strings.Split(input, "|")
	if len(segments) != 2 {
		return nil, fmt.Errorf("expected an ID in the format `{workspaceId}|{id} but got %q", input)
	}

	workspaceId, err := SynapseWorkspaceID(segments[0])
	if err != nil {
		return nil, err
	}

	return &SynapseRoleAssignmentId{
		Workspace: workspaceId,
		Id:        segments[1],
	}, nil
}
