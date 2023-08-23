package runbook

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type RunbookProvisioningState string

const (
	RunbookProvisioningStateSucceeded RunbookProvisioningState = "Succeeded"
)

func PossibleValuesForRunbookProvisioningState() []string {
	return []string{
		string(RunbookProvisioningStateSucceeded),
	}
}

func (s *RunbookProvisioningState) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseRunbookProvisioningState(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseRunbookProvisioningState(input string) (*RunbookProvisioningState, error) {
	vals := map[string]RunbookProvisioningState{
		"succeeded": RunbookProvisioningStateSucceeded,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := RunbookProvisioningState(input)
	return &out, nil
}

type RunbookState string

const (
	RunbookStateEdit      RunbookState = "Edit"
	RunbookStateNew       RunbookState = "New"
	RunbookStatePublished RunbookState = "Published"
)

func PossibleValuesForRunbookState() []string {
	return []string{
		string(RunbookStateEdit),
		string(RunbookStateNew),
		string(RunbookStatePublished),
	}
}

func (s *RunbookState) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseRunbookState(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseRunbookState(input string) (*RunbookState, error) {
	vals := map[string]RunbookState{
		"edit":      RunbookStateEdit,
		"new":       RunbookStateNew,
		"published": RunbookStatePublished,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := RunbookState(input)
	return &out, nil
}

type RunbookTypeEnum string

const (
	RunbookTypeEnumGraph                   RunbookTypeEnum = "Graph"
	RunbookTypeEnumGraphPowerShell         RunbookTypeEnum = "GraphPowerShell"
	RunbookTypeEnumGraphPowerShellWorkflow RunbookTypeEnum = "GraphPowerShellWorkflow"
	RunbookTypeEnumPowerShell              RunbookTypeEnum = "PowerShell"
	RunbookTypeEnumPowerShellWorkflow      RunbookTypeEnum = "PowerShellWorkflow"
	RunbookTypeEnumPythonThree             RunbookTypeEnum = "Python3"
	RunbookTypeEnumPythonTwo               RunbookTypeEnum = "Python2"
	RunbookTypeEnumScript                  RunbookTypeEnum = "Script"
)

func PossibleValuesForRunbookTypeEnum() []string {
	return []string{
		string(RunbookTypeEnumGraph),
		string(RunbookTypeEnumGraphPowerShell),
		string(RunbookTypeEnumGraphPowerShellWorkflow),
		string(RunbookTypeEnumPowerShell),
		string(RunbookTypeEnumPowerShellWorkflow),
		string(RunbookTypeEnumPythonThree),
		string(RunbookTypeEnumPythonTwo),
		string(RunbookTypeEnumScript),
	}
}

func (s *RunbookTypeEnum) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseRunbookTypeEnum(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseRunbookTypeEnum(input string) (*RunbookTypeEnum, error) {
	vals := map[string]RunbookTypeEnum{
		"graph":                   RunbookTypeEnumGraph,
		"graphpowershell":         RunbookTypeEnumGraphPowerShell,
		"graphpowershellworkflow": RunbookTypeEnumGraphPowerShellWorkflow,
		"powershell":              RunbookTypeEnumPowerShell,
		"powershellworkflow":      RunbookTypeEnumPowerShellWorkflow,
		"python3":                 RunbookTypeEnumPythonThree,
		"python2":                 RunbookTypeEnumPythonTwo,
		"script":                  RunbookTypeEnumScript,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := RunbookTypeEnum(input)
	return &out, nil
}
