package partnerregistrations

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

type PartnerRegistrationProvisioningState string

const (
	PartnerRegistrationProvisioningStateCanceled  PartnerRegistrationProvisioningState = "Canceled"
	PartnerRegistrationProvisioningStateCreating  PartnerRegistrationProvisioningState = "Creating"
	PartnerRegistrationProvisioningStateDeleting  PartnerRegistrationProvisioningState = "Deleting"
	PartnerRegistrationProvisioningStateFailed    PartnerRegistrationProvisioningState = "Failed"
	PartnerRegistrationProvisioningStateSucceeded PartnerRegistrationProvisioningState = "Succeeded"
	PartnerRegistrationProvisioningStateUpdating  PartnerRegistrationProvisioningState = "Updating"
)

func PossibleValuesForPartnerRegistrationProvisioningState() []string {
	return []string{
		string(PartnerRegistrationProvisioningStateCanceled),
		string(PartnerRegistrationProvisioningStateCreating),
		string(PartnerRegistrationProvisioningStateDeleting),
		string(PartnerRegistrationProvisioningStateFailed),
		string(PartnerRegistrationProvisioningStateSucceeded),
		string(PartnerRegistrationProvisioningStateUpdating),
	}
}

func parsePartnerRegistrationProvisioningState(input string) (*PartnerRegistrationProvisioningState, error) {
	vals := map[string]PartnerRegistrationProvisioningState{
		"canceled":  PartnerRegistrationProvisioningStateCanceled,
		"creating":  PartnerRegistrationProvisioningStateCreating,
		"deleting":  PartnerRegistrationProvisioningStateDeleting,
		"failed":    PartnerRegistrationProvisioningStateFailed,
		"succeeded": PartnerRegistrationProvisioningStateSucceeded,
		"updating":  PartnerRegistrationProvisioningStateUpdating,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := PartnerRegistrationProvisioningState(input)
	return &out, nil
}

type PartnerRegistrationVisibilityState string

const (
	PartnerRegistrationVisibilityStateGenerallyAvailable PartnerRegistrationVisibilityState = "GenerallyAvailable"
	PartnerRegistrationVisibilityStateHidden             PartnerRegistrationVisibilityState = "Hidden"
	PartnerRegistrationVisibilityStatePublicPreview      PartnerRegistrationVisibilityState = "PublicPreview"
)

func PossibleValuesForPartnerRegistrationVisibilityState() []string {
	return []string{
		string(PartnerRegistrationVisibilityStateGenerallyAvailable),
		string(PartnerRegistrationVisibilityStateHidden),
		string(PartnerRegistrationVisibilityStatePublicPreview),
	}
}

func parsePartnerRegistrationVisibilityState(input string) (*PartnerRegistrationVisibilityState, error) {
	vals := map[string]PartnerRegistrationVisibilityState{
		"generallyavailable": PartnerRegistrationVisibilityStateGenerallyAvailable,
		"hidden":             PartnerRegistrationVisibilityStateHidden,
		"publicpreview":      PartnerRegistrationVisibilityStatePublicPreview,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := PartnerRegistrationVisibilityState(input)
	return &out, nil
}
