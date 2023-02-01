# Change History

## Breaking Changes

### Removed Constants

1. AlertIntent.CommandAndControl
1. AlertIntent.Exploitation
1. AlertIntent.PreAttack
1. AlertIntent.Probing
1. AlertsToAdmins.AlertsToAdminsOff
1. AlertsToAdmins.AlertsToAdminsOn
1. AuthorizationState.Authorized
1. AuthorizationState.Unauthorized
1. AutoProvision.AutoProvisionOff
1. AutoProvision.AutoProvisionOn
1. ConfigurationStatus.Configured
1. ConfigurationStatus.NoStatus
1. ConfigurationStatus.NotConfigured
1. CreatedByType.Application
1. CreatedByType.Key
1. CreatedByType.ManagedIdentity
1. CreatedByType.User
1. DeviceCriticality.Important
1. DeviceStatus.DeviceStatusActive
1. DeviceStatus.DeviceStatusRemoved
1. EnforcementSupport.EnforcementSupportNotSupported
1. EnforcementSupport.EnforcementSupportSupported
1. EnforcementSupport.EnforcementSupportUnknown
1. MacSignificance.Primary
1. MacSignificance.Secondary
1. ManagementState.Managed
1. ManagementState.Unmanaged
1. OnboardingKind.Evaluation
1. OnboardingKind.MigratedToAzure
1. OnboardingKind.Purchased
1. Operator.Contains
1. Operator.EndsWith
1. Operator.GreaterThan
1. Operator.GreaterThanOrEqualTo
1. Operator.LesserThan
1. Operator.LesserThanOrEqualTo
1. Operator.NotEquals
1. Operator.StartsWith
1. PricingTier.PricingTierFree
1. PricingTier.PricingTierStandard
1. ProgrammingState.NotProgrammingDevice
1. ProgrammingState.ProgrammingDevice
1. ProvisioningState.ProvisioningStateFailed
1. ProvisioningState.ProvisioningStateSucceeded
1. ProvisioningState.ProvisioningStateUpdating
1. PurdueLevel.Enterprise
1. PurdueLevel.ProcessControl
1. PurdueLevel.Supervisory
1. RecommendationSeverity.RecommendationSeverityHealthy
1. RecommendationSeverity.RecommendationSeverityHigh
1. RecommendationSeverity.RecommendationSeverityLow
1. RecommendationSeverity.RecommendationSeverityMedium
1. RecommendationSeverity.RecommendationSeverityNotApplicable
1. RecommendationSeverity.RecommendationSeverityOffByPolicy
1. RecommendationSeverity.RecommendationSeverityUnknown
1. RelationToIPStatus.Certain
1. RelationToIPStatus.Guess
1. ScanningFunctionality.NotScannerDevice
1. ScanningFunctionality.ScannerDevice
1. SensorStatus.Disconnected
1. SensorStatus.Ok
1. SensorStatus.Unavailable
1. SensorType.SensorTypeEnterprise
1. SensorType.SensorTypeOt
1. TiStatus.TiStatusFailed
1. TiStatus.TiStatusInProgress
1. TiStatus.TiStatusOk
1. TiStatus.TiStatusUpdateAvailable
1. VersionKind.Latest
1. VersionKind.Preview
1. VersionKind.Previous

### Removed Type Aliases

1. string.AlertIntent
1. string.AlertNotifications
1. string.AlertsToAdmins
1. string.AuthorizationState
1. string.DeviceCriticality
1. string.DeviceStatus
1. string.MacSignificance
1. string.ManagementState
1. string.OnboardingKind
1. string.ProgrammingState
1. string.PurdueLevel
1. string.RecommendationSeverity
1. string.RelationToIPStatus
1. string.ScanningFunctionality
1. string.SensorStatus
1. string.SensorType
1. string.TiStatus
1. string.VersionKind

### Removed Funcs

1. *AssessmentMetadataListIterator.Next() error
1. *AssessmentMetadataListIterator.NextWithContext(context.Context) error
1. *AssessmentMetadataListPage.Next() error
1. *AssessmentMetadataListPage.NextWithContext(context.Context) error
1. *Device.UnmarshalJSON([]byte) error
1. *DeviceListIterator.Next() error
1. *DeviceListIterator.NextWithContext(context.Context) error
1. *DeviceListPage.Next() error
1. *DeviceListPage.NextWithContext(context.Context) error
1. *IotAlertListModelIterator.Next() error
1. *IotAlertListModelIterator.NextWithContext(context.Context) error
1. *IotAlertListModelPage.Next() error
1. *IotAlertListModelPage.NextWithContext(context.Context) error
1. *IotAlertModel.UnmarshalJSON([]byte) error
1. *IotAlertType.UnmarshalJSON([]byte) error
1. *IotDefenderSettingsModel.UnmarshalJSON([]byte) error
1. *IotRecommendationListModelIterator.Next() error
1. *IotRecommendationListModelIterator.NextWithContext(context.Context) error
1. *IotRecommendationListModelPage.Next() error
1. *IotRecommendationListModelPage.NextWithContext(context.Context) error
1. *IotRecommendationModel.UnmarshalJSON([]byte) error
1. *IotRecommendationType.UnmarshalJSON([]byte) error
1. *IotSensorsModel.UnmarshalJSON([]byte) error
1. *IotSitesModel.UnmarshalJSON([]byte) error
1. AssessmentMetadataList.IsEmpty() bool
1. AssessmentMetadataList.MarshalJSON() ([]byte, error)
1. AssessmentMetadataListIterator.NotDone() bool
1. AssessmentMetadataListIterator.Response() AssessmentMetadataList
1. AssessmentMetadataListIterator.Value() AssessmentMetadata
1. AssessmentMetadataListPage.NotDone() bool
1. AssessmentMetadataListPage.Response() AssessmentMetadataList
1. AssessmentMetadataListPage.Values() []AssessmentMetadata
1. ContactsClient.Update(context.Context, string, Contact) (Contact, error)
1. ContactsClient.UpdatePreparer(context.Context, string, Contact) (*http.Request, error)
1. ContactsClient.UpdateResponder(*http.Response) (Contact, error)
1. ContactsClient.UpdateSender(*http.Request) (*http.Response, error)
1. Device.MarshalJSON() ([]byte, error)
1. DeviceClient.Get(context.Context, string, string) (Device, error)
1. DeviceClient.GetPreparer(context.Context, string, string) (*http.Request, error)
1. DeviceClient.GetResponder(*http.Response) (Device, error)
1. DeviceClient.GetSender(*http.Request) (*http.Response, error)
1. DeviceList.IsEmpty() bool
1. DeviceList.MarshalJSON() ([]byte, error)
1. DeviceListIterator.NotDone() bool
1. DeviceListIterator.Response() DeviceList
1. DeviceListIterator.Value() Device
1. DeviceListPage.NotDone() bool
1. DeviceListPage.Response() DeviceList
1. DeviceListPage.Values() []Device
1. DeviceProperties.MarshalJSON() ([]byte, error)
1. DevicesForHubClient.List(context.Context, string, *int32, string, ManagementState) (DeviceListPage, error)
1. DevicesForHubClient.ListComplete(context.Context, string, *int32, string, ManagementState) (DeviceListIterator, error)
1. DevicesForHubClient.ListPreparer(context.Context, string, *int32, string, ManagementState) (*http.Request, error)
1. DevicesForHubClient.ListResponder(*http.Response) (DeviceList, error)
1. DevicesForHubClient.ListSender(*http.Request) (*http.Response, error)
1. DevicesForSubscriptionClient.List(context.Context, *int32, string, ManagementState) (DeviceListPage, error)
1. DevicesForSubscriptionClient.ListComplete(context.Context, *int32, string, ManagementState) (DeviceListIterator, error)
1. DevicesForSubscriptionClient.ListPreparer(context.Context, *int32, string, ManagementState) (*http.Request, error)
1. DevicesForSubscriptionClient.ListResponder(*http.Response) (DeviceList, error)
1. DevicesForSubscriptionClient.ListSender(*http.Request) (*http.Response, error)
1. Firmware.MarshalJSON() ([]byte, error)
1. IPAddress.MarshalJSON() ([]byte, error)
1. IotAlertListModel.IsEmpty() bool
1. IotAlertListModel.MarshalJSON() ([]byte, error)
1. IotAlertListModelIterator.NotDone() bool
1. IotAlertListModelIterator.Response() IotAlertListModel
1. IotAlertListModelIterator.Value() IotAlertModel
1. IotAlertListModelPage.NotDone() bool
1. IotAlertListModelPage.Response() IotAlertListModel
1. IotAlertListModelPage.Values() []IotAlertModel
1. IotAlertModel.MarshalJSON() ([]byte, error)
1. IotAlertPropertiesModel.MarshalJSON() ([]byte, error)
1. IotAlertType.MarshalJSON() ([]byte, error)
1. IotAlertTypeProperties.MarshalJSON() ([]byte, error)
1. IotAlertTypesClient.Get(context.Context, string) (IotAlertType, error)
1. IotAlertTypesClient.GetPreparer(context.Context, string) (*http.Request, error)
1. IotAlertTypesClient.GetResponder(*http.Response) (IotAlertType, error)
1. IotAlertTypesClient.GetSender(*http.Request) (*http.Response, error)
1. IotAlertTypesClient.List(context.Context) (IotAlertTypeList, error)
1. IotAlertTypesClient.ListPreparer(context.Context) (*http.Request, error)
1. IotAlertTypesClient.ListResponder(*http.Response) (IotAlertTypeList, error)
1. IotAlertTypesClient.ListSender(*http.Request) (*http.Response, error)
1. IotAlertsClient.Get(context.Context, string, string) (IotAlertModel, error)
1. IotAlertsClient.GetPreparer(context.Context, string, string) (*http.Request, error)
1. IotAlertsClient.GetResponder(*http.Response) (IotAlertModel, error)
1. IotAlertsClient.GetSender(*http.Request) (*http.Response, error)
1. IotAlertsClient.List(context.Context, string, string, string, string, ManagementState, string, *int32, string) (IotAlertListModelPage, error)
1. IotAlertsClient.ListComplete(context.Context, string, string, string, string, ManagementState, string, *int32, string) (IotAlertListModelIterator, error)
1. IotAlertsClient.ListPreparer(context.Context, string, string, string, string, ManagementState, string, *int32, string) (*http.Request, error)
1. IotAlertsClient.ListResponder(*http.Response) (IotAlertListModel, error)
1. IotAlertsClient.ListSender(*http.Request) (*http.Response, error)
1. IotDefenderSettingsClient.CreateOrUpdate(context.Context, IotDefenderSettingsModel) (IotDefenderSettingsModel, error)
1. IotDefenderSettingsClient.CreateOrUpdatePreparer(context.Context, IotDefenderSettingsModel) (*http.Request, error)
1. IotDefenderSettingsClient.CreateOrUpdateResponder(*http.Response) (IotDefenderSettingsModel, error)
1. IotDefenderSettingsClient.CreateOrUpdateSender(*http.Request) (*http.Response, error)
1. IotDefenderSettingsClient.Delete(context.Context) (autorest.Response, error)
1. IotDefenderSettingsClient.DeletePreparer(context.Context) (*http.Request, error)
1. IotDefenderSettingsClient.DeleteResponder(*http.Response) (autorest.Response, error)
1. IotDefenderSettingsClient.DeleteSender(*http.Request) (*http.Response, error)
1. IotDefenderSettingsClient.DownloadManagerActivation(context.Context) (ReadCloser, error)
1. IotDefenderSettingsClient.DownloadManagerActivationPreparer(context.Context) (*http.Request, error)
1. IotDefenderSettingsClient.DownloadManagerActivationResponder(*http.Response) (ReadCloser, error)
1. IotDefenderSettingsClient.DownloadManagerActivationSender(*http.Request) (*http.Response, error)
1. IotDefenderSettingsClient.Get(context.Context) (IotDefenderSettingsModel, error)
1. IotDefenderSettingsClient.GetPreparer(context.Context) (*http.Request, error)
1. IotDefenderSettingsClient.GetResponder(*http.Response) (IotDefenderSettingsModel, error)
1. IotDefenderSettingsClient.GetSender(*http.Request) (*http.Response, error)
1. IotDefenderSettingsClient.List(context.Context) (IotDefenderSettingsList, error)
1. IotDefenderSettingsClient.ListPreparer(context.Context) (*http.Request, error)
1. IotDefenderSettingsClient.ListResponder(*http.Response) (IotDefenderSettingsList, error)
1. IotDefenderSettingsClient.ListSender(*http.Request) (*http.Response, error)
1. IotDefenderSettingsClient.PackageDownloadsMethod(context.Context) (PackageDownloads, error)
1. IotDefenderSettingsClient.PackageDownloadsMethodPreparer(context.Context) (*http.Request, error)
1. IotDefenderSettingsClient.PackageDownloadsMethodResponder(*http.Response) (PackageDownloads, error)
1. IotDefenderSettingsClient.PackageDownloadsMethodSender(*http.Request) (*http.Response, error)
1. IotDefenderSettingsList.MarshalJSON() ([]byte, error)
1. IotDefenderSettingsModel.MarshalJSON() ([]byte, error)
1. IotDefenderSettingsProperties.MarshalJSON() ([]byte, error)
1. IotRecommendationListModel.IsEmpty() bool
1. IotRecommendationListModel.MarshalJSON() ([]byte, error)
1. IotRecommendationListModelIterator.NotDone() bool
1. IotRecommendationListModelIterator.Response() IotRecommendationListModel
1. IotRecommendationListModelIterator.Value() IotRecommendationModel
1. IotRecommendationListModelPage.NotDone() bool
1. IotRecommendationListModelPage.Response() IotRecommendationListModel
1. IotRecommendationListModelPage.Values() []IotRecommendationModel
1. IotRecommendationModel.MarshalJSON() ([]byte, error)
1. IotRecommendationPropertiesModel.MarshalJSON() ([]byte, error)
1. IotRecommendationType.MarshalJSON() ([]byte, error)
1. IotRecommendationTypeProperties.MarshalJSON() ([]byte, error)
1. IotRecommendationTypesClient.Get(context.Context, string) (IotRecommendationType, error)
1. IotRecommendationTypesClient.GetPreparer(context.Context, string) (*http.Request, error)
1. IotRecommendationTypesClient.GetResponder(*http.Response) (IotRecommendationType, error)
1. IotRecommendationTypesClient.GetSender(*http.Request) (*http.Response, error)
1. IotRecommendationTypesClient.List(context.Context) (IotRecommendationTypeList, error)
1. IotRecommendationTypesClient.ListPreparer(context.Context) (*http.Request, error)
1. IotRecommendationTypesClient.ListResponder(*http.Response) (IotRecommendationTypeList, error)
1. IotRecommendationTypesClient.ListSender(*http.Request) (*http.Response, error)
1. IotRecommendationsClient.Get(context.Context, string, string) (IotRecommendationModel, error)
1. IotRecommendationsClient.GetPreparer(context.Context, string, string) (*http.Request, error)
1. IotRecommendationsClient.GetResponder(*http.Response) (IotRecommendationModel, error)
1. IotRecommendationsClient.GetSender(*http.Request) (*http.Response, error)
1. IotRecommendationsClient.List(context.Context, string, string, string, *int32, string) (IotRecommendationListModelPage, error)
1. IotRecommendationsClient.ListComplete(context.Context, string, string, string, *int32, string) (IotRecommendationListModelIterator, error)
1. IotRecommendationsClient.ListPreparer(context.Context, string, string, string, *int32, string) (*http.Request, error)
1. IotRecommendationsClient.ListResponder(*http.Response) (IotRecommendationListModel, error)
1. IotRecommendationsClient.ListSender(*http.Request) (*http.Response, error)
1. IotSensorProperties.MarshalJSON() ([]byte, error)
1. IotSensorsClient.CreateOrUpdate(context.Context, string, string, IotSensorsModel) (IotSensorsModel, error)
1. IotSensorsClient.CreateOrUpdatePreparer(context.Context, string, string, IotSensorsModel) (*http.Request, error)
1. IotSensorsClient.CreateOrUpdateResponder(*http.Response) (IotSensorsModel, error)
1. IotSensorsClient.CreateOrUpdateSender(*http.Request) (*http.Response, error)
1. IotSensorsClient.Delete(context.Context, string, string) (autorest.Response, error)
1. IotSensorsClient.DeletePreparer(context.Context, string, string) (*http.Request, error)
1. IotSensorsClient.DeleteResponder(*http.Response) (autorest.Response, error)
1. IotSensorsClient.DeleteSender(*http.Request) (*http.Response, error)
1. IotSensorsClient.DownloadActivation(context.Context, string, string) (ReadCloser, error)
1. IotSensorsClient.DownloadActivationPreparer(context.Context, string, string) (*http.Request, error)
1. IotSensorsClient.DownloadActivationResponder(*http.Response) (ReadCloser, error)
1. IotSensorsClient.DownloadActivationSender(*http.Request) (*http.Response, error)
1. IotSensorsClient.DownloadResetPassword(context.Context, string, string, ResetPasswordInput) (ReadCloser, error)
1. IotSensorsClient.DownloadResetPasswordPreparer(context.Context, string, string, ResetPasswordInput) (*http.Request, error)
1. IotSensorsClient.DownloadResetPasswordResponder(*http.Response) (ReadCloser, error)
1. IotSensorsClient.DownloadResetPasswordSender(*http.Request) (*http.Response, error)
1. IotSensorsClient.Get(context.Context, string, string) (IotSensorsModel, error)
1. IotSensorsClient.GetPreparer(context.Context, string, string) (*http.Request, error)
1. IotSensorsClient.GetResponder(*http.Response) (IotSensorsModel, error)
1. IotSensorsClient.GetSender(*http.Request) (*http.Response, error)
1. IotSensorsClient.List(context.Context, string) (IotSensorsList, error)
1. IotSensorsClient.ListPreparer(context.Context, string) (*http.Request, error)
1. IotSensorsClient.ListResponder(*http.Response) (IotSensorsList, error)
1. IotSensorsClient.ListSender(*http.Request) (*http.Response, error)
1. IotSensorsClient.TriggerTiPackageUpdate(context.Context, string, string) (autorest.Response, error)
1. IotSensorsClient.TriggerTiPackageUpdatePreparer(context.Context, string, string) (*http.Request, error)
1. IotSensorsClient.TriggerTiPackageUpdateResponder(*http.Response) (autorest.Response, error)
1. IotSensorsClient.TriggerTiPackageUpdateSender(*http.Request) (*http.Response, error)
1. IotSensorsList.MarshalJSON() ([]byte, error)
1. IotSensorsModel.MarshalJSON() ([]byte, error)
1. IotSiteProperties.MarshalJSON() ([]byte, error)
1. IotSitesClient.CreateOrUpdate(context.Context, string, IotSitesModel) (IotSitesModel, error)
1. IotSitesClient.CreateOrUpdatePreparer(context.Context, string, IotSitesModel) (*http.Request, error)
1. IotSitesClient.CreateOrUpdateResponder(*http.Response) (IotSitesModel, error)
1. IotSitesClient.CreateOrUpdateSender(*http.Request) (*http.Response, error)
1. IotSitesClient.Delete(context.Context, string) (autorest.Response, error)
1. IotSitesClient.DeletePreparer(context.Context, string) (*http.Request, error)
1. IotSitesClient.DeleteResponder(*http.Response) (autorest.Response, error)
1. IotSitesClient.DeleteSender(*http.Request) (*http.Response, error)
1. IotSitesClient.Get(context.Context, string) (IotSitesModel, error)
1. IotSitesClient.GetPreparer(context.Context, string) (*http.Request, error)
1. IotSitesClient.GetResponder(*http.Response) (IotSitesModel, error)
1. IotSitesClient.GetSender(*http.Request) (*http.Response, error)
1. IotSitesClient.List(context.Context, string) (IotSitesList, error)
1. IotSitesClient.ListPreparer(context.Context, string) (*http.Request, error)
1. IotSitesClient.ListResponder(*http.Response) (IotSitesList, error)
1. IotSitesClient.ListSender(*http.Request) (*http.Response, error)
1. IotSitesList.MarshalJSON() ([]byte, error)
1. IotSitesModel.MarshalJSON() ([]byte, error)
1. MacAddress.MarshalJSON() ([]byte, error)
1. NetworkInterface.MarshalJSON() ([]byte, error)
1. NewAssessmentMetadataListIterator(AssessmentMetadataListPage) AssessmentMetadataListIterator
1. NewAssessmentMetadataListPage(AssessmentMetadataList, func(context.Context, AssessmentMetadataList) (AssessmentMetadataList, error)) AssessmentMetadataListPage
1. NewDeviceClient(string, string) DeviceClient
1. NewDeviceClientWithBaseURI(string, string, string) DeviceClient
1. NewDeviceListIterator(DeviceListPage) DeviceListIterator
1. NewDeviceListPage(DeviceList, func(context.Context, DeviceList) (DeviceList, error)) DeviceListPage
1. NewDevicesForHubClient(string, string) DevicesForHubClient
1. NewDevicesForHubClientWithBaseURI(string, string, string) DevicesForHubClient
1. NewDevicesForSubscriptionClient(string, string) DevicesForSubscriptionClient
1. NewDevicesForSubscriptionClientWithBaseURI(string, string, string) DevicesForSubscriptionClient
1. NewIotAlertListModelIterator(IotAlertListModelPage) IotAlertListModelIterator
1. NewIotAlertListModelPage(IotAlertListModel, func(context.Context, IotAlertListModel) (IotAlertListModel, error)) IotAlertListModelPage
1. NewIotAlertTypesClient(string, string) IotAlertTypesClient
1. NewIotAlertTypesClientWithBaseURI(string, string, string) IotAlertTypesClient
1. NewIotAlertsClient(string, string) IotAlertsClient
1. NewIotAlertsClientWithBaseURI(string, string, string) IotAlertsClient
1. NewIotDefenderSettingsClient(string, string) IotDefenderSettingsClient
1. NewIotDefenderSettingsClientWithBaseURI(string, string, string) IotDefenderSettingsClient
1. NewIotRecommendationListModelIterator(IotRecommendationListModelPage) IotRecommendationListModelIterator
1. NewIotRecommendationListModelPage(IotRecommendationListModel, func(context.Context, IotRecommendationListModel) (IotRecommendationListModel, error)) IotRecommendationListModelPage
1. NewIotRecommendationTypesClient(string, string) IotRecommendationTypesClient
1. NewIotRecommendationTypesClientWithBaseURI(string, string, string) IotRecommendationTypesClient
1. NewIotRecommendationsClient(string, string) IotRecommendationsClient
1. NewIotRecommendationsClientWithBaseURI(string, string, string) IotRecommendationsClient
1. NewIotSensorsClient(string, string) IotSensorsClient
1. NewIotSensorsClientWithBaseURI(string, string, string) IotSensorsClient
1. NewIotSitesClient(string, string) IotSitesClient
1. NewIotSitesClientWithBaseURI(string, string, string) IotSitesClient
1. NewOnPremiseIotSensorsClient(string, string) OnPremiseIotSensorsClient
1. NewOnPremiseIotSensorsClientWithBaseURI(string, string, string) OnPremiseIotSensorsClient
1. OnPremiseIotSensor.MarshalJSON() ([]byte, error)
1. OnPremiseIotSensorsClient.CreateOrUpdate(context.Context, string) (OnPremiseIotSensor, error)
1. OnPremiseIotSensorsClient.CreateOrUpdatePreparer(context.Context, string) (*http.Request, error)
1. OnPremiseIotSensorsClient.CreateOrUpdateResponder(*http.Response) (OnPremiseIotSensor, error)
1. OnPremiseIotSensorsClient.CreateOrUpdateSender(*http.Request) (*http.Response, error)
1. OnPremiseIotSensorsClient.Delete(context.Context, string) (autorest.Response, error)
1. OnPremiseIotSensorsClient.DeletePreparer(context.Context, string) (*http.Request, error)
1. OnPremiseIotSensorsClient.DeleteResponder(*http.Response) (autorest.Response, error)
1. OnPremiseIotSensorsClient.DeleteSender(*http.Request) (*http.Response, error)
1. OnPremiseIotSensorsClient.DownloadActivation(context.Context, string) (ReadCloser, error)
1. OnPremiseIotSensorsClient.DownloadActivationPreparer(context.Context, string) (*http.Request, error)
1. OnPremiseIotSensorsClient.DownloadActivationResponder(*http.Response) (ReadCloser, error)
1. OnPremiseIotSensorsClient.DownloadActivationSender(*http.Request) (*http.Response, error)
1. OnPremiseIotSensorsClient.DownloadResetPassword(context.Context, string, ResetPasswordInput) (ReadCloser, error)
1. OnPremiseIotSensorsClient.DownloadResetPasswordPreparer(context.Context, string, ResetPasswordInput) (*http.Request, error)
1. OnPremiseIotSensorsClient.DownloadResetPasswordResponder(*http.Response) (ReadCloser, error)
1. OnPremiseIotSensorsClient.DownloadResetPasswordSender(*http.Request) (*http.Response, error)
1. OnPremiseIotSensorsClient.Get(context.Context, string) (OnPremiseIotSensor, error)
1. OnPremiseIotSensorsClient.GetPreparer(context.Context, string) (*http.Request, error)
1. OnPremiseIotSensorsClient.GetResponder(*http.Response) (OnPremiseIotSensor, error)
1. OnPremiseIotSensorsClient.GetSender(*http.Request) (*http.Response, error)
1. OnPremiseIotSensorsClient.List(context.Context) (OnPremiseIotSensorsList, error)
1. OnPremiseIotSensorsClient.ListPreparer(context.Context) (*http.Request, error)
1. OnPremiseIotSensorsClient.ListResponder(*http.Response) (OnPremiseIotSensorsList, error)
1. OnPremiseIotSensorsClient.ListSender(*http.Request) (*http.Response, error)
1. OnPremiseIotSensorsList.MarshalJSON() ([]byte, error)
1. PackageDownloadInfo.MarshalJSON() ([]byte, error)
1. PackageDownloads.MarshalJSON() ([]byte, error)
1. PackageDownloadsCentralManager.MarshalJSON() ([]byte, error)
1. PackageDownloadsCentralManagerFull.MarshalJSON() ([]byte, error)
1. PackageDownloadsCentralManagerFullOvf.MarshalJSON() ([]byte, error)
1. PackageDownloadsSensor.MarshalJSON() ([]byte, error)
1. PackageDownloadsSensorFull.MarshalJSON() ([]byte, error)
1. PackageDownloadsSensorFullOvf.MarshalJSON() ([]byte, error)
1. PossibleAlertIntentValues() []AlertIntent
1. PossibleAlertNotificationsValues() []AlertNotifications
1. PossibleAlertsToAdminsValues() []AlertsToAdmins
1. PossibleAuthorizationStateValues() []AuthorizationState
1. PossibleDeviceCriticalityValues() []DeviceCriticality
1. PossibleDeviceStatusValues() []DeviceStatus
1. PossibleMacSignificanceValues() []MacSignificance
1. PossibleManagementStateValues() []ManagementState
1. PossibleOnboardingKindValues() []OnboardingKind
1. PossibleProgrammingStateValues() []ProgrammingState
1. PossiblePurdueLevelValues() []PurdueLevel
1. PossibleRecommendationSeverityValues() []RecommendationSeverity
1. PossibleRelationToIPStatusValues() []RelationToIPStatus
1. PossibleScanningFunctionalityValues() []ScanningFunctionality
1. PossibleSensorStatusValues() []SensorStatus
1. PossibleSensorTypeValues() []SensorType
1. PossibleTiStatusValues() []TiStatus
1. PossibleVersionKindValues() []VersionKind
1. Protocol1.MarshalJSON() ([]byte, error)
1. Sensor.MarshalJSON() ([]byte, error)
1. Site.MarshalJSON() ([]byte, error)
1. UpgradePackageDownloadInfo.MarshalJSON() ([]byte, error)

