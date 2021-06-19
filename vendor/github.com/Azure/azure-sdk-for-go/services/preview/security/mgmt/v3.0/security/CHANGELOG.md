# Change History

## Breaking Changes

### Removed Constants

1. EnforcementMode.Audit
1. EnforcementMode.Enforce
1. EventSource.Assessments
1. EventSource.SecureScoreControls
1. EventSource.SecureScores
1. EventSource.SubAssessments

### Removed Funcs

1. *IotAlert.UnmarshalJSON([]byte) error
1. *IotAlertListIterator.Next() error
1. *IotAlertListIterator.NextWithContext(context.Context) error
1. *IotAlertListPage.Next() error
1. *IotAlertListPage.NextWithContext(context.Context) error
1. *IotRecommendation.UnmarshalJSON([]byte) error
1. *IotRecommendationListIterator.Next() error
1. *IotRecommendationListIterator.NextWithContext(context.Context) error
1. *IotRecommendationListPage.Next() error
1. *IotRecommendationListPage.NextWithContext(context.Context) error
1. AlertsClient.GetResourceGroupLevelAlerts(context.Context, string, string) (Alert, error)
1. AlertsClient.GetResourceGroupLevelAlertsPreparer(context.Context, string, string) (*http.Request, error)
1. AlertsClient.GetResourceGroupLevelAlertsResponder(*http.Response) (Alert, error)
1. AlertsClient.GetResourceGroupLevelAlertsSender(*http.Request) (*http.Response, error)
1. AlertsClient.GetSubscriptionLevelAlert(context.Context, string) (Alert, error)
1. AlertsClient.GetSubscriptionLevelAlertPreparer(context.Context, string) (*http.Request, error)
1. AlertsClient.GetSubscriptionLevelAlertResponder(*http.Response) (Alert, error)
1. AlertsClient.GetSubscriptionLevelAlertSender(*http.Request) (*http.Response, error)
1. AlertsClient.ListResourceGroupLevelAlertsByRegion(context.Context, string) (AlertListPage, error)
1. AlertsClient.ListResourceGroupLevelAlertsByRegionComplete(context.Context, string) (AlertListIterator, error)
1. AlertsClient.ListResourceGroupLevelAlertsByRegionPreparer(context.Context, string) (*http.Request, error)
1. AlertsClient.ListResourceGroupLevelAlertsByRegionResponder(*http.Response) (AlertList, error)
1. AlertsClient.ListResourceGroupLevelAlertsByRegionSender(*http.Request) (*http.Response, error)
1. AlertsClient.ListSubscriptionLevelAlertsByRegion(context.Context) (AlertListPage, error)
1. AlertsClient.ListSubscriptionLevelAlertsByRegionComplete(context.Context) (AlertListIterator, error)
1. AlertsClient.ListSubscriptionLevelAlertsByRegionPreparer(context.Context) (*http.Request, error)
1. AlertsClient.ListSubscriptionLevelAlertsByRegionResponder(*http.Response) (AlertList, error)
1. AlertsClient.ListSubscriptionLevelAlertsByRegionSender(*http.Request) (*http.Response, error)
1. AlertsClient.UpdateResourceGroupLevelAlertStateToDismiss(context.Context, string, string) (autorest.Response, error)
1. AlertsClient.UpdateResourceGroupLevelAlertStateToDismissPreparer(context.Context, string, string) (*http.Request, error)
1. AlertsClient.UpdateResourceGroupLevelAlertStateToDismissResponder(*http.Response) (autorest.Response, error)
1. AlertsClient.UpdateResourceGroupLevelAlertStateToDismissSender(*http.Request) (*http.Response, error)
1. AlertsClient.UpdateResourceGroupLevelAlertStateToReactivate(context.Context, string, string) (autorest.Response, error)
1. AlertsClient.UpdateResourceGroupLevelAlertStateToReactivatePreparer(context.Context, string, string) (*http.Request, error)
1. AlertsClient.UpdateResourceGroupLevelAlertStateToReactivateResponder(*http.Response) (autorest.Response, error)
1. AlertsClient.UpdateResourceGroupLevelAlertStateToReactivateSender(*http.Request) (*http.Response, error)
1. AlertsClient.UpdateSubscriptionLevelAlertStateToDismiss(context.Context, string) (autorest.Response, error)
1. AlertsClient.UpdateSubscriptionLevelAlertStateToDismissPreparer(context.Context, string) (*http.Request, error)
1. AlertsClient.UpdateSubscriptionLevelAlertStateToDismissResponder(*http.Response) (autorest.Response, error)
1. AlertsClient.UpdateSubscriptionLevelAlertStateToDismissSender(*http.Request) (*http.Response, error)
1. AlertsClient.UpdateSubscriptionLevelAlertStateToReactivate(context.Context, string) (autorest.Response, error)
1. AlertsClient.UpdateSubscriptionLevelAlertStateToReactivatePreparer(context.Context, string) (*http.Request, error)
1. AlertsClient.UpdateSubscriptionLevelAlertStateToReactivateResponder(*http.Response) (autorest.Response, error)
1. AlertsClient.UpdateSubscriptionLevelAlertStateToReactivateSender(*http.Request) (*http.Response, error)
1. IotAlert.MarshalJSON() ([]byte, error)
1. IotAlertList.IsEmpty() bool
1. IotAlertList.MarshalJSON() ([]byte, error)
1. IotAlertListIterator.NotDone() bool
1. IotAlertListIterator.Response() IotAlertList
1. IotAlertListIterator.Value() IotAlert
1. IotAlertListPage.NotDone() bool
1. IotAlertListPage.Response() IotAlertList
1. IotAlertListPage.Values() []IotAlert
1. IotAlertProperties.MarshalJSON() ([]byte, error)
1. IotAlertTypesClient.Get1(context.Context, string) (IotAlertType, error)
1. IotAlertTypesClient.Get1Preparer(context.Context, string) (*http.Request, error)
1. IotAlertTypesClient.Get1Responder(*http.Response) (IotAlertType, error)
1. IotAlertTypesClient.Get1Sender(*http.Request) (*http.Response, error)
1. IotAlertTypesClient.List1(context.Context) (IotAlertTypeList, error)
1. IotAlertTypesClient.List1Preparer(context.Context) (*http.Request, error)
1. IotAlertTypesClient.List1Responder(*http.Response) (IotAlertTypeList, error)
1. IotAlertTypesClient.List1Sender(*http.Request) (*http.Response, error)
1. IotAlertsClient.Get1(context.Context, string, string) (IotAlertModel, error)
1. IotAlertsClient.Get1Preparer(context.Context, string, string) (*http.Request, error)
1. IotAlertsClient.Get1Responder(*http.Response) (IotAlertModel, error)
1. IotAlertsClient.Get1Sender(*http.Request) (*http.Response, error)
1. IotAlertsClient.List1(context.Context, string, string, string, string, ManagementState, string, *int32, string) (IotAlertListModelPage, error)
1. IotAlertsClient.List1Complete(context.Context, string, string, string, string, ManagementState, string, *int32, string) (IotAlertListModelIterator, error)
1. IotAlertsClient.List1Preparer(context.Context, string, string, string, string, ManagementState, string, *int32, string) (*http.Request, error)
1. IotAlertsClient.List1Responder(*http.Response) (IotAlertListModel, error)
1. IotAlertsClient.List1Sender(*http.Request) (*http.Response, error)
1. IotRecommendation.MarshalJSON() ([]byte, error)
1. IotRecommendationList.IsEmpty() bool
1. IotRecommendationList.MarshalJSON() ([]byte, error)
1. IotRecommendationListIterator.NotDone() bool
1. IotRecommendationListIterator.Response() IotRecommendationList
1. IotRecommendationListIterator.Value() IotRecommendation
1. IotRecommendationListPage.NotDone() bool
1. IotRecommendationListPage.Response() IotRecommendationList
1. IotRecommendationListPage.Values() []IotRecommendation
1. IotRecommendationProperties.MarshalJSON() ([]byte, error)
1. IotRecommendationTypesClient.Get1(context.Context, string) (IotRecommendationType, error)
1. IotRecommendationTypesClient.Get1Preparer(context.Context, string) (*http.Request, error)
1. IotRecommendationTypesClient.Get1Responder(*http.Response) (IotRecommendationType, error)
1. IotRecommendationTypesClient.Get1Sender(*http.Request) (*http.Response, error)
1. IotRecommendationTypesClient.List1(context.Context) (IotRecommendationTypeList, error)
1. IotRecommendationTypesClient.List1Preparer(context.Context) (*http.Request, error)
1. IotRecommendationTypesClient.List1Responder(*http.Response) (IotRecommendationTypeList, error)
1. IotRecommendationTypesClient.List1Sender(*http.Request) (*http.Response, error)
1. IotRecommendationsClient.Get1(context.Context, string, string) (IotRecommendationModel, error)
1. IotRecommendationsClient.Get1Preparer(context.Context, string, string) (*http.Request, error)
1. IotRecommendationsClient.Get1Responder(*http.Response) (IotRecommendationModel, error)
1. IotRecommendationsClient.Get1Sender(*http.Request) (*http.Response, error)
1. IotRecommendationsClient.List1(context.Context, string, string, string, *int32, string) (IotRecommendationListModelPage, error)
1. IotRecommendationsClient.List1Complete(context.Context, string, string, string, *int32, string) (IotRecommendationListModelIterator, error)
1. IotRecommendationsClient.List1Preparer(context.Context, string, string, string, *int32, string) (*http.Request, error)
1. IotRecommendationsClient.List1Responder(*http.Response) (IotRecommendationListModel, error)
1. IotRecommendationsClient.List1Sender(*http.Request) (*http.Response, error)
1. NewIotAlertListIterator(IotAlertListPage) IotAlertListIterator
1. NewIotAlertListPage(IotAlertList, func(context.Context, IotAlertList) (IotAlertList, error)) IotAlertListPage
1. NewIotRecommendationListIterator(IotRecommendationListPage) IotRecommendationListIterator
1. NewIotRecommendationListPage(IotRecommendationList, func(context.Context, IotRecommendationList) (IotRecommendationList, error)) IotRecommendationListPage
1. PossibleCategoryValues() []Category

