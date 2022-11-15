# Change History

## Breaking Changes

### Removed Constants

1. AFDEndpointProtocols.HTTP
1. AFDEndpointProtocols.HTTPS
1. ActionType.Allow
1. ActionType.Block
1. ActionType.Log
1. ActionType.Redirect
1. AfdMinimumTLSVersion.TLS10
1. AfdMinimumTLSVersion.TLS12
1. AfdProvisioningState.Creating
1. AfdProvisioningState.Deleting
1. AfdProvisioningState.Failed
1. AfdProvisioningState.Succeeded
1. AfdProvisioningState.Updating
1. AfdQueryStringCachingBehavior.IgnoreQueryString
1. AfdQueryStringCachingBehavior.NotSet
1. AfdQueryStringCachingBehavior.UseQueryString
1. Algorithm.SHA256
1. CacheBehavior.BypassCache
1. CacheBehavior.Override
1. CacheBehavior.SetIfMissing
1. CertificateType.Dedicated
1. CertificateType.Shared
1. CookiesOperator.Any
1. CookiesOperator.BeginsWith
1. CookiesOperator.Contains
1. CookiesOperator.EndsWith
1. CookiesOperator.Equal
1. CookiesOperator.GreaterThan
1. CookiesOperator.GreaterThanOrEqual
1. CookiesOperator.LessThan
1. CookiesOperator.LessThanOrEqual
1. CookiesOperator.RegEx
1. CustomHTTPSProvisioningSubstate.CertificateDeleted
1. CustomHTTPSProvisioningSubstate.CertificateDeployed
1. CustomHTTPSProvisioningSubstate.DeletingCertificate
1. CustomHTTPSProvisioningSubstate.DeployingCertificate
1. CustomHTTPSProvisioningSubstate.DomainControlValidationRequestApproved
1. CustomHTTPSProvisioningSubstate.DomainControlValidationRequestRejected
1. CustomHTTPSProvisioningSubstate.DomainControlValidationRequestTimedOut
1. CustomHTTPSProvisioningSubstate.IssuingCertificate
1. CustomHTTPSProvisioningSubstate.PendingDomainControlValidationREquestApproval
1. CustomHTTPSProvisioningSubstate.SubmittingDomainControlValidationRequest
1. CustomRuleEnabledState.Disabled
1. CustomRuleEnabledState.Enabled
1. DomainValidationState.Approved
1. DomainValidationState.Pending
1. DomainValidationState.PendingRevalidation
1. DomainValidationState.Submitting
1. DomainValidationState.TimedOut
1. DomainValidationState.Unknown
1. ForwardingProtocol.HTTPOnly
1. ForwardingProtocol.HTTPSOnly
1. ForwardingProtocol.MatchRequest
1. Granularity.P1D
1. Granularity.PT1H
1. Granularity.PT5M
1. HeaderAction.Append
1. HeaderAction.Delete
1. HeaderAction.Overwrite
1. IdentityType.Application
1. IdentityType.Key
1. IdentityType.ManagedIdentity
1. IdentityType.User
1. MatchProcessingBehavior.Continue
1. MatchProcessingBehavior.Stop
1. MatchVariable.Cookies
1. MatchVariable.PostArgs
1. MatchVariable.QueryString
1. MatchVariable.RemoteAddr
1. MatchVariable.RequestBody
1. MatchVariable.RequestHeader
1. MatchVariable.RequestMethod
1. MatchVariable.RequestURI
1. MatchVariable.SocketAddr
1. NameBasicDeliveryRuleAction.NameCacheExpiration
1. NameBasicDeliveryRuleAction.NameCacheKeyQueryString
1. NameBasicDeliveryRuleAction.NameDeliveryRuleAction
1. NameBasicDeliveryRuleAction.NameModifyRequestHeader
1. NameBasicDeliveryRuleAction.NameModifyResponseHeader
1. NameBasicDeliveryRuleAction.NameOriginGroupOverride
1. NameBasicDeliveryRuleAction.NameURLRedirect
1. NameBasicDeliveryRuleAction.NameURLRewrite
1. NameBasicDeliveryRuleAction.NameURLSigning
1. OptimizationType.DynamicSiteAcceleration
1. OptimizationType.GeneralMediaStreaming
1. OptimizationType.GeneralWebDelivery
1. OptimizationType.LargeFileDownload
1. OptimizationType.VideoOnDemandMediaStreaming
1. ParamIndicator.Expires
1. ParamIndicator.KeyID
1. ParamIndicator.Signature
1. PolicyMode.Detection
1. PolicyMode.Prevention
1. ProtocolType.IPBased
1. ProtocolType.ServerNameIndication
1. QueryStringBehavior.Exclude
1. QueryStringBehavior.ExcludeAll
1. QueryStringBehavior.Include
1. QueryStringBehavior.IncludeAll
1. RedirectType.Found
1. RedirectType.Moved
1. RedirectType.PermanentRedirect
1. RedirectType.TemporaryRedirect
1. ResourceType.MicrosoftCdnProfilesEndpoints
1. ResponseBasedDetectedErrorTypes.None
1. ResponseBasedDetectedErrorTypes.TCPAndHTTPErrors
1. ResponseBasedDetectedErrorTypes.TCPErrorsOnly
1. SkuName.CustomVerizon
1. SkuName.PremiumAzureFrontDoor
1. SkuName.PremiumChinaCdn
1. SkuName.PremiumVerizon
1. SkuName.Standard955BandWidthChinaCdn
1. SkuName.StandardAkamai
1. SkuName.StandardAvgBandWidthChinaCdn
1. SkuName.StandardAzureFrontDoor
1. SkuName.StandardChinaCdn
1. SkuName.StandardMicrosoft
1. SkuName.StandardPlus955BandWidthChinaCdn
1. SkuName.StandardPlusAvgBandWidthChinaCdn
1. SkuName.StandardPlusChinaCdn
1. SkuName.StandardVerizon
1. Status.AccessDenied
1. Status.CertificateExpired
1. Status.Invalid
1. Status.Valid
1. Transform.Lowercase
1. Transform.Uppercase
1. TypeBasicSecretParameters.TypeCustomerCertificate
1. TypeBasicSecretParameters.TypeManagedCertificate
1. TypeBasicSecretParameters.TypeSecretParameters
1. TypeBasicSecretParameters.TypeURLSigningKey
1. Unit.BitsPerSecond
1. Unit.Bytes
1. Unit.Count

