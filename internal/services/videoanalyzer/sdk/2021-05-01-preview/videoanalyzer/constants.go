package videoanalyzer

import "strings"

type AccessPolicyEccAlgo string

const (
	AccessPolicyEccAlgoESFiveOneTwo     AccessPolicyEccAlgo = "ES512"
	AccessPolicyEccAlgoESThreeEightFour AccessPolicyEccAlgo = "ES384"
	AccessPolicyEccAlgoESTwoFiveSix     AccessPolicyEccAlgo = "ES256"
)

func PossibleValuesForAccessPolicyEccAlgo() []string {
	return []string{
		"ES512",
		"ES384",
		"ES256",
	}
}

func parseAccessPolicyEccAlgo(input string) (*AccessPolicyEccAlgo, error) {
	vals := map[string]AccessPolicyEccAlgo{
		"esfiveonetwo":     "ES512",
		"esthreeeightfour": "ES384",
		"estwofivesix":     "ES256",
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// it could be a new value - best effort convert this
	v := input

	out := AccessPolicyEccAlgo(v)
	return &out, nil
}

type AccessPolicyRole string

const (
	AccessPolicyRoleReader AccessPolicyRole = "Reader"
)

func PossibleValuesForAccessPolicyRole() []string {
	return []string{
		"Reader",
	}
}

func parseAccessPolicyRole(input string) (*AccessPolicyRole, error) {
	vals := map[string]AccessPolicyRole{
		"reader": "Reader",
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// it could be a new value - best effort convert this
	v := input

	out := AccessPolicyRole(v)
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
		"RS512",
		"RS384",
		"RS256",
	}
}

func parseAccessPolicyRsaAlgo(input string) (*AccessPolicyRsaAlgo, error) {
	vals := map[string]AccessPolicyRsaAlgo{
		"rsfiveonetwo":     "RS512",
		"rsthreeeightfour": "RS384",
		"rstwofivesix":     "RS256",
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// it could be a new value - best effort convert this
	v := input

	out := AccessPolicyRsaAlgo(v)
	return &out, nil
}

type AccountEncryptionKeyType string

const (
	AccountEncryptionKeyTypeCustomerKey AccountEncryptionKeyType = "CustomerKey"
	AccountEncryptionKeyTypeSystemKey   AccountEncryptionKeyType = "SystemKey"
)

func PossibleValuesForAccountEncryptionKeyType() []string {
	return []string{
		"CustomerKey",
		"SystemKey",
	}
}

func parseAccountEncryptionKeyType(input string) (*AccountEncryptionKeyType, error) {
	vals := map[string]AccountEncryptionKeyType{
		"customerkey": "CustomerKey",
		"systemkey":   "SystemKey",
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// it could be a new value - best effort convert this
	v := input

	out := AccountEncryptionKeyType(v)
	return &out, nil
}

type CheckNameAvailabilityReason string

const (
	CheckNameAvailabilityReasonAlreadyExists CheckNameAvailabilityReason = "AlreadyExists"
	CheckNameAvailabilityReasonInvalid       CheckNameAvailabilityReason = "Invalid"
)

func PossibleValuesForCheckNameAvailabilityReason() []string {
	return []string{
		"AlreadyExists",
		"Invalid",
	}
}

func parseCheckNameAvailabilityReason(input string) (*CheckNameAvailabilityReason, error) {
	vals := map[string]CheckNameAvailabilityReason{
		"alreadyexists": "AlreadyExists",
		"invalid":       "Invalid",
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// it could be a new value - best effort convert this
	v := input

	out := CheckNameAvailabilityReason(v)
	return &out, nil
}

type CreatedByType string

const (
	CreatedByTypeApplication     CreatedByType = "Application"
	CreatedByTypeKey             CreatedByType = "Key"
	CreatedByTypeManagedIdentity CreatedByType = "ManagedIdentity"
	CreatedByTypeUser            CreatedByType = "User"
)

func PossibleValuesForCreatedByType() []string {
	return []string{
		"Application",
		"Key",
		"ManagedIdentity",
		"User",
	}
}

func parseCreatedByType(input string) (*CreatedByType, error) {
	vals := map[string]CreatedByType{
		"application":     "Application",
		"key":             "Key",
		"managedidentity": "ManagedIdentity",
		"user":            "User",
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// it could be a new value - best effort convert this
	v := input

	out := CreatedByType(v)
	return &out, nil
}

type VideoAnalyzerEndpointType string

const (
	VideoAnalyzerEndpointTypeClientApi VideoAnalyzerEndpointType = "ClientApi"
)

func PossibleValuesForVideoAnalyzerEndpointType() []string {
	return []string{
		"ClientApi",
	}
}

func parseVideoAnalyzerEndpointType(input string) (*VideoAnalyzerEndpointType, error) {
	vals := map[string]VideoAnalyzerEndpointType{
		"clientapi": "ClientApi",
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// it could be a new value - best effort convert this
	v := input

	out := VideoAnalyzerEndpointType(v)
	return &out, nil
}

type VideoType string

const (
	VideoTypeArchive VideoType = "Archive"
)

func PossibleValuesForVideoType() []string {
	return []string{
		"Archive",
	}
}

func parseVideoType(input string) (*VideoType, error) {
	vals := map[string]VideoType{
		"archive": "Archive",
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// it could be a new value - best effort convert this
	v := input

	out := VideoType(v)
	return &out, nil
}
