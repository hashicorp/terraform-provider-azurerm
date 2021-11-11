package disasterrecoveryconfigs

import "strings"

type CreatedByType string

const (
	CreatedByTypeApplication     CreatedByType = "Application"
	CreatedByTypeKey             CreatedByType = "Key"
	CreatedByTypeManagedIdentity CreatedByType = "ManagedIdentity"
	CreatedByTypeUser            CreatedByType = "User"
)

func PossibleValuesForCreatedByType() []string {
	return []string{
		string(CreatedByTypeApplication),
		string(CreatedByTypeKey),
		string(CreatedByTypeManagedIdentity),
		string(CreatedByTypeUser),
	}
}

func parseCreatedByType(input string) (*CreatedByType, error) {
	vals := map[string]CreatedByType{
		"application":     CreatedByTypeApplication,
		"key":             CreatedByTypeKey,
		"managedidentity": CreatedByTypeManagedIdentity,
		"user":            CreatedByTypeUser,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := CreatedByType(input)
	return &out, nil
}

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
