package gallerysharingupdate

import "strings"

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type SharingProfileGroupTypes string

const (
	SharingProfileGroupTypesAADTenants    SharingProfileGroupTypes = "AADTenants"
	SharingProfileGroupTypesSubscriptions SharingProfileGroupTypes = "Subscriptions"
)

func PossibleValuesForSharingProfileGroupTypes() []string {
	return []string{
		string(SharingProfileGroupTypesAADTenants),
		string(SharingProfileGroupTypesSubscriptions),
	}
}

func parseSharingProfileGroupTypes(input string) (*SharingProfileGroupTypes, error) {
	vals := map[string]SharingProfileGroupTypes{
		"aadtenants":    SharingProfileGroupTypesAADTenants,
		"subscriptions": SharingProfileGroupTypesSubscriptions,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := SharingProfileGroupTypes(input)
	return &out, nil
}

type SharingUpdateOperationTypes string

const (
	SharingUpdateOperationTypesAdd             SharingUpdateOperationTypes = "Add"
	SharingUpdateOperationTypesEnableCommunity SharingUpdateOperationTypes = "EnableCommunity"
	SharingUpdateOperationTypesRemove          SharingUpdateOperationTypes = "Remove"
	SharingUpdateOperationTypesReset           SharingUpdateOperationTypes = "Reset"
)

func PossibleValuesForSharingUpdateOperationTypes() []string {
	return []string{
		string(SharingUpdateOperationTypesAdd),
		string(SharingUpdateOperationTypesEnableCommunity),
		string(SharingUpdateOperationTypesRemove),
		string(SharingUpdateOperationTypesReset),
	}
}

func parseSharingUpdateOperationTypes(input string) (*SharingUpdateOperationTypes, error) {
	vals := map[string]SharingUpdateOperationTypes{
		"add":             SharingUpdateOperationTypesAdd,
		"enablecommunity": SharingUpdateOperationTypesEnableCommunity,
		"remove":          SharingUpdateOperationTypesRemove,
		"reset":           SharingUpdateOperationTypesReset,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := SharingUpdateOperationTypes(input)
	return &out, nil
}
