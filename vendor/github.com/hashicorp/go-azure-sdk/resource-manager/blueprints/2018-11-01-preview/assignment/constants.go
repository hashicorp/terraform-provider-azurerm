package assignment

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type AssignmentDeleteBehavior string

const (
	AssignmentDeleteBehaviorAll  AssignmentDeleteBehavior = "all"
	AssignmentDeleteBehaviorNone AssignmentDeleteBehavior = "none"
)

func PossibleValuesForAssignmentDeleteBehavior() []string {
	return []string{
		string(AssignmentDeleteBehaviorAll),
		string(AssignmentDeleteBehaviorNone),
	}
}

func (s *AssignmentDeleteBehavior) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseAssignmentDeleteBehavior(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseAssignmentDeleteBehavior(input string) (*AssignmentDeleteBehavior, error) {
	vals := map[string]AssignmentDeleteBehavior{
		"all":  AssignmentDeleteBehaviorAll,
		"none": AssignmentDeleteBehaviorNone,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := AssignmentDeleteBehavior(input)
	return &out, nil
}

type AssignmentLockMode string

const (
	AssignmentLockModeAllResourcesDoNotDelete AssignmentLockMode = "AllResourcesDoNotDelete"
	AssignmentLockModeAllResourcesReadOnly    AssignmentLockMode = "AllResourcesReadOnly"
	AssignmentLockModeNone                    AssignmentLockMode = "None"
)

func PossibleValuesForAssignmentLockMode() []string {
	return []string{
		string(AssignmentLockModeAllResourcesDoNotDelete),
		string(AssignmentLockModeAllResourcesReadOnly),
		string(AssignmentLockModeNone),
	}
}

func (s *AssignmentLockMode) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseAssignmentLockMode(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseAssignmentLockMode(input string) (*AssignmentLockMode, error) {
	vals := map[string]AssignmentLockMode{
		"allresourcesdonotdelete": AssignmentLockModeAllResourcesDoNotDelete,
		"allresourcesreadonly":    AssignmentLockModeAllResourcesReadOnly,
		"none":                    AssignmentLockModeNone,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := AssignmentLockMode(input)
	return &out, nil
}

type AssignmentProvisioningState string

const (
	AssignmentProvisioningStateCanceled   AssignmentProvisioningState = "canceled"
	AssignmentProvisioningStateCancelling AssignmentProvisioningState = "cancelling"
	AssignmentProvisioningStateCreating   AssignmentProvisioningState = "creating"
	AssignmentProvisioningStateDeleting   AssignmentProvisioningState = "deleting"
	AssignmentProvisioningStateDeploying  AssignmentProvisioningState = "deploying"
	AssignmentProvisioningStateFailed     AssignmentProvisioningState = "failed"
	AssignmentProvisioningStateLocking    AssignmentProvisioningState = "locking"
	AssignmentProvisioningStateSucceeded  AssignmentProvisioningState = "succeeded"
	AssignmentProvisioningStateValidating AssignmentProvisioningState = "validating"
	AssignmentProvisioningStateWaiting    AssignmentProvisioningState = "waiting"
)

func PossibleValuesForAssignmentProvisioningState() []string {
	return []string{
		string(AssignmentProvisioningStateCanceled),
		string(AssignmentProvisioningStateCancelling),
		string(AssignmentProvisioningStateCreating),
		string(AssignmentProvisioningStateDeleting),
		string(AssignmentProvisioningStateDeploying),
		string(AssignmentProvisioningStateFailed),
		string(AssignmentProvisioningStateLocking),
		string(AssignmentProvisioningStateSucceeded),
		string(AssignmentProvisioningStateValidating),
		string(AssignmentProvisioningStateWaiting),
	}
}

func (s *AssignmentProvisioningState) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseAssignmentProvisioningState(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseAssignmentProvisioningState(input string) (*AssignmentProvisioningState, error) {
	vals := map[string]AssignmentProvisioningState{
		"canceled":   AssignmentProvisioningStateCanceled,
		"cancelling": AssignmentProvisioningStateCancelling,
		"creating":   AssignmentProvisioningStateCreating,
		"deleting":   AssignmentProvisioningStateDeleting,
		"deploying":  AssignmentProvisioningStateDeploying,
		"failed":     AssignmentProvisioningStateFailed,
		"locking":    AssignmentProvisioningStateLocking,
		"succeeded":  AssignmentProvisioningStateSucceeded,
		"validating": AssignmentProvisioningStateValidating,
		"waiting":    AssignmentProvisioningStateWaiting,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := AssignmentProvisioningState(input)
	return &out, nil
}
