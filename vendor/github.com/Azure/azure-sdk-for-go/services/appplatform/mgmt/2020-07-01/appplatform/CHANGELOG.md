Generated from https://github.com/Azure/azure-rest-api-specs/tree/3c764635e7d442b3e74caf593029fcd440b3ef82

Code generator @microsoft.azure/autorest.go@~2.1.161

## Breaking Changes

- Function `NewBindingResourceCollectionPage` parameter(s) have been changed from `(func(context.Context, BindingResourceCollection) (BindingResourceCollection, error))` to `(BindingResourceCollection, func(context.Context, BindingResourceCollection) (BindingResourceCollection, error))`
- Function `NewCertificateResourceCollectionPage` parameter(s) have been changed from `(func(context.Context, CertificateResourceCollection) (CertificateResourceCollection, error))` to `(CertificateResourceCollection, func(context.Context, CertificateResourceCollection) (CertificateResourceCollection, error))`
- Function `NewResourceSkuCollectionPage` parameter(s) have been changed from `(func(context.Context, ResourceSkuCollection) (ResourceSkuCollection, error))` to `(ResourceSkuCollection, func(context.Context, ResourceSkuCollection) (ResourceSkuCollection, error))`
- Function `NewDeploymentResourceCollectionPage` parameter(s) have been changed from `(func(context.Context, DeploymentResourceCollection) (DeploymentResourceCollection, error))` to `(DeploymentResourceCollection, func(context.Context, DeploymentResourceCollection) (DeploymentResourceCollection, error))`
- Function `NewCustomDomainResourceCollectionPage` parameter(s) have been changed from `(func(context.Context, CustomDomainResourceCollection) (CustomDomainResourceCollection, error))` to `(CustomDomainResourceCollection, func(context.Context, CustomDomainResourceCollection) (CustomDomainResourceCollection, error))`
- Function `NewAvailableOperationsPage` parameter(s) have been changed from `(func(context.Context, AvailableOperations) (AvailableOperations, error))` to `(AvailableOperations, func(context.Context, AvailableOperations) (AvailableOperations, error))`
- Function `NewAppResourceCollectionPage` parameter(s) have been changed from `(func(context.Context, AppResourceCollection) (AppResourceCollection, error))` to `(AppResourceCollection, func(context.Context, AppResourceCollection) (AppResourceCollection, error))`
- Function `NewServiceResourceListPage` parameter(s) have been changed from `(func(context.Context, ServiceResourceList) (ServiceResourceList, error))` to `(ServiceResourceList, func(context.Context, ServiceResourceList) (ServiceResourceList, error))`

## New Content

- New const `SupportedRuntimeValueJava11`
- New const `NETCore`
- New const `SupportedRuntimeValueNetCore31`
- New const `NetCoreZip`
- New const `SupportedRuntimeValueJava8`
- New const `Java`
- New const `NetCore31`
- New function `RuntimeVersionsClient.ListRuntimeVersions(context.Context) (AvailableRuntimeVersions, error)`
- New function `PossibleSupportedRuntimeValueValues() []SupportedRuntimeValue`
- New function `RuntimeVersionsClient.ListRuntimeVersionsResponder(*http.Response) (AvailableRuntimeVersions, error)`
- New function `ConfigServersClient.ValidateResponder(*http.Response) (ConfigServerSettingsValidateResult, error)`
- New function `RuntimeVersionsClient.ListRuntimeVersionsPreparer(context.Context) (*http.Request, error)`
- New function `ConfigServersClient.ValidateSender(*http.Request) (ConfigServersValidateFuture, error)`
- New function `NewRuntimeVersionsClient(string) RuntimeVersionsClient`
- New function `ConfigServersClient.Validate(context.Context, string, string, ConfigServerSettings) (ConfigServersValidateFuture, error)`
- New function `NewRuntimeVersionsClientWithBaseURI(string, string) RuntimeVersionsClient`
- New function `ConfigServersClient.ValidatePreparer(context.Context, string, string, ConfigServerSettings) (*http.Request, error)`
- New function `NetworkProfile.MarshalJSON() ([]byte, error)`
- New function `PossibleSupportedRuntimePlatformValues() []SupportedRuntimePlatform`
- New function `RuntimeVersionsClient.ListRuntimeVersionsSender(*http.Request) (*http.Response, error)`
- New function `*ConfigServersValidateFuture.Result(ConfigServersClient) (ConfigServerSettingsValidateResult, error)`
- New struct `AvailableRuntimeVersions`
- New struct `ConfigServerSettingsErrorRecord`
- New struct `ConfigServerSettingsValidateResult`
- New struct `ConfigServersValidateFuture`
- New struct `NetworkProfileOutboundIPs`
- New struct `RuntimeVersionsClient`
- New struct `SupportedRuntimeVersion`
- New field `NetCoreMainEntryPath` in struct `DeploymentSettings`
- New field `StartTime` in struct `DeploymentInstance`
- New field `OutboundIPs` in struct `NetworkProfile`
