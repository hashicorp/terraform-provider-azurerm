Generated from https://github.com/Azure/azure-rest-api-specs/tree/3c764635e7d442b3e74caf593029fcd440b3ef82

Code generator @microsoft.azure/autorest.go@~2.1.161

## Breaking Changes

- Function `NewDashboardListResultPage` parameter(s) have been changed from `(func(context.Context, DashboardListResult) (DashboardListResult, error))` to `(DashboardListResult, func(context.Context, DashboardListResult) (DashboardListResult, error))`
- Function `NewResourceProviderOperationListPage` parameter(s) have been changed from `(func(context.Context, ResourceProviderOperationList) (ResourceProviderOperationList, error))` to `(ResourceProviderOperationList, func(context.Context, ResourceProviderOperationList) (ResourceProviderOperationList, error))`
- Type of `ErrorDefinition.Code` has been changed from `*string` to `*int32`

## New Content

- New function `TenantConfigurationsClient.Delete(context.Context) (autorest.Response, error)`
- New function `TenantConfigurationsClient.GetResponder(*http.Response) (Configuration, error)`
- New function `TrackedResource.MarshalJSON() ([]byte, error)`
- New function `NewTenantConfigurationsClient(string) TenantConfigurationsClient`
- New function `TenantConfigurationsClient.Create(context.Context, Configuration) (Configuration, error)`
- New function `TenantConfigurationsClient.ListSender(*http.Request) (*http.Response, error)`
- New function `TenantConfigurationsClient.CreateSender(*http.Request) (*http.Response, error)`
- New function `TenantConfigurationsClient.ListResponder(*http.Response) (ConfigurationList, error)`
- New function `TenantConfigurationsClient.Get(context.Context) (Configuration, error)`
- New function `TenantConfigurationsClient.DeletePreparer(context.Context) (*http.Request, error)`
- New function `Configuration.MarshalJSON() ([]byte, error)`
- New function `NewTenantConfigurationsClientWithBaseURI(string, string) TenantConfigurationsClient`
- New function `TenantConfigurationsClient.GetPreparer(context.Context) (*http.Request, error)`
- New function `TenantConfigurationsClient.ListPreparer(context.Context) (*http.Request, error)`
- New function `TenantConfigurationsClient.GetSender(*http.Request) (*http.Response, error)`
- New function `TenantConfigurationsClient.CreateResponder(*http.Response) (Configuration, error)`
- New function `*Configuration.UnmarshalJSON([]byte) error`
- New function `TenantConfigurationsClient.List(context.Context) (ConfigurationList, error)`
- New function `TenantConfigurationsClient.CreatePreparer(context.Context, Configuration) (*http.Request, error)`
- New function `TenantConfigurationsClient.DeleteSender(*http.Request) (*http.Response, error)`
- New function `TenantConfigurationsClient.DeleteResponder(*http.Response) (autorest.Response, error)`
- New struct `AzureEntityResource`
- New struct `Configuration`
- New struct `ConfigurationList`
- New struct `ConfigurationProperties`
- New struct `ProxyResource`
- New struct `Resource`
- New struct `TenantConfigurationsClient`
- New struct `TrackedResource`
