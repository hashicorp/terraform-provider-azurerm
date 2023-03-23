package hostpool

import "strings"

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type HostPoolType string

const (
	HostPoolTypeBYODesktop HostPoolType = "BYODesktop"
	HostPoolTypePersonal   HostPoolType = "Personal"
	HostPoolTypePooled     HostPoolType = "Pooled"
)

func PossibleValuesForHostPoolType() []string {
	return []string{
		string(HostPoolTypeBYODesktop),
		string(HostPoolTypePersonal),
		string(HostPoolTypePooled),
	}
}

func parseHostPoolType(input string) (*HostPoolType, error) {
	vals := map[string]HostPoolType{
		"byodesktop": HostPoolTypeBYODesktop,
		"personal":   HostPoolTypePersonal,
		"pooled":     HostPoolTypePooled,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := HostPoolType(input)
	return &out, nil
}

type LoadBalancerType string

const (
	LoadBalancerTypeBreadthFirst LoadBalancerType = "BreadthFirst"
	LoadBalancerTypeDepthFirst   LoadBalancerType = "DepthFirst"
	LoadBalancerTypePersistent   LoadBalancerType = "Persistent"
)

func PossibleValuesForLoadBalancerType() []string {
	return []string{
		string(LoadBalancerTypeBreadthFirst),
		string(LoadBalancerTypeDepthFirst),
		string(LoadBalancerTypePersistent),
	}
}

func parseLoadBalancerType(input string) (*LoadBalancerType, error) {
	vals := map[string]LoadBalancerType{
		"breadthfirst": LoadBalancerTypeBreadthFirst,
		"depthfirst":   LoadBalancerTypeDepthFirst,
		"persistent":   LoadBalancerTypePersistent,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := LoadBalancerType(input)
	return &out, nil
}

type Operation string

const (
	OperationComplete Operation = "Complete"
	OperationHide     Operation = "Hide"
	OperationRevoke   Operation = "Revoke"
	OperationStart    Operation = "Start"
	OperationUnhide   Operation = "Unhide"
)

func PossibleValuesForOperation() []string {
	return []string{
		string(OperationComplete),
		string(OperationHide),
		string(OperationRevoke),
		string(OperationStart),
		string(OperationUnhide),
	}
}

func parseOperation(input string) (*Operation, error) {
	vals := map[string]Operation{
		"complete": OperationComplete,
		"hide":     OperationHide,
		"revoke":   OperationRevoke,
		"start":    OperationStart,
		"unhide":   OperationUnhide,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := Operation(input)
	return &out, nil
}

type PersonalDesktopAssignmentType string

const (
	PersonalDesktopAssignmentTypeAutomatic PersonalDesktopAssignmentType = "Automatic"
	PersonalDesktopAssignmentTypeDirect    PersonalDesktopAssignmentType = "Direct"
)

func PossibleValuesForPersonalDesktopAssignmentType() []string {
	return []string{
		string(PersonalDesktopAssignmentTypeAutomatic),
		string(PersonalDesktopAssignmentTypeDirect),
	}
}

func parsePersonalDesktopAssignmentType(input string) (*PersonalDesktopAssignmentType, error) {
	vals := map[string]PersonalDesktopAssignmentType{
		"automatic": PersonalDesktopAssignmentTypeAutomatic,
		"direct":    PersonalDesktopAssignmentTypeDirect,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := PersonalDesktopAssignmentType(input)
	return &out, nil
}

type PreferredAppGroupType string

const (
	PreferredAppGroupTypeDesktop          PreferredAppGroupType = "Desktop"
	PreferredAppGroupTypeNone             PreferredAppGroupType = "None"
	PreferredAppGroupTypeRailApplications PreferredAppGroupType = "RailApplications"
)

func PossibleValuesForPreferredAppGroupType() []string {
	return []string{
		string(PreferredAppGroupTypeDesktop),
		string(PreferredAppGroupTypeNone),
		string(PreferredAppGroupTypeRailApplications),
	}
}

func parsePreferredAppGroupType(input string) (*PreferredAppGroupType, error) {
	vals := map[string]PreferredAppGroupType{
		"desktop":          PreferredAppGroupTypeDesktop,
		"none":             PreferredAppGroupTypeNone,
		"railapplications": PreferredAppGroupTypeRailApplications,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := PreferredAppGroupType(input)
	return &out, nil
}

type PublicNetworkAccess string

const (
	PublicNetworkAccessDisabled PublicNetworkAccess = "Disabled"
	PublicNetworkAccessEnabled  PublicNetworkAccess = "Enabled"
)

func PossibleValuesForPublicNetworkAccess() []string {
	return []string{
		string(PublicNetworkAccessDisabled),
		string(PublicNetworkAccessEnabled),
	}
}

func parsePublicNetworkAccess(input string) (*PublicNetworkAccess, error) {
	vals := map[string]PublicNetworkAccess{
		"disabled": PublicNetworkAccessDisabled,
		"enabled":  PublicNetworkAccessEnabled,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := PublicNetworkAccess(input)
	return &out, nil
}

type RegistrationTokenOperation string

const (
	RegistrationTokenOperationDelete RegistrationTokenOperation = "Delete"
	RegistrationTokenOperationNone   RegistrationTokenOperation = "None"
	RegistrationTokenOperationUpdate RegistrationTokenOperation = "Update"
)

func PossibleValuesForRegistrationTokenOperation() []string {
	return []string{
		string(RegistrationTokenOperationDelete),
		string(RegistrationTokenOperationNone),
		string(RegistrationTokenOperationUpdate),
	}
}

func parseRegistrationTokenOperation(input string) (*RegistrationTokenOperation, error) {
	vals := map[string]RegistrationTokenOperation{
		"delete": RegistrationTokenOperationDelete,
		"none":   RegistrationTokenOperationNone,
		"update": RegistrationTokenOperationUpdate,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := RegistrationTokenOperation(input)
	return &out, nil
}

type SSOSecretType string

const (
	SSOSecretTypeCertificate           SSOSecretType = "Certificate"
	SSOSecretTypeCertificateInKeyVault SSOSecretType = "CertificateInKeyVault"
	SSOSecretTypeSharedKey             SSOSecretType = "SharedKey"
	SSOSecretTypeSharedKeyInKeyVault   SSOSecretType = "SharedKeyInKeyVault"
)

func PossibleValuesForSSOSecretType() []string {
	return []string{
		string(SSOSecretTypeCertificate),
		string(SSOSecretTypeCertificateInKeyVault),
		string(SSOSecretTypeSharedKey),
		string(SSOSecretTypeSharedKeyInKeyVault),
	}
}

func parseSSOSecretType(input string) (*SSOSecretType, error) {
	vals := map[string]SSOSecretType{
		"certificate":           SSOSecretTypeCertificate,
		"certificateinkeyvault": SSOSecretTypeCertificateInKeyVault,
		"sharedkey":             SSOSecretTypeSharedKey,
		"sharedkeyinkeyvault":   SSOSecretTypeSharedKeyInKeyVault,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := SSOSecretType(input)
	return &out, nil
}

type SkuTier string

const (
	SkuTierBasic    SkuTier = "Basic"
	SkuTierFree     SkuTier = "Free"
	SkuTierPremium  SkuTier = "Premium"
	SkuTierStandard SkuTier = "Standard"
)

func PossibleValuesForSkuTier() []string {
	return []string{
		string(SkuTierBasic),
		string(SkuTierFree),
		string(SkuTierPremium),
		string(SkuTierStandard),
	}
}

func parseSkuTier(input string) (*SkuTier, error) {
	vals := map[string]SkuTier{
		"basic":    SkuTierBasic,
		"free":     SkuTierFree,
		"premium":  SkuTierPremium,
		"standard": SkuTierStandard,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := SkuTier(input)
	return &out, nil
}
