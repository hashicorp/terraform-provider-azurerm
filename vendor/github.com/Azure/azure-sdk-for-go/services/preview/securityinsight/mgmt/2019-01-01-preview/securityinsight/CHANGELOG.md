Generated from https://github.com/Azure/azure-rest-api-specs/tree/e5839301dfd424559851119c99ef0a2699fbd228/specification/securityinsights/resource-manager/readme.md tag: `package-2019-01-preview-only`

Code generator @microsoft.azure/autorest.go@2.1.180


## Breaking Changes

### Removed Constants

1. AlertRuleKind.AlertRuleKindAnomaly
1. KindBasicDataConnector.KindAmazonWebServicesCloudTrail
1. KindBasicDataConnector.KindAzureActiveDirectory
1. KindBasicDataConnector.KindAzureAdvancedThreatProtection
1. KindBasicDataConnector.KindAzureSecurityCenter
1. KindBasicDataConnector.KindDataConnector
1. KindBasicDataConnector.KindDynamics365
1. KindBasicDataConnector.KindMicrosoftCloudAppSecurity
1. KindBasicDataConnector.KindMicrosoftDefenderAdvancedThreatProtection
1. KindBasicDataConnector.KindOffice365
1. KindBasicDataConnector.KindOfficeATP
1. KindBasicDataConnector.KindThreatIntelligenceTaxii
1. KindBasicEntityQueryItem.KindEntityQueryItem
1. KindBasicEntityQueryItem.KindInsight

### Removed Funcs

1. *GetQueriesResponse.UnmarshalJSON([]byte) error
1. EntityQueryItem.AsBasicEntityQueryItem() (BasicEntityQueryItem, bool)
1. EntityQueryItem.AsEntityQueryItem() (*EntityQueryItem, bool)
1. EntityQueryItem.AsInsightQueryItem() (*InsightQueryItem, bool)
1. InsightQueryItem.AsBasicEntityQueryItem() (BasicEntityQueryItem, bool)
1. InsightQueryItem.AsEntityQueryItem() (*EntityQueryItem, bool)
1. InsightQueryItem.AsInsightQueryItem() (*InsightQueryItem, bool)
1. NewWatchlistItemClient(string) WatchlistItemClient
1. NewWatchlistItemClientWithBaseURI(string, string) WatchlistItemClient
1. OfficeConsentProperties.MarshalJSON() ([]byte, error)
1. PossibleKindBasicEntityQueryItemValues() []KindBasicEntityQueryItem
1. WatchlistItemClient.CreateOrUpdate(context.Context, string, string, string, string, string, WatchlistItem) (WatchlistItem, error)
1. WatchlistItemClient.CreateOrUpdatePreparer(context.Context, string, string, string, string, string, WatchlistItem) (*http.Request, error)
1. WatchlistItemClient.CreateOrUpdateResponder(*http.Response) (WatchlistItem, error)
1. WatchlistItemClient.CreateOrUpdateSender(*http.Request) (*http.Response, error)
1. WatchlistItemClient.Delete(context.Context, string, string, string, string, string) (autorest.Response, error)
1. WatchlistItemClient.DeletePreparer(context.Context, string, string, string, string, string) (*http.Request, error)
1. WatchlistItemClient.DeleteResponder(*http.Response) (autorest.Response, error)
1. WatchlistItemClient.DeleteSender(*http.Request) (*http.Response, error)

## Struct Changes

### Removed Structs

1. WatchlistItemClient

### Removed Struct Fields

1. MailMessageEntityProperties.ReceivedDate
1. OfficeConsentProperties.TenantName

## Signature Changes

### Const Types

1. KindThreatIntelligence changed type from KindBasicDataConnector to KindBasicAlertRule

### Struct Fields

1. EntityQueryItem.ID changed type from *uuid.UUID to *string
1. EntityQueryItem.Kind changed type from KindBasicEntityQueryItem to EntityQueryKind
1. GetQueriesResponse.Value changed type from *[]BasicEntityQueryItem to *[]EntityQueryItem
1. InsightQueryItem.ID changed type from *uuid.UUID to *string
1. InsightQueryItem.Kind changed type from KindBasicEntityQueryItem to EntityQueryKind
1. MailMessageEntityProperties.ThreatDetectionMethods changed type from *string to *[]string
1. MailMessageEntityProperties.Urls changed type from *string to *[]string
1. ThreatIntelligenceIndicatorProperties.ExternalReferences changed type from *[]string to *[]ThreatIntelligenceExternalReference
1. ThreatIntelligenceKillChainPhase.PhaseName changed type from *int32 to *string

### New Constants

