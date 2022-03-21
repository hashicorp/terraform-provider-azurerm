package parse

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

var _ resourceids.Id = LighthouseDefinitionId{}

type LighthouseDefinitionId struct {
	Scope                  string
	LighthouseDefinitionID string
}

func NewLighthouseDefinitionID(scope, lighthouseDefinitionId string) LighthouseDefinitionId {
	return LighthouseDefinitionId{
		Scope:                  scope,
		LighthouseDefinitionID: lighthouseDefinitionId,
	}
}

func (id LighthouseDefinitionId) ID() string {
	fmtStr := "%s/providers/Microsoft.ManagedServices/registrationDefinitions/%s"
	return fmt.Sprintf(fmtStr, id.Scope, id.LighthouseDefinitionID)
}

func (id LighthouseDefinitionId) String() string {
	segments := []string{
		fmt.Sprintf("Lighthouse Definition ID %q", id.LighthouseDefinitionID),
		fmt.Sprintf("Scope %q", id.Scope),
	}
	segmentsStr := strings.Join(segments, " / ")
	return fmt.Sprintf("%s: (%s)", "Lighthouse Assignment", segmentsStr)
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
