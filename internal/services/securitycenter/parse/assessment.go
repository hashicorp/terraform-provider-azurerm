// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package parse

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
)

var _ resourceids.Id = AssessmentId{}

type AssessmentId struct {
	TargetResourceID string
	Name             string
}

func (id AssessmentId) String() string {
	components := []string{
		fmt.Sprintf("Target Resource %q", id.TargetResourceID),
		fmt.Sprintf("Name %q", id.Name),
	}
	return fmt.Sprintf("Assessment %s", strings.Join(components, " / "))
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
