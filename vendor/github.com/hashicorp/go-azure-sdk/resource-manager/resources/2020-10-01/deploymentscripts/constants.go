package deploymentscripts

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type CleanupOptions string

const (
	CleanupOptionsAlways       CleanupOptions = "Always"
	CleanupOptionsOnExpiration CleanupOptions = "OnExpiration"
	CleanupOptionsOnSuccess    CleanupOptions = "OnSuccess"
)

func PossibleValuesForCleanupOptions() []string {
	return []string{
		string(CleanupOptionsAlways),
		string(CleanupOptionsOnExpiration),
		string(CleanupOptionsOnSuccess),
	}
}

func (s *CleanupOptions) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseCleanupOptions(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseCleanupOptions(input string) (*CleanupOptions, error) {
	vals := map[string]CleanupOptions{
		"always":       CleanupOptionsAlways,
		"onexpiration": CleanupOptionsOnExpiration,
		"onsuccess":    CleanupOptionsOnSuccess,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := CleanupOptions(input)
	return &out, nil
}

type ScriptProvisioningState string

const (
	ScriptProvisioningStateCanceled              ScriptProvisioningState = "Canceled"
	ScriptProvisioningStateCreating              ScriptProvisioningState = "Creating"
	ScriptProvisioningStateFailed                ScriptProvisioningState = "Failed"
	ScriptProvisioningStateProvisioningResources ScriptProvisioningState = "ProvisioningResources"
	ScriptProvisioningStateRunning               ScriptProvisioningState = "Running"
	ScriptProvisioningStateSucceeded             ScriptProvisioningState = "Succeeded"
)

func PossibleValuesForScriptProvisioningState() []string {
	return []string{
		string(ScriptProvisioningStateCanceled),
		string(ScriptProvisioningStateCreating),
		string(ScriptProvisioningStateFailed),
		string(ScriptProvisioningStateProvisioningResources),
		string(ScriptProvisioningStateRunning),
		string(ScriptProvisioningStateSucceeded),
	}
}

func (s *ScriptProvisioningState) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseScriptProvisioningState(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseScriptProvisioningState(input string) (*ScriptProvisioningState, error) {
	vals := map[string]ScriptProvisioningState{
		"canceled":              ScriptProvisioningStateCanceled,
		"creating":              ScriptProvisioningStateCreating,
		"failed":                ScriptProvisioningStateFailed,
		"provisioningresources": ScriptProvisioningStateProvisioningResources,
		"running":               ScriptProvisioningStateRunning,
		"succeeded":             ScriptProvisioningStateSucceeded,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := ScriptProvisioningState(input)
	return &out, nil
}

type ScriptType string

const (
	ScriptTypeAzureCLI        ScriptType = "AzureCLI"
	ScriptTypeAzurePowerShell ScriptType = "AzurePowerShell"
)

func PossibleValuesForScriptType() []string {
	return []string{
		string(ScriptTypeAzureCLI),
		string(ScriptTypeAzurePowerShell),
	}
}

func (s *ScriptType) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseScriptType(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseScriptType(input string) (*ScriptType, error) {
	vals := map[string]ScriptType{
		"azurecli":        ScriptTypeAzureCLI,
		"azurepowershell": ScriptTypeAzurePowerShell,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := ScriptType(input)
	return &out, nil
}
