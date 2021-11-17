package partnertopics

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

type PartnerTopicActivationState string

const (
	PartnerTopicActivationStateActivated      PartnerTopicActivationState = "Activated"
	PartnerTopicActivationStateDeactivated    PartnerTopicActivationState = "Deactivated"
	PartnerTopicActivationStateNeverActivated PartnerTopicActivationState = "NeverActivated"
)

func PossibleValuesForPartnerTopicActivationState() []string {
	return []string{
		string(PartnerTopicActivationStateActivated),
		string(PartnerTopicActivationStateDeactivated),
		string(PartnerTopicActivationStateNeverActivated),
	}
}

func parsePartnerTopicActivationState(input string) (*PartnerTopicActivationState, error) {
	vals := map[string]PartnerTopicActivationState{
		"activated":      PartnerTopicActivationStateActivated,
		"deactivated":    PartnerTopicActivationStateDeactivated,
		"neveractivated": PartnerTopicActivationStateNeverActivated,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := PartnerTopicActivationState(input)
	return &out, nil
}

type PartnerTopicProvisioningState string

const (
	PartnerTopicProvisioningStateCanceled  PartnerTopicProvisioningState = "Canceled"
	PartnerTopicProvisioningStateCreating  PartnerTopicProvisioningState = "Creating"
	PartnerTopicProvisioningStateDeleting  PartnerTopicProvisioningState = "Deleting"
	PartnerTopicProvisioningStateFailed    PartnerTopicProvisioningState = "Failed"
	PartnerTopicProvisioningStateSucceeded PartnerTopicProvisioningState = "Succeeded"
	PartnerTopicProvisioningStateUpdating  PartnerTopicProvisioningState = "Updating"
)

func PossibleValuesForPartnerTopicProvisioningState() []string {
	return []string{
		string(PartnerTopicProvisioningStateCanceled),
		string(PartnerTopicProvisioningStateCreating),
		string(PartnerTopicProvisioningStateDeleting),
		string(PartnerTopicProvisioningStateFailed),
		string(PartnerTopicProvisioningStateSucceeded),
		string(PartnerTopicProvisioningStateUpdating),
	}
}

func parsePartnerTopicProvisioningState(input string) (*PartnerTopicProvisioningState, error) {
	vals := map[string]PartnerTopicProvisioningState{
		"canceled":  PartnerTopicProvisioningStateCanceled,
		"creating":  PartnerTopicProvisioningStateCreating,
		"deleting":  PartnerTopicProvisioningStateDeleting,
		"failed":    PartnerTopicProvisioningStateFailed,
		"succeeded": PartnerTopicProvisioningStateSucceeded,
		"updating":  PartnerTopicProvisioningStateUpdating,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := PartnerTopicProvisioningState(input)
	return &out, nil
}
