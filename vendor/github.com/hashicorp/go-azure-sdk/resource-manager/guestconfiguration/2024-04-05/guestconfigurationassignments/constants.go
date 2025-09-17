package guestconfigurationassignments

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ActionAfterReboot string

const (
	ActionAfterRebootContinueConfiguration ActionAfterReboot = "ContinueConfiguration"
	ActionAfterRebootStopConfiguration     ActionAfterReboot = "StopConfiguration"
)

func PossibleValuesForActionAfterReboot() []string {
	return []string{
		string(ActionAfterRebootContinueConfiguration),
		string(ActionAfterRebootStopConfiguration),
	}
}

func (s *ActionAfterReboot) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseActionAfterReboot(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseActionAfterReboot(input string) (*ActionAfterReboot, error) {
	vals := map[string]ActionAfterReboot{
		"continueconfiguration": ActionAfterRebootContinueConfiguration,
		"stopconfiguration":     ActionAfterRebootStopConfiguration,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := ActionAfterReboot(input)
	return &out, nil
}

type AssignmentType string

const (
	AssignmentTypeApplyAndAutoCorrect  AssignmentType = "ApplyAndAutoCorrect"
	AssignmentTypeApplyAndMonitor      AssignmentType = "ApplyAndMonitor"
	AssignmentTypeAudit                AssignmentType = "Audit"
	AssignmentTypeDeployAndAutoCorrect AssignmentType = "DeployAndAutoCorrect"
)

func PossibleValuesForAssignmentType() []string {
	return []string{
		string(AssignmentTypeApplyAndAutoCorrect),
		string(AssignmentTypeApplyAndMonitor),
		string(AssignmentTypeAudit),
		string(AssignmentTypeDeployAndAutoCorrect),
	}
}

func (s *AssignmentType) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseAssignmentType(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseAssignmentType(input string) (*AssignmentType, error) {
	vals := map[string]AssignmentType{
		"applyandautocorrect":  AssignmentTypeApplyAndAutoCorrect,
		"applyandmonitor":      AssignmentTypeApplyAndMonitor,
		"audit":                AssignmentTypeAudit,
		"deployandautocorrect": AssignmentTypeDeployAndAutoCorrect,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := AssignmentType(input)
	return &out, nil
}

type ComplianceStatus string

const (
	ComplianceStatusCompliant    ComplianceStatus = "Compliant"
	ComplianceStatusNonCompliant ComplianceStatus = "NonCompliant"
	ComplianceStatusPending      ComplianceStatus = "Pending"
)

func PossibleValuesForComplianceStatus() []string {
	return []string{
		string(ComplianceStatusCompliant),
		string(ComplianceStatusNonCompliant),
		string(ComplianceStatusPending),
	}
}

func (s *ComplianceStatus) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseComplianceStatus(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseComplianceStatus(input string) (*ComplianceStatus, error) {
	vals := map[string]ComplianceStatus{
		"compliant":    ComplianceStatusCompliant,
		"noncompliant": ComplianceStatusNonCompliant,
		"pending":      ComplianceStatusPending,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := ComplianceStatus(input)
	return &out, nil
}

type ConfigurationMode string

const (
	ConfigurationModeApplyAndAutoCorrect ConfigurationMode = "ApplyAndAutoCorrect"
	ConfigurationModeApplyAndMonitor     ConfigurationMode = "ApplyAndMonitor"
	ConfigurationModeApplyOnly           ConfigurationMode = "ApplyOnly"
)

func PossibleValuesForConfigurationMode() []string {
	return []string{
		string(ConfigurationModeApplyAndAutoCorrect),
		string(ConfigurationModeApplyAndMonitor),
		string(ConfigurationModeApplyOnly),
	}
}

func (s *ConfigurationMode) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseConfigurationMode(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseConfigurationMode(input string) (*ConfigurationMode, error) {
	vals := map[string]ConfigurationMode{
		"applyandautocorrect": ConfigurationModeApplyAndAutoCorrect,
		"applyandmonitor":     ConfigurationModeApplyAndMonitor,
		"applyonly":           ConfigurationModeApplyOnly,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := ConfigurationMode(input)
	return &out, nil
}

type Kind string

const (
	KindDSC Kind = "DSC"
)

func PossibleValuesForKind() []string {
	return []string{
		string(KindDSC),
	}
}

func (s *Kind) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseKind(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseKind(input string) (*Kind, error) {
	vals := map[string]Kind{
		"dsc": KindDSC,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := Kind(input)
	return &out, nil
}

type ProvisioningState string

const (
	ProvisioningStateCanceled  ProvisioningState = "Canceled"
	ProvisioningStateCreated   ProvisioningState = "Created"
	ProvisioningStateFailed    ProvisioningState = "Failed"
	ProvisioningStateSucceeded ProvisioningState = "Succeeded"
)

func PossibleValuesForProvisioningState() []string {
	return []string{
		string(ProvisioningStateCanceled),
		string(ProvisioningStateCreated),
		string(ProvisioningStateFailed),
		string(ProvisioningStateSucceeded),
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
		"created":   ProvisioningStateCreated,
		"failed":    ProvisioningStateFailed,
		"succeeded": ProvisioningStateSucceeded,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := ProvisioningState(input)
	return &out, nil
}

type Type string

const (
	TypeConsistency Type = "Consistency"
	TypeInitial     Type = "Initial"
)

func PossibleValuesForType() []string {
	return []string{
		string(TypeConsistency),
		string(TypeInitial),
	}
}

func (s *Type) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseType(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseType(input string) (*Type, error) {
	vals := map[string]Type{
		"consistency": TypeConsistency,
		"initial":     TypeInitial,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := Type(input)
	return &out, nil
}
