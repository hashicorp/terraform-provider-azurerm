package fhirservices

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

type FhirServiceKind string

const (
	FhirServiceKindFhirNegativeRFour    FhirServiceKind = "fhir-R4"
	FhirServiceKindFhirNegativeStuThree FhirServiceKind = "fhir-Stu3"
)

func PossibleValuesForFhirServiceKind() []string {
	return []string{
		string(FhirServiceKindFhirNegativeRFour),
		string(FhirServiceKindFhirNegativeStuThree),
	}
}

func parseFhirServiceKind(input string) (*FhirServiceKind, error) {
	vals := map[string]FhirServiceKind{
		"fhir-r4":   FhirServiceKindFhirNegativeRFour,
		"fhir-stu3": FhirServiceKindFhirNegativeStuThree,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := FhirServiceKind(input)
	return &out, nil
}

type ManagedServiceIdentityType string

const (
	ManagedServiceIdentityTypeNone           ManagedServiceIdentityType = "None"
	ManagedServiceIdentityTypeSystemAssigned ManagedServiceIdentityType = "SystemAssigned"
)

func PossibleValuesForManagedServiceIdentityType() []string {
	return []string{
		string(ManagedServiceIdentityTypeNone),
		string(ManagedServiceIdentityTypeSystemAssigned),
	}
}

func parseManagedServiceIdentityType(input string) (*ManagedServiceIdentityType, error) {
	vals := map[string]ManagedServiceIdentityType{
		"none":           ManagedServiceIdentityTypeNone,
		"systemassigned": ManagedServiceIdentityTypeSystemAssigned,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := ManagedServiceIdentityType(input)
	return &out, nil
}

type ProvisioningState string

const (
	ProvisioningStateAccepted          ProvisioningState = "Accepted"
	ProvisioningStateCanceled          ProvisioningState = "Canceled"
	ProvisioningStateCreating          ProvisioningState = "Creating"
	ProvisioningStateDeleting          ProvisioningState = "Deleting"
	ProvisioningStateDeprovisioned     ProvisioningState = "Deprovisioned"
	ProvisioningStateFailed            ProvisioningState = "Failed"
	ProvisioningStateMoving            ProvisioningState = "Moving"
	ProvisioningStateSucceeded         ProvisioningState = "Succeeded"
	ProvisioningStateSuspended         ProvisioningState = "Suspended"
	ProvisioningStateSystemMaintenance ProvisioningState = "SystemMaintenance"
	ProvisioningStateUpdating          ProvisioningState = "Updating"
	ProvisioningStateVerifying         ProvisioningState = "Verifying"
	ProvisioningStateWarned            ProvisioningState = "Warned"
)

func PossibleValuesForProvisioningState() []string {
	return []string{
		string(ProvisioningStateAccepted),
		string(ProvisioningStateCanceled),
		string(ProvisioningStateCreating),
		string(ProvisioningStateDeleting),
		string(ProvisioningStateDeprovisioned),
		string(ProvisioningStateFailed),
		string(ProvisioningStateMoving),
		string(ProvisioningStateSucceeded),
		string(ProvisioningStateSuspended),
		string(ProvisioningStateSystemMaintenance),
		string(ProvisioningStateUpdating),
		string(ProvisioningStateVerifying),
		string(ProvisioningStateWarned),
	}
}

func parseProvisioningState(input string) (*ProvisioningState, error) {
	vals := map[string]ProvisioningState{
		"accepted":          ProvisioningStateAccepted,
		"canceled":          ProvisioningStateCanceled,
		"creating":          ProvisioningStateCreating,
		"deleting":          ProvisioningStateDeleting,
		"deprovisioned":     ProvisioningStateDeprovisioned,
		"failed":            ProvisioningStateFailed,
		"moving":            ProvisioningStateMoving,
		"succeeded":         ProvisioningStateSucceeded,
		"suspended":         ProvisioningStateSuspended,
		"systemmaintenance": ProvisioningStateSystemMaintenance,
		"updating":          ProvisioningStateUpdating,
		"verifying":         ProvisioningStateVerifying,
		"warned":            ProvisioningStateWarned,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := ProvisioningState(input)
	return &out, nil
}
