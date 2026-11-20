package agents

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type AgentAccessLevel string

const (
	AgentAccessLevelHigh AgentAccessLevel = "High"
	AgentAccessLevelLow  AgentAccessLevel = "Low"
)

func PossibleValuesForAgentAccessLevel() []string {
	return []string{
		string(AgentAccessLevelHigh),
		string(AgentAccessLevelLow),
	}
}

func (s *AgentAccessLevel) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseAgentAccessLevel(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseAgentAccessLevel(input string) (*AgentAccessLevel, error) {
	vals := map[string]AgentAccessLevel{
		"high": AgentAccessLevelHigh,
		"low":  AgentAccessLevelLow,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := AgentAccessLevel(input)
	return &out, nil
}

type AgentMode string

const (
	AgentModeAutonomous AgentMode = "Autonomous"
	AgentModeReadOnly   AgentMode = "ReadOnly"
	AgentModeReview     AgentMode = "Review"
)

func PossibleValuesForAgentMode() []string {
	return []string{
		string(AgentModeAutonomous),
		string(AgentModeReadOnly),
		string(AgentModeReview),
	}
}

func (s *AgentMode) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseAgentMode(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseAgentMode(input string) (*AgentMode, error) {
	vals := map[string]AgentMode{
		"autonomous": AgentModeAutonomous,
		"readonly":   AgentModeReadOnly,
		"review":     AgentModeReview,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := AgentMode(input)
	return &out, nil
}

type AgentPowerState string

const (
	AgentPowerStateRunning AgentPowerState = "Running"
	AgentPowerStateStopped AgentPowerState = "Stopped"
)

func PossibleValuesForAgentPowerState() []string {
	return []string{
		string(AgentPowerStateRunning),
		string(AgentPowerStateStopped),
	}
}

func (s *AgentPowerState) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseAgentPowerState(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseAgentPowerState(input string) (*AgentPowerState, error) {
	vals := map[string]AgentPowerState{
		"running": AgentPowerStateRunning,
		"stopped": AgentPowerStateStopped,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := AgentPowerState(input)
	return &out, nil
}

type AgentProvisioningState string

const (
	AgentProvisioningStateCanceled   AgentProvisioningState = "Canceled"
	AgentProvisioningStateDeleting   AgentProvisioningState = "Deleting"
	AgentProvisioningStateFailed     AgentProvisioningState = "Failed"
	AgentProvisioningStateInProgress AgentProvisioningState = "InProgress"
	AgentProvisioningStateSucceeded  AgentProvisioningState = "Succeeded"
)

func PossibleValuesForAgentProvisioningState() []string {
	return []string{
		string(AgentProvisioningStateCanceled),
		string(AgentProvisioningStateDeleting),
		string(AgentProvisioningStateFailed),
		string(AgentProvisioningStateInProgress),
		string(AgentProvisioningStateSucceeded),
	}
}

func (s *AgentProvisioningState) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseAgentProvisioningState(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseAgentProvisioningState(input string) (*AgentProvisioningState, error) {
	vals := map[string]AgentProvisioningState{
		"canceled":   AgentProvisioningStateCanceled,
		"deleting":   AgentProvisioningStateDeleting,
		"failed":     AgentProvisioningStateFailed,
		"inprogress": AgentProvisioningStateInProgress,
		"succeeded":  AgentProvisioningStateSucceeded,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := AgentProvisioningState(input)
	return &out, nil
}

type SandboxEgressMode string

const (
	SandboxEgressModeAzureVNet    SandboxEgressMode = "AzureVNet"
	SandboxEgressModeLimited      SandboxEgressMode = "Limited"
	SandboxEgressModeUnrestricted SandboxEgressMode = "Unrestricted"
)

func PossibleValuesForSandboxEgressMode() []string {
	return []string{
		string(SandboxEgressModeAzureVNet),
		string(SandboxEgressModeLimited),
		string(SandboxEgressModeUnrestricted),
	}
}

func (s *SandboxEgressMode) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseSandboxEgressMode(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseSandboxEgressMode(input string) (*SandboxEgressMode, error) {
	vals := map[string]SandboxEgressMode{
		"azurevnet":    SandboxEgressModeAzureVNet,
		"limited":      SandboxEgressModeLimited,
		"unrestricted": SandboxEgressModeUnrestricted,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := SandboxEgressMode(input)
	return &out, nil
}

type SandboxPackageManager string

const (
	SandboxPackageManagerNpm   SandboxPackageManager = "npm"
	SandboxPackageManagerNuget SandboxPackageManager = "nuget"
	SandboxPackageManagerPip   SandboxPackageManager = "pip"
)

func PossibleValuesForSandboxPackageManager() []string {
	return []string{
		string(SandboxPackageManagerNpm),
		string(SandboxPackageManagerNuget),
		string(SandboxPackageManagerPip),
	}
}

func (s *SandboxPackageManager) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseSandboxPackageManager(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseSandboxPackageManager(input string) (*SandboxPackageManager, error) {
	vals := map[string]SandboxPackageManager{
		"npm":   SandboxPackageManagerNpm,
		"nuget": SandboxPackageManagerNuget,
		"pip":   SandboxPackageManagerPip,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := SandboxPackageManager(input)
	return &out, nil
}

type UpgradeChannel string

const (
	UpgradeChannelPreview UpgradeChannel = "Preview"
	UpgradeChannelStable  UpgradeChannel = "Stable"
)

func PossibleValuesForUpgradeChannel() []string {
	return []string{
		string(UpgradeChannelPreview),
		string(UpgradeChannelStable),
	}
}

func (s *UpgradeChannel) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseUpgradeChannel(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseUpgradeChannel(input string) (*UpgradeChannel, error) {
	vals := map[string]UpgradeChannel{
		"preview": UpgradeChannelPreview,
		"stable":  UpgradeChannelStable,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := UpgradeChannel(input)
	return &out, nil
}
