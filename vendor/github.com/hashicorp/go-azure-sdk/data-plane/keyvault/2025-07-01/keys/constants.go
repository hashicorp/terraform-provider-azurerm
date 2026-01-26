package keys

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

type JsonWebKeyEncryptionAlgorithm string

const (
	JsonWebKeyEncryptionAlgorithmAOneNineTwoCBC                    JsonWebKeyEncryptionAlgorithm = "A192CBC"
	JsonWebKeyEncryptionAlgorithmAOneNineTwoCBCPAD                 JsonWebKeyEncryptionAlgorithm = "A192CBCPAD"
	JsonWebKeyEncryptionAlgorithmAOneNineTwoGCM                    JsonWebKeyEncryptionAlgorithm = "A192GCM"
	JsonWebKeyEncryptionAlgorithmAOneNineTwoKW                     JsonWebKeyEncryptionAlgorithm = "A192KW"
	JsonWebKeyEncryptionAlgorithmAOneTwoEightCBC                   JsonWebKeyEncryptionAlgorithm = "A128CBC"
	JsonWebKeyEncryptionAlgorithmAOneTwoEightCBCPAD                JsonWebKeyEncryptionAlgorithm = "A128CBCPAD"
	JsonWebKeyEncryptionAlgorithmAOneTwoEightGCM                   JsonWebKeyEncryptionAlgorithm = "A128GCM"
	JsonWebKeyEncryptionAlgorithmAOneTwoEightKW                    JsonWebKeyEncryptionAlgorithm = "A128KW"
	JsonWebKeyEncryptionAlgorithmATwoFiveSixCBC                    JsonWebKeyEncryptionAlgorithm = "A256CBC"
	JsonWebKeyEncryptionAlgorithmATwoFiveSixCBCPAD                 JsonWebKeyEncryptionAlgorithm = "A256CBCPAD"
	JsonWebKeyEncryptionAlgorithmATwoFiveSixGCM                    JsonWebKeyEncryptionAlgorithm = "A256GCM"
	JsonWebKeyEncryptionAlgorithmATwoFiveSixKW                     JsonWebKeyEncryptionAlgorithm = "A256KW"
	JsonWebKeyEncryptionAlgorithmCKMAESKEYWRAP                     JsonWebKeyEncryptionAlgorithm = "CKM_AES_KEY_WRAP"
	JsonWebKeyEncryptionAlgorithmCKMAESKEYWRAPPAD                  JsonWebKeyEncryptionAlgorithm = "CKM_AES_KEY_WRAP_PAD"
	JsonWebKeyEncryptionAlgorithmRSANegativeOAEP                   JsonWebKeyEncryptionAlgorithm = "RSA-OAEP"
	JsonWebKeyEncryptionAlgorithmRSANegativeOAEPNegativeTwoFiveSix JsonWebKeyEncryptionAlgorithm = "RSA-OAEP-256"
	JsonWebKeyEncryptionAlgorithmRSAOneFive                        JsonWebKeyEncryptionAlgorithm = "RSA1_5"
)

func PossibleValuesForJsonWebKeyEncryptionAlgorithm() []string {
	return []string{
		string(JsonWebKeyEncryptionAlgorithmAOneNineTwoCBC),
		string(JsonWebKeyEncryptionAlgorithmAOneNineTwoCBCPAD),
		string(JsonWebKeyEncryptionAlgorithmAOneNineTwoGCM),
		string(JsonWebKeyEncryptionAlgorithmAOneNineTwoKW),
		string(JsonWebKeyEncryptionAlgorithmAOneTwoEightCBC),
		string(JsonWebKeyEncryptionAlgorithmAOneTwoEightCBCPAD),
		string(JsonWebKeyEncryptionAlgorithmAOneTwoEightGCM),
		string(JsonWebKeyEncryptionAlgorithmAOneTwoEightKW),
		string(JsonWebKeyEncryptionAlgorithmATwoFiveSixCBC),
		string(JsonWebKeyEncryptionAlgorithmATwoFiveSixCBCPAD),
		string(JsonWebKeyEncryptionAlgorithmATwoFiveSixGCM),
		string(JsonWebKeyEncryptionAlgorithmATwoFiveSixKW),
		string(JsonWebKeyEncryptionAlgorithmCKMAESKEYWRAP),
		string(JsonWebKeyEncryptionAlgorithmCKMAESKEYWRAPPAD),
		string(JsonWebKeyEncryptionAlgorithmRSANegativeOAEP),
		string(JsonWebKeyEncryptionAlgorithmRSANegativeOAEPNegativeTwoFiveSix),
		string(JsonWebKeyEncryptionAlgorithmRSAOneFive),
	}
}

