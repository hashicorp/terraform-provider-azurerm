package skus

import "strings"

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ExtendedLocationType string

const (
	ExtendedLocationTypeEdgeZone ExtendedLocationType = "EdgeZone"
)

func PossibleValuesForExtendedLocationType() []string {
	return []string{
		string(ExtendedLocationTypeEdgeZone),
	}
}

func parseExtendedLocationType(input string) (*ExtendedLocationType, error) {
	vals := map[string]ExtendedLocationType{
		"edgezone": ExtendedLocationTypeEdgeZone,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := ExtendedLocationType(input)
	return &out, nil
}

type ResourceSkuCapacityScaleType string

const (
	ResourceSkuCapacityScaleTypeAutomatic ResourceSkuCapacityScaleType = "Automatic"
	ResourceSkuCapacityScaleTypeManual    ResourceSkuCapacityScaleType = "Manual"
	ResourceSkuCapacityScaleTypeNone      ResourceSkuCapacityScaleType = "None"
)

func PossibleValuesForResourceSkuCapacityScaleType() []string {
	return []string{
		string(ResourceSkuCapacityScaleTypeAutomatic),
		string(ResourceSkuCapacityScaleTypeManual),
		string(ResourceSkuCapacityScaleTypeNone),
	}
}

func parseResourceSkuCapacityScaleType(input string) (*ResourceSkuCapacityScaleType, error) {
	vals := map[string]ResourceSkuCapacityScaleType{
		"automatic": ResourceSkuCapacityScaleTypeAutomatic,
		"manual":    ResourceSkuCapacityScaleTypeManual,
		"none":      ResourceSkuCapacityScaleTypeNone,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := ResourceSkuCapacityScaleType(input)
	return &out, nil
}

type ResourceSkuRestrictionsReasonCode string

const (
	ResourceSkuRestrictionsReasonCodeNotAvailableForSubscription ResourceSkuRestrictionsReasonCode = "NotAvailableForSubscription"
	ResourceSkuRestrictionsReasonCodeQuotaId                     ResourceSkuRestrictionsReasonCode = "QuotaId"
)

func PossibleValuesForResourceSkuRestrictionsReasonCode() []string {
	return []string{
		string(ResourceSkuRestrictionsReasonCodeNotAvailableForSubscription),
		string(ResourceSkuRestrictionsReasonCodeQuotaId),
	}
}

func parseResourceSkuRestrictionsReasonCode(input string) (*ResourceSkuRestrictionsReasonCode, error) {
	vals := map[string]ResourceSkuRestrictionsReasonCode{
		"notavailableforsubscription": ResourceSkuRestrictionsReasonCodeNotAvailableForSubscription,
		"quotaid":                     ResourceSkuRestrictionsReasonCodeQuotaId,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := ResourceSkuRestrictionsReasonCode(input)
	return &out, nil
}

type ResourceSkuRestrictionsType string

const (
	ResourceSkuRestrictionsTypeLocation ResourceSkuRestrictionsType = "Location"
	ResourceSkuRestrictionsTypeZone     ResourceSkuRestrictionsType = "Zone"
)

func PossibleValuesForResourceSkuRestrictionsType() []string {
	return []string{
		string(ResourceSkuRestrictionsTypeLocation),
		string(ResourceSkuRestrictionsTypeZone),
	}
}

func parseResourceSkuRestrictionsType(input string) (*ResourceSkuRestrictionsType, error) {
	vals := map[string]ResourceSkuRestrictionsType{
		"location": ResourceSkuRestrictionsTypeLocation,
		"zone":     ResourceSkuRestrictionsTypeZone,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := ResourceSkuRestrictionsType(input)
	return &out, nil
}
