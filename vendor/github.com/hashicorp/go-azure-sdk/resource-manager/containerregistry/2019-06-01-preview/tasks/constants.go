package tasks

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

type BaseImageDependencyType string

const (
	BaseImageDependencyTypeBuildTime BaseImageDependencyType = "BuildTime"
	BaseImageDependencyTypeRunTime   BaseImageDependencyType = "RunTime"
)

func PossibleValuesForBaseImageDependencyType() []string {
	return []string{
		string(BaseImageDependencyTypeBuildTime),
		string(BaseImageDependencyTypeRunTime),
	}
}

func (s *BaseImageDependencyType) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseBaseImageDependencyType(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseBaseImageDependencyType(input string) (*BaseImageDependencyType, error) {
	vals := map[string]BaseImageDependencyType{
		"buildtime": BaseImageDependencyTypeBuildTime,
		"runtime":   BaseImageDependencyTypeRunTime,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := BaseImageDependencyType(input)
	return &out, nil
}

type BaseImageTriggerType string

const (
	BaseImageTriggerTypeAll     BaseImageTriggerType = "All"
	BaseImageTriggerTypeRuntime BaseImageTriggerType = "Runtime"
)

func PossibleValuesForBaseImageTriggerType() []string {
	return []string{
		string(BaseImageTriggerTypeAll),
		string(BaseImageTriggerTypeRuntime),
	}
}

func (s *BaseImageTriggerType) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseBaseImageTriggerType(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseBaseImageTriggerType(input string) (*BaseImageTriggerType, error) {
	vals := map[string]BaseImageTriggerType{
		"all":     BaseImageTriggerTypeAll,
		"runtime": BaseImageTriggerTypeRuntime,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := BaseImageTriggerType(input)
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

type SecretObjectType string

const (
	SecretObjectTypeOpaque      SecretObjectType = "Opaque"
	SecretObjectTypeVaultsecret SecretObjectType = "Vaultsecret"
)

func PossibleValuesForSecretObjectType() []string {
	return []string{
		string(SecretObjectTypeOpaque),
		string(SecretObjectTypeVaultsecret),
	}
}

func (s *SecretObjectType) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseSecretObjectType(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseSecretObjectType(input string) (*SecretObjectType, error) {
	vals := map[string]SecretObjectType{
		"opaque":      SecretObjectTypeOpaque,
		"vaultsecret": SecretObjectTypeVaultsecret,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := SecretObjectType(input)
	return &out, nil
}

type SourceControlType string

const (
	SourceControlTypeGithub                  SourceControlType = "Github"
	SourceControlTypeVisualStudioTeamService SourceControlType = "VisualStudioTeamService"
)

func PossibleValuesForSourceControlType() []string {
	return []string{
		string(SourceControlTypeGithub),
		string(SourceControlTypeVisualStudioTeamService),
	}
}

func (s *SourceControlType) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseSourceControlType(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseSourceControlType(input string) (*SourceControlType, error) {
	vals := map[string]SourceControlType{
		"github":                  SourceControlTypeGithub,
		"visualstudioteamservice": SourceControlTypeVisualStudioTeamService,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := SourceControlType(input)
	return &out, nil
}

type SourceRegistryLoginMode string

const (
	SourceRegistryLoginModeDefault SourceRegistryLoginMode = "Default"
	SourceRegistryLoginModeNone    SourceRegistryLoginMode = "None"
)

func PossibleValuesForSourceRegistryLoginMode() []string {
	return []string{
		string(SourceRegistryLoginModeDefault),
		string(SourceRegistryLoginModeNone),
	}
}

func (s *SourceRegistryLoginMode) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseSourceRegistryLoginMode(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseSourceRegistryLoginMode(input string) (*SourceRegistryLoginMode, error) {
	vals := map[string]SourceRegistryLoginMode{
		"default": SourceRegistryLoginModeDefault,
		"none":    SourceRegistryLoginModeNone,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := SourceRegistryLoginMode(input)
	return &out, nil
}

type SourceTriggerEvent string

const (
	SourceTriggerEventCommit      SourceTriggerEvent = "commit"
	SourceTriggerEventPullrequest SourceTriggerEvent = "pullrequest"
)

func PossibleValuesForSourceTriggerEvent() []string {
	return []string{
		string(SourceTriggerEventCommit),
		string(SourceTriggerEventPullrequest),
	}
}

func (s *SourceTriggerEvent) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseSourceTriggerEvent(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseSourceTriggerEvent(input string) (*SourceTriggerEvent, error) {
	vals := map[string]SourceTriggerEvent{
		"commit":      SourceTriggerEventCommit,
		"pullrequest": SourceTriggerEventPullrequest,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := SourceTriggerEvent(input)
	return &out, nil
}

type StepType string

const (
	StepTypeDocker      StepType = "Docker"
	StepTypeEncodedTask StepType = "EncodedTask"
	StepTypeFileTask    StepType = "FileTask"
)

func PossibleValuesForStepType() []string {
	return []string{
		string(StepTypeDocker),
		string(StepTypeEncodedTask),
		string(StepTypeFileTask),
	}
}

func (s *StepType) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseStepType(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseStepType(input string) (*StepType, error) {
	vals := map[string]StepType{
		"docker":      StepTypeDocker,
		"encodedtask": StepTypeEncodedTask,
		"filetask":    StepTypeFileTask,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := StepType(input)
	return &out, nil
}

type TaskStatus string

const (
	TaskStatusDisabled TaskStatus = "Disabled"
	TaskStatusEnabled  TaskStatus = "Enabled"
)

func PossibleValuesForTaskStatus() []string {
	return []string{
		string(TaskStatusDisabled),
		string(TaskStatusEnabled),
	}
}

func (s *TaskStatus) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseTaskStatus(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseTaskStatus(input string) (*TaskStatus, error) {
	vals := map[string]TaskStatus{
		"disabled": TaskStatusDisabled,
		"enabled":  TaskStatusEnabled,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := TaskStatus(input)
	return &out, nil
}

type TokenType string

const (
	TokenTypeOAuth TokenType = "OAuth"
	TokenTypePAT   TokenType = "PAT"
)

func PossibleValuesForTokenType() []string {
	return []string{
		string(TokenTypeOAuth),
		string(TokenTypePAT),
	}
}

func (s *TokenType) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseTokenType(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseTokenType(input string) (*TokenType, error) {
	vals := map[string]TokenType{
		"oauth": TokenTypeOAuth,
		"pat":   TokenTypePAT,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := TokenType(input)
	return &out, nil
}

type TriggerStatus string

const (
	TriggerStatusDisabled TriggerStatus = "Disabled"
	TriggerStatusEnabled  TriggerStatus = "Enabled"
)

func PossibleValuesForTriggerStatus() []string {
	return []string{
		string(TriggerStatusDisabled),
		string(TriggerStatusEnabled),
	}
}

func (s *TriggerStatus) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseTriggerStatus(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseTriggerStatus(input string) (*TriggerStatus, error) {
	vals := map[string]TriggerStatus{
		"disabled": TriggerStatusDisabled,
		"enabled":  TriggerStatusEnabled,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := TriggerStatus(input)
	return &out, nil
}

type UpdateTriggerPayloadType string

const (
	UpdateTriggerPayloadTypeDefault UpdateTriggerPayloadType = "Default"
	UpdateTriggerPayloadTypeToken   UpdateTriggerPayloadType = "Token"
)

func PossibleValuesForUpdateTriggerPayloadType() []string {
	return []string{
		string(UpdateTriggerPayloadTypeDefault),
		string(UpdateTriggerPayloadTypeToken),
	}
}

func (s *UpdateTriggerPayloadType) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseUpdateTriggerPayloadType(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseUpdateTriggerPayloadType(input string) (*UpdateTriggerPayloadType, error) {
	vals := map[string]UpdateTriggerPayloadType{
		"default": UpdateTriggerPayloadTypeDefault,
		"token":   UpdateTriggerPayloadTypeToken,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := UpdateTriggerPayloadType(input)
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
