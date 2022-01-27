package customdomains

import "strings"

type CertificateSource string

const (
	CertificateSourceAzureKeyVault CertificateSource = "AzureKeyVault"
	CertificateSourceCdn           CertificateSource = "Cdn"
)

func PossibleValuesForCertificateSource() []string {
	return []string{
		string(CertificateSourceAzureKeyVault),
		string(CertificateSourceCdn),
	}
}

func parseCertificateSource(input string) (*CertificateSource, error) {
	vals := map[string]CertificateSource{
		"azurekeyvault": CertificateSourceAzureKeyVault,
		"cdn":           CertificateSourceCdn,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := CertificateSource(input)
	return &out, nil
}

type CertificateType string

const (
	CertificateTypeDedicated CertificateType = "Dedicated"
	CertificateTypeShared    CertificateType = "Shared"
)

func PossibleValuesForCertificateType() []string {
	return []string{
		string(CertificateTypeDedicated),
		string(CertificateTypeShared),
	}
}

func parseCertificateType(input string) (*CertificateType, error) {
	vals := map[string]CertificateType{
		"dedicated": CertificateTypeDedicated,
		"shared":    CertificateTypeShared,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := CertificateType(input)
	return &out, nil
}

type CustomDomainResourceState string

const (
	CustomDomainResourceStateActive   CustomDomainResourceState = "Active"
	CustomDomainResourceStateCreating CustomDomainResourceState = "Creating"
	CustomDomainResourceStateDeleting CustomDomainResourceState = "Deleting"
)

func PossibleValuesForCustomDomainResourceState() []string {
	return []string{
		string(CustomDomainResourceStateActive),
		string(CustomDomainResourceStateCreating),
		string(CustomDomainResourceStateDeleting),
	}
}

func parseCustomDomainResourceState(input string) (*CustomDomainResourceState, error) {
	vals := map[string]CustomDomainResourceState{
		"active":   CustomDomainResourceStateActive,
		"creating": CustomDomainResourceStateCreating,
		"deleting": CustomDomainResourceStateDeleting,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := CustomDomainResourceState(input)
	return &out, nil
}

type CustomHttpsProvisioningState string

const (
	CustomHttpsProvisioningStateDisabled  CustomHttpsProvisioningState = "Disabled"
	CustomHttpsProvisioningStateDisabling CustomHttpsProvisioningState = "Disabling"
	CustomHttpsProvisioningStateEnabled   CustomHttpsProvisioningState = "Enabled"
	CustomHttpsProvisioningStateEnabling  CustomHttpsProvisioningState = "Enabling"
	CustomHttpsProvisioningStateFailed    CustomHttpsProvisioningState = "Failed"
)

func PossibleValuesForCustomHttpsProvisioningState() []string {
	return []string{
		string(CustomHttpsProvisioningStateDisabled),
		string(CustomHttpsProvisioningStateDisabling),
		string(CustomHttpsProvisioningStateEnabled),
		string(CustomHttpsProvisioningStateEnabling),
		string(CustomHttpsProvisioningStateFailed),
	}
}

func parseCustomHttpsProvisioningState(input string) (*CustomHttpsProvisioningState, error) {
	vals := map[string]CustomHttpsProvisioningState{
		"disabled":  CustomHttpsProvisioningStateDisabled,
		"disabling": CustomHttpsProvisioningStateDisabling,
		"enabled":   CustomHttpsProvisioningStateEnabled,
		"enabling":  CustomHttpsProvisioningStateEnabling,
		"failed":    CustomHttpsProvisioningStateFailed,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := CustomHttpsProvisioningState(input)
	return &out, nil
}

type CustomHttpsProvisioningSubstate string

const (
	CustomHttpsProvisioningSubstateCertificateDeleted                            CustomHttpsProvisioningSubstate = "CertificateDeleted"
	CustomHttpsProvisioningSubstateCertificateDeployed                           CustomHttpsProvisioningSubstate = "CertificateDeployed"
	CustomHttpsProvisioningSubstateDeletingCertificate                           CustomHttpsProvisioningSubstate = "DeletingCertificate"
	CustomHttpsProvisioningSubstateDeployingCertificate                          CustomHttpsProvisioningSubstate = "DeployingCertificate"
	CustomHttpsProvisioningSubstateDomainControlValidationRequestApproved        CustomHttpsProvisioningSubstate = "DomainControlValidationRequestApproved"
	CustomHttpsProvisioningSubstateDomainControlValidationRequestRejected        CustomHttpsProvisioningSubstate = "DomainControlValidationRequestRejected"
	CustomHttpsProvisioningSubstateDomainControlValidationRequestTimedOut        CustomHttpsProvisioningSubstate = "DomainControlValidationRequestTimedOut"
	CustomHttpsProvisioningSubstateIssuingCertificate                            CustomHttpsProvisioningSubstate = "IssuingCertificate"
	CustomHttpsProvisioningSubstatePendingDomainControlValidationREquestApproval CustomHttpsProvisioningSubstate = "PendingDomainControlValidationREquestApproval"
	CustomHttpsProvisioningSubstateSubmittingDomainControlValidationRequest      CustomHttpsProvisioningSubstate = "SubmittingDomainControlValidationRequest"
)

func PossibleValuesForCustomHttpsProvisioningSubstate() []string {
	return []string{
		string(CustomHttpsProvisioningSubstateCertificateDeleted),
		string(CustomHttpsProvisioningSubstateCertificateDeployed),
		string(CustomHttpsProvisioningSubstateDeletingCertificate),
		string(CustomHttpsProvisioningSubstateDeployingCertificate),
		string(CustomHttpsProvisioningSubstateDomainControlValidationRequestApproved),
		string(CustomHttpsProvisioningSubstateDomainControlValidationRequestRejected),
		string(CustomHttpsProvisioningSubstateDomainControlValidationRequestTimedOut),
		string(CustomHttpsProvisioningSubstateIssuingCertificate),
		string(CustomHttpsProvisioningSubstatePendingDomainControlValidationREquestApproval),
		string(CustomHttpsProvisioningSubstateSubmittingDomainControlValidationRequest),
	}
}

func parseCustomHttpsProvisioningSubstate(input string) (*CustomHttpsProvisioningSubstate, error) {
	vals := map[string]CustomHttpsProvisioningSubstate{
		"certificatedeleted":                            CustomHttpsProvisioningSubstateCertificateDeleted,
		"certificatedeployed":                           CustomHttpsProvisioningSubstateCertificateDeployed,
		"deletingcertificate":                           CustomHttpsProvisioningSubstateDeletingCertificate,
		"deployingcertificate":                          CustomHttpsProvisioningSubstateDeployingCertificate,
		"domaincontrolvalidationrequestapproved":        CustomHttpsProvisioningSubstateDomainControlValidationRequestApproved,
		"domaincontrolvalidationrequestrejected":        CustomHttpsProvisioningSubstateDomainControlValidationRequestRejected,
		"domaincontrolvalidationrequesttimedout":        CustomHttpsProvisioningSubstateDomainControlValidationRequestTimedOut,
		"issuingcertificate":                            CustomHttpsProvisioningSubstateIssuingCertificate,
		"pendingdomaincontrolvalidationrequestapproval": CustomHttpsProvisioningSubstatePendingDomainControlValidationREquestApproval,
		"submittingdomaincontrolvalidationrequest":      CustomHttpsProvisioningSubstateSubmittingDomainControlValidationRequest,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := CustomHttpsProvisioningSubstate(input)
	return &out, nil
}

type DeleteRule string

const (
	DeleteRuleNoAction DeleteRule = "NoAction"
)

func PossibleValuesForDeleteRule() []string {
	return []string{
		string(DeleteRuleNoAction),
	}
}

func parseDeleteRule(input string) (*DeleteRule, error) {
	vals := map[string]DeleteRule{
		"noaction": DeleteRuleNoAction,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := DeleteRule(input)
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

type MinimumTlsVersion string

const (
	MinimumTlsVersionNone       MinimumTlsVersion = "None"
	MinimumTlsVersionTLSOneTwo  MinimumTlsVersion = "TLS12"
	MinimumTlsVersionTLSOneZero MinimumTlsVersion = "TLS10"
)

func PossibleValuesForMinimumTlsVersion() []string {
	return []string{
		string(MinimumTlsVersionNone),
		string(MinimumTlsVersionTLSOneTwo),
		string(MinimumTlsVersionTLSOneZero),
	}
}

func parseMinimumTlsVersion(input string) (*MinimumTlsVersion, error) {
	vals := map[string]MinimumTlsVersion{
		"none":  MinimumTlsVersionNone,
		"tls12": MinimumTlsVersionTLSOneTwo,
		"tls10": MinimumTlsVersionTLSOneZero,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := MinimumTlsVersion(input)
	return &out, nil
}

type ProtocolType string

const (
	ProtocolTypeIPBased              ProtocolType = "IPBased"
	ProtocolTypeServerNameIndication ProtocolType = "ServerNameIndication"
)

func PossibleValuesForProtocolType() []string {
	return []string{
		string(ProtocolTypeIPBased),
		string(ProtocolTypeServerNameIndication),
	}
}

func parseProtocolType(input string) (*ProtocolType, error) {
	vals := map[string]ProtocolType{
		"ipbased":              ProtocolTypeIPBased,
		"servernameindication": ProtocolTypeServerNameIndication,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := ProtocolType(input)
	return &out, nil
}

type TypeName string

const (
	TypeNameCdnCertificateSourceParameters TypeName = "CdnCertificateSourceParameters"
)

func PossibleValuesForTypeName() []string {
	return []string{
		string(TypeNameCdnCertificateSourceParameters),
	}
}

func parseTypeName(input string) (*TypeName, error) {
	vals := map[string]TypeName{
		"cdncertificatesourceparameters": TypeNameCdnCertificateSourceParameters,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := TypeName(input)
	return &out, nil
}

type UpdateRule string

const (
	UpdateRuleNoAction UpdateRule = "NoAction"
)

func PossibleValuesForUpdateRule() []string {
	return []string{
		string(UpdateRuleNoAction),
	}
}

func parseUpdateRule(input string) (*UpdateRule, error) {
	vals := map[string]UpdateRule{
		"noaction": UpdateRuleNoAction,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := UpdateRule(input)
	return &out, nil
}
