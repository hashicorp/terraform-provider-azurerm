package defaultparametergroups

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ParameterValueSource string

const (
	ParameterValueSourceSystemDefault ParameterValueSource = "SystemDefault"
	ParameterValueSourceUserOverride  ParameterValueSource = "UserOverride"
)

func PossibleValuesForParameterValueSource() []string {
	return []string{
		string(ParameterValueSourceSystemDefault),
		string(ParameterValueSourceUserOverride),
	}
}

func (s *ParameterValueSource) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseParameterValueSource(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseParameterValueSource(input string) (*ParameterValueSource, error) {
	vals := map[string]ParameterValueSource{
		"systemdefault": ParameterValueSourceSystemDefault,
		"useroverride":  ParameterValueSourceUserOverride,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := ParameterValueSource(input)
	return &out, nil
}

type PostgreSqlVersion string

const (
	PostgreSqlVersionOneEight PostgreSqlVersion = "18"
	PostgreSqlVersionOneSeven PostgreSqlVersion = "17"
)

func PossibleValuesForPostgreSqlVersion() []string {
	return []string{
		string(PostgreSqlVersionOneEight),
		string(PostgreSqlVersionOneSeven),
	}
}

func (s *PostgreSqlVersion) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parsePostgreSqlVersion(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parsePostgreSqlVersion(input string) (*PostgreSqlVersion, error) {
	vals := map[string]PostgreSqlVersion{
		"18": PostgreSqlVersionOneEight,
		"17": PostgreSqlVersionOneSeven,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := PostgreSqlVersion(input)
	return &out, nil
}

type ProvisioningState string

const (
	ProvisioningStateCanceled     ProvisioningState = "Canceled"
	ProvisioningStateFailed       ProvisioningState = "Failed"
	ProvisioningStateInProgress   ProvisioningState = "InProgress"
	ProvisioningStateProvisioning ProvisioningState = "Provisioning"
	ProvisioningStateSucceeded    ProvisioningState = "Succeeded"
	ProvisioningStateUpdating     ProvisioningState = "Updating"
)

func PossibleValuesForProvisioningState() []string {
	return []string{
		string(ProvisioningStateCanceled),
		string(ProvisioningStateFailed),
		string(ProvisioningStateInProgress),
		string(ProvisioningStateProvisioning),
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
		"canceled":     ProvisioningStateCanceled,
		"failed":       ProvisioningStateFailed,
		"inprogress":   ProvisioningStateInProgress,
		"provisioning": ProvisioningStateProvisioning,
		"succeeded":    ProvisioningStateSucceeded,
		"updating":     ProvisioningStateUpdating,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := ProvisioningState(input)
	return &out, nil
}
