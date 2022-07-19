package videoanalyzer

import "strings"

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type AccessPolicyEccAlgo string

const (
	AccessPolicyEccAlgoESFiveOneTwo     AccessPolicyEccAlgo = "ES512"
	AccessPolicyEccAlgoESThreeEightFour AccessPolicyEccAlgo = "ES384"
	AccessPolicyEccAlgoESTwoFiveSix     AccessPolicyEccAlgo = "ES256"
)

func PossibleValuesForAccessPolicyEccAlgo() []string {
	return []string{
		string(AccessPolicyEccAlgoESFiveOneTwo),
		string(AccessPolicyEccAlgoESThreeEightFour),
		string(AccessPolicyEccAlgoESTwoFiveSix),
	}
}

func parseAccessPolicyEccAlgo(input string) (*AccessPolicyEccAlgo, error) {
	vals := map[string]AccessPolicyEccAlgo{
		"es512": AccessPolicyEccAlgoESFiveOneTwo,
		"es384": AccessPolicyEccAlgoESThreeEightFour,
		"es256": AccessPolicyEccAlgoESTwoFiveSix,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := AccessPolicyEccAlgo(input)
	return &out, nil
}

type AccessPolicyRole string

const (
	AccessPolicyRoleReader AccessPolicyRole = "Reader"
)

func PossibleValuesForAccessPolicyRole() []string {
	return []string{
		string(AccessPolicyRoleReader),
	}
}

func parseAccessPolicyRole(input string) (*AccessPolicyRole, error) {
	vals := map[string]AccessPolicyRole{
		"reader": AccessPolicyRoleReader,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := AccessPolicyRole(input)
	return &out, nil
}

type AccessPolicyRsaAlgo string

const (
	AccessPolicyRsaAlgoRSFiveOneTwo     AccessPolicyRsaAlgo = "RS512"
	AccessPolicyRsaAlgoRSThreeEightFour AccessPolicyRsaAlgo = "RS384"
	AccessPolicyRsaAlgoRSTwoFiveSix     AccessPolicyRsaAlgo = "RS256"
)

func PossibleValuesForAccessPolicyRsaAlgo() []string {
	return []string{
		string(AccessPolicyRsaAlgoRSFiveOneTwo),
		string(AccessPolicyRsaAlgoRSThreeEightFour),
		string(AccessPolicyRsaAlgoRSTwoFiveSix),
	}
}

func parseAccessPolicyRsaAlgo(input string) (*AccessPolicyRsaAlgo, error) {
	vals := map[string]AccessPolicyRsaAlgo{
		"rs512": AccessPolicyRsaAlgoRSFiveOneTwo,
		"rs384": AccessPolicyRsaAlgoRSThreeEightFour,
		"rs256": AccessPolicyRsaAlgoRSTwoFiveSix,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := AccessPolicyRsaAlgo(input)
	return &out, nil
}

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

type VideoType string

const (
	VideoTypeArchive VideoType = "Archive"
)

func PossibleValuesForVideoType() []string {
	return []string{
		string(VideoTypeArchive),
	}
}

func parseVideoType(input string) (*VideoType, error) {
	vals := map[string]VideoType{
		"archive": VideoTypeArchive,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := VideoType(input)
	return &out, nil
}
