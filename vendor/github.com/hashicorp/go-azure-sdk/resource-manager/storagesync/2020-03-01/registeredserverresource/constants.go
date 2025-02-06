package registeredserverresource

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type RegisteredServerAgentVersionStatus string

const (
	RegisteredServerAgentVersionStatusBlocked    RegisteredServerAgentVersionStatus = "Blocked"
	RegisteredServerAgentVersionStatusExpired    RegisteredServerAgentVersionStatus = "Expired"
	RegisteredServerAgentVersionStatusNearExpiry RegisteredServerAgentVersionStatus = "NearExpiry"
	RegisteredServerAgentVersionStatusOk         RegisteredServerAgentVersionStatus = "Ok"
)

func PossibleValuesForRegisteredServerAgentVersionStatus() []string {
	return []string{
		string(RegisteredServerAgentVersionStatusBlocked),
		string(RegisteredServerAgentVersionStatusExpired),
		string(RegisteredServerAgentVersionStatusNearExpiry),
		string(RegisteredServerAgentVersionStatusOk),
	}
}

func (s *RegisteredServerAgentVersionStatus) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseRegisteredServerAgentVersionStatus(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseRegisteredServerAgentVersionStatus(input string) (*RegisteredServerAgentVersionStatus, error) {
	vals := map[string]RegisteredServerAgentVersionStatus{
		"blocked":    RegisteredServerAgentVersionStatusBlocked,
		"expired":    RegisteredServerAgentVersionStatusExpired,
		"nearexpiry": RegisteredServerAgentVersionStatusNearExpiry,
		"ok":         RegisteredServerAgentVersionStatusOk,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := RegisteredServerAgentVersionStatus(input)
	return &out, nil
}
