package virtualmachineimagetemplate

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type AutoRunState string

const (
	AutoRunStateDisabled AutoRunState = "Disabled"
	AutoRunStateEnabled  AutoRunState = "Enabled"
)

func PossibleValuesForAutoRunState() []string {
	return []string{
		string(AutoRunStateDisabled),
		string(AutoRunStateEnabled),
	}
}

func (s *AutoRunState) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseAutoRunState(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseAutoRunState(input string) (*AutoRunState, error) {
	vals := map[string]AutoRunState{
		"disabled": AutoRunStateDisabled,
		"enabled":  AutoRunStateEnabled,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := AutoRunState(input)
	return &out, nil
}

type OnBuildError string

const (
	OnBuildErrorAbort   OnBuildError = "abort"
	OnBuildErrorCleanup OnBuildError = "cleanup"
)

func PossibleValuesForOnBuildError() []string {
	return []string{
		string(OnBuildErrorAbort),
		string(OnBuildErrorCleanup),
	}
}

func (s *OnBuildError) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseOnBuildError(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseOnBuildError(input string) (*OnBuildError, error) {
	vals := map[string]OnBuildError{
		"abort":   OnBuildErrorAbort,
		"cleanup": OnBuildErrorCleanup,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := OnBuildError(input)
	return &out, nil
}

type ProvisioningErrorCode string

const (
	ProvisioningErrorCodeBadCustomizerType           ProvisioningErrorCode = "BadCustomizerType"
	ProvisioningErrorCodeBadDistributeType           ProvisioningErrorCode = "BadDistributeType"
	ProvisioningErrorCodeBadManagedImageSource       ProvisioningErrorCode = "BadManagedImageSource"
	ProvisioningErrorCodeBadPIRSource                ProvisioningErrorCode = "BadPIRSource"
	ProvisioningErrorCodeBadSharedImageDistribute    ProvisioningErrorCode = "BadSharedImageDistribute"
	ProvisioningErrorCodeBadSharedImageVersionSource ProvisioningErrorCode = "BadSharedImageVersionSource"
	ProvisioningErrorCodeBadSourceType               ProvisioningErrorCode = "BadSourceType"
	ProvisioningErrorCodeBadStagingResourceGroup     ProvisioningErrorCode = "BadStagingResourceGroup"
	ProvisioningErrorCodeBadValidatorType            ProvisioningErrorCode = "BadValidatorType"
	ProvisioningErrorCodeNoCustomizerScript          ProvisioningErrorCode = "NoCustomizerScript"
	ProvisioningErrorCodeNoValidatorScript           ProvisioningErrorCode = "NoValidatorScript"
	ProvisioningErrorCodeOther                       ProvisioningErrorCode = "Other"
	ProvisioningErrorCodeServerError                 ProvisioningErrorCode = "ServerError"
	ProvisioningErrorCodeUnsupportedCustomizerType   ProvisioningErrorCode = "UnsupportedCustomizerType"
	ProvisioningErrorCodeUnsupportedValidatorType    ProvisioningErrorCode = "UnsupportedValidatorType"
)

func PossibleValuesForProvisioningErrorCode() []string {
	return []string{
		string(ProvisioningErrorCodeBadCustomizerType),
		string(ProvisioningErrorCodeBadDistributeType),
		string(ProvisioningErrorCodeBadManagedImageSource),
		string(ProvisioningErrorCodeBadPIRSource),
		string(ProvisioningErrorCodeBadSharedImageDistribute),
		string(ProvisioningErrorCodeBadSharedImageVersionSource),
		string(ProvisioningErrorCodeBadSourceType),
		string(ProvisioningErrorCodeBadStagingResourceGroup),
		string(ProvisioningErrorCodeBadValidatorType),
		string(ProvisioningErrorCodeNoCustomizerScript),
		string(ProvisioningErrorCodeNoValidatorScript),
		string(ProvisioningErrorCodeOther),
		string(ProvisioningErrorCodeServerError),
		string(ProvisioningErrorCodeUnsupportedCustomizerType),
		string(ProvisioningErrorCodeUnsupportedValidatorType),
	}
}

func (s *ProvisioningErrorCode) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseProvisioningErrorCode(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseProvisioningErrorCode(input string) (*ProvisioningErrorCode, error) {
	vals := map[string]ProvisioningErrorCode{
		"badcustomizertype":           ProvisioningErrorCodeBadCustomizerType,
		"baddistributetype":           ProvisioningErrorCodeBadDistributeType,
		"badmanagedimagesource":       ProvisioningErrorCodeBadManagedImageSource,
		"badpirsource":                ProvisioningErrorCodeBadPIRSource,
		"badsharedimagedistribute":    ProvisioningErrorCodeBadSharedImageDistribute,
		"badsharedimageversionsource": ProvisioningErrorCodeBadSharedImageVersionSource,
		"badsourcetype":               ProvisioningErrorCodeBadSourceType,
		"badstagingresourcegroup":     ProvisioningErrorCodeBadStagingResourceGroup,
		"badvalidatortype":            ProvisioningErrorCodeBadValidatorType,
		"nocustomizerscript":          ProvisioningErrorCodeNoCustomizerScript,
		"novalidatorscript":           ProvisioningErrorCodeNoValidatorScript,
		"other":                       ProvisioningErrorCodeOther,
		"servererror":                 ProvisioningErrorCodeServerError,
		"unsupportedcustomizertype":   ProvisioningErrorCodeUnsupportedCustomizerType,
		"unsupportedvalidatortype":    ProvisioningErrorCodeUnsupportedValidatorType,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := ProvisioningErrorCode(input)
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

type RunState string

const (
	RunStateCanceled           RunState = "Canceled"
	RunStateCanceling          RunState = "Canceling"
	RunStateFailed             RunState = "Failed"
	RunStatePartiallySucceeded RunState = "PartiallySucceeded"
	RunStateRunning            RunState = "Running"
	RunStateSucceeded          RunState = "Succeeded"
)

func PossibleValuesForRunState() []string {
	return []string{
		string(RunStateCanceled),
		string(RunStateCanceling),
		string(RunStateFailed),
		string(RunStatePartiallySucceeded),
		string(RunStateRunning),
		string(RunStateSucceeded),
	}
}

func (s *RunState) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseRunState(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseRunState(input string) (*RunState, error) {
	vals := map[string]RunState{
		"canceled":           RunStateCanceled,
		"canceling":          RunStateCanceling,
		"failed":             RunStateFailed,
		"partiallysucceeded": RunStatePartiallySucceeded,
		"running":            RunStateRunning,
		"succeeded":          RunStateSucceeded,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := RunState(input)
	return &out, nil
}

type RunSubState string

const (
	RunSubStateBuilding     RunSubState = "Building"
	RunSubStateCustomizing  RunSubState = "Customizing"
	RunSubStateDistributing RunSubState = "Distributing"
	RunSubStateOptimizing   RunSubState = "Optimizing"
	RunSubStateQueued       RunSubState = "Queued"
	RunSubStateValidating   RunSubState = "Validating"
)

func PossibleValuesForRunSubState() []string {
	return []string{
		string(RunSubStateBuilding),
		string(RunSubStateCustomizing),
		string(RunSubStateDistributing),
		string(RunSubStateOptimizing),
		string(RunSubStateQueued),
		string(RunSubStateValidating),
	}
}

func (s *RunSubState) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseRunSubState(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseRunSubState(input string) (*RunSubState, error) {
	vals := map[string]RunSubState{
		"building":     RunSubStateBuilding,
		"customizing":  RunSubStateCustomizing,
		"distributing": RunSubStateDistributing,
		"optimizing":   RunSubStateOptimizing,
		"queued":       RunSubStateQueued,
		"validating":   RunSubStateValidating,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := RunSubState(input)
	return &out, nil
}

type SharedImageStorageAccountType string

const (
	SharedImageStorageAccountTypePremiumLRS  SharedImageStorageAccountType = "Premium_LRS"
	SharedImageStorageAccountTypeStandardLRS SharedImageStorageAccountType = "Standard_LRS"
	SharedImageStorageAccountTypeStandardZRS SharedImageStorageAccountType = "Standard_ZRS"
)

func PossibleValuesForSharedImageStorageAccountType() []string {
	return []string{
		string(SharedImageStorageAccountTypePremiumLRS),
		string(SharedImageStorageAccountTypeStandardLRS),
		string(SharedImageStorageAccountTypeStandardZRS),
	}
}

func (s *SharedImageStorageAccountType) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseSharedImageStorageAccountType(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseSharedImageStorageAccountType(input string) (*SharedImageStorageAccountType, error) {
	vals := map[string]SharedImageStorageAccountType{
		"premium_lrs":  SharedImageStorageAccountTypePremiumLRS,
		"standard_lrs": SharedImageStorageAccountTypeStandardLRS,
		"standard_zrs": SharedImageStorageAccountTypeStandardZRS,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := SharedImageStorageAccountType(input)
	return &out, nil
}

type VMBootOptimizationState string

const (
	VMBootOptimizationStateDisabled VMBootOptimizationState = "Disabled"
	VMBootOptimizationStateEnabled  VMBootOptimizationState = "Enabled"
)

func PossibleValuesForVMBootOptimizationState() []string {
	return []string{
		string(VMBootOptimizationStateDisabled),
		string(VMBootOptimizationStateEnabled),
	}
}

func (s *VMBootOptimizationState) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseVMBootOptimizationState(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseVMBootOptimizationState(input string) (*VMBootOptimizationState, error) {
	vals := map[string]VMBootOptimizationState{
		"disabled": VMBootOptimizationStateDisabled,
		"enabled":  VMBootOptimizationStateEnabled,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := VMBootOptimizationState(input)
	return &out, nil
}
