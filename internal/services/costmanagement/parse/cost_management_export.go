// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package parse

import (
	"fmt"
	"regexp"
	"strings"
)

type CostManagementExportId struct {
	Name  string
	Scope string
}

func (id CostManagementExportId) String() string {
	segments := []string{
		fmt.Sprintf("Cost Management Export Name %q", id.Name),
		fmt.Sprintf("Scope %q", id.Scope),
	}
	segmentsStr := strings.Join(segments, " / ")
	return fmt.Sprintf("%s: (%s)", "Cost Management Export ID", segmentsStr)
}

func (id CostManagementExportId) ID() string {
	fmtString := "%s/providers/Microsoft.CostManagement/exports/%s"
	return fmt.Sprintf(fmtString, id.Scope, id.Name)
}

func NewCostManagementExportId(scope, name string) CostManagementExportId {
	return CostManagementExportId{
		Name:  name,
		Scope: scope,
	}
}

func CostManagementExportID(input string) (*CostManagementExportId, error) {
	// in general, the id of a assignment should be:
	// {scope}/providers/Microsoft.CostManagement/exports/{name}
	regex := regexp.MustCompile(`/providers/Microsoft\.CostManagement/exports/`)
	if !regex.MatchString(input) {
		return nil, fmt.Errorf("unable to parse Cost Management Export ID %q", input)
	}

	segments := regex.Split(input, -1)

	if len(segments) != 2 {
		return nil, fmt.Errorf("unable to parse Cost Management Export ID %q: Expected 2 segments after split", input)
	}

	scope := segments[0]
	name := segments[1]
	if name == "" {
		return nil, fmt.Errorf("unable to parse Cost Management Export ID %q: export name is empty", input)
	}

	return &CostManagementExportId{
		Name:  name,
		Scope: scope,
	}, nil
}
