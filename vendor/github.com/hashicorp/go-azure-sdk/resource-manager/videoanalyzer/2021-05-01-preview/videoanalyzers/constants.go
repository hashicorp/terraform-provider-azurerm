package videoanalyzers

import "strings"

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type AccountEncryptionKeyType string

const (
	AccountEncryptionKeyTypeCustomerKey AccountEncryptionKeyType = "CustomerKey"
	AccountEncryptionKeyTypeSystemKey   AccountEncryptionKeyType = "SystemKey"
)

func PossibleValuesForAccountEncryptionKeyType() []string {
	return []string{
		string(AccountEncryptionKeyTypeCustomerKey),
		string(AccountEncryptionKeyTypeSystemKey),
	}
}

func parseAccountEncryptionKeyType(input string) (*AccountEncryptionKeyType, error) {
	vals := map[string]AccountEncryptionKeyType{
		"customerkey": AccountEncryptionKeyTypeCustomerKey,
		"systemkey":   AccountEncryptionKeyTypeSystemKey,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := AccountEncryptionKeyType(input)
	return &out, nil
}

type CheckNameAvailabilityReason string

const (
	CheckNameAvailabilityReasonAlreadyExists CheckNameAvailabilityReason = "AlreadyExists"
	CheckNameAvailabilityReasonInvalid       CheckNameAvailabilityReason = "Invalid"
)

func PossibleValuesForCheckNameAvailabilityReason() []string {
	return []string{
		string(CheckNameAvailabilityReasonAlreadyExists),
		string(CheckNameAvailabilityReasonInvalid),
	}
}

func parseCheckNameAvailabilityReason(input string) (*CheckNameAvailabilityReason, error) {
	vals := map[string]CheckNameAvailabilityReason{
		"alreadyexists": CheckNameAvailabilityReasonAlreadyExists,
		"invalid":       CheckNameAvailabilityReasonInvalid,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := CheckNameAvailabilityReason(input)
	return &out, nil
}

type VideoAnalyzerEndpointType string

const (
	VideoAnalyzerEndpointTypeClientApi VideoAnalyzerEndpointType = "ClientApi"
)

func PossibleValuesForVideoAnalyzerEndpointType() []string {
	return []string{
		string(VideoAnalyzerEndpointTypeClientApi),
	}
}

func parseVideoAnalyzerEndpointType(input string) (*VideoAnalyzerEndpointType, error) {
	vals := map[string]VideoAnalyzerEndpointType{
		"clientapi": VideoAnalyzerEndpointTypeClientApi,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := VideoAnalyzerEndpointType(input)
	return &out, nil
}
