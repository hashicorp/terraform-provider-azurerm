package afdcustomdomains

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type AfdCertificateType string

const (
	AfdCertificateTypeAzureFirstPartyManagedCertificate AfdCertificateType = "AzureFirstPartyManagedCertificate"
	AfdCertificateTypeCustomerCertificate               AfdCertificateType = "CustomerCertificate"
	AfdCertificateTypeManagedCertificate                AfdCertificateType = "ManagedCertificate"
)

func PossibleValuesForAfdCertificateType() []string {
	return []string{
		string(AfdCertificateTypeAzureFirstPartyManagedCertificate),
		string(AfdCertificateTypeCustomerCertificate),
		string(AfdCertificateTypeManagedCertificate),
	}
}

func (s *AfdCertificateType) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseAfdCertificateType(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseAfdCertificateType(input string) (*AfdCertificateType, error) {
	vals := map[string]AfdCertificateType{
		"azurefirstpartymanagedcertificate": AfdCertificateTypeAzureFirstPartyManagedCertificate,
		"customercertificate":               AfdCertificateTypeCustomerCertificate,
		"managedcertificate":                AfdCertificateTypeManagedCertificate,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := AfdCertificateType(input)
	return &out, nil
}

type AfdCipherSuiteSetType string

const (
	AfdCipherSuiteSetTypeCustomized               AfdCipherSuiteSetType = "Customized"
	AfdCipherSuiteSetTypeTLSOneTwoTwoZeroTwoThree AfdCipherSuiteSetType = "TLS12_2023"
	AfdCipherSuiteSetTypeTLSOneTwoTwoZeroTwoTwo   AfdCipherSuiteSetType = "TLS12_2022"
	AfdCipherSuiteSetTypeTLSOneZeroTwoZeroOneNine AfdCipherSuiteSetType = "TLS10_2019"
)

func PossibleValuesForAfdCipherSuiteSetType() []string {
	return []string{
		string(AfdCipherSuiteSetTypeCustomized),
		string(AfdCipherSuiteSetTypeTLSOneTwoTwoZeroTwoThree),
		string(AfdCipherSuiteSetTypeTLSOneTwoTwoZeroTwoTwo),
		string(AfdCipherSuiteSetTypeTLSOneZeroTwoZeroOneNine),
	}
}

func (s *AfdCipherSuiteSetType) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseAfdCipherSuiteSetType(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseAfdCipherSuiteSetType(input string) (*AfdCipherSuiteSetType, error) {
	vals := map[string]AfdCipherSuiteSetType{
		"customized": AfdCipherSuiteSetTypeCustomized,
		"tls12_2023": AfdCipherSuiteSetTypeTLSOneTwoTwoZeroTwoThree,
		"tls12_2022": AfdCipherSuiteSetTypeTLSOneTwoTwoZeroTwoTwo,
		"tls10_2019": AfdCipherSuiteSetTypeTLSOneZeroTwoZeroOneNine,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := AfdCipherSuiteSetType(input)
	return &out, nil
}

type AfdCustomizedCipherSuiteForTls12 string

const (
	AfdCustomizedCipherSuiteForTls12DHERSAAESOneTwoEightGCMSHATwoFiveSix      AfdCustomizedCipherSuiteForTls12 = "DHE_RSA_AES128_GCM_SHA256"
	AfdCustomizedCipherSuiteForTls12DHERSAAESTwoFiveSixGCMSHAThreeEightFour   AfdCustomizedCipherSuiteForTls12 = "DHE_RSA_AES256_GCM_SHA384"
	AfdCustomizedCipherSuiteForTls12ECDHERSAAESOneTwoEightGCMSHATwoFiveSix    AfdCustomizedCipherSuiteForTls12 = "ECDHE_RSA_AES128_GCM_SHA256"
	AfdCustomizedCipherSuiteForTls12ECDHERSAAESOneTwoEightSHATwoFiveSix       AfdCustomizedCipherSuiteForTls12 = "ECDHE_RSA_AES128_SHA256"
	AfdCustomizedCipherSuiteForTls12ECDHERSAAESTwoFiveSixGCMSHAThreeEightFour AfdCustomizedCipherSuiteForTls12 = "ECDHE_RSA_AES256_GCM_SHA384"
	AfdCustomizedCipherSuiteForTls12ECDHERSAAESTwoFiveSixSHAThreeEightFour    AfdCustomizedCipherSuiteForTls12 = "ECDHE_RSA_AES256_SHA384"
)

func PossibleValuesForAfdCustomizedCipherSuiteForTls12() []string {
	return []string{
		string(AfdCustomizedCipherSuiteForTls12DHERSAAESOneTwoEightGCMSHATwoFiveSix),
		string(AfdCustomizedCipherSuiteForTls12DHERSAAESTwoFiveSixGCMSHAThreeEightFour),
		string(AfdCustomizedCipherSuiteForTls12ECDHERSAAESOneTwoEightGCMSHATwoFiveSix),
		string(AfdCustomizedCipherSuiteForTls12ECDHERSAAESOneTwoEightSHATwoFiveSix),
		string(AfdCustomizedCipherSuiteForTls12ECDHERSAAESTwoFiveSixGCMSHAThreeEightFour),
		string(AfdCustomizedCipherSuiteForTls12ECDHERSAAESTwoFiveSixSHAThreeEightFour),
	}
}

func (s *AfdCustomizedCipherSuiteForTls12) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseAfdCustomizedCipherSuiteForTls12(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseAfdCustomizedCipherSuiteForTls12(input string) (*AfdCustomizedCipherSuiteForTls12, error) {
	vals := map[string]AfdCustomizedCipherSuiteForTls12{
		"dhe_rsa_aes128_gcm_sha256":   AfdCustomizedCipherSuiteForTls12DHERSAAESOneTwoEightGCMSHATwoFiveSix,
		"dhe_rsa_aes256_gcm_sha384":   AfdCustomizedCipherSuiteForTls12DHERSAAESTwoFiveSixGCMSHAThreeEightFour,
		"ecdhe_rsa_aes128_gcm_sha256": AfdCustomizedCipherSuiteForTls12ECDHERSAAESOneTwoEightGCMSHATwoFiveSix,
		"ecdhe_rsa_aes128_sha256":     AfdCustomizedCipherSuiteForTls12ECDHERSAAESOneTwoEightSHATwoFiveSix,
		"ecdhe_rsa_aes256_gcm_sha384": AfdCustomizedCipherSuiteForTls12ECDHERSAAESTwoFiveSixGCMSHAThreeEightFour,
		"ecdhe_rsa_aes256_sha384":     AfdCustomizedCipherSuiteForTls12ECDHERSAAESTwoFiveSixSHAThreeEightFour,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := AfdCustomizedCipherSuiteForTls12(input)
	return &out, nil
}

type AfdCustomizedCipherSuiteForTls13 string

const (
	AfdCustomizedCipherSuiteForTls13TLSAESOneTwoEightGCMSHATwoFiveSix    AfdCustomizedCipherSuiteForTls13 = "TLS_AES_128_GCM_SHA256"
	AfdCustomizedCipherSuiteForTls13TLSAESTwoFiveSixGCMSHAThreeEightFour AfdCustomizedCipherSuiteForTls13 = "TLS_AES_256_GCM_SHA384"
)

func PossibleValuesForAfdCustomizedCipherSuiteForTls13() []string {
	return []string{
		string(AfdCustomizedCipherSuiteForTls13TLSAESOneTwoEightGCMSHATwoFiveSix),
		string(AfdCustomizedCipherSuiteForTls13TLSAESTwoFiveSixGCMSHAThreeEightFour),
	}
}

func (s *AfdCustomizedCipherSuiteForTls13) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseAfdCustomizedCipherSuiteForTls13(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseAfdCustomizedCipherSuiteForTls13(input string) (*AfdCustomizedCipherSuiteForTls13, error) {
	vals := map[string]AfdCustomizedCipherSuiteForTls13{
		"tls_aes_128_gcm_sha256": AfdCustomizedCipherSuiteForTls13TLSAESOneTwoEightGCMSHATwoFiveSix,
		"tls_aes_256_gcm_sha384": AfdCustomizedCipherSuiteForTls13TLSAESTwoFiveSixGCMSHAThreeEightFour,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := AfdCustomizedCipherSuiteForTls13(input)
	return &out, nil
}

type AfdMinimumTlsVersion string

const (
	AfdMinimumTlsVersionTLSOneThree AfdMinimumTlsVersion = "TLS13"
	AfdMinimumTlsVersionTLSOneTwo   AfdMinimumTlsVersion = "TLS12"
	AfdMinimumTlsVersionTLSOneZero  AfdMinimumTlsVersion = "TLS10"
)

func PossibleValuesForAfdMinimumTlsVersion() []string {
	return []string{
		string(AfdMinimumTlsVersionTLSOneThree),
		string(AfdMinimumTlsVersionTLSOneTwo),
		string(AfdMinimumTlsVersionTLSOneZero),
	}
}

func (s *AfdMinimumTlsVersion) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseAfdMinimumTlsVersion(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseAfdMinimumTlsVersion(input string) (*AfdMinimumTlsVersion, error) {
	vals := map[string]AfdMinimumTlsVersion{
		"tls13": AfdMinimumTlsVersionTLSOneThree,
		"tls12": AfdMinimumTlsVersionTLSOneTwo,
		"tls10": AfdMinimumTlsVersionTLSOneZero,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := AfdMinimumTlsVersion(input)
	return &out, nil
}

type AfdProvisioningState string

const (
	AfdProvisioningStateCreating  AfdProvisioningState = "Creating"
	AfdProvisioningStateDeleting  AfdProvisioningState = "Deleting"
	AfdProvisioningStateFailed    AfdProvisioningState = "Failed"
	AfdProvisioningStateSucceeded AfdProvisioningState = "Succeeded"
	AfdProvisioningStateUpdating  AfdProvisioningState = "Updating"
)

func PossibleValuesForAfdProvisioningState() []string {
	return []string{
		string(AfdProvisioningStateCreating),
		string(AfdProvisioningStateDeleting),
		string(AfdProvisioningStateFailed),
		string(AfdProvisioningStateSucceeded),
		string(AfdProvisioningStateUpdating),
	}
}

func (s *AfdProvisioningState) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseAfdProvisioningState(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseAfdProvisioningState(input string) (*AfdProvisioningState, error) {
	vals := map[string]AfdProvisioningState{
		"creating":  AfdProvisioningStateCreating,
		"deleting":  AfdProvisioningStateDeleting,
		"failed":    AfdProvisioningStateFailed,
		"succeeded": AfdProvisioningStateSucceeded,
		"updating":  AfdProvisioningStateUpdating,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := AfdProvisioningState(input)
	return &out, nil
}

type DeploymentStatus string

const (
	DeploymentStatusFailed     DeploymentStatus = "Failed"
	DeploymentStatusInProgress DeploymentStatus = "InProgress"
	DeploymentStatusNotStarted DeploymentStatus = "NotStarted"
	DeploymentStatusSucceeded  DeploymentStatus = "Succeeded"
)

func PossibleValuesForDeploymentStatus() []string {
	return []string{
		string(DeploymentStatusFailed),
		string(DeploymentStatusInProgress),
		string(DeploymentStatusNotStarted),
		string(DeploymentStatusSucceeded),
	}
}

func (s *DeploymentStatus) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseDeploymentStatus(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseDeploymentStatus(input string) (*DeploymentStatus, error) {
	vals := map[string]DeploymentStatus{
		"failed":     DeploymentStatusFailed,
		"inprogress": DeploymentStatusInProgress,
		"notstarted": DeploymentStatusNotStarted,
		"succeeded":  DeploymentStatusSucceeded,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := DeploymentStatus(input)
	return &out, nil
}

type DomainValidationState string

const (
	DomainValidationStateApproved                  DomainValidationState = "Approved"
	DomainValidationStateInternalError             DomainValidationState = "InternalError"
	DomainValidationStatePending                   DomainValidationState = "Pending"
	DomainValidationStatePendingRevalidation       DomainValidationState = "PendingRevalidation"
	DomainValidationStateRefreshingValidationToken DomainValidationState = "RefreshingValidationToken"
	DomainValidationStateRejected                  DomainValidationState = "Rejected"
	DomainValidationStateSubmitting                DomainValidationState = "Submitting"
	DomainValidationStateTimedOut                  DomainValidationState = "TimedOut"
	DomainValidationStateUnknown                   DomainValidationState = "Unknown"
)

func PossibleValuesForDomainValidationState() []string {
	return []string{
		string(DomainValidationStateApproved),
		string(DomainValidationStateInternalError),
		string(DomainValidationStatePending),
		string(DomainValidationStatePendingRevalidation),
		string(DomainValidationStateRefreshingValidationToken),
		string(DomainValidationStateRejected),
		string(DomainValidationStateSubmitting),
		string(DomainValidationStateTimedOut),
		string(DomainValidationStateUnknown),
	}
}

func (s *DomainValidationState) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseDomainValidationState(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseDomainValidationState(input string) (*DomainValidationState, error) {
	vals := map[string]DomainValidationState{
		"approved":                  DomainValidationStateApproved,
		"internalerror":             DomainValidationStateInternalError,
		"pending":                   DomainValidationStatePending,
		"pendingrevalidation":       DomainValidationStatePendingRevalidation,
		"refreshingvalidationtoken": DomainValidationStateRefreshingValidationToken,
		"rejected":                  DomainValidationStateRejected,
		"submitting":                DomainValidationStateSubmitting,
		"timedout":                  DomainValidationStateTimedOut,
		"unknown":                   DomainValidationStateUnknown,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := DomainValidationState(input)
	return &out, nil
}
