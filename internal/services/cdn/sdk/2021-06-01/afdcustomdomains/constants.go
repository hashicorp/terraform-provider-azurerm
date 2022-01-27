package afdcustomdomains

import "strings"

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

type AfdMinimumTlsVersion string

const (
	AfdMinimumTlsVersionTLSOneTwo  AfdMinimumTlsVersion = "TLS12"
	AfdMinimumTlsVersionTLSOneZero AfdMinimumTlsVersion = "TLS10"
)

func PossibleValuesForAfdMinimumTlsVersion() []string {
	return []string{
		string(AfdMinimumTlsVersionTLSOneTwo),
		string(AfdMinimumTlsVersionTLSOneZero),
	}
}

func parseAfdMinimumTlsVersion(input string) (*AfdMinimumTlsVersion, error) {
	vals := map[string]AfdMinimumTlsVersion{
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
	DomainValidationStateApproved            DomainValidationState = "Approved"
	DomainValidationStatePending             DomainValidationState = "Pending"
	DomainValidationStatePendingRevalidation DomainValidationState = "PendingRevalidation"
	DomainValidationStateRejected            DomainValidationState = "Rejected"
	DomainValidationStateSubmitting          DomainValidationState = "Submitting"
	DomainValidationStateTimedOut            DomainValidationState = "TimedOut"
	DomainValidationStateUnknown             DomainValidationState = "Unknown"
)

func PossibleValuesForDomainValidationState() []string {
	return []string{
		string(DomainValidationStateApproved),
		string(DomainValidationStatePending),
		string(DomainValidationStatePendingRevalidation),
		string(DomainValidationStateRejected),
		string(DomainValidationStateSubmitting),
		string(DomainValidationStateTimedOut),
		string(DomainValidationStateUnknown),
	}
}

func parseDomainValidationState(input string) (*DomainValidationState, error) {
	vals := map[string]DomainValidationState{
		"approved":            DomainValidationStateApproved,
		"pending":             DomainValidationStatePending,
		"pendingrevalidation": DomainValidationStatePendingRevalidation,
		"rejected":            DomainValidationStateRejected,
		"submitting":          DomainValidationStateSubmitting,
		"timedout":            DomainValidationStateTimedOut,
		"unknown":             DomainValidationStateUnknown,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := DomainValidationState(input)
	return &out, nil
}

type IdentityType string

const (
	IdentityTypeApplication     IdentityType = "application"
	IdentityTypeKey             IdentityType = "key"
	IdentityTypeManagedIdentity IdentityType = "managedIdentity"
	IdentityTypeUser            IdentityType = "user"
)

func PossibleValuesForIdentityType() []string {
	return []string{
		string(IdentityTypeApplication),
		string(IdentityTypeKey),
		string(IdentityTypeManagedIdentity),
		string(IdentityTypeUser),
	}
}

func parseIdentityType(input string) (*IdentityType, error) {
	vals := map[string]IdentityType{
		"application":     IdentityTypeApplication,
		"key":             IdentityTypeKey,
		"managedidentity": IdentityTypeManagedIdentity,
		"user":            IdentityTypeUser,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := IdentityType(input)
	return &out, nil
}
