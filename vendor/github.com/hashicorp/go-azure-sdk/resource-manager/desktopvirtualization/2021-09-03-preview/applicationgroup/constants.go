package applicationgroup

import "strings"

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ApplicationGroupType string

const (
	ApplicationGroupTypeDesktop   ApplicationGroupType = "Desktop"
	ApplicationGroupTypeRemoteApp ApplicationGroupType = "RemoteApp"
)

func PossibleValuesForApplicationGroupType() []string {
	return []string{
		string(ApplicationGroupTypeDesktop),
		string(ApplicationGroupTypeRemoteApp),
	}
}

func parseApplicationGroupType(input string) (*ApplicationGroupType, error) {
	vals := map[string]ApplicationGroupType{
		"desktop":   ApplicationGroupTypeDesktop,
		"remoteapp": ApplicationGroupTypeRemoteApp,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := ApplicationGroupType(input)
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
