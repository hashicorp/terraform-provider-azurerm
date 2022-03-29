package parse

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

var _ resourceids.Id = LighthouseAssignmentId{}

type LighthouseAssignmentId struct {
	Scope string
	Name  string
}

func NewLighthouseAssignmentID(scope, name string) LighthouseAssignmentId {
	return LighthouseAssignmentId{
		Scope: scope,
		Name:  name,
	}
}

func (id LighthouseAssignmentId) ID() string {
	fmtStr := "%s/providers/Microsoft.ManagedServices/registrationAssignments/%s"
	return fmt.Sprintf(fmtStr, id.Scope, id.Name)
}

func (id LighthouseAssignmentId) String() string {
	segments := []string{
		fmt.Sprintf("Name %q", id.Name),
		fmt.Sprintf("Scope %q", id.Scope),
	}
	segmentsStr := strings.Join(segments, " / ")
	return fmt.Sprintf("%s: (%s)", "Lighthouse Assignment", segmentsStr)
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