1. ActionType.ActionTypeAutomationRuleAction
1. ActionType.ActionTypeModifyProperties
1. ActionType.ActionTypeRunPlaybook
1. AutomationRulePropertyConditionSupportedOperator.Contains
1. AutomationRulePropertyConditionSupportedOperator.EndsWith
1. AutomationRulePropertyConditionSupportedOperator.Equals
1. AutomationRulePropertyConditionSupportedOperator.NotContains
1. AutomationRulePropertyConditionSupportedOperator.NotEndsWith
1. AutomationRulePropertyConditionSupportedOperator.NotEquals
1. AutomationRulePropertyConditionSupportedOperator.NotStartsWith
1. AutomationRulePropertyConditionSupportedOperator.StartsWith
1. AutomationRulePropertyConditionSupportedProperty.AutomationRulePropertyConditionSupportedPropertyAccountAadTenantID
1. AutomationRulePropertyConditionSupportedProperty.AutomationRulePropertyConditionSupportedPropertyAccountAadUserID
1. AutomationRulePropertyConditionSupportedProperty.AutomationRulePropertyConditionSupportedPropertyAccountNTDomain
1. AutomationRulePropertyConditionSupportedProperty.AutomationRulePropertyConditionSupportedPropertyAccountName
1. AutomationRulePropertyConditionSupportedProperty.AutomationRulePropertyConditionSupportedPropertyAccountObjectGUID
1. AutomationRulePropertyConditionSupportedProperty.AutomationRulePropertyConditionSupportedPropertyAccountPUID
1. AutomationRulePropertyConditionSupportedProperty.AutomationRulePropertyConditionSupportedPropertyAccountSid
1. AutomationRulePropertyConditionSupportedProperty.AutomationRulePropertyConditionSupportedPropertyAccountUPNSuffix
1. AutomationRulePropertyConditionSupportedProperty.AutomationRulePropertyConditionSupportedPropertyAzureResourceResourceID
1. AutomationRulePropertyConditionSupportedProperty.AutomationRulePropertyConditionSupportedPropertyAzureResourceSubscriptionID
1. AutomationRulePropertyConditionSupportedProperty.AutomationRulePropertyConditionSupportedPropertyCloudApplicationAppID
1. AutomationRulePropertyConditionSupportedProperty.AutomationRulePropertyConditionSupportedPropertyCloudApplicationAppName
1. AutomationRulePropertyConditionSupportedProperty.AutomationRulePropertyConditionSupportedPropertyDNSDomainName
1. AutomationRulePropertyConditionSupportedProperty.AutomationRulePropertyConditionSupportedPropertyFileDirectory
1. AutomationRulePropertyConditionSupportedProperty.AutomationRulePropertyConditionSupportedPropertyFileHashValue
1. AutomationRulePropertyConditionSupportedProperty.AutomationRulePropertyConditionSupportedPropertyFileName
1. AutomationRulePropertyConditionSupportedProperty.AutomationRulePropertyConditionSupportedPropertyHostAzureID
1. AutomationRulePropertyConditionSupportedProperty.AutomationRulePropertyConditionSupportedPropertyHostNTDomain
1. AutomationRulePropertyConditionSupportedProperty.AutomationRulePropertyConditionSupportedPropertyHostName
1. AutomationRulePropertyConditionSupportedProperty.AutomationRulePropertyConditionSupportedPropertyHostNetBiosName
1. AutomationRulePropertyConditionSupportedProperty.AutomationRulePropertyConditionSupportedPropertyHostOSVersion
1. AutomationRulePropertyConditionSupportedProperty.AutomationRulePropertyConditionSupportedPropertyIPAddress
1. AutomationRulePropertyConditionSupportedProperty.AutomationRulePropertyConditionSupportedPropertyIncidentDescription
1. AutomationRulePropertyConditionSupportedProperty.AutomationRulePropertyConditionSupportedPropertyIncidentProviderName
1. AutomationRulePropertyConditionSupportedProperty.AutomationRulePropertyConditionSupportedPropertyIncidentRelatedAnalyticRuleIds
1. AutomationRulePropertyConditionSupportedProperty.AutomationRulePropertyConditionSupportedPropertyIncidentSeverity
1. AutomationRulePropertyConditionSupportedProperty.AutomationRulePropertyConditionSupportedPropertyIncidentStatus
1. AutomationRulePropertyConditionSupportedProperty.AutomationRulePropertyConditionSupportedPropertyIncidentTactics
1. AutomationRulePropertyConditionSupportedProperty.AutomationRulePropertyConditionSupportedPropertyIncidentTitle
1. AutomationRulePropertyConditionSupportedProperty.AutomationRulePropertyConditionSupportedPropertyIoTDeviceID
1. AutomationRulePropertyConditionSupportedProperty.AutomationRulePropertyConditionSupportedPropertyIoTDeviceModel
1. AutomationRulePropertyConditionSupportedProperty.AutomationRulePropertyConditionSupportedPropertyIoTDeviceName
1. AutomationRulePropertyConditionSupportedProperty.AutomationRulePropertyConditionSupportedPropertyIoTDeviceOperatingSystem
1. AutomationRulePropertyConditionSupportedProperty.AutomationRulePropertyConditionSupportedPropertyIoTDeviceType
1. AutomationRulePropertyConditionSupportedProperty.AutomationRulePropertyConditionSupportedPropertyIoTDeviceVendor
1. AutomationRulePropertyConditionSupportedProperty.AutomationRulePropertyConditionSupportedPropertyMailMessageDeliveryAction
1. AutomationRulePropertyConditionSupportedProperty.AutomationRulePropertyConditionSupportedPropertyMailMessageDeliveryLocation
1. AutomationRulePropertyConditionSupportedProperty.AutomationRulePropertyConditionSupportedPropertyMailMessageP1Sender
1. AutomationRulePropertyConditionSupportedProperty.AutomationRulePropertyConditionSupportedPropertyMailMessageP2Sender
1. AutomationRulePropertyConditionSupportedProperty.AutomationRulePropertyConditionSupportedPropertyMailMessageRecipient
1. AutomationRulePropertyConditionSupportedProperty.AutomationRulePropertyConditionSupportedPropertyMailMessageSenderIP
1. AutomationRulePropertyConditionSupportedProperty.AutomationRulePropertyConditionSupportedPropertyMailMessageSubject
1. AutomationRulePropertyConditionSupportedProperty.AutomationRulePropertyConditionSupportedPropertyMailboxDisplayName
1. AutomationRulePropertyConditionSupportedProperty.AutomationRulePropertyConditionSupportedPropertyMailboxPrimaryAddress
1. AutomationRulePropertyConditionSupportedProperty.AutomationRulePropertyConditionSupportedPropertyMailboxUPN
1. AutomationRulePropertyConditionSupportedProperty.AutomationRulePropertyConditionSupportedPropertyMalwareCategory
1. AutomationRulePropertyConditionSupportedProperty.AutomationRulePropertyConditionSupportedPropertyMalwareName
1. AutomationRulePropertyConditionSupportedProperty.AutomationRulePropertyConditionSupportedPropertyProcessCommandLine
1. AutomationRulePropertyConditionSupportedProperty.AutomationRulePropertyConditionSupportedPropertyProcessID
1. AutomationRulePropertyConditionSupportedProperty.AutomationRulePropertyConditionSupportedPropertyRegistryKey
1. AutomationRulePropertyConditionSupportedProperty.AutomationRulePropertyConditionSupportedPropertyRegistryValueData
1. AutomationRulePropertyConditionSupportedProperty.AutomationRulePropertyConditionSupportedPropertyURL
1. ConditionType.ConditionTypeAutomationRuleCondition
1. ConditionType.ConditionTypeProperty
1. DataConnectorKind.DataConnectorKindMicrosoftThreatIntelligence
1. DataConnectorKind.DataConnectorKindMicrosoftThreatProtection
1. KindBasicAlertRuleTemplate.KindBasicAlertRuleTemplateKindThreatIntelligence
1. KindBasicDataConnector.KindBasicDataConnectorKindAmazonWebServicesCloudTrail
1. KindBasicDataConnector.KindBasicDataConnectorKindAzureActiveDirectory
1. KindBasicDataConnector.KindBasicDataConnectorKindAzureAdvancedThreatProtection
1. KindBasicDataConnector.KindBasicDataConnectorKindAzureSecurityCenter
1. KindBasicDataConnector.KindBasicDataConnectorKindDataConnector
1. KindBasicDataConnector.KindBasicDataConnectorKindDynamics365
1. KindBasicDataConnector.KindBasicDataConnectorKindMicrosoftCloudAppSecurity
1. KindBasicDataConnector.KindBasicDataConnectorKindMicrosoftDefenderAdvancedThreatProtection
1. KindBasicDataConnector.KindBasicDataConnectorKindMicrosoftThreatIntelligence
1. KindBasicDataConnector.KindBasicDataConnectorKindMicrosoftThreatProtection
1. KindBasicDataConnector.KindBasicDataConnectorKindOffice365
1. KindBasicDataConnector.KindBasicDataConnectorKindOfficeATP
1. KindBasicDataConnector.KindBasicDataConnectorKindThreatIntelligence
1. KindBasicDataConnector.KindBasicDataConnectorKindThreatIntelligenceTaxii
1. KindBasicDataConnectorsCheckRequirements.KindBasicDataConnectorsCheckRequirementsKindMicrosoftThreatIntelligence
1. KindBasicDataConnectorsCheckRequirements.KindBasicDataConnectorsCheckRequirementsKindMicrosoftThreatProtection
1. KindBasicSettings.KindIPSyncer
1. PollingFrequency.OnceADay
1. PollingFrequency.OnceAMinute
1. PollingFrequency.OnceAnHour

