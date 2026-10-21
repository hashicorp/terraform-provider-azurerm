package rules

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

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

type Algorithm string

const (
	AlgorithmSHATwoFiveSix Algorithm = "SHA256"
)

func PossibleValuesForAlgorithm() []string {
	return []string{
		string(AlgorithmSHATwoFiveSix),
	}
}

func (s *Algorithm) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseAlgorithm(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseAlgorithm(input string) (*Algorithm, error) {
	vals := map[string]Algorithm{
		"sha256": AlgorithmSHATwoFiveSix,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := Algorithm(input)
	return &out, nil
}

type CacheBehavior string

const (
	CacheBehaviorBypassCache  CacheBehavior = "BypassCache"
	CacheBehaviorOverride     CacheBehavior = "Override"
	CacheBehaviorSetIfMissing CacheBehavior = "SetIfMissing"
)

func PossibleValuesForCacheBehavior() []string {
	return []string{
		string(CacheBehaviorBypassCache),
		string(CacheBehaviorOverride),
		string(CacheBehaviorSetIfMissing),
	}
}

func (s *CacheBehavior) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseCacheBehavior(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseCacheBehavior(input string) (*CacheBehavior, error) {
	vals := map[string]CacheBehavior{
		"bypasscache":  CacheBehaviorBypassCache,
		"override":     CacheBehaviorOverride,
		"setifmissing": CacheBehaviorSetIfMissing,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := CacheBehavior(input)
	return &out, nil
}

type CacheType string

const (
	CacheTypeAll CacheType = "All"
)

func PossibleValuesForCacheType() []string {
	return []string{
		string(CacheTypeAll),
	}
}

func (s *CacheType) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseCacheType(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseCacheType(input string) (*CacheType, error) {
	vals := map[string]CacheType{
		"all": CacheTypeAll,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := CacheType(input)
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

func (s *ClientPortOperator) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseClientPortOperator(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
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

func (s *CookiesOperator) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseCookiesOperator(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
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

type DeliveryRuleActionName string

const (
	DeliveryRuleActionNameCacheExpiration            DeliveryRuleActionName = "CacheExpiration"
	DeliveryRuleActionNameCacheKeyQueryString        DeliveryRuleActionName = "CacheKeyQueryString"
	DeliveryRuleActionNameModifyRequestHeader        DeliveryRuleActionName = "ModifyRequestHeader"
	DeliveryRuleActionNameModifyResponseHeader       DeliveryRuleActionName = "ModifyResponseHeader"
	DeliveryRuleActionNameOriginGroupOverride        DeliveryRuleActionName = "OriginGroupOverride"
	DeliveryRuleActionNameRouteConfigurationOverride DeliveryRuleActionName = "RouteConfigurationOverride"
	DeliveryRuleActionNameURLRedirect                DeliveryRuleActionName = "UrlRedirect"
	DeliveryRuleActionNameURLRewrite                 DeliveryRuleActionName = "UrlRewrite"
	DeliveryRuleActionNameURLSigning                 DeliveryRuleActionName = "UrlSigning"
)

func PossibleValuesForDeliveryRuleActionName() []string {
	return []string{
		string(DeliveryRuleActionNameCacheExpiration),
		string(DeliveryRuleActionNameCacheKeyQueryString),
		string(DeliveryRuleActionNameModifyRequestHeader),
		string(DeliveryRuleActionNameModifyResponseHeader),
		string(DeliveryRuleActionNameOriginGroupOverride),
		string(DeliveryRuleActionNameRouteConfigurationOverride),
		string(DeliveryRuleActionNameURLRedirect),
		string(DeliveryRuleActionNameURLRewrite),
		string(DeliveryRuleActionNameURLSigning),
	}
}

func (s *DeliveryRuleActionName) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseDeliveryRuleActionName(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseDeliveryRuleActionName(input string) (*DeliveryRuleActionName, error) {
	vals := map[string]DeliveryRuleActionName{
		"cacheexpiration":            DeliveryRuleActionNameCacheExpiration,
		"cachekeyquerystring":        DeliveryRuleActionNameCacheKeyQueryString,
		"modifyrequestheader":        DeliveryRuleActionNameModifyRequestHeader,
		"modifyresponseheader":       DeliveryRuleActionNameModifyResponseHeader,
		"origingroupoverride":        DeliveryRuleActionNameOriginGroupOverride,
		"routeconfigurationoverride": DeliveryRuleActionNameRouteConfigurationOverride,
		"urlredirect":                DeliveryRuleActionNameURLRedirect,
		"urlrewrite":                 DeliveryRuleActionNameURLRewrite,
		"urlsigning":                 DeliveryRuleActionNameURLSigning,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := DeliveryRuleActionName(input)
	return &out, nil
}

type DeliveryRuleActionParametersType string

const (
	DeliveryRuleActionParametersTypeDeliveryRuleCacheExpirationActionParameters             DeliveryRuleActionParametersType = "DeliveryRuleCacheExpirationActionParameters"
	DeliveryRuleActionParametersTypeDeliveryRuleCacheKeyQueryStringBehaviorActionParameters DeliveryRuleActionParametersType = "DeliveryRuleCacheKeyQueryStringBehaviorActionParameters"
	DeliveryRuleActionParametersTypeDeliveryRuleHeaderActionParameters                      DeliveryRuleActionParametersType = "DeliveryRuleHeaderActionParameters"
	DeliveryRuleActionParametersTypeDeliveryRuleOriginGroupOverrideActionParameters         DeliveryRuleActionParametersType = "DeliveryRuleOriginGroupOverrideActionParameters"
	DeliveryRuleActionParametersTypeDeliveryRuleRouteConfigurationOverrideActionParameters  DeliveryRuleActionParametersType = "DeliveryRuleRouteConfigurationOverrideActionParameters"
	DeliveryRuleActionParametersTypeDeliveryRuleURLRedirectActionParameters                 DeliveryRuleActionParametersType = "DeliveryRuleUrlRedirectActionParameters"
	DeliveryRuleActionParametersTypeDeliveryRuleURLRewriteActionParameters                  DeliveryRuleActionParametersType = "DeliveryRuleUrlRewriteActionParameters"
	DeliveryRuleActionParametersTypeDeliveryRuleURLSigningActionParameters                  DeliveryRuleActionParametersType = "DeliveryRuleUrlSigningActionParameters"
)

func PossibleValuesForDeliveryRuleActionParametersType() []string {
	return []string{
		string(DeliveryRuleActionParametersTypeDeliveryRuleCacheExpirationActionParameters),
		string(DeliveryRuleActionParametersTypeDeliveryRuleCacheKeyQueryStringBehaviorActionParameters),
		string(DeliveryRuleActionParametersTypeDeliveryRuleHeaderActionParameters),
		string(DeliveryRuleActionParametersTypeDeliveryRuleOriginGroupOverrideActionParameters),
		string(DeliveryRuleActionParametersTypeDeliveryRuleRouteConfigurationOverrideActionParameters),
		string(DeliveryRuleActionParametersTypeDeliveryRuleURLRedirectActionParameters),
		string(DeliveryRuleActionParametersTypeDeliveryRuleURLRewriteActionParameters),
		string(DeliveryRuleActionParametersTypeDeliveryRuleURLSigningActionParameters),
	}
}

func (s *DeliveryRuleActionParametersType) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseDeliveryRuleActionParametersType(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseDeliveryRuleActionParametersType(input string) (*DeliveryRuleActionParametersType, error) {
	vals := map[string]DeliveryRuleActionParametersType{
		"deliveryrulecacheexpirationactionparameters":             DeliveryRuleActionParametersTypeDeliveryRuleCacheExpirationActionParameters,
		"deliveryrulecachekeyquerystringbehavioractionparameters": DeliveryRuleActionParametersTypeDeliveryRuleCacheKeyQueryStringBehaviorActionParameters,
		"deliveryruleheaderactionparameters":                      DeliveryRuleActionParametersTypeDeliveryRuleHeaderActionParameters,
		"deliveryruleorigingroupoverrideactionparameters":         DeliveryRuleActionParametersTypeDeliveryRuleOriginGroupOverrideActionParameters,
		"deliveryrulerouteconfigurationoverrideactionparameters":  DeliveryRuleActionParametersTypeDeliveryRuleRouteConfigurationOverrideActionParameters,
		"deliveryruleurlredirectactionparameters":                 DeliveryRuleActionParametersTypeDeliveryRuleURLRedirectActionParameters,
		"deliveryruleurlrewriteactionparameters":                  DeliveryRuleActionParametersTypeDeliveryRuleURLRewriteActionParameters,
		"deliveryruleurlsigningactionparameters":                  DeliveryRuleActionParametersTypeDeliveryRuleURLSigningActionParameters,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := DeliveryRuleActionParametersType(input)
	return &out, nil
}

type DeliveryRuleConditionParametersType string

const (
	DeliveryRuleConditionParametersTypeDeliveryRuleClientPortConditionParameters            DeliveryRuleConditionParametersType = "DeliveryRuleClientPortConditionParameters"
	DeliveryRuleConditionParametersTypeDeliveryRuleCookiesConditionParameters               DeliveryRuleConditionParametersType = "DeliveryRuleCookiesConditionParameters"
	DeliveryRuleConditionParametersTypeDeliveryRuleHTTPVersionConditionParameters           DeliveryRuleConditionParametersType = "DeliveryRuleHttpVersionConditionParameters"
	DeliveryRuleConditionParametersTypeDeliveryRuleHostNameConditionParameters              DeliveryRuleConditionParametersType = "DeliveryRuleHostNameConditionParameters"
	DeliveryRuleConditionParametersTypeDeliveryRuleIsDeviceConditionParameters              DeliveryRuleConditionParametersType = "DeliveryRuleIsDeviceConditionParameters"
	DeliveryRuleConditionParametersTypeDeliveryRulePostArgsConditionParameters              DeliveryRuleConditionParametersType = "DeliveryRulePostArgsConditionParameters"
	DeliveryRuleConditionParametersTypeDeliveryRuleQueryStringConditionParameters           DeliveryRuleConditionParametersType = "DeliveryRuleQueryStringConditionParameters"
	DeliveryRuleConditionParametersTypeDeliveryRuleRemoteAddressConditionParameters         DeliveryRuleConditionParametersType = "DeliveryRuleRemoteAddressConditionParameters"
	DeliveryRuleConditionParametersTypeDeliveryRuleRequestBodyConditionParameters           DeliveryRuleConditionParametersType = "DeliveryRuleRequestBodyConditionParameters"
	DeliveryRuleConditionParametersTypeDeliveryRuleRequestHeaderConditionParameters         DeliveryRuleConditionParametersType = "DeliveryRuleRequestHeaderConditionParameters"
	DeliveryRuleConditionParametersTypeDeliveryRuleRequestMethodConditionParameters         DeliveryRuleConditionParametersType = "DeliveryRuleRequestMethodConditionParameters"
	DeliveryRuleConditionParametersTypeDeliveryRuleRequestSchemeConditionParameters         DeliveryRuleConditionParametersType = "DeliveryRuleRequestSchemeConditionParameters"
	DeliveryRuleConditionParametersTypeDeliveryRuleRequestUriConditionParameters            DeliveryRuleConditionParametersType = "DeliveryRuleRequestUriConditionParameters"
	DeliveryRuleConditionParametersTypeDeliveryRuleServerPortConditionParameters            DeliveryRuleConditionParametersType = "DeliveryRuleServerPortConditionParameters"
	DeliveryRuleConditionParametersTypeDeliveryRuleSocketAddrConditionParameters            DeliveryRuleConditionParametersType = "DeliveryRuleSocketAddrConditionParameters"
	DeliveryRuleConditionParametersTypeDeliveryRuleSslProtocolConditionParameters           DeliveryRuleConditionParametersType = "DeliveryRuleSslProtocolConditionParameters"
	DeliveryRuleConditionParametersTypeDeliveryRuleURLFileExtensionMatchConditionParameters DeliveryRuleConditionParametersType = "DeliveryRuleUrlFileExtensionMatchConditionParameters"
	DeliveryRuleConditionParametersTypeDeliveryRuleURLFilenameConditionParameters           DeliveryRuleConditionParametersType = "DeliveryRuleUrlFilenameConditionParameters"
	DeliveryRuleConditionParametersTypeDeliveryRuleURLPathMatchConditionParameters          DeliveryRuleConditionParametersType = "DeliveryRuleUrlPathMatchConditionParameters"
)

func PossibleValuesForDeliveryRuleConditionParametersType() []string {
	return []string{
		string(DeliveryRuleConditionParametersTypeDeliveryRuleClientPortConditionParameters),
		string(DeliveryRuleConditionParametersTypeDeliveryRuleCookiesConditionParameters),
		string(DeliveryRuleConditionParametersTypeDeliveryRuleHTTPVersionConditionParameters),
		string(DeliveryRuleConditionParametersTypeDeliveryRuleHostNameConditionParameters),
		string(DeliveryRuleConditionParametersTypeDeliveryRuleIsDeviceConditionParameters),
		string(DeliveryRuleConditionParametersTypeDeliveryRulePostArgsConditionParameters),
		string(DeliveryRuleConditionParametersTypeDeliveryRuleQueryStringConditionParameters),
		string(DeliveryRuleConditionParametersTypeDeliveryRuleRemoteAddressConditionParameters),
		string(DeliveryRuleConditionParametersTypeDeliveryRuleRequestBodyConditionParameters),
		string(DeliveryRuleConditionParametersTypeDeliveryRuleRequestHeaderConditionParameters),
		string(DeliveryRuleConditionParametersTypeDeliveryRuleRequestMethodConditionParameters),
		string(DeliveryRuleConditionParametersTypeDeliveryRuleRequestSchemeConditionParameters),
		string(DeliveryRuleConditionParametersTypeDeliveryRuleRequestUriConditionParameters),
		string(DeliveryRuleConditionParametersTypeDeliveryRuleServerPortConditionParameters),
		string(DeliveryRuleConditionParametersTypeDeliveryRuleSocketAddrConditionParameters),
		string(DeliveryRuleConditionParametersTypeDeliveryRuleSslProtocolConditionParameters),
		string(DeliveryRuleConditionParametersTypeDeliveryRuleURLFileExtensionMatchConditionParameters),
		string(DeliveryRuleConditionParametersTypeDeliveryRuleURLFilenameConditionParameters),
		string(DeliveryRuleConditionParametersTypeDeliveryRuleURLPathMatchConditionParameters),
	}
}

func (s *DeliveryRuleConditionParametersType) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseDeliveryRuleConditionParametersType(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseDeliveryRuleConditionParametersType(input string) (*DeliveryRuleConditionParametersType, error) {
	vals := map[string]DeliveryRuleConditionParametersType{
		"deliveryruleclientportconditionparameters":            DeliveryRuleConditionParametersTypeDeliveryRuleClientPortConditionParameters,
		"deliveryrulecookiesconditionparameters":               DeliveryRuleConditionParametersTypeDeliveryRuleCookiesConditionParameters,
		"deliveryrulehttpversionconditionparameters":           DeliveryRuleConditionParametersTypeDeliveryRuleHTTPVersionConditionParameters,
		"deliveryrulehostnameconditionparameters":              DeliveryRuleConditionParametersTypeDeliveryRuleHostNameConditionParameters,
		"deliveryruleisdeviceconditionparameters":              DeliveryRuleConditionParametersTypeDeliveryRuleIsDeviceConditionParameters,
		"deliveryrulepostargsconditionparameters":              DeliveryRuleConditionParametersTypeDeliveryRulePostArgsConditionParameters,
		"deliveryrulequerystringconditionparameters":           DeliveryRuleConditionParametersTypeDeliveryRuleQueryStringConditionParameters,
		"deliveryruleremoteaddressconditionparameters":         DeliveryRuleConditionParametersTypeDeliveryRuleRemoteAddressConditionParameters,
		"deliveryrulerequestbodyconditionparameters":           DeliveryRuleConditionParametersTypeDeliveryRuleRequestBodyConditionParameters,
		"deliveryrulerequestheaderconditionparameters":         DeliveryRuleConditionParametersTypeDeliveryRuleRequestHeaderConditionParameters,
		"deliveryrulerequestmethodconditionparameters":         DeliveryRuleConditionParametersTypeDeliveryRuleRequestMethodConditionParameters,
		"deliveryrulerequestschemeconditionparameters":         DeliveryRuleConditionParametersTypeDeliveryRuleRequestSchemeConditionParameters,
		"deliveryrulerequesturiconditionparameters":            DeliveryRuleConditionParametersTypeDeliveryRuleRequestUriConditionParameters,
		"deliveryruleserverportconditionparameters":            DeliveryRuleConditionParametersTypeDeliveryRuleServerPortConditionParameters,
		"deliveryrulesocketaddrconditionparameters":            DeliveryRuleConditionParametersTypeDeliveryRuleSocketAddrConditionParameters,
		"deliveryrulesslprotocolconditionparameters":           DeliveryRuleConditionParametersTypeDeliveryRuleSslProtocolConditionParameters,
		"deliveryruleurlfileextensionmatchconditionparameters": DeliveryRuleConditionParametersTypeDeliveryRuleURLFileExtensionMatchConditionParameters,
		"deliveryruleurlfilenameconditionparameters":           DeliveryRuleConditionParametersTypeDeliveryRuleURLFilenameConditionParameters,
		"deliveryruleurlpathmatchconditionparameters":          DeliveryRuleConditionParametersTypeDeliveryRuleURLPathMatchConditionParameters,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := DeliveryRuleConditionParametersType(input)
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

type DestinationProtocol string

const (
	DestinationProtocolHTTP         DestinationProtocol = "Http"
	DestinationProtocolHTTPS        DestinationProtocol = "Https"
	DestinationProtocolMatchRequest DestinationProtocol = "MatchRequest"
)

func PossibleValuesForDestinationProtocol() []string {
	return []string{
		string(DestinationProtocolHTTP),
		string(DestinationProtocolHTTPS),
		string(DestinationProtocolMatchRequest),
	}
}

func (s *DestinationProtocol) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseDestinationProtocol(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseDestinationProtocol(input string) (*DestinationProtocol, error) {
	vals := map[string]DestinationProtocol{
		"http":         DestinationProtocolHTTP,
		"https":        DestinationProtocolHTTPS,
		"matchrequest": DestinationProtocolMatchRequest,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := DestinationProtocol(input)
	return &out, nil
}

type ForwardingProtocol string

const (
	ForwardingProtocolHTTPOnly     ForwardingProtocol = "HttpOnly"
	ForwardingProtocolHTTPSOnly    ForwardingProtocol = "HttpsOnly"
	ForwardingProtocolMatchRequest ForwardingProtocol = "MatchRequest"
)

func PossibleValuesForForwardingProtocol() []string {
	return []string{
		string(ForwardingProtocolHTTPOnly),
		string(ForwardingProtocolHTTPSOnly),
		string(ForwardingProtocolMatchRequest),
	}
}

func (s *ForwardingProtocol) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseForwardingProtocol(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseForwardingProtocol(input string) (*ForwardingProtocol, error) {
	vals := map[string]ForwardingProtocol{
		"httponly":     ForwardingProtocolHTTPOnly,
		"httpsonly":    ForwardingProtocolHTTPSOnly,
		"matchrequest": ForwardingProtocolMatchRequest,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := ForwardingProtocol(input)
	return &out, nil
}

type HTTPVersionOperator string

const (
	HTTPVersionOperatorEqual HTTPVersionOperator = "Equal"
)

func PossibleValuesForHTTPVersionOperator() []string {
	return []string{
		string(HTTPVersionOperatorEqual),
	}
}

func (s *HTTPVersionOperator) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseHTTPVersionOperator(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseHTTPVersionOperator(input string) (*HTTPVersionOperator, error) {
	vals := map[string]HTTPVersionOperator{
		"equal": HTTPVersionOperatorEqual,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := HTTPVersionOperator(input)
	return &out, nil
}

type HeaderAction string

const (
	HeaderActionAppend    HeaderAction = "Append"
	HeaderActionDelete    HeaderAction = "Delete"
	HeaderActionOverwrite HeaderAction = "Overwrite"
)

func PossibleValuesForHeaderAction() []string {
	return []string{
		string(HeaderActionAppend),
		string(HeaderActionDelete),
		string(HeaderActionOverwrite),
	}
}

func (s *HeaderAction) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseHeaderAction(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseHeaderAction(input string) (*HeaderAction, error) {
	vals := map[string]HeaderAction{
		"append":    HeaderActionAppend,
		"delete":    HeaderActionDelete,
		"overwrite": HeaderActionOverwrite,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := HeaderAction(input)
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

func (s *HostNameOperator) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseHostNameOperator(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
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

type IsDeviceMatchValue string

const (
	IsDeviceMatchValueDesktop IsDeviceMatchValue = "Desktop"
	IsDeviceMatchValueMobile  IsDeviceMatchValue = "Mobile"
)

func PossibleValuesForIsDeviceMatchValue() []string {
	return []string{
		string(IsDeviceMatchValueDesktop),
		string(IsDeviceMatchValueMobile),
	}
}

func (s *IsDeviceMatchValue) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseIsDeviceMatchValue(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseIsDeviceMatchValue(input string) (*IsDeviceMatchValue, error) {
	vals := map[string]IsDeviceMatchValue{
		"desktop": IsDeviceMatchValueDesktop,
		"mobile":  IsDeviceMatchValueMobile,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := IsDeviceMatchValue(input)
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

func (s *IsDeviceOperator) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseIsDeviceOperator(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
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

func (s *MatchProcessingBehavior) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseMatchProcessingBehavior(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
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
	MatchVariableHTTPVersion      MatchVariable = "HttpVersion"
	MatchVariableHostName         MatchVariable = "HostName"
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
	MatchVariableURLFileExtension MatchVariable = "UrlFileExtension"
	MatchVariableURLFileName      MatchVariable = "UrlFileName"
	MatchVariableURLPath          MatchVariable = "UrlPath"
)

func PossibleValuesForMatchVariable() []string {
	return []string{
		string(MatchVariableClientPort),
		string(MatchVariableCookies),
		string(MatchVariableHTTPVersion),
		string(MatchVariableHostName),
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
		string(MatchVariableURLFileExtension),
		string(MatchVariableURLFileName),
		string(MatchVariableURLPath),
	}
}

func (s *MatchVariable) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseMatchVariable(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseMatchVariable(input string) (*MatchVariable, error) {
	vals := map[string]MatchVariable{
		"clientport":       MatchVariableClientPort,
		"cookies":          MatchVariableCookies,
		"httpversion":      MatchVariableHTTPVersion,
		"hostname":         MatchVariableHostName,
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
		"urlfileextension": MatchVariableURLFileExtension,
		"urlfilename":      MatchVariableURLFileName,
		"urlpath":          MatchVariableURLPath,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := MatchVariable(input)
	return &out, nil
}

type ParamIndicator string

const (
	ParamIndicatorExpires   ParamIndicator = "Expires"
	ParamIndicatorKeyId     ParamIndicator = "KeyId"
	ParamIndicatorSignature ParamIndicator = "Signature"
)

func PossibleValuesForParamIndicator() []string {
	return []string{
		string(ParamIndicatorExpires),
		string(ParamIndicatorKeyId),
		string(ParamIndicatorSignature),
	}
}

func (s *ParamIndicator) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseParamIndicator(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseParamIndicator(input string) (*ParamIndicator, error) {
	vals := map[string]ParamIndicator{
		"expires":   ParamIndicatorExpires,
		"keyid":     ParamIndicatorKeyId,
		"signature": ParamIndicatorSignature,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := ParamIndicator(input)
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

func (s *PostArgsOperator) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parsePostArgsOperator(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
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

type QueryStringBehavior string

const (
	QueryStringBehaviorExclude    QueryStringBehavior = "Exclude"
	QueryStringBehaviorExcludeAll QueryStringBehavior = "ExcludeAll"
	QueryStringBehaviorInclude    QueryStringBehavior = "Include"
	QueryStringBehaviorIncludeAll QueryStringBehavior = "IncludeAll"
)

func PossibleValuesForQueryStringBehavior() []string {
	return []string{
		string(QueryStringBehaviorExclude),
		string(QueryStringBehaviorExcludeAll),
		string(QueryStringBehaviorInclude),
		string(QueryStringBehaviorIncludeAll),
	}
}

func (s *QueryStringBehavior) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseQueryStringBehavior(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseQueryStringBehavior(input string) (*QueryStringBehavior, error) {
	vals := map[string]QueryStringBehavior{
		"exclude":    QueryStringBehaviorExclude,
		"excludeall": QueryStringBehaviorExcludeAll,
		"include":    QueryStringBehaviorInclude,
		"includeall": QueryStringBehaviorIncludeAll,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := QueryStringBehavior(input)
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

func (s *QueryStringOperator) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseQueryStringOperator(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
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

type RedirectType string

const (
	RedirectTypeFound             RedirectType = "Found"
	RedirectTypeMoved             RedirectType = "Moved"
	RedirectTypePermanentRedirect RedirectType = "PermanentRedirect"
	RedirectTypeTemporaryRedirect RedirectType = "TemporaryRedirect"
)

func PossibleValuesForRedirectType() []string {
	return []string{
		string(RedirectTypeFound),
		string(RedirectTypeMoved),
		string(RedirectTypePermanentRedirect),
		string(RedirectTypeTemporaryRedirect),
	}
}

func (s *RedirectType) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseRedirectType(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseRedirectType(input string) (*RedirectType, error) {
	vals := map[string]RedirectType{
		"found":             RedirectTypeFound,
		"moved":             RedirectTypeMoved,
		"permanentredirect": RedirectTypePermanentRedirect,
		"temporaryredirect": RedirectTypeTemporaryRedirect,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := RedirectType(input)
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

func (s *RemoteAddressOperator) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseRemoteAddressOperator(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
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

func (s *RequestBodyOperator) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseRequestBodyOperator(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
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

func (s *RequestHeaderOperator) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseRequestHeaderOperator(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
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

type RequestMethodMatchValue string

const (
	RequestMethodMatchValueDELETE  RequestMethodMatchValue = "DELETE"
	RequestMethodMatchValueGET     RequestMethodMatchValue = "GET"
	RequestMethodMatchValueHEAD    RequestMethodMatchValue = "HEAD"
	RequestMethodMatchValueOPTIONS RequestMethodMatchValue = "OPTIONS"
	RequestMethodMatchValuePOST    RequestMethodMatchValue = "POST"
	RequestMethodMatchValuePUT     RequestMethodMatchValue = "PUT"
	RequestMethodMatchValueTRACE   RequestMethodMatchValue = "TRACE"
)

func PossibleValuesForRequestMethodMatchValue() []string {
	return []string{
		string(RequestMethodMatchValueDELETE),
		string(RequestMethodMatchValueGET),
		string(RequestMethodMatchValueHEAD),
		string(RequestMethodMatchValueOPTIONS),
		string(RequestMethodMatchValuePOST),
		string(RequestMethodMatchValuePUT),
		string(RequestMethodMatchValueTRACE),
	}
}

func (s *RequestMethodMatchValue) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseRequestMethodMatchValue(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseRequestMethodMatchValue(input string) (*RequestMethodMatchValue, error) {
	vals := map[string]RequestMethodMatchValue{
		"delete":  RequestMethodMatchValueDELETE,
		"get":     RequestMethodMatchValueGET,
		"head":    RequestMethodMatchValueHEAD,
		"options": RequestMethodMatchValueOPTIONS,
		"post":    RequestMethodMatchValuePOST,
		"put":     RequestMethodMatchValuePUT,
		"trace":   RequestMethodMatchValueTRACE,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := RequestMethodMatchValue(input)
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

func (s *RequestMethodOperator) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseRequestMethodOperator(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
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

type RequestSchemeMatchConditionParametersOperator string

const (
	RequestSchemeMatchConditionParametersOperatorEqual RequestSchemeMatchConditionParametersOperator = "Equal"
)

func PossibleValuesForRequestSchemeMatchConditionParametersOperator() []string {
	return []string{
		string(RequestSchemeMatchConditionParametersOperatorEqual),
	}
}

func (s *RequestSchemeMatchConditionParametersOperator) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseRequestSchemeMatchConditionParametersOperator(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseRequestSchemeMatchConditionParametersOperator(input string) (*RequestSchemeMatchConditionParametersOperator, error) {
	vals := map[string]RequestSchemeMatchConditionParametersOperator{
		"equal": RequestSchemeMatchConditionParametersOperatorEqual,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := RequestSchemeMatchConditionParametersOperator(input)
	return &out, nil
}

type RequestSchemeMatchValue string

const (
	RequestSchemeMatchValueHTTP  RequestSchemeMatchValue = "HTTP"
	RequestSchemeMatchValueHTTPS RequestSchemeMatchValue = "HTTPS"
)

func PossibleValuesForRequestSchemeMatchValue() []string {
	return []string{
		string(RequestSchemeMatchValueHTTP),
		string(RequestSchemeMatchValueHTTPS),
	}
}

func (s *RequestSchemeMatchValue) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseRequestSchemeMatchValue(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseRequestSchemeMatchValue(input string) (*RequestSchemeMatchValue, error) {
	vals := map[string]RequestSchemeMatchValue{
		"http":  RequestSchemeMatchValueHTTP,
		"https": RequestSchemeMatchValueHTTPS,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := RequestSchemeMatchValue(input)
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

func (s *RequestUriOperator) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseRequestUriOperator(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
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

type RuleCacheBehavior string

const (
	RuleCacheBehaviorHonorOrigin             RuleCacheBehavior = "HonorOrigin"
	RuleCacheBehaviorOverrideAlways          RuleCacheBehavior = "OverrideAlways"
	RuleCacheBehaviorOverrideIfOriginMissing RuleCacheBehavior = "OverrideIfOriginMissing"
)

func PossibleValuesForRuleCacheBehavior() []string {
	return []string{
		string(RuleCacheBehaviorHonorOrigin),
		string(RuleCacheBehaviorOverrideAlways),
		string(RuleCacheBehaviorOverrideIfOriginMissing),
	}
}

func (s *RuleCacheBehavior) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseRuleCacheBehavior(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseRuleCacheBehavior(input string) (*RuleCacheBehavior, error) {
	vals := map[string]RuleCacheBehavior{
		"honororigin":             RuleCacheBehaviorHonorOrigin,
		"overridealways":          RuleCacheBehaviorOverrideAlways,
		"overrideiforiginmissing": RuleCacheBehaviorOverrideIfOriginMissing,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := RuleCacheBehavior(input)
	return &out, nil
}

type RuleIsCompressionEnabled string

const (
	RuleIsCompressionEnabledDisabled RuleIsCompressionEnabled = "Disabled"
	RuleIsCompressionEnabledEnabled  RuleIsCompressionEnabled = "Enabled"
)

func PossibleValuesForRuleIsCompressionEnabled() []string {
	return []string{
		string(RuleIsCompressionEnabledDisabled),
		string(RuleIsCompressionEnabledEnabled),
	}
}

func (s *RuleIsCompressionEnabled) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseRuleIsCompressionEnabled(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseRuleIsCompressionEnabled(input string) (*RuleIsCompressionEnabled, error) {
	vals := map[string]RuleIsCompressionEnabled{
		"disabled": RuleIsCompressionEnabledDisabled,
		"enabled":  RuleIsCompressionEnabledEnabled,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := RuleIsCompressionEnabled(input)
	return &out, nil
}

type RuleQueryStringCachingBehavior string

const (
	RuleQueryStringCachingBehaviorIgnoreQueryString            RuleQueryStringCachingBehavior = "IgnoreQueryString"
	RuleQueryStringCachingBehaviorIgnoreSpecifiedQueryStrings  RuleQueryStringCachingBehavior = "IgnoreSpecifiedQueryStrings"
	RuleQueryStringCachingBehaviorIncludeSpecifiedQueryStrings RuleQueryStringCachingBehavior = "IncludeSpecifiedQueryStrings"
	RuleQueryStringCachingBehaviorUseQueryString               RuleQueryStringCachingBehavior = "UseQueryString"
)

func PossibleValuesForRuleQueryStringCachingBehavior() []string {
	return []string{
		string(RuleQueryStringCachingBehaviorIgnoreQueryString),
		string(RuleQueryStringCachingBehaviorIgnoreSpecifiedQueryStrings),
		string(RuleQueryStringCachingBehaviorIncludeSpecifiedQueryStrings),
		string(RuleQueryStringCachingBehaviorUseQueryString),
	}
}

func (s *RuleQueryStringCachingBehavior) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseRuleQueryStringCachingBehavior(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseRuleQueryStringCachingBehavior(input string) (*RuleQueryStringCachingBehavior, error) {
	vals := map[string]RuleQueryStringCachingBehavior{
		"ignorequerystring":            RuleQueryStringCachingBehaviorIgnoreQueryString,
		"ignorespecifiedquerystrings":  RuleQueryStringCachingBehaviorIgnoreSpecifiedQueryStrings,
		"includespecifiedquerystrings": RuleQueryStringCachingBehaviorIncludeSpecifiedQueryStrings,
		"usequerystring":               RuleQueryStringCachingBehaviorUseQueryString,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := RuleQueryStringCachingBehavior(input)
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

func (s *ServerPortOperator) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseServerPortOperator(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
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

func (s *SocketAddrOperator) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseSocketAddrOperator(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
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

func (s *SslProtocol) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseSslProtocol(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
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

func (s *SslProtocolOperator) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseSslProtocolOperator(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
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
	TransformURLDecode   Transform = "UrlDecode"
	TransformURLEncode   Transform = "UrlEncode"
	TransformUppercase   Transform = "Uppercase"
)

func PossibleValuesForTransform() []string {
	return []string{
		string(TransformLowercase),
		string(TransformRemoveNulls),
		string(TransformTrim),
		string(TransformURLDecode),
		string(TransformURLEncode),
		string(TransformUppercase),
	}
}

func (s *Transform) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseTransform(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseTransform(input string) (*Transform, error) {
	vals := map[string]Transform{
		"lowercase":   TransformLowercase,
		"removenulls": TransformRemoveNulls,
		"trim":        TransformTrim,
		"urldecode":   TransformURLDecode,
		"urlencode":   TransformURLEncode,
		"uppercase":   TransformUppercase,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := Transform(input)
	return &out, nil
}

type URLFileExtensionOperator string

const (
	URLFileExtensionOperatorAny                URLFileExtensionOperator = "Any"
	URLFileExtensionOperatorBeginsWith         URLFileExtensionOperator = "BeginsWith"
	URLFileExtensionOperatorContains           URLFileExtensionOperator = "Contains"
	URLFileExtensionOperatorEndsWith           URLFileExtensionOperator = "EndsWith"
	URLFileExtensionOperatorEqual              URLFileExtensionOperator = "Equal"
	URLFileExtensionOperatorGreaterThan        URLFileExtensionOperator = "GreaterThan"
	URLFileExtensionOperatorGreaterThanOrEqual URLFileExtensionOperator = "GreaterThanOrEqual"
	URLFileExtensionOperatorLessThan           URLFileExtensionOperator = "LessThan"
	URLFileExtensionOperatorLessThanOrEqual    URLFileExtensionOperator = "LessThanOrEqual"
	URLFileExtensionOperatorRegEx              URLFileExtensionOperator = "RegEx"
)

func PossibleValuesForURLFileExtensionOperator() []string {
	return []string{
		string(URLFileExtensionOperatorAny),
		string(URLFileExtensionOperatorBeginsWith),
		string(URLFileExtensionOperatorContains),
		string(URLFileExtensionOperatorEndsWith),
		string(URLFileExtensionOperatorEqual),
		string(URLFileExtensionOperatorGreaterThan),
		string(URLFileExtensionOperatorGreaterThanOrEqual),
		string(URLFileExtensionOperatorLessThan),
		string(URLFileExtensionOperatorLessThanOrEqual),
		string(URLFileExtensionOperatorRegEx),
	}
}

func (s *URLFileExtensionOperator) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseURLFileExtensionOperator(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseURLFileExtensionOperator(input string) (*URLFileExtensionOperator, error) {
	vals := map[string]URLFileExtensionOperator{
		"any":                URLFileExtensionOperatorAny,
		"beginswith":         URLFileExtensionOperatorBeginsWith,
		"contains":           URLFileExtensionOperatorContains,
		"endswith":           URLFileExtensionOperatorEndsWith,
		"equal":              URLFileExtensionOperatorEqual,
		"greaterthan":        URLFileExtensionOperatorGreaterThan,
		"greaterthanorequal": URLFileExtensionOperatorGreaterThanOrEqual,
		"lessthan":           URLFileExtensionOperatorLessThan,
		"lessthanorequal":    URLFileExtensionOperatorLessThanOrEqual,
		"regex":              URLFileExtensionOperatorRegEx,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := URLFileExtensionOperator(input)
	return &out, nil
}

type URLFileNameOperator string

const (
	URLFileNameOperatorAny                URLFileNameOperator = "Any"
	URLFileNameOperatorBeginsWith         URLFileNameOperator = "BeginsWith"
	URLFileNameOperatorContains           URLFileNameOperator = "Contains"
	URLFileNameOperatorEndsWith           URLFileNameOperator = "EndsWith"
	URLFileNameOperatorEqual              URLFileNameOperator = "Equal"
	URLFileNameOperatorGreaterThan        URLFileNameOperator = "GreaterThan"
	URLFileNameOperatorGreaterThanOrEqual URLFileNameOperator = "GreaterThanOrEqual"
	URLFileNameOperatorLessThan           URLFileNameOperator = "LessThan"
	URLFileNameOperatorLessThanOrEqual    URLFileNameOperator = "LessThanOrEqual"
	URLFileNameOperatorRegEx              URLFileNameOperator = "RegEx"
)

func PossibleValuesForURLFileNameOperator() []string {
	return []string{
		string(URLFileNameOperatorAny),
		string(URLFileNameOperatorBeginsWith),
		string(URLFileNameOperatorContains),
		string(URLFileNameOperatorEndsWith),
		string(URLFileNameOperatorEqual),
		string(URLFileNameOperatorGreaterThan),
		string(URLFileNameOperatorGreaterThanOrEqual),
		string(URLFileNameOperatorLessThan),
		string(URLFileNameOperatorLessThanOrEqual),
		string(URLFileNameOperatorRegEx),
	}
}

func (s *URLFileNameOperator) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseURLFileNameOperator(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseURLFileNameOperator(input string) (*URLFileNameOperator, error) {
	vals := map[string]URLFileNameOperator{
		"any":                URLFileNameOperatorAny,
		"beginswith":         URLFileNameOperatorBeginsWith,
		"contains":           URLFileNameOperatorContains,
		"endswith":           URLFileNameOperatorEndsWith,
		"equal":              URLFileNameOperatorEqual,
		"greaterthan":        URLFileNameOperatorGreaterThan,
		"greaterthanorequal": URLFileNameOperatorGreaterThanOrEqual,
		"lessthan":           URLFileNameOperatorLessThan,
		"lessthanorequal":    URLFileNameOperatorLessThanOrEqual,
		"regex":              URLFileNameOperatorRegEx,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := URLFileNameOperator(input)
	return &out, nil
}

type URLPathOperator string

const (
	URLPathOperatorAny                URLPathOperator = "Any"
	URLPathOperatorBeginsWith         URLPathOperator = "BeginsWith"
	URLPathOperatorContains           URLPathOperator = "Contains"
	URLPathOperatorEndsWith           URLPathOperator = "EndsWith"
	URLPathOperatorEqual              URLPathOperator = "Equal"
	URLPathOperatorGreaterThan        URLPathOperator = "GreaterThan"
	URLPathOperatorGreaterThanOrEqual URLPathOperator = "GreaterThanOrEqual"
	URLPathOperatorLessThan           URLPathOperator = "LessThan"
	URLPathOperatorLessThanOrEqual    URLPathOperator = "LessThanOrEqual"
	URLPathOperatorRegEx              URLPathOperator = "RegEx"
	URLPathOperatorWildcard           URLPathOperator = "Wildcard"
)

func PossibleValuesForURLPathOperator() []string {
	return []string{
		string(URLPathOperatorAny),
		string(URLPathOperatorBeginsWith),
		string(URLPathOperatorContains),
		string(URLPathOperatorEndsWith),
		string(URLPathOperatorEqual),
		string(URLPathOperatorGreaterThan),
		string(URLPathOperatorGreaterThanOrEqual),
		string(URLPathOperatorLessThan),
		string(URLPathOperatorLessThanOrEqual),
		string(URLPathOperatorRegEx),
		string(URLPathOperatorWildcard),
	}
}

func (s *URLPathOperator) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseURLPathOperator(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseURLPathOperator(input string) (*URLPathOperator, error) {
	vals := map[string]URLPathOperator{
		"any":                URLPathOperatorAny,
		"beginswith":         URLPathOperatorBeginsWith,
		"contains":           URLPathOperatorContains,
		"endswith":           URLPathOperatorEndsWith,
		"equal":              URLPathOperatorEqual,
		"greaterthan":        URLPathOperatorGreaterThan,
		"greaterthanorequal": URLPathOperatorGreaterThanOrEqual,
		"lessthan":           URLPathOperatorLessThan,
		"lessthanorequal":    URLPathOperatorLessThanOrEqual,
		"regex":              URLPathOperatorRegEx,
		"wildcard":           URLPathOperatorWildcard,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := URLPathOperator(input)
	return &out, nil
}
