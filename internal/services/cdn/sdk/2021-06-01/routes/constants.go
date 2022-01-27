package routes

import "strings"

type AFDEndpointProtocols string

const (
	AFDEndpointProtocolsHttp  AFDEndpointProtocols = "Http"
	AFDEndpointProtocolsHttps AFDEndpointProtocols = "Https"
)

func PossibleValuesForAFDEndpointProtocols() []string {
	return []string{
		string(AFDEndpointProtocolsHttp),
		string(AFDEndpointProtocolsHttps),
	}
}

func parseAFDEndpointProtocols(input string) (*AFDEndpointProtocols, error) {
	vals := map[string]AFDEndpointProtocols{
		"http":  AFDEndpointProtocolsHttp,
		"https": AFDEndpointProtocolsHttps,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := AFDEndpointProtocols(input)
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

type AfdQueryStringCachingBehavior string

const (
	AfdQueryStringCachingBehaviorIgnoreQueryString            AfdQueryStringCachingBehavior = "IgnoreQueryString"
	AfdQueryStringCachingBehaviorIgnoreSpecifiedQueryStrings  AfdQueryStringCachingBehavior = "IgnoreSpecifiedQueryStrings"
	AfdQueryStringCachingBehaviorIncludeSpecifiedQueryStrings AfdQueryStringCachingBehavior = "IncludeSpecifiedQueryStrings"
	AfdQueryStringCachingBehaviorUseQueryString               AfdQueryStringCachingBehavior = "UseQueryString"
)

func PossibleValuesForAfdQueryStringCachingBehavior() []string {
	return []string{
		string(AfdQueryStringCachingBehaviorIgnoreQueryString),
		string(AfdQueryStringCachingBehaviorIgnoreSpecifiedQueryStrings),
		string(AfdQueryStringCachingBehaviorIncludeSpecifiedQueryStrings),
		string(AfdQueryStringCachingBehaviorUseQueryString),
	}
}

func parseAfdQueryStringCachingBehavior(input string) (*AfdQueryStringCachingBehavior, error) {
	vals := map[string]AfdQueryStringCachingBehavior{
		"ignorequerystring":            AfdQueryStringCachingBehaviorIgnoreQueryString,
		"ignorespecifiedquerystrings":  AfdQueryStringCachingBehaviorIgnoreSpecifiedQueryStrings,
		"includespecifiedquerystrings": AfdQueryStringCachingBehaviorIncludeSpecifiedQueryStrings,
		"usequerystring":               AfdQueryStringCachingBehaviorUseQueryString,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := AfdQueryStringCachingBehavior(input)
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

type EnabledState string

const (
	EnabledStateDisabled EnabledState = "Disabled"
	EnabledStateEnabled  EnabledState = "Enabled"
)

func PossibleValuesForEnabledState() []string {
	return []string{
		string(EnabledStateDisabled),
		string(EnabledStateEnabled),
	}
}

func parseEnabledState(input string) (*EnabledState, error) {
	vals := map[string]EnabledState{
		"disabled": EnabledStateDisabled,
		"enabled":  EnabledStateEnabled,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := EnabledState(input)
	return &out, nil
}

type ForwardingProtocol string

const (
	ForwardingProtocolHttpOnly     ForwardingProtocol = "HttpOnly"
	ForwardingProtocolHttpsOnly    ForwardingProtocol = "HttpsOnly"
	ForwardingProtocolMatchRequest ForwardingProtocol = "MatchRequest"
)

func PossibleValuesForForwardingProtocol() []string {
	return []string{
		string(ForwardingProtocolHttpOnly),
		string(ForwardingProtocolHttpsOnly),
		string(ForwardingProtocolMatchRequest),
	}
}

func parseForwardingProtocol(input string) (*ForwardingProtocol, error) {
	vals := map[string]ForwardingProtocol{
		"httponly":     ForwardingProtocolHttpOnly,
		"httpsonly":    ForwardingProtocolHttpsOnly,
		"matchrequest": ForwardingProtocolMatchRequest,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := ForwardingProtocol(input)
	return &out, nil
}

type HttpsRedirect string

const (
	HttpsRedirectDisabled HttpsRedirect = "Disabled"
	HttpsRedirectEnabled  HttpsRedirect = "Enabled"
)

func PossibleValuesForHttpsRedirect() []string {
	return []string{
		string(HttpsRedirectDisabled),
		string(HttpsRedirectEnabled),
	}
}

func parseHttpsRedirect(input string) (*HttpsRedirect, error) {
	vals := map[string]HttpsRedirect{
		"disabled": HttpsRedirectDisabled,
		"enabled":  HttpsRedirectEnabled,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := HttpsRedirect(input)
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

type LinkToDefaultDomain string

const (
	LinkToDefaultDomainDisabled LinkToDefaultDomain = "Disabled"
	LinkToDefaultDomainEnabled  LinkToDefaultDomain = "Enabled"
)

func PossibleValuesForLinkToDefaultDomain() []string {
	return []string{
		string(LinkToDefaultDomainDisabled),
		string(LinkToDefaultDomainEnabled),
	}
}

func parseLinkToDefaultDomain(input string) (*LinkToDefaultDomain, error) {
	vals := map[string]LinkToDefaultDomain{
		"disabled": LinkToDefaultDomainDisabled,
		"enabled":  LinkToDefaultDomainEnabled,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := LinkToDefaultDomain(input)
	return &out, nil
}
