package disasterrecoveryconfigs

import "strings"

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
