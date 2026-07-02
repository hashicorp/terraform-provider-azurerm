package storagetasks

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type OnFailure string

const (
	OnFailureBreak OnFailure = "break"
)

func PossibleValuesForOnFailure() []string {
	return []string{
		string(OnFailureBreak),
	}
}

func (s *OnFailure) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseOnFailure(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseOnFailure(input string) (*OnFailure, error) {
	vals := map[string]OnFailure{
		"break": OnFailureBreak,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := OnFailure(input)
	return &out, nil
}

type OnSuccess string

const (
	OnSuccessContinue OnSuccess = "continue"
)

func PossibleValuesForOnSuccess() []string {
	return []string{
		string(OnSuccessContinue),
	}
}

func (s *OnSuccess) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseOnSuccess(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseOnSuccess(input string) (*OnSuccess, error) {
	vals := map[string]OnSuccess{
		"continue": OnSuccessContinue,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := OnSuccess(input)
	return &out, nil
}

type ProvisioningState string

const (
	ProvisioningStateAccepted                       ProvisioningState = "Accepted"
	ProvisioningStateCanceled                       ProvisioningState = "Canceled"
	ProvisioningStateCreating                       ProvisioningState = "Creating"
	ProvisioningStateDeleting                       ProvisioningState = "Deleting"
	ProvisioningStateFailed                         ProvisioningState = "Failed"
	ProvisioningStateSucceeded                      ProvisioningState = "Succeeded"
	ProvisioningStateValidateSubscriptionQuotaBegin ProvisioningState = "ValidateSubscriptionQuotaBegin"
	ProvisioningStateValidateSubscriptionQuotaEnd   ProvisioningState = "ValidateSubscriptionQuotaEnd"
)

func PossibleValuesForProvisioningState() []string {
	return []string{
		string(ProvisioningStateAccepted),
		string(ProvisioningStateCanceled),
		string(ProvisioningStateCreating),
		string(ProvisioningStateDeleting),
		string(ProvisioningStateFailed),
		string(ProvisioningStateSucceeded),
		string(ProvisioningStateValidateSubscriptionQuotaBegin),
		string(ProvisioningStateValidateSubscriptionQuotaEnd),
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
		"accepted":                       ProvisioningStateAccepted,
		"canceled":                       ProvisioningStateCanceled,
		"creating":                       ProvisioningStateCreating,
		"deleting":                       ProvisioningStateDeleting,
		"failed":                         ProvisioningStateFailed,
		"succeeded":                      ProvisioningStateSucceeded,
		"validatesubscriptionquotabegin": ProvisioningStateValidateSubscriptionQuotaBegin,
		"validatesubscriptionquotaend":   ProvisioningStateValidateSubscriptionQuotaEnd,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := ProvisioningState(input)
	return &out, nil
}

type RunResult string

const (
	RunResultFailed    RunResult = "Failed"
	RunResultSucceeded RunResult = "Succeeded"
)

func PossibleValuesForRunResult() []string {
	return []string{
		string(RunResultFailed),
		string(RunResultSucceeded),
	}
}

func (s *RunResult) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseRunResult(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseRunResult(input string) (*RunResult, error) {
	vals := map[string]RunResult{
		"failed":    RunResultFailed,
		"succeeded": RunResultSucceeded,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := RunResult(input)
	return &out, nil
}

type RunStatusEnum string

const (
	RunStatusEnumFinished   RunStatusEnum = "Finished"
	RunStatusEnumInProgress RunStatusEnum = "InProgress"
)

func PossibleValuesForRunStatusEnum() []string {
	return []string{
		string(RunStatusEnumFinished),
		string(RunStatusEnumInProgress),
	}
}

func (s *RunStatusEnum) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseRunStatusEnum(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseRunStatusEnum(input string) (*RunStatusEnum, error) {
	vals := map[string]RunStatusEnum{
		"finished":   RunStatusEnumFinished,
		"inprogress": RunStatusEnumInProgress,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := RunStatusEnum(input)
	return &out, nil
}

type StorageTaskOperationName string

const (
	StorageTaskOperationNameDeleteBlob                StorageTaskOperationName = "DeleteBlob"
	StorageTaskOperationNameSetBlobExpiry             StorageTaskOperationName = "SetBlobExpiry"
	StorageTaskOperationNameSetBlobImmutabilityPolicy StorageTaskOperationName = "SetBlobImmutabilityPolicy"
	StorageTaskOperationNameSetBlobLegalHold          StorageTaskOperationName = "SetBlobLegalHold"
	StorageTaskOperationNameSetBlobTags               StorageTaskOperationName = "SetBlobTags"
	StorageTaskOperationNameSetBlobTier               StorageTaskOperationName = "SetBlobTier"
	StorageTaskOperationNameUndeleteBlob              StorageTaskOperationName = "UndeleteBlob"
)

func PossibleValuesForStorageTaskOperationName() []string {
	return []string{
		string(StorageTaskOperationNameDeleteBlob),
		string(StorageTaskOperationNameSetBlobExpiry),
		string(StorageTaskOperationNameSetBlobImmutabilityPolicy),
		string(StorageTaskOperationNameSetBlobLegalHold),
		string(StorageTaskOperationNameSetBlobTags),
		string(StorageTaskOperationNameSetBlobTier),
		string(StorageTaskOperationNameUndeleteBlob),
	}
}

func (s *StorageTaskOperationName) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseStorageTaskOperationName(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseStorageTaskOperationName(input string) (*StorageTaskOperationName, error) {
	vals := map[string]StorageTaskOperationName{
		"deleteblob":                StorageTaskOperationNameDeleteBlob,
		"setblobexpiry":             StorageTaskOperationNameSetBlobExpiry,
		"setblobimmutabilitypolicy": StorageTaskOperationNameSetBlobImmutabilityPolicy,
		"setbloblegalhold":          StorageTaskOperationNameSetBlobLegalHold,
		"setblobtags":               StorageTaskOperationNameSetBlobTags,
		"setblobtier":               StorageTaskOperationNameSetBlobTier,
		"undeleteblob":              StorageTaskOperationNameUndeleteBlob,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := StorageTaskOperationName(input)
	return &out, nil
}
