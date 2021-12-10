package parse

import (
	"fmt"
	"strings"
)

type LighthouseDefinitionId struct {
	Scope                  string
	LighthouseDefinitionID string
}

func LighthouseDefinitionID(id string) (*LighthouseDefinitionId, error) {
	segments := strings.Split(id, "/providers/Microsoft.ManagedServices/registrationDefinitions/")

	if len(segments) != 2 {
		return nil, fmt.Errorf("Expected ID to be in the format `{scope}/providers/Microsoft.ManagedServices/registrationDefinitions/{name} - got %d segments", len(segments))
	}

	lighthouseDefinitionId := LighthouseDefinitionId{
		Scope:                  segments[0],
		LighthouseDefinitionID: segments[1],
	}

	return &lighthouseDefinitionId, nil
}
