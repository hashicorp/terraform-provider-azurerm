package saplandscapemonitor

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type SapLandscapeMonitorProvisioningState string

const (
	SapLandscapeMonitorProvisioningStateAccepted  SapLandscapeMonitorProvisioningState = "Accepted"
	SapLandscapeMonitorProvisioningStateCanceled  SapLandscapeMonitorProvisioningState = "Canceled"
	SapLandscapeMonitorProvisioningStateCreated   SapLandscapeMonitorProvisioningState = "Created"
	SapLandscapeMonitorProvisioningStateFailed    SapLandscapeMonitorProvisioningState = "Failed"
	SapLandscapeMonitorProvisioningStateSucceeded SapLandscapeMonitorProvisioningState = "Succeeded"
)

func PossibleValuesForSapLandscapeMonitorProvisioningState() []string {
	return []string{
		string(SapLandscapeMonitorProvisioningStateAccepted),
		string(SapLandscapeMonitorProvisioningStateCanceled),
		string(SapLandscapeMonitorProvisioningStateCreated),
		string(SapLandscapeMonitorProvisioningStateFailed),
		string(SapLandscapeMonitorProvisioningStateSucceeded),
	}
}

func (s *SapLandscapeMonitorProvisioningState) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseSapLandscapeMonitorProvisioningState(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseSapLandscapeMonitorProvisioningState(input string) (*SapLandscapeMonitorProvisioningState, error) {
	vals := map[string]SapLandscapeMonitorProvisioningState{
		"accepted":  SapLandscapeMonitorProvisioningStateAccepted,
		"canceled":  SapLandscapeMonitorProvisioningStateCanceled,
		"created":   SapLandscapeMonitorProvisioningStateCreated,
		"failed":    SapLandscapeMonitorProvisioningStateFailed,
		"succeeded": SapLandscapeMonitorProvisioningStateSucceeded,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := SapLandscapeMonitorProvisioningState(input)
	return &out, nil
}
