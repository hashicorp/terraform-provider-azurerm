package backupprotectableitems

import (
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type AzureFileShareType string

const (
	AzureFileShareTypeInvalid AzureFileShareType = "Invalid"
	AzureFileShareTypeXSMB    AzureFileShareType = "XSMB"
	AzureFileShareTypeXSync   AzureFileShareType = "XSync"
)

func PossibleValuesForAzureFileShareType() []string {
	return []string{
		string(AzureFileShareTypeInvalid),
		string(AzureFileShareTypeXSMB),
		string(AzureFileShareTypeXSync),
	}
}

func parseAzureFileShareType(input string) (*AzureFileShareType, error) {
	vals := map[string]AzureFileShareType{
		"invalid": AzureFileShareTypeInvalid,
		"xsmb":    AzureFileShareTypeXSMB,
		"xsync":   AzureFileShareTypeXSync,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := AzureFileShareType(input)
	return &out, nil
}

type InquiryStatus string

const (
	InquiryStatusFailed  InquiryStatus = "Failed"
	InquiryStatusInvalid InquiryStatus = "Invalid"
	InquiryStatusSuccess InquiryStatus = "Success"
)

func PossibleValuesForInquiryStatus() []string {
	return []string{
		string(InquiryStatusFailed),
		string(InquiryStatusInvalid),
		string(InquiryStatusSuccess),
	}
}

func parseInquiryStatus(input string) (*InquiryStatus, error) {
	vals := map[string]InquiryStatus{
		"failed":  InquiryStatusFailed,
		"invalid": InquiryStatusInvalid,
		"success": InquiryStatusSuccess,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := InquiryStatus(input)
	return &out, nil
}

type ProtectionStatus string

const (
	ProtectionStatusInvalid          ProtectionStatus = "Invalid"
	ProtectionStatusNotProtected     ProtectionStatus = "NotProtected"
	ProtectionStatusProtected        ProtectionStatus = "Protected"
	ProtectionStatusProtecting       ProtectionStatus = "Protecting"
	ProtectionStatusProtectionFailed ProtectionStatus = "ProtectionFailed"
)

func PossibleValuesForProtectionStatus() []string {
	return []string{
		string(ProtectionStatusInvalid),
		string(ProtectionStatusNotProtected),
		string(ProtectionStatusProtected),
		string(ProtectionStatusProtecting),
		string(ProtectionStatusProtectionFailed),
	}
}

func parseProtectionStatus(input string) (*ProtectionStatus, error) {
	vals := map[string]ProtectionStatus{
		"invalid":          ProtectionStatusInvalid,
		"notprotected":     ProtectionStatusNotProtected,
		"protected":        ProtectionStatusProtected,
		"protecting":       ProtectionStatusProtecting,
		"protectionfailed": ProtectionStatusProtectionFailed,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := ProtectionStatus(input)
	return &out, nil
}
