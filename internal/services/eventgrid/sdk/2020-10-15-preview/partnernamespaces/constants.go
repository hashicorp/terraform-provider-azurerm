package partnernamespaces

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

type PartnerNamespaceProvisioningState string

const (
	PartnerNamespaceProvisioningStateCanceled  PartnerNamespaceProvisioningState = "Canceled"
	PartnerNamespaceProvisioningStateCreating  PartnerNamespaceProvisioningState = "Creating"
	PartnerNamespaceProvisioningStateDeleting  PartnerNamespaceProvisioningState = "Deleting"
	PartnerNamespaceProvisioningStateFailed    PartnerNamespaceProvisioningState = "Failed"
	PartnerNamespaceProvisioningStateSucceeded PartnerNamespaceProvisioningState = "Succeeded"
	PartnerNamespaceProvisioningStateUpdating  PartnerNamespaceProvisioningState = "Updating"
)

func PossibleValuesForPartnerNamespaceProvisioningState() []string {
	return []string{
		string(PartnerNamespaceProvisioningStateCanceled),
		string(PartnerNamespaceProvisioningStateCreating),
		string(PartnerNamespaceProvisioningStateDeleting),
		string(PartnerNamespaceProvisioningStateFailed),
		string(PartnerNamespaceProvisioningStateSucceeded),
		string(PartnerNamespaceProvisioningStateUpdating),
	}
}

func parsePartnerNamespaceProvisioningState(input string) (*PartnerNamespaceProvisioningState, error) {
	vals := map[string]PartnerNamespaceProvisioningState{
		"canceled":  PartnerNamespaceProvisioningStateCanceled,
		"creating":  PartnerNamespaceProvisioningStateCreating,
		"deleting":  PartnerNamespaceProvisioningStateDeleting,
		"failed":    PartnerNamespaceProvisioningStateFailed,
		"succeeded": PartnerNamespaceProvisioningStateSucceeded,
		"updating":  PartnerNamespaceProvisioningStateUpdating,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := PartnerNamespaceProvisioningState(input)
	return &out, nil
}
