package secrets

import "strings"

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

type SecretType string

const (
	SecretTypeAzureFirstPartyManagedCertificate SecretType = "AzureFirstPartyManagedCertificate"
	SecretTypeCustomerCertificate               SecretType = "CustomerCertificate"
	SecretTypeManagedCertificate                SecretType = "ManagedCertificate"
	SecretTypeUrlSigningKey                     SecretType = "UrlSigningKey"
)

func PossibleValuesForSecretType() []string {
	return []string{
		string(SecretTypeAzureFirstPartyManagedCertificate),
		string(SecretTypeCustomerCertificate),
		string(SecretTypeManagedCertificate),
		string(SecretTypeUrlSigningKey),
	}
}

func parseSecretType(input string) (*SecretType, error) {
	vals := map[string]SecretType{
		"azurefirstpartymanagedcertificate": SecretTypeAzureFirstPartyManagedCertificate,
		"customercertificate":               SecretTypeCustomerCertificate,
		"managedcertificate":                SecretTypeManagedCertificate,
		"urlsigningkey":                     SecretTypeUrlSigningKey,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := SecretType(input)
	return &out, nil
}

type Status string

const (
	StatusAccessDenied       Status = "AccessDenied"
	StatusCertificateExpired Status = "CertificateExpired"
	StatusInvalid            Status = "Invalid"
	StatusValid              Status = "Valid"
)

func PossibleValuesForStatus() []string {
	return []string{
		string(StatusAccessDenied),
		string(StatusCertificateExpired),
		string(StatusInvalid),
		string(StatusValid),
	}
}

func parseStatus(input string) (*Status, error) {
	vals := map[string]Status{
		"accessdenied":       StatusAccessDenied,
		"certificateexpired": StatusCertificateExpired,
		"invalid":            StatusInvalid,
		"valid":              StatusValid,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := Status(input)
	return &out, nil
}