### Struct Changes

#### Removed Structs

1. IotAlert
1. IotAlertList
1. IotAlertListIterator
1. IotAlertListPage
1. IotAlertProperties
1. IotRecommendation
1. IotRecommendationList
1. IotRecommendationListIterator
1. IotRecommendationListPage
1. IotRecommendationProperties

#### Removed Struct Fields

1. AssessmentMetadataProperties.Category

### Signature Changes

#### Const Types

1. Alerts changed type from EventSource to AdditionalWorkspaceDataType
1. Compute changed type from Category to Categories
1. Data changed type from Category to Categories
1. IdentityAndAccess changed type from Category to Categories
1. IoT changed type from Category to Categories
1. Networking changed type from Category to Categories
1. None changed type from EnforcementMode to EndOfSupportStatus
1. RawEvents changed type from ExportData to AdditionalWorkspaceDataType

#### Funcs

1. IotAlertTypesClient.Get
	- Params
		- From: context.Context, string, string, string
		- To: context.Context, string
1. IotAlertTypesClient.GetPreparer
	- Params
		- From: context.Context, string, string, string
		- To: context.Context, string
1. IotAlertTypesClient.List
	- Params
		- From: context.Context, string, string
		- To: context.Context
1. IotAlertTypesClient.ListPreparer
	- Params
		- From: context.Context, string, string
		- To: context.Context
