package broker

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type BrokerMemoryProfile string

const (
	BrokerMemoryProfileHigh   BrokerMemoryProfile = "High"
	BrokerMemoryProfileLow    BrokerMemoryProfile = "Low"
	BrokerMemoryProfileMedium BrokerMemoryProfile = "Medium"
	BrokerMemoryProfileTiny   BrokerMemoryProfile = "Tiny"
)

func PossibleValuesForBrokerMemoryProfile() []string {
	return []string{
		string(BrokerMemoryProfileHigh),
		string(BrokerMemoryProfileLow),
		string(BrokerMemoryProfileMedium),
		string(BrokerMemoryProfileTiny),
	}
}

func (s *BrokerMemoryProfile) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseBrokerMemoryProfile(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseBrokerMemoryProfile(input string) (*BrokerMemoryProfile, error) {
	vals := map[string]BrokerMemoryProfile{
		"high":   BrokerMemoryProfileHigh,
		"low":    BrokerMemoryProfileLow,
		"medium": BrokerMemoryProfileMedium,
		"tiny":   BrokerMemoryProfileTiny,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := BrokerMemoryProfile(input)
	return &out, nil
}

type ExtendedLocationType string

const (
	ExtendedLocationTypeCustomLocation ExtendedLocationType = "CustomLocation"
)

func PossibleValuesForExtendedLocationType() []string {
	return []string{
		string(ExtendedLocationTypeCustomLocation),
	}
}

func (s *ExtendedLocationType) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseExtendedLocationType(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseExtendedLocationType(input string) (*ExtendedLocationType, error) {
	vals := map[string]ExtendedLocationType{
		"customlocation": ExtendedLocationTypeCustomLocation,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := ExtendedLocationType(input)
	return &out, nil
}

type OperationalMode string

const (
	OperationalModeDisabled OperationalMode = "Disabled"
	OperationalModeEnabled  OperationalMode = "Enabled"
)

func PossibleValuesForOperationalMode() []string {
	return []string{
		string(OperationalModeDisabled),
		string(OperationalModeEnabled),
	}
}

func (s *OperationalMode) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseOperationalMode(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseOperationalMode(input string) (*OperationalMode, error) {
	vals := map[string]OperationalMode{
		"disabled": OperationalModeDisabled,
		"enabled":  OperationalModeEnabled,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := OperationalMode(input)
	return &out, nil
}

type OperatorValues string

const (
	OperatorValuesDoesNotExist OperatorValues = "DoesNotExist"
	OperatorValuesExists       OperatorValues = "Exists"
	OperatorValuesIn           OperatorValues = "In"
	OperatorValuesNotIn        OperatorValues = "NotIn"
)

func PossibleValuesForOperatorValues() []string {
	return []string{
		string(OperatorValuesDoesNotExist),
		string(OperatorValuesExists),
		string(OperatorValuesIn),
		string(OperatorValuesNotIn),
	}
}

func (s *OperatorValues) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseOperatorValues(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseOperatorValues(input string) (*OperatorValues, error) {
	vals := map[string]OperatorValues{
		"doesnotexist": OperatorValuesDoesNotExist,
		"exists":       OperatorValuesExists,
		"in":           OperatorValuesIn,
		"notin":        OperatorValuesNotIn,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := OperatorValues(input)
	return &out, nil
}

type PrivateKeyAlgorithm string

const (
	PrivateKeyAlgorithmEcFiveTwoOne         PrivateKeyAlgorithm = "Ec521"
	PrivateKeyAlgorithmEcThreeEightFour     PrivateKeyAlgorithm = "Ec384"
	PrivateKeyAlgorithmEcTwoFiveSix         PrivateKeyAlgorithm = "Ec256"
	PrivateKeyAlgorithmEdTwoFiveFiveOneNine PrivateKeyAlgorithm = "Ed25519"
	PrivateKeyAlgorithmRsaEightOneNineTwo   PrivateKeyAlgorithm = "Rsa8192"
	PrivateKeyAlgorithmRsaFourZeroNineSix   PrivateKeyAlgorithm = "Rsa4096"
	PrivateKeyAlgorithmRsaTwoZeroFourEight  PrivateKeyAlgorithm = "Rsa2048"
)

func PossibleValuesForPrivateKeyAlgorithm() []string {
	return []string{
		string(PrivateKeyAlgorithmEcFiveTwoOne),
		string(PrivateKeyAlgorithmEcThreeEightFour),
		string(PrivateKeyAlgorithmEcTwoFiveSix),
		string(PrivateKeyAlgorithmEdTwoFiveFiveOneNine),
		string(PrivateKeyAlgorithmRsaEightOneNineTwo),
		string(PrivateKeyAlgorithmRsaFourZeroNineSix),
		string(PrivateKeyAlgorithmRsaTwoZeroFourEight),
	}
}

func (s *PrivateKeyAlgorithm) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parsePrivateKeyAlgorithm(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parsePrivateKeyAlgorithm(input string) (*PrivateKeyAlgorithm, error) {
	vals := map[string]PrivateKeyAlgorithm{
		"ec521":   PrivateKeyAlgorithmEcFiveTwoOne,
		"ec384":   PrivateKeyAlgorithmEcThreeEightFour,
		"ec256":   PrivateKeyAlgorithmEcTwoFiveSix,
		"ed25519": PrivateKeyAlgorithmEdTwoFiveFiveOneNine,
		"rsa8192": PrivateKeyAlgorithmRsaEightOneNineTwo,
		"rsa4096": PrivateKeyAlgorithmRsaFourZeroNineSix,
		"rsa2048": PrivateKeyAlgorithmRsaTwoZeroFourEight,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := PrivateKeyAlgorithm(input)
	return &out, nil
}

type PrivateKeyRotationPolicy string

const (
	PrivateKeyRotationPolicyAlways PrivateKeyRotationPolicy = "Always"
	PrivateKeyRotationPolicyNever  PrivateKeyRotationPolicy = "Never"
)

func PossibleValuesForPrivateKeyRotationPolicy() []string {
	return []string{
		string(PrivateKeyRotationPolicyAlways),
		string(PrivateKeyRotationPolicyNever),
	}
}

func (s *PrivateKeyRotationPolicy) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parsePrivateKeyRotationPolicy(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parsePrivateKeyRotationPolicy(input string) (*PrivateKeyRotationPolicy, error) {
	vals := map[string]PrivateKeyRotationPolicy{
		"always": PrivateKeyRotationPolicyAlways,
		"never":  PrivateKeyRotationPolicyNever,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := PrivateKeyRotationPolicy(input)
	return &out, nil
}

type ProvisioningState string

const (
	ProvisioningStateAccepted     ProvisioningState = "Accepted"
	ProvisioningStateCanceled     ProvisioningState = "Canceled"
	ProvisioningStateDeleting     ProvisioningState = "Deleting"
	ProvisioningStateFailed       ProvisioningState = "Failed"
	ProvisioningStateProvisioning ProvisioningState = "Provisioning"
	ProvisioningStateSucceeded    ProvisioningState = "Succeeded"
	ProvisioningStateUpdating     ProvisioningState = "Updating"
)

func PossibleValuesForProvisioningState() []string {
	return []string{
		string(ProvisioningStateAccepted),
		string(ProvisioningStateCanceled),
		string(ProvisioningStateDeleting),
		string(ProvisioningStateFailed),
		string(ProvisioningStateProvisioning),
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
		"accepted":     ProvisioningStateAccepted,
		"canceled":     ProvisioningStateCanceled,
		"deleting":     ProvisioningStateDeleting,
		"failed":       ProvisioningStateFailed,
		"provisioning": ProvisioningStateProvisioning,
		"succeeded":    ProvisioningStateSucceeded,
		"updating":     ProvisioningStateUpdating,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := ProvisioningState(input)
	return &out, nil
}

type SubscriberMessageDropStrategy string

const (
	SubscriberMessageDropStrategyDropOldest SubscriberMessageDropStrategy = "DropOldest"
	SubscriberMessageDropStrategyNone       SubscriberMessageDropStrategy = "None"
)

func PossibleValuesForSubscriberMessageDropStrategy() []string {
	return []string{
		string(SubscriberMessageDropStrategyDropOldest),
		string(SubscriberMessageDropStrategyNone),
	}
}

func (s *SubscriberMessageDropStrategy) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseSubscriberMessageDropStrategy(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseSubscriberMessageDropStrategy(input string) (*SubscriberMessageDropStrategy, error) {
	vals := map[string]SubscriberMessageDropStrategy{
		"dropoldest": SubscriberMessageDropStrategyDropOldest,
		"none":       SubscriberMessageDropStrategyNone,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := SubscriberMessageDropStrategy(input)
	return &out, nil
}