### Struct Changes

#### Removed Structs

1. AssessmentMetadataList
1. AssessmentMetadataListIterator
1. AssessmentMetadataListPage
1. Device
1. DeviceClient
1. DeviceList
1. DeviceListIterator
1. DeviceListPage
1. DeviceProperties
1. DevicesForHubClient
1. DevicesForSubscriptionClient
1. Firmware
1. IPAddress
1. IotAlertListModel
1. IotAlertListModelIterator
1. IotAlertListModelPage
1. IotAlertModel
1. IotAlertPropertiesModel
1. IotAlertType
1. IotAlertTypeList
1. IotAlertTypeProperties
1. IotAlertTypesClient
1. IotAlertsClient
1. IotDefenderSettingsClient
1. IotDefenderSettingsList
1. IotDefenderSettingsModel
1. IotDefenderSettingsProperties
1. IotRecommendationListModel
1. IotRecommendationListModelIterator
1. IotRecommendationListModelPage
1. IotRecommendationModel
1. IotRecommendationPropertiesModel
1. IotRecommendationType
1. IotRecommendationTypeList
1. IotRecommendationTypeProperties
1. IotRecommendationTypesClient
1. IotRecommendationsClient
1. IotSensorProperties
1. IotSensorsClient
1. IotSensorsList
1. IotSensorsModel
1. IotSiteProperties
1. IotSitesClient
1. IotSitesList
1. IotSitesModel
1. MacAddress
1. NetworkInterface
1. OnPremiseIotSensor
1. OnPremiseIotSensorsClient
1. OnPremiseIotSensorsList
1. PackageDownloadInfo
1. PackageDownloads
1. PackageDownloadsCentralManager
1. PackageDownloadsCentralManagerFull
1. PackageDownloadsCentralManagerFullOvf
1. PackageDownloadsSensor
1. PackageDownloadsSensorFull
1. PackageDownloadsSensorFullOvf
1. Protocol1
1. ReadCloser
1. ResetPasswordInput
1. Sensor
1. Site
1. UpgradePackageDownloadInfo

#### Removed Struct Fields

1. Assessment.autorest.Response
1. AssessmentMetadata.autorest.Response
1. BaseClient.AscLocation
1. ContactProperties.AlertsToAdmins
1. ContactProperties.Email

### Signature Changes

#### Const Types

1. BinarySignature changed type from Type to Type1
1. Collection changed type from AlertIntent to Tactics
1. CredentialAccess changed type from AlertIntent to Tactics
1. Default changed type from OnboardingKind to ScanningMode
1. DefenseEvasion changed type from AlertIntent to Tactics
1. Discovery changed type from AlertIntent to Tactics
1. Equals changed type from Operator to ApplicationConditionOperator
1. Execution changed type from AlertIntent to Tactics
1. Exfiltration changed type from AlertIntent to Tactics
1. Failed changed type from ConfigurationStatus to ProvisioningState
1. File changed type from Type to Type1
1. FileHash changed type from Type to Type1
1. Impact changed type from AlertIntent to Tactics
1. InProgress changed type from ConfigurationStatus to AlertStatus
1. InitialAccess changed type from AlertIntent to Tactics
1. LateralMovement changed type from AlertIntent to Tactics
1. Off changed type from AlertNotifications to AutoProvision
1. On changed type from AlertNotifications to AutoProvision
1. Persistence changed type from AlertIntent to Tactics
1. PrivilegeEscalation changed type from AlertIntent to Tactics
1. ProductSignature changed type from Type to Type1
1. PublisherSignature changed type from Type to Type1
1. Standard changed type from DeviceCriticality to PricingTier
1. Unknown changed type from AlertIntent to EnforcementSupport
1. VersionAndAboveSignature changed type from Type to Type1

#### Funcs

1. AdaptiveApplicationControlsClient.Delete
	- Params
		- From: context.Context, string
		- To: context.Context, string, string
1. AdaptiveApplicationControlsClient.DeletePreparer
	- Params
		- From: context.Context, string
		- To: context.Context, string, string
1. AdaptiveApplicationControlsClient.Get
	- Params
		- From: context.Context, string
		- To: context.Context, string, string
1. AdaptiveApplicationControlsClient.GetPreparer
	- Params
		- From: context.Context, string
		- To: context.Context, string, string
1. AdaptiveApplicationControlsClient.Put
	- Params
		- From: context.Context, string, AdaptiveApplicationControlGroup
		- To: context.Context, string, string, AdaptiveApplicationControlGroup
1. AdaptiveApplicationControlsClient.PutPreparer
	- Params
		- From: context.Context, string, AdaptiveApplicationControlGroup
		- To: context.Context, string, string, AdaptiveApplicationControlGroup
1. AlertsClient.GetResourceGroupLevel
	- Params
		- From: context.Context, string, string
		- To: context.Context, string, string, string
1. AlertsClient.GetResourceGroupLevelPreparer
	- Params
		- From: context.Context, string, string
		- To: context.Context, string, string, string
1. AlertsClient.GetSubscriptionLevel
	- Params
		- From: context.Context, string
		- To: context.Context, string, string
1. AlertsClient.GetSubscriptionLevelPreparer
	- Params
		- From: context.Context, string
		- To: context.Context, string, string
1. AlertsClient.ListResourceGroupLevelByRegion
	- Params
		- From: context.Context, string
		- To: context.Context, string, string
1. AlertsClient.ListResourceGroupLevelByRegionComplete
	- Params
		- From: context.Context, string
		- To: context.Context, string, string
1. AlertsClient.ListResourceGroupLevelByRegionPreparer
	- Params
		- From: context.Context, string
		- To: context.Context, string, string
1. AlertsClient.ListSubscriptionLevelByRegion
	- Params
		- From: context.Context
		- To: context.Context, string
1. AlertsClient.ListSubscriptionLevelByRegionComplete
	- Params
		- From: context.Context
		- To: context.Context, string
1. AlertsClient.ListSubscriptionLevelByRegionPreparer
	- Params
		- From: context.Context
		- To: context.Context, string
1. AlertsClient.Simulate
	- Params
		- From: context.Context, AlertSimulatorRequestBody
		- To: context.Context, string, AlertSimulatorRequestBody
1. AlertsClient.SimulatePreparer
	- Params
		- From: context.Context, AlertSimulatorRequestBody
		- To: context.Context, string, AlertSimulatorRequestBody
1. AlertsClient.UpdateResourceGroupLevelStateToActivate
	- Params
		- From: context.Context, string, string
		- To: context.Context, string, string, string
1. AlertsClient.UpdateResourceGroupLevelStateToActivatePreparer
	- Params
		- From: context.Context, string, string
		- To: context.Context, string, string, string
1. AlertsClient.UpdateResourceGroupLevelStateToDismiss
	- Params
		- From: context.Context, string, string
		- To: context.Context, string, string, string
1. AlertsClient.UpdateResourceGroupLevelStateToDismissPreparer
	- Params
		- From: context.Context, string, string
		- To: context.Context, string, string, string
1. AlertsClient.UpdateResourceGroupLevelStateToResolve
	- Params
		- From: context.Context, string, string
		- To: context.Context, string, string, string
1. AlertsClient.UpdateResourceGroupLevelStateToResolvePreparer
	- Params
		- From: context.Context, string, string
		- To: context.Context, string, string, string
1. AlertsClient.UpdateSubscriptionLevelStateToActivate
	- Params
		- From: context.Context, string
		- To: context.Context, string, string
1. AlertsClient.UpdateSubscriptionLevelStateToActivatePreparer
	- Params
		- From: context.Context, string
		- To: context.Context, string, string
1. AlertsClient.UpdateSubscriptionLevelStateToDismiss
	- Params
		- From: context.Context, string
		- To: context.Context, string, string
1. AlertsClient.UpdateSubscriptionLevelStateToDismissPreparer
	- Params
		- From: context.Context, string
		- To: context.Context, string, string
1. AlertsClient.UpdateSubscriptionLevelStateToResolve
	- Params
		- From: context.Context, string
		- To: context.Context, string, string
1. AlertsClient.UpdateSubscriptionLevelStateToResolvePreparer
	- Params
		- From: context.Context, string
		- To: context.Context, string, string
1. AllowedConnectionsClient.Get
	- Params
		- From: context.Context, string, ConnectionType
		- To: context.Context, string, string, ConnectionType
1. AllowedConnectionsClient.GetPreparer
	- Params
		- From: context.Context, string, ConnectionType
		- To: context.Context, string, string, ConnectionType
1. AllowedConnectionsClient.ListByHomeRegion
	- Params
		- From: context.Context
		- To: context.Context, string
1. AllowedConnectionsClient.ListByHomeRegionComplete
	- Params
		- From: context.Context
		- To: context.Context, string
1. AllowedConnectionsClient.ListByHomeRegionPreparer
	- Params
		- From: context.Context
		- To: context.Context, string
1. AssessmentListIterator.Value
	- Returns
		- From: Assessment
		- To: AssessmentResponse
1. AssessmentListPage.Values
	- Returns
		- From: []Assessment
		- To: []AssessmentResponse
1. AssessmentsClient.CreateOrUpdate
	- Returns
		- From: Assessment, error
		- To: AssessmentResponse, error
1. AssessmentsClient.CreateOrUpdateResponder
	- Returns
		- From: Assessment, error
		- To: AssessmentResponse, error
1. AssessmentsClient.Get
	- Returns
		- From: Assessment, error
		- To: AssessmentResponse, error
1. AssessmentsClient.GetResponder
	- Returns
		- From: Assessment, error
		- To: AssessmentResponse, error
1. AssessmentsMetadataClient.CreateInSubscription
	- Params
		- From: context.Context, string, AssessmentMetadata
		- To: context.Context, string, AssessmentMetadataResponse
	- Returns
		- From: AssessmentMetadata, error
		- To: AssessmentMetadataResponse, error
1. AssessmentsMetadataClient.CreateInSubscriptionPreparer
	- Params
		- From: context.Context, string, AssessmentMetadata
		- To: context.Context, string, AssessmentMetadataResponse
1. AssessmentsMetadataClient.CreateInSubscriptionResponder
	- Returns
		- From: AssessmentMetadata, error
		- To: AssessmentMetadataResponse, error
1. AssessmentsMetadataClient.Get
	- Returns
		- From: AssessmentMetadata, error
		- To: AssessmentMetadataResponse, error
1. AssessmentsMetadataClient.GetInSubscription
	- Returns
		- From: AssessmentMetadata, error
		- To: AssessmentMetadataResponse, error
1. AssessmentsMetadataClient.GetInSubscriptionResponder
	- Returns
		- From: AssessmentMetadata, error
		- To: AssessmentMetadataResponse, error
1. AssessmentsMetadataClient.GetResponder
	- Returns
		- From: AssessmentMetadata, error
		- To: AssessmentMetadataResponse, error
1. AssessmentsMetadataClient.List
	- Returns
		- From: AssessmentMetadataListPage, error
		- To: AssessmentMetadataResponseListPage, error
1. AssessmentsMetadataClient.ListBySubscription
	- Returns
		- From: AssessmentMetadataListPage, error
		- To: AssessmentMetadataResponseListPage, error
1. AssessmentsMetadataClient.ListBySubscriptionComplete
	- Returns
		- From: AssessmentMetadataListIterator, error
		- To: AssessmentMetadataResponseListIterator, error
1. AssessmentsMetadataClient.ListBySubscriptionResponder
	- Returns
		- From: AssessmentMetadataList, error
		- To: AssessmentMetadataResponseList, error
1. AssessmentsMetadataClient.ListComplete
	- Returns
		- From: AssessmentMetadataListIterator, error
		- To: AssessmentMetadataResponseListIterator, error
