package disasterrecoveryconfigs

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ProvisioningStateDR string

const (
	ProvisioningStateDRAccepted  ProvisioningStateDR = "Accepted"
	ProvisioningStateDRFailed    ProvisioningStateDR = "Failed"
	ProvisioningStateDRSucceeded ProvisioningStateDR = "Succeeded"
)

func PossibleValuesForProvisioningStateDR() []string {
	return []string{
		string(ProvisioningStateDRAccepted),
		string(ProvisioningStateDRFailed),
		string(ProvisioningStateDRSucceeded),
	}
}

func (s *ProvisioningStateDR) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseProvisioningStateDR(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseProvisioningStateDR(input string) (*ProvisioningStateDR, error) {
	vals := map[string]ProvisioningStateDR{
		"accepted":  ProvisioningStateDRAccepted,
		"failed":    ProvisioningStateDRFailed,
		"succeeded": ProvisioningStateDRSucceeded,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := ProvisioningStateDR(input)
	return &out, nil
}

type RoleDisasterRecovery string

const (
	RoleDisasterRecoveryPrimary               RoleDisasterRecovery = "Primary"
	RoleDisasterRecoveryPrimaryNotReplicating RoleDisasterRecovery = "PrimaryNotReplicating"
	RoleDisasterRecoverySecondary             RoleDisasterRecovery = "Secondary"
)

func PossibleValuesForRoleDisasterRecovery() []string {
	return []string{
		string(RoleDisasterRecoveryPrimary),
		string(RoleDisasterRecoveryPrimaryNotReplicating),
		string(RoleDisasterRecoverySecondary),
	}
}

func (s *RoleDisasterRecovery) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseRoleDisasterRecovery(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseRoleDisasterRecovery(input string) (*RoleDisasterRecovery, error) {
	vals := map[string]RoleDisasterRecovery{
		"primary":               RoleDisasterRecoveryPrimary,
		"primarynotreplicating": RoleDisasterRecoveryPrimaryNotReplicating,
		"secondary":             RoleDisasterRecoverySecondary,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := RoleDisasterRecovery(input)
	return &out, nil
}
