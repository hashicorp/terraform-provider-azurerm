package endpoints

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

type ClientPortOperator string

const (
	ClientPortOperatorAny                ClientPortOperator = "Any"
	ClientPortOperatorBeginsWith         ClientPortOperator = "BeginsWith"
	ClientPortOperatorContains           ClientPortOperator = "Contains"
	ClientPortOperatorEndsWith           ClientPortOperator = "EndsWith"
	ClientPortOperatorEqual              ClientPortOperator = "Equal"
	ClientPortOperatorGreaterThan        ClientPortOperator = "GreaterThan"
	ClientPortOperatorGreaterThanOrEqual ClientPortOperator = "GreaterThanOrEqual"
	ClientPortOperatorLessThan           ClientPortOperator = "LessThan"
	ClientPortOperatorLessThanOrEqual    ClientPortOperator = "LessThanOrEqual"
	ClientPortOperatorRegEx              ClientPortOperator = "RegEx"
)

func PossibleValuesForClientPortOperator() []string {
	return []string{
		string(ClientPortOperatorAny),
		string(ClientPortOperatorBeginsWith),
		string(ClientPortOperatorContains),
		string(ClientPortOperatorEndsWith),
		string(ClientPortOperatorEqual),
		string(ClientPortOperatorGreaterThan),
		string(ClientPortOperatorGreaterThanOrEqual),
		string(ClientPortOperatorLessThan),
		string(ClientPortOperatorLessThanOrEqual),
		string(ClientPortOperatorRegEx),
	}
}

