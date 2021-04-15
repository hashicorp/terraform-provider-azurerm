package parse

import (
	"fmt"
	"strings"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
)

type AssessmentId struct {
	TargetResourceID string
	Name             string
}

func NewAssessmentID(targetResourceID, name string) AssessmentId {
	return AssessmentId{
		TargetResourceID: targetResourceID,
		Name:             name,
	}
}

func (id AssessmentId) ID() string {
	fmtString := "%s/providers/Microsoft.Security/assessments/%s"
	return fmt.Sprintf(fmtString, id.TargetResourceID, id.Name)
}

func AssessmentID(input string) (*AssessmentId, error) {
	id, err := azure.ParseAzureResourceID(input)
	if err != nil {
		return nil, fmt.Errorf("parsing Security Assessment ID %q: %+v", input, err)
	}

	parts := strings.Split(input, "/providers/Microsoft.Security/assessments/")
	if len(parts) != 2 {
		return nil, fmt.Errorf("parsing Security Assessment ID: %q", id)
	}

	return &AssessmentId{
		TargetResourceID: parts[0],
		Name:             parts[1],
	}, nil
}
