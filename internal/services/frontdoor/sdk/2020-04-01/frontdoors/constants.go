package frontdoors

type BackendEnabledState string

const (
	BackendEnabledStateDisabled BackendEnabledState = "Disabled"
	BackendEnabledStateEnabled  BackendEnabledState = "Enabled"
)

type CustomHttpsProvisioningState string

const (
	CustomHttpsProvisioningStateDisabled  CustomHttpsProvisioningState = "Disabled"
	CustomHttpsProvisioningStateDisabling CustomHttpsProvisioningState = "Disabling"
	CustomHttpsProvisioningStateEnabled   CustomHttpsProvisioningState = "Enabled"
	CustomHttpsProvisioningStateEnabling  CustomHttpsProvisioningState = "Enabling"
	CustomHttpsProvisioningStateFailed    CustomHttpsProvisioningState = "Failed"
)

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

type DynamicCompressionEnabled string

const (
	DynamicCompressionEnabledDisabled DynamicCompressionEnabled = "Disabled"
	DynamicCompressionEnabledEnabled  DynamicCompressionEnabled = "Enabled"
)

type EnforceCertificateNameCheckEnabledState string

const (
	EnforceCertificateNameCheckEnabledStateDisabled EnforceCertificateNameCheckEnabledState = "Disabled"
	EnforceCertificateNameCheckEnabledStateEnabled  EnforceCertificateNameCheckEnabledState = "Enabled"
)

type FrontDoorCertificateSource string

const (
	FrontDoorCertificateSourceAzureKeyVault FrontDoorCertificateSource = "AzureKeyVault"
	FrontDoorCertificateSourceFrontDoor     FrontDoorCertificateSource = "FrontDoor"
)

type FrontDoorCertificateType string

const (
	FrontDoorCertificateTypeDedicated FrontDoorCertificateType = "Dedicated"
)

type FrontDoorEnabledState string

const (
	FrontDoorEnabledStateDisabled FrontDoorEnabledState = "Disabled"
	FrontDoorEnabledStateEnabled  FrontDoorEnabledState = "Enabled"
)

type FrontDoorForwardingProtocol string

const (
	FrontDoorForwardingProtocolHttpOnly     FrontDoorForwardingProtocol = "HttpOnly"
	FrontDoorForwardingProtocolHttpsOnly    FrontDoorForwardingProtocol = "HttpsOnly"
	FrontDoorForwardingProtocolMatchRequest FrontDoorForwardingProtocol = "MatchRequest"
)

type FrontDoorHealthProbeMethod string

const (
	FrontDoorHealthProbeMethodGET  FrontDoorHealthProbeMethod = "GET"
	FrontDoorHealthProbeMethodHEAD FrontDoorHealthProbeMethod = "HEAD"
)

type FrontDoorProtocol string

const (
	FrontDoorProtocolHttp  FrontDoorProtocol = "Http"
	FrontDoorProtocolHttps FrontDoorProtocol = "Https"
)

type FrontDoorQuery string

const (
	FrontDoorQueryStripAll       FrontDoorQuery = "StripAll"
	FrontDoorQueryStripAllExcept FrontDoorQuery = "StripAllExcept"
	FrontDoorQueryStripNone      FrontDoorQuery = "StripNone"
	FrontDoorQueryStripOnly      FrontDoorQuery = "StripOnly"
)

type FrontDoorRedirectProtocol string

const (
	FrontDoorRedirectProtocolHttpOnly     FrontDoorRedirectProtocol = "HttpOnly"
	FrontDoorRedirectProtocolHttpsOnly    FrontDoorRedirectProtocol = "HttpsOnly"
	FrontDoorRedirectProtocolMatchRequest FrontDoorRedirectProtocol = "MatchRequest"
)

type FrontDoorRedirectType string

const (
	FrontDoorRedirectTypeFound             FrontDoorRedirectType = "Found"
	FrontDoorRedirectTypeMoved             FrontDoorRedirectType = "Moved"
	FrontDoorRedirectTypePermanentRedirect FrontDoorRedirectType = "PermanentRedirect"
	FrontDoorRedirectTypeTemporaryRedirect FrontDoorRedirectType = "TemporaryRedirect"
)

type FrontDoorResourceState string

const (
	FrontDoorResourceStateCreating  FrontDoorResourceState = "Creating"
	FrontDoorResourceStateDeleting  FrontDoorResourceState = "Deleting"
	FrontDoorResourceStateDisabled  FrontDoorResourceState = "Disabled"
	FrontDoorResourceStateDisabling FrontDoorResourceState = "Disabling"
	FrontDoorResourceStateEnabled   FrontDoorResourceState = "Enabled"
	FrontDoorResourceStateEnabling  FrontDoorResourceState = "Enabling"
)

type FrontDoorTlsProtocolType string

const (
	FrontDoorTlsProtocolTypeServerNameIndication FrontDoorTlsProtocolType = "ServerNameIndication"
)

type HeaderActionType string

const (
	HeaderActionTypeAppend    HeaderActionType = "Append"
	HeaderActionTypeDelete    HeaderActionType = "Delete"
	HeaderActionTypeOverwrite HeaderActionType = "Overwrite"
)

type HealthProbeEnabled string

const (
	HealthProbeEnabledDisabled HealthProbeEnabled = "Disabled"
	HealthProbeEnabledEnabled  HealthProbeEnabled = "Enabled"
)

type MatchProcessingBehavior string

const (
	MatchProcessingBehaviorContinue MatchProcessingBehavior = "Continue"
	MatchProcessingBehaviorStop     MatchProcessingBehavior = "Stop"
)

type MinimumTLSVersion string

const (
	MinimumTLSVersionOnePointTwo  MinimumTLSVersion = "1.2"
	MinimumTLSVersionOnePointZero MinimumTLSVersion = "1.0"
)

type PrivateEndpointStatus string

const (
	PrivateEndpointStatusApproved     PrivateEndpointStatus = "Approved"
	PrivateEndpointStatusDisconnected PrivateEndpointStatus = "Disconnected"
	PrivateEndpointStatusPending      PrivateEndpointStatus = "Pending"
	PrivateEndpointStatusRejected     PrivateEndpointStatus = "Rejected"
	PrivateEndpointStatusTimeout      PrivateEndpointStatus = "Timeout"
)

type RoutingRuleEnabledState string

const (
	RoutingRuleEnabledStateDisabled RoutingRuleEnabledState = "Disabled"
	RoutingRuleEnabledStateEnabled  RoutingRuleEnabledState = "Enabled"
)

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

type SessionAffinityEnabledState string

const (
	SessionAffinityEnabledStateDisabled SessionAffinityEnabledState = "Disabled"
	SessionAffinityEnabledStateEnabled  SessionAffinityEnabledState = "Enabled"
)

type Transform string

const (
	TransformLowercase   Transform = "Lowercase"
	TransformRemoveNulls Transform = "RemoveNulls"
	TransformTrim        Transform = "Trim"
	TransformUppercase   Transform = "Uppercase"
	TransformUrlDecode   Transform = "UrlDecode"
	TransformUrlEncode   Transform = "UrlEncode"
)
