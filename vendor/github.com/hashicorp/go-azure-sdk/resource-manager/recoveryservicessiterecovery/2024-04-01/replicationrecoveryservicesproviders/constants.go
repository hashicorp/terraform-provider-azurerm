package replicationrecoveryservicesproviders

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type AgentVersionStatus string

const (
	AgentVersionStatusDeprecated             AgentVersionStatus = "Deprecated"
	AgentVersionStatusNotSupported           AgentVersionStatus = "NotSupported"
	AgentVersionStatusSecurityUpdateRequired AgentVersionStatus = "SecurityUpdateRequired"
	AgentVersionStatusSupported              AgentVersionStatus = "Supported"
	AgentVersionStatusUpdateRequired         AgentVersionStatus = "UpdateRequired"
)

func PossibleValuesForAgentVersionStatus() []string {
	return []string{
		string(AgentVersionStatusDeprecated),
		string(AgentVersionStatusNotSupported),
		string(AgentVersionStatusSecurityUpdateRequired),
		string(AgentVersionStatusSupported),
		string(AgentVersionStatusUpdateRequired),
	}
}

func (s *AgentVersionStatus) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseAgentVersionStatus(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseAgentVersionStatus(input string) (*AgentVersionStatus, error) {
	vals := map[string]AgentVersionStatus{
		"deprecated":             AgentVersionStatusDeprecated,
		"notsupported":           AgentVersionStatusNotSupported,
		"securityupdaterequired": AgentVersionStatusSecurityUpdateRequired,
		"supported":              AgentVersionStatusSupported,
		"updaterequired":         AgentVersionStatusUpdateRequired,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := AgentVersionStatus(input)
	return &out, nil
}

type HealthErrorCustomerResolvability string

const (
	HealthErrorCustomerResolvabilityAllowed    HealthErrorCustomerResolvability = "Allowed"
	HealthErrorCustomerResolvabilityNotAllowed HealthErrorCustomerResolvability = "NotAllowed"
)

func PossibleValuesForHealthErrorCustomerResolvability() []string {
	return []string{
		string(HealthErrorCustomerResolvabilityAllowed),
		string(HealthErrorCustomerResolvabilityNotAllowed),
	}
}

func (s *HealthErrorCustomerResolvability) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseHealthErrorCustomerResolvability(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseHealthErrorCustomerResolvability(input string) (*HealthErrorCustomerResolvability, error) {
	vals := map[string]HealthErrorCustomerResolvability{
		"allowed":    HealthErrorCustomerResolvabilityAllowed,
		"notallowed": HealthErrorCustomerResolvabilityNotAllowed,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := HealthErrorCustomerResolvability(input)
	return &out, nil
}
