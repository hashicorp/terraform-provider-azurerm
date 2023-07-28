package batchdeployment

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type BatchLoggingLevel string

const (
	BatchLoggingLevelDebug   BatchLoggingLevel = "Debug"
	BatchLoggingLevelInfo    BatchLoggingLevel = "Info"
	BatchLoggingLevelWarning BatchLoggingLevel = "Warning"
)

func PossibleValuesForBatchLoggingLevel() []string {
	return []string{
		string(BatchLoggingLevelDebug),
		string(BatchLoggingLevelInfo),
		string(BatchLoggingLevelWarning),
	}
}

func (s *BatchLoggingLevel) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseBatchLoggingLevel(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseBatchLoggingLevel(input string) (*BatchLoggingLevel, error) {
	vals := map[string]BatchLoggingLevel{
		"debug":   BatchLoggingLevelDebug,
		"info":    BatchLoggingLevelInfo,
		"warning": BatchLoggingLevelWarning,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := BatchLoggingLevel(input)
	return &out, nil
}

type BatchOutputAction string

const (
	BatchOutputActionAppendRow   BatchOutputAction = "AppendRow"
	BatchOutputActionSummaryOnly BatchOutputAction = "SummaryOnly"
)

func PossibleValuesForBatchOutputAction() []string {
	return []string{
		string(BatchOutputActionAppendRow),
		string(BatchOutputActionSummaryOnly),
	}
}

func (s *BatchOutputAction) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseBatchOutputAction(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseBatchOutputAction(input string) (*BatchOutputAction, error) {
	vals := map[string]BatchOutputAction{
		"appendrow":   BatchOutputActionAppendRow,
		"summaryonly": BatchOutputActionSummaryOnly,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := BatchOutputAction(input)
	return &out, nil
}

type DeploymentProvisioningState string

const (
	DeploymentProvisioningStateCanceled  DeploymentProvisioningState = "Canceled"
	DeploymentProvisioningStateCreating  DeploymentProvisioningState = "Creating"
	DeploymentProvisioningStateDeleting  DeploymentProvisioningState = "Deleting"
	DeploymentProvisioningStateFailed    DeploymentProvisioningState = "Failed"
	DeploymentProvisioningStateScaling   DeploymentProvisioningState = "Scaling"
	DeploymentProvisioningStateSucceeded DeploymentProvisioningState = "Succeeded"
	DeploymentProvisioningStateUpdating  DeploymentProvisioningState = "Updating"
)

func PossibleValuesForDeploymentProvisioningState() []string {
	return []string{
		string(DeploymentProvisioningStateCanceled),
		string(DeploymentProvisioningStateCreating),
		string(DeploymentProvisioningStateDeleting),
		string(DeploymentProvisioningStateFailed),
		string(DeploymentProvisioningStateScaling),
		string(DeploymentProvisioningStateSucceeded),
		string(DeploymentProvisioningStateUpdating),
	}
}

func (s *DeploymentProvisioningState) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseDeploymentProvisioningState(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseDeploymentProvisioningState(input string) (*DeploymentProvisioningState, error) {
	vals := map[string]DeploymentProvisioningState{
		"canceled":  DeploymentProvisioningStateCanceled,
		"creating":  DeploymentProvisioningStateCreating,
		"deleting":  DeploymentProvisioningStateDeleting,
		"failed":    DeploymentProvisioningStateFailed,
		"scaling":   DeploymentProvisioningStateScaling,
		"succeeded": DeploymentProvisioningStateSucceeded,
		"updating":  DeploymentProvisioningStateUpdating,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := DeploymentProvisioningState(input)
	return &out, nil
}

type ReferenceType string

const (
	ReferenceTypeDataPath   ReferenceType = "DataPath"
	ReferenceTypeId         ReferenceType = "Id"
	ReferenceTypeOutputPath ReferenceType = "OutputPath"
)

func PossibleValuesForReferenceType() []string {
	return []string{
		string(ReferenceTypeDataPath),
		string(ReferenceTypeId),
		string(ReferenceTypeOutputPath),
	}
}

func (s *ReferenceType) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseReferenceType(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseReferenceType(input string) (*ReferenceType, error) {
	vals := map[string]ReferenceType{
		"datapath":   ReferenceTypeDataPath,
		"id":         ReferenceTypeId,
		"outputpath": ReferenceTypeOutputPath,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := ReferenceType(input)
	return &out, nil
}

type SkuTier string

const (
	SkuTierBasic    SkuTier = "Basic"
	SkuTierFree     SkuTier = "Free"
	SkuTierPremium  SkuTier = "Premium"
	SkuTierStandard SkuTier = "Standard"
)

func PossibleValuesForSkuTier() []string {
	return []string{
		string(SkuTierBasic),
		string(SkuTierFree),
		string(SkuTierPremium),
		string(SkuTierStandard),
	}
}

func (s *SkuTier) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseSkuTier(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseSkuTier(input string) (*SkuTier, error) {
	vals := map[string]SkuTier{
		"basic":    SkuTierBasic,
		"free":     SkuTierFree,
		"premium":  SkuTierPremium,
		"standard": SkuTierStandard,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := SkuTier(input)
	return &out, nil
}