1. IotAlertsClient.Get
	- Params
		- From: context.Context, string, string, string
		- To: context.Context, string, string
	- Returns
		- From: IotAlert, error
		- To: IotAlertModel, error
1. IotAlertsClient.GetPreparer
	- Params
		- From: context.Context, string, string, string
		- To: context.Context, string, string
1. IotAlertsClient.GetResponder
	- Returns
		- From: IotAlert, error
		- To: IotAlertModel, error
1. IotAlertsClient.List
	- Params
		- From: context.Context, string, string, string, string, string, string, *int32, string
		- To: context.Context, string, string, string, string, ManagementState, string, *int32, string
	- Returns
		- From: IotAlertListPage, error
		- To: IotAlertListModelPage, error
1. IotAlertsClient.ListComplete
	- Params
		- From: context.Context, string, string, string, string, string, string, *int32, string
		- To: context.Context, string, string, string, string, ManagementState, string, *int32, string
	- Returns
		- From: IotAlertListIterator, error
		- To: IotAlertListModelIterator, error
1. IotAlertsClient.ListPreparer
	- Params
		- From: context.Context, string, string, string, string, string, string, *int32, string
		- To: context.Context, string, string, string, string, ManagementState, string, *int32, string
1. IotAlertsClient.ListResponder
	- Returns
		- From: IotAlertList, error
		- To: IotAlertListModel, error
1. IotRecommendationTypesClient.Get
	- Params
		- From: context.Context, string, string, string
		- To: context.Context, string
1. IotRecommendationTypesClient.GetPreparer
	- Params
		- From: context.Context, string, string, string
		- To: context.Context, string
1. IotRecommendationTypesClient.List
	- Params
		- From: context.Context, string, string
		- To: context.Context
1. IotRecommendationTypesClient.ListPreparer
	- Params
		- From: context.Context, string, string
		- To: context.Context
1. IotRecommendationsClient.Get
	- Params
		- From: context.Context, string, string, string
		- To: context.Context, string, string
	- Returns
		- From: IotRecommendation, error
		- To: IotRecommendationModel, error
1. IotRecommendationsClient.GetPreparer
	- Params
		- From: context.Context, string, string, string
		- To: context.Context, string, string
1. IotRecommendationsClient.GetResponder
	- Returns
		- From: IotRecommendation, error
		- To: IotRecommendationModel, error
1. IotRecommendationsClient.List
	- Params
		- From: context.Context, string, string, string, string, *int32, string
		- To: context.Context, string, string, string, *int32, string
	- Returns
		- From: IotRecommendationListPage, error
		- To: IotRecommendationListModelPage, error
1. IotRecommendationsClient.ListComplete
	- Params
		- From: context.Context, string, string, string, string, *int32, string
		- To: context.Context, string, string, string, *int32, string
	- Returns
		- From: IotRecommendationListIterator, error
		- To: IotRecommendationListModelIterator, error
1. IotRecommendationsClient.ListPreparer
	- Params
		- From: context.Context, string, string, string, string, *int32, string
		- To: context.Context, string, string, string, *int32, string
1. IotRecommendationsClient.ListResponder
	- Returns
		- From: IotRecommendationList, error
		- To: IotRecommendationListModel, error
1. ServerVulnerabilityAssessmentClient.Delete
	- Returns
		- From: autorest.Response, error
		- To: ServerVulnerabilityAssessmentDeleteFuture, error
1. ServerVulnerabilityAssessmentClient.DeleteSender
	- Returns
		- From: *http.Response, error
		- To: ServerVulnerabilityAssessmentDeleteFuture, error

#### Struct Fields

