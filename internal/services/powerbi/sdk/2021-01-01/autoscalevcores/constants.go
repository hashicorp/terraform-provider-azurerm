package autoscalevcores

import "strings"

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

type VCoreProvisioningState string

const (
	VCoreProvisioningStateSucceeded VCoreProvisioningState = "Succeeded"
)

func PossibleValuesForVCoreProvisioningState() []string {
	return []string{
		"Succeeded",
	}
}

func parseVCoreProvisioningState(input string) (*VCoreProvisioningState, error) {
	vals := map[string]VCoreProvisioningState{
		"succeeded": "Succeeded",
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// it could be a new value - best effort convert this
	v := input

	out := VCoreProvisioningState(v)
	return &out, nil
}

type VCoreSkuTier string

const (
	VCoreSkuTierAutoScale VCoreSkuTier = "AutoScale"
)

func PossibleValuesForVCoreSkuTier() []string {
	return []string{
		"AutoScale",
	}
}

func parseVCoreSkuTier(input string) (*VCoreSkuTier, error) {
	vals := map[string]VCoreSkuTier{
		"autoscale": "AutoScale",
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// it could be a new value - best effort convert this
	v := input

	out := VCoreSkuTier(v)
	return &out, nil
}
