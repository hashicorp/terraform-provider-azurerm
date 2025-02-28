package networkstatus

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ConnectivityStatusType string

const (
	ConnectivityStatusTypeFailure      ConnectivityStatusType = "failure"
	ConnectivityStatusTypeInitializing ConnectivityStatusType = "initializing"
	ConnectivityStatusTypeSuccess      ConnectivityStatusType = "success"
)

func PossibleValuesForConnectivityStatusType() []string {
	return []string{
		string(ConnectivityStatusTypeFailure),
		string(ConnectivityStatusTypeInitializing),
		string(ConnectivityStatusTypeSuccess),
	}
}

func (s *ConnectivityStatusType) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseConnectivityStatusType(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseConnectivityStatusType(input string) (*ConnectivityStatusType, error) {
	vals := map[string]ConnectivityStatusType{
		"failure":      ConnectivityStatusTypeFailure,
		"initializing": ConnectivityStatusTypeInitializing,
		"success":      ConnectivityStatusTypeSuccess,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := ConnectivityStatusType(input)
	return &out, nil
}