1. IoTSecurityAggregatedAlertProperties.Count changed type from *int32 to *int64
1. IoTSecurityAggregatedAlertPropertiesTopDevicesListItem.AlertsCount changed type from *int32 to *int64
1. IoTSecurityAggregatedRecommendationProperties.HealthyDevices changed type from *int32 to *int64
1. IoTSecurityAggregatedRecommendationProperties.UnhealthyDeviceCount changed type from *int32 to *int64
1. IoTSecurityAlertedDevice.AlertsCount changed type from *int32 to *int64
1. IoTSecurityDeviceAlert.AlertsCount changed type from *int32 to *int64
1. IoTSecurityDeviceRecommendation.DevicesCount changed type from *int32 to *int64
1. IoTSecuritySolutionAnalyticsModelProperties.UnhealthyDeviceCount changed type from *int32 to *int64
1. IoTSeverityMetrics.High changed type from *int32 to *int64
1. IoTSeverityMetrics.Low changed type from *int32 to *int64
1. IoTSeverityMetrics.Medium changed type from *int32 to *int64

## Additive Changes

### New Constants

1. AdditionalWorkspaceType.Sentinel
1. BundleType.AppServices
1. BundleType.DNS
1. BundleType.KeyVaults
1. BundleType.KubernetesService
1. BundleType.ResourceManager
1. BundleType.SQLServers
1. BundleType.StorageAccounts
1. BundleType.VirtualMachines
1. CreatedByType.Application
1. CreatedByType.Key
1. CreatedByType.ManagedIdentity
1. CreatedByType.User
1. EndOfSupportStatus.NoLongerSupported
1. EndOfSupportStatus.UpcomingNoLongerSupported
1. EndOfSupportStatus.UpcomingVersionNoLongerSupported
1. EndOfSupportStatus.VersionNoLongerSupported
1. EnforcementMode.EnforcementModeAudit
1. EnforcementMode.EnforcementModeEnforce
1. EnforcementMode.EnforcementModeNone
1. EventSource.EventSourceAlerts
1. EventSource.EventSourceAssessments
1. EventSource.EventSourceRegulatoryComplianceAssessment
1. EventSource.EventSourceRegulatoryComplianceAssessmentSnapshot
1. EventSource.EventSourceSecureScoreControls
1. EventSource.EventSourceSecureScoreControlsSnapshot
1. EventSource.EventSourceSecureScores
1. EventSource.EventSourceSecureScoresSnapshot
1. EventSource.EventSourceSubAssessments
1. ExportData.ExportDataRawEvents
1. KindEnum2.KindAlertSimulatorRequestProperties
1. KindEnum2.KindBundles
1. OnboardingKind.Default
1. OnboardingKind.MigratedToAzure
1. RuleTypeBasicCustomAlertRule.RuleTypeConnectionFromIPNotAllowed
1. SensorType.SensorTypeEnterprise
1. SensorType.SensorTypeOt

### New Funcs

