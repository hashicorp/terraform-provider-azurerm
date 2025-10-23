package daprcomponents

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type DaprComponentProvisioningState string

const (
	DaprComponentProvisioningStateCanceled   DaprComponentProvisioningState = "Canceled"
	DaprComponentProvisioningStateDeleting   DaprComponentProvisioningState = "Deleting"
	DaprComponentProvisioningStateFailed     DaprComponentProvisioningState = "Failed"
	DaprComponentProvisioningStateInProgress DaprComponentProvisioningState = "InProgress"
	DaprComponentProvisioningStateSucceeded  DaprComponentProvisioningState = "Succeeded"
)

func PossibleValuesForDaprComponentProvisioningState() []string {
	return []string{
		string(DaprComponentProvisioningStateCanceled),
		string(DaprComponentProvisioningStateDeleting),
		string(DaprComponentProvisioningStateFailed),
		string(DaprComponentProvisioningStateInProgress),
		string(DaprComponentProvisioningStateSucceeded),
	}
}

func (s *DaprComponentProvisioningState) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseDaprComponentProvisioningState(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseDaprComponentProvisioningState(input string) (*DaprComponentProvisioningState, error) {
	vals := map[string]DaprComponentProvisioningState{
		"canceled":   DaprComponentProvisioningStateCanceled,
		"deleting":   DaprComponentProvisioningStateDeleting,
		"failed":     DaprComponentProvisioningStateFailed,
		"inprogress": DaprComponentProvisioningStateInProgress,
		"succeeded":  DaprComponentProvisioningStateSucceeded,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := DaprComponentProvisioningState(input)
	return &out, nil
}
