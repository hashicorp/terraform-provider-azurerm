package transparentdataencryptions

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type TransparentDataEncryptionScanState string

const (
	TransparentDataEncryptionScanStateAborted   TransparentDataEncryptionScanState = "Aborted"
	TransparentDataEncryptionScanStateCompleted TransparentDataEncryptionScanState = "Completed"
	TransparentDataEncryptionScanStateNone      TransparentDataEncryptionScanState = "None"
	TransparentDataEncryptionScanStateResume    TransparentDataEncryptionScanState = "Resume"
	TransparentDataEncryptionScanStateRunning   TransparentDataEncryptionScanState = "Running"
	TransparentDataEncryptionScanStateSuspend   TransparentDataEncryptionScanState = "Suspend"
)

func PossibleValuesForTransparentDataEncryptionScanState() []string {
	return []string{
		string(TransparentDataEncryptionScanStateAborted),
		string(TransparentDataEncryptionScanStateCompleted),
		string(TransparentDataEncryptionScanStateNone),
		string(TransparentDataEncryptionScanStateResume),
		string(TransparentDataEncryptionScanStateRunning),
		string(TransparentDataEncryptionScanStateSuspend),
	}
}

func (s *TransparentDataEncryptionScanState) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseTransparentDataEncryptionScanState(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseTransparentDataEncryptionScanState(input string) (*TransparentDataEncryptionScanState, error) {
	vals := map[string]TransparentDataEncryptionScanState{
		"aborted":   TransparentDataEncryptionScanStateAborted,
		"completed": TransparentDataEncryptionScanStateCompleted,
		"none":      TransparentDataEncryptionScanStateNone,
		"resume":    TransparentDataEncryptionScanStateResume,
		"running":   TransparentDataEncryptionScanStateRunning,
		"suspend":   TransparentDataEncryptionScanStateSuspend,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := TransparentDataEncryptionScanState(input)
	return &out, nil
}

type TransparentDataEncryptionState string

const (
	TransparentDataEncryptionStateDisabled TransparentDataEncryptionState = "Disabled"
	TransparentDataEncryptionStateEnabled  TransparentDataEncryptionState = "Enabled"
)

func PossibleValuesForTransparentDataEncryptionState() []string {
	return []string{
		string(TransparentDataEncryptionStateDisabled),
		string(TransparentDataEncryptionStateEnabled),
	}
}

func (s *TransparentDataEncryptionState) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseTransparentDataEncryptionState(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseTransparentDataEncryptionState(input string) (*TransparentDataEncryptionState, error) {
	vals := map[string]TransparentDataEncryptionState{
		"disabled": TransparentDataEncryptionStateDisabled,
		"enabled":  TransparentDataEncryptionStateEnabled,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := TransparentDataEncryptionState(input)
	return &out, nil
}
