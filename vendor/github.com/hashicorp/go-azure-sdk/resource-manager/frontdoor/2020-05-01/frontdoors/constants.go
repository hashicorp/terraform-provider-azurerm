package frontdoors

import "strings"

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type BackendEnabledState string

const (
	BackendEnabledStateDisabled BackendEnabledState = "Disabled"
	BackendEnabledStateEnabled  BackendEnabledState = "Enabled"
)

func PossibleValuesForBackendEnabledState() []string {
	return []string{
		string(BackendEnabledStateDisabled),
		string(BackendEnabledStateEnabled),
	}
}

func parseBackendEnabledState(input string) (*BackendEnabledState, error) {
	vals := map[string]BackendEnabledState{
		"disabled": BackendEnabledStateDisabled,
		"enabled":  BackendEnabledStateEnabled,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := BackendEnabledState(input)
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

type DynamicCompressionEnabled string

const (
	DynamicCompressionEnabledDisabled DynamicCompressionEnabled = "Disabled"
	DynamicCompressionEnabledEnabled  DynamicCompressionEnabled = "Enabled"
)

func PossibleValuesForDynamicCompressionEnabled() []string {
	return []string{
		string(DynamicCompressionEnabledDisabled),
		string(DynamicCompressionEnabledEnabled),
	}
}

func parseDynamicCompressionEnabled(input string) (*DynamicCompressionEnabled, error) {
	vals := map[string]DynamicCompressionEnabled{
		"disabled": DynamicCompressionEnabledDisabled,
		"enabled":  DynamicCompressionEnabledEnabled,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := DynamicCompressionEnabled(input)
	return &out, nil
}

type EnforceCertificateNameCheckEnabledState string

const (
	EnforceCertificateNameCheckEnabledStateDisabled EnforceCertificateNameCheckEnabledState = "Disabled"
	EnforceCertificateNameCheckEnabledStateEnabled  EnforceCertificateNameCheckEnabledState = "Enabled"
)

func PossibleValuesForEnforceCertificateNameCheckEnabledState() []string {
	return []string{
		string(EnforceCertificateNameCheckEnabledStateDisabled),
		string(EnforceCertificateNameCheckEnabledStateEnabled),
	}
}

func parseEnforceCertificateNameCheckEnabledState(input string) (*EnforceCertificateNameCheckEnabledState, error) {
	vals := map[string]EnforceCertificateNameCheckEnabledState{
		"disabled": EnforceCertificateNameCheckEnabledStateDisabled,
		"enabled":  EnforceCertificateNameCheckEnabledStateEnabled,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := EnforceCertificateNameCheckEnabledState(input)
	return &out, nil
}

type FrontDoorCertificateSource string

const (
	FrontDoorCertificateSourceAzureKeyVault FrontDoorCertificateSource = "AzureKeyVault"
	FrontDoorCertificateSourceFrontDoor     FrontDoorCertificateSource = "FrontDoor"
)

func PossibleValuesForFrontDoorCertificateSource() []string {
	return []string{
		string(FrontDoorCertificateSourceAzureKeyVault),
		string(FrontDoorCertificateSourceFrontDoor),
	}
}

func parseFrontDoorCertificateSource(input string) (*FrontDoorCertificateSource, error) {
	vals := map[string]FrontDoorCertificateSource{
		"azurekeyvault": FrontDoorCertificateSourceAzureKeyVault,
		"frontdoor":     FrontDoorCertificateSourceFrontDoor,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := FrontDoorCertificateSource(input)
	return &out, nil
}

type FrontDoorCertificateType string

const (
	FrontDoorCertificateTypeDedicated FrontDoorCertificateType = "Dedicated"
)

func PossibleValuesForFrontDoorCertificateType() []string {
	return []string{
		string(FrontDoorCertificateTypeDedicated),
	}
}

func parseFrontDoorCertificateType(input string) (*FrontDoorCertificateType, error) {
	vals := map[string]FrontDoorCertificateType{
		"dedicated": FrontDoorCertificateTypeDedicated,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := FrontDoorCertificateType(input)
	return &out, nil
}

type FrontDoorEnabledState string

const (
	FrontDoorEnabledStateDisabled FrontDoorEnabledState = "Disabled"
	FrontDoorEnabledStateEnabled  FrontDoorEnabledState = "Enabled"
)

func PossibleValuesForFrontDoorEnabledState() []string {
	return []string{
		string(FrontDoorEnabledStateDisabled),
		string(FrontDoorEnabledStateEnabled),
	}
}

func parseFrontDoorEnabledState(input string) (*FrontDoorEnabledState, error) {
	vals := map[string]FrontDoorEnabledState{
		"disabled": FrontDoorEnabledStateDisabled,
		"enabled":  FrontDoorEnabledStateEnabled,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := FrontDoorEnabledState(input)
	return &out, nil
}

type FrontDoorForwardingProtocol string

const (
	FrontDoorForwardingProtocolHTTPOnly     FrontDoorForwardingProtocol = "HttpOnly"
	FrontDoorForwardingProtocolHTTPSOnly    FrontDoorForwardingProtocol = "HttpsOnly"
	FrontDoorForwardingProtocolMatchRequest FrontDoorForwardingProtocol = "MatchRequest"
)

func PossibleValuesForFrontDoorForwardingProtocol() []string {
	return []string{
		string(FrontDoorForwardingProtocolHTTPOnly),
		string(FrontDoorForwardingProtocolHTTPSOnly),
		string(FrontDoorForwardingProtocolMatchRequest),
	}
}

func parseFrontDoorForwardingProtocol(input string) (*FrontDoorForwardingProtocol, error) {
	vals := map[string]FrontDoorForwardingProtocol{
		"httponly":     FrontDoorForwardingProtocolHTTPOnly,
		"httpsonly":    FrontDoorForwardingProtocolHTTPSOnly,
		"matchrequest": FrontDoorForwardingProtocolMatchRequest,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := FrontDoorForwardingProtocol(input)
	return &out, nil
}

type FrontDoorHealthProbeMethod string

const (
	FrontDoorHealthProbeMethodGET  FrontDoorHealthProbeMethod = "GET"
	FrontDoorHealthProbeMethodHEAD FrontDoorHealthProbeMethod = "HEAD"
)

func PossibleValuesForFrontDoorHealthProbeMethod() []string {
	return []string{
		string(FrontDoorHealthProbeMethodGET),
		string(FrontDoorHealthProbeMethodHEAD),
	}
}

func parseFrontDoorHealthProbeMethod(input string) (*FrontDoorHealthProbeMethod, error) {
	vals := map[string]FrontDoorHealthProbeMethod{
		"get":  FrontDoorHealthProbeMethodGET,
		"head": FrontDoorHealthProbeMethodHEAD,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := FrontDoorHealthProbeMethod(input)
	return &out, nil
}

type FrontDoorProtocol string

const (
	FrontDoorProtocolHTTP  FrontDoorProtocol = "Http"
	FrontDoorProtocolHTTPS FrontDoorProtocol = "Https"
)

func PossibleValuesForFrontDoorProtocol() []string {
	return []string{
		string(FrontDoorProtocolHTTP),
		string(FrontDoorProtocolHTTPS),
	}
}

func parseFrontDoorProtocol(input string) (*FrontDoorProtocol, error) {
	vals := map[string]FrontDoorProtocol{
		"http":  FrontDoorProtocolHTTP,
		"https": FrontDoorProtocolHTTPS,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := FrontDoorProtocol(input)
	return &out, nil
}

type FrontDoorQuery string

const (
	FrontDoorQueryStripAll       FrontDoorQuery = "StripAll"
	FrontDoorQueryStripAllExcept FrontDoorQuery = "StripAllExcept"
	FrontDoorQueryStripNone      FrontDoorQuery = "StripNone"
	FrontDoorQueryStripOnly      FrontDoorQuery = "StripOnly"
)

func PossibleValuesForFrontDoorQuery() []string {
	return []string{
		string(FrontDoorQueryStripAll),
		string(FrontDoorQueryStripAllExcept),
		string(FrontDoorQueryStripNone),
		string(FrontDoorQueryStripOnly),
	}
}

func parseFrontDoorQuery(input string) (*FrontDoorQuery, error) {
	vals := map[string]FrontDoorQuery{
		"stripall":       FrontDoorQueryStripAll,
		"stripallexcept": FrontDoorQueryStripAllExcept,
		"stripnone":      FrontDoorQueryStripNone,
		"striponly":      FrontDoorQueryStripOnly,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := FrontDoorQuery(input)
	return &out, nil
}

type FrontDoorRedirectProtocol string

const (
	FrontDoorRedirectProtocolHTTPOnly     FrontDoorRedirectProtocol = "HttpOnly"
	FrontDoorRedirectProtocolHTTPSOnly    FrontDoorRedirectProtocol = "HttpsOnly"
	FrontDoorRedirectProtocolMatchRequest FrontDoorRedirectProtocol = "MatchRequest"
)

func PossibleValuesForFrontDoorRedirectProtocol() []string {
	return []string{
		string(FrontDoorRedirectProtocolHTTPOnly),
		string(FrontDoorRedirectProtocolHTTPSOnly),
		string(FrontDoorRedirectProtocolMatchRequest),
	}
}

func parseFrontDoorRedirectProtocol(input string) (*FrontDoorRedirectProtocol, error) {
	vals := map[string]FrontDoorRedirectProtocol{
		"httponly":     FrontDoorRedirectProtocolHTTPOnly,
		"httpsonly":    FrontDoorRedirectProtocolHTTPSOnly,
		"matchrequest": FrontDoorRedirectProtocolMatchRequest,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := FrontDoorRedirectProtocol(input)
	return &out, nil
}

type FrontDoorRedirectType string

const (
	FrontDoorRedirectTypeFound             FrontDoorRedirectType = "Found"
	FrontDoorRedirectTypeMoved             FrontDoorRedirectType = "Moved"
	FrontDoorRedirectTypePermanentRedirect FrontDoorRedirectType = "PermanentRedirect"
	FrontDoorRedirectTypeTemporaryRedirect FrontDoorRedirectType = "TemporaryRedirect"
)

func PossibleValuesForFrontDoorRedirectType() []string {
	return []string{
		string(FrontDoorRedirectTypeFound),
		string(FrontDoorRedirectTypeMoved),
		string(FrontDoorRedirectTypePermanentRedirect),
		string(FrontDoorRedirectTypeTemporaryRedirect),
	}
}

func parseFrontDoorRedirectType(input string) (*FrontDoorRedirectType, error) {
	vals := map[string]FrontDoorRedirectType{
		"found":             FrontDoorRedirectTypeFound,
		"moved":             FrontDoorRedirectTypeMoved,
		"permanentredirect": FrontDoorRedirectTypePermanentRedirect,
		"temporaryredirect": FrontDoorRedirectTypeTemporaryRedirect,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := FrontDoorRedirectType(input)
	return &out, nil
}

type FrontDoorResourceState string

const (
	FrontDoorResourceStateCreating  FrontDoorResourceState = "Creating"
	FrontDoorResourceStateDeleting  FrontDoorResourceState = "Deleting"
	FrontDoorResourceStateDisabled  FrontDoorResourceState = "Disabled"
	FrontDoorResourceStateDisabling FrontDoorResourceState = "Disabling"
	FrontDoorResourceStateEnabled   FrontDoorResourceState = "Enabled"
	FrontDoorResourceStateEnabling  FrontDoorResourceState = "Enabling"
)

func PossibleValuesForFrontDoorResourceState() []string {
	return []string{
		string(FrontDoorResourceStateCreating),
		string(FrontDoorResourceStateDeleting),
		string(FrontDoorResourceStateDisabled),
		string(FrontDoorResourceStateDisabling),
		string(FrontDoorResourceStateEnabled),
		string(FrontDoorResourceStateEnabling),
	}
}

func parseFrontDoorResourceState(input string) (*FrontDoorResourceState, error) {
	vals := map[string]FrontDoorResourceState{
		"creating":  FrontDoorResourceStateCreating,
		"deleting":  FrontDoorResourceStateDeleting,
		"disabled":  FrontDoorResourceStateDisabled,
		"disabling": FrontDoorResourceStateDisabling,
		"enabled":   FrontDoorResourceStateEnabled,
		"enabling":  FrontDoorResourceStateEnabling,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := FrontDoorResourceState(input)
	return &out, nil
}

type FrontDoorTlsProtocolType string

const (
	FrontDoorTlsProtocolTypeServerNameIndication FrontDoorTlsProtocolType = "ServerNameIndication"
)

func PossibleValuesForFrontDoorTlsProtocolType() []string {
	return []string{
		string(FrontDoorTlsProtocolTypeServerNameIndication),
	}
}

func parseFrontDoorTlsProtocolType(input string) (*FrontDoorTlsProtocolType, error) {
	vals := map[string]FrontDoorTlsProtocolType{
		"servernameindication": FrontDoorTlsProtocolTypeServerNameIndication,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := FrontDoorTlsProtocolType(input)
	return &out, nil
}

type HeaderActionType string

const (
	HeaderActionTypeAppend    HeaderActionType = "Append"
	HeaderActionTypeDelete    HeaderActionType = "Delete"
	HeaderActionTypeOverwrite HeaderActionType = "Overwrite"
)

func PossibleValuesForHeaderActionType() []string {
	return []string{
		string(HeaderActionTypeAppend),
		string(HeaderActionTypeDelete),
		string(HeaderActionTypeOverwrite),
	}
}

func parseHeaderActionType(input string) (*HeaderActionType, error) {
	vals := map[string]HeaderActionType{
		"append":    HeaderActionTypeAppend,
		"delete":    HeaderActionTypeDelete,
		"overwrite": HeaderActionTypeOverwrite,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := HeaderActionType(input)
	return &out, nil
}

type HealthProbeEnabled string

const (
	HealthProbeEnabledDisabled HealthProbeEnabled = "Disabled"
	HealthProbeEnabledEnabled  HealthProbeEnabled = "Enabled"
)

func PossibleValuesForHealthProbeEnabled() []string {
	return []string{
		string(HealthProbeEnabledDisabled),
		string(HealthProbeEnabledEnabled),
	}
}

func parseHealthProbeEnabled(input string) (*HealthProbeEnabled, error) {
	vals := map[string]HealthProbeEnabled{
		"disabled": HealthProbeEnabledDisabled,
		"enabled":  HealthProbeEnabledEnabled,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := HealthProbeEnabled(input)
	return &out, nil
}

type MatchProcessingBehavior string

const (
	MatchProcessingBehaviorContinue MatchProcessingBehavior = "Continue"
	MatchProcessingBehaviorStop     MatchProcessingBehavior = "Stop"
)

func PossibleValuesForMatchProcessingBehavior() []string {
	return []string{
		string(MatchProcessingBehaviorContinue),
		string(MatchProcessingBehaviorStop),
	}
}

func parseMatchProcessingBehavior(input string) (*MatchProcessingBehavior, error) {
	vals := map[string]MatchProcessingBehavior{
		"continue": MatchProcessingBehaviorContinue,
		"stop":     MatchProcessingBehaviorStop,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := MatchProcessingBehavior(input)
	return &out, nil
}

type MinimumTLSVersion string

const (
	MinimumTLSVersionOnePointTwo  MinimumTLSVersion = "1.2"
	MinimumTLSVersionOnePointZero MinimumTLSVersion = "1.0"
)

func PossibleValuesForMinimumTLSVersion() []string {
	return []string{
		string(MinimumTLSVersionOnePointTwo),
		string(MinimumTLSVersionOnePointZero),
	}
}

func parseMinimumTLSVersion(input string) (*MinimumTLSVersion, error) {
	vals := map[string]MinimumTLSVersion{
		"1.2": MinimumTLSVersionOnePointTwo,
		"1.0": MinimumTLSVersionOnePointZero,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := MinimumTLSVersion(input)
	return &out, nil
}

type PrivateEndpointStatus string

const (
	PrivateEndpointStatusApproved     PrivateEndpointStatus = "Approved"
	PrivateEndpointStatusDisconnected PrivateEndpointStatus = "Disconnected"
	PrivateEndpointStatusPending      PrivateEndpointStatus = "Pending"
	PrivateEndpointStatusRejected     PrivateEndpointStatus = "Rejected"
	PrivateEndpointStatusTimeout      PrivateEndpointStatus = "Timeout"
)

func PossibleValuesForPrivateEndpointStatus() []string {
	return []string{
		string(PrivateEndpointStatusApproved),
		string(PrivateEndpointStatusDisconnected),
		string(PrivateEndpointStatusPending),
		string(PrivateEndpointStatusRejected),
		string(PrivateEndpointStatusTimeout),
	}
}

func parsePrivateEndpointStatus(input string) (*PrivateEndpointStatus, error) {
	vals := map[string]PrivateEndpointStatus{
		"approved":     PrivateEndpointStatusApproved,
		"disconnected": PrivateEndpointStatusDisconnected,
		"pending":      PrivateEndpointStatusPending,
		"rejected":     PrivateEndpointStatusRejected,
		"timeout":      PrivateEndpointStatusTimeout,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := PrivateEndpointStatus(input)
	return &out, nil
}

type RoutingRuleEnabledState string

const (
	RoutingRuleEnabledStateDisabled RoutingRuleEnabledState = "Disabled"
	RoutingRuleEnabledStateEnabled  RoutingRuleEnabledState = "Enabled"
)

func PossibleValuesForRoutingRuleEnabledState() []string {
	return []string{
		string(RoutingRuleEnabledStateDisabled),
		string(RoutingRuleEnabledStateEnabled),
	}
}

func parseRoutingRuleEnabledState(input string) (*RoutingRuleEnabledState, error) {
	vals := map[string]RoutingRuleEnabledState{
		"disabled": RoutingRuleEnabledStateDisabled,
		"enabled":  RoutingRuleEnabledStateEnabled,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := RoutingRuleEnabledState(input)
	return &out, nil
}

type RulesEngineMatchVariable string

const (
	RulesEngineMatchVariableIsMobile                 RulesEngineMatchVariable = "IsMobile"
	RulesEngineMatchVariablePostArgs                 RulesEngineMatchVariable = "PostArgs"
	RulesEngineMatchVariableQueryString              RulesEngineMatchVariable = "QueryString"
	RulesEngineMatchVariableRemoteAddr               RulesEngineMatchVariable = "RemoteAddr"
	RulesEngineMatchVariableRequestBody              RulesEngineMatchVariable = "RequestBody"
	RulesEngineMatchVariableRequestFilename          RulesEngineMatchVariable = "RequestFilename"
	RulesEngineMatchVariableRequestFilenameExtension RulesEngineMatchVariable = "RequestFilenameExtension"
	RulesEngineMatchVariableRequestHeader            RulesEngineMatchVariable = "RequestHeader"
	RulesEngineMatchVariableRequestMethod            RulesEngineMatchVariable = "RequestMethod"
	RulesEngineMatchVariableRequestPath              RulesEngineMatchVariable = "RequestPath"
	RulesEngineMatchVariableRequestScheme            RulesEngineMatchVariable = "RequestScheme"
	RulesEngineMatchVariableRequestUri               RulesEngineMatchVariable = "RequestUri"
)

func PossibleValuesForRulesEngineMatchVariable() []string {
	return []string{
		string(RulesEngineMatchVariableIsMobile),
		string(RulesEngineMatchVariablePostArgs),
		string(RulesEngineMatchVariableQueryString),
		string(RulesEngineMatchVariableRemoteAddr),
		string(RulesEngineMatchVariableRequestBody),
		string(RulesEngineMatchVariableRequestFilename),
		string(RulesEngineMatchVariableRequestFilenameExtension),
		string(RulesEngineMatchVariableRequestHeader),
		string(RulesEngineMatchVariableRequestMethod),
		string(RulesEngineMatchVariableRequestPath),
		string(RulesEngineMatchVariableRequestScheme),
		string(RulesEngineMatchVariableRequestUri),
	}
}

func parseRulesEngineMatchVariable(input string) (*RulesEngineMatchVariable, error) {
	vals := map[string]RulesEngineMatchVariable{
		"ismobile":                 RulesEngineMatchVariableIsMobile,
		"postargs":                 RulesEngineMatchVariablePostArgs,
		"querystring":              RulesEngineMatchVariableQueryString,
		"remoteaddr":               RulesEngineMatchVariableRemoteAddr,
		"requestbody":              RulesEngineMatchVariableRequestBody,
		"requestfilename":          RulesEngineMatchVariableRequestFilename,
		"requestfilenameextension": RulesEngineMatchVariableRequestFilenameExtension,
		"requestheader":            RulesEngineMatchVariableRequestHeader,
		"requestmethod":            RulesEngineMatchVariableRequestMethod,
		"requestpath":              RulesEngineMatchVariableRequestPath,
		"requestscheme":            RulesEngineMatchVariableRequestScheme,
		"requesturi":               RulesEngineMatchVariableRequestUri,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := RulesEngineMatchVariable(input)
	return &out, nil
}

type RulesEngineOperator string

const (
	RulesEngineOperatorAny                RulesEngineOperator = "Any"
	RulesEngineOperatorBeginsWith         RulesEngineOperator = "BeginsWith"
	RulesEngineOperatorContains           RulesEngineOperator = "Contains"
	RulesEngineOperatorEndsWith           RulesEngineOperator = "EndsWith"
	RulesEngineOperatorEqual              RulesEngineOperator = "Equal"
	RulesEngineOperatorGeoMatch           RulesEngineOperator = "GeoMatch"
	RulesEngineOperatorGreaterThan        RulesEngineOperator = "GreaterThan"
	RulesEngineOperatorGreaterThanOrEqual RulesEngineOperator = "GreaterThanOrEqual"
	RulesEngineOperatorIPMatch            RulesEngineOperator = "IPMatch"
	RulesEngineOperatorLessThan           RulesEngineOperator = "LessThan"
	RulesEngineOperatorLessThanOrEqual    RulesEngineOperator = "LessThanOrEqual"
)

func PossibleValuesForRulesEngineOperator() []string {
	return []string{
		string(RulesEngineOperatorAny),
		string(RulesEngineOperatorBeginsWith),
		string(RulesEngineOperatorContains),
		string(RulesEngineOperatorEndsWith),
		string(RulesEngineOperatorEqual),
		string(RulesEngineOperatorGeoMatch),
		string(RulesEngineOperatorGreaterThan),
		string(RulesEngineOperatorGreaterThanOrEqual),
		string(RulesEngineOperatorIPMatch),
		string(RulesEngineOperatorLessThan),
		string(RulesEngineOperatorLessThanOrEqual),
	}
}

func parseRulesEngineOperator(input string) (*RulesEngineOperator, error) {
	vals := map[string]RulesEngineOperator{
		"any":                RulesEngineOperatorAny,
		"beginswith":         RulesEngineOperatorBeginsWith,
		"contains":           RulesEngineOperatorContains,
		"endswith":           RulesEngineOperatorEndsWith,
		"equal":              RulesEngineOperatorEqual,
		"geomatch":           RulesEngineOperatorGeoMatch,
		"greaterthan":        RulesEngineOperatorGreaterThan,
		"greaterthanorequal": RulesEngineOperatorGreaterThanOrEqual,
		"ipmatch":            RulesEngineOperatorIPMatch,
		"lessthan":           RulesEngineOperatorLessThan,
		"lessthanorequal":    RulesEngineOperatorLessThanOrEqual,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := RulesEngineOperator(input)
	return &out, nil
}

type SessionAffinityEnabledState string

const (
	SessionAffinityEnabledStateDisabled SessionAffinityEnabledState = "Disabled"
	SessionAffinityEnabledStateEnabled  SessionAffinityEnabledState = "Enabled"
)

func PossibleValuesForSessionAffinityEnabledState() []string {
	return []string{
		string(SessionAffinityEnabledStateDisabled),
		string(SessionAffinityEnabledStateEnabled),
	}
}

func parseSessionAffinityEnabledState(input string) (*SessionAffinityEnabledState, error) {
	vals := map[string]SessionAffinityEnabledState{
		"disabled": SessionAffinityEnabledStateDisabled,
		"enabled":  SessionAffinityEnabledStateEnabled,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := SessionAffinityEnabledState(input)
	return &out, nil
}

type Transform string

const (
	TransformLowercase   Transform = "Lowercase"
	TransformRemoveNulls Transform = "RemoveNulls"
	TransformTrim        Transform = "Trim"
	TransformUppercase   Transform = "Uppercase"
	TransformUrlDecode   Transform = "UrlDecode"
	TransformUrlEncode   Transform = "UrlEncode"
)

func PossibleValuesForTransform() []string {
	return []string{
		string(TransformLowercase),
		string(TransformRemoveNulls),
		string(TransformTrim),
		string(TransformUppercase),
		string(TransformUrlDecode),
		string(TransformUrlEncode),
	}
}

func parseTransform(input string) (*Transform, error) {
	vals := map[string]Transform{
		"lowercase":   TransformLowercase,
		"removenulls": TransformRemoveNulls,
		"trim":        TransformTrim,
		"uppercase":   TransformUppercase,
		"urldecode":   TransformUrlDecode,
		"urlencode":   TransformUrlEncode,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := Transform(input)
	return &out, nil
}