1. *AlertSimulatorBundlesRequestProperties.UnmarshalJSON([]byte) error
1. *AlertSimulatorRequestBody.UnmarshalJSON([]byte) error
1. *AlertSimulatorRequestProperties.UnmarshalJSON([]byte) error
1. *AlertsSimulateFuture.UnmarshalJSON([]byte) error
1. *IngestionSettingListIterator.Next() error
1. *IngestionSettingListIterator.NextWithContext(context.Context) error
1. *IngestionSettingListPage.Next() error
1. *IngestionSettingListPage.NextWithContext(context.Context) error
1. *ServerVulnerabilityAssessmentDeleteFuture.UnmarshalJSON([]byte) error
1. *Software.UnmarshalJSON([]byte) error
1. *SoftwaresListIterator.Next() error
1. *SoftwaresListIterator.NextWithContext(context.Context) error
1. *SoftwaresListPage.Next() error
1. *SoftwaresListPage.NextWithContext(context.Context) error
1. ActiveConnectionsNotInAllowedRange.AsConnectionFromIPNotAllowed() (*ConnectionFromIPNotAllowed, bool)
1. AlertSimulatorBundlesRequestProperties.AsAlertSimulatorBundlesRequestProperties() (*AlertSimulatorBundlesRequestProperties, bool)
1. AlertSimulatorBundlesRequestProperties.AsAlertSimulatorRequestProperties() (*AlertSimulatorRequestProperties, bool)
1. AlertSimulatorBundlesRequestProperties.AsBasicAlertSimulatorRequestProperties() (BasicAlertSimulatorRequestProperties, bool)
1. AlertSimulatorBundlesRequestProperties.MarshalJSON() ([]byte, error)
1. AlertSimulatorRequestProperties.AsAlertSimulatorBundlesRequestProperties() (*AlertSimulatorBundlesRequestProperties, bool)
1. AlertSimulatorRequestProperties.AsAlertSimulatorRequestProperties() (*AlertSimulatorRequestProperties, bool)
1. AlertSimulatorRequestProperties.AsBasicAlertSimulatorRequestProperties() (BasicAlertSimulatorRequestProperties, bool)
1. AlertSimulatorRequestProperties.MarshalJSON() ([]byte, error)
1. AlertsClient.GetResourceGroupLevel(context.Context, string, string) (Alert, error)
1. AlertsClient.GetResourceGroupLevelPreparer(context.Context, string, string) (*http.Request, error)
1. AlertsClient.GetResourceGroupLevelResponder(*http.Response) (Alert, error)
1. AlertsClient.GetResourceGroupLevelSender(*http.Request) (*http.Response, error)
1. AlertsClient.GetSubscriptionLevel(context.Context, string) (Alert, error)
1. AlertsClient.GetSubscriptionLevelPreparer(context.Context, string) (*http.Request, error)
1. AlertsClient.GetSubscriptionLevelResponder(*http.Response) (Alert, error)
1. AlertsClient.GetSubscriptionLevelSender(*http.Request) (*http.Response, error)
1. AlertsClient.ListResourceGroupLevelByRegion(context.Context, string) (AlertListPage, error)
1. AlertsClient.ListResourceGroupLevelByRegionComplete(context.Context, string) (AlertListIterator, error)
1. AlertsClient.ListResourceGroupLevelByRegionPreparer(context.Context, string) (*http.Request, error)
1. AlertsClient.ListResourceGroupLevelByRegionResponder(*http.Response) (AlertList, error)
1. AlertsClient.ListResourceGroupLevelByRegionSender(*http.Request) (*http.Response, error)
1. AlertsClient.ListSubscriptionLevelByRegion(context.Context) (AlertListPage, error)
1. AlertsClient.ListSubscriptionLevelByRegionComplete(context.Context) (AlertListIterator, error)
1. AlertsClient.ListSubscriptionLevelByRegionPreparer(context.Context) (*http.Request, error)
1. AlertsClient.ListSubscriptionLevelByRegionResponder(*http.Response) (AlertList, error)
1. AlertsClient.ListSubscriptionLevelByRegionSender(*http.Request) (*http.Response, error)
1. AlertsClient.Simulate(context.Context, AlertSimulatorRequestBody) (AlertsSimulateFuture, error)
1. AlertsClient.SimulatePreparer(context.Context, AlertSimulatorRequestBody) (*http.Request, error)
1. AlertsClient.SimulateResponder(*http.Response) (autorest.Response, error)
1. AlertsClient.SimulateSender(*http.Request) (AlertsSimulateFuture, error)
1. AlertsClient.UpdateResourceGroupLevelStateToActivate(context.Context, string, string) (autorest.Response, error)
1. AlertsClient.UpdateResourceGroupLevelStateToActivatePreparer(context.Context, string, string) (*http.Request, error)
1. AlertsClient.UpdateResourceGroupLevelStateToActivateResponder(*http.Response) (autorest.Response, error)
1. AlertsClient.UpdateResourceGroupLevelStateToActivateSender(*http.Request) (*http.Response, error)
1. AlertsClient.UpdateResourceGroupLevelStateToDismiss(context.Context, string, string) (autorest.Response, error)
1. AlertsClient.UpdateResourceGroupLevelStateToDismissPreparer(context.Context, string, string) (*http.Request, error)
1. AlertsClient.UpdateResourceGroupLevelStateToDismissResponder(*http.Response) (autorest.Response, error)
1. AlertsClient.UpdateResourceGroupLevelStateToDismissSender(*http.Request) (*http.Response, error)
1. AlertsClient.UpdateSubscriptionLevelStateToActivate(context.Context, string) (autorest.Response, error)
1. AlertsClient.UpdateSubscriptionLevelStateToActivatePreparer(context.Context, string) (*http.Request, error)
1. AlertsClient.UpdateSubscriptionLevelStateToActivateResponder(*http.Response) (autorest.Response, error)
1. AlertsClient.UpdateSubscriptionLevelStateToActivateSender(*http.Request) (*http.Response, error)
1. AlertsClient.UpdateSubscriptionLevelStateToDismiss(context.Context, string) (autorest.Response, error)
1. AlertsClient.UpdateSubscriptionLevelStateToDismissPreparer(context.Context, string) (*http.Request, error)
1. AlertsClient.UpdateSubscriptionLevelStateToDismissResponder(*http.Response) (autorest.Response, error)
1. AlertsClient.UpdateSubscriptionLevelStateToDismissSender(*http.Request) (*http.Response, error)
1. AllowlistCustomAlertRule.AsConnectionFromIPNotAllowed() (*ConnectionFromIPNotAllowed, bool)
1. AmqpC2DMessagesNotInAllowedRange.AsConnectionFromIPNotAllowed() (*ConnectionFromIPNotAllowed, bool)
1. AmqpC2DRejectedMessagesNotInAllowedRange.AsConnectionFromIPNotAllowed() (*ConnectionFromIPNotAllowed, bool)
1. AmqpD2CMessagesNotInAllowedRange.AsConnectionFromIPNotAllowed() (*ConnectionFromIPNotAllowed, bool)
1. ConnectionFromIPNotAllowed.AsActiveConnectionsNotInAllowedRange() (*ActiveConnectionsNotInAllowedRange, bool)
1. ConnectionFromIPNotAllowed.AsAllowlistCustomAlertRule() (*AllowlistCustomAlertRule, bool)
1. ConnectionFromIPNotAllowed.AsAmqpC2DMessagesNotInAllowedRange() (*AmqpC2DMessagesNotInAllowedRange, bool)
1. ConnectionFromIPNotAllowed.AsAmqpC2DRejectedMessagesNotInAllowedRange() (*AmqpC2DRejectedMessagesNotInAllowedRange, bool)
1. ConnectionFromIPNotAllowed.AsAmqpD2CMessagesNotInAllowedRange() (*AmqpD2CMessagesNotInAllowedRange, bool)
1. ConnectionFromIPNotAllowed.AsBasicAllowlistCustomAlertRule() (BasicAllowlistCustomAlertRule, bool)
1. ConnectionFromIPNotAllowed.AsBasicCustomAlertRule() (BasicCustomAlertRule, bool)
1. ConnectionFromIPNotAllowed.AsBasicListCustomAlertRule() (BasicListCustomAlertRule, bool)
1. ConnectionFromIPNotAllowed.AsBasicThresholdCustomAlertRule() (BasicThresholdCustomAlertRule, bool)
1. ConnectionFromIPNotAllowed.AsBasicTimeWindowCustomAlertRule() (BasicTimeWindowCustomAlertRule, bool)
1. ConnectionFromIPNotAllowed.AsConnectionFromIPNotAllowed() (*ConnectionFromIPNotAllowed, bool)
1. ConnectionFromIPNotAllowed.AsConnectionToIPNotAllowed() (*ConnectionToIPNotAllowed, bool)
1. ConnectionFromIPNotAllowed.AsCustomAlertRule() (*CustomAlertRule, bool)
1. ConnectionFromIPNotAllowed.AsDenylistCustomAlertRule() (*DenylistCustomAlertRule, bool)
1. ConnectionFromIPNotAllowed.AsDirectMethodInvokesNotInAllowedRange() (*DirectMethodInvokesNotInAllowedRange, bool)
1. ConnectionFromIPNotAllowed.AsFailedLocalLoginsNotInAllowedRange() (*FailedLocalLoginsNotInAllowedRange, bool)
1. ConnectionFromIPNotAllowed.AsFileUploadsNotInAllowedRange() (*FileUploadsNotInAllowedRange, bool)
1. ConnectionFromIPNotAllowed.AsHTTPC2DMessagesNotInAllowedRange() (*HTTPC2DMessagesNotInAllowedRange, bool)
1. ConnectionFromIPNotAllowed.AsHTTPC2DRejectedMessagesNotInAllowedRange() (*HTTPC2DRejectedMessagesNotInAllowedRange, bool)
1. ConnectionFromIPNotAllowed.AsHTTPD2CMessagesNotInAllowedRange() (*HTTPD2CMessagesNotInAllowedRange, bool)
1. ConnectionFromIPNotAllowed.AsListCustomAlertRule() (*ListCustomAlertRule, bool)
1. ConnectionFromIPNotAllowed.AsLocalUserNotAllowed() (*LocalUserNotAllowed, bool)
1. ConnectionFromIPNotAllowed.AsMqttC2DMessagesNotInAllowedRange() (*MqttC2DMessagesNotInAllowedRange, bool)
1. ConnectionFromIPNotAllowed.AsMqttC2DRejectedMessagesNotInAllowedRange() (*MqttC2DRejectedMessagesNotInAllowedRange, bool)
1. ConnectionFromIPNotAllowed.AsMqttD2CMessagesNotInAllowedRange() (*MqttD2CMessagesNotInAllowedRange, bool)
1. ConnectionFromIPNotAllowed.AsProcessNotAllowed() (*ProcessNotAllowed, bool)
1. ConnectionFromIPNotAllowed.AsQueuePurgesNotInAllowedRange() (*QueuePurgesNotInAllowedRange, bool)
1. ConnectionFromIPNotAllowed.AsThresholdCustomAlertRule() (*ThresholdCustomAlertRule, bool)
1. ConnectionFromIPNotAllowed.AsTimeWindowCustomAlertRule() (*TimeWindowCustomAlertRule, bool)
1. ConnectionFromIPNotAllowed.AsTwinUpdatesNotInAllowedRange() (*TwinUpdatesNotInAllowedRange, bool)
1. ConnectionFromIPNotAllowed.AsUnauthorizedOperationsNotInAllowedRange() (*UnauthorizedOperationsNotInAllowedRange, bool)
1. ConnectionFromIPNotAllowed.MarshalJSON() ([]byte, error)
1. ConnectionToIPNotAllowed.AsConnectionFromIPNotAllowed() (*ConnectionFromIPNotAllowed, bool)
1. CustomAlertRule.AsConnectionFromIPNotAllowed() (*ConnectionFromIPNotAllowed, bool)
1. DenylistCustomAlertRule.AsConnectionFromIPNotAllowed() (*ConnectionFromIPNotAllowed, bool)
1. DirectMethodInvokesNotInAllowedRange.AsConnectionFromIPNotAllowed() (*ConnectionFromIPNotAllowed, bool)
1. ErrorAdditionalInfo.MarshalJSON() ([]byte, error)
1. FailedLocalLoginsNotInAllowedRange.AsConnectionFromIPNotAllowed() (*ConnectionFromIPNotAllowed, bool)
1. FileUploadsNotInAllowedRange.AsConnectionFromIPNotAllowed() (*ConnectionFromIPNotAllowed, bool)
1. HTTPC2DMessagesNotInAllowedRange.AsConnectionFromIPNotAllowed() (*ConnectionFromIPNotAllowed, bool)
1. HTTPC2DRejectedMessagesNotInAllowedRange.AsConnectionFromIPNotAllowed() (*ConnectionFromIPNotAllowed, bool)
1. HTTPD2CMessagesNotInAllowedRange.AsConnectionFromIPNotAllowed() (*ConnectionFromIPNotAllowed, bool)
1. IngestionConnectionString.MarshalJSON() ([]byte, error)
1. IngestionSetting.MarshalJSON() ([]byte, error)
1. IngestionSettingList.IsEmpty() bool
1. IngestionSettingList.MarshalJSON() ([]byte, error)
1. IngestionSettingListIterator.NotDone() bool
1. IngestionSettingListIterator.Response() IngestionSettingList
1. IngestionSettingListIterator.Value() IngestionSetting
1. IngestionSettingListPage.NotDone() bool
1. IngestionSettingListPage.Response() IngestionSettingList
1. IngestionSettingListPage.Values() []IngestionSetting
1. IngestionSettingToken.MarshalJSON() ([]byte, error)
1. IngestionSettingsClient.Create(context.Context, string, IngestionSetting) (IngestionSetting, error)
1. IngestionSettingsClient.CreatePreparer(context.Context, string, IngestionSetting) (*http.Request, error)
1. IngestionSettingsClient.CreateResponder(*http.Response) (IngestionSetting, error)
1. IngestionSettingsClient.CreateSender(*http.Request) (*http.Response, error)
1. IngestionSettingsClient.Delete(context.Context, string) (autorest.Response, error)
1. IngestionSettingsClient.DeletePreparer(context.Context, string) (*http.Request, error)
1. IngestionSettingsClient.DeleteResponder(*http.Response) (autorest.Response, error)
1. IngestionSettingsClient.DeleteSender(*http.Request) (*http.Response, error)
1. IngestionSettingsClient.Get(context.Context, string) (IngestionSetting, error)
1. IngestionSettingsClient.GetPreparer(context.Context, string) (*http.Request, error)
1. IngestionSettingsClient.GetResponder(*http.Response) (IngestionSetting, error)
1. IngestionSettingsClient.GetSender(*http.Request) (*http.Response, error)
1. IngestionSettingsClient.List(context.Context) (IngestionSettingListPage, error)
1. IngestionSettingsClient.ListComplete(context.Context) (IngestionSettingListIterator, error)
1. IngestionSettingsClient.ListConnectionStrings(context.Context, string) (ConnectionStrings, error)
1. IngestionSettingsClient.ListConnectionStringsPreparer(context.Context, string) (*http.Request, error)
1. IngestionSettingsClient.ListConnectionStringsResponder(*http.Response) (ConnectionStrings, error)
1. IngestionSettingsClient.ListConnectionStringsSender(*http.Request) (*http.Response, error)
1. IngestionSettingsClient.ListPreparer(context.Context) (*http.Request, error)
1. IngestionSettingsClient.ListResponder(*http.Response) (IngestionSettingList, error)
1. IngestionSettingsClient.ListSender(*http.Request) (*http.Response, error)
1. IngestionSettingsClient.ListTokens(context.Context, string) (IngestionSettingToken, error)
1. IngestionSettingsClient.ListTokensPreparer(context.Context, string) (*http.Request, error)
1. IngestionSettingsClient.ListTokensResponder(*http.Response) (IngestionSettingToken, error)
1. IngestionSettingsClient.ListTokensSender(*http.Request) (*http.Response, error)
1. ListCustomAlertRule.AsConnectionFromIPNotAllowed() (*ConnectionFromIPNotAllowed, bool)
1. LocalUserNotAllowed.AsConnectionFromIPNotAllowed() (*ConnectionFromIPNotAllowed, bool)
1. MqttC2DMessagesNotInAllowedRange.AsConnectionFromIPNotAllowed() (*ConnectionFromIPNotAllowed, bool)
1. MqttC2DRejectedMessagesNotInAllowedRange.AsConnectionFromIPNotAllowed() (*ConnectionFromIPNotAllowed, bool)
1. MqttD2CMessagesNotInAllowedRange.AsConnectionFromIPNotAllowed() (*ConnectionFromIPNotAllowed, bool)
1. NewIngestionSettingListIterator(IngestionSettingListPage) IngestionSettingListIterator
1. NewIngestionSettingListPage(IngestionSettingList, func(context.Context, IngestionSettingList) (IngestionSettingList, error)) IngestionSettingListPage
1. NewIngestionSettingsClient(string, string) IngestionSettingsClient
1. NewIngestionSettingsClientWithBaseURI(string, string, string) IngestionSettingsClient
1. NewSoftwareInventoriesClient(string, string) SoftwareInventoriesClient
1. NewSoftwareInventoriesClientWithBaseURI(string, string, string) SoftwareInventoriesClient
1. NewSoftwaresListIterator(SoftwaresListPage) SoftwaresListIterator
1. NewSoftwaresListPage(SoftwaresList, func(context.Context, SoftwaresList) (SoftwaresList, error)) SoftwaresListPage
1. PossibleAdditionalWorkspaceDataTypeValues() []AdditionalWorkspaceDataType
1. PossibleAdditionalWorkspaceTypeValues() []AdditionalWorkspaceType
1. PossibleBundleTypeValues() []BundleType
1. PossibleCategoriesValues() []Categories
1. PossibleCreatedByTypeValues() []CreatedByType
1. PossibleEndOfSupportStatusValues() []EndOfSupportStatus
1. PossibleKindEnum2Values() []KindEnum2
1. PossibleOnboardingKindValues() []OnboardingKind
1. PossibleSensorTypeValues() []SensorType
1. ProcessNotAllowed.AsConnectionFromIPNotAllowed() (*ConnectionFromIPNotAllowed, bool)
1. QueuePurgesNotInAllowedRange.AsConnectionFromIPNotAllowed() (*ConnectionFromIPNotAllowed, bool)
1. Software.MarshalJSON() ([]byte, error)
1. SoftwareInventoriesClient.Get(context.Context, string, string, string, string, string) (Software, error)
1. SoftwareInventoriesClient.GetPreparer(context.Context, string, string, string, string, string) (*http.Request, error)
1. SoftwareInventoriesClient.GetResponder(*http.Response) (Software, error)
1. SoftwareInventoriesClient.GetSender(*http.Request) (*http.Response, error)
1. SoftwareInventoriesClient.ListByExtendedResource(context.Context, string, string, string, string) (SoftwaresListPage, error)
1. SoftwareInventoriesClient.ListByExtendedResourceComplete(context.Context, string, string, string, string) (SoftwaresListIterator, error)
1. SoftwareInventoriesClient.ListByExtendedResourcePreparer(context.Context, string, string, string, string) (*http.Request, error)
1. SoftwareInventoriesClient.ListByExtendedResourceResponder(*http.Response) (SoftwaresList, error)
1. SoftwareInventoriesClient.ListByExtendedResourceSender(*http.Request) (*http.Response, error)
1. SoftwareInventoriesClient.ListBySubscription(context.Context) (SoftwaresListPage, error)
1. SoftwareInventoriesClient.ListBySubscriptionComplete(context.Context) (SoftwaresListIterator, error)
1. SoftwareInventoriesClient.ListBySubscriptionPreparer(context.Context) (*http.Request, error)
1. SoftwareInventoriesClient.ListBySubscriptionResponder(*http.Response) (SoftwaresList, error)
1. SoftwareInventoriesClient.ListBySubscriptionSender(*http.Request) (*http.Response, error)
1. SoftwaresList.IsEmpty() bool
1. SoftwaresList.MarshalJSON() ([]byte, error)
1. SoftwaresListIterator.NotDone() bool
1. SoftwaresListIterator.Response() SoftwaresList
1. SoftwaresListIterator.Value() Software
1. SoftwaresListPage.NotDone() bool
1. SoftwaresListPage.Response() SoftwaresList
1. SoftwaresListPage.Values() []Software
1. ThresholdCustomAlertRule.AsConnectionFromIPNotAllowed() (*ConnectionFromIPNotAllowed, bool)
1. TimeWindowCustomAlertRule.AsConnectionFromIPNotAllowed() (*ConnectionFromIPNotAllowed, bool)
1. TwinUpdatesNotInAllowedRange.AsConnectionFromIPNotAllowed() (*ConnectionFromIPNotAllowed, bool)
1. UnauthorizedOperationsNotInAllowedRange.AsConnectionFromIPNotAllowed() (*ConnectionFromIPNotAllowed, bool)

