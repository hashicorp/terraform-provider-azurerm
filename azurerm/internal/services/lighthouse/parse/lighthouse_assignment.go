package parse

import (
	"fmt"
	"strings"
)

type LighthouseAssignmentId struct {
	Scope string
	Name  string
}

func LighthouseAssignmentID(id string) (*LighthouseAssignmentId, error) {
	segments := strings.Split(id, "/providers/Microsoft.ManagedServices/registrationAssignments/")
	if len(segments) != 2 {
		return nil, fmt.Errorf("Expected ID to be in the format `{scope}/providers/Microsoft.ManagedServices/registrationAssignments/{name} - got %d segments", len(segments))
	}

	lighthouseAssignmentId := LighthouseAssignmentId{
		Scope: segments[0],
		Name:  segments[1],
	}

	return &lighthouseAssignmentId, nil
}