func parseClientPortOperator(input string) (*ClientPortOperator, error) {
	vals := map[string]ClientPortOperator{
		"any":                ClientPortOperatorAny,
		"beginswith":         ClientPortOperatorBeginsWith,
		"contains":           ClientPortOperatorContains,
		"endswith":           ClientPortOperatorEndsWith,
		"equal":              ClientPortOperatorEqual,
		"greaterthan":        ClientPortOperatorGreaterThan,
		"greaterthanorequal": ClientPortOperatorGreaterThanOrEqual,
		"lessthan":           ClientPortOperatorLessThan,
		"lessthanorequal":    ClientPortOperatorLessThanOrEqual,
		"regex":              ClientPortOperatorRegEx,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := ClientPortOperator(input)
	return &out, nil
}

type CookiesOperator string

const (
	CookiesOperatorAny                CookiesOperator = "Any"
	CookiesOperatorBeginsWith         CookiesOperator = "BeginsWith"
	CookiesOperatorContains           CookiesOperator = "Contains"
	CookiesOperatorEndsWith           CookiesOperator = "EndsWith"
	CookiesOperatorEqual              CookiesOperator = "Equal"
	CookiesOperatorGreaterThan        CookiesOperator = "GreaterThan"
	CookiesOperatorGreaterThanOrEqual CookiesOperator = "GreaterThanOrEqual"
	CookiesOperatorLessThan           CookiesOperator = "LessThan"
	CookiesOperatorLessThanOrEqual    CookiesOperator = "LessThanOrEqual"
	CookiesOperatorRegEx              CookiesOperator = "RegEx"
)

func PossibleValuesForCookiesOperator() []string {
	return []string{
		string(CookiesOperatorAny),
		string(CookiesOperatorBeginsWith),
		string(CookiesOperatorContains),
		string(CookiesOperatorEndsWith),
		string(CookiesOperatorEqual),
		string(CookiesOperatorGreaterThan),
		string(CookiesOperatorGreaterThanOrEqual),
		string(CookiesOperatorLessThan),
		string(CookiesOperatorLessThanOrEqual),
		string(CookiesOperatorRegEx),
	}
}

func parseCookiesOperator(input string) (*CookiesOperator, error) {
	vals := map[string]CookiesOperator{
		"any":                CookiesOperatorAny,
		"beginswith":         CookiesOperatorBeginsWith,
		"contains":           CookiesOperatorContains,
		"endswith":           CookiesOperatorEndsWith,
		"equal":              CookiesOperatorEqual,
		"greaterthan":        CookiesOperatorGreaterThan,
		"greaterthanorequal": CookiesOperatorGreaterThanOrEqual,
		"lessthan":           CookiesOperatorLessThan,
		"lessthanorequal":    CookiesOperatorLessThanOrEqual,
		"regex":              CookiesOperatorRegEx,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := CookiesOperator(input)
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

type DeliveryRuleAction string

const (
	DeliveryRuleActionCacheExpiration            DeliveryRuleAction = "CacheExpiration"
	DeliveryRuleActionCacheKeyQueryString        DeliveryRuleAction = "CacheKeyQueryString"
	DeliveryRuleActionModifyRequestHeader        DeliveryRuleAction = "ModifyRequestHeader"
	DeliveryRuleActionModifyResponseHeader       DeliveryRuleAction = "ModifyResponseHeader"
	DeliveryRuleActionOriginGroupOverride        DeliveryRuleAction = "OriginGroupOverride"
	DeliveryRuleActionRouteConfigurationOverride DeliveryRuleAction = "RouteConfigurationOverride"
	DeliveryRuleActionUrlRedirect                DeliveryRuleAction = "UrlRedirect"
	DeliveryRuleActionUrlRewrite                 DeliveryRuleAction = "UrlRewrite"
	DeliveryRuleActionUrlSigning                 DeliveryRuleAction = "UrlSigning"
)

func PossibleValuesForDeliveryRuleAction() []string {
	return []string{
		string(DeliveryRuleActionCacheExpiration),
		string(DeliveryRuleActionCacheKeyQueryString),
		string(DeliveryRuleActionModifyRequestHeader),
		string(DeliveryRuleActionModifyResponseHeader),
		string(DeliveryRuleActionOriginGroupOverride),
		string(DeliveryRuleActionRouteConfigurationOverride),
		string(DeliveryRuleActionUrlRedirect),
		string(DeliveryRuleActionUrlRewrite),
		string(DeliveryRuleActionUrlSigning),
	}
}

func parseDeliveryRuleAction(input string) (*DeliveryRuleAction, error) {
	vals := map[string]DeliveryRuleAction{
		"cacheexpiration":            DeliveryRuleActionCacheExpiration,
		"cachekeyquerystring":        DeliveryRuleActionCacheKeyQueryString,
		"modifyrequestheader":        DeliveryRuleActionModifyRequestHeader,
		"modifyresponseheader":       DeliveryRuleActionModifyResponseHeader,
		"origingroupoverride":        DeliveryRuleActionOriginGroupOverride,
		"routeconfigurationoverride": DeliveryRuleActionRouteConfigurationOverride,
		"urlredirect":                DeliveryRuleActionUrlRedirect,
		"urlrewrite":                 DeliveryRuleActionUrlRewrite,
		"urlsigning":                 DeliveryRuleActionUrlSigning,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := DeliveryRuleAction(input)
	return &out, nil
}

type EndpointResourceState string

const (
	EndpointResourceStateCreating EndpointResourceState = "Creating"
	EndpointResourceStateDeleting EndpointResourceState = "Deleting"
	EndpointResourceStateRunning  EndpointResourceState = "Running"
	EndpointResourceStateStarting EndpointResourceState = "Starting"
	EndpointResourceStateStopped  EndpointResourceState = "Stopped"
	EndpointResourceStateStopping EndpointResourceState = "Stopping"
)

func PossibleValuesForEndpointResourceState() []string {
	return []string{
		string(EndpointResourceStateCreating),
		string(EndpointResourceStateDeleting),
		string(EndpointResourceStateRunning),
		string(EndpointResourceStateStarting),
		string(EndpointResourceStateStopped),
		string(EndpointResourceStateStopping),
	}
}

func parseEndpointResourceState(input string) (*EndpointResourceState, error) {
	vals := map[string]EndpointResourceState{
		"creating": EndpointResourceStateCreating,
		"deleting": EndpointResourceStateDeleting,
		"running":  EndpointResourceStateRunning,
		"starting": EndpointResourceStateStarting,
		"stopped":  EndpointResourceStateStopped,
		"stopping": EndpointResourceStateStopping,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := EndpointResourceState(input)
	return &out, nil
}

type GeoFilterActions string

const (
	GeoFilterActionsAllow GeoFilterActions = "Allow"
	GeoFilterActionsBlock GeoFilterActions = "Block"
)

func PossibleValuesForGeoFilterActions() []string {
	return []string{
		string(GeoFilterActionsAllow),
		string(GeoFilterActionsBlock),
	}
}

func parseGeoFilterActions(input string) (*GeoFilterActions, error) {
	vals := map[string]GeoFilterActions{
		"allow": GeoFilterActionsAllow,
		"block": GeoFilterActionsBlock,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := GeoFilterActions(input)
	return &out, nil
}

type HealthProbeRequestType string

const (
	HealthProbeRequestTypeGET    HealthProbeRequestType = "GET"
	HealthProbeRequestTypeHEAD   HealthProbeRequestType = "HEAD"
	HealthProbeRequestTypeNotSet HealthProbeRequestType = "NotSet"
)

func PossibleValuesForHealthProbeRequestType() []string {
	return []string{
		string(HealthProbeRequestTypeGET),
		string(HealthProbeRequestTypeHEAD),
		string(HealthProbeRequestTypeNotSet),
	}
}

func parseHealthProbeRequestType(input string) (*HealthProbeRequestType, error) {
	vals := map[string]HealthProbeRequestType{
		"get":    HealthProbeRequestTypeGET,
		"head":   HealthProbeRequestTypeHEAD,
		"notset": HealthProbeRequestTypeNotSet,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := HealthProbeRequestType(input)
	return &out, nil
}

type HostNameOperator string

const (
	HostNameOperatorAny                HostNameOperator = "Any"
	HostNameOperatorBeginsWith         HostNameOperator = "BeginsWith"
	HostNameOperatorContains           HostNameOperator = "Contains"
	HostNameOperatorEndsWith           HostNameOperator = "EndsWith"
	HostNameOperatorEqual              HostNameOperator = "Equal"
	HostNameOperatorGreaterThan        HostNameOperator = "GreaterThan"
	HostNameOperatorGreaterThanOrEqual HostNameOperator = "GreaterThanOrEqual"
	HostNameOperatorLessThan           HostNameOperator = "LessThan"
	HostNameOperatorLessThanOrEqual    HostNameOperator = "LessThanOrEqual"
	HostNameOperatorRegEx              HostNameOperator = "RegEx"
)

func PossibleValuesForHostNameOperator() []string {
	return []string{
		string(HostNameOperatorAny),
		string(HostNameOperatorBeginsWith),
		string(HostNameOperatorContains),
		string(HostNameOperatorEndsWith),
		string(HostNameOperatorEqual),
		string(HostNameOperatorGreaterThan),
		string(HostNameOperatorGreaterThanOrEqual),
		string(HostNameOperatorLessThan),
		string(HostNameOperatorLessThanOrEqual),
		string(HostNameOperatorRegEx),
	}
}

func parseHostNameOperator(input string) (*HostNameOperator, error) {
	vals := map[string]HostNameOperator{
		"any":                HostNameOperatorAny,
		"beginswith":         HostNameOperatorBeginsWith,
		"contains":           HostNameOperatorContains,
		"endswith":           HostNameOperatorEndsWith,
		"equal":              HostNameOperatorEqual,
		"greaterthan":        HostNameOperatorGreaterThan,
		"greaterthanorequal": HostNameOperatorGreaterThanOrEqual,
		"lessthan":           HostNameOperatorLessThan,
		"lessthanorequal":    HostNameOperatorLessThanOrEqual,
		"regex":              HostNameOperatorRegEx,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := HostNameOperator(input)
	return &out, nil
}

type HttpVersionOperator string

const (
	HttpVersionOperatorEqual HttpVersionOperator = "Equal"
)

func PossibleValuesForHttpVersionOperator() []string {
	return []string{
		string(HttpVersionOperatorEqual),
	}
}

func parseHttpVersionOperator(input string) (*HttpVersionOperator, error) {
	vals := map[string]HttpVersionOperator{
		"equal": HttpVersionOperatorEqual,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := HttpVersionOperator(input)
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

type IsDeviceOperator string

const (
	IsDeviceOperatorEqual IsDeviceOperator = "Equal"
)

func PossibleValuesForIsDeviceOperator() []string {
	return []string{
		string(IsDeviceOperatorEqual),
	}
}

func parseIsDeviceOperator(input string) (*IsDeviceOperator, error) {
	vals := map[string]IsDeviceOperator{
		"equal": IsDeviceOperatorEqual,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := IsDeviceOperator(input)
	return &out, nil
}

type MatchValues string

const (
	MatchValuesDELETE  MatchValues = "DELETE"
	MatchValuesGET     MatchValues = "GET"
	MatchValuesHEAD    MatchValues = "HEAD"
	MatchValuesOPTIONS MatchValues = "OPTIONS"
	MatchValuesPOST    MatchValues = "POST"
	MatchValuesPUT     MatchValues = "PUT"
	MatchValuesTRACE   MatchValues = "TRACE"
)

func PossibleValuesForMatchValues() []string {
	return []string{
		string(MatchValuesDELETE),
		string(MatchValuesGET),
		string(MatchValuesHEAD),
		string(MatchValuesOPTIONS),
		string(MatchValuesPOST),
		string(MatchValuesPUT),
		string(MatchValuesTRACE),
	}
}

func parseMatchValues(input string) (*MatchValues, error) {
	vals := map[string]MatchValues{
		"delete":  MatchValuesDELETE,
		"get":     MatchValuesGET,
		"head":    MatchValuesHEAD,
		"options": MatchValuesOPTIONS,
		"post":    MatchValuesPOST,
		"put":     MatchValuesPUT,
		"trace":   MatchValuesTRACE,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := MatchValues(input)
	return &out, nil
}

type MatchVariable string

const (
	MatchVariableClientPort       MatchVariable = "ClientPort"
	MatchVariableCookies          MatchVariable = "Cookies"
	MatchVariableHostName         MatchVariable = "HostName"
	MatchVariableHttpVersion      MatchVariable = "HttpVersion"
	MatchVariableIsDevice         MatchVariable = "IsDevice"
	MatchVariablePostArgs         MatchVariable = "PostArgs"
	MatchVariableQueryString      MatchVariable = "QueryString"
	MatchVariableRemoteAddress    MatchVariable = "RemoteAddress"
	MatchVariableRequestBody      MatchVariable = "RequestBody"
	MatchVariableRequestHeader    MatchVariable = "RequestHeader"
	MatchVariableRequestMethod    MatchVariable = "RequestMethod"
	MatchVariableRequestScheme    MatchVariable = "RequestScheme"
	MatchVariableRequestUri       MatchVariable = "RequestUri"
	MatchVariableServerPort       MatchVariable = "ServerPort"
	MatchVariableSocketAddr       MatchVariable = "SocketAddr"
	MatchVariableSslProtocol      MatchVariable = "SslProtocol"
	MatchVariableUrlFileExtension MatchVariable = "UrlFileExtension"
	MatchVariableUrlFileName      MatchVariable = "UrlFileName"
	MatchVariableUrlPath          MatchVariable = "UrlPath"
)

func PossibleValuesForMatchVariable() []string {
	return []string{
		string(MatchVariableClientPort),
		string(MatchVariableCookies),
		string(MatchVariableHostName),
		string(MatchVariableHttpVersion),
		string(MatchVariableIsDevice),
		string(MatchVariablePostArgs),
		string(MatchVariableQueryString),
		string(MatchVariableRemoteAddress),
		string(MatchVariableRequestBody),
		string(MatchVariableRequestHeader),
		string(MatchVariableRequestMethod),
		string(MatchVariableRequestScheme),
		string(MatchVariableRequestUri),
		string(MatchVariableServerPort),
		string(MatchVariableSocketAddr),
		string(MatchVariableSslProtocol),
		string(MatchVariableUrlFileExtension),
		string(MatchVariableUrlFileName),
		string(MatchVariableUrlPath),
	}
}

func parseMatchVariable(input string) (*MatchVariable, error) {
	vals := map[string]MatchVariable{
		"clientport":       MatchVariableClientPort,
		"cookies":          MatchVariableCookies,
		"hostname":         MatchVariableHostName,
		"httpversion":      MatchVariableHttpVersion,
		"isdevice":         MatchVariableIsDevice,
		"postargs":         MatchVariablePostArgs,
		"querystring":      MatchVariableQueryString,
		"remoteaddress":    MatchVariableRemoteAddress,
		"requestbody":      MatchVariableRequestBody,
		"requestheader":    MatchVariableRequestHeader,
		"requestmethod":    MatchVariableRequestMethod,
		"requestscheme":    MatchVariableRequestScheme,
		"requesturi":       MatchVariableRequestUri,
		"serverport":       MatchVariableServerPort,
		"socketaddr":       MatchVariableSocketAddr,
		"sslprotocol":      MatchVariableSslProtocol,
		"urlfileextension": MatchVariableUrlFileExtension,
		"urlfilename":      MatchVariableUrlFileName,
		"urlpath":          MatchVariableUrlPath,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := MatchVariable(input)
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

type Operator string

const (
	OperatorEqual Operator = "Equal"
)

func PossibleValuesForOperator() []string {
	return []string{
		string(OperatorEqual),
	}
}

func parseOperator(input string) (*Operator, error) {
	vals := map[string]Operator{
		"equal": OperatorEqual,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := Operator(input)
	return &out, nil
}

type OptimizationType string

const (
	OptimizationTypeDynamicSiteAcceleration     OptimizationType = "DynamicSiteAcceleration"
	OptimizationTypeGeneralMediaStreaming       OptimizationType = "GeneralMediaStreaming"
	OptimizationTypeGeneralWebDelivery          OptimizationType = "GeneralWebDelivery"
	OptimizationTypeLargeFileDownload           OptimizationType = "LargeFileDownload"
	OptimizationTypeVideoOnDemandMediaStreaming OptimizationType = "VideoOnDemandMediaStreaming"
)

func PossibleValuesForOptimizationType() []string {
	return []string{
		string(OptimizationTypeDynamicSiteAcceleration),
		string(OptimizationTypeGeneralMediaStreaming),
		string(OptimizationTypeGeneralWebDelivery),
		string(OptimizationTypeLargeFileDownload),
		string(OptimizationTypeVideoOnDemandMediaStreaming),
	}
}

func parseOptimizationType(input string) (*OptimizationType, error) {
	vals := map[string]OptimizationType{
		"dynamicsiteacceleration":     OptimizationTypeDynamicSiteAcceleration,
		"generalmediastreaming":       OptimizationTypeGeneralMediaStreaming,
		"generalwebdelivery":          OptimizationTypeGeneralWebDelivery,
		"largefiledownload":           OptimizationTypeLargeFileDownload,
		"videoondemandmediastreaming": OptimizationTypeVideoOnDemandMediaStreaming,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := OptimizationType(input)
	return &out, nil
}

type PostArgsOperator string

const (
	PostArgsOperatorAny                PostArgsOperator = "Any"
	PostArgsOperatorBeginsWith         PostArgsOperator = "BeginsWith"
	PostArgsOperatorContains           PostArgsOperator = "Contains"
	PostArgsOperatorEndsWith           PostArgsOperator = "EndsWith"
	PostArgsOperatorEqual              PostArgsOperator = "Equal"
	PostArgsOperatorGreaterThan        PostArgsOperator = "GreaterThan"
	PostArgsOperatorGreaterThanOrEqual PostArgsOperator = "GreaterThanOrEqual"
	PostArgsOperatorLessThan           PostArgsOperator = "LessThan"
	PostArgsOperatorLessThanOrEqual    PostArgsOperator = "LessThanOrEqual"
	PostArgsOperatorRegEx              PostArgsOperator = "RegEx"
)

func PossibleValuesForPostArgsOperator() []string {
	return []string{
		string(PostArgsOperatorAny),
		string(PostArgsOperatorBeginsWith),
		string(PostArgsOperatorContains),
		string(PostArgsOperatorEndsWith),
		string(PostArgsOperatorEqual),
		string(PostArgsOperatorGreaterThan),
		string(PostArgsOperatorGreaterThanOrEqual),
		string(PostArgsOperatorLessThan),
		string(PostArgsOperatorLessThanOrEqual),
		string(PostArgsOperatorRegEx),
	}
}

func parsePostArgsOperator(input string) (*PostArgsOperator, error) {
	vals := map[string]PostArgsOperator{
		"any":                PostArgsOperatorAny,
		"beginswith":         PostArgsOperatorBeginsWith,
		"contains":           PostArgsOperatorContains,
		"endswith":           PostArgsOperatorEndsWith,
		"equal":              PostArgsOperatorEqual,
		"greaterthan":        PostArgsOperatorGreaterThan,
		"greaterthanorequal": PostArgsOperatorGreaterThanOrEqual,
		"lessthan":           PostArgsOperatorLessThan,
		"lessthanorequal":    PostArgsOperatorLessThanOrEqual,
		"regex":              PostArgsOperatorRegEx,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := PostArgsOperator(input)
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

type ProbeProtocol string

const (
	ProbeProtocolHttp   ProbeProtocol = "Http"
	ProbeProtocolHttps  ProbeProtocol = "Https"
	ProbeProtocolNotSet ProbeProtocol = "NotSet"
)

func PossibleValuesForProbeProtocol() []string {
	return []string{
		string(ProbeProtocolHttp),
		string(ProbeProtocolHttps),
		string(ProbeProtocolNotSet),
	}
}

func parseProbeProtocol(input string) (*ProbeProtocol, error) {
	vals := map[string]ProbeProtocol{
		"http":   ProbeProtocolHttp,
		"https":  ProbeProtocolHttps,
		"notset": ProbeProtocolNotSet,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := ProbeProtocol(input)
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

type QueryStringCachingBehavior string

const (
	QueryStringCachingBehaviorBypassCaching     QueryStringCachingBehavior = "BypassCaching"
	QueryStringCachingBehaviorIgnoreQueryString QueryStringCachingBehavior = "IgnoreQueryString"
	QueryStringCachingBehaviorNotSet            QueryStringCachingBehavior = "NotSet"
	QueryStringCachingBehaviorUseQueryString    QueryStringCachingBehavior = "UseQueryString"
)

func PossibleValuesForQueryStringCachingBehavior() []string {
	return []string{
		string(QueryStringCachingBehaviorBypassCaching),
		string(QueryStringCachingBehaviorIgnoreQueryString),
		string(QueryStringCachingBehaviorNotSet),
		string(QueryStringCachingBehaviorUseQueryString),
	}
}

func parseQueryStringCachingBehavior(input string) (*QueryStringCachingBehavior, error) {
	vals := map[string]QueryStringCachingBehavior{
		"bypasscaching":     QueryStringCachingBehaviorBypassCaching,
		"ignorequerystring": QueryStringCachingBehaviorIgnoreQueryString,
		"notset":            QueryStringCachingBehaviorNotSet,
		"usequerystring":    QueryStringCachingBehaviorUseQueryString,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := QueryStringCachingBehavior(input)
	return &out, nil
}

type QueryStringOperator string

const (
	QueryStringOperatorAny                QueryStringOperator = "Any"
	QueryStringOperatorBeginsWith         QueryStringOperator = "BeginsWith"
	QueryStringOperatorContains           QueryStringOperator = "Contains"
	QueryStringOperatorEndsWith           QueryStringOperator = "EndsWith"
	QueryStringOperatorEqual              QueryStringOperator = "Equal"
	QueryStringOperatorGreaterThan        QueryStringOperator = "GreaterThan"
	QueryStringOperatorGreaterThanOrEqual QueryStringOperator = "GreaterThanOrEqual"
	QueryStringOperatorLessThan           QueryStringOperator = "LessThan"
	QueryStringOperatorLessThanOrEqual    QueryStringOperator = "LessThanOrEqual"
	QueryStringOperatorRegEx              QueryStringOperator = "RegEx"
)

func PossibleValuesForQueryStringOperator() []string {
	return []string{
		string(QueryStringOperatorAny),
		string(QueryStringOperatorBeginsWith),
		string(QueryStringOperatorContains),
		string(QueryStringOperatorEndsWith),
		string(QueryStringOperatorEqual),
		string(QueryStringOperatorGreaterThan),
		string(QueryStringOperatorGreaterThanOrEqual),
		string(QueryStringOperatorLessThan),
		string(QueryStringOperatorLessThanOrEqual),
		string(QueryStringOperatorRegEx),
	}
}

func parseQueryStringOperator(input string) (*QueryStringOperator, error) {
	vals := map[string]QueryStringOperator{
		"any":                QueryStringOperatorAny,
		"beginswith":         QueryStringOperatorBeginsWith,
		"contains":           QueryStringOperatorContains,
		"endswith":           QueryStringOperatorEndsWith,
		"equal":              QueryStringOperatorEqual,
		"greaterthan":        QueryStringOperatorGreaterThan,
		"greaterthanorequal": QueryStringOperatorGreaterThanOrEqual,
		"lessthan":           QueryStringOperatorLessThan,
		"lessthanorequal":    QueryStringOperatorLessThanOrEqual,
		"regex":              QueryStringOperatorRegEx,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := QueryStringOperator(input)
	return &out, nil
}

type RemoteAddressOperator string

const (
	RemoteAddressOperatorAny      RemoteAddressOperator = "Any"
	RemoteAddressOperatorGeoMatch RemoteAddressOperator = "GeoMatch"
	RemoteAddressOperatorIPMatch  RemoteAddressOperator = "IPMatch"
)

func PossibleValuesForRemoteAddressOperator() []string {
	return []string{
		string(RemoteAddressOperatorAny),
		string(RemoteAddressOperatorGeoMatch),
		string(RemoteAddressOperatorIPMatch),
	}
}

func parseRemoteAddressOperator(input string) (*RemoteAddressOperator, error) {
	vals := map[string]RemoteAddressOperator{
		"any":      RemoteAddressOperatorAny,
		"geomatch": RemoteAddressOperatorGeoMatch,
		"ipmatch":  RemoteAddressOperatorIPMatch,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := RemoteAddressOperator(input)
	return &out, nil
}

type RequestBodyOperator string

const (
	RequestBodyOperatorAny                RequestBodyOperator = "Any"
	RequestBodyOperatorBeginsWith         RequestBodyOperator = "BeginsWith"
	RequestBodyOperatorContains           RequestBodyOperator = "Contains"
	RequestBodyOperatorEndsWith           RequestBodyOperator = "EndsWith"
	RequestBodyOperatorEqual              RequestBodyOperator = "Equal"
	RequestBodyOperatorGreaterThan        RequestBodyOperator = "GreaterThan"
	RequestBodyOperatorGreaterThanOrEqual RequestBodyOperator = "GreaterThanOrEqual"
	RequestBodyOperatorLessThan           RequestBodyOperator = "LessThan"
	RequestBodyOperatorLessThanOrEqual    RequestBodyOperator = "LessThanOrEqual"
	RequestBodyOperatorRegEx              RequestBodyOperator = "RegEx"
)

func PossibleValuesForRequestBodyOperator() []string {
	return []string{
		string(RequestBodyOperatorAny),
		string(RequestBodyOperatorBeginsWith),
		string(RequestBodyOperatorContains),
		string(RequestBodyOperatorEndsWith),
		string(RequestBodyOperatorEqual),
		string(RequestBodyOperatorGreaterThan),
		string(RequestBodyOperatorGreaterThanOrEqual),
		string(RequestBodyOperatorLessThan),
		string(RequestBodyOperatorLessThanOrEqual),
		string(RequestBodyOperatorRegEx),
	}
}

func parseRequestBodyOperator(input string) (*RequestBodyOperator, error) {
	vals := map[string]RequestBodyOperator{
		"any":                RequestBodyOperatorAny,
		"beginswith":         RequestBodyOperatorBeginsWith,
		"contains":           RequestBodyOperatorContains,
		"endswith":           RequestBodyOperatorEndsWith,
		"equal":              RequestBodyOperatorEqual,
		"greaterthan":        RequestBodyOperatorGreaterThan,
		"greaterthanorequal": RequestBodyOperatorGreaterThanOrEqual,
		"lessthan":           RequestBodyOperatorLessThan,
		"lessthanorequal":    RequestBodyOperatorLessThanOrEqual,
		"regex":              RequestBodyOperatorRegEx,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := RequestBodyOperator(input)
	return &out, nil
}

type RequestHeaderOperator string

const (
	RequestHeaderOperatorAny                RequestHeaderOperator = "Any"
	RequestHeaderOperatorBeginsWith         RequestHeaderOperator = "BeginsWith"
	RequestHeaderOperatorContains           RequestHeaderOperator = "Contains"
	RequestHeaderOperatorEndsWith           RequestHeaderOperator = "EndsWith"
	RequestHeaderOperatorEqual              RequestHeaderOperator = "Equal"
	RequestHeaderOperatorGreaterThan        RequestHeaderOperator = "GreaterThan"
	RequestHeaderOperatorGreaterThanOrEqual RequestHeaderOperator = "GreaterThanOrEqual"
	RequestHeaderOperatorLessThan           RequestHeaderOperator = "LessThan"
	RequestHeaderOperatorLessThanOrEqual    RequestHeaderOperator = "LessThanOrEqual"
	RequestHeaderOperatorRegEx              RequestHeaderOperator = "RegEx"
)

func PossibleValuesForRequestHeaderOperator() []string {
	return []string{
		string(RequestHeaderOperatorAny),
		string(RequestHeaderOperatorBeginsWith),
		string(RequestHeaderOperatorContains),
		string(RequestHeaderOperatorEndsWith),
		string(RequestHeaderOperatorEqual),
		string(RequestHeaderOperatorGreaterThan),
		string(RequestHeaderOperatorGreaterThanOrEqual),
		string(RequestHeaderOperatorLessThan),
		string(RequestHeaderOperatorLessThanOrEqual),
		string(RequestHeaderOperatorRegEx),
	}
}

func parseRequestHeaderOperator(input string) (*RequestHeaderOperator, error) {
	vals := map[string]RequestHeaderOperator{
		"any":                RequestHeaderOperatorAny,
		"beginswith":         RequestHeaderOperatorBeginsWith,
		"contains":           RequestHeaderOperatorContains,
		"endswith":           RequestHeaderOperatorEndsWith,
		"equal":              RequestHeaderOperatorEqual,
		"greaterthan":        RequestHeaderOperatorGreaterThan,
		"greaterthanorequal": RequestHeaderOperatorGreaterThanOrEqual,
		"lessthan":           RequestHeaderOperatorLessThan,
		"lessthanorequal":    RequestHeaderOperatorLessThanOrEqual,
		"regex":              RequestHeaderOperatorRegEx,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := RequestHeaderOperator(input)
	return &out, nil
}

type RequestMethodOperator string

const (
	RequestMethodOperatorEqual RequestMethodOperator = "Equal"
)

func PossibleValuesForRequestMethodOperator() []string {
	return []string{
		string(RequestMethodOperatorEqual),
	}
}

func parseRequestMethodOperator(input string) (*RequestMethodOperator, error) {
	vals := map[string]RequestMethodOperator{
		"equal": RequestMethodOperatorEqual,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := RequestMethodOperator(input)
	return &out, nil
}

type RequestUriOperator string

const (
	RequestUriOperatorAny                RequestUriOperator = "Any"
	RequestUriOperatorBeginsWith         RequestUriOperator = "BeginsWith"
	RequestUriOperatorContains           RequestUriOperator = "Contains"
	RequestUriOperatorEndsWith           RequestUriOperator = "EndsWith"
	RequestUriOperatorEqual              RequestUriOperator = "Equal"
	RequestUriOperatorGreaterThan        RequestUriOperator = "GreaterThan"
	RequestUriOperatorGreaterThanOrEqual RequestUriOperator = "GreaterThanOrEqual"
	RequestUriOperatorLessThan           RequestUriOperator = "LessThan"
	RequestUriOperatorLessThanOrEqual    RequestUriOperator = "LessThanOrEqual"
	RequestUriOperatorRegEx              RequestUriOperator = "RegEx"
)

func PossibleValuesForRequestUriOperator() []string {
	return []string{
		string(RequestUriOperatorAny),
		string(RequestUriOperatorBeginsWith),
		string(RequestUriOperatorContains),
		string(RequestUriOperatorEndsWith),
		string(RequestUriOperatorEqual),
		string(RequestUriOperatorGreaterThan),
		string(RequestUriOperatorGreaterThanOrEqual),
		string(RequestUriOperatorLessThan),
		string(RequestUriOperatorLessThanOrEqual),
		string(RequestUriOperatorRegEx),
	}
}

func parseRequestUriOperator(input string) (*RequestUriOperator, error) {
	vals := map[string]RequestUriOperator{
		"any":                RequestUriOperatorAny,
		"beginswith":         RequestUriOperatorBeginsWith,
		"contains":           RequestUriOperatorContains,
		"endswith":           RequestUriOperatorEndsWith,
		"equal":              RequestUriOperatorEqual,
		"greaterthan":        RequestUriOperatorGreaterThan,
		"greaterthanorequal": RequestUriOperatorGreaterThanOrEqual,
		"lessthan":           RequestUriOperatorLessThan,
		"lessthanorequal":    RequestUriOperatorLessThanOrEqual,
		"regex":              RequestUriOperatorRegEx,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := RequestUriOperator(input)
	return &out, nil
}

type ResponseBasedDetectedErrorTypes string

const (
	ResponseBasedDetectedErrorTypesNone             ResponseBasedDetectedErrorTypes = "None"
	ResponseBasedDetectedErrorTypesTcpAndHttpErrors ResponseBasedDetectedErrorTypes = "TcpAndHttpErrors"
	ResponseBasedDetectedErrorTypesTcpErrorsOnly    ResponseBasedDetectedErrorTypes = "TcpErrorsOnly"
)

func PossibleValuesForResponseBasedDetectedErrorTypes() []string {
	return []string{
		string(ResponseBasedDetectedErrorTypesNone),
		string(ResponseBasedDetectedErrorTypesTcpAndHttpErrors),
		string(ResponseBasedDetectedErrorTypesTcpErrorsOnly),
	}
}

func parseResponseBasedDetectedErrorTypes(input string) (*ResponseBasedDetectedErrorTypes, error) {
	vals := map[string]ResponseBasedDetectedErrorTypes{
		"none":             ResponseBasedDetectedErrorTypesNone,
		"tcpandhttperrors": ResponseBasedDetectedErrorTypesTcpAndHttpErrors,
		"tcperrorsonly":    ResponseBasedDetectedErrorTypesTcpErrorsOnly,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := ResponseBasedDetectedErrorTypes(input)
	return &out, nil
}

type ServerPortOperator string

const (
	ServerPortOperatorAny                ServerPortOperator = "Any"
	ServerPortOperatorBeginsWith         ServerPortOperator = "BeginsWith"
	ServerPortOperatorContains           ServerPortOperator = "Contains"
	ServerPortOperatorEndsWith           ServerPortOperator = "EndsWith"
	ServerPortOperatorEqual              ServerPortOperator = "Equal"
	ServerPortOperatorGreaterThan        ServerPortOperator = "GreaterThan"
	ServerPortOperatorGreaterThanOrEqual ServerPortOperator = "GreaterThanOrEqual"
	ServerPortOperatorLessThan           ServerPortOperator = "LessThan"
	ServerPortOperatorLessThanOrEqual    ServerPortOperator = "LessThanOrEqual"
	ServerPortOperatorRegEx              ServerPortOperator = "RegEx"
)

func PossibleValuesForServerPortOperator() []string {
	return []string{
		string(ServerPortOperatorAny),
		string(ServerPortOperatorBeginsWith),
		string(ServerPortOperatorContains),
		string(ServerPortOperatorEndsWith),
		string(ServerPortOperatorEqual),
		string(ServerPortOperatorGreaterThan),
		string(ServerPortOperatorGreaterThanOrEqual),
		string(ServerPortOperatorLessThan),
		string(ServerPortOperatorLessThanOrEqual),
		string(ServerPortOperatorRegEx),
	}
}

func parseServerPortOperator(input string) (*ServerPortOperator, error) {
	vals := map[string]ServerPortOperator{
		"any":                ServerPortOperatorAny,
		"beginswith":         ServerPortOperatorBeginsWith,
		"contains":           ServerPortOperatorContains,
		"endswith":           ServerPortOperatorEndsWith,
		"equal":              ServerPortOperatorEqual,
		"greaterthan":        ServerPortOperatorGreaterThan,
		"greaterthanorequal": ServerPortOperatorGreaterThanOrEqual,
		"lessthan":           ServerPortOperatorLessThan,
		"lessthanorequal":    ServerPortOperatorLessThanOrEqual,
		"regex":              ServerPortOperatorRegEx,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := ServerPortOperator(input)
	return &out, nil
}

type SocketAddrOperator string

const (
	SocketAddrOperatorAny     SocketAddrOperator = "Any"
	SocketAddrOperatorIPMatch SocketAddrOperator = "IPMatch"
)

func PossibleValuesForSocketAddrOperator() []string {
	return []string{
		string(SocketAddrOperatorAny),
		string(SocketAddrOperatorIPMatch),
	}
}

func parseSocketAddrOperator(input string) (*SocketAddrOperator, error) {
	vals := map[string]SocketAddrOperator{
		"any":     SocketAddrOperatorAny,
		"ipmatch": SocketAddrOperatorIPMatch,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := SocketAddrOperator(input)
	return &out, nil
}

type SslProtocol string

const (
	SslProtocolTLSvOne         SslProtocol = "TLSv1"
	SslProtocolTLSvOnePointOne SslProtocol = "TLSv1.1"
	SslProtocolTLSvOnePointTwo SslProtocol = "TLSv1.2"
)

func PossibleValuesForSslProtocol() []string {
	return []string{
		string(SslProtocolTLSvOne),
		string(SslProtocolTLSvOnePointOne),
		string(SslProtocolTLSvOnePointTwo),
	}
}

func parseSslProtocol(input string) (*SslProtocol, error) {
	vals := map[string]SslProtocol{
		"tlsv1":   SslProtocolTLSvOne,
		"tlsv1.1": SslProtocolTLSvOnePointOne,
		"tlsv1.2": SslProtocolTLSvOnePointTwo,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := SslProtocol(input)
	return &out, nil
}

type SslProtocolOperator string

const (
	SslProtocolOperatorEqual SslProtocolOperator = "Equal"
)

func PossibleValuesForSslProtocolOperator() []string {
	return []string{
		string(SslProtocolOperatorEqual),
	}
}

func parseSslProtocolOperator(input string) (*SslProtocolOperator, error) {
	vals := map[string]SslProtocolOperator{
		"equal": SslProtocolOperatorEqual,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := SslProtocolOperator(input)
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

type UrlFileExtensionOperator string

const (
	UrlFileExtensionOperatorAny                UrlFileExtensionOperator = "Any"
	UrlFileExtensionOperatorBeginsWith         UrlFileExtensionOperator = "BeginsWith"
	UrlFileExtensionOperatorContains           UrlFileExtensionOperator = "Contains"
	UrlFileExtensionOperatorEndsWith           UrlFileExtensionOperator = "EndsWith"
	UrlFileExtensionOperatorEqual              UrlFileExtensionOperator = "Equal"
	UrlFileExtensionOperatorGreaterThan        UrlFileExtensionOperator = "GreaterThan"
	UrlFileExtensionOperatorGreaterThanOrEqual UrlFileExtensionOperator = "GreaterThanOrEqual"
	UrlFileExtensionOperatorLessThan           UrlFileExtensionOperator = "LessThan"
	UrlFileExtensionOperatorLessThanOrEqual    UrlFileExtensionOperator = "LessThanOrEqual"
	UrlFileExtensionOperatorRegEx              UrlFileExtensionOperator = "RegEx"
)

func PossibleValuesForUrlFileExtensionOperator() []string {
	return []string{
		string(UrlFileExtensionOperatorAny),
		string(UrlFileExtensionOperatorBeginsWith),
		string(UrlFileExtensionOperatorContains),
		string(UrlFileExtensionOperatorEndsWith),
		string(UrlFileExtensionOperatorEqual),
		string(UrlFileExtensionOperatorGreaterThan),
		string(UrlFileExtensionOperatorGreaterThanOrEqual),
		string(UrlFileExtensionOperatorLessThan),
		string(UrlFileExtensionOperatorLessThanOrEqual),
		string(UrlFileExtensionOperatorRegEx),
	}
}

func parseUrlFileExtensionOperator(input string) (*UrlFileExtensionOperator, error) {
	vals := map[string]UrlFileExtensionOperator{
		"any":                UrlFileExtensionOperatorAny,
		"beginswith":         UrlFileExtensionOperatorBeginsWith,
		"contains":           UrlFileExtensionOperatorContains,
		"endswith":           UrlFileExtensionOperatorEndsWith,
		"equal":              UrlFileExtensionOperatorEqual,
		"greaterthan":        UrlFileExtensionOperatorGreaterThan,
		"greaterthanorequal": UrlFileExtensionOperatorGreaterThanOrEqual,
		"lessthan":           UrlFileExtensionOperatorLessThan,
		"lessthanorequal":    UrlFileExtensionOperatorLessThanOrEqual,
		"regex":              UrlFileExtensionOperatorRegEx,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := UrlFileExtensionOperator(input)
	return &out, nil
}

type UrlFileNameOperator string

const (
	UrlFileNameOperatorAny                UrlFileNameOperator = "Any"
	UrlFileNameOperatorBeginsWith         UrlFileNameOperator = "BeginsWith"
	UrlFileNameOperatorContains           UrlFileNameOperator = "Contains"
	UrlFileNameOperatorEndsWith           UrlFileNameOperator = "EndsWith"
	UrlFileNameOperatorEqual              UrlFileNameOperator = "Equal"
	UrlFileNameOperatorGreaterThan        UrlFileNameOperator = "GreaterThan"
	UrlFileNameOperatorGreaterThanOrEqual UrlFileNameOperator = "GreaterThanOrEqual"
	UrlFileNameOperatorLessThan           UrlFileNameOperator = "LessThan"
	UrlFileNameOperatorLessThanOrEqual    UrlFileNameOperator = "LessThanOrEqual"
	UrlFileNameOperatorRegEx              UrlFileNameOperator = "RegEx"
)

func PossibleValuesForUrlFileNameOperator() []string {
	return []string{
		string(UrlFileNameOperatorAny),
		string(UrlFileNameOperatorBeginsWith),
		string(UrlFileNameOperatorContains),
		string(UrlFileNameOperatorEndsWith),
		string(UrlFileNameOperatorEqual),
		string(UrlFileNameOperatorGreaterThan),
		string(UrlFileNameOperatorGreaterThanOrEqual),
		string(UrlFileNameOperatorLessThan),
		string(UrlFileNameOperatorLessThanOrEqual),
		string(UrlFileNameOperatorRegEx),
	}
}

func parseUrlFileNameOperator(input string) (*UrlFileNameOperator, error) {
	vals := map[string]UrlFileNameOperator{
		"any":                UrlFileNameOperatorAny,
		"beginswith":         UrlFileNameOperatorBeginsWith,
		"contains":           UrlFileNameOperatorContains,
		"endswith":           UrlFileNameOperatorEndsWith,
		"equal":              UrlFileNameOperatorEqual,
		"greaterthan":        UrlFileNameOperatorGreaterThan,
		"greaterthanorequal": UrlFileNameOperatorGreaterThanOrEqual,
		"lessthan":           UrlFileNameOperatorLessThan,
		"lessthanorequal":    UrlFileNameOperatorLessThanOrEqual,
		"regex":              UrlFileNameOperatorRegEx,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := UrlFileNameOperator(input)
	return &out, nil
}

type UrlPathOperator string

const (
	UrlPathOperatorAny                UrlPathOperator = "Any"
	UrlPathOperatorBeginsWith         UrlPathOperator = "BeginsWith"
	UrlPathOperatorContains           UrlPathOperator = "Contains"
	UrlPathOperatorEndsWith           UrlPathOperator = "EndsWith"
	UrlPathOperatorEqual              UrlPathOperator = "Equal"
	UrlPathOperatorGreaterThan        UrlPathOperator = "GreaterThan"
	UrlPathOperatorGreaterThanOrEqual UrlPathOperator = "GreaterThanOrEqual"
	UrlPathOperatorLessThan           UrlPathOperator = "LessThan"
	UrlPathOperatorLessThanOrEqual    UrlPathOperator = "LessThanOrEqual"
	UrlPathOperatorRegEx              UrlPathOperator = "RegEx"
	UrlPathOperatorWildcard           UrlPathOperator = "Wildcard"
)

func PossibleValuesForUrlPathOperator() []string {
	return []string{
		string(UrlPathOperatorAny),
		string(UrlPathOperatorBeginsWith),
		string(UrlPathOperatorContains),
		string(UrlPathOperatorEndsWith),
		string(UrlPathOperatorEqual),
		string(UrlPathOperatorGreaterThan),
		string(UrlPathOperatorGreaterThanOrEqual),
		string(UrlPathOperatorLessThan),
		string(UrlPathOperatorLessThanOrEqual),
		string(UrlPathOperatorRegEx),
		string(UrlPathOperatorWildcard),
	}
}

func parseUrlPathOperator(input string) (*UrlPathOperator, error) {
	vals := map[string]UrlPathOperator{
		"any":                UrlPathOperatorAny,
		"beginswith":         UrlPathOperatorBeginsWith,
		"contains":           UrlPathOperatorContains,
		"endswith":           UrlPathOperatorEndsWith,
		"equal":              UrlPathOperatorEqual,
		"greaterthan":        UrlPathOperatorGreaterThan,
		"greaterthanorequal": UrlPathOperatorGreaterThanOrEqual,
		"lessthan":           UrlPathOperatorLessThan,
		"lessthanorequal":    UrlPathOperatorLessThanOrEqual,
		"regex":              UrlPathOperatorRegEx,
		"wildcard":           UrlPathOperatorWildcard,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := UrlPathOperator(input)
	return &out, nil
}