1. AssessmentsMetadataClient.ListResponder
	- Returns
		- From: AssessmentMetadataList, error
		- To: AssessmentMetadataResponseList, error
1. DiscoveredSecuritySolutionsClient.Get
	- Params
		- From: context.Context, string, string
		- To: context.Context, string, string, string
1. DiscoveredSecuritySolutionsClient.GetPreparer
	- Params
		- From: context.Context, string, string
		- To: context.Context, string, string, string
1. DiscoveredSecuritySolutionsClient.ListByHomeRegion
	- Params
		- From: context.Context
		- To: context.Context, string
1. DiscoveredSecuritySolutionsClient.ListByHomeRegionComplete
	- Params
		- From: context.Context
		- To: context.Context, string
1. DiscoveredSecuritySolutionsClient.ListByHomeRegionPreparer
	- Params
		- From: context.Context
		- To: context.Context, string
1. ExternalSecuritySolutionsClient.Get
	- Params
		- From: context.Context, string, string
		- To: context.Context, string, string, string
1. ExternalSecuritySolutionsClient.GetPreparer
	- Params
		- From: context.Context, string, string
		- To: context.Context, string, string, string
1. ExternalSecuritySolutionsClient.ListByHomeRegion
	- Params
		- From: context.Context
		- To: context.Context, string
1. ExternalSecuritySolutionsClient.ListByHomeRegionComplete
	- Params
		- From: context.Context
		- To: context.Context, string
1. ExternalSecuritySolutionsClient.ListByHomeRegionPreparer
	- Params
		- From: context.Context
		- To: context.Context, string
1. InformationProtectionPoliciesClient.CreateOrUpdate
	- Params
		- From: context.Context, string, string, InformationProtectionPolicy
		- To: context.Context, string, InformationProtectionPolicyName, InformationProtectionPolicy
1. InformationProtectionPoliciesClient.CreateOrUpdatePreparer
	- Params
		- From: context.Context, string, string, InformationProtectionPolicy
		- To: context.Context, string, InformationProtectionPolicyName, InformationProtectionPolicy
1. InformationProtectionPoliciesClient.Get
	- Params
		- From: context.Context, string, string
		- To: context.Context, string, InformationProtectionPolicyName
1. InformationProtectionPoliciesClient.GetPreparer
	- Params
		- From: context.Context, string, string
		- To: context.Context, string, InformationProtectionPolicyName
1. JitNetworkAccessPoliciesClient.CreateOrUpdate
	- Params
		- From: context.Context, string, string, JitNetworkAccessPolicy
		- To: context.Context, string, string, string, JitNetworkAccessPolicy
1. JitNetworkAccessPoliciesClient.CreateOrUpdatePreparer
	- Params
		- From: context.Context, string, string, JitNetworkAccessPolicy
		- To: context.Context, string, string, string, JitNetworkAccessPolicy
1. JitNetworkAccessPoliciesClient.Delete
	- Params
		- From: context.Context, string, string
		- To: context.Context, string, string, string
1. JitNetworkAccessPoliciesClient.DeletePreparer
	- Params
		- From: context.Context, string, string
		- To: context.Context, string, string, string
1. JitNetworkAccessPoliciesClient.Get
	- Params
		- From: context.Context, string, string
		- To: context.Context, string, string, string
1. JitNetworkAccessPoliciesClient.GetPreparer
	- Params
		- From: context.Context, string, string
		- To: context.Context, string, string, string
1. JitNetworkAccessPoliciesClient.Initiate
	- Params
		- From: context.Context, string, string, JitNetworkAccessPolicyInitiateRequest
		- To: context.Context, string, string, string, JitNetworkAccessPolicyInitiateRequest
1. JitNetworkAccessPoliciesClient.InitiatePreparer
	- Params
		- From: context.Context, string, string, JitNetworkAccessPolicyInitiateRequest
		- To: context.Context, string, string, string, JitNetworkAccessPolicyInitiateRequest
1. JitNetworkAccessPoliciesClient.ListByRegion
	- Params
		- From: context.Context
		- To: context.Context, string
1. JitNetworkAccessPoliciesClient.ListByRegionComplete
	- Params
		- From: context.Context
		- To: context.Context, string
1. JitNetworkAccessPoliciesClient.ListByRegionPreparer
	- Params
		- From: context.Context
		- To: context.Context, string
1. JitNetworkAccessPoliciesClient.ListByResourceGroupAndRegion
	- Params
		- From: context.Context, string
		- To: context.Context, string, string
1. JitNetworkAccessPoliciesClient.ListByResourceGroupAndRegionComplete
	- Params
		- From: context.Context, string
		- To: context.Context, string, string
1. JitNetworkAccessPoliciesClient.ListByResourceGroupAndRegionPreparer
	- Params
		- From: context.Context, string
		- To: context.Context, string, string
1. LocationsClient.Get
	- Params
		- From: context.Context
		- To: context.Context, string
1. LocationsClient.GetPreparer
	- Params
		- From: context.Context
		- To: context.Context, string
1. New
	- Params
		- From: string, string
		- To: string
1. NewAdaptiveApplicationControlsClient
	- Params
		- From: string, string
		- To: string
1. NewAdaptiveApplicationControlsClientWithBaseURI
	- Params
		- From: string, string, string
		- To: string, string
1. NewAdaptiveNetworkHardeningsClient
	- Params
		- From: string, string
		- To: string
1. NewAdaptiveNetworkHardeningsClientWithBaseURI
	- Params
		- From: string, string, string
		- To: string, string
1. NewAdvancedThreatProtectionClient
	- Params
		- From: string, string
		- To: string
1. NewAdvancedThreatProtectionClientWithBaseURI
	- Params
		- From: string, string, string
		- To: string, string
1. NewAlertsClient
	- Params
		- From: string, string
		- To: string
1. NewAlertsClientWithBaseURI
	- Params
		- From: string, string, string
		- To: string, string
1. NewAlertsSuppressionRulesClient
	- Params
		- From: string, string
		- To: string
1. NewAlertsSuppressionRulesClientWithBaseURI
	- Params
		- From: string, string, string
		- To: string, string
1. NewAllowedConnectionsClient
	- Params
		- From: string, string
		- To: string
1. NewAllowedConnectionsClientWithBaseURI
	- Params
		- From: string, string, string
		- To: string, string
1. NewAssessmentsClient
	- Params
		- From: string, string
		- To: string
1. NewAssessmentsClientWithBaseURI
	- Params
		- From: string, string, string
		- To: string, string
1. NewAssessmentsMetadataClient
	- Params
		- From: string, string
		- To: string
1. NewAssessmentsMetadataClientWithBaseURI
	- Params
		- From: string, string, string
		- To: string, string
1. NewAutoProvisioningSettingsClient
	- Params
		- From: string, string
		- To: string
1. NewAutoProvisioningSettingsClientWithBaseURI
	- Params
		- From: string, string, string
		- To: string, string
1. NewAutomationsClient
	- Params
		- From: string, string
		- To: string
1. NewAutomationsClientWithBaseURI
	- Params
		- From: string, string, string
		- To: string, string
1. NewComplianceResultsClient
	- Params
		- From: string, string
		- To: string
1. NewComplianceResultsClientWithBaseURI
	- Params
		- From: string, string, string
		- To: string, string
1. NewCompliancesClient
	- Params
		- From: string, string
		- To: string
1. NewCompliancesClientWithBaseURI
	- Params
		- From: string, string, string
		- To: string, string
1. NewConnectorsClient
	- Params
		- From: string, string
		- To: string
1. NewConnectorsClientWithBaseURI
	- Params
		- From: string, string, string
		- To: string, string
1. NewContactsClient
	- Params
		- From: string, string
		- To: string
1. NewContactsClientWithBaseURI
	- Params
		- From: string, string, string
		- To: string, string
1. NewDeviceSecurityGroupsClient
	- Params
		- From: string, string
		- To: string
1. NewDeviceSecurityGroupsClientWithBaseURI
	- Params
		- From: string, string, string
		- To: string, string
1. NewDiscoveredSecuritySolutionsClient
	- Params
		- From: string, string
		- To: string
1. NewDiscoveredSecuritySolutionsClientWithBaseURI
	- Params
		- From: string, string, string
		- To: string, string
1. NewExternalSecuritySolutionsClient
	- Params
		- From: string, string
		- To: string
1. NewExternalSecuritySolutionsClientWithBaseURI
	- Params
		- From: string, string, string
		- To: string, string
1. NewInformationProtectionPoliciesClient
	- Params
		- From: string, string
		- To: string
1. NewInformationProtectionPoliciesClientWithBaseURI
	- Params
		- From: string, string, string
		- To: string, string
1. NewIngestionSettingsClient
	- Params
		- From: string, string
		- To: string
1. NewIngestionSettingsClientWithBaseURI
	- Params
		- From: string, string, string
		- To: string, string
1. NewIotSecuritySolutionAnalyticsClient
	- Params
		- From: string, string
		- To: string
1. NewIotSecuritySolutionAnalyticsClientWithBaseURI
	- Params
		- From: string, string, string
		- To: string, string
1. NewIotSecuritySolutionClient
	- Params
		- From: string, string
		- To: string
1. NewIotSecuritySolutionClientWithBaseURI
	- Params
		- From: string, string, string
		- To: string, string
1. NewIotSecuritySolutionsAnalyticsAggregatedAlertClient
	- Params
		- From: string, string
		- To: string
1. NewIotSecuritySolutionsAnalyticsAggregatedAlertClientWithBaseURI
	- Params
		- From: string, string, string
		- To: string, string
1. NewIotSecuritySolutionsAnalyticsRecommendationClient
	- Params
		- From: string, string
		- To: string
1. NewIotSecuritySolutionsAnalyticsRecommendationClientWithBaseURI
	- Params
		- From: string, string, string
		- To: string, string
1. NewJitNetworkAccessPoliciesClient
	- Params
		- From: string, string
		- To: string
1. NewJitNetworkAccessPoliciesClientWithBaseURI
	- Params
		- From: string, string, string
		- To: string, string
1. NewLocationsClient
	- Params
		- From: string, string
		- To: string
1. NewLocationsClientWithBaseURI
	- Params
		- From: string, string, string
		- To: string, string
1. NewOperationsClient
	- Params
		- From: string, string
		- To: string
1. NewOperationsClientWithBaseURI
	- Params
		- From: string, string, string
		- To: string, string
1. NewPricingsClient
	- Params
		- From: string, string
		- To: string
1. NewPricingsClientWithBaseURI
	- Params
		- From: string, string, string
		- To: string, string
1. NewRegulatoryComplianceAssessmentsClient
	- Params
		- From: string, string
		- To: string
1. NewRegulatoryComplianceAssessmentsClientWithBaseURI
	- Params
		- From: string, string, string
		- To: string, string
1. NewRegulatoryComplianceControlsClient
	- Params
		- From: string, string
		- To: string
1. NewRegulatoryComplianceControlsClientWithBaseURI
	- Params
		- From: string, string, string
		- To: string, string
1. NewRegulatoryComplianceStandardsClient
	- Params
		- From: string, string
		- To: string
1. NewRegulatoryComplianceStandardsClientWithBaseURI
	- Params
		- From: string, string, string
		- To: string, string
1. NewSQLVulnerabilityAssessmentBaselineRulesClient
	- Params
		- From: string, string
		- To: string
1. NewSQLVulnerabilityAssessmentBaselineRulesClientWithBaseURI
	- Params
		- From: string, string, string
		- To: string, string
1. NewSQLVulnerabilityAssessmentScanResultsClient
	- Params
		- From: string, string
		- To: string
1. NewSQLVulnerabilityAssessmentScanResultsClientWithBaseURI
	- Params
		- From: string, string, string
		- To: string, string
1. NewSQLVulnerabilityAssessmentScansClient
	- Params
		- From: string, string
		- To: string
1. NewSQLVulnerabilityAssessmentScansClientWithBaseURI
	- Params
		- From: string, string, string
		- To: string, string
1. NewSecureScoreControlDefinitionsClient
	- Params
		- From: string, string
		- To: string
1. NewSecureScoreControlDefinitionsClientWithBaseURI
	- Params
		- From: string, string, string
		- To: string, string
1. NewSecureScoreControlsClient
	- Params
		- From: string, string
		- To: string
1. NewSecureScoreControlsClientWithBaseURI
	- Params
		- From: string, string, string
		- To: string, string
1. NewSecureScoresClient
	- Params
		- From: string, string
		- To: string
1. NewSecureScoresClientWithBaseURI
	- Params
		- From: string, string, string
		- To: string, string
1. NewServerVulnerabilityAssessmentClient
	- Params
		- From: string, string
		- To: string
1. NewServerVulnerabilityAssessmentClientWithBaseURI
	- Params
		- From: string, string, string
		- To: string, string
1. NewSettingsClient
	- Params
		- From: string, string
		- To: string
1. NewSettingsClientWithBaseURI
	- Params
		- From: string, string, string
		- To: string, string
1. NewSoftwareInventoriesClient
	- Params
		- From: string, string
		- To: string
1. NewSoftwareInventoriesClientWithBaseURI
	- Params
		- From: string, string, string
		- To: string, string
1. NewSolutionsClient
	- Params
		- From: string, string
		- To: string
1. NewSolutionsClientWithBaseURI
	- Params
		- From: string, string, string
		- To: string, string
1. NewSolutionsReferenceDataClient
	- Params
		- From: string, string
		- To: string
1. NewSolutionsReferenceDataClientWithBaseURI
	- Params
		- From: string, string, string
		- To: string, string
1. NewSubAssessmentsClient
	- Params
		- From: string, string
		- To: string
1. NewSubAssessmentsClientWithBaseURI
	- Params
		- From: string, string, string
		- To: string, string
1. NewTasksClient
	- Params
		- From: string, string
		- To: string
1. NewTasksClientWithBaseURI
	- Params
		- From: string, string, string
		- To: string, string
1. NewTopologyClient
	- Params
		- From: string, string
		- To: string
1. NewTopologyClientWithBaseURI
	- Params
		- From: string, string, string
		- To: string, string
1. NewWithBaseURI
	- Params
		- From: string, string, string
		- To: string, string
1. NewWorkspaceSettingsClient
	- Params
		- From: string, string
		- To: string
1. NewWorkspaceSettingsClientWithBaseURI
	- Params
		- From: string, string, string
		- To: string, string
1. SettingsClient.Get
	- Params
		- From: context.Context, string
		- To: context.Context, SettingName4
1. SettingsClient.GetPreparer
	- Params
		- From: context.Context, string
		- To: context.Context, SettingName4
1. SettingsClient.Update
	- Params
		- From: context.Context, string, BasicSetting
		- To: context.Context, SettingName5, BasicSetting
1. SettingsClient.UpdatePreparer
	- Params
		- From: context.Context, string, BasicSetting
		- To: context.Context, SettingName5, BasicSetting
1. SolutionsClient.Get
	- Params
		- From: context.Context, string, string
		- To: context.Context, string, string, string
1. SolutionsClient.GetPreparer
	- Params
		- From: context.Context, string, string
		- To: context.Context, string, string, string
1. SolutionsReferenceDataClient.ListByHomeRegion
	- Params
		- From: context.Context
		- To: context.Context, string
1. SolutionsReferenceDataClient.ListByHomeRegionPreparer
	- Params
		- From: context.Context
		- To: context.Context, string
1. TasksClient.GetResourceGroupLevelTask
	- Params
		- From: context.Context, string, string
		- To: context.Context, string, string, string
1. TasksClient.GetResourceGroupLevelTaskPreparer
	- Params
		- From: context.Context, string, string
		- To: context.Context, string, string, string
1. TasksClient.GetSubscriptionLevelTask
	- Params
		- From: context.Context, string
		- To: context.Context, string, string
1. TasksClient.GetSubscriptionLevelTaskPreparer
	- Params
		- From: context.Context, string
		- To: context.Context, string, string
1. TasksClient.ListByHomeRegion
	- Params
		- From: context.Context, string
		- To: context.Context, string, string
1. TasksClient.ListByHomeRegionComplete
	- Params
		- From: context.Context, string
		- To: context.Context, string, string
1. TasksClient.ListByHomeRegionPreparer
	- Params
		- From: context.Context, string
		- To: context.Context, string, string
1. TasksClient.ListByResourceGroup
	- Params
		- From: context.Context, string, string
		- To: context.Context, string, string, string
1. TasksClient.ListByResourceGroupComplete
	- Params
		- From: context.Context, string, string
		- To: context.Context, string, string, string
1. TasksClient.ListByResourceGroupPreparer
	- Params
		- From: context.Context, string, string
		- To: context.Context, string, string, string
1. TasksClient.UpdateResourceGroupLevelTaskState
	- Params
		- From: context.Context, string, string, string
		- To: context.Context, string, string, string, TaskUpdateActionType
1. TasksClient.UpdateResourceGroupLevelTaskStatePreparer
	- Params
		- From: context.Context, string, string, string
		- To: context.Context, string, string, string, TaskUpdateActionType
1. TasksClient.UpdateSubscriptionLevelTaskState
	- Params
		- From: context.Context, string, string
		- To: context.Context, string, string, TaskUpdateActionType
1. TasksClient.UpdateSubscriptionLevelTaskStatePreparer
	- Params
		- From: context.Context, string, string
		- To: context.Context, string, string, TaskUpdateActionType
1. TopologyClient.Get
	- Params
		- From: context.Context, string, string
		- To: context.Context, string, string, string
1. TopologyClient.GetPreparer
	- Params
		- From: context.Context, string, string
		- To: context.Context, string, string, string
1. TopologyClient.ListByHomeRegion
	- Params
		- From: context.Context
		- To: context.Context, string
1. TopologyClient.ListByHomeRegionComplete
	- Params
		- From: context.Context
		- To: context.Context, string
1. TopologyClient.ListByHomeRegionPreparer
	- Params
		- From: context.Context
		- To: context.Context, string

#### Struct Fields

1. AssessmentList.Value changed type from *[]Assessment to *[]AssessmentResponse
1. ContactProperties.AlertNotifications changed type from AlertNotifications to *ContactPropertiesAlertNotifications
1. PathRecommendation.Type changed type from Type to Type1

## Additive Changes

### New Constants