func (s *JsonWebKeyEncryptionAlgorithm) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseJsonWebKeyEncryptionAlgorithm(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseJsonWebKeyEncryptionAlgorithm(input string) (*JsonWebKeyEncryptionAlgorithm, error) {
	vals := map[string]JsonWebKeyEncryptionAlgorithm{
		"a192cbc":              JsonWebKeyEncryptionAlgorithmAOneNineTwoCBC,
		"a192cbcpad":           JsonWebKeyEncryptionAlgorithmAOneNineTwoCBCPAD,
		"a192gcm":              JsonWebKeyEncryptionAlgorithmAOneNineTwoGCM,
		"a192kw":               JsonWebKeyEncryptionAlgorithmAOneNineTwoKW,
		"a128cbc":              JsonWebKeyEncryptionAlgorithmAOneTwoEightCBC,
		"a128cbcpad":           JsonWebKeyEncryptionAlgorithmAOneTwoEightCBCPAD,
		"a128gcm":              JsonWebKeyEncryptionAlgorithmAOneTwoEightGCM,
		"a128kw":               JsonWebKeyEncryptionAlgorithmAOneTwoEightKW,
		"a256cbc":              JsonWebKeyEncryptionAlgorithmATwoFiveSixCBC,
		"a256cbcpad":           JsonWebKeyEncryptionAlgorithmATwoFiveSixCBCPAD,
		"a256gcm":              JsonWebKeyEncryptionAlgorithmATwoFiveSixGCM,
		"a256kw":               JsonWebKeyEncryptionAlgorithmATwoFiveSixKW,
		"ckm_aes_key_wrap":     JsonWebKeyEncryptionAlgorithmCKMAESKEYWRAP,
		"ckm_aes_key_wrap_pad": JsonWebKeyEncryptionAlgorithmCKMAESKEYWRAPPAD,
		"rsa-oaep":             JsonWebKeyEncryptionAlgorithmRSANegativeOAEP,
		"rsa-oaep-256":         JsonWebKeyEncryptionAlgorithmRSANegativeOAEPNegativeTwoFiveSix,
		"rsa1_5":               JsonWebKeyEncryptionAlgorithmRSAOneFive,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := JsonWebKeyEncryptionAlgorithm(input)
	return &out, nil
}

type JsonWebKeyOperation string

const (
	JsonWebKeyOperationDecrypt   JsonWebKeyOperation = "decrypt"
	JsonWebKeyOperationEncrypt   JsonWebKeyOperation = "encrypt"
	JsonWebKeyOperationExport    JsonWebKeyOperation = "export"
	JsonWebKeyOperationImport    JsonWebKeyOperation = "import"
	JsonWebKeyOperationSign      JsonWebKeyOperation = "sign"
	JsonWebKeyOperationUnwrapKey JsonWebKeyOperation = "unwrapKey"
	JsonWebKeyOperationVerify    JsonWebKeyOperation = "verify"
	JsonWebKeyOperationWrapKey   JsonWebKeyOperation = "wrapKey"
)

func PossibleValuesForJsonWebKeyOperation() []string {
	return []string{
		string(JsonWebKeyOperationDecrypt),
		string(JsonWebKeyOperationEncrypt),
		string(JsonWebKeyOperationExport),
		string(JsonWebKeyOperationImport),
		string(JsonWebKeyOperationSign),
		string(JsonWebKeyOperationUnwrapKey),
		string(JsonWebKeyOperationVerify),
		string(JsonWebKeyOperationWrapKey),
	}
}

func (s *JsonWebKeyOperation) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseJsonWebKeyOperation(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseJsonWebKeyOperation(input string) (*JsonWebKeyOperation, error) {
	vals := map[string]JsonWebKeyOperation{
		"decrypt":   JsonWebKeyOperationDecrypt,
		"encrypt":   JsonWebKeyOperationEncrypt,
		"export":    JsonWebKeyOperationExport,
		"import":    JsonWebKeyOperationImport,
		"sign":      JsonWebKeyOperationSign,
		"unwrapkey": JsonWebKeyOperationUnwrapKey,
		"verify":    JsonWebKeyOperationVerify,
		"wrapkey":   JsonWebKeyOperationWrapKey,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := JsonWebKeyOperation(input)
	return &out, nil
}

type JsonWebKeySignatureAlgorithm string

const (
	JsonWebKeySignatureAlgorithmESFiveOneTwo     JsonWebKeySignatureAlgorithm = "ES512"
	JsonWebKeySignatureAlgorithmESThreeEightFour JsonWebKeySignatureAlgorithm = "ES384"
	JsonWebKeySignatureAlgorithmESTwoFiveSix     JsonWebKeySignatureAlgorithm = "ES256"
	JsonWebKeySignatureAlgorithmESTwoFiveSixK    JsonWebKeySignatureAlgorithm = "ES256K"
	JsonWebKeySignatureAlgorithmHSFiveOneTwo     JsonWebKeySignatureAlgorithm = "HS512"
	JsonWebKeySignatureAlgorithmHSThreeEightFour JsonWebKeySignatureAlgorithm = "HS384"
	JsonWebKeySignatureAlgorithmHSTwoFiveSix     JsonWebKeySignatureAlgorithm = "HS256"
	JsonWebKeySignatureAlgorithmPSFiveOneTwo     JsonWebKeySignatureAlgorithm = "PS512"
	JsonWebKeySignatureAlgorithmPSThreeEightFour JsonWebKeySignatureAlgorithm = "PS384"
	JsonWebKeySignatureAlgorithmPSTwoFiveSix     JsonWebKeySignatureAlgorithm = "PS256"
	JsonWebKeySignatureAlgorithmRSFiveOneTwo     JsonWebKeySignatureAlgorithm = "RS512"
	JsonWebKeySignatureAlgorithmRSNULL           JsonWebKeySignatureAlgorithm = "RSNULL"
	JsonWebKeySignatureAlgorithmRSThreeEightFour JsonWebKeySignatureAlgorithm = "RS384"
	JsonWebKeySignatureAlgorithmRSTwoFiveSix     JsonWebKeySignatureAlgorithm = "RS256"
)

func PossibleValuesForJsonWebKeySignatureAlgorithm() []string {
	return []string{
		string(JsonWebKeySignatureAlgorithmESFiveOneTwo),
		string(JsonWebKeySignatureAlgorithmESThreeEightFour),
		string(JsonWebKeySignatureAlgorithmESTwoFiveSix),
		string(JsonWebKeySignatureAlgorithmESTwoFiveSixK),
		string(JsonWebKeySignatureAlgorithmHSFiveOneTwo),
		string(JsonWebKeySignatureAlgorithmHSThreeEightFour),
		string(JsonWebKeySignatureAlgorithmHSTwoFiveSix),
		string(JsonWebKeySignatureAlgorithmPSFiveOneTwo),
		string(JsonWebKeySignatureAlgorithmPSThreeEightFour),
		string(JsonWebKeySignatureAlgorithmPSTwoFiveSix),
		string(JsonWebKeySignatureAlgorithmRSFiveOneTwo),
		string(JsonWebKeySignatureAlgorithmRSNULL),
		string(JsonWebKeySignatureAlgorithmRSThreeEightFour),
		string(JsonWebKeySignatureAlgorithmRSTwoFiveSix),
	}
}

func (s *JsonWebKeySignatureAlgorithm) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseJsonWebKeySignatureAlgorithm(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseJsonWebKeySignatureAlgorithm(input string) (*JsonWebKeySignatureAlgorithm, error) {
	vals := map[string]JsonWebKeySignatureAlgorithm{
		"es512":  JsonWebKeySignatureAlgorithmESFiveOneTwo,
		"es384":  JsonWebKeySignatureAlgorithmESThreeEightFour,
		"es256":  JsonWebKeySignatureAlgorithmESTwoFiveSix,
		"es256k": JsonWebKeySignatureAlgorithmESTwoFiveSixK,
		"hs512":  JsonWebKeySignatureAlgorithmHSFiveOneTwo,
		"hs384":  JsonWebKeySignatureAlgorithmHSThreeEightFour,
		"hs256":  JsonWebKeySignatureAlgorithmHSTwoFiveSix,
		"ps512":  JsonWebKeySignatureAlgorithmPSFiveOneTwo,
		"ps384":  JsonWebKeySignatureAlgorithmPSThreeEightFour,
		"ps256":  JsonWebKeySignatureAlgorithmPSTwoFiveSix,
		"rs512":  JsonWebKeySignatureAlgorithmRSFiveOneTwo,
		"rsnull": JsonWebKeySignatureAlgorithmRSNULL,
		"rs384":  JsonWebKeySignatureAlgorithmRSThreeEightFour,
		"rs256":  JsonWebKeySignatureAlgorithmRSTwoFiveSix,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := JsonWebKeySignatureAlgorithm(input)
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

type KeyEncryptionAlgorithm string

const (
	KeyEncryptionAlgorithmCKMRSAAESKEYWRAP            KeyEncryptionAlgorithm = "CKM_RSA_AES_KEY_WRAP"
	KeyEncryptionAlgorithmRSAAESKEYWRAPThreeEightFour KeyEncryptionAlgorithm = "RSA_AES_KEY_WRAP_384"
	KeyEncryptionAlgorithmRSAAESKEYWRAPTwoFiveSix     KeyEncryptionAlgorithm = "RSA_AES_KEY_WRAP_256"
)

func PossibleValuesForKeyEncryptionAlgorithm() []string {
	return []string{
		string(KeyEncryptionAlgorithmCKMRSAAESKEYWRAP),
		string(KeyEncryptionAlgorithmRSAAESKEYWRAPThreeEightFour),
		string(KeyEncryptionAlgorithmRSAAESKEYWRAPTwoFiveSix),
	}
}

func (s *KeyEncryptionAlgorithm) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseKeyEncryptionAlgorithm(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseKeyEncryptionAlgorithm(input string) (*KeyEncryptionAlgorithm, error) {
	vals := map[string]KeyEncryptionAlgorithm{
		"ckm_rsa_aes_key_wrap": KeyEncryptionAlgorithmCKMRSAAESKEYWRAP,
		"rsa_aes_key_wrap_384": KeyEncryptionAlgorithmRSAAESKEYWRAPThreeEightFour,
		"rsa_aes_key_wrap_256": KeyEncryptionAlgorithmRSAAESKEYWRAPTwoFiveSix,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := KeyEncryptionAlgorithm(input)
	return &out, nil
}

type KeyRotationPolicyAction string

const (
	KeyRotationPolicyActionNotify KeyRotationPolicyAction = "Notify"
	KeyRotationPolicyActionRotate KeyRotationPolicyAction = "Rotate"
)

func PossibleValuesForKeyRotationPolicyAction() []string {
	return []string{
		string(KeyRotationPolicyActionNotify),
		string(KeyRotationPolicyActionRotate),
	}
}

func (s *KeyRotationPolicyAction) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseKeyRotationPolicyAction(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseKeyRotationPolicyAction(input string) (*KeyRotationPolicyAction, error) {
	vals := map[string]KeyRotationPolicyAction{
		"notify": KeyRotationPolicyActionNotify,
		"rotate": KeyRotationPolicyActionRotate,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := KeyRotationPolicyAction(input)
	return &out, nil
}