### New Funcs

1. *AutomationRule.UnmarshalJSON([]byte) error
1. *AutomationRuleProperties.UnmarshalJSON([]byte) error
1. *AutomationRuleTriggeringLogic.UnmarshalJSON([]byte) error
1. *AutomationRulesListIterator.Next() error
1. *AutomationRulesListIterator.NextWithContext(context.Context) error
1. *AutomationRulesListPage.Next() error
1. *AutomationRulesListPage.NextWithContext(context.Context) error
1. *IPSyncer.UnmarshalJSON([]byte) error
1. *MSTICheckRequirements.UnmarshalJSON([]byte) error
1. *MSTIDataConnector.UnmarshalJSON([]byte) error
1. *MTPDataConnector.UnmarshalJSON([]byte) error
1. *MtpCheckRequirements.UnmarshalJSON([]byte) error
1. *ThreatIntelligenceAlertRule.UnmarshalJSON([]byte) error
1. *ThreatIntelligenceAlertRuleTemplate.UnmarshalJSON([]byte) error
1. *WatchlistItemListIterator.Next() error
1. *WatchlistItemListIterator.NextWithContext(context.Context) error
1. *WatchlistItemListPage.Next() error
1. *WatchlistItemListPage.NextWithContext(context.Context) error
1. AADCheckRequirements.AsMSTICheckRequirements() (*MSTICheckRequirements, bool)
1. AADCheckRequirements.AsMtpCheckRequirements() (*MtpCheckRequirements, bool)
1. AADDataConnector.AsMSTIDataConnector() (*MSTIDataConnector, bool)
1. AADDataConnector.AsMTPDataConnector() (*MTPDataConnector, bool)
1. AATPCheckRequirements.AsMSTICheckRequirements() (*MSTICheckRequirements, bool)
1. AATPCheckRequirements.AsMtpCheckRequirements() (*MtpCheckRequirements, bool)
1. AATPDataConnector.AsMSTIDataConnector() (*MSTIDataConnector, bool)
1. AATPDataConnector.AsMTPDataConnector() (*MTPDataConnector, bool)
1. ASCCheckRequirements.AsMSTICheckRequirements() (*MSTICheckRequirements, bool)
1. ASCCheckRequirements.AsMtpCheckRequirements() (*MtpCheckRequirements, bool)
1. ASCDataConnector.AsMSTIDataConnector() (*MSTIDataConnector, bool)
1. ASCDataConnector.AsMTPDataConnector() (*MTPDataConnector, bool)
1. AlertRule.AsThreatIntelligenceAlertRule() (*ThreatIntelligenceAlertRule, bool)
1. AlertRuleTemplate.AsThreatIntelligenceAlertRuleTemplate() (*ThreatIntelligenceAlertRuleTemplate, bool)
1. AutomationRule.MarshalJSON() ([]byte, error)
1. AutomationRuleAction.AsAutomationRuleAction() (*AutomationRuleAction, bool)
1. AutomationRuleAction.AsAutomationRuleModifyPropertiesAction() (*AutomationRuleModifyPropertiesAction, bool)
1. AutomationRuleAction.AsAutomationRuleRunPlaybookAction() (*AutomationRuleRunPlaybookAction, bool)
1. AutomationRuleAction.AsBasicAutomationRuleAction() (BasicAutomationRuleAction, bool)
1. AutomationRuleAction.MarshalJSON() ([]byte, error)
1. AutomationRuleCondition.AsAutomationRuleCondition() (*AutomationRuleCondition, bool)
1. AutomationRuleCondition.AsAutomationRulePropertyValuesCondition() (*AutomationRulePropertyValuesCondition, bool)
1. AutomationRuleCondition.AsBasicAutomationRuleCondition() (BasicAutomationRuleCondition, bool)
1. AutomationRuleCondition.MarshalJSON() ([]byte, error)
1. AutomationRuleModifyPropertiesAction.AsAutomationRuleAction() (*AutomationRuleAction, bool)
1. AutomationRuleModifyPropertiesAction.AsAutomationRuleModifyPropertiesAction() (*AutomationRuleModifyPropertiesAction, bool)
1. AutomationRuleModifyPropertiesAction.AsAutomationRuleRunPlaybookAction() (*AutomationRuleRunPlaybookAction, bool)
1. AutomationRuleModifyPropertiesAction.AsBasicAutomationRuleAction() (BasicAutomationRuleAction, bool)
1. AutomationRuleModifyPropertiesAction.MarshalJSON() ([]byte, error)
1. AutomationRuleProperties.MarshalJSON() ([]byte, error)
1. AutomationRulePropertyValuesCondition.AsAutomationRuleCondition() (*AutomationRuleCondition, bool)
1. AutomationRulePropertyValuesCondition.AsAutomationRulePropertyValuesCondition() (*AutomationRulePropertyValuesCondition, bool)
1. AutomationRulePropertyValuesCondition.AsBasicAutomationRuleCondition() (BasicAutomationRuleCondition, bool)
1. AutomationRulePropertyValuesCondition.MarshalJSON() ([]byte, error)
1. AutomationRuleRunPlaybookAction.AsAutomationRuleAction() (*AutomationRuleAction, bool)
1. AutomationRuleRunPlaybookAction.AsAutomationRuleModifyPropertiesAction() (*AutomationRuleModifyPropertiesAction, bool)
1. AutomationRuleRunPlaybookAction.AsAutomationRuleRunPlaybookAction() (*AutomationRuleRunPlaybookAction, bool)
1. AutomationRuleRunPlaybookAction.AsBasicAutomationRuleAction() (BasicAutomationRuleAction, bool)
1. AutomationRuleRunPlaybookAction.MarshalJSON() ([]byte, error)
1. AutomationRulesClient.CreateOrUpdate(context.Context, string, string, string, string, AutomationRule) (AutomationRule, error)
1. AutomationRulesClient.CreateOrUpdatePreparer(context.Context, string, string, string, string, AutomationRule) (*http.Request, error)
1. AutomationRulesClient.CreateOrUpdateResponder(*http.Response) (AutomationRule, error)
1. AutomationRulesClient.CreateOrUpdateSender(*http.Request) (*http.Response, error)
1. AutomationRulesClient.Delete(context.Context, string, string, string, string) (autorest.Response, error)
1. AutomationRulesClient.DeletePreparer(context.Context, string, string, string, string) (*http.Request, error)
1. AutomationRulesClient.DeleteResponder(*http.Response) (autorest.Response, error)
1. AutomationRulesClient.DeleteSender(*http.Request) (*http.Response, error)
1. AutomationRulesClient.Get(context.Context, string, string, string, string) (AutomationRule, error)
1. AutomationRulesClient.GetPreparer(context.Context, string, string, string, string) (*http.Request, error)
1. AutomationRulesClient.GetResponder(*http.Response) (AutomationRule, error)
1. AutomationRulesClient.GetSender(*http.Request) (*http.Response, error)
1. AutomationRulesClient.List(context.Context, string, string, string) (AutomationRulesListPage, error)
1. AutomationRulesClient.ListComplete(context.Context, string, string, string) (AutomationRulesListIterator, error)
1. AutomationRulesClient.ListPreparer(context.Context, string, string, string) (*http.Request, error)
1. AutomationRulesClient.ListResponder(*http.Response) (AutomationRulesList, error)
1. AutomationRulesClient.ListSender(*http.Request) (*http.Response, error)
1. AutomationRulesList.IsEmpty() bool
1. AutomationRulesList.MarshalJSON() ([]byte, error)
1. AutomationRulesListIterator.NotDone() bool
1. AutomationRulesListIterator.Response() AutomationRulesList
1. AutomationRulesListIterator.Value() AutomationRule
1. AutomationRulesListPage.NotDone() bool
1. AutomationRulesListPage.Response() AutomationRulesList
1. AutomationRulesListPage.Values() []AutomationRule
1. AwsCloudTrailCheckRequirements.AsMSTICheckRequirements() (*MSTICheckRequirements, bool)
1. AwsCloudTrailCheckRequirements.AsMtpCheckRequirements() (*MtpCheckRequirements, bool)
1. AwsCloudTrailDataConnector.AsMSTIDataConnector() (*MSTIDataConnector, bool)
1. AwsCloudTrailDataConnector.AsMTPDataConnector() (*MTPDataConnector, bool)
1. DataConnector.AsMSTIDataConnector() (*MSTIDataConnector, bool)
1. DataConnector.AsMTPDataConnector() (*MTPDataConnector, bool)
1. DataConnectorsCheckRequirements.AsMSTICheckRequirements() (*MSTICheckRequirements, bool)
1. DataConnectorsCheckRequirements.AsMtpCheckRequirements() (*MtpCheckRequirements, bool)
1. Dynamics365CheckRequirements.AsMSTICheckRequirements() (*MSTICheckRequirements, bool)
1. Dynamics365CheckRequirements.AsMtpCheckRequirements() (*MtpCheckRequirements, bool)
1. Dynamics365DataConnector.AsMSTIDataConnector() (*MSTIDataConnector, bool)
1. Dynamics365DataConnector.AsMTPDataConnector() (*MTPDataConnector, bool)
1. EntityAnalytics.AsIPSyncer() (*IPSyncer, bool)
1. EntityEdges.MarshalJSON() ([]byte, error)
1. EyesOn.AsIPSyncer() (*IPSyncer, bool)
1. FusionAlertRule.AsThreatIntelligenceAlertRule() (*ThreatIntelligenceAlertRule, bool)
1. FusionAlertRuleTemplate.AsThreatIntelligenceAlertRuleTemplate() (*ThreatIntelligenceAlertRuleTemplate, bool)
1. IPSyncer.AsBasicSettings() (BasicSettings, bool)
1. IPSyncer.AsEntityAnalytics() (*EntityAnalytics, bool)
1. IPSyncer.AsEyesOn() (*EyesOn, bool)
1. IPSyncer.AsIPSyncer() (*IPSyncer, bool)
1. IPSyncer.AsSettings() (*Settings, bool)
1. IPSyncer.AsUeba() (*Ueba, bool)
1. IPSyncer.MarshalJSON() ([]byte, error)
1. MCASCheckRequirements.AsMSTICheckRequirements() (*MSTICheckRequirements, bool)
1. MCASCheckRequirements.AsMtpCheckRequirements() (*MtpCheckRequirements, bool)
1. MCASDataConnector.AsMSTIDataConnector() (*MSTIDataConnector, bool)
1. MCASDataConnector.AsMTPDataConnector() (*MTPDataConnector, bool)
1. MDATPCheckRequirements.AsMSTICheckRequirements() (*MSTICheckRequirements, bool)
1. MDATPCheckRequirements.AsMtpCheckRequirements() (*MtpCheckRequirements, bool)
1. MDATPDataConnector.AsMSTIDataConnector() (*MSTIDataConnector, bool)
1. MDATPDataConnector.AsMTPDataConnector() (*MTPDataConnector, bool)
1. MLBehaviorAnalyticsAlertRule.AsThreatIntelligenceAlertRule() (*ThreatIntelligenceAlertRule, bool)
1. MLBehaviorAnalyticsAlertRuleTemplate.AsThreatIntelligenceAlertRuleTemplate() (*ThreatIntelligenceAlertRuleTemplate, bool)
1. MSTICheckRequirements.AsAADCheckRequirements() (*AADCheckRequirements, bool)
1. MSTICheckRequirements.AsAATPCheckRequirements() (*AATPCheckRequirements, bool)
1. MSTICheckRequirements.AsASCCheckRequirements() (*ASCCheckRequirements, bool)
1. MSTICheckRequirements.AsAwsCloudTrailCheckRequirements() (*AwsCloudTrailCheckRequirements, bool)
1. MSTICheckRequirements.AsBasicDataConnectorsCheckRequirements() (BasicDataConnectorsCheckRequirements, bool)
1. MSTICheckRequirements.AsDataConnectorsCheckRequirements() (*DataConnectorsCheckRequirements, bool)
1. MSTICheckRequirements.AsDynamics365CheckRequirements() (*Dynamics365CheckRequirements, bool)
1. MSTICheckRequirements.AsMCASCheckRequirements() (*MCASCheckRequirements, bool)
1. MSTICheckRequirements.AsMDATPCheckRequirements() (*MDATPCheckRequirements, bool)
1. MSTICheckRequirements.AsMSTICheckRequirements() (*MSTICheckRequirements, bool)
1. MSTICheckRequirements.AsMtpCheckRequirements() (*MtpCheckRequirements, bool)
1. MSTICheckRequirements.AsOfficeATPCheckRequirements() (*OfficeATPCheckRequirements, bool)
1. MSTICheckRequirements.AsTICheckRequirements() (*TICheckRequirements, bool)
1. MSTICheckRequirements.AsTiTaxiiCheckRequirements() (*TiTaxiiCheckRequirements, bool)
1. MSTICheckRequirements.MarshalJSON() ([]byte, error)
1. MSTIDataConnector.AsAADDataConnector() (*AADDataConnector, bool)
1. MSTIDataConnector.AsAATPDataConnector() (*AATPDataConnector, bool)
1. MSTIDataConnector.AsASCDataConnector() (*ASCDataConnector, bool)
1. MSTIDataConnector.AsAwsCloudTrailDataConnector() (*AwsCloudTrailDataConnector, bool)
1. MSTIDataConnector.AsBasicDataConnector() (BasicDataConnector, bool)
1. MSTIDataConnector.AsDataConnector() (*DataConnector, bool)
1. MSTIDataConnector.AsDynamics365DataConnector() (*Dynamics365DataConnector, bool)
1. MSTIDataConnector.AsMCASDataConnector() (*MCASDataConnector, bool)
1. MSTIDataConnector.AsMDATPDataConnector() (*MDATPDataConnector, bool)
1. MSTIDataConnector.AsMSTIDataConnector() (*MSTIDataConnector, bool)
1. MSTIDataConnector.AsMTPDataConnector() (*MTPDataConnector, bool)
1. MSTIDataConnector.AsOfficeATPDataConnector() (*OfficeATPDataConnector, bool)
1. MSTIDataConnector.AsOfficeDataConnector() (*OfficeDataConnector, bool)
1. MSTIDataConnector.AsTIDataConnector() (*TIDataConnector, bool)
1. MSTIDataConnector.AsTiTaxiiDataConnector() (*TiTaxiiDataConnector, bool)
1. MSTIDataConnector.MarshalJSON() ([]byte, error)
1. MTPDataConnector.AsAADDataConnector() (*AADDataConnector, bool)
1. MTPDataConnector.AsAATPDataConnector() (*AATPDataConnector, bool)
1. MTPDataConnector.AsASCDataConnector() (*ASCDataConnector, bool)
1. MTPDataConnector.AsAwsCloudTrailDataConnector() (*AwsCloudTrailDataConnector, bool)
1. MTPDataConnector.AsBasicDataConnector() (BasicDataConnector, bool)
1. MTPDataConnector.AsDataConnector() (*DataConnector, bool)
1. MTPDataConnector.AsDynamics365DataConnector() (*Dynamics365DataConnector, bool)
1. MTPDataConnector.AsMCASDataConnector() (*MCASDataConnector, bool)
1. MTPDataConnector.AsMDATPDataConnector() (*MDATPDataConnector, bool)
1. MTPDataConnector.AsMSTIDataConnector() (*MSTIDataConnector, bool)
1. MTPDataConnector.AsMTPDataConnector() (*MTPDataConnector, bool)
1. MTPDataConnector.AsOfficeATPDataConnector() (*OfficeATPDataConnector, bool)
1. MTPDataConnector.AsOfficeDataConnector() (*OfficeDataConnector, bool)
1. MTPDataConnector.AsTIDataConnector() (*TIDataConnector, bool)
1. MTPDataConnector.AsTiTaxiiDataConnector() (*TiTaxiiDataConnector, bool)
1. MTPDataConnector.MarshalJSON() ([]byte, error)
1. MicrosoftSecurityIncidentCreationAlertRule.AsThreatIntelligenceAlertRule() (*ThreatIntelligenceAlertRule, bool)
1. MicrosoftSecurityIncidentCreationAlertRuleTemplate.AsThreatIntelligenceAlertRuleTemplate() (*ThreatIntelligenceAlertRuleTemplate, bool)
1. MtpCheckRequirements.AsAADCheckRequirements() (*AADCheckRequirements, bool)
1. MtpCheckRequirements.AsAATPCheckRequirements() (*AATPCheckRequirements, bool)
1. MtpCheckRequirements.AsASCCheckRequirements() (*ASCCheckRequirements, bool)
1. MtpCheckRequirements.AsAwsCloudTrailCheckRequirements() (*AwsCloudTrailCheckRequirements, bool)
1. MtpCheckRequirements.AsBasicDataConnectorsCheckRequirements() (BasicDataConnectorsCheckRequirements, bool)
1. MtpCheckRequirements.AsDataConnectorsCheckRequirements() (*DataConnectorsCheckRequirements, bool)
1. MtpCheckRequirements.AsDynamics365CheckRequirements() (*Dynamics365CheckRequirements, bool)
1. MtpCheckRequirements.AsMCASCheckRequirements() (*MCASCheckRequirements, bool)
1. MtpCheckRequirements.AsMDATPCheckRequirements() (*MDATPCheckRequirements, bool)
1. MtpCheckRequirements.AsMSTICheckRequirements() (*MSTICheckRequirements, bool)
1. MtpCheckRequirements.AsMtpCheckRequirements() (*MtpCheckRequirements, bool)
1. MtpCheckRequirements.AsOfficeATPCheckRequirements() (*OfficeATPCheckRequirements, bool)
1. MtpCheckRequirements.AsTICheckRequirements() (*TICheckRequirements, bool)
1. MtpCheckRequirements.AsTiTaxiiCheckRequirements() (*TiTaxiiCheckRequirements, bool)
1. MtpCheckRequirements.MarshalJSON() ([]byte, error)
1. NewAutomationRulesClient(string) AutomationRulesClient
1. NewAutomationRulesClientWithBaseURI(string, string) AutomationRulesClient
1. NewAutomationRulesListIterator(AutomationRulesListPage) AutomationRulesListIterator
1. NewAutomationRulesListPage(AutomationRulesList, func(context.Context, AutomationRulesList) (AutomationRulesList, error)) AutomationRulesListPage
1. NewWatchlistItemListIterator(WatchlistItemListPage) WatchlistItemListIterator
1. NewWatchlistItemListPage(WatchlistItemList, func(context.Context, WatchlistItemList) (WatchlistItemList, error)) WatchlistItemListPage
1. NewWatchlistItemsClient(string) WatchlistItemsClient
1. NewWatchlistItemsClientWithBaseURI(string, string) WatchlistItemsClient
1. OfficeATPCheckRequirements.AsMSTICheckRequirements() (*MSTICheckRequirements, bool)
1. OfficeATPCheckRequirements.AsMtpCheckRequirements() (*MtpCheckRequirements, bool)
1. OfficeATPDataConnector.AsMSTIDataConnector() (*MSTIDataConnector, bool)
1. OfficeATPDataConnector.AsMTPDataConnector() (*MTPDataConnector, bool)
1. OfficeDataConnector.AsMSTIDataConnector() (*MSTIDataConnector, bool)
1. OfficeDataConnector.AsMTPDataConnector() (*MTPDataConnector, bool)
1. PossibleActionTypeValues() []ActionType
1. PossibleAutomationRulePropertyConditionSupportedOperatorValues() []AutomationRulePropertyConditionSupportedOperator
1. PossibleAutomationRulePropertyConditionSupportedPropertyValues() []AutomationRulePropertyConditionSupportedProperty
1. PossibleConditionTypeValues() []ConditionType
1. PossiblePollingFrequencyValues() []PollingFrequency
1. ScheduledAlertRule.AsThreatIntelligenceAlertRule() (*ThreatIntelligenceAlertRule, bool)
1. ScheduledAlertRuleTemplate.AsThreatIntelligenceAlertRuleTemplate() (*ThreatIntelligenceAlertRuleTemplate, bool)
1. Settings.AsIPSyncer() (*IPSyncer, bool)
1. TICheckRequirements.AsMSTICheckRequirements() (*MSTICheckRequirements, bool)
1. TICheckRequirements.AsMtpCheckRequirements() (*MtpCheckRequirements, bool)
1. TIDataConnector.AsMSTIDataConnector() (*MSTIDataConnector, bool)
1. TIDataConnector.AsMTPDataConnector() (*MTPDataConnector, bool)
1. ThreatIntelligenceAlertRule.AsAlertRule() (*AlertRule, bool)
1. ThreatIntelligenceAlertRule.AsBasicAlertRule() (BasicAlertRule, bool)
1. ThreatIntelligenceAlertRule.AsFusionAlertRule() (*FusionAlertRule, bool)
1. ThreatIntelligenceAlertRule.AsMLBehaviorAnalyticsAlertRule() (*MLBehaviorAnalyticsAlertRule, bool)
1. ThreatIntelligenceAlertRule.AsMicrosoftSecurityIncidentCreationAlertRule() (*MicrosoftSecurityIncidentCreationAlertRule, bool)
1. ThreatIntelligenceAlertRule.AsScheduledAlertRule() (*ScheduledAlertRule, bool)
1. ThreatIntelligenceAlertRule.AsThreatIntelligenceAlertRule() (*ThreatIntelligenceAlertRule, bool)
1. ThreatIntelligenceAlertRule.MarshalJSON() ([]byte, error)
1. ThreatIntelligenceAlertRuleProperties.MarshalJSON() ([]byte, error)
1. ThreatIntelligenceAlertRuleTemplate.AsAlertRuleTemplate() (*AlertRuleTemplate, bool)
1. ThreatIntelligenceAlertRuleTemplate.AsBasicAlertRuleTemplate() (BasicAlertRuleTemplate, bool)
1. ThreatIntelligenceAlertRuleTemplate.AsFusionAlertRuleTemplate() (*FusionAlertRuleTemplate, bool)
1. ThreatIntelligenceAlertRuleTemplate.AsMLBehaviorAnalyticsAlertRuleTemplate() (*MLBehaviorAnalyticsAlertRuleTemplate, bool)
1. ThreatIntelligenceAlertRuleTemplate.AsMicrosoftSecurityIncidentCreationAlertRuleTemplate() (*MicrosoftSecurityIncidentCreationAlertRuleTemplate, bool)
1. ThreatIntelligenceAlertRuleTemplate.AsScheduledAlertRuleTemplate() (*ScheduledAlertRuleTemplate, bool)
1. ThreatIntelligenceAlertRuleTemplate.AsThreatIntelligenceAlertRuleTemplate() (*ThreatIntelligenceAlertRuleTemplate, bool)
1. ThreatIntelligenceAlertRuleTemplate.MarshalJSON() ([]byte, error)
1. ThreatIntelligenceAlertRuleTemplateProperties.MarshalJSON() ([]byte, error)
1. ThreatIntelligenceExternalReference.MarshalJSON() ([]byte, error)
1. TiTaxiiCheckRequirements.AsMSTICheckRequirements() (*MSTICheckRequirements, bool)
1. TiTaxiiCheckRequirements.AsMtpCheckRequirements() (*MtpCheckRequirements, bool)
1. TiTaxiiDataConnector.AsMSTIDataConnector() (*MSTIDataConnector, bool)
1. TiTaxiiDataConnector.AsMTPDataConnector() (*MTPDataConnector, bool)
1. Ueba.AsIPSyncer() (*IPSyncer, bool)
1. WatchlistItemList.IsEmpty() bool
1. WatchlistItemList.MarshalJSON() ([]byte, error)
1. WatchlistItemListIterator.NotDone() bool
1. WatchlistItemListIterator.Response() WatchlistItemList
1. WatchlistItemListIterator.Value() WatchlistItem
1. WatchlistItemListPage.NotDone() bool
1. WatchlistItemListPage.Response() WatchlistItemList
1. WatchlistItemListPage.Values() []WatchlistItem
1. WatchlistItemsClient.CreateOrUpdate(context.Context, string, string, string, string, string, WatchlistItem) (WatchlistItem, error)
1. WatchlistItemsClient.CreateOrUpdatePreparer(context.Context, string, string, string, string, string, WatchlistItem) (*http.Request, error)
1. WatchlistItemsClient.CreateOrUpdateResponder(*http.Response) (WatchlistItem, error)
1. WatchlistItemsClient.CreateOrUpdateSender(*http.Request) (*http.Response, error)
1. WatchlistItemsClient.Delete(context.Context, string, string, string, string, string) (autorest.Response, error)
1. WatchlistItemsClient.DeletePreparer(context.Context, string, string, string, string, string) (*http.Request, error)
1. WatchlistItemsClient.DeleteResponder(*http.Response) (autorest.Response, error)
1. WatchlistItemsClient.DeleteSender(*http.Request) (*http.Response, error)
1. WatchlistItemsClient.Get(context.Context, string, string, string, string, string) (WatchlistItem, error)
1. WatchlistItemsClient.GetPreparer(context.Context, string, string, string, string, string) (*http.Request, error)
1. WatchlistItemsClient.GetResponder(*http.Response) (WatchlistItem, error)
1. WatchlistItemsClient.GetSender(*http.Request) (*http.Response, error)
1. WatchlistItemsClient.List(context.Context, string, string, string, string) (WatchlistItemListPage, error)
1. WatchlistItemsClient.ListComplete(context.Context, string, string, string, string) (WatchlistItemListIterator, error)
1. WatchlistItemsClient.ListPreparer(context.Context, string, string, string, string) (*http.Request, error)
1. WatchlistItemsClient.ListResponder(*http.Response) (WatchlistItemList, error)
1. WatchlistItemsClient.ListSender(*http.Request) (*http.Response, error)

