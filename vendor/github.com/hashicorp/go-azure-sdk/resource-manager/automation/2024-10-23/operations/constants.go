package operations

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type GraphRunbookType string

const (
	GraphRunbookTypeGraphPowerShell         GraphRunbookType = "GraphPowerShell"
	GraphRunbookTypeGraphPowerShellWorkflow GraphRunbookType = "GraphPowerShellWorkflow"
)

func PossibleValuesForGraphRunbookType() []string {
	return []string{
		string(GraphRunbookTypeGraphPowerShell),
		string(GraphRunbookTypeGraphPowerShellWorkflow),
	}
}

func (s *GraphRunbookType) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseGraphRunbookType(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseGraphRunbookType(input string) (*GraphRunbookType, error) {
	vals := map[string]GraphRunbookType{
		"graphpowershell":         GraphRunbookTypeGraphPowerShell,
		"graphpowershellworkflow": GraphRunbookTypeGraphPowerShellWorkflow,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := GraphRunbookType(input)
	return &out, nil
}
