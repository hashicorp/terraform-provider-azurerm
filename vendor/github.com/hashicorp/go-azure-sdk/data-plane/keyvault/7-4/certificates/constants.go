package certificates

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type CertificatePolicyAction string

const (
	CertificatePolicyActionAutoRenew     CertificatePolicyAction = "AutoRenew"
	CertificatePolicyActionEmailContacts CertificatePolicyAction = "EmailContacts"
)

func PossibleValuesForCertificatePolicyAction() []string {
	return []string{
		string(CertificatePolicyActionAutoRenew),
		string(CertificatePolicyActionEmailContacts),
	}
}

func (s *CertificatePolicyAction) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseCertificatePolicyAction(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseCertificatePolicyAction(input string) (*CertificatePolicyAction, error) {
	vals := map[string]CertificatePolicyAction{
		"autorenew":     CertificatePolicyActionAutoRenew,
		"emailcontacts": CertificatePolicyActionEmailContacts,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := CertificatePolicyAction(input)
	return &out, nil
}

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

type KeyUsageType string

const (
	KeyUsageTypeCRLSign          KeyUsageType = "cRLSign"
	KeyUsageTypeDataEncipherment KeyUsageType = "dataEncipherment"
	KeyUsageTypeDecipherOnly     KeyUsageType = "decipherOnly"
	KeyUsageTypeDigitalSignature KeyUsageType = "digitalSignature"
	KeyUsageTypeEncipherOnly     KeyUsageType = "encipherOnly"
	KeyUsageTypeKeyAgreement     KeyUsageType = "keyAgreement"
	KeyUsageTypeKeyCertSign      KeyUsageType = "keyCertSign"
	KeyUsageTypeKeyEncipherment  KeyUsageType = "keyEncipherment"
	KeyUsageTypeNonRepudiation   KeyUsageType = "nonRepudiation"
)

func PossibleValuesForKeyUsageType() []string {
	return []string{
		string(KeyUsageTypeCRLSign),
		string(KeyUsageTypeDataEncipherment),
		string(KeyUsageTypeDecipherOnly),
		string(KeyUsageTypeDigitalSignature),
		string(KeyUsageTypeEncipherOnly),
		string(KeyUsageTypeKeyAgreement),
		string(KeyUsageTypeKeyCertSign),
		string(KeyUsageTypeKeyEncipherment),
		string(KeyUsageTypeNonRepudiation),
	}
}

func (s *KeyUsageType) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseKeyUsageType(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseKeyUsageType(input string) (*KeyUsageType, error) {
	vals := map[string]KeyUsageType{
		"crlsign":          KeyUsageTypeCRLSign,
		"dataencipherment": KeyUsageTypeDataEncipherment,
		"decipheronly":     KeyUsageTypeDecipherOnly,
		"digitalsignature": KeyUsageTypeDigitalSignature,
		"encipheronly":     KeyUsageTypeEncipherOnly,
		"keyagreement":     KeyUsageTypeKeyAgreement,
		"keycertsign":      KeyUsageTypeKeyCertSign,
		"keyencipherment":  KeyUsageTypeKeyEncipherment,
		"nonrepudiation":   KeyUsageTypeNonRepudiation,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := KeyUsageType(input)
	return &out, nil
}