### Struct Changes

#### New Structs

1. AdditionalWorkspacesProperties
1. AlertSimulatorBundlesRequestProperties
1. AlertSimulatorRequestBody
1. AlertSimulatorRequestProperties
1. AlertsSimulateFuture
1. ConnectionFromIPNotAllowed
1. ConnectionStrings
1. ErrorAdditionalInfo
1. IngestionConnectionString
1. IngestionSetting
1. IngestionSettingList
1. IngestionSettingListIterator
1. IngestionSettingListPage
1. IngestionSettingToken
1. IngestionSettingsClient
1. ServerVulnerabilityAssessmentDeleteFuture
1. Software
1. SoftwareInventoriesClient
1. SoftwareProperties
1. SoftwaresList
1. SoftwaresListIterator
1. SoftwaresListPage
1. SystemData

#### New Struct Fields

1. AssessmentMetadataProperties.Categories
1. CloudErrorBody.AdditionalInfo
1. CloudErrorBody.Details
1. CloudErrorBody.Target
1. IoTSecuritySolutionModel.SystemData
1. IoTSecuritySolutionProperties.AdditionalWorkspaces
1. IotAlertModel.ID
1. IotAlertModel.Name
1. IotAlertModel.Type
1. IotDefenderSettingsProperties.OnboardingKind
1. IotSensorProperties.SensorType
