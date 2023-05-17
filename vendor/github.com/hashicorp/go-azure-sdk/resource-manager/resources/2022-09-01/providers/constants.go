package providers

import "strings"

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type AliasPathAttributes string

const (
	AliasPathAttributesModifiable AliasPathAttributes = "Modifiable"
	AliasPathAttributesNone       AliasPathAttributes = "None"
)

func PossibleValuesForAliasPathAttributes() []string {
	return []string{
		string(AliasPathAttributesModifiable),
		string(AliasPathAttributesNone),
	}
}

func parseAliasPathAttributes(input string) (*AliasPathAttributes, error) {
	vals := map[string]AliasPathAttributes{
		"modifiable": AliasPathAttributesModifiable,
		"none":       AliasPathAttributesNone,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := AliasPathAttributes(input)
	return &out, nil
}

type AliasPathTokenType string

const (
	AliasPathTokenTypeAny          AliasPathTokenType = "Any"
	AliasPathTokenTypeArray        AliasPathTokenType = "Array"
	AliasPathTokenTypeBoolean      AliasPathTokenType = "Boolean"
	AliasPathTokenTypeInteger      AliasPathTokenType = "Integer"
	AliasPathTokenTypeNotSpecified AliasPathTokenType = "NotSpecified"
	AliasPathTokenTypeNumber       AliasPathTokenType = "Number"
	AliasPathTokenTypeObject       AliasPathTokenType = "Object"
	AliasPathTokenTypeString       AliasPathTokenType = "String"
)

func PossibleValuesForAliasPathTokenType() []string {
	return []string{
		string(AliasPathTokenTypeAny),
		string(AliasPathTokenTypeArray),
		string(AliasPathTokenTypeBoolean),
		string(AliasPathTokenTypeInteger),
		string(AliasPathTokenTypeNotSpecified),
		string(AliasPathTokenTypeNumber),
		string(AliasPathTokenTypeObject),
		string(AliasPathTokenTypeString),
	}
}

func parseAliasPathTokenType(input string) (*AliasPathTokenType, error) {
	vals := map[string]AliasPathTokenType{
		"any":          AliasPathTokenTypeAny,
		"array":        AliasPathTokenTypeArray,
		"boolean":      AliasPathTokenTypeBoolean,
		"integer":      AliasPathTokenTypeInteger,
		"notspecified": AliasPathTokenTypeNotSpecified,
		"number":       AliasPathTokenTypeNumber,
		"object":       AliasPathTokenTypeObject,
		"string":       AliasPathTokenTypeString,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := AliasPathTokenType(input)
	return &out, nil
}

type AliasPatternType string

const (
	AliasPatternTypeExtract      AliasPatternType = "Extract"
	AliasPatternTypeNotSpecified AliasPatternType = "NotSpecified"
)

func PossibleValuesForAliasPatternType() []string {
	return []string{
		string(AliasPatternTypeExtract),
		string(AliasPatternTypeNotSpecified),
	}
}

func parseAliasPatternType(input string) (*AliasPatternType, error) {
	vals := map[string]AliasPatternType{
		"extract":      AliasPatternTypeExtract,
		"notspecified": AliasPatternTypeNotSpecified,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := AliasPatternType(input)
	return &out, nil
}

type AliasType string

const (
	AliasTypeMask         AliasType = "Mask"
	AliasTypeNotSpecified AliasType = "NotSpecified"
	AliasTypePlainText    AliasType = "PlainText"
)

func PossibleValuesForAliasType() []string {
	return []string{
		string(AliasTypeMask),
		string(AliasTypeNotSpecified),
		string(AliasTypePlainText),
	}
}

func parseAliasType(input string) (*AliasType, error) {
	vals := map[string]AliasType{
		"mask":         AliasTypeMask,
		"notspecified": AliasTypeNotSpecified,
		"plaintext":    AliasTypePlainText,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := AliasType(input)
	return &out, nil
}

type ProviderAuthorizationConsentState string

const (
	ProviderAuthorizationConsentStateConsented    ProviderAuthorizationConsentState = "Consented"
	ProviderAuthorizationConsentStateNotRequired  ProviderAuthorizationConsentState = "NotRequired"
	ProviderAuthorizationConsentStateNotSpecified ProviderAuthorizationConsentState = "NotSpecified"
	ProviderAuthorizationConsentStateRequired     ProviderAuthorizationConsentState = "Required"
)

func PossibleValuesForProviderAuthorizationConsentState() []string {
	return []string{
		string(ProviderAuthorizationConsentStateConsented),
		string(ProviderAuthorizationConsentStateNotRequired),
		string(ProviderAuthorizationConsentStateNotSpecified),
		string(ProviderAuthorizationConsentStateRequired),
	}
}

func parseProviderAuthorizationConsentState(input string) (*ProviderAuthorizationConsentState, error) {
	vals := map[string]ProviderAuthorizationConsentState{
		"consented":    ProviderAuthorizationConsentStateConsented,
		"notrequired":  ProviderAuthorizationConsentStateNotRequired,
		"notspecified": ProviderAuthorizationConsentStateNotSpecified,
		"required":     ProviderAuthorizationConsentStateRequired,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := ProviderAuthorizationConsentState(input)
	return &out, nil
}
