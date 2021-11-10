package capacities

import "strings"

type CapacityProvisioningState string

const (
	CapacityProvisioningStateDeleting     CapacityProvisioningState = "Deleting"
	CapacityProvisioningStateFailed       CapacityProvisioningState = "Failed"
	CapacityProvisioningStatePaused       CapacityProvisioningState = "Paused"
	CapacityProvisioningStatePausing      CapacityProvisioningState = "Pausing"
	CapacityProvisioningStatePreparing    CapacityProvisioningState = "Preparing"
	CapacityProvisioningStateProvisioning CapacityProvisioningState = "Provisioning"
	CapacityProvisioningStateResuming     CapacityProvisioningState = "Resuming"
	CapacityProvisioningStateScaling      CapacityProvisioningState = "Scaling"
	CapacityProvisioningStateSucceeded    CapacityProvisioningState = "Succeeded"
	CapacityProvisioningStateSuspended    CapacityProvisioningState = "Suspended"
	CapacityProvisioningStateSuspending   CapacityProvisioningState = "Suspending"
	CapacityProvisioningStateUpdating     CapacityProvisioningState = "Updating"
)

func PossibleValuesForCapacityProvisioningState() []string {
	return []string{
		"Deleting",
		"Failed",
		"Paused",
		"Pausing",
		"Preparing",
		"Provisioning",
		"Resuming",
		"Scaling",
		"Succeeded",
		"Suspended",
		"Suspending",
		"Updating",
	}
}

func parseCapacityProvisioningState(input string) (*CapacityProvisioningState, error) {
	vals := map[string]CapacityProvisioningState{
		"deleting":     "Deleting",
		"failed":       "Failed",
		"paused":       "Paused",
		"pausing":      "Pausing",
		"preparing":    "Preparing",
		"provisioning": "Provisioning",
		"resuming":     "Resuming",
		"scaling":      "Scaling",
		"succeeded":    "Succeeded",
		"suspended":    "Suspended",
		"suspending":   "Suspending",
		"updating":     "Updating",
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// it could be a new value - best effort convert this
	v := input

	out := CapacityProvisioningState(v)
	return &out, nil
}

type CapacitySkuTier string

const (
	CapacitySkuTierAutoPremiumHost CapacitySkuTier = "AutoPremiumHost"
	CapacitySkuTierPBIEAzure       CapacitySkuTier = "PBIE_Azure"
	CapacitySkuTierPremium         CapacitySkuTier = "Premium"
)

func PossibleValuesForCapacitySkuTier() []string {
	return []string{
		"AutoPremiumHost",
		"PBIE_Azure",
		"Premium",
	}
}

func parseCapacitySkuTier(input string) (*CapacitySkuTier, error) {
	vals := map[string]CapacitySkuTier{
		"autopremiumhost": "AutoPremiumHost",
		"pbieazure":       "PBIE_Azure",
		"premium":         "Premium",
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// it could be a new value - best effort convert this
	v := input

	out := CapacitySkuTier(v)
	return &out, nil
}

type IdentityType string

const (
	IdentityTypeApplication     IdentityType = "Application"
	IdentityTypeKey             IdentityType = "Key"
	IdentityTypeManagedIdentity IdentityType = "ManagedIdentity"
	IdentityTypeUser            IdentityType = "User"
)

func PossibleValuesForIdentityType() []string {
	return []string{
		"Application",
		"Key",
		"ManagedIdentity",
		"User",
	}
}

func parseIdentityType(input string) (*IdentityType, error) {
	vals := map[string]IdentityType{
		"application":     "Application",
		"key":             "Key",
		"managedidentity": "ManagedIdentity",
		"user":            "User",
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// it could be a new value - best effort convert this
	v := input

	out := IdentityType(v)
	return &out, nil
}

type Mode string

const (
	ModeGenOne Mode = "Gen1"
	ModeGenTwo Mode = "Gen2"
)

func PossibleValuesForMode() []string {
	return []string{
		"Gen1",
		"Gen2",
	}
}

func parseMode(input string) (*Mode, error) {
	vals := map[string]Mode{
		"genone": "Gen1",
		"gentwo": "Gen2",
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// it could be a new value - best effort convert this
	v := input

	out := Mode(v)
	return &out, nil
}

type State string

const (
	StateDeleting     State = "Deleting"
	StateFailed       State = "Failed"
	StatePaused       State = "Paused"
	StatePausing      State = "Pausing"
	StatePreparing    State = "Preparing"
	StateProvisioning State = "Provisioning"
	StateResuming     State = "Resuming"
	StateScaling      State = "Scaling"
	StateSucceeded    State = "Succeeded"
	StateSuspended    State = "Suspended"
	StateSuspending   State = "Suspending"
	StateUpdating     State = "Updating"
)

func PossibleValuesForState() []string {
	return []string{
		"Deleting",
		"Failed",
		"Paused",
		"Pausing",
		"Preparing",
		"Provisioning",
		"Resuming",
		"Scaling",
		"Succeeded",
		"Suspended",
		"Suspending",
		"Updating",
	}
}

func parseState(input string) (*State, error) {
	vals := map[string]State{
		"deleting":     "Deleting",
		"failed":       "Failed",
		"paused":       "Paused",
		"pausing":      "Pausing",
		"preparing":    "Preparing",
		"provisioning": "Provisioning",
		"resuming":     "Resuming",
		"scaling":      "Scaling",
		"succeeded":    "Succeeded",
		"suspended":    "Suspended",
		"suspending":   "Suspending",
		"updating":     "Updating",
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// it could be a new value - best effort convert this
	v := input

	out := State(v)
	return &out, nil
}