### Signature Changes

#### Funcs

1. CustomDomainsClient.DisableCustomHTTPS
	- Returns
		- From: CustomDomain, error
		- To: CustomDomainsDisableCustomHTTPSFuture, error
1. CustomDomainsClient.DisableCustomHTTPSSender
	- Returns
		- From: *http.Response, error
		- To: CustomDomainsDisableCustomHTTPSFuture, error
1. CustomDomainsClient.EnableCustomHTTPS
	- Returns
		- From: CustomDomain, error
		- To: CustomDomainsEnableCustomHTTPSFuture, error
1. CustomDomainsClient.EnableCustomHTTPSSender
	- Returns
		- From: *http.Response, error
		- To: CustomDomainsEnableCustomHTTPSFuture, error
1. LogAnalyticsClient.GetLogAnalyticsMetrics
	- Params
		- From: context.Context, string, string, []string, date.Time, date.Time, string, []string, []string, []string, []string, []string
		- To: context.Context, string, string, []LogMetric, date.Time, date.Time, LogMetricsGranularity, []string, []string, []LogMetricsGroupBy, []string, []string
1. LogAnalyticsClient.GetLogAnalyticsMetricsPreparer
	- Params
		- From: context.Context, string, string, []string, date.Time, date.Time, string, []string, []string, []string, []string, []string
		- To: context.Context, string, string, []LogMetric, date.Time, date.Time, LogMetricsGranularity, []string, []string, []LogMetricsGroupBy, []string, []string