## Struct Changes

### New Structs

1. AutomationRule
1. AutomationRuleAction
1. AutomationRuleCondition
1. AutomationRuleModifyPropertiesAction
1. AutomationRuleModifyPropertiesActionActionConfiguration
1. AutomationRuleProperties
1. AutomationRulePropertyValuesCondition
1. AutomationRulePropertyValuesConditionConditionProperties
1. AutomationRuleRunPlaybookAction
1. AutomationRuleRunPlaybookActionActionConfiguration
1. AutomationRuleTriggeringLogic
1. AutomationRulesClient
1. AutomationRulesList
1. AutomationRulesListIterator
1. AutomationRulesListPage
1. ConnectedEntity
1. EntityEdges
1. IPSyncer
1. IPSyncerSettingsProperties
1. MSTICheckRequirements
1. MSTICheckRequirementsProperties
1. MSTIDataConnector
1. MSTIDataConnectorDataTypes
1. MSTIDataConnectorDataTypesBingSafetyPhishingURL
1. MSTIDataConnectorDataTypesMicrosoftEmergingThreatFeed
1. MSTIDataConnectorProperties
1. MTPCheckRequirementsProperties
1. MTPDataConnector
1. MTPDataConnectorDataTypes
1. MTPDataConnectorDataTypesIncidents
1. MTPDataConnectorProperties
1. MtpCheckRequirements
1. ThreatIntelligenceAlertRule
1. ThreatIntelligenceAlertRuleProperties
1. ThreatIntelligenceAlertRuleTemplate
1. ThreatIntelligenceAlertRuleTemplateProperties
1. ThreatIntelligenceExternalReference
1. ThreatIntelligenceParsedPattern
1. ThreatIntelligenceParsedPatternTypeValue
1. WatchlistItemList
1. WatchlistItemListIterator
1. WatchlistItemListPage
1. WatchlistItemsClient

