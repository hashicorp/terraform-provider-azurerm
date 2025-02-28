package skus

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ApiManagementSkuCapacityScaleType string

const (
	ApiManagementSkuCapacityScaleTypeAutomatic ApiManagementSkuCapacityScaleType = "Automatic"
	ApiManagementSkuCapacityScaleTypeManual    ApiManagementSkuCapacityScaleType = "Manual"
	ApiManagementSkuCapacityScaleTypeNone      ApiManagementSkuCapacityScaleType = "None"
)

func PossibleValuesForApiManagementSkuCapacityScaleType() []string {
	return []string{
		string(ApiManagementSkuCapacityScaleTypeAutomatic),
		string(ApiManagementSkuCapacityScaleTypeManual),
		string(ApiManagementSkuCapacityScaleTypeNone),
	}
}

func (s *ApiManagementSkuCapacityScaleType) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseApiManagementSkuCapacityScaleType(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseApiManagementSkuCapacityScaleType(input string) (*ApiManagementSkuCapacityScaleType, error) {
	vals := map[string]ApiManagementSkuCapacityScaleType{
		"automatic": ApiManagementSkuCapacityScaleTypeAutomatic,
		"manual":    ApiManagementSkuCapacityScaleTypeManual,
		"none":      ApiManagementSkuCapacityScaleTypeNone,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := ApiManagementSkuCapacityScaleType(input)
	return &out, nil
}

type ApiManagementSkuRestrictionsReasonCode string

const (
	ApiManagementSkuRestrictionsReasonCodeNotAvailableForSubscription ApiManagementSkuRestrictionsReasonCode = "NotAvailableForSubscription"
	ApiManagementSkuRestrictionsReasonCodeQuotaId                     ApiManagementSkuRestrictionsReasonCode = "QuotaId"
)

func PossibleValuesForApiManagementSkuRestrictionsReasonCode() []string {
	return []string{
		string(ApiManagementSkuRestrictionsReasonCodeNotAvailableForSubscription),
		string(ApiManagementSkuRestrictionsReasonCodeQuotaId),
	}
}

func (s *ApiManagementSkuRestrictionsReasonCode) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseApiManagementSkuRestrictionsReasonCode(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseApiManagementSkuRestrictionsReasonCode(input string) (*ApiManagementSkuRestrictionsReasonCode, error) {
	vals := map[string]ApiManagementSkuRestrictionsReasonCode{
		"notavailableforsubscription": ApiManagementSkuRestrictionsReasonCodeNotAvailableForSubscription,
		"quotaid":                     ApiManagementSkuRestrictionsReasonCodeQuotaId,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := ApiManagementSkuRestrictionsReasonCode(input)
	return &out, nil
}

type ApiManagementSkuRestrictionsType string

const (
	ApiManagementSkuRestrictionsTypeLocation ApiManagementSkuRestrictionsType = "Location"
	ApiManagementSkuRestrictionsTypeZone     ApiManagementSkuRestrictionsType = "Zone"
)

func PossibleValuesForApiManagementSkuRestrictionsType() []string {
	return []string{
		string(ApiManagementSkuRestrictionsTypeLocation),
		string(ApiManagementSkuRestrictionsTypeZone),
	}
}

func (s *ApiManagementSkuRestrictionsType) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseApiManagementSkuRestrictionsType(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseApiManagementSkuRestrictionsType(input string) (*ApiManagementSkuRestrictionsType, error) {
	vals := map[string]ApiManagementSkuRestrictionsType{
		"location": ApiManagementSkuRestrictionsTypeLocation,
		"zone":     ApiManagementSkuRestrictionsTypeZone,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := ApiManagementSkuRestrictionsType(input)
	return &out, nil
}