1. ApplicationConditionOperator.In
1. BundleType.CosmosDbs
1. CloudName.AWS
1. CloudName.Azure
1. CloudName.AzureDevOps
1. CloudName.GCP
1. CloudName.Github
1. ConfigurationStatus.ConfigurationStatusConfigured
1. ConfigurationStatus.ConfigurationStatusFailed
1. ConfigurationStatus.ConfigurationStatusInProgress
1. ConfigurationStatus.ConfigurationStatusNoStatus
1. ConfigurationStatus.ConfigurationStatusNotConfigured
1. CreatedByType.CreatedByTypeApplication
1. CreatedByType.CreatedByTypeKey
1. CreatedByType.CreatedByTypeManagedIdentity
1. CreatedByType.CreatedByTypeUser
1. EnforcementSupport.NotSupported
1. EnforcementSupport.Supported
1. EnvironmentType.EnvironmentTypeAwsAccount
1. EnvironmentType.EnvironmentTypeAzureDevOpsScope
1. EnvironmentType.EnvironmentTypeEnvironmentData
1. EnvironmentType.EnvironmentTypeGcpProject
1. EnvironmentType.EnvironmentTypeGithubScope
1. EventSource.EventSourceAssessmentsSnapshot
1. EventSource.EventSourceSubAssessmentsSnapshot
1. GovernanceRuleConditionOperator.GovernanceRuleConditionOperatorEquals
1. GovernanceRuleConditionOperator.GovernanceRuleConditionOperatorIn
1. GovernanceRuleOwnerSourceType.ByTag
1. GovernanceRuleOwnerSourceType.Manually
1. GovernanceRuleType.Integrated
1. GovernanceRuleType.ServiceNow
1. InformationProtectionPolicyName.Custom
1. InformationProtectionPolicyName.Effective
1. MinimalSeverity.MinimalSeverityHigh
1. MinimalSeverity.MinimalSeverityLow
1. MinimalSeverity.MinimalSeverityMedium
1. OfferingType.OfferingTypeCloudOffering
1. OfferingType.OfferingTypeCspmMonitorAws
1. OfferingType.OfferingTypeCspmMonitorAzureDevOps
1. OfferingType.OfferingTypeCspmMonitorGcp
1. OfferingType.OfferingTypeCspmMonitorGithub
1. OfferingType.OfferingTypeDefenderCspmAws
1. OfferingType.OfferingTypeDefenderCspmGcp
1. OfferingType.OfferingTypeDefenderForContainersAws
1. OfferingType.OfferingTypeDefenderForContainersGcp
1. OfferingType.OfferingTypeDefenderForDatabasesAws
1. OfferingType.OfferingTypeDefenderForDatabasesGcp
1. OfferingType.OfferingTypeDefenderForDevOpsAzureDevOps
1. OfferingType.OfferingTypeDefenderForDevOpsGithub
1. OfferingType.OfferingTypeDefenderForServersAws
1. OfferingType.OfferingTypeDefenderForServersGcp
1. OfferingType.OfferingTypeInformationProtectionAws
1. Operator.OperatorContains
1. Operator.OperatorEndsWith
1. Operator.OperatorEquals
1. Operator.OperatorGreaterThan
1. Operator.OperatorGreaterThanOrEqualTo
1. Operator.OperatorLesserThan
1. Operator.OperatorLesserThanOrEqualTo
1. Operator.OperatorNotEquals
1. Operator.OperatorStartsWith
1. OrganizationMembershipType.OrganizationMembershipTypeAwsOrganizationalData
1. OrganizationMembershipType.OrganizationMembershipTypeMember
1. OrganizationMembershipType.OrganizationMembershipTypeOrganization
1. OrganizationMembershipTypeBasicGcpOrganizationalData.OrganizationMembershipTypeBasicGcpOrganizationalDataOrganizationMembershipTypeGcpOrganizationalData
1. OrganizationMembershipTypeBasicGcpOrganizationalData.OrganizationMembershipTypeBasicGcpOrganizationalDataOrganizationMembershipTypeMember
1. OrganizationMembershipTypeBasicGcpOrganizationalData.OrganizationMembershipTypeBasicGcpOrganizationalDataOrganizationMembershipTypeOrganization
1. PricingTier.Free
1. ProvisioningState.Succeeded
1. ProvisioningState.Updating
1. Roles.AccountAdmin
1. Roles.Contributor
1. Roles.Owner
1. Roles.ServiceAdmin
1. SettingName2.SettingName2MCAS
1. SettingName2.SettingName2Sentinel
1. SettingName2.SettingName2WDATP
1. SettingName2.SettingName2WDATPEXCLUDELINUXPUBLICPREVIEW
1. SettingName2.SettingName2WDATPUNIFIEDSOLUTION
1. SettingName4.SettingName4MCAS
1. SettingName4.SettingName4Sentinel
1. SettingName4.SettingName4WDATP
1. SettingName4.SettingName4WDATPEXCLUDELINUXPUBLICPREVIEW
1. SettingName4.SettingName4WDATPUNIFIEDSOLUTION
1. SettingName5.SettingName5MCAS
1. SettingName5.SettingName5Sentinel
1. SettingName5.SettingName5WDATP
1. SettingName5.SettingName5WDATPEXCLUDELINUXPUBLICPREVIEW
1. SettingName5.SettingName5WDATPUNIFIEDSOLUTION
1. SeverityEnum.SeverityEnumHigh
1. SeverityEnum.SeverityEnumLow
1. SeverityEnum.SeverityEnumMedium
1. StateForAlertNotifications.StateForAlertNotificationsOff
1. StateForAlertNotifications.StateForAlertNotificationsOn
1. StateForNotificationsByRole.StateForNotificationsByRoleOff
1. StateForNotificationsByRole.StateForNotificationsByRoleOn
1. SubPlan.P1
1. SubPlan.P2
1. SupportedCloudEnum.SupportedCloudEnumAWS
1. SupportedCloudEnum.SupportedCloudEnumGCP
1. Tactics.CommandandControl
1. Tactics.Reconnaissance
1. Tactics.ResourceDevelopment
1. TaskUpdateActionType.Activate
1. TaskUpdateActionType.Close
1. TaskUpdateActionType.Dismiss
1. TaskUpdateActionType.Resolve
1. TaskUpdateActionType.Start
1. Techniques.AbuseElevationControlMechanism
1. Techniques.AccessTokenManipulation
1. Techniques.AccountDiscovery
1. Techniques.AccountManipulation
1. Techniques.ActiveScanning
1. Techniques.ApplicationLayerProtocol
1. Techniques.AudioCapture
1. Techniques.BootorLogonAutostartExecution
1. Techniques.BootorLogonInitializationScripts
1. Techniques.BruteForce
1. Techniques.CloudInfrastructureDiscovery
1. Techniques.CloudServiceDashboard
1. Techniques.CloudServiceDiscovery
1. Techniques.CommandandScriptingInterpreter
1. Techniques.CompromiseClientSoftwareBinary
1. Techniques.CompromiseInfrastructure
1. Techniques.ContainerandResourceDiscovery
1. Techniques.CreateAccount
1. Techniques.CreateorModifySystemProcess
1. Techniques.CredentialsfromPasswordStores
1. Techniques.DataDestruction
1. Techniques.DataEncryptedforImpact
1. Techniques.DataManipulation
1. Techniques.DataStaged
1. Techniques.DatafromCloudStorageObject
1. Techniques.DatafromConfigurationRepository
1. Techniques.DatafromInformationRepositories
1. Techniques.DatafromLocalSystem
1. Techniques.Defacement
1. Techniques.DeobfuscateDecodeFilesorInformation
1. Techniques.DiskWipe
1. Techniques.DomainTrustDiscovery
1. Techniques.DriveByCompromise
1. Techniques.DynamicResolution
1. Techniques.EndpointDenialofService
1. Techniques.EventTriggeredExecution
1. Techniques.ExfiltrationOverAlternativeProtocol
1. Techniques.ExploitPublicFacingApplication
1. Techniques.ExploitationforClientExecution
1. Techniques.ExploitationforCredentialAccess
1. Techniques.ExploitationforDefenseEvasion
1. Techniques.ExploitationforPrivilegeEscalation
1. Techniques.ExploitationofRemoteServices
1. Techniques.ExternalRemoteServices
1. Techniques.FallbackChannels
1. Techniques.FileandDirectoryDiscovery
1. Techniques.FileandDirectoryPermissionsModification
1. Techniques.GatherVictimNetworkInformation
1. Techniques.HideArtifacts
1. Techniques.HijackExecutionFlow
1. Techniques.ImpairDefenses
1. Techniques.ImplantContainerImage
1. Techniques.IndicatorRemovalonHost
1. Techniques.IndirectCommandExecution
1. Techniques.IngressToolTransfer
1. Techniques.InputCapture
1. Techniques.InterProcessCommunication
1. Techniques.LateralToolTransfer
1. Techniques.ManInTheMiddle
1. Techniques.Masquerading
1. Techniques.ModifyAuthenticationProcess
1. Techniques.ModifyRegistry
1. Techniques.NetworkDenialofService
1. Techniques.NetworkServiceScanning
1. Techniques.NetworkSniffing
1. Techniques.NonApplicationLayerProtocol
1. Techniques.NonStandardPort
1. Techniques.OSCredentialDumping
1. Techniques.ObfuscatedFilesorInformation
1. Techniques.ObtainCapabilities
1. Techniques.OfficeApplicationStartup
1. Techniques.PermissionGroupsDiscovery
1. Techniques.Phishing
1. Techniques.PreOSBoot
1. Techniques.ProcessDiscovery
1. Techniques.ProcessInjection
1. Techniques.ProtocolTunneling
1. Techniques.Proxy
1. Techniques.QueryRegistry
1. Techniques.RemoteAccessSoftware
1. Techniques.RemoteServiceSessionHijacking
1. Techniques.RemoteServices
1. Techniques.RemoteSystemDiscovery
1. Techniques.ResourceHijacking
1. Techniques.SQLStoredProcedures
1. Techniques.ScheduledTaskJob
1. Techniques.ScreenCapture
1. Techniques.SearchVictimOwnedWebsites
1. Techniques.ServerSoftwareComponent
1. Techniques.ServiceStop
1. Techniques.SignedBinaryProxyExecution
1. Techniques.SoftwareDeploymentTools
1. Techniques.StealorForgeKerberosTickets
1. Techniques.SubvertTrustControls
1. Techniques.SupplyChainCompromise
1. Techniques.SystemInformationDiscovery
1. Techniques.TaintSharedContent
1. Techniques.TrafficSignaling
1. Techniques.TransferDatatoCloudAccount
1. Techniques.TrustedRelationship
1. Techniques.UnsecuredCredentials
1. Techniques.UserExecution
1. Techniques.ValidAccounts
1. Techniques.WindowsManagementInstrumentation
1. Type.Qualys
1. Type.TVM

### New Type Aliases

1. string.ApplicationConditionOperator
1. string.CloudName
1. string.EnvironmentType
1. string.GovernanceRuleConditionOperator
1. string.GovernanceRuleOwnerSourceType
1. string.GovernanceRuleType
1. string.InformationProtectionPolicyName
1. string.MinimalSeverity
1. string.OfferingType
1. string.OrganizationMembershipType
1. string.OrganizationMembershipTypeBasicGcpOrganizationalData
1. string.Roles
1. string.ScanningMode
1. string.SettingName2
1. string.SettingName4
1. string.SettingName5
1. string.SeverityEnum
1. string.StateForAlertNotifications
1. string.StateForNotificationsByRole
1. string.SubPlan
1. string.SupportedCloudEnum
1. string.Tactics
1. string.TaskUpdateActionType
1. string.Techniques
1. string.Type1

### New Funcs