1. LogAnalyticsClient.GetLogAnalyticsRankings
	- Params
		- From: context.Context, string, string, []string, []string, int32, date.Time, date.Time, []string
		- To: context.Context, string, string, []LogRanking, []LogRankingMetric, int32, date.Time, date.Time, []string
1. LogAnalyticsClient.GetLogAnalyticsRankingsPreparer
	- Params
		- From: context.Context, string, string, []string, []string, int32, date.Time, date.Time, []string
		- To: context.Context, string, string, []LogRanking, []LogRankingMetric, int32, date.Time, date.Time, []string
1. LogAnalyticsClient.GetWafLogAnalyticsMetrics
	- Params
		- From: context.Context, string, string, []string, date.Time, date.Time, string, []string, []string, []string
		- To: context.Context, string, string, []WafMetric, date.Time, date.Time, WafGranularity, []WafAction, []WafRankingGroupBy, []WafRuleType
1. LogAnalyticsClient.GetWafLogAnalyticsMetricsPreparer
	- Params
		- From: context.Context, string, string, []string, date.Time, date.Time, string, []string, []string, []string
		- To: context.Context, string, string, []WafMetric, date.Time, date.Time, WafGranularity, []WafAction, []WafRankingGroupBy, []WafRuleType
1. LogAnalyticsClient.GetWafLogAnalyticsRankings
	- Params
		- From: context.Context, string, string, []string, date.Time, date.Time, int32, []string, []string, []string
		- To: context.Context, string, string, []WafMetric, date.Time, date.Time, int32, []WafRankingType, []WafAction, []WafRuleType
1. LogAnalyticsClient.GetWafLogAnalyticsRankingsPreparer
	- Params
		- From: context.Context, string, string, []string, date.Time, date.Time, int32, []string, []string, []string
		- To: context.Context, string, string, []WafMetric, date.Time, date.Time, int32, []WafRankingType, []WafAction, []WafRuleType
1. RuleSetsClient.Create
	- Params
		- From: context.Context, string, string, string, RuleSet
		- To: context.Context, string, string, string
1. RuleSetsClient.CreatePreparer
	- Params
		- From: context.Context, string, string, string, RuleSet
		- To: context.Context, string, string, string

## Additive Changes

### New Constants

