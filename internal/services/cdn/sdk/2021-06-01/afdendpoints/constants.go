package afdendpoints

import "strings"

type AfdProvisioningState string

const (
	AfdProvisioningStateCreating  AfdProvisioningState = "Creating"
	AfdProvisioningStateDeleting  AfdProvisioningState = "Deleting"
	AfdProvisioningStateFailed    AfdProvisioningState = "Failed"
	AfdProvisioningStateSucceeded AfdProvisioningState = "Succeeded"
	AfdProvisioningStateUpdating  AfdProvisioningState = "Updating"
)

func PossibleValuesForAfdProvisioningState() []string {
	return []string{
		string(AfdProvisioningStateCreating),
		string(AfdProvisioningStateDeleting),
		string(AfdProvisioningStateFailed),
		string(AfdProvisioningStateSucceeded),
		string(AfdProvisioningStateUpdating),
	}
}

func parseAfdProvisioningState(input string) (*AfdProvisioningState, error) {
	vals := map[string]AfdProvisioningState{
		"creating":  AfdProvisioningStateCreating,
		"deleting":  AfdProvisioningStateDeleting,
		"failed":    AfdProvisioningStateFailed,
		"succeeded": AfdProvisioningStateSucceeded,
		"updating":  AfdProvisioningStateUpdating,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := AfdProvisioningState(input)
	return &out, nil
}

type DeploymentStatus string

const (
	DeploymentStatusFailed     DeploymentStatus = "Failed"
	DeploymentStatusInProgress DeploymentStatus = "InProgress"
	DeploymentStatusNotStarted DeploymentStatus = "NotStarted"
	DeploymentStatusSucceeded  DeploymentStatus = "Succeeded"
)

func PossibleValuesForDeploymentStatus() []string {
	return []string{
		string(DeploymentStatusFailed),
		string(DeploymentStatusInProgress),
		string(DeploymentStatusNotStarted),
		string(DeploymentStatusSucceeded),
	}
}

func parseDeploymentStatus(input string) (*DeploymentStatus, error) {
	vals := map[string]DeploymentStatus{
		"failed":     DeploymentStatusFailed,
		"inprogress": DeploymentStatusInProgress,
		"notstarted": DeploymentStatusNotStarted,
		"succeeded":  DeploymentStatusSucceeded,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := DeploymentStatus(input)
	return &out, nil
}

type EnabledState string

const (
	EnabledStateDisabled EnabledState = "Disabled"
	EnabledStateEnabled  EnabledState = "Enabled"
)

func PossibleValuesForEnabledState() []string {
	return []string{
		string(EnabledStateDisabled),
		string(EnabledStateEnabled),
	}
}

func parseEnabledState(input string) (*EnabledState, error) {
	vals := map[string]EnabledState{
		"disabled": EnabledStateDisabled,
		"enabled":  EnabledStateEnabled,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := EnabledState(input)
	return &out, nil
}

type IdentityType string

const (
	IdentityTypeApplication     IdentityType = "application"
	IdentityTypeKey             IdentityType = "key"
	IdentityTypeManagedIdentity IdentityType = "managedIdentity"
	IdentityTypeUser            IdentityType = "user"
)

func PossibleValuesForIdentityType() []string {
	return []string{
		string(IdentityTypeApplication),
		string(IdentityTypeKey),
		string(IdentityTypeManagedIdentity),
		string(IdentityTypeUser),
	}
}

func parseIdentityType(input string) (*IdentityType, error) {
	vals := map[string]IdentityType{
		"application":     IdentityTypeApplication,
		"key":             IdentityTypeKey,
		"managedidentity": IdentityTypeManagedIdentity,
		"user":            IdentityTypeUser,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := IdentityType(input)
	return &out, nil
}

type UsageUnit string

const (
	UsageUnitCount UsageUnit = "Count"
)

func PossibleValuesForUsageUnit() []string {
	return []string{
		string(UsageUnitCount),
	}
}

func parseUsageUnit(input string) (*UsageUnit, error) {
	vals := map[string]UsageUnit{
		"count": UsageUnitCount,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := UsageUnit(input)
	return &out, nil
}
