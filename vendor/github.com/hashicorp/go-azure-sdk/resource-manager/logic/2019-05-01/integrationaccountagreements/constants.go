package integrationaccountagreements

import "strings"

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type AgreementType string

const (
	AgreementTypeASTwo        AgreementType = "AS2"
	AgreementTypeEdifact      AgreementType = "Edifact"
	AgreementTypeNotSpecified AgreementType = "NotSpecified"
	AgreementTypeXOneTwo      AgreementType = "X12"
)

func PossibleValuesForAgreementType() []string {
	return []string{
		string(AgreementTypeASTwo),
		string(AgreementTypeEdifact),
		string(AgreementTypeNotSpecified),
		string(AgreementTypeXOneTwo),
	}
}

func parseAgreementType(input string) (*AgreementType, error) {
	vals := map[string]AgreementType{
		"as2":          AgreementTypeASTwo,
		"edifact":      AgreementTypeEdifact,
		"notspecified": AgreementTypeNotSpecified,
		"x12":          AgreementTypeXOneTwo,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := AgreementType(input)
	return &out, nil
}

type EdifactCharacterSet string

const (
	EdifactCharacterSetKECA         EdifactCharacterSet = "KECA"
	EdifactCharacterSetNotSpecified EdifactCharacterSet = "NotSpecified"
	EdifactCharacterSetUNOA         EdifactCharacterSet = "UNOA"
	EdifactCharacterSetUNOB         EdifactCharacterSet = "UNOB"
	EdifactCharacterSetUNOC         EdifactCharacterSet = "UNOC"
	EdifactCharacterSetUNOD         EdifactCharacterSet = "UNOD"
	EdifactCharacterSetUNOE         EdifactCharacterSet = "UNOE"
	EdifactCharacterSetUNOF         EdifactCharacterSet = "UNOF"
	EdifactCharacterSetUNOG         EdifactCharacterSet = "UNOG"
	EdifactCharacterSetUNOH         EdifactCharacterSet = "UNOH"
	EdifactCharacterSetUNOI         EdifactCharacterSet = "UNOI"
	EdifactCharacterSetUNOJ         EdifactCharacterSet = "UNOJ"
	EdifactCharacterSetUNOK         EdifactCharacterSet = "UNOK"
	EdifactCharacterSetUNOX         EdifactCharacterSet = "UNOX"
	EdifactCharacterSetUNOY         EdifactCharacterSet = "UNOY"
)

func PossibleValuesForEdifactCharacterSet() []string {
	return []string{
		string(EdifactCharacterSetKECA),
		string(EdifactCharacterSetNotSpecified),
		string(EdifactCharacterSetUNOA),
		string(EdifactCharacterSetUNOB),
		string(EdifactCharacterSetUNOC),
		string(EdifactCharacterSetUNOD),
		string(EdifactCharacterSetUNOE),
		string(EdifactCharacterSetUNOF),
		string(EdifactCharacterSetUNOG),
		string(EdifactCharacterSetUNOH),
		string(EdifactCharacterSetUNOI),
		string(EdifactCharacterSetUNOJ),
		string(EdifactCharacterSetUNOK),
		string(EdifactCharacterSetUNOX),
		string(EdifactCharacterSetUNOY),
	}
}

func parseEdifactCharacterSet(input string) (*EdifactCharacterSet, error) {
	vals := map[string]EdifactCharacterSet{
		"keca":         EdifactCharacterSetKECA,
		"notspecified": EdifactCharacterSetNotSpecified,
		"unoa":         EdifactCharacterSetUNOA,
		"unob":         EdifactCharacterSetUNOB,
		"unoc":         EdifactCharacterSetUNOC,
		"unod":         EdifactCharacterSetUNOD,
		"unoe":         EdifactCharacterSetUNOE,
		"unof":         EdifactCharacterSetUNOF,
		"unog":         EdifactCharacterSetUNOG,
		"unoh":         EdifactCharacterSetUNOH,
		"unoi":         EdifactCharacterSetUNOI,
		"unoj":         EdifactCharacterSetUNOJ,
		"unok":         EdifactCharacterSetUNOK,
		"unox":         EdifactCharacterSetUNOX,
		"unoy":         EdifactCharacterSetUNOY,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := EdifactCharacterSet(input)
	return &out, nil
}

type EdifactDecimalIndicator string

const (
	EdifactDecimalIndicatorComma        EdifactDecimalIndicator = "Comma"
	EdifactDecimalIndicatorDecimal      EdifactDecimalIndicator = "Decimal"
	EdifactDecimalIndicatorNotSpecified EdifactDecimalIndicator = "NotSpecified"
)

func PossibleValuesForEdifactDecimalIndicator() []string {
	return []string{
		string(EdifactDecimalIndicatorComma),
		string(EdifactDecimalIndicatorDecimal),
		string(EdifactDecimalIndicatorNotSpecified),
	}
}

func parseEdifactDecimalIndicator(input string) (*EdifactDecimalIndicator, error) {
	vals := map[string]EdifactDecimalIndicator{
		"comma":        EdifactDecimalIndicatorComma,
		"decimal":      EdifactDecimalIndicatorDecimal,
		"notspecified": EdifactDecimalIndicatorNotSpecified,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := EdifactDecimalIndicator(input)
	return &out, nil
}

type EncryptionAlgorithm string

const (
	EncryptionAlgorithmAESOneNineTwo  EncryptionAlgorithm = "AES192"
	EncryptionAlgorithmAESOneTwoEight EncryptionAlgorithm = "AES128"
	EncryptionAlgorithmAESTwoFiveSix  EncryptionAlgorithm = "AES256"
	EncryptionAlgorithmDESThree       EncryptionAlgorithm = "DES3"
	EncryptionAlgorithmNone           EncryptionAlgorithm = "None"
	EncryptionAlgorithmNotSpecified   EncryptionAlgorithm = "NotSpecified"
	EncryptionAlgorithmRCTwo          EncryptionAlgorithm = "RC2"
)

func PossibleValuesForEncryptionAlgorithm() []string {
	return []string{
		string(EncryptionAlgorithmAESOneNineTwo),
		string(EncryptionAlgorithmAESOneTwoEight),
		string(EncryptionAlgorithmAESTwoFiveSix),
		string(EncryptionAlgorithmDESThree),
		string(EncryptionAlgorithmNone),
		string(EncryptionAlgorithmNotSpecified),
		string(EncryptionAlgorithmRCTwo),
	}
}

func parseEncryptionAlgorithm(input string) (*EncryptionAlgorithm, error) {
	vals := map[string]EncryptionAlgorithm{
		"aes192":       EncryptionAlgorithmAESOneNineTwo,
		"aes128":       EncryptionAlgorithmAESOneTwoEight,
		"aes256":       EncryptionAlgorithmAESTwoFiveSix,
		"des3":         EncryptionAlgorithmDESThree,
		"none":         EncryptionAlgorithmNone,
		"notspecified": EncryptionAlgorithmNotSpecified,
		"rc2":          EncryptionAlgorithmRCTwo,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := EncryptionAlgorithm(input)
	return &out, nil
}

type HashingAlgorithm string

const (
	HashingAlgorithmMDFive               HashingAlgorithm = "MD5"
	HashingAlgorithmNone                 HashingAlgorithm = "None"
	HashingAlgorithmNotSpecified         HashingAlgorithm = "NotSpecified"
	HashingAlgorithmSHAOne               HashingAlgorithm = "SHA1"
	HashingAlgorithmSHATwoFiveOneTwo     HashingAlgorithm = "SHA2512"
	HashingAlgorithmSHATwoThreeEightFour HashingAlgorithm = "SHA2384"
	HashingAlgorithmSHATwoTwoFiveSix     HashingAlgorithm = "SHA2256"
)

func PossibleValuesForHashingAlgorithm() []string {
	return []string{
		string(HashingAlgorithmMDFive),
		string(HashingAlgorithmNone),
		string(HashingAlgorithmNotSpecified),
		string(HashingAlgorithmSHAOne),
		string(HashingAlgorithmSHATwoFiveOneTwo),
		string(HashingAlgorithmSHATwoThreeEightFour),
		string(HashingAlgorithmSHATwoTwoFiveSix),
	}
}

func parseHashingAlgorithm(input string) (*HashingAlgorithm, error) {
	vals := map[string]HashingAlgorithm{
		"md5":          HashingAlgorithmMDFive,
		"none":         HashingAlgorithmNone,
		"notspecified": HashingAlgorithmNotSpecified,
		"sha1":         HashingAlgorithmSHAOne,
		"sha2512":      HashingAlgorithmSHATwoFiveOneTwo,
		"sha2384":      HashingAlgorithmSHATwoThreeEightFour,
		"sha2256":      HashingAlgorithmSHATwoTwoFiveSix,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := HashingAlgorithm(input)
	return &out, nil
}

type KeyType string

const (
	KeyTypeNotSpecified KeyType = "NotSpecified"
	KeyTypePrimary      KeyType = "Primary"
	KeyTypeSecondary    KeyType = "Secondary"
)

func PossibleValuesForKeyType() []string {
	return []string{
		string(KeyTypeNotSpecified),
		string(KeyTypePrimary),
		string(KeyTypeSecondary),
	}
}

func parseKeyType(input string) (*KeyType, error) {
	vals := map[string]KeyType{
		"notspecified": KeyTypeNotSpecified,
		"primary":      KeyTypePrimary,
		"secondary":    KeyTypeSecondary,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := KeyType(input)
	return &out, nil
}

type MessageFilterType string

const (
	MessageFilterTypeExclude      MessageFilterType = "Exclude"
	MessageFilterTypeInclude      MessageFilterType = "Include"
	MessageFilterTypeNotSpecified MessageFilterType = "NotSpecified"
)

func PossibleValuesForMessageFilterType() []string {
	return []string{
		string(MessageFilterTypeExclude),
		string(MessageFilterTypeInclude),
		string(MessageFilterTypeNotSpecified),
	}
}

func parseMessageFilterType(input string) (*MessageFilterType, error) {
	vals := map[string]MessageFilterType{
		"exclude":      MessageFilterTypeExclude,
		"include":      MessageFilterTypeInclude,
		"notspecified": MessageFilterTypeNotSpecified,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := MessageFilterType(input)
	return &out, nil
}

type SegmentTerminatorSuffix string

const (
	SegmentTerminatorSuffixCR           SegmentTerminatorSuffix = "CR"
	SegmentTerminatorSuffixCRLF         SegmentTerminatorSuffix = "CRLF"
	SegmentTerminatorSuffixLF           SegmentTerminatorSuffix = "LF"
	SegmentTerminatorSuffixNone         SegmentTerminatorSuffix = "None"
	SegmentTerminatorSuffixNotSpecified SegmentTerminatorSuffix = "NotSpecified"
)

func PossibleValuesForSegmentTerminatorSuffix() []string {
	return []string{
		string(SegmentTerminatorSuffixCR),
		string(SegmentTerminatorSuffixCRLF),
		string(SegmentTerminatorSuffixLF),
		string(SegmentTerminatorSuffixNone),
		string(SegmentTerminatorSuffixNotSpecified),
	}
}

func parseSegmentTerminatorSuffix(input string) (*SegmentTerminatorSuffix, error) {
	vals := map[string]SegmentTerminatorSuffix{
		"cr":           SegmentTerminatorSuffixCR,
		"crlf":         SegmentTerminatorSuffixCRLF,
		"lf":           SegmentTerminatorSuffixLF,
		"none":         SegmentTerminatorSuffixNone,
		"notspecified": SegmentTerminatorSuffixNotSpecified,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := SegmentTerminatorSuffix(input)
	return &out, nil
}

type SigningAlgorithm string

const (
	SigningAlgorithmDefault              SigningAlgorithm = "Default"
	SigningAlgorithmNotSpecified         SigningAlgorithm = "NotSpecified"
	SigningAlgorithmSHAOne               SigningAlgorithm = "SHA1"
	SigningAlgorithmSHATwoFiveOneTwo     SigningAlgorithm = "SHA2512"
	SigningAlgorithmSHATwoThreeEightFour SigningAlgorithm = "SHA2384"
	SigningAlgorithmSHATwoTwoFiveSix     SigningAlgorithm = "SHA2256"
)

func PossibleValuesForSigningAlgorithm() []string {
	return []string{
		string(SigningAlgorithmDefault),
		string(SigningAlgorithmNotSpecified),
		string(SigningAlgorithmSHAOne),
		string(SigningAlgorithmSHATwoFiveOneTwo),
		string(SigningAlgorithmSHATwoThreeEightFour),
		string(SigningAlgorithmSHATwoTwoFiveSix),
	}
}

func parseSigningAlgorithm(input string) (*SigningAlgorithm, error) {
	vals := map[string]SigningAlgorithm{
		"default":      SigningAlgorithmDefault,
		"notspecified": SigningAlgorithmNotSpecified,
		"sha1":         SigningAlgorithmSHAOne,
		"sha2512":      SigningAlgorithmSHATwoFiveOneTwo,
		"sha2384":      SigningAlgorithmSHATwoThreeEightFour,
		"sha2256":      SigningAlgorithmSHATwoTwoFiveSix,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := SigningAlgorithm(input)
	return &out, nil
}

type TrailingSeparatorPolicy string

const (
	TrailingSeparatorPolicyMandatory    TrailingSeparatorPolicy = "Mandatory"
	TrailingSeparatorPolicyNotAllowed   TrailingSeparatorPolicy = "NotAllowed"
	TrailingSeparatorPolicyNotSpecified TrailingSeparatorPolicy = "NotSpecified"
	TrailingSeparatorPolicyOptional     TrailingSeparatorPolicy = "Optional"
)

func PossibleValuesForTrailingSeparatorPolicy() []string {
	return []string{
		string(TrailingSeparatorPolicyMandatory),
		string(TrailingSeparatorPolicyNotAllowed),
		string(TrailingSeparatorPolicyNotSpecified),
		string(TrailingSeparatorPolicyOptional),
	}
}

func parseTrailingSeparatorPolicy(input string) (*TrailingSeparatorPolicy, error) {
	vals := map[string]TrailingSeparatorPolicy{
		"mandatory":    TrailingSeparatorPolicyMandatory,
		"notallowed":   TrailingSeparatorPolicyNotAllowed,
		"notspecified": TrailingSeparatorPolicyNotSpecified,
		"optional":     TrailingSeparatorPolicyOptional,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := TrailingSeparatorPolicy(input)
	return &out, nil
}

type UsageIndicator string

const (
	UsageIndicatorInformation  UsageIndicator = "Information"
	UsageIndicatorNotSpecified UsageIndicator = "NotSpecified"
	UsageIndicatorProduction   UsageIndicator = "Production"
	UsageIndicatorTest         UsageIndicator = "Test"
)

func PossibleValuesForUsageIndicator() []string {
	return []string{
		string(UsageIndicatorInformation),
		string(UsageIndicatorNotSpecified),
		string(UsageIndicatorProduction),
		string(UsageIndicatorTest),
	}
}

func parseUsageIndicator(input string) (*UsageIndicator, error) {
	vals := map[string]UsageIndicator{
		"information":  UsageIndicatorInformation,
		"notspecified": UsageIndicatorNotSpecified,
		"production":   UsageIndicatorProduction,
		"test":         UsageIndicatorTest,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := UsageIndicator(input)
	return &out, nil
}

type X12CharacterSet string

const (
	X12CharacterSetBasic        X12CharacterSet = "Basic"
	X12CharacterSetExtended     X12CharacterSet = "Extended"
	X12CharacterSetNotSpecified X12CharacterSet = "NotSpecified"
	X12CharacterSetUTFEight     X12CharacterSet = "UTF8"
)

func PossibleValuesForX12CharacterSet() []string {
	return []string{
		string(X12CharacterSetBasic),
		string(X12CharacterSetExtended),
		string(X12CharacterSetNotSpecified),
		string(X12CharacterSetUTFEight),
	}
}

func parseX12CharacterSet(input string) (*X12CharacterSet, error) {
	vals := map[string]X12CharacterSet{
		"basic":        X12CharacterSetBasic,
		"extended":     X12CharacterSetExtended,
		"notspecified": X12CharacterSetNotSpecified,
		"utf8":         X12CharacterSetUTFEight,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := X12CharacterSet(input)
	return &out, nil
}

type X12DateFormat string

const (
	X12DateFormatCCYYMMDD     X12DateFormat = "CCYYMMDD"
	X12DateFormatNotSpecified X12DateFormat = "NotSpecified"
	X12DateFormatYYMMDD       X12DateFormat = "YYMMDD"
)

func PossibleValuesForX12DateFormat() []string {
	return []string{
		string(X12DateFormatCCYYMMDD),
		string(X12DateFormatNotSpecified),
		string(X12DateFormatYYMMDD),
	}
}

func parseX12DateFormat(input string) (*X12DateFormat, error) {
	vals := map[string]X12DateFormat{
		"ccyymmdd":     X12DateFormatCCYYMMDD,
		"notspecified": X12DateFormatNotSpecified,
		"yymmdd":       X12DateFormatYYMMDD,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := X12DateFormat(input)
	return &out, nil
}

type X12TimeFormat string

const (
	X12TimeFormatHHMM         X12TimeFormat = "HHMM"
	X12TimeFormatHHMMSS       X12TimeFormat = "HHMMSS"
	X12TimeFormatHHMMSSd      X12TimeFormat = "HHMMSSd"
	X12TimeFormatHHMMSSdd     X12TimeFormat = "HHMMSSdd"
	X12TimeFormatNotSpecified X12TimeFormat = "NotSpecified"
)

func PossibleValuesForX12TimeFormat() []string {
	return []string{
		string(X12TimeFormatHHMM),
		string(X12TimeFormatHHMMSS),
		string(X12TimeFormatHHMMSSd),
		string(X12TimeFormatHHMMSSdd),
		string(X12TimeFormatNotSpecified),
	}
}

func parseX12TimeFormat(input string) (*X12TimeFormat, error) {
	vals := map[string]X12TimeFormat{
		"hhmm":         X12TimeFormatHHMM,
		"hhmmss":       X12TimeFormatHHMMSS,
		"hhmmssd":      X12TimeFormatHHMMSSd,
		"hhmmssdd":     X12TimeFormatHHMMSSdd,
		"notspecified": X12TimeFormatNotSpecified,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := X12TimeFormat(input)
	return &out, nil
}
