package customdomains

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

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

func (s *CertificateSource) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseCertificateSource(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
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

type CertificateSourceParametersType string

const (
	CertificateSourceParametersTypeCdnCertificateSourceParameters      CertificateSourceParametersType = "CdnCertificateSourceParameters"
	CertificateSourceParametersTypeKeyVaultCertificateSourceParameters CertificateSourceParametersType = "KeyVaultCertificateSourceParameters"
)

func PossibleValuesForCertificateSourceParametersType() []string {
	return []string{
		string(CertificateSourceParametersTypeCdnCertificateSourceParameters),
		string(CertificateSourceParametersTypeKeyVaultCertificateSourceParameters),
	}
}

func (s *CertificateSourceParametersType) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseCertificateSourceParametersType(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseCertificateSourceParametersType(input string) (*CertificateSourceParametersType, error) {
	vals := map[string]CertificateSourceParametersType{
		"cdncertificatesourceparameters":      CertificateSourceParametersTypeCdnCertificateSourceParameters,
		"keyvaultcertificatesourceparameters": CertificateSourceParametersTypeKeyVaultCertificateSourceParameters,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := CertificateSourceParametersType(input)
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

func (s *CertificateType) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseCertificateType(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
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

func (s *CustomDomainResourceState) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseCustomDomainResourceState(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
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

type CustomHTTPSProvisioningState string

const (
	CustomHTTPSProvisioningStateDisabled  CustomHTTPSProvisioningState = "Disabled"
	CustomHTTPSProvisioningStateDisabling CustomHTTPSProvisioningState = "Disabling"
	CustomHTTPSProvisioningStateEnabled   CustomHTTPSProvisioningState = "Enabled"
	CustomHTTPSProvisioningStateEnabling  CustomHTTPSProvisioningState = "Enabling"
	CustomHTTPSProvisioningStateFailed    CustomHTTPSProvisioningState = "Failed"
)

func PossibleValuesForCustomHTTPSProvisioningState() []string {
	return []string{
		string(CustomHTTPSProvisioningStateDisabled),
		string(CustomHTTPSProvisioningStateDisabling),
		string(CustomHTTPSProvisioningStateEnabled),
		string(CustomHTTPSProvisioningStateEnabling),
		string(CustomHTTPSProvisioningStateFailed),
	}
}

func (s *CustomHTTPSProvisioningState) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseCustomHTTPSProvisioningState(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseCustomHTTPSProvisioningState(input string) (*CustomHTTPSProvisioningState, error) {
	vals := map[string]CustomHTTPSProvisioningState{
		"disabled":  CustomHTTPSProvisioningStateDisabled,
		"disabling": CustomHTTPSProvisioningStateDisabling,
		"enabled":   CustomHTTPSProvisioningStateEnabled,
		"enabling":  CustomHTTPSProvisioningStateEnabling,
		"failed":    CustomHTTPSProvisioningStateFailed,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := CustomHTTPSProvisioningState(input)
	return &out, nil
}

type CustomHTTPSProvisioningSubstate string

const (
	CustomHTTPSProvisioningSubstateCertificateDeleted                            CustomHTTPSProvisioningSubstate = "CertificateDeleted"
	CustomHTTPSProvisioningSubstateCertificateDeployed                           CustomHTTPSProvisioningSubstate = "CertificateDeployed"
	CustomHTTPSProvisioningSubstateDeletingCertificate                           CustomHTTPSProvisioningSubstate = "DeletingCertificate"
	CustomHTTPSProvisioningSubstateDeployingCertificate                          CustomHTTPSProvisioningSubstate = "DeployingCertificate"
	CustomHTTPSProvisioningSubstateDomainControlValidationRequestApproved        CustomHTTPSProvisioningSubstate = "DomainControlValidationRequestApproved"
	CustomHTTPSProvisioningSubstateDomainControlValidationRequestRejected        CustomHTTPSProvisioningSubstate = "DomainControlValidationRequestRejected"
	CustomHTTPSProvisioningSubstateDomainControlValidationRequestTimedOut        CustomHTTPSProvisioningSubstate = "DomainControlValidationRequestTimedOut"
	CustomHTTPSProvisioningSubstateIssuingCertificate                            CustomHTTPSProvisioningSubstate = "IssuingCertificate"
	CustomHTTPSProvisioningSubstatePendingDomainControlValidationREquestApproval CustomHTTPSProvisioningSubstate = "PendingDomainControlValidationREquestApproval"
	CustomHTTPSProvisioningSubstateSubmittingDomainControlValidationRequest      CustomHTTPSProvisioningSubstate = "SubmittingDomainControlValidationRequest"
)

func PossibleValuesForCustomHTTPSProvisioningSubstate() []string {
	return []string{
		string(CustomHTTPSProvisioningSubstateCertificateDeleted),
		string(CustomHTTPSProvisioningSubstateCertificateDeployed),
		string(CustomHTTPSProvisioningSubstateDeletingCertificate),
		string(CustomHTTPSProvisioningSubstateDeployingCertificate),
		string(CustomHTTPSProvisioningSubstateDomainControlValidationRequestApproved),
		string(CustomHTTPSProvisioningSubstateDomainControlValidationRequestRejected),
		string(CustomHTTPSProvisioningSubstateDomainControlValidationRequestTimedOut),
		string(CustomHTTPSProvisioningSubstateIssuingCertificate),
		string(CustomHTTPSProvisioningSubstatePendingDomainControlValidationREquestApproval),
		string(CustomHTTPSProvisioningSubstateSubmittingDomainControlValidationRequest),
	}
}

func (s *CustomHTTPSProvisioningSubstate) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseCustomHTTPSProvisioningSubstate(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseCustomHTTPSProvisioningSubstate(input string) (*CustomHTTPSProvisioningSubstate, error) {
	vals := map[string]CustomHTTPSProvisioningSubstate{
		"certificatedeleted":                            CustomHTTPSProvisioningSubstateCertificateDeleted,
		"certificatedeployed":                           CustomHTTPSProvisioningSubstateCertificateDeployed,
		"deletingcertificate":                           CustomHTTPSProvisioningSubstateDeletingCertificate,
		"deployingcertificate":                          CustomHTTPSProvisioningSubstateDeployingCertificate,
		"domaincontrolvalidationrequestapproved":        CustomHTTPSProvisioningSubstateDomainControlValidationRequestApproved,
		"domaincontrolvalidationrequestrejected":        CustomHTTPSProvisioningSubstateDomainControlValidationRequestRejected,
		"domaincontrolvalidationrequesttimedout":        CustomHTTPSProvisioningSubstateDomainControlValidationRequestTimedOut,
		"issuingcertificate":                            CustomHTTPSProvisioningSubstateIssuingCertificate,
		"pendingdomaincontrolvalidationrequestapproval": CustomHTTPSProvisioningSubstatePendingDomainControlValidationREquestApproval,
		"submittingdomaincontrolvalidationrequest":      CustomHTTPSProvisioningSubstateSubmittingDomainControlValidationRequest,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := CustomHTTPSProvisioningSubstate(input)
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

func (s *DeleteRule) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseDeleteRule(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
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

func (s *MinimumTlsVersion) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseMinimumTlsVersion(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
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

func (s *ProtocolType) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseProtocolType(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
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

type UpdateRule string

const (
	UpdateRuleNoAction UpdateRule = "NoAction"
)

func PossibleValuesForUpdateRule() []string {
	return []string{
		string(UpdateRuleNoAction),
	}
}

func (s *UpdateRule) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseUpdateRule(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
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
