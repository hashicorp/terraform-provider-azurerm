package fluidrelayservers

import "strings"

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type CmkIdentityType string

const (
	CmkIdentityTypeSystemAssigned CmkIdentityType = "SystemAssigned"
	CmkIdentityTypeUserAssigned   CmkIdentityType = "UserAssigned"
)

func PossibleValuesForCmkIdentityType() []string {
	return []string{
		string(CmkIdentityTypeSystemAssigned),
		string(CmkIdentityTypeUserAssigned),
	}
}

func parseCmkIdentityType(input string) (*CmkIdentityType, error) {
	vals := map[string]CmkIdentityType{
		"systemassigned": CmkIdentityTypeSystemAssigned,
		"userassigned":   CmkIdentityTypeUserAssigned,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := CmkIdentityType(input)
	return &out, nil
}

type KeyName string

const (
	KeyNameKeyOne KeyName = "key1"
	KeyNameKeyTwo KeyName = "key2"
)

func PossibleValuesForKeyName() []string {
	return []string{
		string(KeyNameKeyOne),
		string(KeyNameKeyTwo),
	}
}

func parseKeyName(input string) (*KeyName, error) {
	vals := map[string]KeyName{
		"key1": KeyNameKeyOne,
		"key2": KeyNameKeyTwo,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := KeyName(input)
	return &out, nil
}

type ProvisioningState string

const (
	ProvisioningStateCanceled  ProvisioningState = "Canceled"
	ProvisioningStateFailed    ProvisioningState = "Failed"
	ProvisioningStateSucceeded ProvisioningState = "Succeeded"
)

func PossibleValuesForProvisioningState() []string {
	return []string{
		string(ProvisioningStateCanceled),
		string(ProvisioningStateFailed),
		string(ProvisioningStateSucceeded),
	}
}

func parseProvisioningState(input string) (*ProvisioningState, error) {
	vals := map[string]ProvisioningState{
		"canceled":  ProvisioningStateCanceled,
		"failed":    ProvisioningStateFailed,
		"succeeded": ProvisioningStateSucceeded,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := ProvisioningState(input)
	return &out, nil
}

type StorageSKU string

const (
	StorageSKUBasic    StorageSKU = "basic"
	StorageSKUStandard StorageSKU = "standard"
)

func PossibleValuesForStorageSKU() []string {
	return []string{
		string(StorageSKUBasic),
		string(StorageSKUStandard),
	}
}

func parseStorageSKU(input string) (*StorageSKU, error) {
	vals := map[string]StorageSKU{
		"basic":    StorageSKUBasic,
		"standard": StorageSKUStandard,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := StorageSKU(input)
	return &out, nil
}
