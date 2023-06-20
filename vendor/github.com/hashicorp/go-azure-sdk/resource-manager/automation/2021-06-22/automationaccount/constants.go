package automationaccount

import "strings"

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type AutomationAccountState string

const (
	AutomationAccountStateOk          AutomationAccountState = "Ok"
	AutomationAccountStateSuspended   AutomationAccountState = "Suspended"
	AutomationAccountStateUnavailable AutomationAccountState = "Unavailable"
)

func PossibleValuesForAutomationAccountState() []string {
	return []string{
		string(AutomationAccountStateOk),
		string(AutomationAccountStateSuspended),
		string(AutomationAccountStateUnavailable),
	}
}

func parseAutomationAccountState(input string) (*AutomationAccountState, error) {
	vals := map[string]AutomationAccountState{
		"ok":          AutomationAccountStateOk,
		"suspended":   AutomationAccountStateSuspended,
		"unavailable": AutomationAccountStateUnavailable,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := AutomationAccountState(input)
	return &out, nil
}

type EncryptionKeySourceType string

const (
	EncryptionKeySourceTypeMicrosoftPointAutomation EncryptionKeySourceType = "Microsoft.Automation"
	EncryptionKeySourceTypeMicrosoftPointKeyvault   EncryptionKeySourceType = "Microsoft.Keyvault"
)

func PossibleValuesForEncryptionKeySourceType() []string {
	return []string{
		string(EncryptionKeySourceTypeMicrosoftPointAutomation),
		string(EncryptionKeySourceTypeMicrosoftPointKeyvault),
	}
}

func parseEncryptionKeySourceType(input string) (*EncryptionKeySourceType, error) {
	vals := map[string]EncryptionKeySourceType{
		"microsoft.automation": EncryptionKeySourceTypeMicrosoftPointAutomation,
		"microsoft.keyvault":   EncryptionKeySourceTypeMicrosoftPointKeyvault,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := EncryptionKeySourceType(input)
	return &out, nil
}

type SkuNameEnum string

const (
	SkuNameEnumBasic SkuNameEnum = "Basic"
	SkuNameEnumFree  SkuNameEnum = "Free"
)

func PossibleValuesForSkuNameEnum() []string {
	return []string{
		string(SkuNameEnumBasic),
		string(SkuNameEnumFree),
	}
}

func parseSkuNameEnum(input string) (*SkuNameEnum, error) {
	vals := map[string]SkuNameEnum{
		"basic": SkuNameEnumBasic,
		"free":  SkuNameEnumFree,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := SkuNameEnum(input)
	return &out, nil
}
