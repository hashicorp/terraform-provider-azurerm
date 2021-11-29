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

type VCoreProvisioningState string

const (
	VCoreProvisioningStateSucceeded VCoreProvisioningState = "Succeeded"
)

func PossibleValuesForVCoreProvisioningState() []string {
	return []string{
		string(VCoreProvisioningStateSucceeded),
	}
}

func parseVCoreProvisioningState(input string) (*VCoreProvisioningState, error) {
	vals := map[string]VCoreProvisioningState{
		"succeeded": VCoreProvisioningStateSucceeded,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := VCoreProvisioningState(input)
	return &out, nil
}

type VCoreSkuTier string

const (
	VCoreSkuTierAutoScale VCoreSkuTier = "AutoScale"
)

func PossibleValuesForVCoreSkuTier() []string {
	return []string{
		string(VCoreSkuTierAutoScale),
	}
}

func parseVCoreSkuTier(input string) (*VCoreSkuTier, error) {
	vals := map[string]VCoreSkuTier{
		"autoscale": VCoreSkuTierAutoScale,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := VCoreSkuTier(input)
	return &out, nil
}
