package virtualmachinesizes

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type BillingCurrency string

const (
	BillingCurrencyUSD BillingCurrency = "USD"
)

func PossibleValuesForBillingCurrency() []string {
	return []string{
		string(BillingCurrencyUSD),
	}
}

func (s *BillingCurrency) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseBillingCurrency(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseBillingCurrency(input string) (*BillingCurrency, error) {
	vals := map[string]BillingCurrency{
		"usd": BillingCurrencyUSD,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := BillingCurrency(input)
	return &out, nil
}

type UnitOfMeasure string

const (
	UnitOfMeasureOneHour UnitOfMeasure = "OneHour"
)

func PossibleValuesForUnitOfMeasure() []string {
	return []string{
		string(UnitOfMeasureOneHour),
	}
}

func (s *UnitOfMeasure) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseUnitOfMeasure(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseUnitOfMeasure(input string) (*UnitOfMeasure, error) {
	vals := map[string]UnitOfMeasure{
		"onehour": UnitOfMeasureOneHour,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := UnitOfMeasure(input)
	return &out, nil
}

type VMPriceOSType string

const (
	VMPriceOSTypeLinux   VMPriceOSType = "Linux"
	VMPriceOSTypeWindows VMPriceOSType = "Windows"
)

func PossibleValuesForVMPriceOSType() []string {
	return []string{
		string(VMPriceOSTypeLinux),
		string(VMPriceOSTypeWindows),
	}
}

func (s *VMPriceOSType) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseVMPriceOSType(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseVMPriceOSType(input string) (*VMPriceOSType, error) {
	vals := map[string]VMPriceOSType{
		"linux":   VMPriceOSTypeLinux,
		"windows": VMPriceOSTypeWindows,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := VMPriceOSType(input)
	return &out, nil
}

type VMTier string

const (
	VMTierLowPriority VMTier = "LowPriority"
	VMTierSpot        VMTier = "Spot"
	VMTierStandard    VMTier = "Standard"
)

func PossibleValuesForVMTier() []string {
	return []string{
		string(VMTierLowPriority),
		string(VMTierSpot),
		string(VMTierStandard),
	}
}

func (s *VMTier) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseVMTier(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseVMTier(input string) (*VMTier, error) {
	vals := map[string]VMTier{
		"lowpriority": VMTierLowPriority,
		"spot":        VMTierSpot,
		"standard":    VMTierStandard,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := VMTier(input)
	return &out, nil
}
