package availableservicetiers

import "strings"

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type SkuNameEnum string

const (
	SkuNameEnumCapacityReservation  SkuNameEnum = "CapacityReservation"
	SkuNameEnumFree                 SkuNameEnum = "Free"
	SkuNameEnumPerGBTwoZeroOneEight SkuNameEnum = "PerGB2018"
	SkuNameEnumPerNode              SkuNameEnum = "PerNode"
	SkuNameEnumPremium              SkuNameEnum = "Premium"
	SkuNameEnumStandalone           SkuNameEnum = "Standalone"
	SkuNameEnumStandard             SkuNameEnum = "Standard"
)

func PossibleValuesForSkuNameEnum() []string {
	return []string{
		string(SkuNameEnumCapacityReservation),
		string(SkuNameEnumFree),
		string(SkuNameEnumPerGBTwoZeroOneEight),
		string(SkuNameEnumPerNode),
		string(SkuNameEnumPremium),
		string(SkuNameEnumStandalone),
		string(SkuNameEnumStandard),
	}
}

func parseSkuNameEnum(input string) (*SkuNameEnum, error) {
	vals := map[string]SkuNameEnum{
		"capacityreservation": SkuNameEnumCapacityReservation,
		"free":                SkuNameEnumFree,
		"pergb2018":           SkuNameEnumPerGBTwoZeroOneEight,
		"pernode":             SkuNameEnumPerNode,
		"premium":             SkuNameEnumPremium,
		"standalone":          SkuNameEnumStandalone,
		"standard":            SkuNameEnumStandard,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := SkuNameEnum(input)
	return &out, nil
}
