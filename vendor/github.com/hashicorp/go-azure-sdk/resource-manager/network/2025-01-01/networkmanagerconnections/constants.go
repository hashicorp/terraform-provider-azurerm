package networkmanagerconnections

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ScopeConnectionState string

const (
	ScopeConnectionStateConflict  ScopeConnectionState = "Conflict"
	ScopeConnectionStateConnected ScopeConnectionState = "Connected"
	ScopeConnectionStatePending   ScopeConnectionState = "Pending"
	ScopeConnectionStateRejected  ScopeConnectionState = "Rejected"
	ScopeConnectionStateRevoked   ScopeConnectionState = "Revoked"
)

func PossibleValuesForScopeConnectionState() []string {
	return []string{
		string(ScopeConnectionStateConflict),
		string(ScopeConnectionStateConnected),
		string(ScopeConnectionStatePending),
		string(ScopeConnectionStateRejected),
		string(ScopeConnectionStateRevoked),
	}
}

func (s *ScopeConnectionState) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseScopeConnectionState(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseScopeConnectionState(input string) (*ScopeConnectionState, error) {
	vals := map[string]ScopeConnectionState{
		"conflict":  ScopeConnectionStateConflict,
		"connected": ScopeConnectionStateConnected,
		"pending":   ScopeConnectionStatePending,
		"rejected":  ScopeConnectionStateRejected,
		"revoked":   ScopeConnectionStateRevoked,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := ScopeConnectionState(input)
	return &out, nil
}
