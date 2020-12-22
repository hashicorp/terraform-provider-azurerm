Generated from https://github.com/Azure/azure-rest-api-specs/tree/3a3a9452f965a227ce43e6b545035b99dd175f23

Code generator @microsoft.azure/autorest.go@~2.1.165

## Breaking Changes

### Removed Constants

1. KindEnum.KindSettingResource

### Removed Funcs

1. DataExportSettings.AsBasicSettingResource() (BasicSettingResource, bool)
1. DataExportSettings.AsSettingResource() (*SettingResource, bool)
1. IotSensor.MarshalJSON() ([]byte, error)
1. Setting.AsBasicSettingResource() (BasicSettingResource, bool)
1. Setting.AsSettingResource() (*SettingResource, bool)
1. SettingResource.AsBasicSetting() (BasicSetting, bool)
1. SettingResource.AsBasicSettingResource() (BasicSettingResource, bool)
1. SettingResource.AsDataExportSettings() (*DataExportSettings, bool)
1. SettingResource.AsSetting() (*Setting, bool)
1. SettingResource.AsSettingResource() (*SettingResource, bool)
1. SettingResource.MarshalJSON() ([]byte, error)

## Struct Changes

### Removed Structs

1. IotSensor
1. SettingResource

## Signature Changes

### Funcs

1. IotSensorsClient.CreateOrUpdate
	- Params
		- From: context.Context, string, string
		- To: context.Context, string, string, IotSensorsModel
	- Returns
		- From: IotSensor, error
		- To: IotSensorsModel, error
1. IotSensorsClient.CreateOrUpdatePreparer
	- Params
		- From: context.Context, string, string
		- To: context.Context, string, string, IotSensorsModel
1. IotSensorsClient.CreateOrUpdateResponder
	- Returns
		- From: IotSensor, error
		- To: IotSensorsModel, error
1. IotSensorsClient.Get
	- Returns
		- From: IotSensor, error
		- To: IotSensorsModel, error
1. IotSensorsClient.GetResponder
	- Returns
		- From: IotSensor, error
		- To: IotSensorsModel, error
1. SettingsClient.Get
	- Returns
		- From: Setting, error
		- To: SettingModel, error
1. SettingsClient.GetResponder
	- Returns
		- From: Setting, error
		- To: SettingModel, error
1. SettingsClient.Update
	- Returns
		- From: Setting, error
		- To: SettingModel, error
1. SettingsClient.UpdateResponder
	- Returns
		- From: Setting, error
		- To: SettingModel, error

### Struct Fields

1. IotSensorsList.Value changed type from *[]IotSensor to *[]IotSensorsModel

### New Constants

1. SensorStatus.Disconnected
1. SensorStatus.Ok
1. SensorStatus.Unavailable
1. TiStatus.TiStatusFailed
1. TiStatus.TiStatusInProgress
1. TiStatus.TiStatusOk
1. TiStatus.TiStatusUpdateAvailable

### New Funcs

1. *IotSensorsModel.UnmarshalJSON([]byte) error
1. *IotSitesModel.UnmarshalJSON([]byte) error
1. *SettingModel.UnmarshalJSON([]byte) error
1. IotSensorProperties.MarshalJSON() ([]byte, error)
1. IotSensorsClient.TriggerTiPackageUpdate(context.Context, string, string) (autorest.Response, error)
1. IotSensorsClient.TriggerTiPackageUpdatePreparer(context.Context, string, string) (*http.Request, error)
1. IotSensorsClient.TriggerTiPackageUpdateResponder(*http.Response) (autorest.Response, error)
1. IotSensorsClient.TriggerTiPackageUpdateSender(*http.Request) (*http.Response, error)
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
1. IotSitesModel.MarshalJSON() ([]byte, error)
1. NewIotSitesClient(string, string) IotSitesClient
1. NewIotSitesClientWithBaseURI(string, string, string) IotSitesClient
1. PossibleSensorStatusValues() []SensorStatus
1. PossibleTiStatusValues() []TiStatus

## Struct Changes

### New Structs

1. IotSensorProperties
1. IotSensorsModel
1. IotSiteProperties
1. IotSitesClient
1. IotSitesList
1. IotSitesModel
1. SettingModel
