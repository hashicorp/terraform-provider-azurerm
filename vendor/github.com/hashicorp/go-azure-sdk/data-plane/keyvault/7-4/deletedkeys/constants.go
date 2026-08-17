package deletedkeys

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type DeletionRecoveryLevel string

const (
	DeletionRecoveryLevelCustomizedRecoverable                              DeletionRecoveryLevel = "CustomizedRecoverable"
	DeletionRecoveryLevelCustomizedRecoverablePositiveProtectedSubscription DeletionRecoveryLevel = "CustomizedRecoverable+ProtectedSubscription"
	DeletionRecoveryLevelCustomizedRecoverablePositivePurgeable             DeletionRecoveryLevel = "CustomizedRecoverable+Purgeable"
	DeletionRecoveryLevelPurgeable                                          DeletionRecoveryLevel = "Purgeable"
	DeletionRecoveryLevelRecoverable                                        DeletionRecoveryLevel = "Recoverable"
	DeletionRecoveryLevelRecoverablePositiveProtectedSubscription           DeletionRecoveryLevel = "Recoverable+ProtectedSubscription"
	DeletionRecoveryLevelRecoverablePositivePurgeable                       DeletionRecoveryLevel = "Recoverable+Purgeable"
)

func PossibleValuesForDeletionRecoveryLevel() []string {
	return []string{
		string(DeletionRecoveryLevelCustomizedRecoverable),
		string(DeletionRecoveryLevelCustomizedRecoverablePositiveProtectedSubscription),
		string(DeletionRecoveryLevelCustomizedRecoverablePositivePurgeable),
		string(DeletionRecoveryLevelPurgeable),
		string(DeletionRecoveryLevelRecoverable),
		string(DeletionRecoveryLevelRecoverablePositiveProtectedSubscription),
		string(DeletionRecoveryLevelRecoverablePositivePurgeable),
	}
}

func (s *DeletionRecoveryLevel) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseDeletionRecoveryLevel(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseDeletionRecoveryLevel(input string) (*DeletionRecoveryLevel, error) {
	vals := map[string]DeletionRecoveryLevel{
		"customizedrecoverable":                       DeletionRecoveryLevelCustomizedRecoverable,
		"customizedrecoverable+protectedsubscription": DeletionRecoveryLevelCustomizedRecoverablePositiveProtectedSubscription,
		"customizedrecoverable+purgeable":             DeletionRecoveryLevelCustomizedRecoverablePositivePurgeable,
		"purgeable":                                   DeletionRecoveryLevelPurgeable,
		"recoverable":                                 DeletionRecoveryLevelRecoverable,
		"recoverable+protectedsubscription":           DeletionRecoveryLevelRecoverablePositiveProtectedSubscription,
		"recoverable+purgeable":                       DeletionRecoveryLevelRecoverablePositivePurgeable,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := DeletionRecoveryLevel(input)
	return &out, nil
}

type JsonWebKeyCurveName string

const (
	JsonWebKeyCurveNamePNegativeFiveTwoOne     JsonWebKeyCurveName = "P-521"
	JsonWebKeyCurveNamePNegativeThreeEightFour JsonWebKeyCurveName = "P-384"
	JsonWebKeyCurveNamePNegativeTwoFiveSix     JsonWebKeyCurveName = "P-256"
	JsonWebKeyCurveNamePNegativeTwoFiveSixK    JsonWebKeyCurveName = "P-256K"
)

func PossibleValuesForJsonWebKeyCurveName() []string {
	return []string{
		string(JsonWebKeyCurveNamePNegativeFiveTwoOne),
		string(JsonWebKeyCurveNamePNegativeThreeEightFour),
		string(JsonWebKeyCurveNamePNegativeTwoFiveSix),
		string(JsonWebKeyCurveNamePNegativeTwoFiveSixK),
	}
}

func (s *JsonWebKeyCurveName) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseJsonWebKeyCurveName(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseJsonWebKeyCurveName(input string) (*JsonWebKeyCurveName, error) {
	vals := map[string]JsonWebKeyCurveName{
		"p-521":  JsonWebKeyCurveNamePNegativeFiveTwoOne,
		"p-384":  JsonWebKeyCurveNamePNegativeThreeEightFour,
		"p-256":  JsonWebKeyCurveNamePNegativeTwoFiveSix,
		"p-256k": JsonWebKeyCurveNamePNegativeTwoFiveSixK,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := JsonWebKeyCurveName(input)
	return &out, nil
}

type JsonWebKeyType string

const (
	JsonWebKeyTypeEC             JsonWebKeyType = "EC"
	JsonWebKeyTypeECNegativeHSM  JsonWebKeyType = "EC-HSM"
	JsonWebKeyTypeOct            JsonWebKeyType = "oct"
	JsonWebKeyTypeOctNegativeHSM JsonWebKeyType = "oct-HSM"
	JsonWebKeyTypeRSA            JsonWebKeyType = "RSA"
	JsonWebKeyTypeRSANegativeHSM JsonWebKeyType = "RSA-HSM"
)

func PossibleValuesForJsonWebKeyType() []string {
	return []string{
		string(JsonWebKeyTypeEC),
		string(JsonWebKeyTypeECNegativeHSM),
		string(JsonWebKeyTypeOct),
		string(JsonWebKeyTypeOctNegativeHSM),
		string(JsonWebKeyTypeRSA),
		string(JsonWebKeyTypeRSANegativeHSM),
	}
}

func (s *JsonWebKeyType) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseJsonWebKeyType(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseJsonWebKeyType(input string) (*JsonWebKeyType, error) {
	vals := map[string]JsonWebKeyType{
		"ec":      JsonWebKeyTypeEC,
		"ec-hsm":  JsonWebKeyTypeECNegativeHSM,
		"oct":     JsonWebKeyTypeOct,
		"oct-hsm": JsonWebKeyTypeOctNegativeHSM,
		"rsa":     JsonWebKeyTypeRSA,
		"rsa-hsm": JsonWebKeyTypeRSANegativeHSM,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := JsonWebKeyType(input)
	return &out, nil
}
