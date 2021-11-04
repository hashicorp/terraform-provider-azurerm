package webapplicationfirewallpolicies

type ActionType string

const (
	ActionTypeAllow    ActionType = "Allow"
	ActionTypeBlock    ActionType = "Block"
	ActionTypeLog      ActionType = "Log"
	ActionTypeRedirect ActionType = "Redirect"
)

type CustomRuleEnabledState string

const (
	CustomRuleEnabledStateDisabled CustomRuleEnabledState = "Disabled"
	CustomRuleEnabledStateEnabled  CustomRuleEnabledState = "Enabled"
)

type ManagedRuleEnabledState string

const (
	ManagedRuleEnabledStateDisabled ManagedRuleEnabledState = "Disabled"
	ManagedRuleEnabledStateEnabled  ManagedRuleEnabledState = "Enabled"
)

type ManagedRuleExclusionMatchVariable string

const (
	ManagedRuleExclusionMatchVariableQueryStringArgNames     ManagedRuleExclusionMatchVariable = "QueryStringArgNames"
	ManagedRuleExclusionMatchVariableRequestBodyPostArgNames ManagedRuleExclusionMatchVariable = "RequestBodyPostArgNames"
	ManagedRuleExclusionMatchVariableRequestCookieNames      ManagedRuleExclusionMatchVariable = "RequestCookieNames"
	ManagedRuleExclusionMatchVariableRequestHeaderNames      ManagedRuleExclusionMatchVariable = "RequestHeaderNames"
)

type ManagedRuleExclusionSelectorMatchOperator string

const (
	ManagedRuleExclusionSelectorMatchOperatorContains   ManagedRuleExclusionSelectorMatchOperator = "Contains"
	ManagedRuleExclusionSelectorMatchOperatorEndsWith   ManagedRuleExclusionSelectorMatchOperator = "EndsWith"
	ManagedRuleExclusionSelectorMatchOperatorEquals     ManagedRuleExclusionSelectorMatchOperator = "Equals"
	ManagedRuleExclusionSelectorMatchOperatorEqualsAny  ManagedRuleExclusionSelectorMatchOperator = "EqualsAny"
	ManagedRuleExclusionSelectorMatchOperatorStartsWith ManagedRuleExclusionSelectorMatchOperator = "StartsWith"
)

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

type PolicyEnabledState string

const (
	PolicyEnabledStateDisabled PolicyEnabledState = "Disabled"
	PolicyEnabledStateEnabled  PolicyEnabledState = "Enabled"
)

type PolicyMode string

const (
	PolicyModeDetection  PolicyMode = "Detection"
	PolicyModePrevention PolicyMode = "Prevention"
)

type PolicyResourceState string

const (
	PolicyResourceStateCreating  PolicyResourceState = "Creating"
	PolicyResourceStateDeleting  PolicyResourceState = "Deleting"
	PolicyResourceStateDisabled  PolicyResourceState = "Disabled"
	PolicyResourceStateDisabling PolicyResourceState = "Disabling"
	PolicyResourceStateEnabled   PolicyResourceState = "Enabled"
	PolicyResourceStateEnabling  PolicyResourceState = "Enabling"
)

type RuleType string

const (
	RuleTypeMatchRule     RuleType = "MatchRule"
	RuleTypeRateLimitRule RuleType = "RateLimitRule"
)

type TransformType string

const (
	TransformTypeLowercase   TransformType = "Lowercase"
	TransformTypeRemoveNulls TransformType = "RemoveNulls"
	TransformTypeTrim        TransformType = "Trim"
	TransformTypeUppercase   TransformType = "Uppercase"
	TransformTypeUrlDecode   TransformType = "UrlDecode"
	TransformTypeUrlEncode   TransformType = "UrlEncode"
)