1. *AlertPropertiesSupportingEvidence.UnmarshalJSON([]byte) error
1. *Application.UnmarshalJSON([]byte) error
1. *ApplicationsListIterator.Next() error
1. *ApplicationsListIterator.NextWithContext(context.Context) error
1. *ApplicationsListPage.Next() error
1. *ApplicationsListPage.NextWithContext(context.Context) error
1. *AssessmentMetadataResponse.UnmarshalJSON([]byte) error
1. *AssessmentMetadataResponseListIterator.Next() error
1. *AssessmentMetadataResponseListIterator.NextWithContext(context.Context) error
1. *AssessmentMetadataResponseListPage.Next() error
1. *AssessmentMetadataResponseListPage.NextWithContext(context.Context) error
1. *AssessmentPropertiesBase.UnmarshalJSON([]byte) error
1. *AssessmentPropertiesResponse.UnmarshalJSON([]byte) error
1. *AssessmentResponse.UnmarshalJSON([]byte) error
1. *AwsEnvironmentData.UnmarshalJSON([]byte) error
1. *Connector.UnmarshalJSON([]byte) error
1. *ConnectorGovernanceRulesExecuteStatusGetFuture.UnmarshalJSON([]byte) error
1. *ConnectorProperties.UnmarshalJSON([]byte) error
1. *ConnectorsListIterator.Next() error
1. *ConnectorsListIterator.NextWithContext(context.Context) error
1. *ConnectorsListPage.Next() error
1. *ConnectorsListPage.NextWithContext(context.Context) error
1. *CustomAssessmentAutomation.UnmarshalJSON([]byte) error
1. *CustomAssessmentAutomationRequest.UnmarshalJSON([]byte) error
1. *CustomAssessmentAutomationsListResultIterator.Next() error
1. *CustomAssessmentAutomationsListResultIterator.NextWithContext(context.Context) error
1. *CustomAssessmentAutomationsListResultPage.Next() error
1. *CustomAssessmentAutomationsListResultPage.NextWithContext(context.Context) error
1. *CustomEntityStoreAssignment.UnmarshalJSON([]byte) error
1. *CustomEntityStoreAssignmentRequest.UnmarshalJSON([]byte) error
1. *CustomEntityStoreAssignmentsListResultIterator.Next() error
1. *CustomEntityStoreAssignmentsListResultIterator.NextWithContext(context.Context) error
1. *CustomEntityStoreAssignmentsListResultPage.Next() error
1. *CustomEntityStoreAssignmentsListResultPage.NextWithContext(context.Context) error
1. *GcpProjectEnvironmentData.UnmarshalJSON([]byte) error
1. *GovernanceAssignment.UnmarshalJSON([]byte) error
1. *GovernanceAssignmentsListIterator.Next() error
1. *GovernanceAssignmentsListIterator.NextWithContext(context.Context) error
1. *GovernanceAssignmentsListPage.Next() error
1. *GovernanceAssignmentsListPage.NextWithContext(context.Context) error
1. *GovernanceRule.UnmarshalJSON([]byte) error
1. *GovernanceRuleListIterator.Next() error
1. *GovernanceRuleListIterator.NextWithContext(context.Context) error
1. *GovernanceRuleListPage.Next() error
1. *GovernanceRuleListPage.NextWithContext(context.Context) error
1. *GovernanceRulesRuleIDExecuteSingleSecurityConnectorFuture.UnmarshalJSON([]byte) error
1. *GovernanceRulesRuleIDExecuteSingleSubscriptionFuture.UnmarshalJSON([]byte) error
1. *MdeOnboardingData.UnmarshalJSON([]byte) error
1. *SubscriptionGovernanceRulesExecuteStatusGetFuture.UnmarshalJSON([]byte) error
1. AlertPropertiesSupportingEvidence.MarshalJSON() ([]byte, error)
1. AlertsClient.UpdateResourceGroupLevelStateToInProgress(context.Context, string, string, string) (autorest.Response, error)
1. AlertsClient.UpdateResourceGroupLevelStateToInProgressPreparer(context.Context, string, string, string) (*http.Request, error)
1. AlertsClient.UpdateResourceGroupLevelStateToInProgressResponder(*http.Response) (autorest.Response, error)
1. AlertsClient.UpdateResourceGroupLevelStateToInProgressSender(*http.Request) (*http.Response, error)
1. AlertsClient.UpdateSubscriptionLevelStateToInProgress(context.Context, string, string) (autorest.Response, error)
1. AlertsClient.UpdateSubscriptionLevelStateToInProgressPreparer(context.Context, string, string) (*http.Request, error)
1. AlertsClient.UpdateSubscriptionLevelStateToInProgressResponder(*http.Response) (autorest.Response, error)
1. AlertsClient.UpdateSubscriptionLevelStateToInProgressSender(*http.Request) (*http.Response, error)
1. Application.MarshalJSON() ([]byte, error)
1. ApplicationClient.CreateOrUpdate(context.Context, string, Application) (Application, error)
1. ApplicationClient.CreateOrUpdatePreparer(context.Context, string, Application) (*http.Request, error)
1. ApplicationClient.CreateOrUpdateResponder(*http.Response) (Application, error)
1. ApplicationClient.CreateOrUpdateSender(*http.Request) (*http.Response, error)
1. ApplicationClient.Delete(context.Context, string) (autorest.Response, error)
1. ApplicationClient.DeletePreparer(context.Context, string) (*http.Request, error)
1. ApplicationClient.DeleteResponder(*http.Response) (autorest.Response, error)
1. ApplicationClient.DeleteSender(*http.Request) (*http.Response, error)
1. ApplicationClient.Get(context.Context, string) (Application, error)
1. ApplicationClient.GetPreparer(context.Context, string) (*http.Request, error)
1. ApplicationClient.GetResponder(*http.Response) (Application, error)
1. ApplicationClient.GetSender(*http.Request) (*http.Response, error)
1. ApplicationsClient.List(context.Context) (ApplicationsListPage, error)
1. ApplicationsClient.ListComplete(context.Context) (ApplicationsListIterator, error)
1. ApplicationsClient.ListPreparer(context.Context) (*http.Request, error)
1. ApplicationsClient.ListResponder(*http.Response) (ApplicationsList, error)
1. ApplicationsClient.ListSender(*http.Request) (*http.Response, error)
1. ApplicationsList.IsEmpty() bool
1. ApplicationsList.MarshalJSON() ([]byte, error)
1. ApplicationsListIterator.NotDone() bool
1. ApplicationsListIterator.Response() ApplicationsList
1. ApplicationsListIterator.Value() Application
1. ApplicationsListPage.NotDone() bool
1. ApplicationsListPage.Response() ApplicationsList
1. ApplicationsListPage.Values() []Application
1. AssessmentMetadataPropertiesResponse.MarshalJSON() ([]byte, error)
1. AssessmentMetadataResponse.MarshalJSON() ([]byte, error)
1. AssessmentMetadataResponseList.IsEmpty() bool
1. AssessmentMetadataResponseList.MarshalJSON() ([]byte, error)
1. AssessmentMetadataResponseListIterator.NotDone() bool
1. AssessmentMetadataResponseListIterator.Response() AssessmentMetadataResponseList
1. AssessmentMetadataResponseListIterator.Value() AssessmentMetadataResponse
1. AssessmentMetadataResponseListPage.NotDone() bool
1. AssessmentMetadataResponseListPage.Response() AssessmentMetadataResponseList
1. AssessmentMetadataResponseListPage.Values() []AssessmentMetadataResponse
1. AssessmentPropertiesBase.MarshalJSON() ([]byte, error)
1. AssessmentPropertiesResponse.MarshalJSON() ([]byte, error)
1. AssessmentResponse.MarshalJSON() ([]byte, error)
1. AssessmentStatusResponse.MarshalJSON() ([]byte, error)
1. AwsEnvironmentData.AsAwsEnvironmentData() (*AwsEnvironmentData, bool)
1. AwsEnvironmentData.AsAzureDevOpsScopeEnvironmentData() (*AzureDevOpsScopeEnvironmentData, bool)
1. AwsEnvironmentData.AsBasicEnvironmentData() (BasicEnvironmentData, bool)
1. AwsEnvironmentData.AsEnvironmentData() (*EnvironmentData, bool)
1. AwsEnvironmentData.AsGcpProjectEnvironmentData() (*GcpProjectEnvironmentData, bool)
1. AwsEnvironmentData.AsGithubScopeEnvironmentData() (*GithubScopeEnvironmentData, bool)
1. AwsEnvironmentData.MarshalJSON() ([]byte, error)
1. AwsOrganizationalData.AsAwsOrganizationalData() (*AwsOrganizationalData, bool)
1. AwsOrganizationalData.AsAwsOrganizationalDataMaster() (*AwsOrganizationalDataMaster, bool)
1. AwsOrganizationalData.AsAwsOrganizationalDataMember() (*AwsOrganizationalDataMember, bool)
1. AwsOrganizationalData.AsBasicAwsOrganizationalData() (BasicAwsOrganizationalData, bool)
1. AwsOrganizationalData.MarshalJSON() ([]byte, error)
1. AwsOrganizationalDataMaster.AsAwsOrganizationalData() (*AwsOrganizationalData, bool)
1. AwsOrganizationalDataMaster.AsAwsOrganizationalDataMaster() (*AwsOrganizationalDataMaster, bool)
1. AwsOrganizationalDataMaster.AsAwsOrganizationalDataMember() (*AwsOrganizationalDataMember, bool)
1. AwsOrganizationalDataMaster.AsBasicAwsOrganizationalData() (BasicAwsOrganizationalData, bool)
1. AwsOrganizationalDataMaster.MarshalJSON() ([]byte, error)
1. AwsOrganizationalDataMember.AsAwsOrganizationalData() (*AwsOrganizationalData, bool)
1. AwsOrganizationalDataMember.AsAwsOrganizationalDataMaster() (*AwsOrganizationalDataMaster, bool)
1. AwsOrganizationalDataMember.AsAwsOrganizationalDataMember() (*AwsOrganizationalDataMember, bool)
1. AwsOrganizationalDataMember.AsBasicAwsOrganizationalData() (BasicAwsOrganizationalData, bool)
1. AwsOrganizationalDataMember.MarshalJSON() ([]byte, error)
1. AzureDevOpsScopeEnvironmentData.AsAwsEnvironmentData() (*AwsEnvironmentData, bool)
1. AzureDevOpsScopeEnvironmentData.AsAzureDevOpsScopeEnvironmentData() (*AzureDevOpsScopeEnvironmentData, bool)
1. AzureDevOpsScopeEnvironmentData.AsBasicEnvironmentData() (BasicEnvironmentData, bool)
1. AzureDevOpsScopeEnvironmentData.AsEnvironmentData() (*EnvironmentData, bool)
1. AzureDevOpsScopeEnvironmentData.AsGcpProjectEnvironmentData() (*GcpProjectEnvironmentData, bool)
1. AzureDevOpsScopeEnvironmentData.AsGithubScopeEnvironmentData() (*GithubScopeEnvironmentData, bool)
1. AzureDevOpsScopeEnvironmentData.MarshalJSON() ([]byte, error)
1. CloudOffering.AsBasicCloudOffering() (BasicCloudOffering, bool)
1. CloudOffering.AsCloudOffering() (*CloudOffering, bool)
1. CloudOffering.AsCspmMonitorAwsOffering() (*CspmMonitorAwsOffering, bool)
1. CloudOffering.AsCspmMonitorAzureDevOpsOffering() (*CspmMonitorAzureDevOpsOffering, bool)
1. CloudOffering.AsCspmMonitorGcpOffering() (*CspmMonitorGcpOffering, bool)
1. CloudOffering.AsCspmMonitorGithubOffering() (*CspmMonitorGithubOffering, bool)
1. CloudOffering.AsDefenderCspmAwsOffering() (*DefenderCspmAwsOffering, bool)
1. CloudOffering.AsDefenderCspmGcpOffering() (*DefenderCspmGcpOffering, bool)
1. CloudOffering.AsDefenderFoDatabasesAwsOffering() (*DefenderFoDatabasesAwsOffering, bool)
1. CloudOffering.AsDefenderForContainersAwsOffering() (*DefenderForContainersAwsOffering, bool)
1. CloudOffering.AsDefenderForContainersGcpOffering() (*DefenderForContainersGcpOffering, bool)
1. CloudOffering.AsDefenderForDatabasesGcpOffering() (*DefenderForDatabasesGcpOffering, bool)
1. CloudOffering.AsDefenderForDevOpsAzureDevOpsOffering() (*DefenderForDevOpsAzureDevOpsOffering, bool)
1. CloudOffering.AsDefenderForDevOpsGithubOffering() (*DefenderForDevOpsGithubOffering, bool)
1. CloudOffering.AsDefenderForServersAwsOffering() (*DefenderForServersAwsOffering, bool)
1. CloudOffering.AsDefenderForServersGcpOffering() (*DefenderForServersGcpOffering, bool)
1. CloudOffering.AsInformationProtectionAwsOffering() (*InformationProtectionAwsOffering, bool)
1. CloudOffering.MarshalJSON() ([]byte, error)
1. Connector.MarshalJSON() ([]byte, error)
1. ConnectorApplicationClient.CreateOrUpdate(context.Context, string, string, string, Application) (Application, error)
1. ConnectorApplicationClient.CreateOrUpdatePreparer(context.Context, string, string, string, Application) (*http.Request, error)
1. ConnectorApplicationClient.CreateOrUpdateResponder(*http.Response) (Application, error)
1. ConnectorApplicationClient.CreateOrUpdateSender(*http.Request) (*http.Response, error)
1. ConnectorApplicationClient.Delete(context.Context, string, string, string) (autorest.Response, error)
1. ConnectorApplicationClient.DeletePreparer(context.Context, string, string, string) (*http.Request, error)
1. ConnectorApplicationClient.DeleteResponder(*http.Response) (autorest.Response, error)
1. ConnectorApplicationClient.DeleteSender(*http.Request) (*http.Response, error)
1. ConnectorApplicationClient.Get(context.Context, string, string, string) (Application, error)
1. ConnectorApplicationClient.GetPreparer(context.Context, string, string, string) (*http.Request, error)
1. ConnectorApplicationClient.GetResponder(*http.Response) (Application, error)
1. ConnectorApplicationClient.GetSender(*http.Request) (*http.Response, error)
1. ConnectorApplicationsClient.List(context.Context, string, string) (ApplicationsListPage, error)
1. ConnectorApplicationsClient.ListComplete(context.Context, string, string) (ApplicationsListIterator, error)
1. ConnectorApplicationsClient.ListPreparer(context.Context, string, string) (*http.Request, error)
1. ConnectorApplicationsClient.ListResponder(*http.Response) (ApplicationsList, error)
1. ConnectorApplicationsClient.ListSender(*http.Request) (*http.Response, error)
1. ConnectorGovernanceRuleClient.List(context.Context, string, string) (GovernanceRuleListPage, error)
1. ConnectorGovernanceRuleClient.ListComplete(context.Context, string, string) (GovernanceRuleListIterator, error)
1. ConnectorGovernanceRuleClient.ListPreparer(context.Context, string, string) (*http.Request, error)
1. ConnectorGovernanceRuleClient.ListResponder(*http.Response) (GovernanceRuleList, error)
1. ConnectorGovernanceRuleClient.ListSender(*http.Request) (*http.Response, error)
1. ConnectorGovernanceRulesClient.CreateOrUpdate(context.Context, string, string, string, GovernanceRule) (GovernanceRule, error)
1. ConnectorGovernanceRulesClient.CreateOrUpdatePreparer(context.Context, string, string, string, GovernanceRule) (*http.Request, error)
1. ConnectorGovernanceRulesClient.CreateOrUpdateResponder(*http.Response) (GovernanceRule, error)
1. ConnectorGovernanceRulesClient.CreateOrUpdateSender(*http.Request) (*http.Response, error)
1. ConnectorGovernanceRulesClient.Delete(context.Context, string, string, string) (autorest.Response, error)
1. ConnectorGovernanceRulesClient.DeletePreparer(context.Context, string, string, string) (*http.Request, error)
1. ConnectorGovernanceRulesClient.DeleteResponder(*http.Response) (autorest.Response, error)
1. ConnectorGovernanceRulesClient.DeleteSender(*http.Request) (*http.Response, error)
1. ConnectorGovernanceRulesClient.Get(context.Context, string, string, string) (GovernanceRule, error)
1. ConnectorGovernanceRulesClient.GetPreparer(context.Context, string, string, string) (*http.Request, error)
1. ConnectorGovernanceRulesClient.GetResponder(*http.Response) (GovernanceRule, error)
1. ConnectorGovernanceRulesClient.GetSender(*http.Request) (*http.Response, error)
1. ConnectorGovernanceRulesExecuteStatusClient.Get(context.Context, string, string, string, string) (ConnectorGovernanceRulesExecuteStatusGetFuture, error)
1. ConnectorGovernanceRulesExecuteStatusClient.GetPreparer(context.Context, string, string, string, string) (*http.Request, error)
1. ConnectorGovernanceRulesExecuteStatusClient.GetResponder(*http.Response) (ExecuteRuleStatus, error)
1. ConnectorGovernanceRulesExecuteStatusClient.GetSender(*http.Request) (ConnectorGovernanceRulesExecuteStatusGetFuture, error)
1. ConnectorProperties.MarshalJSON() ([]byte, error)
1. ConnectorsGroupClient.CreateOrUpdate(context.Context, string, string, Connector) (Connector, error)
1. ConnectorsGroupClient.CreateOrUpdatePreparer(context.Context, string, string, Connector) (*http.Request, error)
1. ConnectorsGroupClient.CreateOrUpdateResponder(*http.Response) (Connector, error)
1. ConnectorsGroupClient.CreateOrUpdateSender(*http.Request) (*http.Response, error)
1. ConnectorsGroupClient.Delete(context.Context, string, string) (autorest.Response, error)
1. ConnectorsGroupClient.DeletePreparer(context.Context, string, string) (*http.Request, error)
1. ConnectorsGroupClient.DeleteResponder(*http.Response) (autorest.Response, error)
1. ConnectorsGroupClient.DeleteSender(*http.Request) (*http.Response, error)
1. ConnectorsGroupClient.Get(context.Context, string, string) (Connector, error)
1. ConnectorsGroupClient.GetPreparer(context.Context, string, string) (*http.Request, error)
1. ConnectorsGroupClient.GetResponder(*http.Response) (Connector, error)
1. ConnectorsGroupClient.GetSender(*http.Request) (*http.Response, error)
1. ConnectorsGroupClient.List(context.Context) (ConnectorsListPage, error)
1. ConnectorsGroupClient.ListByResourceGroup(context.Context, string) (ConnectorsListPage, error)
1. ConnectorsGroupClient.ListByResourceGroupComplete(context.Context, string) (ConnectorsListIterator, error)
1. ConnectorsGroupClient.ListByResourceGroupPreparer(context.Context, string) (*http.Request, error)
1. ConnectorsGroupClient.ListByResourceGroupResponder(*http.Response) (ConnectorsList, error)
1. ConnectorsGroupClient.ListByResourceGroupSender(*http.Request) (*http.Response, error)
1. ConnectorsGroupClient.ListComplete(context.Context) (ConnectorsListIterator, error)
1. ConnectorsGroupClient.ListPreparer(context.Context) (*http.Request, error)
1. ConnectorsGroupClient.ListResponder(*http.Response) (ConnectorsList, error)
1. ConnectorsGroupClient.ListSender(*http.Request) (*http.Response, error)
1. ConnectorsGroupClient.Update(context.Context, string, string, Connector) (Connector, error)
1. ConnectorsGroupClient.UpdatePreparer(context.Context, string, string, Connector) (*http.Request, error)
1. ConnectorsGroupClient.UpdateResponder(*http.Response) (Connector, error)
1. ConnectorsGroupClient.UpdateSender(*http.Request) (*http.Response, error)
1. ConnectorsList.IsEmpty() bool
1. ConnectorsList.MarshalJSON() ([]byte, error)
1. ConnectorsListIterator.NotDone() bool
1. ConnectorsListIterator.Response() ConnectorsList
1. ConnectorsListIterator.Value() Connector
1. ConnectorsListPage.NotDone() bool
1. ConnectorsListPage.Response() ConnectorsList
1. ConnectorsListPage.Values() []Connector
1. CspmMonitorAwsOffering.AsBasicCloudOffering() (BasicCloudOffering, bool)
1. CspmMonitorAwsOffering.AsCloudOffering() (*CloudOffering, bool)
1. CspmMonitorAwsOffering.AsCspmMonitorAwsOffering() (*CspmMonitorAwsOffering, bool)
1. CspmMonitorAwsOffering.AsCspmMonitorAzureDevOpsOffering() (*CspmMonitorAzureDevOpsOffering, bool)
1. CspmMonitorAwsOffering.AsCspmMonitorGcpOffering() (*CspmMonitorGcpOffering, bool)
1. CspmMonitorAwsOffering.AsCspmMonitorGithubOffering() (*CspmMonitorGithubOffering, bool)
1. CspmMonitorAwsOffering.AsDefenderCspmAwsOffering() (*DefenderCspmAwsOffering, bool)
1. CspmMonitorAwsOffering.AsDefenderCspmGcpOffering() (*DefenderCspmGcpOffering, bool)
1. CspmMonitorAwsOffering.AsDefenderFoDatabasesAwsOffering() (*DefenderFoDatabasesAwsOffering, bool)
1. CspmMonitorAwsOffering.AsDefenderForContainersAwsOffering() (*DefenderForContainersAwsOffering, bool)
1. CspmMonitorAwsOffering.AsDefenderForContainersGcpOffering() (*DefenderForContainersGcpOffering, bool)
1. CspmMonitorAwsOffering.AsDefenderForDatabasesGcpOffering() (*DefenderForDatabasesGcpOffering, bool)
1. CspmMonitorAwsOffering.AsDefenderForDevOpsAzureDevOpsOffering() (*DefenderForDevOpsAzureDevOpsOffering, bool)
1. CspmMonitorAwsOffering.AsDefenderForDevOpsGithubOffering() (*DefenderForDevOpsGithubOffering, bool)
1. CspmMonitorAwsOffering.AsDefenderForServersAwsOffering() (*DefenderForServersAwsOffering, bool)
1. CspmMonitorAwsOffering.AsDefenderForServersGcpOffering() (*DefenderForServersGcpOffering, bool)
1. CspmMonitorAwsOffering.AsInformationProtectionAwsOffering() (*InformationProtectionAwsOffering, bool)
1. CspmMonitorAwsOffering.MarshalJSON() ([]byte, error)
1. CspmMonitorAzureDevOpsOffering.AsBasicCloudOffering() (BasicCloudOffering, bool)
1. CspmMonitorAzureDevOpsOffering.AsCloudOffering() (*CloudOffering, bool)
1. CspmMonitorAzureDevOpsOffering.AsCspmMonitorAwsOffering() (*CspmMonitorAwsOffering, bool)
1. CspmMonitorAzureDevOpsOffering.AsCspmMonitorAzureDevOpsOffering() (*CspmMonitorAzureDevOpsOffering, bool)
1. CspmMonitorAzureDevOpsOffering.AsCspmMonitorGcpOffering() (*CspmMonitorGcpOffering, bool)
1. CspmMonitorAzureDevOpsOffering.AsCspmMonitorGithubOffering() (*CspmMonitorGithubOffering, bool)
1. CspmMonitorAzureDevOpsOffering.AsDefenderCspmAwsOffering() (*DefenderCspmAwsOffering, bool)
1. CspmMonitorAzureDevOpsOffering.AsDefenderCspmGcpOffering() (*DefenderCspmGcpOffering, bool)
1. CspmMonitorAzureDevOpsOffering.AsDefenderFoDatabasesAwsOffering() (*DefenderFoDatabasesAwsOffering, bool)
1. CspmMonitorAzureDevOpsOffering.AsDefenderForContainersAwsOffering() (*DefenderForContainersAwsOffering, bool)
1. CspmMonitorAzureDevOpsOffering.AsDefenderForContainersGcpOffering() (*DefenderForContainersGcpOffering, bool)
1. CspmMonitorAzureDevOpsOffering.AsDefenderForDatabasesGcpOffering() (*DefenderForDatabasesGcpOffering, bool)
1. CspmMonitorAzureDevOpsOffering.AsDefenderForDevOpsAzureDevOpsOffering() (*DefenderForDevOpsAzureDevOpsOffering, bool)
1. CspmMonitorAzureDevOpsOffering.AsDefenderForDevOpsGithubOffering() (*DefenderForDevOpsGithubOffering, bool)
1. CspmMonitorAzureDevOpsOffering.AsDefenderForServersAwsOffering() (*DefenderForServersAwsOffering, bool)
1. CspmMonitorAzureDevOpsOffering.AsDefenderForServersGcpOffering() (*DefenderForServersGcpOffering, bool)
1. CspmMonitorAzureDevOpsOffering.AsInformationProtectionAwsOffering() (*InformationProtectionAwsOffering, bool)
1. CspmMonitorAzureDevOpsOffering.MarshalJSON() ([]byte, error)
1. CspmMonitorGcpOffering.AsBasicCloudOffering() (BasicCloudOffering, bool)
1. CspmMonitorGcpOffering.AsCloudOffering() (*CloudOffering, bool)
1. CspmMonitorGcpOffering.AsCspmMonitorAwsOffering() (*CspmMonitorAwsOffering, bool)
1. CspmMonitorGcpOffering.AsCspmMonitorAzureDevOpsOffering() (*CspmMonitorAzureDevOpsOffering, bool)
1. CspmMonitorGcpOffering.AsCspmMonitorGcpOffering() (*CspmMonitorGcpOffering, bool)
1. CspmMonitorGcpOffering.AsCspmMonitorGithubOffering() (*CspmMonitorGithubOffering, bool)
1. CspmMonitorGcpOffering.AsDefenderCspmAwsOffering() (*DefenderCspmAwsOffering, bool)
1. CspmMonitorGcpOffering.AsDefenderCspmGcpOffering() (*DefenderCspmGcpOffering, bool)
1. CspmMonitorGcpOffering.AsDefenderFoDatabasesAwsOffering() (*DefenderFoDatabasesAwsOffering, bool)
1. CspmMonitorGcpOffering.AsDefenderForContainersAwsOffering() (*DefenderForContainersAwsOffering, bool)
1. CspmMonitorGcpOffering.AsDefenderForContainersGcpOffering() (*DefenderForContainersGcpOffering, bool)
1. CspmMonitorGcpOffering.AsDefenderForDatabasesGcpOffering() (*DefenderForDatabasesGcpOffering, bool)
1. CspmMonitorGcpOffering.AsDefenderForDevOpsAzureDevOpsOffering() (*DefenderForDevOpsAzureDevOpsOffering, bool)
1. CspmMonitorGcpOffering.AsDefenderForDevOpsGithubOffering() (*DefenderForDevOpsGithubOffering, bool)
1. CspmMonitorGcpOffering.AsDefenderForServersAwsOffering() (*DefenderForServersAwsOffering, bool)
1. CspmMonitorGcpOffering.AsDefenderForServersGcpOffering() (*DefenderForServersGcpOffering, bool)
1. CspmMonitorGcpOffering.AsInformationProtectionAwsOffering() (*InformationProtectionAwsOffering, bool)
1. CspmMonitorGcpOffering.MarshalJSON() ([]byte, error)
1. CspmMonitorGithubOffering.AsBasicCloudOffering() (BasicCloudOffering, bool)
1. CspmMonitorGithubOffering.AsCloudOffering() (*CloudOffering, bool)
1. CspmMonitorGithubOffering.AsCspmMonitorAwsOffering() (*CspmMonitorAwsOffering, bool)
1. CspmMonitorGithubOffering.AsCspmMonitorAzureDevOpsOffering() (*CspmMonitorAzureDevOpsOffering, bool)
1. CspmMonitorGithubOffering.AsCspmMonitorGcpOffering() (*CspmMonitorGcpOffering, bool)
1. CspmMonitorGithubOffering.AsCspmMonitorGithubOffering() (*CspmMonitorGithubOffering, bool)
1. CspmMonitorGithubOffering.AsDefenderCspmAwsOffering() (*DefenderCspmAwsOffering, bool)
1. CspmMonitorGithubOffering.AsDefenderCspmGcpOffering() (*DefenderCspmGcpOffering, bool)
1. CspmMonitorGithubOffering.AsDefenderFoDatabasesAwsOffering() (*DefenderFoDatabasesAwsOffering, bool)
1. CspmMonitorGithubOffering.AsDefenderForContainersAwsOffering() (*DefenderForContainersAwsOffering, bool)
1. CspmMonitorGithubOffering.AsDefenderForContainersGcpOffering() (*DefenderForContainersGcpOffering, bool)
1. CspmMonitorGithubOffering.AsDefenderForDatabasesGcpOffering() (*DefenderForDatabasesGcpOffering, bool)
1. CspmMonitorGithubOffering.AsDefenderForDevOpsAzureDevOpsOffering() (*DefenderForDevOpsAzureDevOpsOffering, bool)
1. CspmMonitorGithubOffering.AsDefenderForDevOpsGithubOffering() (*DefenderForDevOpsGithubOffering, bool)
1. CspmMonitorGithubOffering.AsDefenderForServersAwsOffering() (*DefenderForServersAwsOffering, bool)
1. CspmMonitorGithubOffering.AsDefenderForServersGcpOffering() (*DefenderForServersGcpOffering, bool)
1. CspmMonitorGithubOffering.AsInformationProtectionAwsOffering() (*InformationProtectionAwsOffering, bool)
1. CspmMonitorGithubOffering.MarshalJSON() ([]byte, error)
1. CustomAssessmentAutomation.MarshalJSON() ([]byte, error)
1. CustomAssessmentAutomationRequest.MarshalJSON() ([]byte, error)
1. CustomAssessmentAutomationsClient.Create(context.Context, string, string, CustomAssessmentAutomationRequest) (CustomAssessmentAutomation, error)
1. CustomAssessmentAutomationsClient.CreatePreparer(context.Context, string, string, CustomAssessmentAutomationRequest) (*http.Request, error)
1. CustomAssessmentAutomationsClient.CreateResponder(*http.Response) (CustomAssessmentAutomation, error)
1. CustomAssessmentAutomationsClient.CreateSender(*http.Request) (*http.Response, error)
1. CustomAssessmentAutomationsClient.Delete(context.Context, string, string) (autorest.Response, error)
1. CustomAssessmentAutomationsClient.DeletePreparer(context.Context, string, string) (*http.Request, error)
1. CustomAssessmentAutomationsClient.DeleteResponder(*http.Response) (autorest.Response, error)
1. CustomAssessmentAutomationsClient.DeleteSender(*http.Request) (*http.Response, error)
1. CustomAssessmentAutomationsClient.Get(context.Context, string, string) (CustomAssessmentAutomation, error)
1. CustomAssessmentAutomationsClient.GetPreparer(context.Context, string, string) (*http.Request, error)
1. CustomAssessmentAutomationsClient.GetResponder(*http.Response) (CustomAssessmentAutomation, error)
1. CustomAssessmentAutomationsClient.GetSender(*http.Request) (*http.Response, error)
1. CustomAssessmentAutomationsClient.ListByResourceGroup(context.Context, string) (CustomAssessmentAutomationsListResultPage, error)
1. CustomAssessmentAutomationsClient.ListByResourceGroupComplete(context.Context, string) (CustomAssessmentAutomationsListResultIterator, error)
1. CustomAssessmentAutomationsClient.ListByResourceGroupPreparer(context.Context, string) (*http.Request, error)
1. CustomAssessmentAutomationsClient.ListByResourceGroupResponder(*http.Response) (CustomAssessmentAutomationsListResult, error)
1. CustomAssessmentAutomationsClient.ListByResourceGroupSender(*http.Request) (*http.Response, error)
1. CustomAssessmentAutomationsClient.ListBySubscription(context.Context) (CustomAssessmentAutomationsListResultPage, error)
1. CustomAssessmentAutomationsClient.ListBySubscriptionComplete(context.Context) (CustomAssessmentAutomationsListResultIterator, error)
1. CustomAssessmentAutomationsClient.ListBySubscriptionPreparer(context.Context) (*http.Request, error)
1. CustomAssessmentAutomationsClient.ListBySubscriptionResponder(*http.Response) (CustomAssessmentAutomationsListResult, error)
1. CustomAssessmentAutomationsClient.ListBySubscriptionSender(*http.Request) (*http.Response, error)
1. CustomAssessmentAutomationsListResult.IsEmpty() bool
1. CustomAssessmentAutomationsListResult.MarshalJSON() ([]byte, error)
1. CustomAssessmentAutomationsListResultIterator.NotDone() bool
1. CustomAssessmentAutomationsListResultIterator.Response() CustomAssessmentAutomationsListResult
1. CustomAssessmentAutomationsListResultIterator.Value() CustomAssessmentAutomation
1. CustomAssessmentAutomationsListResultPage.NotDone() bool
1. CustomAssessmentAutomationsListResultPage.Response() CustomAssessmentAutomationsListResult
1. CustomAssessmentAutomationsListResultPage.Values() []CustomAssessmentAutomation
1. CustomEntityStoreAssignment.MarshalJSON() ([]byte, error)
1. CustomEntityStoreAssignmentRequest.MarshalJSON() ([]byte, error)
1. CustomEntityStoreAssignmentsClient.Create(context.Context, string, string, CustomEntityStoreAssignmentRequest) (CustomEntityStoreAssignment, error)
1. CustomEntityStoreAssignmentsClient.CreatePreparer(context.Context, string, string, CustomEntityStoreAssignmentRequest) (*http.Request, error)
1. CustomEntityStoreAssignmentsClient.CreateResponder(*http.Response) (CustomEntityStoreAssignment, error)
1. CustomEntityStoreAssignmentsClient.CreateSender(*http.Request) (*http.Response, error)
1. CustomEntityStoreAssignmentsClient.Delete(context.Context, string, string) (autorest.Response, error)
1. CustomEntityStoreAssignmentsClient.DeletePreparer(context.Context, string, string) (*http.Request, error)
1. CustomEntityStoreAssignmentsClient.DeleteResponder(*http.Response) (autorest.Response, error)
1. CustomEntityStoreAssignmentsClient.DeleteSender(*http.Request) (*http.Response, error)
1. CustomEntityStoreAssignmentsClient.Get(context.Context, string, string) (CustomEntityStoreAssignment, error)
1. CustomEntityStoreAssignmentsClient.GetPreparer(context.Context, string, string) (*http.Request, error)
1. CustomEntityStoreAssignmentsClient.GetResponder(*http.Response) (CustomEntityStoreAssignment, error)
1. CustomEntityStoreAssignmentsClient.GetSender(*http.Request) (*http.Response, error)
1. CustomEntityStoreAssignmentsClient.ListByResourceGroup(context.Context, string) (CustomEntityStoreAssignmentsListResultPage, error)
1. CustomEntityStoreAssignmentsClient.ListByResourceGroupComplete(context.Context, string) (CustomEntityStoreAssignmentsListResultIterator, error)
1. CustomEntityStoreAssignmentsClient.ListByResourceGroupPreparer(context.Context, string) (*http.Request, error)
1. CustomEntityStoreAssignmentsClient.ListByResourceGroupResponder(*http.Response) (CustomEntityStoreAssignmentsListResult, error)
1. CustomEntityStoreAssignmentsClient.ListByResourceGroupSender(*http.Request) (*http.Response, error)
1. CustomEntityStoreAssignmentsClient.ListBySubscription(context.Context) (CustomEntityStoreAssignmentsListResultPage, error)
1. CustomEntityStoreAssignmentsClient.ListBySubscriptionComplete(context.Context) (CustomEntityStoreAssignmentsListResultIterator, error)
1. CustomEntityStoreAssignmentsClient.ListBySubscriptionPreparer(context.Context) (*http.Request, error)
1. CustomEntityStoreAssignmentsClient.ListBySubscriptionResponder(*http.Response) (CustomEntityStoreAssignmentsListResult, error)
1. CustomEntityStoreAssignmentsClient.ListBySubscriptionSender(*http.Request) (*http.Response, error)
1. CustomEntityStoreAssignmentsListResult.IsEmpty() bool
1. CustomEntityStoreAssignmentsListResult.MarshalJSON() ([]byte, error)
1. CustomEntityStoreAssignmentsListResultIterator.NotDone() bool
1. CustomEntityStoreAssignmentsListResultIterator.Response() CustomEntityStoreAssignmentsListResult
1. CustomEntityStoreAssignmentsListResultIterator.Value() CustomEntityStoreAssignment
1. CustomEntityStoreAssignmentsListResultPage.NotDone() bool
1. CustomEntityStoreAssignmentsListResultPage.Response() CustomEntityStoreAssignmentsListResult
1. CustomEntityStoreAssignmentsListResultPage.Values() []CustomEntityStoreAssignment
1. DefenderCspmAwsOffering.AsBasicCloudOffering() (BasicCloudOffering, bool)
1. DefenderCspmAwsOffering.AsCloudOffering() (*CloudOffering, bool)
1. DefenderCspmAwsOffering.AsCspmMonitorAwsOffering() (*CspmMonitorAwsOffering, bool)
1. DefenderCspmAwsOffering.AsCspmMonitorAzureDevOpsOffering() (*CspmMonitorAzureDevOpsOffering, bool)
1. DefenderCspmAwsOffering.AsCspmMonitorGcpOffering() (*CspmMonitorGcpOffering, bool)
1. DefenderCspmAwsOffering.AsCspmMonitorGithubOffering() (*CspmMonitorGithubOffering, bool)
1. DefenderCspmAwsOffering.AsDefenderCspmAwsOffering() (*DefenderCspmAwsOffering, bool)
1. DefenderCspmAwsOffering.AsDefenderCspmGcpOffering() (*DefenderCspmGcpOffering, bool)
1. DefenderCspmAwsOffering.AsDefenderFoDatabasesAwsOffering() (*DefenderFoDatabasesAwsOffering, bool)
1. DefenderCspmAwsOffering.AsDefenderForContainersAwsOffering() (*DefenderForContainersAwsOffering, bool)
1. DefenderCspmAwsOffering.AsDefenderForContainersGcpOffering() (*DefenderForContainersGcpOffering, bool)
1. DefenderCspmAwsOffering.AsDefenderForDatabasesGcpOffering() (*DefenderForDatabasesGcpOffering, bool)
1. DefenderCspmAwsOffering.AsDefenderForDevOpsAzureDevOpsOffering() (*DefenderForDevOpsAzureDevOpsOffering, bool)
1. DefenderCspmAwsOffering.AsDefenderForDevOpsGithubOffering() (*DefenderForDevOpsGithubOffering, bool)
1. DefenderCspmAwsOffering.AsDefenderForServersAwsOffering() (*DefenderForServersAwsOffering, bool)
1. DefenderCspmAwsOffering.AsDefenderForServersGcpOffering() (*DefenderForServersGcpOffering, bool)
1. DefenderCspmAwsOffering.AsInformationProtectionAwsOffering() (*InformationProtectionAwsOffering, bool)
1. DefenderCspmAwsOffering.MarshalJSON() ([]byte, error)
1. DefenderCspmAwsOfferingVMScannersConfiguration.MarshalJSON() ([]byte, error)
1. DefenderCspmGcpOffering.AsBasicCloudOffering() (BasicCloudOffering, bool)
1. DefenderCspmGcpOffering.AsCloudOffering() (*CloudOffering, bool)
1. DefenderCspmGcpOffering.AsCspmMonitorAwsOffering() (*CspmMonitorAwsOffering, bool)
1. DefenderCspmGcpOffering.AsCspmMonitorAzureDevOpsOffering() (*CspmMonitorAzureDevOpsOffering, bool)
1. DefenderCspmGcpOffering.AsCspmMonitorGcpOffering() (*CspmMonitorGcpOffering, bool)
1. DefenderCspmGcpOffering.AsCspmMonitorGithubOffering() (*CspmMonitorGithubOffering, bool)
1. DefenderCspmGcpOffering.AsDefenderCspmAwsOffering() (*DefenderCspmAwsOffering, bool)
1. DefenderCspmGcpOffering.AsDefenderCspmGcpOffering() (*DefenderCspmGcpOffering, bool)
1. DefenderCspmGcpOffering.AsDefenderFoDatabasesAwsOffering() (*DefenderFoDatabasesAwsOffering, bool)
1. DefenderCspmGcpOffering.AsDefenderForContainersAwsOffering() (*DefenderForContainersAwsOffering, bool)
1. DefenderCspmGcpOffering.AsDefenderForContainersGcpOffering() (*DefenderForContainersGcpOffering, bool)
1. DefenderCspmGcpOffering.AsDefenderForDatabasesGcpOffering() (*DefenderForDatabasesGcpOffering, bool)
1. DefenderCspmGcpOffering.AsDefenderForDevOpsAzureDevOpsOffering() (*DefenderForDevOpsAzureDevOpsOffering, bool)
1. DefenderCspmGcpOffering.AsDefenderForDevOpsGithubOffering() (*DefenderForDevOpsGithubOffering, bool)
1. DefenderCspmGcpOffering.AsDefenderForServersAwsOffering() (*DefenderForServersAwsOffering, bool)
1. DefenderCspmGcpOffering.AsDefenderForServersGcpOffering() (*DefenderForServersGcpOffering, bool)
1. DefenderCspmGcpOffering.AsInformationProtectionAwsOffering() (*InformationProtectionAwsOffering, bool)
1. DefenderCspmGcpOffering.MarshalJSON() ([]byte, error)
1. DefenderFoDatabasesAwsOffering.AsBasicCloudOffering() (BasicCloudOffering, bool)
1. DefenderFoDatabasesAwsOffering.AsCloudOffering() (*CloudOffering, bool)
1. DefenderFoDatabasesAwsOffering.AsCspmMonitorAwsOffering() (*CspmMonitorAwsOffering, bool)
1. DefenderFoDatabasesAwsOffering.AsCspmMonitorAzureDevOpsOffering() (*CspmMonitorAzureDevOpsOffering, bool)
1. DefenderFoDatabasesAwsOffering.AsCspmMonitorGcpOffering() (*CspmMonitorGcpOffering, bool)
1. DefenderFoDatabasesAwsOffering.AsCspmMonitorGithubOffering() (*CspmMonitorGithubOffering, bool)
1. DefenderFoDatabasesAwsOffering.AsDefenderCspmAwsOffering() (*DefenderCspmAwsOffering, bool)
1. DefenderFoDatabasesAwsOffering.AsDefenderCspmGcpOffering() (*DefenderCspmGcpOffering, bool)
1. DefenderFoDatabasesAwsOffering.AsDefenderFoDatabasesAwsOffering() (*DefenderFoDatabasesAwsOffering, bool)
1. DefenderFoDatabasesAwsOffering.AsDefenderForContainersAwsOffering() (*DefenderForContainersAwsOffering, bool)
1. DefenderFoDatabasesAwsOffering.AsDefenderForContainersGcpOffering() (*DefenderForContainersGcpOffering, bool)
1. DefenderFoDatabasesAwsOffering.AsDefenderForDatabasesGcpOffering() (*DefenderForDatabasesGcpOffering, bool)
1. DefenderFoDatabasesAwsOffering.AsDefenderForDevOpsAzureDevOpsOffering() (*DefenderForDevOpsAzureDevOpsOffering, bool)
1. DefenderFoDatabasesAwsOffering.AsDefenderForDevOpsGithubOffering() (*DefenderForDevOpsGithubOffering, bool)
1. DefenderFoDatabasesAwsOffering.AsDefenderForServersAwsOffering() (*DefenderForServersAwsOffering, bool)
1. DefenderFoDatabasesAwsOffering.AsDefenderForServersGcpOffering() (*DefenderForServersGcpOffering, bool)
1. DefenderFoDatabasesAwsOffering.AsInformationProtectionAwsOffering() (*InformationProtectionAwsOffering, bool)
1. DefenderFoDatabasesAwsOffering.MarshalJSON() ([]byte, error)
1. DefenderForContainersAwsOffering.AsBasicCloudOffering() (BasicCloudOffering, bool)
1. DefenderForContainersAwsOffering.AsCloudOffering() (*CloudOffering, bool)
1. DefenderForContainersAwsOffering.AsCspmMonitorAwsOffering() (*CspmMonitorAwsOffering, bool)
1. DefenderForContainersAwsOffering.AsCspmMonitorAzureDevOpsOffering() (*CspmMonitorAzureDevOpsOffering, bool)
1. DefenderForContainersAwsOffering.AsCspmMonitorGcpOffering() (*CspmMonitorGcpOffering, bool)
1. DefenderForContainersAwsOffering.AsCspmMonitorGithubOffering() (*CspmMonitorGithubOffering, bool)
1. DefenderForContainersAwsOffering.AsDefenderCspmAwsOffering() (*DefenderCspmAwsOffering, bool)
1. DefenderForContainersAwsOffering.AsDefenderCspmGcpOffering() (*DefenderCspmGcpOffering, bool)
1. DefenderForContainersAwsOffering.AsDefenderFoDatabasesAwsOffering() (*DefenderFoDatabasesAwsOffering, bool)
1. DefenderForContainersAwsOffering.AsDefenderForContainersAwsOffering() (*DefenderForContainersAwsOffering, bool)
1. DefenderForContainersAwsOffering.AsDefenderForContainersGcpOffering() (*DefenderForContainersGcpOffering, bool)
1. DefenderForContainersAwsOffering.AsDefenderForDatabasesGcpOffering() (*DefenderForDatabasesGcpOffering, bool)
1. DefenderForContainersAwsOffering.AsDefenderForDevOpsAzureDevOpsOffering() (*DefenderForDevOpsAzureDevOpsOffering, bool)
1. DefenderForContainersAwsOffering.AsDefenderForDevOpsGithubOffering() (*DefenderForDevOpsGithubOffering, bool)
1. DefenderForContainersAwsOffering.AsDefenderForServersAwsOffering() (*DefenderForServersAwsOffering, bool)
1. DefenderForContainersAwsOffering.AsDefenderForServersGcpOffering() (*DefenderForServersGcpOffering, bool)
1. DefenderForContainersAwsOffering.AsInformationProtectionAwsOffering() (*InformationProtectionAwsOffering, bool)
1. DefenderForContainersAwsOffering.MarshalJSON() ([]byte, error)
1. DefenderForContainersGcpOffering.AsBasicCloudOffering() (BasicCloudOffering, bool)
1. DefenderForContainersGcpOffering.AsCloudOffering() (*CloudOffering, bool)
1. DefenderForContainersGcpOffering.AsCspmMonitorAwsOffering() (*CspmMonitorAwsOffering, bool)
1. DefenderForContainersGcpOffering.AsCspmMonitorAzureDevOpsOffering() (*CspmMonitorAzureDevOpsOffering, bool)
1. DefenderForContainersGcpOffering.AsCspmMonitorGcpOffering() (*CspmMonitorGcpOffering, bool)
1. DefenderForContainersGcpOffering.AsCspmMonitorGithubOffering() (*CspmMonitorGithubOffering, bool)
1. DefenderForContainersGcpOffering.AsDefenderCspmAwsOffering() (*DefenderCspmAwsOffering, bool)
1. DefenderForContainersGcpOffering.AsDefenderCspmGcpOffering() (*DefenderCspmGcpOffering, bool)
1. DefenderForContainersGcpOffering.AsDefenderFoDatabasesAwsOffering() (*DefenderFoDatabasesAwsOffering, bool)
1. DefenderForContainersGcpOffering.AsDefenderForContainersAwsOffering() (*DefenderForContainersAwsOffering, bool)
1. DefenderForContainersGcpOffering.AsDefenderForContainersGcpOffering() (*DefenderForContainersGcpOffering, bool)
1. DefenderForContainersGcpOffering.AsDefenderForDatabasesGcpOffering() (*DefenderForDatabasesGcpOffering, bool)
1. DefenderForContainersGcpOffering.AsDefenderForDevOpsAzureDevOpsOffering() (*DefenderForDevOpsAzureDevOpsOffering, bool)
1. DefenderForContainersGcpOffering.AsDefenderForDevOpsGithubOffering() (*DefenderForDevOpsGithubOffering, bool)
1. DefenderForContainersGcpOffering.AsDefenderForServersAwsOffering() (*DefenderForServersAwsOffering, bool)
1. DefenderForContainersGcpOffering.AsDefenderForServersGcpOffering() (*DefenderForServersGcpOffering, bool)
1. DefenderForContainersGcpOffering.AsInformationProtectionAwsOffering() (*InformationProtectionAwsOffering, bool)
1. DefenderForContainersGcpOffering.MarshalJSON() ([]byte, error)
1. DefenderForDatabasesGcpOffering.AsBasicCloudOffering() (BasicCloudOffering, bool)
1. DefenderForDatabasesGcpOffering.AsCloudOffering() (*CloudOffering, bool)
1. DefenderForDatabasesGcpOffering.AsCspmMonitorAwsOffering() (*CspmMonitorAwsOffering, bool)
1. DefenderForDatabasesGcpOffering.AsCspmMonitorAzureDevOpsOffering() (*CspmMonitorAzureDevOpsOffering, bool)
1. DefenderForDatabasesGcpOffering.AsCspmMonitorGcpOffering() (*CspmMonitorGcpOffering, bool)
1. DefenderForDatabasesGcpOffering.AsCspmMonitorGithubOffering() (*CspmMonitorGithubOffering, bool)
1. DefenderForDatabasesGcpOffering.AsDefenderCspmAwsOffering() (*DefenderCspmAwsOffering, bool)
1. DefenderForDatabasesGcpOffering.AsDefenderCspmGcpOffering() (*DefenderCspmGcpOffering, bool)
1. DefenderForDatabasesGcpOffering.AsDefenderFoDatabasesAwsOffering() (*DefenderFoDatabasesAwsOffering, bool)
1. DefenderForDatabasesGcpOffering.AsDefenderForContainersAwsOffering() (*DefenderForContainersAwsOffering, bool)
1. DefenderForDatabasesGcpOffering.AsDefenderForContainersGcpOffering() (*DefenderForContainersGcpOffering, bool)
1. DefenderForDatabasesGcpOffering.AsDefenderForDatabasesGcpOffering() (*DefenderForDatabasesGcpOffering, bool)
1. DefenderForDatabasesGcpOffering.AsDefenderForDevOpsAzureDevOpsOffering() (*DefenderForDevOpsAzureDevOpsOffering, bool)
1. DefenderForDatabasesGcpOffering.AsDefenderForDevOpsGithubOffering() (*DefenderForDevOpsGithubOffering, bool)
1. DefenderForDatabasesGcpOffering.AsDefenderForServersAwsOffering() (*DefenderForServersAwsOffering, bool)
1. DefenderForDatabasesGcpOffering.AsDefenderForServersGcpOffering() (*DefenderForServersGcpOffering, bool)
1. DefenderForDatabasesGcpOffering.AsInformationProtectionAwsOffering() (*InformationProtectionAwsOffering, bool)
1. DefenderForDatabasesGcpOffering.MarshalJSON() ([]byte, error)
1. DefenderForDevOpsAzureDevOpsOffering.AsBasicCloudOffering() (BasicCloudOffering, bool)
1. DefenderForDevOpsAzureDevOpsOffering.AsCloudOffering() (*CloudOffering, bool)
1. DefenderForDevOpsAzureDevOpsOffering.AsCspmMonitorAwsOffering() (*CspmMonitorAwsOffering, bool)
1. DefenderForDevOpsAzureDevOpsOffering.AsCspmMonitorAzureDevOpsOffering() (*CspmMonitorAzureDevOpsOffering, bool)
1. DefenderForDevOpsAzureDevOpsOffering.AsCspmMonitorGcpOffering() (*CspmMonitorGcpOffering, bool)
1. DefenderForDevOpsAzureDevOpsOffering.AsCspmMonitorGithubOffering() (*CspmMonitorGithubOffering, bool)
1. DefenderForDevOpsAzureDevOpsOffering.AsDefenderCspmAwsOffering() (*DefenderCspmAwsOffering, bool)
1. DefenderForDevOpsAzureDevOpsOffering.AsDefenderCspmGcpOffering() (*DefenderCspmGcpOffering, bool)
1. DefenderForDevOpsAzureDevOpsOffering.AsDefenderFoDatabasesAwsOffering() (*DefenderFoDatabasesAwsOffering, bool)
1. DefenderForDevOpsAzureDevOpsOffering.AsDefenderForContainersAwsOffering() (*DefenderForContainersAwsOffering, bool)
1. DefenderForDevOpsAzureDevOpsOffering.AsDefenderForContainersGcpOffering() (*DefenderForContainersGcpOffering, bool)
1. DefenderForDevOpsAzureDevOpsOffering.AsDefenderForDatabasesGcpOffering() (*DefenderForDatabasesGcpOffering, bool)
1. DefenderForDevOpsAzureDevOpsOffering.AsDefenderForDevOpsAzureDevOpsOffering() (*DefenderForDevOpsAzureDevOpsOffering, bool)
1. DefenderForDevOpsAzureDevOpsOffering.AsDefenderForDevOpsGithubOffering() (*DefenderForDevOpsGithubOffering, bool)
1. DefenderForDevOpsAzureDevOpsOffering.AsDefenderForServersAwsOffering() (*DefenderForServersAwsOffering, bool)
1. DefenderForDevOpsAzureDevOpsOffering.AsDefenderForServersGcpOffering() (*DefenderForServersGcpOffering, bool)
1. DefenderForDevOpsAzureDevOpsOffering.AsInformationProtectionAwsOffering() (*InformationProtectionAwsOffering, bool)
1. DefenderForDevOpsAzureDevOpsOffering.MarshalJSON() ([]byte, error)
1. DefenderForDevOpsGithubOffering.AsBasicCloudOffering() (BasicCloudOffering, bool)
1. DefenderForDevOpsGithubOffering.AsCloudOffering() (*CloudOffering, bool)
1. DefenderForDevOpsGithubOffering.AsCspmMonitorAwsOffering() (*CspmMonitorAwsOffering, bool)
1. DefenderForDevOpsGithubOffering.AsCspmMonitorAzureDevOpsOffering() (*CspmMonitorAzureDevOpsOffering, bool)
1. DefenderForDevOpsGithubOffering.AsCspmMonitorGcpOffering() (*CspmMonitorGcpOffering, bool)
1. DefenderForDevOpsGithubOffering.AsCspmMonitorGithubOffering() (*CspmMonitorGithubOffering, bool)
1. DefenderForDevOpsGithubOffering.AsDefenderCspmAwsOffering() (*DefenderCspmAwsOffering, bool)
1. DefenderForDevOpsGithubOffering.AsDefenderCspmGcpOffering() (*DefenderCspmGcpOffering, bool)
1. DefenderForDevOpsGithubOffering.AsDefenderFoDatabasesAwsOffering() (*DefenderFoDatabasesAwsOffering, bool)
1. DefenderForDevOpsGithubOffering.AsDefenderForContainersAwsOffering() (*DefenderForContainersAwsOffering, bool)
1. DefenderForDevOpsGithubOffering.AsDefenderForContainersGcpOffering() (*DefenderForContainersGcpOffering, bool)
1. DefenderForDevOpsGithubOffering.AsDefenderForDatabasesGcpOffering() (*DefenderForDatabasesGcpOffering, bool)
1. DefenderForDevOpsGithubOffering.AsDefenderForDevOpsAzureDevOpsOffering() (*DefenderForDevOpsAzureDevOpsOffering, bool)
1. DefenderForDevOpsGithubOffering.AsDefenderForDevOpsGithubOffering() (*DefenderForDevOpsGithubOffering, bool)
1. DefenderForDevOpsGithubOffering.AsDefenderForServersAwsOffering() (*DefenderForServersAwsOffering, bool)
1. DefenderForDevOpsGithubOffering.AsDefenderForServersGcpOffering() (*DefenderForServersGcpOffering, bool)
1. DefenderForDevOpsGithubOffering.AsInformationProtectionAwsOffering() (*InformationProtectionAwsOffering, bool)
1. DefenderForDevOpsGithubOffering.MarshalJSON() ([]byte, error)
1. DefenderForServersAwsOffering.AsBasicCloudOffering() (BasicCloudOffering, bool)
1. DefenderForServersAwsOffering.AsCloudOffering() (*CloudOffering, bool)
1. DefenderForServersAwsOffering.AsCspmMonitorAwsOffering() (*CspmMonitorAwsOffering, bool)
1. DefenderForServersAwsOffering.AsCspmMonitorAzureDevOpsOffering() (*CspmMonitorAzureDevOpsOffering, bool)
1. DefenderForServersAwsOffering.AsCspmMonitorGcpOffering() (*CspmMonitorGcpOffering, bool)
1. DefenderForServersAwsOffering.AsCspmMonitorGithubOffering() (*CspmMonitorGithubOffering, bool)
1. DefenderForServersAwsOffering.AsDefenderCspmAwsOffering() (*DefenderCspmAwsOffering, bool)
1. DefenderForServersAwsOffering.AsDefenderCspmGcpOffering() (*DefenderCspmGcpOffering, bool)
1. DefenderForServersAwsOffering.AsDefenderFoDatabasesAwsOffering() (*DefenderFoDatabasesAwsOffering, bool)
1. DefenderForServersAwsOffering.AsDefenderForContainersAwsOffering() (*DefenderForContainersAwsOffering, bool)
1. DefenderForServersAwsOffering.AsDefenderForContainersGcpOffering() (*DefenderForContainersGcpOffering, bool)
1. DefenderForServersAwsOffering.AsDefenderForDatabasesGcpOffering() (*DefenderForDatabasesGcpOffering, bool)
1. DefenderForServersAwsOffering.AsDefenderForDevOpsAzureDevOpsOffering() (*DefenderForDevOpsAzureDevOpsOffering, bool)
1. DefenderForServersAwsOffering.AsDefenderForDevOpsGithubOffering() (*DefenderForDevOpsGithubOffering, bool)
1. DefenderForServersAwsOffering.AsDefenderForServersAwsOffering() (*DefenderForServersAwsOffering, bool)
1. DefenderForServersAwsOffering.AsDefenderForServersGcpOffering() (*DefenderForServersGcpOffering, bool)
1. DefenderForServersAwsOffering.AsInformationProtectionAwsOffering() (*InformationProtectionAwsOffering, bool)
1. DefenderForServersAwsOffering.MarshalJSON() ([]byte, error)
1. DefenderForServersAwsOfferingVMScannersConfiguration.MarshalJSON() ([]byte, error)
1. DefenderForServersGcpOffering.AsBasicCloudOffering() (BasicCloudOffering, bool)
1. DefenderForServersGcpOffering.AsCloudOffering() (*CloudOffering, bool)
1. DefenderForServersGcpOffering.AsCspmMonitorAwsOffering() (*CspmMonitorAwsOffering, bool)
1. DefenderForServersGcpOffering.AsCspmMonitorAzureDevOpsOffering() (*CspmMonitorAzureDevOpsOffering, bool)
1. DefenderForServersGcpOffering.AsCspmMonitorGcpOffering() (*CspmMonitorGcpOffering, bool)
1. DefenderForServersGcpOffering.AsCspmMonitorGithubOffering() (*CspmMonitorGithubOffering, bool)
1. DefenderForServersGcpOffering.AsDefenderCspmAwsOffering() (*DefenderCspmAwsOffering, bool)
1. DefenderForServersGcpOffering.AsDefenderCspmGcpOffering() (*DefenderCspmGcpOffering, bool)
1. DefenderForServersGcpOffering.AsDefenderFoDatabasesAwsOffering() (*DefenderFoDatabasesAwsOffering, bool)
1. DefenderForServersGcpOffering.AsDefenderForContainersAwsOffering() (*DefenderForContainersAwsOffering, bool)
1. DefenderForServersGcpOffering.AsDefenderForContainersGcpOffering() (*DefenderForContainersGcpOffering, bool)
1. DefenderForServersGcpOffering.AsDefenderForDatabasesGcpOffering() (*DefenderForDatabasesGcpOffering, bool)
1. DefenderForServersGcpOffering.AsDefenderForDevOpsAzureDevOpsOffering() (*DefenderForDevOpsAzureDevOpsOffering, bool)
1. DefenderForServersGcpOffering.AsDefenderForDevOpsGithubOffering() (*DefenderForDevOpsGithubOffering, bool)
1. DefenderForServersGcpOffering.AsDefenderForServersAwsOffering() (*DefenderForServersAwsOffering, bool)
1. DefenderForServersGcpOffering.AsDefenderForServersGcpOffering() (*DefenderForServersGcpOffering, bool)
1. DefenderForServersGcpOffering.AsInformationProtectionAwsOffering() (*InformationProtectionAwsOffering, bool)
1. DefenderForServersGcpOffering.MarshalJSON() ([]byte, error)
1. EnvironmentData.AsAwsEnvironmentData() (*AwsEnvironmentData, bool)
1. EnvironmentData.AsAzureDevOpsScopeEnvironmentData() (*AzureDevOpsScopeEnvironmentData, bool)
1. EnvironmentData.AsBasicEnvironmentData() (BasicEnvironmentData, bool)
1. EnvironmentData.AsEnvironmentData() (*EnvironmentData, bool)
1. EnvironmentData.AsGcpProjectEnvironmentData() (*GcpProjectEnvironmentData, bool)
1. EnvironmentData.AsGithubScopeEnvironmentData() (*GithubScopeEnvironmentData, bool)
1. EnvironmentData.MarshalJSON() ([]byte, error)
1. ExecuteRuleStatus.MarshalJSON() ([]byte, error)
1. GcpOrganizationalData.AsBasicGcpOrganizationalData() (BasicGcpOrganizationalData, bool)
1. GcpOrganizationalData.AsGcpOrganizationalData() (*GcpOrganizationalData, bool)
1. GcpOrganizationalData.AsGcpOrganizationalDataMember() (*GcpOrganizationalDataMember, bool)
1. GcpOrganizationalData.AsGcpOrganizationalDataOrganization() (*GcpOrganizationalDataOrganization, bool)
1. GcpOrganizationalData.MarshalJSON() ([]byte, error)
1. GcpOrganizationalDataMember.AsBasicGcpOrganizationalData() (BasicGcpOrganizationalData, bool)
1. GcpOrganizationalDataMember.AsGcpOrganizationalData() (*GcpOrganizationalData, bool)
1. GcpOrganizationalDataMember.AsGcpOrganizationalDataMember() (*GcpOrganizationalDataMember, bool)
1. GcpOrganizationalDataMember.AsGcpOrganizationalDataOrganization() (*GcpOrganizationalDataOrganization, bool)
1. GcpOrganizationalDataMember.MarshalJSON() ([]byte, error)
1. GcpOrganizationalDataOrganization.AsBasicGcpOrganizationalData() (BasicGcpOrganizationalData, bool)
1. GcpOrganizationalDataOrganization.AsGcpOrganizationalData() (*GcpOrganizationalData, bool)
1. GcpOrganizationalDataOrganization.AsGcpOrganizationalDataMember() (*GcpOrganizationalDataMember, bool)
1. GcpOrganizationalDataOrganization.AsGcpOrganizationalDataOrganization() (*GcpOrganizationalDataOrganization, bool)
1. GcpOrganizationalDataOrganization.MarshalJSON() ([]byte, error)
1. GcpProjectDetails.MarshalJSON() ([]byte, error)
1. GcpProjectEnvironmentData.AsAwsEnvironmentData() (*AwsEnvironmentData, bool)
1. GcpProjectEnvironmentData.AsAzureDevOpsScopeEnvironmentData() (*AzureDevOpsScopeEnvironmentData, bool)
1. GcpProjectEnvironmentData.AsBasicEnvironmentData() (BasicEnvironmentData, bool)
1. GcpProjectEnvironmentData.AsEnvironmentData() (*EnvironmentData, bool)
1. GcpProjectEnvironmentData.AsGcpProjectEnvironmentData() (*GcpProjectEnvironmentData, bool)
1. GcpProjectEnvironmentData.AsGithubScopeEnvironmentData() (*GithubScopeEnvironmentData, bool)
1. GcpProjectEnvironmentData.MarshalJSON() ([]byte, error)
1. GithubScopeEnvironmentData.AsAwsEnvironmentData() (*AwsEnvironmentData, bool)
1. GithubScopeEnvironmentData.AsAzureDevOpsScopeEnvironmentData() (*AzureDevOpsScopeEnvironmentData, bool)
1. GithubScopeEnvironmentData.AsBasicEnvironmentData() (BasicEnvironmentData, bool)
1. GithubScopeEnvironmentData.AsEnvironmentData() (*EnvironmentData, bool)
1. GithubScopeEnvironmentData.AsGcpProjectEnvironmentData() (*GcpProjectEnvironmentData, bool)
1. GithubScopeEnvironmentData.AsGithubScopeEnvironmentData() (*GithubScopeEnvironmentData, bool)
1. GithubScopeEnvironmentData.MarshalJSON() ([]byte, error)
1. GovernanceAssignment.MarshalJSON() ([]byte, error)
1. GovernanceAssignmentsClient.CreateOrUpdate(context.Context, string, string, string, GovernanceAssignment) (GovernanceAssignment, error)
1. GovernanceAssignmentsClient.CreateOrUpdatePreparer(context.Context, string, string, string, GovernanceAssignment) (*http.Request, error)
1. GovernanceAssignmentsClient.CreateOrUpdateResponder(*http.Response) (GovernanceAssignment, error)
1. GovernanceAssignmentsClient.CreateOrUpdateSender(*http.Request) (*http.Response, error)
1. GovernanceAssignmentsClient.Delete(context.Context, string, string, string) (autorest.Response, error)
1. GovernanceAssignmentsClient.DeletePreparer(context.Context, string, string, string) (*http.Request, error)
1. GovernanceAssignmentsClient.DeleteResponder(*http.Response) (autorest.Response, error)
1. GovernanceAssignmentsClient.DeleteSender(*http.Request) (*http.Response, error)
1. GovernanceAssignmentsClient.Get(context.Context, string, string, string) (GovernanceAssignment, error)
1. GovernanceAssignmentsClient.GetPreparer(context.Context, string, string, string) (*http.Request, error)
1. GovernanceAssignmentsClient.GetResponder(*http.Response) (GovernanceAssignment, error)
1. GovernanceAssignmentsClient.GetSender(*http.Request) (*http.Response, error)
1. GovernanceAssignmentsClient.List(context.Context, string, string) (GovernanceAssignmentsListPage, error)
1. GovernanceAssignmentsClient.ListComplete(context.Context, string, string) (GovernanceAssignmentsListIterator, error)
1. GovernanceAssignmentsClient.ListPreparer(context.Context, string, string) (*http.Request, error)
1. GovernanceAssignmentsClient.ListResponder(*http.Response) (GovernanceAssignmentsList, error)
1. GovernanceAssignmentsClient.ListSender(*http.Request) (*http.Response, error)
1. GovernanceAssignmentsList.IsEmpty() bool
1. GovernanceAssignmentsList.MarshalJSON() ([]byte, error)
1. GovernanceAssignmentsListIterator.NotDone() bool
1. GovernanceAssignmentsListIterator.Response() GovernanceAssignmentsList
1. GovernanceAssignmentsListIterator.Value() GovernanceAssignment
1. GovernanceAssignmentsListPage.NotDone() bool
1. GovernanceAssignmentsListPage.Response() GovernanceAssignmentsList
1. GovernanceAssignmentsListPage.Values() []GovernanceAssignment
1. GovernanceRule.MarshalJSON() ([]byte, error)
1. GovernanceRuleClient.List(context.Context) (GovernanceRuleListPage, error)
1. GovernanceRuleClient.ListComplete(context.Context) (GovernanceRuleListIterator, error)
1. GovernanceRuleClient.ListPreparer(context.Context) (*http.Request, error)
1. GovernanceRuleClient.ListResponder(*http.Response) (GovernanceRuleList, error)
1. GovernanceRuleClient.ListSender(*http.Request) (*http.Response, error)
1. GovernanceRuleList.IsEmpty() bool
1. GovernanceRuleList.MarshalJSON() ([]byte, error)
1. GovernanceRuleListIterator.NotDone() bool
1. GovernanceRuleListIterator.Response() GovernanceRuleList
1. GovernanceRuleListIterator.Value() GovernanceRule
1. GovernanceRuleListPage.NotDone() bool
1. GovernanceRuleListPage.Response() GovernanceRuleList
1. GovernanceRuleListPage.Values() []GovernanceRule
1. GovernanceRulesClient.CreateOrUpdate(context.Context, string, GovernanceRule) (GovernanceRule, error)
1. GovernanceRulesClient.CreateOrUpdatePreparer(context.Context, string, GovernanceRule) (*http.Request, error)
1. GovernanceRulesClient.CreateOrUpdateResponder(*http.Response) (GovernanceRule, error)
1. GovernanceRulesClient.CreateOrUpdateSender(*http.Request) (*http.Response, error)
1. GovernanceRulesClient.Delete(context.Context, string) (autorest.Response, error)
1. GovernanceRulesClient.DeletePreparer(context.Context, string) (*http.Request, error)
1. GovernanceRulesClient.DeleteResponder(*http.Response) (autorest.Response, error)
1. GovernanceRulesClient.DeleteSender(*http.Request) (*http.Response, error)
1. GovernanceRulesClient.Get(context.Context, string) (GovernanceRule, error)
1. GovernanceRulesClient.GetPreparer(context.Context, string) (*http.Request, error)
1. GovernanceRulesClient.GetResponder(*http.Response) (GovernanceRule, error)
1. GovernanceRulesClient.GetSender(*http.Request) (*http.Response, error)
1. GovernanceRulesClient.RuleIDExecuteSingleSecurityConnector(context.Context, string, string, string, *ExecuteGovernanceRuleParams) (GovernanceRulesRuleIDExecuteSingleSecurityConnectorFuture, error)
1. GovernanceRulesClient.RuleIDExecuteSingleSecurityConnectorPreparer(context.Context, string, string, string, *ExecuteGovernanceRuleParams) (*http.Request, error)
1. GovernanceRulesClient.RuleIDExecuteSingleSecurityConnectorResponder(*http.Response) (autorest.Response, error)
1. GovernanceRulesClient.RuleIDExecuteSingleSecurityConnectorSender(*http.Request) (GovernanceRulesRuleIDExecuteSingleSecurityConnectorFuture, error)
1. GovernanceRulesClient.RuleIDExecuteSingleSubscription(context.Context, string, *ExecuteGovernanceRuleParams) (GovernanceRulesRuleIDExecuteSingleSubscriptionFuture, error)
1. GovernanceRulesClient.RuleIDExecuteSingleSubscriptionPreparer(context.Context, string, *ExecuteGovernanceRuleParams) (*http.Request, error)
1. GovernanceRulesClient.RuleIDExecuteSingleSubscriptionResponder(*http.Response) (autorest.Response, error)
1. GovernanceRulesClient.RuleIDExecuteSingleSubscriptionSender(*http.Request) (GovernanceRulesRuleIDExecuteSingleSubscriptionFuture, error)
1. InformationProtectionAwsOffering.AsBasicCloudOffering() (BasicCloudOffering, bool)
1. InformationProtectionAwsOffering.AsCloudOffering() (*CloudOffering, bool)
1. InformationProtectionAwsOffering.AsCspmMonitorAwsOffering() (*CspmMonitorAwsOffering, bool)
1. InformationProtectionAwsOffering.AsCspmMonitorAzureDevOpsOffering() (*CspmMonitorAzureDevOpsOffering, bool)
1. InformationProtectionAwsOffering.AsCspmMonitorGcpOffering() (*CspmMonitorGcpOffering, bool)
1. InformationProtectionAwsOffering.AsCspmMonitorGithubOffering() (*CspmMonitorGithubOffering, bool)
1. InformationProtectionAwsOffering.AsDefenderCspmAwsOffering() (*DefenderCspmAwsOffering, bool)
1. InformationProtectionAwsOffering.AsDefenderCspmGcpOffering() (*DefenderCspmGcpOffering, bool)
1. InformationProtectionAwsOffering.AsDefenderFoDatabasesAwsOffering() (*DefenderFoDatabasesAwsOffering, bool)
1. InformationProtectionAwsOffering.AsDefenderForContainersAwsOffering() (*DefenderForContainersAwsOffering, bool)
1. InformationProtectionAwsOffering.AsDefenderForContainersGcpOffering() (*DefenderForContainersGcpOffering, bool)
1. InformationProtectionAwsOffering.AsDefenderForDatabasesGcpOffering() (*DefenderForDatabasesGcpOffering, bool)
1. InformationProtectionAwsOffering.AsDefenderForDevOpsAzureDevOpsOffering() (*DefenderForDevOpsAzureDevOpsOffering, bool)
1. InformationProtectionAwsOffering.AsDefenderForDevOpsGithubOffering() (*DefenderForDevOpsGithubOffering, bool)
1. InformationProtectionAwsOffering.AsDefenderForServersAwsOffering() (*DefenderForServersAwsOffering, bool)
1. InformationProtectionAwsOffering.AsDefenderForServersGcpOffering() (*DefenderForServersGcpOffering, bool)
1. InformationProtectionAwsOffering.AsInformationProtectionAwsOffering() (*InformationProtectionAwsOffering, bool)
1. InformationProtectionAwsOffering.MarshalJSON() ([]byte, error)
1. MdeOnboardingData.MarshalJSON() ([]byte, error)
1. MdeOnboardingsClient.Get(context.Context) (MdeOnboardingData, error)
1. MdeOnboardingsClient.GetPreparer(context.Context) (*http.Request, error)
1. MdeOnboardingsClient.GetResponder(*http.Response) (MdeOnboardingData, error)
1. MdeOnboardingsClient.GetSender(*http.Request) (*http.Response, error)
1. MdeOnboardingsClient.List(context.Context) (MdeOnboardingDataList, error)
1. MdeOnboardingsClient.ListPreparer(context.Context) (*http.Request, error)
1. MdeOnboardingsClient.ListResponder(*http.Response) (MdeOnboardingDataList, error)
1. MdeOnboardingsClient.ListSender(*http.Request) (*http.Response, error)
1. NewApplicationClient(string) ApplicationClient
1. NewApplicationClientWithBaseURI(string, string) ApplicationClient
1. NewApplicationsClient(string) ApplicationsClient
1. NewApplicationsClientWithBaseURI(string, string) ApplicationsClient
1. NewApplicationsListIterator(ApplicationsListPage) ApplicationsListIterator
1. NewApplicationsListPage(ApplicationsList, func(context.Context, ApplicationsList) (ApplicationsList, error)) ApplicationsListPage
1. NewAssessmentMetadataResponseListIterator(AssessmentMetadataResponseListPage) AssessmentMetadataResponseListIterator
1. NewAssessmentMetadataResponseListPage(AssessmentMetadataResponseList, func(context.Context, AssessmentMetadataResponseList) (AssessmentMetadataResponseList, error)) AssessmentMetadataResponseListPage
1. NewConnectorApplicationClient(string) ConnectorApplicationClient
1. NewConnectorApplicationClientWithBaseURI(string, string) ConnectorApplicationClient
1. NewConnectorApplicationsClient(string) ConnectorApplicationsClient
1. NewConnectorApplicationsClientWithBaseURI(string, string) ConnectorApplicationsClient
1. NewConnectorGovernanceRuleClient(string) ConnectorGovernanceRuleClient
1. NewConnectorGovernanceRuleClientWithBaseURI(string, string) ConnectorGovernanceRuleClient
1. NewConnectorGovernanceRulesClient(string) ConnectorGovernanceRulesClient
1. NewConnectorGovernanceRulesClientWithBaseURI(string, string) ConnectorGovernanceRulesClient
1. NewConnectorGovernanceRulesExecuteStatusClient(string) ConnectorGovernanceRulesExecuteStatusClient
1. NewConnectorGovernanceRulesExecuteStatusClientWithBaseURI(string, string) ConnectorGovernanceRulesExecuteStatusClient
1. NewConnectorsGroupClient(string) ConnectorsGroupClient
1. NewConnectorsGroupClientWithBaseURI(string, string) ConnectorsGroupClient
1. NewConnectorsListIterator(ConnectorsListPage) ConnectorsListIterator
1. NewConnectorsListPage(ConnectorsList, func(context.Context, ConnectorsList) (ConnectorsList, error)) ConnectorsListPage
1. NewCustomAssessmentAutomationsClient(string) CustomAssessmentAutomationsClient
1. NewCustomAssessmentAutomationsClientWithBaseURI(string, string) CustomAssessmentAutomationsClient
1. NewCustomAssessmentAutomationsListResultIterator(CustomAssessmentAutomationsListResultPage) CustomAssessmentAutomationsListResultIterator
1. NewCustomAssessmentAutomationsListResultPage(CustomAssessmentAutomationsListResult, func(context.Context, CustomAssessmentAutomationsListResult) (CustomAssessmentAutomationsListResult, error)) CustomAssessmentAutomationsListResultPage
1. NewCustomEntityStoreAssignmentsClient(string) CustomEntityStoreAssignmentsClient
1. NewCustomEntityStoreAssignmentsClientWithBaseURI(string, string) CustomEntityStoreAssignmentsClient
1. NewCustomEntityStoreAssignmentsListResultIterator(CustomEntityStoreAssignmentsListResultPage) CustomEntityStoreAssignmentsListResultIterator
1. NewCustomEntityStoreAssignmentsListResultPage(CustomEntityStoreAssignmentsListResult, func(context.Context, CustomEntityStoreAssignmentsListResult) (CustomEntityStoreAssignmentsListResult, error)) CustomEntityStoreAssignmentsListResultPage
1. NewGovernanceAssignmentsClient(string) GovernanceAssignmentsClient
1. NewGovernanceAssignmentsClientWithBaseURI(string, string) GovernanceAssignmentsClient
1. NewGovernanceAssignmentsListIterator(GovernanceAssignmentsListPage) GovernanceAssignmentsListIterator
1. NewGovernanceAssignmentsListPage(GovernanceAssignmentsList, func(context.Context, GovernanceAssignmentsList) (GovernanceAssignmentsList, error)) GovernanceAssignmentsListPage
1. NewGovernanceRuleClient(string) GovernanceRuleClient
1. NewGovernanceRuleClientWithBaseURI(string, string) GovernanceRuleClient
1. NewGovernanceRuleListIterator(GovernanceRuleListPage) GovernanceRuleListIterator
1. NewGovernanceRuleListPage(GovernanceRuleList, func(context.Context, GovernanceRuleList) (GovernanceRuleList, error)) GovernanceRuleListPage
1. NewGovernanceRulesClient(string) GovernanceRulesClient
1. NewGovernanceRulesClientWithBaseURI(string, string) GovernanceRulesClient
1. NewMdeOnboardingsClient(string) MdeOnboardingsClient
1. NewMdeOnboardingsClientWithBaseURI(string, string) MdeOnboardingsClient
1. NewSubscriptionGovernanceRulesExecuteStatusClient(string) SubscriptionGovernanceRulesExecuteStatusClient
1. NewSubscriptionGovernanceRulesExecuteStatusClientWithBaseURI(string, string) SubscriptionGovernanceRulesExecuteStatusClient
1. PossibleApplicationConditionOperatorValues() []ApplicationConditionOperator
1. PossibleCloudNameValues() []CloudName
1. PossibleEnvironmentTypeValues() []EnvironmentType
1. PossibleGovernanceRuleConditionOperatorValues() []GovernanceRuleConditionOperator
1. PossibleGovernanceRuleOwnerSourceTypeValues() []GovernanceRuleOwnerSourceType
1. PossibleGovernanceRuleTypeValues() []GovernanceRuleType
1. PossibleInformationProtectionPolicyNameValues() []InformationProtectionPolicyName
1. PossibleMinimalSeverityValues() []MinimalSeverity
1. PossibleOfferingTypeValues() []OfferingType
1. PossibleOrganizationMembershipTypeBasicGcpOrganizationalDataValues() []OrganizationMembershipTypeBasicGcpOrganizationalData
1. PossibleOrganizationMembershipTypeValues() []OrganizationMembershipType
1. PossibleRolesValues() []Roles
1. PossibleScanningModeValues() []ScanningMode
1. PossibleSettingName2Values() []SettingName2
1. PossibleSettingName4Values() []SettingName4
1. PossibleSettingName5Values() []SettingName5
1. PossibleSeverityEnumValues() []SeverityEnum
1. PossibleStateForAlertNotificationsValues() []StateForAlertNotifications
1. PossibleStateForNotificationsByRoleValues() []StateForNotificationsByRole
1. PossibleSubPlanValues() []SubPlan
1. PossibleSupportedCloudEnumValues() []SupportedCloudEnum
1. PossibleTacticsValues() []Tactics
1. PossibleTaskUpdateActionTypeValues() []TaskUpdateActionType
1. PossibleTechniquesValues() []Techniques
1. PossibleType1Values() []Type1
1. SubscriptionGovernanceRulesExecuteStatusClient.Get(context.Context, string, string) (SubscriptionGovernanceRulesExecuteStatusGetFuture, error)
1. SubscriptionGovernanceRulesExecuteStatusClient.GetPreparer(context.Context, string, string) (*http.Request, error)
1. SubscriptionGovernanceRulesExecuteStatusClient.GetResponder(*http.Response) (ExecuteRuleStatus, error)
1. SubscriptionGovernanceRulesExecuteStatusClient.GetSender(*http.Request) (SubscriptionGovernanceRulesExecuteStatusGetFuture, error)