1. AFDEndpointProtocols.AFDEndpointProtocolsHTTP
1. AFDEndpointProtocols.AFDEndpointProtocolsHTTPS
1. ActionType.ActionTypeAllow
1. ActionType.ActionTypeBlock
1. ActionType.ActionTypeLog
1. ActionType.ActionTypeRedirect
1. AfdMinimumTLSVersion.AfdMinimumTLSVersionTLS10
1. AfdMinimumTLSVersion.AfdMinimumTLSVersionTLS12
1. AfdProvisioningState.AfdProvisioningStateCreating
1. AfdProvisioningState.AfdProvisioningStateDeleting
1. AfdProvisioningState.AfdProvisioningStateFailed
1. AfdProvisioningState.AfdProvisioningStateSucceeded
1. AfdProvisioningState.AfdProvisioningStateUpdating
1. AfdQueryStringCachingBehavior.AfdQueryStringCachingBehaviorIgnoreQueryString
1. AfdQueryStringCachingBehavior.AfdQueryStringCachingBehaviorNotSet
1. AfdQueryStringCachingBehavior.AfdQueryStringCachingBehaviorUseQueryString
1. Algorithm.AlgorithmSHA256
1. CacheBehavior.CacheBehaviorBypassCache
1. CacheBehavior.CacheBehaviorOverride
1. CacheBehavior.CacheBehaviorSetIfMissing
1. CertificateType.CertificateTypeDedicated
1. CertificateType.CertificateTypeShared
1. CookiesOperator.CookiesOperatorAny
1. CookiesOperator.CookiesOperatorBeginsWith
1. CookiesOperator.CookiesOperatorContains
1. CookiesOperator.CookiesOperatorEndsWith
1. CookiesOperator.CookiesOperatorEqual
1. CookiesOperator.CookiesOperatorGreaterThan
1. CookiesOperator.CookiesOperatorGreaterThanOrEqual
1. CookiesOperator.CookiesOperatorLessThan
1. CookiesOperator.CookiesOperatorLessThanOrEqual
1. CookiesOperator.CookiesOperatorRegEx
1. CustomHTTPSProvisioningSubstate.CustomHTTPSProvisioningSubstateCertificateDeleted
1. CustomHTTPSProvisioningSubstate.CustomHTTPSProvisioningSubstateCertificateDeployed
1. CustomHTTPSProvisioningSubstate.CustomHTTPSProvisioningSubstateDeletingCertificate
1. CustomHTTPSProvisioningSubstate.CustomHTTPSProvisioningSubstateDeployingCertificate
1. CustomHTTPSProvisioningSubstate.CustomHTTPSProvisioningSubstateDomainControlValidationRequestApproved
1. CustomHTTPSProvisioningSubstate.CustomHTTPSProvisioningSubstateDomainControlValidationRequestRejected
1. CustomHTTPSProvisioningSubstate.CustomHTTPSProvisioningSubstateDomainControlValidationRequestTimedOut
1. CustomHTTPSProvisioningSubstate.CustomHTTPSProvisioningSubstateIssuingCertificate
1. CustomHTTPSProvisioningSubstate.CustomHTTPSProvisioningSubstatePendingDomainControlValidationREquestApproval
1. CustomHTTPSProvisioningSubstate.CustomHTTPSProvisioningSubstateSubmittingDomainControlValidationRequest
1. CustomRuleEnabledState.CustomRuleEnabledStateDisabled
1. CustomRuleEnabledState.CustomRuleEnabledStateEnabled
1. DomainValidationState.DomainValidationStateApproved
1. DomainValidationState.DomainValidationStatePending
1. DomainValidationState.DomainValidationStatePendingRevalidation
1. DomainValidationState.DomainValidationStateSubmitting
1. DomainValidationState.DomainValidationStateTimedOut
1. DomainValidationState.DomainValidationStateUnknown
1. ForwardingProtocol.ForwardingProtocolHTTPOnly
1. ForwardingProtocol.ForwardingProtocolHTTPSOnly
1. ForwardingProtocol.ForwardingProtocolMatchRequest
1. Granularity.GranularityP1D
1. Granularity.GranularityPT1H
1. Granularity.GranularityPT5M
1. HeaderAction.HeaderActionAppend
1. HeaderAction.HeaderActionDelete
1. HeaderAction.HeaderActionOverwrite
1. IdentityType.IdentityTypeApplication
1. IdentityType.IdentityTypeKey
1. IdentityType.IdentityTypeManagedIdentity
1. IdentityType.IdentityTypeUser
1. LogMetric.LogMetricClientRequestBandwidth
1. LogMetric.LogMetricClientRequestCount
1. LogMetric.LogMetricClientRequestTraffic
1. LogMetric.LogMetricOriginRequestBandwidth
1. LogMetric.LogMetricOriginRequestTraffic
1. LogMetric.LogMetricTotalLatency
1. LogMetricsGranularity.LogMetricsGranularityP1D
1. LogMetricsGranularity.LogMetricsGranularityPT1H
1. LogMetricsGranularity.LogMetricsGranularityPT5M
1. LogMetricsGroupBy.LogMetricsGroupByCacheStatus
1. LogMetricsGroupBy.LogMetricsGroupByCountry
1. LogMetricsGroupBy.LogMetricsGroupByCustomDomain
1. LogMetricsGroupBy.LogMetricsGroupByHTTPStatusCode
1. LogMetricsGroupBy.LogMetricsGroupByProtocol
1. LogRanking.LogRankingBrowser
1. LogRanking.LogRankingCountryOrRegion
1. LogRanking.LogRankingReferrer
1. LogRanking.LogRankingURL
1. LogRanking.LogRankingUserAgent
1. LogRankingMetric.LogRankingMetricClientRequestCount
1. LogRankingMetric.LogRankingMetricClientRequestTraffic
1. LogRankingMetric.LogRankingMetricErrorCount
1. LogRankingMetric.LogRankingMetricHitCount
1. LogRankingMetric.LogRankingMetricMissCount
1. LogRankingMetric.LogRankingMetricUserErrorCount
1. MatchProcessingBehavior.MatchProcessingBehaviorContinue
1. MatchProcessingBehavior.MatchProcessingBehaviorStop
1. MatchVariable.MatchVariableCookies
1. MatchVariable.MatchVariablePostArgs
1. MatchVariable.MatchVariableQueryString
1. MatchVariable.MatchVariableRemoteAddr
1. MatchVariable.MatchVariableRequestBody
1. MatchVariable.MatchVariableRequestHeader
1. MatchVariable.MatchVariableRequestMethod
1. MatchVariable.MatchVariableRequestURI
1. MatchVariable.MatchVariableSocketAddr
1. NameBasicDeliveryRuleAction.NameBasicDeliveryRuleActionNameCacheExpiration
1. NameBasicDeliveryRuleAction.NameBasicDeliveryRuleActionNameCacheKeyQueryString
1. NameBasicDeliveryRuleAction.NameBasicDeliveryRuleActionNameDeliveryRuleAction
1. NameBasicDeliveryRuleAction.NameBasicDeliveryRuleActionNameModifyRequestHeader
1. NameBasicDeliveryRuleAction.NameBasicDeliveryRuleActionNameModifyResponseHeader
1. NameBasicDeliveryRuleAction.NameBasicDeliveryRuleActionNameOriginGroupOverride
1. NameBasicDeliveryRuleAction.NameBasicDeliveryRuleActionNameURLRedirect
1. NameBasicDeliveryRuleAction.NameBasicDeliveryRuleActionNameURLRewrite
1. NameBasicDeliveryRuleAction.NameBasicDeliveryRuleActionNameURLSigning
1. OptimizationType.OptimizationTypeDynamicSiteAcceleration
1. OptimizationType.OptimizationTypeGeneralMediaStreaming
1. OptimizationType.OptimizationTypeGeneralWebDelivery
1. OptimizationType.OptimizationTypeLargeFileDownload
1. OptimizationType.OptimizationTypeVideoOnDemandMediaStreaming
1. ParamIndicator.ParamIndicatorExpires
1. ParamIndicator.ParamIndicatorKeyID
1. ParamIndicator.ParamIndicatorSignature
1. PolicyMode.PolicyModeDetection
1. PolicyMode.PolicyModePrevention
1. ProtocolType.ProtocolTypeIPBased
1. ProtocolType.ProtocolTypeServerNameIndication
1. QueryStringBehavior.QueryStringBehaviorExclude
1. QueryStringBehavior.QueryStringBehaviorExcludeAll
1. QueryStringBehavior.QueryStringBehaviorInclude
1. QueryStringBehavior.QueryStringBehaviorIncludeAll
1. RedirectType.RedirectTypeFound
1. RedirectType.RedirectTypeMoved
1. RedirectType.RedirectTypePermanentRedirect
1. RedirectType.RedirectTypeTemporaryRedirect
1. ResourceType.ResourceTypeMicrosoftCdnProfilesEndpoints
1. ResponseBasedDetectedErrorTypes.ResponseBasedDetectedErrorTypesNone
1. ResponseBasedDetectedErrorTypes.ResponseBasedDetectedErrorTypesTCPAndHTTPErrors
1. ResponseBasedDetectedErrorTypes.ResponseBasedDetectedErrorTypesTCPErrorsOnly
1. SkuName.SkuNameCustomVerizon
1. SkuName.SkuNamePremiumAzureFrontDoor
1. SkuName.SkuNamePremiumChinaCdn
1. SkuName.SkuNamePremiumVerizon
1. SkuName.SkuNameStandard955BandWidthChinaCdn
1. SkuName.SkuNameStandardAkamai
1. SkuName.SkuNameStandardAvgBandWidthChinaCdn
1. SkuName.SkuNameStandardAzureFrontDoor
1. SkuName.SkuNameStandardChinaCdn
1. SkuName.SkuNameStandardMicrosoft
1. SkuName.SkuNameStandardPlus955BandWidthChinaCdn
1. SkuName.SkuNameStandardPlusAvgBandWidthChinaCdn
1. SkuName.SkuNameStandardPlusChinaCdn
1. SkuName.SkuNameStandardVerizon
1. Status.StatusAccessDenied
1. Status.StatusCertificateExpired
1. Status.StatusInvalid
1. Status.StatusValid
1. Transform.TransformLowercase
1. Transform.TransformUppercase
1. TypeBasicSecretParameters.TypeBasicSecretParametersTypeCustomerCertificate
1. TypeBasicSecretParameters.TypeBasicSecretParametersTypeManagedCertificate
1. TypeBasicSecretParameters.TypeBasicSecretParametersTypeSecretParameters
1. TypeBasicSecretParameters.TypeBasicSecretParametersTypeURLSigningKey
1. Unit.UnitBitsPerSecond
1. Unit.UnitBytes
1. Unit.UnitCount
1. WafAction.WafActionAllow
1. WafAction.WafActionBlock
1. WafAction.WafActionLog
1. WafAction.WafActionRedirect
1. WafGranularity.WafGranularityP1D
1. WafGranularity.WafGranularityPT1H
1. WafGranularity.WafGranularityPT5M
1. WafMetric.WafMetricClientRequestCount
1. WafRankingGroupBy.WafRankingGroupByCustomDomain
1. WafRankingGroupBy.WafRankingGroupByHTTPStatusCode
1. WafRankingType.WafRankingTypeAction
1. WafRankingType.WafRankingTypeClientIP
1. WafRankingType.WafRankingTypeCountry
1. WafRankingType.WafRankingTypeRuleGroup
1. WafRankingType.WafRankingTypeRuleID
1. WafRankingType.WafRankingTypeRuleType
1. WafRankingType.WafRankingTypeURL
1. WafRankingType.WafRankingTypeUserAgent
1. WafRuleType.WafRuleTypeBot
1. WafRuleType.WafRuleTypeCustom
1. WafRuleType.WafRuleTypeManaged

