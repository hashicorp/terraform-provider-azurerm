package disasterrecoveryconfigs

import "strings"

type ProvisioningStateDR string

const (
	ProvisioningStateDRAccepted  ProvisioningStateDR = "Accepted"
	ProvisioningStateDRFailed    ProvisioningStateDR = "Failed"
	ProvisioningStateDRSucceeded ProvisioningStateDR = "Succeeded"
)

func PossibleValuesForProvisioningStateDR() []string {
	return []string{
		"Accepted",
		"Failed",
		"Succeeded",
	}
}

func parseProvisioningStateDR(input string) (*ProvisioningStateDR, error) {
	vals := map[string]ProvisioningStateDR{
		"accepted":  "Accepted",
		"failed":    "Failed",
		"succeeded": "Succeeded",
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// it could be a new value - best effort convert this
	v := input

	out := ProvisioningStateDR(v)
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
		"Primary",
		"PrimaryNotReplicating",
		"Secondary",
	}
}

func parseRoleDisasterRecovery(input string) (*RoleDisasterRecovery, error) {
	vals := map[string]RoleDisasterRecovery{
		"primary":               "Primary",
		"primarynotreplicating": "PrimaryNotReplicating",
		"secondary":             "Secondary",
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// it could be a new value - best effort convert this
	v := input

	out := RoleDisasterRecovery(v)
	return &out, nil
}
