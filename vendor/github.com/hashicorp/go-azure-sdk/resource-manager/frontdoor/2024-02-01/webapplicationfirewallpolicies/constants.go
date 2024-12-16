package webapplicationfirewallpolicies

import (
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ActionType string

const (
	ActionTypeAllow          ActionType = "Allow"
	ActionTypeAnomalyScoring ActionType = "AnomalyScoring"
	ActionTypeBlock          ActionType = "Block"
	ActionTypeJSChallenge    ActionType = "JSChallenge"
	ActionTypeLog            ActionType = "Log"
	ActionTypeRedirect       ActionType = "Redirect"
)

func PossibleValuesForActionType() []string {
	return []string{
		string(ActionTypeAllow),
		string(ActionTypeAnomalyScoring),
		string(ActionTypeBlock),
		string(ActionTypeJSChallenge),
		string(ActionTypeLog),
		string(ActionTypeRedirect),
	}
}

func parseActionType(input string) (*ActionType, error) {
	vals := map[string]ActionType{
		"allow":          ActionTypeAllow,
		"anomalyscoring": ActionTypeAnomalyScoring,
		"block":          ActionTypeBlock,
		"jschallenge":    ActionTypeJSChallenge,
		"log":            ActionTypeLog,
		"redirect":       ActionTypeRedirect,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := ActionType(input)
	return &out, nil
}

type CustomRuleEnabledState string

const (
	CustomRuleEnabledStateDisabled CustomRuleEnabledState = "Disabled"
	CustomRuleEnabledStateEnabled  CustomRuleEnabledState = "Enabled"
)

func PossibleValuesForCustomRuleEnabledState() []string {
	return []string{
		string(CustomRuleEnabledStateDisabled),
		string(CustomRuleEnabledStateEnabled),
	}
}

func parseCustomRuleEnabledState(input string) (*CustomRuleEnabledState, error) {
	vals := map[string]CustomRuleEnabledState{
		"disabled": CustomRuleEnabledStateDisabled,
		"enabled":  CustomRuleEnabledStateEnabled,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := CustomRuleEnabledState(input)
	return &out, nil
}

type ManagedRuleEnabledState string

const (
	ManagedRuleEnabledStateDisabled ManagedRuleEnabledState = "Disabled"
	ManagedRuleEnabledStateEnabled  ManagedRuleEnabledState = "Enabled"
)

func PossibleValuesForManagedRuleEnabledState() []string {
	return []string{
		string(ManagedRuleEnabledStateDisabled),
		string(ManagedRuleEnabledStateEnabled),
	}
}

func parseManagedRuleEnabledState(input string) (*ManagedRuleEnabledState, error) {
	vals := map[string]ManagedRuleEnabledState{
		"disabled": ManagedRuleEnabledStateDisabled,
		"enabled":  ManagedRuleEnabledStateEnabled,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := ManagedRuleEnabledState(input)
	return &out, nil
}

type ManagedRuleExclusionMatchVariable string

const (
	ManagedRuleExclusionMatchVariableQueryStringArgNames     ManagedRuleExclusionMatchVariable = "QueryStringArgNames"
	ManagedRuleExclusionMatchVariableRequestBodyJsonArgNames ManagedRuleExclusionMatchVariable = "RequestBodyJsonArgNames"
	ManagedRuleExclusionMatchVariableRequestBodyPostArgNames ManagedRuleExclusionMatchVariable = "RequestBodyPostArgNames"
	ManagedRuleExclusionMatchVariableRequestCookieNames      ManagedRuleExclusionMatchVariable = "RequestCookieNames"
	ManagedRuleExclusionMatchVariableRequestHeaderNames      ManagedRuleExclusionMatchVariable = "RequestHeaderNames"
)

func PossibleValuesForManagedRuleExclusionMatchVariable() []string {
	return []string{
		string(ManagedRuleExclusionMatchVariableQueryStringArgNames),
		string(ManagedRuleExclusionMatchVariableRequestBodyJsonArgNames),
		string(ManagedRuleExclusionMatchVariableRequestBodyPostArgNames),
		string(ManagedRuleExclusionMatchVariableRequestCookieNames),
		string(ManagedRuleExclusionMatchVariableRequestHeaderNames),
	}
}

func parseManagedRuleExclusionMatchVariable(input string) (*ManagedRuleExclusionMatchVariable, error) {
	vals := map[string]ManagedRuleExclusionMatchVariable{
		"querystringargnames":     ManagedRuleExclusionMatchVariableQueryStringArgNames,
		"requestbodyjsonargnames": ManagedRuleExclusionMatchVariableRequestBodyJsonArgNames,
		"requestbodypostargnames": ManagedRuleExclusionMatchVariableRequestBodyPostArgNames,
		"requestcookienames":      ManagedRuleExclusionMatchVariableRequestCookieNames,
		"requestheadernames":      ManagedRuleExclusionMatchVariableRequestHeaderNames,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := ManagedRuleExclusionMatchVariable(input)
	return &out, nil
}

type ManagedRuleExclusionSelectorMatchOperator string

const (
	ManagedRuleExclusionSelectorMatchOperatorContains   ManagedRuleExclusionSelectorMatchOperator = "Contains"
	ManagedRuleExclusionSelectorMatchOperatorEndsWith   ManagedRuleExclusionSelectorMatchOperator = "EndsWith"
	ManagedRuleExclusionSelectorMatchOperatorEquals     ManagedRuleExclusionSelectorMatchOperator = "Equals"
	ManagedRuleExclusionSelectorMatchOperatorEqualsAny  ManagedRuleExclusionSelectorMatchOperator = "EqualsAny"
	ManagedRuleExclusionSelectorMatchOperatorStartsWith ManagedRuleExclusionSelectorMatchOperator = "StartsWith"
)

func PossibleValuesForManagedRuleExclusionSelectorMatchOperator() []string {
	return []string{
		string(ManagedRuleExclusionSelectorMatchOperatorContains),
		string(ManagedRuleExclusionSelectorMatchOperatorEndsWith),
		string(ManagedRuleExclusionSelectorMatchOperatorEquals),
		string(ManagedRuleExclusionSelectorMatchOperatorEqualsAny),
		string(ManagedRuleExclusionSelectorMatchOperatorStartsWith),
	}
}

func parseManagedRuleExclusionSelectorMatchOperator(input string) (*ManagedRuleExclusionSelectorMatchOperator, error) {
	vals := map[string]ManagedRuleExclusionSelectorMatchOperator{
		"contains":   ManagedRuleExclusionSelectorMatchOperatorContains,
		"endswith":   ManagedRuleExclusionSelectorMatchOperatorEndsWith,
		"equals":     ManagedRuleExclusionSelectorMatchOperatorEquals,
		"equalsany":  ManagedRuleExclusionSelectorMatchOperatorEqualsAny,
		"startswith": ManagedRuleExclusionSelectorMatchOperatorStartsWith,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := ManagedRuleExclusionSelectorMatchOperator(input)
	return &out, nil
}

type ManagedRuleSetActionType string

const (
	ManagedRuleSetActionTypeBlock    ManagedRuleSetActionType = "Block"
	ManagedRuleSetActionTypeLog      ManagedRuleSetActionType = "Log"
	ManagedRuleSetActionTypeRedirect ManagedRuleSetActionType = "Redirect"
)

func PossibleValuesForManagedRuleSetActionType() []string {
	return []string{
		string(ManagedRuleSetActionTypeBlock),
		string(ManagedRuleSetActionTypeLog),
		string(ManagedRuleSetActionTypeRedirect),
	}
}

func parseManagedRuleSetActionType(input string) (*ManagedRuleSetActionType, error) {
	vals := map[string]ManagedRuleSetActionType{
		"block":    ManagedRuleSetActionTypeBlock,
		"log":      ManagedRuleSetActionTypeLog,
		"redirect": ManagedRuleSetActionTypeRedirect,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := ManagedRuleSetActionType(input)
	return &out, nil
}

type MatchVariable string

const (
	MatchVariableCookies       MatchVariable = "Cookies"
	MatchVariablePostArgs      MatchVariable = "PostArgs"
	MatchVariableQueryString   MatchVariable = "QueryString"
	MatchVariableRemoteAddr    MatchVariable = "RemoteAddr"
	MatchVariableRequestBody   MatchVariable = "RequestBody"
	MatchVariableRequestHeader MatchVariable = "RequestHeader"
	MatchVariableRequestMethod MatchVariable = "RequestMethod"
	MatchVariableRequestUri    MatchVariable = "RequestUri"
	MatchVariableSocketAddr    MatchVariable = "SocketAddr"
)

func PossibleValuesForMatchVariable() []string {
	return []string{
		string(MatchVariableCookies),
		string(MatchVariablePostArgs),
		string(MatchVariableQueryString),
		string(MatchVariableRemoteAddr),
		string(MatchVariableRequestBody),
		string(MatchVariableRequestHeader),
		string(MatchVariableRequestMethod),
		string(MatchVariableRequestUri),
		string(MatchVariableSocketAddr),
	}
}

func parseMatchVariable(input string) (*MatchVariable, error) {
	vals := map[string]MatchVariable{
		"cookies":       MatchVariableCookies,
		"postargs":      MatchVariablePostArgs,
		"querystring":   MatchVariableQueryString,
		"remoteaddr":    MatchVariableRemoteAddr,
		"requestbody":   MatchVariableRequestBody,
		"requestheader": MatchVariableRequestHeader,
		"requestmethod": MatchVariableRequestMethod,
		"requesturi":    MatchVariableRequestUri,
		"socketaddr":    MatchVariableSocketAddr,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := MatchVariable(input)
	return &out, nil
}

type Operator string

const (
	OperatorAny                Operator = "Any"
	OperatorBeginsWith         Operator = "BeginsWith"
	OperatorContains           Operator = "Contains"
	OperatorEndsWith           Operator = "EndsWith"
	OperatorEqual              Operator = "Equal"
	OperatorGeoMatch           Operator = "GeoMatch"
	OperatorGreaterThan        Operator = "GreaterThan"
	OperatorGreaterThanOrEqual Operator = "GreaterThanOrEqual"
	OperatorIPMatch            Operator = "IPMatch"
	OperatorLessThan           Operator = "LessThan"
	OperatorLessThanOrEqual    Operator = "LessThanOrEqual"
	OperatorRegEx              Operator = "RegEx"
)

func PossibleValuesForOperator() []string {
	return []string{
		string(OperatorAny),
		string(OperatorBeginsWith),
		string(OperatorContains),
		string(OperatorEndsWith),
		string(OperatorEqual),
		string(OperatorGeoMatch),
		string(OperatorGreaterThan),
		string(OperatorGreaterThanOrEqual),
		string(OperatorIPMatch),
		string(OperatorLessThan),
		string(OperatorLessThanOrEqual),
		string(OperatorRegEx),
	}
}

func parseOperator(input string) (*Operator, error) {
	vals := map[string]Operator{
		"any":                OperatorAny,
		"beginswith":         OperatorBeginsWith,
		"contains":           OperatorContains,
		"endswith":           OperatorEndsWith,
		"equal":              OperatorEqual,
		"geomatch":           OperatorGeoMatch,
		"greaterthan":        OperatorGreaterThan,
		"greaterthanorequal": OperatorGreaterThanOrEqual,
		"ipmatch":            OperatorIPMatch,
		"lessthan":           OperatorLessThan,
		"lessthanorequal":    OperatorLessThanOrEqual,
		"regex":              OperatorRegEx,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := Operator(input)
	return &out, nil
}

type PolicyEnabledState string

const (
	PolicyEnabledStateDisabled PolicyEnabledState = "Disabled"
	PolicyEnabledStateEnabled  PolicyEnabledState = "Enabled"
)

func PossibleValuesForPolicyEnabledState() []string {
	return []string{
		string(PolicyEnabledStateDisabled),
		string(PolicyEnabledStateEnabled),
	}
}

func parsePolicyEnabledState(input string) (*PolicyEnabledState, error) {
	vals := map[string]PolicyEnabledState{
		"disabled": PolicyEnabledStateDisabled,
		"enabled":  PolicyEnabledStateEnabled,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := PolicyEnabledState(input)
	return &out, nil
}

type PolicyMode string

const (
	PolicyModeDetection  PolicyMode = "Detection"
	PolicyModePrevention PolicyMode = "Prevention"
)

func PossibleValuesForPolicyMode() []string {
	return []string{
		string(PolicyModeDetection),
		string(PolicyModePrevention),
	}
}

func parsePolicyMode(input string) (*PolicyMode, error) {
	vals := map[string]PolicyMode{
		"detection":  PolicyModeDetection,
		"prevention": PolicyModePrevention,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := PolicyMode(input)
	return &out, nil
}

type PolicyRequestBodyCheck string

const (
	PolicyRequestBodyCheckDisabled PolicyRequestBodyCheck = "Disabled"
	PolicyRequestBodyCheckEnabled  PolicyRequestBodyCheck = "Enabled"
)

func PossibleValuesForPolicyRequestBodyCheck() []string {
	return []string{
		string(PolicyRequestBodyCheckDisabled),
		string(PolicyRequestBodyCheckEnabled),
	}
}

func parsePolicyRequestBodyCheck(input string) (*PolicyRequestBodyCheck, error) {
	vals := map[string]PolicyRequestBodyCheck{
		"disabled": PolicyRequestBodyCheckDisabled,
		"enabled":  PolicyRequestBodyCheckEnabled,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := PolicyRequestBodyCheck(input)
	return &out, nil
}

type PolicyResourceState string

const (
	PolicyResourceStateCreating  PolicyResourceState = "Creating"
	PolicyResourceStateDeleting  PolicyResourceState = "Deleting"
	PolicyResourceStateDisabled  PolicyResourceState = "Disabled"
	PolicyResourceStateDisabling PolicyResourceState = "Disabling"
	PolicyResourceStateEnabled   PolicyResourceState = "Enabled"
	PolicyResourceStateEnabling  PolicyResourceState = "Enabling"
)

func PossibleValuesForPolicyResourceState() []string {
	return []string{
		string(PolicyResourceStateCreating),
		string(PolicyResourceStateDeleting),
		string(PolicyResourceStateDisabled),
		string(PolicyResourceStateDisabling),
		string(PolicyResourceStateEnabled),
		string(PolicyResourceStateEnabling),
	}
}

func parsePolicyResourceState(input string) (*PolicyResourceState, error) {
	vals := map[string]PolicyResourceState{
		"creating":  PolicyResourceStateCreating,
		"deleting":  PolicyResourceStateDeleting,
		"disabled":  PolicyResourceStateDisabled,
		"disabling": PolicyResourceStateDisabling,
		"enabled":   PolicyResourceStateEnabled,
		"enabling":  PolicyResourceStateEnabling,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := PolicyResourceState(input)
	return &out, nil
}

type RuleType string

const (
	RuleTypeMatchRule     RuleType = "MatchRule"
	RuleTypeRateLimitRule RuleType = "RateLimitRule"
)

func PossibleValuesForRuleType() []string {
	return []string{
		string(RuleTypeMatchRule),
		string(RuleTypeRateLimitRule),
	}
}

func parseRuleType(input string) (*RuleType, error) {
	vals := map[string]RuleType{
		"matchrule":     RuleTypeMatchRule,
		"ratelimitrule": RuleTypeRateLimitRule,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := RuleType(input)
	return &out, nil
}

type ScrubbingRuleEntryMatchOperator string

const (
	ScrubbingRuleEntryMatchOperatorEquals    ScrubbingRuleEntryMatchOperator = "Equals"
	ScrubbingRuleEntryMatchOperatorEqualsAny ScrubbingRuleEntryMatchOperator = "EqualsAny"
)

func PossibleValuesForScrubbingRuleEntryMatchOperator() []string {
	return []string{
		string(ScrubbingRuleEntryMatchOperatorEquals),
		string(ScrubbingRuleEntryMatchOperatorEqualsAny),
	}
}

func parseScrubbingRuleEntryMatchOperator(input string) (*ScrubbingRuleEntryMatchOperator, error) {
	vals := map[string]ScrubbingRuleEntryMatchOperator{
		"equals":    ScrubbingRuleEntryMatchOperatorEquals,
		"equalsany": ScrubbingRuleEntryMatchOperatorEqualsAny,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := ScrubbingRuleEntryMatchOperator(input)
	return &out, nil
}

type ScrubbingRuleEntryMatchVariable string

const (
	ScrubbingRuleEntryMatchVariableQueryStringArgNames     ScrubbingRuleEntryMatchVariable = "QueryStringArgNames"
	ScrubbingRuleEntryMatchVariableRequestBodyJsonArgNames ScrubbingRuleEntryMatchVariable = "RequestBodyJsonArgNames"
	ScrubbingRuleEntryMatchVariableRequestBodyPostArgNames ScrubbingRuleEntryMatchVariable = "RequestBodyPostArgNames"
	ScrubbingRuleEntryMatchVariableRequestCookieNames      ScrubbingRuleEntryMatchVariable = "RequestCookieNames"
	ScrubbingRuleEntryMatchVariableRequestHeaderNames      ScrubbingRuleEntryMatchVariable = "RequestHeaderNames"
	ScrubbingRuleEntryMatchVariableRequestIPAddress        ScrubbingRuleEntryMatchVariable = "RequestIPAddress"
	ScrubbingRuleEntryMatchVariableRequestUri              ScrubbingRuleEntryMatchVariable = "RequestUri"
)

func PossibleValuesForScrubbingRuleEntryMatchVariable() []string {
	return []string{
		string(ScrubbingRuleEntryMatchVariableQueryStringArgNames),
		string(ScrubbingRuleEntryMatchVariableRequestBodyJsonArgNames),
		string(ScrubbingRuleEntryMatchVariableRequestBodyPostArgNames),
		string(ScrubbingRuleEntryMatchVariableRequestCookieNames),
		string(ScrubbingRuleEntryMatchVariableRequestHeaderNames),
		string(ScrubbingRuleEntryMatchVariableRequestIPAddress),
		string(ScrubbingRuleEntryMatchVariableRequestUri),
	}
}

func parseScrubbingRuleEntryMatchVariable(input string) (*ScrubbingRuleEntryMatchVariable, error) {
	vals := map[string]ScrubbingRuleEntryMatchVariable{
		"querystringargnames":     ScrubbingRuleEntryMatchVariableQueryStringArgNames,
		"requestbodyjsonargnames": ScrubbingRuleEntryMatchVariableRequestBodyJsonArgNames,
		"requestbodypostargnames": ScrubbingRuleEntryMatchVariableRequestBodyPostArgNames,
		"requestcookienames":      ScrubbingRuleEntryMatchVariableRequestCookieNames,
		"requestheadernames":      ScrubbingRuleEntryMatchVariableRequestHeaderNames,
		"requestipaddress":        ScrubbingRuleEntryMatchVariableRequestIPAddress,
		"requesturi":              ScrubbingRuleEntryMatchVariableRequestUri,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := ScrubbingRuleEntryMatchVariable(input)
	return &out, nil
}

type ScrubbingRuleEntryState string

const (
	ScrubbingRuleEntryStateDisabled ScrubbingRuleEntryState = "Disabled"
	ScrubbingRuleEntryStateEnabled  ScrubbingRuleEntryState = "Enabled"
)

func PossibleValuesForScrubbingRuleEntryState() []string {
	return []string{
		string(ScrubbingRuleEntryStateDisabled),
		string(ScrubbingRuleEntryStateEnabled),
	}
}

func parseScrubbingRuleEntryState(input string) (*ScrubbingRuleEntryState, error) {
	vals := map[string]ScrubbingRuleEntryState{
		"disabled": ScrubbingRuleEntryStateDisabled,
		"enabled":  ScrubbingRuleEntryStateEnabled,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := ScrubbingRuleEntryState(input)
	return &out, nil
}

type SkuName string

const (
	SkuNameClassicAzureFrontDoor  SkuName = "Classic_AzureFrontDoor"
	SkuNamePremiumAzureFrontDoor  SkuName = "Premium_AzureFrontDoor"
	SkuNameStandardAzureFrontDoor SkuName = "Standard_AzureFrontDoor"
)

func PossibleValuesForSkuName() []string {
	return []string{
		string(SkuNameClassicAzureFrontDoor),
		string(SkuNamePremiumAzureFrontDoor),
		string(SkuNameStandardAzureFrontDoor),
	}
}

func parseSkuName(input string) (*SkuName, error) {
	vals := map[string]SkuName{
		"classic_azurefrontdoor":  SkuNameClassicAzureFrontDoor,
		"premium_azurefrontdoor":  SkuNamePremiumAzureFrontDoor,
		"standard_azurefrontdoor": SkuNameStandardAzureFrontDoor,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := SkuName(input)
	return &out, nil
}

type TransformType string

const (
	TransformTypeLowercase   TransformType = "Lowercase"
	TransformTypeRemoveNulls TransformType = "RemoveNulls"
	TransformTypeTrim        TransformType = "Trim"
	TransformTypeURLDecode   TransformType = "UrlDecode"
	TransformTypeURLEncode   TransformType = "UrlEncode"
	TransformTypeUppercase   TransformType = "Uppercase"
)

func PossibleValuesForTransformType() []string {
	return []string{
		string(TransformTypeLowercase),
		string(TransformTypeRemoveNulls),
		string(TransformTypeTrim),
		string(TransformTypeURLDecode),
		string(TransformTypeURLEncode),
		string(TransformTypeUppercase),
	}
}

func parseTransformType(input string) (*TransformType, error) {
	vals := map[string]TransformType{
		"lowercase":   TransformTypeLowercase,
		"removenulls": TransformTypeRemoveNulls,
		"trim":        TransformTypeTrim,
		"urldecode":   TransformTypeURLDecode,
		"urlencode":   TransformTypeURLEncode,
		"uppercase":   TransformTypeUppercase,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := TransformType(input)
	return &out, nil
}

type VariableName string

const (
	VariableNameGeoLocation VariableName = "GeoLocation"
	VariableNameNone        VariableName = "None"
	VariableNameSocketAddr  VariableName = "SocketAddr"
)

func PossibleValuesForVariableName() []string {
	return []string{
		string(VariableNameGeoLocation),
		string(VariableNameNone),
		string(VariableNameSocketAddr),
	}
}

func parseVariableName(input string) (*VariableName, error) {
	vals := map[string]VariableName{
		"geolocation": VariableNameGeoLocation,
		"none":        VariableNameNone,
		"socketaddr":  VariableNameSocketAddr,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := VariableName(input)
	return &out, nil
}

type WebApplicationFirewallScrubbingState string

const (
	WebApplicationFirewallScrubbingStateDisabled WebApplicationFirewallScrubbingState = "Disabled"
	WebApplicationFirewallScrubbingStateEnabled  WebApplicationFirewallScrubbingState = "Enabled"
)

func PossibleValuesForWebApplicationFirewallScrubbingState() []string {
	return []string{
		string(WebApplicationFirewallScrubbingStateDisabled),
		string(WebApplicationFirewallScrubbingStateEnabled),
	}
}

func parseWebApplicationFirewallScrubbingState(input string) (*WebApplicationFirewallScrubbingState, error) {
	vals := map[string]WebApplicationFirewallScrubbingState{
		"disabled": WebApplicationFirewallScrubbingStateDisabled,
		"enabled":  WebApplicationFirewallScrubbingStateEnabled,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := WebApplicationFirewallScrubbingState(input)
	return &out, nil
}
