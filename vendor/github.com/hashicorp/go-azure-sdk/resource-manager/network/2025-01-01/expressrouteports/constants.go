package expressrouteports

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ExpressRouteLinkAdminState string

const (
	ExpressRouteLinkAdminStateDisabled ExpressRouteLinkAdminState = "Disabled"
	ExpressRouteLinkAdminStateEnabled  ExpressRouteLinkAdminState = "Enabled"
)

func PossibleValuesForExpressRouteLinkAdminState() []string {
	return []string{
		string(ExpressRouteLinkAdminStateDisabled),
		string(ExpressRouteLinkAdminStateEnabled),
	}
}

func (s *ExpressRouteLinkAdminState) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseExpressRouteLinkAdminState(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseExpressRouteLinkAdminState(input string) (*ExpressRouteLinkAdminState, error) {
	vals := map[string]ExpressRouteLinkAdminState{
		"disabled": ExpressRouteLinkAdminStateDisabled,
		"enabled":  ExpressRouteLinkAdminStateEnabled,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := ExpressRouteLinkAdminState(input)
	return &out, nil
}

type ExpressRouteLinkConnectorType string

const (
	ExpressRouteLinkConnectorTypeLC ExpressRouteLinkConnectorType = "LC"
	ExpressRouteLinkConnectorTypeSC ExpressRouteLinkConnectorType = "SC"
)

func PossibleValuesForExpressRouteLinkConnectorType() []string {
	return []string{
		string(ExpressRouteLinkConnectorTypeLC),
		string(ExpressRouteLinkConnectorTypeSC),
	}
}

func (s *ExpressRouteLinkConnectorType) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseExpressRouteLinkConnectorType(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseExpressRouteLinkConnectorType(input string) (*ExpressRouteLinkConnectorType, error) {
	vals := map[string]ExpressRouteLinkConnectorType{
		"lc": ExpressRouteLinkConnectorTypeLC,
		"sc": ExpressRouteLinkConnectorTypeSC,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := ExpressRouteLinkConnectorType(input)
	return &out, nil
}

type ExpressRouteLinkMacSecCipher string

const (
	ExpressRouteLinkMacSecCipherGcmAesOneTwoEight    ExpressRouteLinkMacSecCipher = "GcmAes128"
	ExpressRouteLinkMacSecCipherGcmAesTwoFiveSix     ExpressRouteLinkMacSecCipher = "GcmAes256"
	ExpressRouteLinkMacSecCipherGcmAesXpnOneTwoEight ExpressRouteLinkMacSecCipher = "GcmAesXpn128"
	ExpressRouteLinkMacSecCipherGcmAesXpnTwoFiveSix  ExpressRouteLinkMacSecCipher = "GcmAesXpn256"
)

func PossibleValuesForExpressRouteLinkMacSecCipher() []string {
	return []string{
		string(ExpressRouteLinkMacSecCipherGcmAesOneTwoEight),
		string(ExpressRouteLinkMacSecCipherGcmAesTwoFiveSix),
		string(ExpressRouteLinkMacSecCipherGcmAesXpnOneTwoEight),
		string(ExpressRouteLinkMacSecCipherGcmAesXpnTwoFiveSix),
	}
}

func (s *ExpressRouteLinkMacSecCipher) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseExpressRouteLinkMacSecCipher(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseExpressRouteLinkMacSecCipher(input string) (*ExpressRouteLinkMacSecCipher, error) {
	vals := map[string]ExpressRouteLinkMacSecCipher{
		"gcmaes128":    ExpressRouteLinkMacSecCipherGcmAesOneTwoEight,
		"gcmaes256":    ExpressRouteLinkMacSecCipherGcmAesTwoFiveSix,
		"gcmaesxpn128": ExpressRouteLinkMacSecCipherGcmAesXpnOneTwoEight,
		"gcmaesxpn256": ExpressRouteLinkMacSecCipherGcmAesXpnTwoFiveSix,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := ExpressRouteLinkMacSecCipher(input)
	return &out, nil
}

type ExpressRouteLinkMacSecSciState string

const (
	ExpressRouteLinkMacSecSciStateDisabled ExpressRouteLinkMacSecSciState = "Disabled"
	ExpressRouteLinkMacSecSciStateEnabled  ExpressRouteLinkMacSecSciState = "Enabled"
)

func PossibleValuesForExpressRouteLinkMacSecSciState() []string {
	return []string{
		string(ExpressRouteLinkMacSecSciStateDisabled),
		string(ExpressRouteLinkMacSecSciStateEnabled),
	}
}

func (s *ExpressRouteLinkMacSecSciState) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseExpressRouteLinkMacSecSciState(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseExpressRouteLinkMacSecSciState(input string) (*ExpressRouteLinkMacSecSciState, error) {
	vals := map[string]ExpressRouteLinkMacSecSciState{
		"disabled": ExpressRouteLinkMacSecSciStateDisabled,
		"enabled":  ExpressRouteLinkMacSecSciStateEnabled,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := ExpressRouteLinkMacSecSciState(input)
	return &out, nil
}

type ExpressRoutePortsBillingType string

const (
	ExpressRoutePortsBillingTypeMeteredData   ExpressRoutePortsBillingType = "MeteredData"
	ExpressRoutePortsBillingTypeUnlimitedData ExpressRoutePortsBillingType = "UnlimitedData"
)

func PossibleValuesForExpressRoutePortsBillingType() []string {
	return []string{
		string(ExpressRoutePortsBillingTypeMeteredData),
		string(ExpressRoutePortsBillingTypeUnlimitedData),
	}
}

func (s *ExpressRoutePortsBillingType) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseExpressRoutePortsBillingType(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseExpressRoutePortsBillingType(input string) (*ExpressRoutePortsBillingType, error) {
	vals := map[string]ExpressRoutePortsBillingType{
		"metereddata":   ExpressRoutePortsBillingTypeMeteredData,
		"unlimiteddata": ExpressRoutePortsBillingTypeUnlimitedData,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := ExpressRoutePortsBillingType(input)
	return &out, nil
}

type ExpressRoutePortsEncapsulation string

const (
	ExpressRoutePortsEncapsulationDotOneQ ExpressRoutePortsEncapsulation = "Dot1Q"
	ExpressRoutePortsEncapsulationQinQ    ExpressRoutePortsEncapsulation = "QinQ"
)

func PossibleValuesForExpressRoutePortsEncapsulation() []string {
	return []string{
		string(ExpressRoutePortsEncapsulationDotOneQ),
		string(ExpressRoutePortsEncapsulationQinQ),
	}
}

func (s *ExpressRoutePortsEncapsulation) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseExpressRoutePortsEncapsulation(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseExpressRoutePortsEncapsulation(input string) (*ExpressRoutePortsEncapsulation, error) {
	vals := map[string]ExpressRoutePortsEncapsulation{
		"dot1q": ExpressRoutePortsEncapsulationDotOneQ,
		"qinq":  ExpressRoutePortsEncapsulationQinQ,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := ExpressRoutePortsEncapsulation(input)
	return &out, nil
}

type ProvisioningState string

const (
	ProvisioningStateDeleting  ProvisioningState = "Deleting"
	ProvisioningStateFailed    ProvisioningState = "Failed"
	ProvisioningStateSucceeded ProvisioningState = "Succeeded"
	ProvisioningStateUpdating  ProvisioningState = "Updating"
)

func PossibleValuesForProvisioningState() []string {
	return []string{
		string(ProvisioningStateDeleting),
		string(ProvisioningStateFailed),
		string(ProvisioningStateSucceeded),
		string(ProvisioningStateUpdating),
	}
}

func (s *ProvisioningState) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseProvisioningState(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseProvisioningState(input string) (*ProvisioningState, error) {
	vals := map[string]ProvisioningState{
		"deleting":  ProvisioningStateDeleting,
		"failed":    ProvisioningStateFailed,
		"succeeded": ProvisioningStateSucceeded,
		"updating":  ProvisioningStateUpdating,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := ProvisioningState(input)
	return &out, nil
}
