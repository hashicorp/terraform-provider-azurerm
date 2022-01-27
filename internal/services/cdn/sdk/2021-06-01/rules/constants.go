package rules

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
