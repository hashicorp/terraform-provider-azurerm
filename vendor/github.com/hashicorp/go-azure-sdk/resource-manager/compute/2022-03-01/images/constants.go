package images

import "strings"

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type CachingTypes string

const (
	CachingTypesNone      CachingTypes = "None"
	CachingTypesReadOnly  CachingTypes = "ReadOnly"
	CachingTypesReadWrite CachingTypes = "ReadWrite"
)

func PossibleValuesForCachingTypes() []string {
	return []string{
		string(CachingTypesNone),
		string(CachingTypesReadOnly),
		string(CachingTypesReadWrite),
	}
}

func parseCachingTypes(input string) (*CachingTypes, error) {
	vals := map[string]CachingTypes{
		"none":      CachingTypesNone,
		"readonly":  CachingTypesReadOnly,
		"readwrite": CachingTypesReadWrite,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := CachingTypes(input)
	return &out, nil
}

type HyperVGenerationTypes string

const (
	HyperVGenerationTypesVOne HyperVGenerationTypes = "V1"
	HyperVGenerationTypesVTwo HyperVGenerationTypes = "V2"
)

func PossibleValuesForHyperVGenerationTypes() []string {
	return []string{
		string(HyperVGenerationTypesVOne),
		string(HyperVGenerationTypesVTwo),
	}
}

func parseHyperVGenerationTypes(input string) (*HyperVGenerationTypes, error) {
	vals := map[string]HyperVGenerationTypes{
		"v1": HyperVGenerationTypesVOne,
		"v2": HyperVGenerationTypesVTwo,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := HyperVGenerationTypes(input)
	return &out, nil
}

type OperatingSystemStateTypes string

const (
	OperatingSystemStateTypesGeneralized OperatingSystemStateTypes = "Generalized"
	OperatingSystemStateTypesSpecialized OperatingSystemStateTypes = "Specialized"
)

func PossibleValuesForOperatingSystemStateTypes() []string {
	return []string{
		string(OperatingSystemStateTypesGeneralized),
		string(OperatingSystemStateTypesSpecialized),
	}
}

func parseOperatingSystemStateTypes(input string) (*OperatingSystemStateTypes, error) {
	vals := map[string]OperatingSystemStateTypes{
		"generalized": OperatingSystemStateTypesGeneralized,
		"specialized": OperatingSystemStateTypesSpecialized,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := OperatingSystemStateTypes(input)
	return &out, nil
}

type OperatingSystemTypes string

const (
	OperatingSystemTypesLinux   OperatingSystemTypes = "Linux"
	OperatingSystemTypesWindows OperatingSystemTypes = "Windows"
)

func PossibleValuesForOperatingSystemTypes() []string {
	return []string{
		string(OperatingSystemTypesLinux),
		string(OperatingSystemTypesWindows),
	}
}

func parseOperatingSystemTypes(input string) (*OperatingSystemTypes, error) {
	vals := map[string]OperatingSystemTypes{
		"linux":   OperatingSystemTypesLinux,
		"windows": OperatingSystemTypesWindows,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := OperatingSystemTypes(input)
	return &out, nil
}

type StorageAccountTypes string

const (
	StorageAccountTypesPremiumLRS     StorageAccountTypes = "Premium_LRS"
	StorageAccountTypesPremiumVTwoLRS StorageAccountTypes = "PremiumV2_LRS"
	StorageAccountTypesPremiumZRS     StorageAccountTypes = "Premium_ZRS"
	StorageAccountTypesStandardLRS    StorageAccountTypes = "Standard_LRS"
	StorageAccountTypesStandardSSDLRS StorageAccountTypes = "StandardSSD_LRS"
	StorageAccountTypesStandardSSDZRS StorageAccountTypes = "StandardSSD_ZRS"
	StorageAccountTypesUltraSSDLRS    StorageAccountTypes = "UltraSSD_LRS"
)

func PossibleValuesForStorageAccountTypes() []string {
	return []string{
		string(StorageAccountTypesPremiumLRS),
		string(StorageAccountTypesPremiumVTwoLRS),
		string(StorageAccountTypesPremiumZRS),
		string(StorageAccountTypesStandardLRS),
		string(StorageAccountTypesStandardSSDLRS),
		string(StorageAccountTypesStandardSSDZRS),
		string(StorageAccountTypesUltraSSDLRS),
	}
}

func parseStorageAccountTypes(input string) (*StorageAccountTypes, error) {
	vals := map[string]StorageAccountTypes{
		"premium_lrs":     StorageAccountTypesPremiumLRS,
		"premiumv2_lrs":   StorageAccountTypesPremiumVTwoLRS,
		"premium_zrs":     StorageAccountTypesPremiumZRS,
		"standard_lrs":    StorageAccountTypesStandardLRS,
		"standardssd_lrs": StorageAccountTypesStandardSSDLRS,
		"standardssd_zrs": StorageAccountTypesStandardSSDZRS,
		"ultrassd_lrs":    StorageAccountTypesUltraSSDLRS,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := StorageAccountTypes(input)
	return &out, nil
}