### New Funcs

1. *CustomDomainProperties.UnmarshalJSON([]byte) error
1. *CustomDomainsDisableCustomHTTPSFuture.UnmarshalJSON([]byte) error
1. *CustomDomainsEnableCustomHTTPSFuture.UnmarshalJSON([]byte) error
1. PossibleLogMetricValues() []LogMetric
1. PossibleLogMetricsGranularityValues() []LogMetricsGranularity
1. PossibleLogMetricsGroupByValues() []LogMetricsGroupBy
1. PossibleLogRankingMetricValues() []LogRankingMetric
1. PossibleLogRankingValues() []LogRanking
1. PossibleWafActionValues() []WafAction
1. PossibleWafGranularityValues() []WafGranularity
1. PossibleWafMetricValues() []WafMetric
1. PossibleWafRankingGroupByValues() []WafRankingGroupBy
1. PossibleWafRankingTypeValues() []WafRankingType
1. PossibleWafRuleTypeValues() []WafRuleType

### Struct Changes

#### New Structs

1. CustomDomainsDisableCustomHTTPSFuture
1. CustomDomainsEnableCustomHTTPSFuture

#### New Struct Fields

1. CustomDomainProperties.CustomHTTPSParameters
1. ManagedRuleSetDefinition.SystemData
1. Resource.SystemData
