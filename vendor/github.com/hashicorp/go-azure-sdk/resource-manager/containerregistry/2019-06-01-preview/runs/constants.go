package runs

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type Architecture string

const (
	ArchitectureAmdSixFour    Architecture = "amd64"
	ArchitectureArm           Architecture = "arm"
	ArchitectureArmSixFour    Architecture = "arm64"
	ArchitectureThreeEightSix Architecture = "386"
	ArchitectureXEightSix     Architecture = "x86"
)

func PossibleValuesForArchitecture() []string {
	return []string{
		string(ArchitectureAmdSixFour),
		string(ArchitectureArm),
		string(ArchitectureArmSixFour),
		string(ArchitectureThreeEightSix),
		string(ArchitectureXEightSix),
	}
}

func (s *Architecture) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseArchitecture(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseArchitecture(input string) (*Architecture, error) {
	vals := map[string]Architecture{
		"amd64": ArchitectureAmdSixFour,
		"arm":   ArchitectureArm,
		"arm64": ArchitectureArmSixFour,
		"386":   ArchitectureThreeEightSix,
		"x86":   ArchitectureXEightSix,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := Architecture(input)
	return &out, nil
}

type OS string

const (
	OSLinux   OS = "Linux"
	OSWindows OS = "Windows"
)

func PossibleValuesForOS() []string {
	return []string{
		string(OSLinux),
		string(OSWindows),
	}
}

func (s *OS) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseOS(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseOS(input string) (*OS, error) {
	vals := map[string]OS{
		"linux":   OSLinux,
		"windows": OSWindows,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := OS(input)
	return &out, nil
}

type ProvisioningState string

const (
	ProvisioningStateCanceled  ProvisioningState = "Canceled"
	ProvisioningStateCreating  ProvisioningState = "Creating"
	ProvisioningStateDeleting  ProvisioningState = "Deleting"
	ProvisioningStateFailed    ProvisioningState = "Failed"
	ProvisioningStateSucceeded ProvisioningState = "Succeeded"
	ProvisioningStateUpdating  ProvisioningState = "Updating"
)

func PossibleValuesForProvisioningState() []string {
	return []string{
		string(ProvisioningStateCanceled),
		string(ProvisioningStateCreating),
		string(ProvisioningStateDeleting),
		string(ProvisioningStateFailed),
		string(ProvisioningStateSucceeded),
		string(ProvisioningStateUpdating),
	}
}

func (s *ProvisioningState) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseProvisioningState(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseProvisioningState(input string) (*ProvisioningState, error) {
	vals := map[string]ProvisioningState{
		"canceled":  ProvisioningStateCanceled,
		"creating":  ProvisioningStateCreating,
		"deleting":  ProvisioningStateDeleting,
		"failed":    ProvisioningStateFailed,
		"succeeded": ProvisioningStateSucceeded,
		"updating":  ProvisioningStateUpdating,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := ProvisioningState(input)
	return &out, nil
}

type RunStatus string

const (
	RunStatusCanceled  RunStatus = "Canceled"
	RunStatusError     RunStatus = "Error"
	RunStatusFailed    RunStatus = "Failed"
	RunStatusQueued    RunStatus = "Queued"
	RunStatusRunning   RunStatus = "Running"
	RunStatusStarted   RunStatus = "Started"
	RunStatusSucceeded RunStatus = "Succeeded"
	RunStatusTimeout   RunStatus = "Timeout"
)

func PossibleValuesForRunStatus() []string {
	return []string{
		string(RunStatusCanceled),
		string(RunStatusError),
		string(RunStatusFailed),
		string(RunStatusQueued),
		string(RunStatusRunning),
		string(RunStatusStarted),
		string(RunStatusSucceeded),
		string(RunStatusTimeout),
	}
}

func (s *RunStatus) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseRunStatus(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseRunStatus(input string) (*RunStatus, error) {
	vals := map[string]RunStatus{
		"canceled":  RunStatusCanceled,
		"error":     RunStatusError,
		"failed":    RunStatusFailed,
		"queued":    RunStatusQueued,
		"running":   RunStatusRunning,
		"started":   RunStatusStarted,
		"succeeded": RunStatusSucceeded,
		"timeout":   RunStatusTimeout,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := RunStatus(input)
	return &out, nil
}

type RunType string

const (
	RunTypeAutoBuild  RunType = "AutoBuild"
	RunTypeAutoRun    RunType = "AutoRun"
	RunTypeQuickBuild RunType = "QuickBuild"
	RunTypeQuickRun   RunType = "QuickRun"
)

func PossibleValuesForRunType() []string {
	return []string{
		string(RunTypeAutoBuild),
		string(RunTypeAutoRun),
		string(RunTypeQuickBuild),
		string(RunTypeQuickRun),
	}
}

func (s *RunType) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseRunType(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseRunType(input string) (*RunType, error) {
	vals := map[string]RunType{
		"autobuild":  RunTypeAutoBuild,
		"autorun":    RunTypeAutoRun,
		"quickbuild": RunTypeQuickBuild,
		"quickrun":   RunTypeQuickRun,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := RunType(input)
	return &out, nil
}

type Variant string

const (
	VariantVEight Variant = "v8"
	VariantVSeven Variant = "v7"
	VariantVSix   Variant = "v6"
)

func PossibleValuesForVariant() []string {
	return []string{
		string(VariantVEight),
		string(VariantVSeven),
		string(VariantVSix),
	}
}

func (s *Variant) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseVariant(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseVariant(input string) (*Variant, error) {
	vals := map[string]Variant{
		"v8": VariantVEight,
		"v7": VariantVSeven,
		"v6": VariantVSix,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := Variant(input)
	return &out, nil
}