### Struct Changes

#### New Structs

1. AlertPropertiesSupportingEvidence
1. Application
1. ApplicationClient
1. ApplicationCondition
1. ApplicationProperties
1. ApplicationsClient
1. ApplicationsList
1. ApplicationsListIterator
1. ApplicationsListPage
1. AssessmentMetadataPropertiesResponse
1. AssessmentMetadataPropertiesResponsePublishDates
1. AssessmentMetadataResponse
1. AssessmentMetadataResponseList
1. AssessmentMetadataResponseListIterator
1. AssessmentMetadataResponseListPage
1. AssessmentPropertiesBase
1. AssessmentPropertiesResponse
1. AssessmentResponse
1. AssessmentStatusResponse
1. AwsEnvironmentData
1. AwsOrganizationalData
1. AwsOrganizationalDataMaster
1. AwsOrganizationalDataMember
1. AzureDevOpsScopeEnvironmentData
1. CloudOffering
1. Condition
1. Connector
1. ConnectorApplicationClient
1. ConnectorApplicationsClient
1. ConnectorGovernanceRuleClient
1. ConnectorGovernanceRulesClient
1. ConnectorGovernanceRulesExecuteStatusClient
1. ConnectorGovernanceRulesExecuteStatusGetFuture
1. ConnectorProperties
1. ConnectorsGroupClient
1. ConnectorsList
1. ConnectorsListIterator
1. ConnectorsListPage
1. ContactPropertiesAlertNotifications
1. ContactPropertiesNotificationsByRole
1. CspmMonitorAwsOffering
1. CspmMonitorAwsOfferingNativeCloudConnection
1. CspmMonitorAzureDevOpsOffering
1. CspmMonitorGcpOffering
1. CspmMonitorGcpOfferingNativeCloudConnection
1. CspmMonitorGithubOffering
1. CustomAssessmentAutomation
1. CustomAssessmentAutomationProperties
1. CustomAssessmentAutomationRequest
1. CustomAssessmentAutomationRequestProperties
1. CustomAssessmentAutomationsClient
1. CustomAssessmentAutomationsListResult
1. CustomAssessmentAutomationsListResultIterator
1. CustomAssessmentAutomationsListResultPage
1. CustomEntityStoreAssignment
1. CustomEntityStoreAssignmentProperties
1. CustomEntityStoreAssignmentRequest
1. CustomEntityStoreAssignmentRequestProperties
1. CustomEntityStoreAssignmentsClient
1. CustomEntityStoreAssignmentsListResult
1. CustomEntityStoreAssignmentsListResultIterator
1. CustomEntityStoreAssignmentsListResultPage
1. DefenderCspmAwsOffering
1. DefenderCspmAwsOfferingVMScanners
1. DefenderCspmAwsOfferingVMScannersConfiguration
1. DefenderCspmGcpOffering
1. DefenderFoDatabasesAwsOffering
1. DefenderFoDatabasesAwsOfferingArcAutoProvisioning
1. DefenderFoDatabasesAwsOfferingRds
1. DefenderForContainersAwsOffering
1. DefenderForContainersAwsOfferingCloudWatchToKinesis
1. DefenderForContainersAwsOfferingContainerVulnerabilityAssessment
1. DefenderForContainersAwsOfferingContainerVulnerabilityAssessmentTask
1. DefenderForContainersAwsOfferingKinesisToS3
1. DefenderForContainersAwsOfferingKubernetesScubaReader
1. DefenderForContainersAwsOfferingKubernetesService
1. DefenderForContainersGcpOffering
1. DefenderForContainersGcpOfferingDataPipelineNativeCloudConnection
1. DefenderForContainersGcpOfferingNativeCloudConnection
1. DefenderForDatabasesGcpOffering
1. DefenderForDatabasesGcpOfferingArcAutoProvisioning
1. DefenderForDatabasesGcpOfferingDefenderForDatabasesArcAutoProvisioning
1. DefenderForDevOpsAzureDevOpsOffering
1. DefenderForDevOpsGithubOffering
1. DefenderForServersAwsOffering
1. DefenderForServersAwsOfferingArcAutoProvisioning
1. DefenderForServersAwsOfferingDefenderForServers
1. DefenderForServersAwsOfferingMdeAutoProvisioning
1. DefenderForServersAwsOfferingSubPlan
1. DefenderForServersAwsOfferingVMScanners
1. DefenderForServersAwsOfferingVMScannersConfiguration
1. DefenderForServersAwsOfferingVaAutoProvisioning
1. DefenderForServersAwsOfferingVaAutoProvisioningConfiguration
1. DefenderForServersGcpOffering
1. DefenderForServersGcpOfferingArcAutoProvisioning
1. DefenderForServersGcpOfferingDefenderForServers
1. DefenderForServersGcpOfferingMdeAutoProvisioning
1. DefenderForServersGcpOfferingSubPlan
1. DefenderForServersGcpOfferingVaAutoProvisioning
1. DefenderForServersGcpOfferingVaAutoProvisioningConfiguration
1. EnvironmentData
1. ExecuteGovernanceRuleParams
1. ExecuteRuleStatus
1. GcpOrganizationalData
1. GcpOrganizationalDataMember
1. GcpOrganizationalDataOrganization
1. GcpProjectDetails
1. GcpProjectEnvironmentData
1. GithubScopeEnvironmentData
1. GovernanceAssignment
1. GovernanceAssignmentAdditionalData
1. GovernanceAssignmentProperties
1. GovernanceAssignmentsClient
1. GovernanceAssignmentsList
1. GovernanceAssignmentsListIterator
1. GovernanceAssignmentsListPage
1. GovernanceEmailNotification
1. GovernanceRule
1. GovernanceRuleClient
1. GovernanceRuleEmailNotification
1. GovernanceRuleList
1. GovernanceRuleListIterator
1. GovernanceRuleListPage
1. GovernanceRuleOwnerSource
1. GovernanceRuleProperties
1. GovernanceRulesClient
1. GovernanceRulesRuleIDExecuteSingleSecurityConnectorFuture
1. GovernanceRulesRuleIDExecuteSingleSubscriptionFuture
1. InformationProtectionAwsOffering
1. InformationProtectionAwsOfferingInformationProtection
1. MdeOnboardingData
1. MdeOnboardingDataList
1. MdeOnboardingDataProperties
1. MdeOnboardingsClient
1. RemediationEta
1. SubscriptionGovernanceRulesExecuteStatusClient
1. SubscriptionGovernanceRulesExecuteStatusGetFuture

#### New Struct Fields

1. AlertProperties.SubTechniques
1. AlertProperties.SupportingEvidence
1. AlertProperties.Techniques
1. AlertProperties.Version
1. ContactProperties.Emails
1. ContactProperties.NotificationsByRole
1. PricingProperties.Deprecated
1. PricingProperties.ReplacedBy
1. PricingProperties.SubPlan