### New Struct Fields

1. BookmarkExpandResponseValue.Edges
1. BookmarkProperties.EventTime
1. BookmarkProperties.QueryEndTime
1. BookmarkProperties.QueryStartTime
1. EntityExpandResponseValue.Edges
1. MailMessageEntityProperties.ReceiveDate
1. OfficeConsentProperties.ConsentID
1. Operation.Origin
1. SecurityAlertTimelineItem.Description
1. SubmissionMailEntityProperties.SubmissionDate
1. SubmissionMailEntityProperties.SubmissionID
1. SubmissionMailEntityProperties.Submitter
1. TIDataConnectorProperties.TipLookbackPeriod
1. ThreatIntelligenceIndicatorProperties.Defanged
1. ThreatIntelligenceIndicatorProperties.Extensions
1. ThreatIntelligenceIndicatorProperties.ExternalLastUpdatedTimeUtc
1. ThreatIntelligenceIndicatorProperties.Language
1. ThreatIntelligenceIndicatorProperties.ObjectMarkingRefs
1. ThreatIntelligenceIndicatorProperties.ParsedPattern
1. ThreatIntelligenceIndicatorProperties.PatternVersion
1. TiTaxiiDataConnectorProperties.PollingFrequency
1. TiTaxiiDataConnectorProperties.TaxiiLookbackPeriod
