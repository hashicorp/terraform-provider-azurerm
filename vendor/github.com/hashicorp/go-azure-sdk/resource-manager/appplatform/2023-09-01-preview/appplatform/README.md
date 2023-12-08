
## `github.com/hashicorp/go-azure-sdk/resource-manager/appplatform/2023-09-01-preview/appplatform` Documentation

The `appplatform` SDK allows for interaction with the Azure Resource Manager Service `appplatform` (API Version `2023-09-01-preview`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/appplatform/2023-09-01-preview/appplatform"
```


### Client Initialization

```go
client := appplatform.NewAppPlatformClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `AppPlatformClient.ApiPortalCustomDomainsCreateOrUpdate`

```go
ctx := context.TODO()
id := appplatform.NewApiPortalDomainID("12345678-1234-9876-4563-123456789012", "example-resource-group", "springValue", "apiPortalValue", "domainValue")

payload := appplatform.ApiPortalCustomDomainResource{
	// ...
}


if err := client.ApiPortalCustomDomainsCreateOrUpdateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `AppPlatformClient.ApiPortalCustomDomainsDelete`

```go
ctx := context.TODO()
id := appplatform.NewApiPortalDomainID("12345678-1234-9876-4563-123456789012", "example-resource-group", "springValue", "apiPortalValue", "domainValue")

if err := client.ApiPortalCustomDomainsDeleteThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `AppPlatformClient.ApiPortalCustomDomainsGet`

```go
ctx := context.TODO()
id := appplatform.NewApiPortalDomainID("12345678-1234-9876-4563-123456789012", "example-resource-group", "springValue", "apiPortalValue", "domainValue")

read, err := client.ApiPortalCustomDomainsGet(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `AppPlatformClient.ApiPortalCustomDomainsList`

```go
ctx := context.TODO()
id := appplatform.NewApiPortalID("12345678-1234-9876-4563-123456789012", "example-resource-group", "springValue", "apiPortalValue")

// alternatively `client.ApiPortalCustomDomainsList(ctx, id)` can be used to do batched pagination
items, err := client.ApiPortalCustomDomainsListComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `AppPlatformClient.ApiPortalsCreateOrUpdate`

```go
ctx := context.TODO()
id := appplatform.NewApiPortalID("12345678-1234-9876-4563-123456789012", "example-resource-group", "springValue", "apiPortalValue")

payload := appplatform.ApiPortalResource{
	// ...
}


if err := client.ApiPortalsCreateOrUpdateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `AppPlatformClient.ApiPortalsDelete`

```go
ctx := context.TODO()
id := appplatform.NewApiPortalID("12345678-1234-9876-4563-123456789012", "example-resource-group", "springValue", "apiPortalValue")

if err := client.ApiPortalsDeleteThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `AppPlatformClient.ApiPortalsGet`

```go
ctx := context.TODO()
id := appplatform.NewApiPortalID("12345678-1234-9876-4563-123456789012", "example-resource-group", "springValue", "apiPortalValue")

read, err := client.ApiPortalsGet(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `AppPlatformClient.ApiPortalsList`

```go
ctx := context.TODO()
id := appplatform.NewSpringID("12345678-1234-9876-4563-123456789012", "example-resource-group", "springValue")

// alternatively `client.ApiPortalsList(ctx, id)` can be used to do batched pagination
items, err := client.ApiPortalsListComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `AppPlatformClient.ApiPortalsValidateDomain`

```go
ctx := context.TODO()
id := appplatform.NewApiPortalID("12345678-1234-9876-4563-123456789012", "example-resource-group", "springValue", "apiPortalValue")

payload := appplatform.CustomDomainValidatePayload{
	// ...
}


read, err := client.ApiPortalsValidateDomain(ctx, id, payload)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `AppPlatformClient.ApmsCreateOrUpdate`

```go
ctx := context.TODO()
id := appplatform.NewApmID("12345678-1234-9876-4563-123456789012", "example-resource-group", "springValue", "apmValue")

payload := appplatform.ApmResource{
	// ...
}


if err := client.ApmsCreateOrUpdateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `AppPlatformClient.ApmsDelete`

```go
ctx := context.TODO()
id := appplatform.NewApmID("12345678-1234-9876-4563-123456789012", "example-resource-group", "springValue", "apmValue")

if err := client.ApmsDeleteThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `AppPlatformClient.ApmsGet`

```go
ctx := context.TODO()
id := appplatform.NewApmID("12345678-1234-9876-4563-123456789012", "example-resource-group", "springValue", "apmValue")

read, err := client.ApmsGet(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `AppPlatformClient.ApmsList`

```go
ctx := context.TODO()
id := appplatform.NewSpringID("12345678-1234-9876-4563-123456789012", "example-resource-group", "springValue")

// alternatively `client.ApmsList(ctx, id)` can be used to do batched pagination
items, err := client.ApmsListComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `AppPlatformClient.ApmsListSecretKeys`

```go
ctx := context.TODO()
id := appplatform.NewApmID("12345678-1234-9876-4563-123456789012", "example-resource-group", "springValue", "apmValue")

read, err := client.ApmsListSecretKeys(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `AppPlatformClient.ApplicationAcceleratorsCreateOrUpdate`

```go
ctx := context.TODO()
id := appplatform.NewApplicationAcceleratorID("12345678-1234-9876-4563-123456789012", "example-resource-group", "springValue", "applicationAcceleratorValue")

payload := appplatform.ApplicationAcceleratorResource{
	// ...
}


if err := client.ApplicationAcceleratorsCreateOrUpdateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `AppPlatformClient.ApplicationAcceleratorsDelete`

```go
ctx := context.TODO()
id := appplatform.NewApplicationAcceleratorID("12345678-1234-9876-4563-123456789012", "example-resource-group", "springValue", "applicationAcceleratorValue")

if err := client.ApplicationAcceleratorsDeleteThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `AppPlatformClient.ApplicationAcceleratorsGet`

```go
ctx := context.TODO()
id := appplatform.NewApplicationAcceleratorID("12345678-1234-9876-4563-123456789012", "example-resource-group", "springValue", "applicationAcceleratorValue")

read, err := client.ApplicationAcceleratorsGet(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `AppPlatformClient.ApplicationAcceleratorsList`

```go
ctx := context.TODO()
id := appplatform.NewSpringID("12345678-1234-9876-4563-123456789012", "example-resource-group", "springValue")

// alternatively `client.ApplicationAcceleratorsList(ctx, id)` can be used to do batched pagination
items, err := client.ApplicationAcceleratorsListComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `AppPlatformClient.ApplicationLiveViewsCreateOrUpdate`

```go
ctx := context.TODO()
id := appplatform.NewApplicationLiveViewID("12345678-1234-9876-4563-123456789012", "example-resource-group", "springValue", "applicationLiveViewValue")

payload := appplatform.ApplicationLiveViewResource{
	// ...
}


if err := client.ApplicationLiveViewsCreateOrUpdateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `AppPlatformClient.ApplicationLiveViewsDelete`

```go
ctx := context.TODO()
id := appplatform.NewApplicationLiveViewID("12345678-1234-9876-4563-123456789012", "example-resource-group", "springValue", "applicationLiveViewValue")

if err := client.ApplicationLiveViewsDeleteThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `AppPlatformClient.ApplicationLiveViewsGet`

```go
ctx := context.TODO()
id := appplatform.NewApplicationLiveViewID("12345678-1234-9876-4563-123456789012", "example-resource-group", "springValue", "applicationLiveViewValue")

read, err := client.ApplicationLiveViewsGet(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `AppPlatformClient.ApplicationLiveViewsList`

```go
ctx := context.TODO()
id := appplatform.NewSpringID("12345678-1234-9876-4563-123456789012", "example-resource-group", "springValue")

// alternatively `client.ApplicationLiveViewsList(ctx, id)` can be used to do batched pagination
items, err := client.ApplicationLiveViewsListComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `AppPlatformClient.AppsCreateOrUpdate`

```go
ctx := context.TODO()
id := appplatform.NewAppID("12345678-1234-9876-4563-123456789012", "example-resource-group", "springValue", "appValue")

payload := appplatform.AppResource{
	// ...
}


if err := client.AppsCreateOrUpdateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `AppPlatformClient.AppsDelete`

```go
ctx := context.TODO()
id := appplatform.NewAppID("12345678-1234-9876-4563-123456789012", "example-resource-group", "springValue", "appValue")

if err := client.AppsDeleteThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `AppPlatformClient.AppsGet`

```go
ctx := context.TODO()
id := appplatform.NewAppID("12345678-1234-9876-4563-123456789012", "example-resource-group", "springValue", "appValue")

read, err := client.AppsGet(ctx, id, appplatform.DefaultAppsGetOperationOptions())
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `AppPlatformClient.AppsGetResourceUploadUrl`

```go
ctx := context.TODO()
id := appplatform.NewAppID("12345678-1234-9876-4563-123456789012", "example-resource-group", "springValue", "appValue")

read, err := client.AppsGetResourceUploadUrl(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `AppPlatformClient.AppsList`

```go
ctx := context.TODO()
id := appplatform.NewSpringID("12345678-1234-9876-4563-123456789012", "example-resource-group", "springValue")

// alternatively `client.AppsList(ctx, id)` can be used to do batched pagination
items, err := client.AppsListComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `AppPlatformClient.AppsSetActiveDeployments`

```go
ctx := context.TODO()
id := appplatform.NewAppID("12345678-1234-9876-4563-123456789012", "example-resource-group", "springValue", "appValue")

payload := appplatform.ActiveDeploymentCollection{
	// ...
}


if err := client.AppsSetActiveDeploymentsThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `AppPlatformClient.AppsUpdate`

```go
ctx := context.TODO()
id := appplatform.NewAppID("12345678-1234-9876-4563-123456789012", "example-resource-group", "springValue", "appValue")

payload := appplatform.AppResource{
	// ...
}


if err := client.AppsUpdateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `AppPlatformClient.AppsValidateDomain`

```go
ctx := context.TODO()
id := appplatform.NewAppID("12345678-1234-9876-4563-123456789012", "example-resource-group", "springValue", "appValue")

payload := appplatform.CustomDomainValidatePayload{
	// ...
}


read, err := client.AppsValidateDomain(ctx, id, payload)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `AppPlatformClient.BindingsCreateOrUpdate`

```go
ctx := context.TODO()
id := appplatform.NewBindingID("12345678-1234-9876-4563-123456789012", "example-resource-group", "springValue", "appValue", "bindingValue")

payload := appplatform.BindingResource{
	// ...
}


if err := client.BindingsCreateOrUpdateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `AppPlatformClient.BindingsDelete`

```go
ctx := context.TODO()
id := appplatform.NewBindingID("12345678-1234-9876-4563-123456789012", "example-resource-group", "springValue", "appValue", "bindingValue")

if err := client.BindingsDeleteThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `AppPlatformClient.BindingsGet`

```go
ctx := context.TODO()
id := appplatform.NewBindingID("12345678-1234-9876-4563-123456789012", "example-resource-group", "springValue", "appValue", "bindingValue")

read, err := client.BindingsGet(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `AppPlatformClient.BindingsList`

```go
ctx := context.TODO()
id := appplatform.NewAppID("12345678-1234-9876-4563-123456789012", "example-resource-group", "springValue", "appValue")

// alternatively `client.BindingsList(ctx, id)` can be used to do batched pagination
items, err := client.BindingsListComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `AppPlatformClient.BindingsUpdate`

```go
ctx := context.TODO()
id := appplatform.NewBindingID("12345678-1234-9876-4563-123456789012", "example-resource-group", "springValue", "appValue", "bindingValue")

payload := appplatform.BindingResource{
	// ...
}


if err := client.BindingsUpdateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `AppPlatformClient.BuildServiceAgentPoolGet`

```go
ctx := context.TODO()
id := appplatform.NewAgentPoolID("12345678-1234-9876-4563-123456789012", "example-resource-group", "springValue", "buildServiceValue", "agentPoolValue")

read, err := client.BuildServiceAgentPoolGet(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `AppPlatformClient.BuildServiceAgentPoolList`

```go
ctx := context.TODO()
id := appplatform.NewBuildServiceID("12345678-1234-9876-4563-123456789012", "example-resource-group", "springValue", "buildServiceValue")

// alternatively `client.BuildServiceAgentPoolList(ctx, id)` can be used to do batched pagination
items, err := client.BuildServiceAgentPoolListComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `AppPlatformClient.BuildServiceAgentPoolUpdatePut`

```go
ctx := context.TODO()
id := appplatform.NewAgentPoolID("12345678-1234-9876-4563-123456789012", "example-resource-group", "springValue", "buildServiceValue", "agentPoolValue")

payload := appplatform.BuildServiceAgentPoolResource{
	// ...
}


if err := client.BuildServiceAgentPoolUpdatePutThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `AppPlatformClient.BuildServiceBuilderCreateOrUpdate`

```go
ctx := context.TODO()
id := appplatform.NewBuilderID("12345678-1234-9876-4563-123456789012", "example-resource-group", "springValue", "buildServiceValue", "builderValue")

payload := appplatform.BuilderResource{
	// ...
}


if err := client.BuildServiceBuilderCreateOrUpdateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `AppPlatformClient.BuildServiceBuilderDelete`

```go
ctx := context.TODO()
id := appplatform.NewBuilderID("12345678-1234-9876-4563-123456789012", "example-resource-group", "springValue", "buildServiceValue", "builderValue")

if err := client.BuildServiceBuilderDeleteThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `AppPlatformClient.BuildServiceBuilderGet`

```go
ctx := context.TODO()
id := appplatform.NewBuilderID("12345678-1234-9876-4563-123456789012", "example-resource-group", "springValue", "buildServiceValue", "builderValue")

read, err := client.BuildServiceBuilderGet(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `AppPlatformClient.BuildServiceBuilderList`

```go
ctx := context.TODO()
id := appplatform.NewBuildServiceID("12345678-1234-9876-4563-123456789012", "example-resource-group", "springValue", "buildServiceValue")

// alternatively `client.BuildServiceBuilderList(ctx, id)` can be used to do batched pagination
items, err := client.BuildServiceBuilderListComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `AppPlatformClient.BuildServiceBuilderListDeployments`

```go
ctx := context.TODO()
id := appplatform.NewBuilderID("12345678-1234-9876-4563-123456789012", "example-resource-group", "springValue", "buildServiceValue", "builderValue")

read, err := client.BuildServiceBuilderListDeployments(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `AppPlatformClient.BuildServiceCreateOrUpdate`

```go
ctx := context.TODO()
id := appplatform.NewBuildServiceID("12345678-1234-9876-4563-123456789012", "example-resource-group", "springValue", "buildServiceValue")

payload := appplatform.BuildService{
	// ...
}


if err := client.BuildServiceCreateOrUpdateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `AppPlatformClient.BuildServiceCreateOrUpdateBuild`

```go
ctx := context.TODO()
id := appplatform.NewBuildID("12345678-1234-9876-4563-123456789012", "example-resource-group", "springValue", "buildServiceValue", "buildValue")

payload := appplatform.Build{
	// ...
}


read, err := client.BuildServiceCreateOrUpdateBuild(ctx, id, payload)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `AppPlatformClient.BuildServiceDeleteBuild`

```go
ctx := context.TODO()
id := appplatform.NewBuildID("12345678-1234-9876-4563-123456789012", "example-resource-group", "springValue", "buildServiceValue", "buildValue")

if err := client.BuildServiceDeleteBuildThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `AppPlatformClient.BuildServiceGetBuild`

```go
ctx := context.TODO()
id := appplatform.NewBuildID("12345678-1234-9876-4563-123456789012", "example-resource-group", "springValue", "buildServiceValue", "buildValue")

read, err := client.BuildServiceGetBuild(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `AppPlatformClient.BuildServiceGetBuildResult`

```go
ctx := context.TODO()
id := appplatform.NewResultID("12345678-1234-9876-4563-123456789012", "example-resource-group", "springValue", "buildServiceValue", "buildValue", "resultValue")

read, err := client.BuildServiceGetBuildResult(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `AppPlatformClient.BuildServiceGetBuildResultLog`

```go
ctx := context.TODO()
id := appplatform.NewResultID("12345678-1234-9876-4563-123456789012", "example-resource-group", "springValue", "buildServiceValue", "buildValue", "resultValue")

read, err := client.BuildServiceGetBuildResultLog(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `AppPlatformClient.BuildServiceGetBuildService`

```go
ctx := context.TODO()
id := appplatform.NewBuildServiceID("12345678-1234-9876-4563-123456789012", "example-resource-group", "springValue", "buildServiceValue")

read, err := client.BuildServiceGetBuildService(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `AppPlatformClient.BuildServiceGetResourceUploadUrl`

```go
ctx := context.TODO()
id := appplatform.NewBuildServiceID("12345678-1234-9876-4563-123456789012", "example-resource-group", "springValue", "buildServiceValue")

read, err := client.BuildServiceGetResourceUploadUrl(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `AppPlatformClient.BuildServiceGetSupportedBuildpack`

```go
ctx := context.TODO()
id := appplatform.NewSupportedBuildPackID("12345678-1234-9876-4563-123456789012", "example-resource-group", "springValue", "buildServiceValue", "supportedBuildPackValue")

read, err := client.BuildServiceGetSupportedBuildpack(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `AppPlatformClient.BuildServiceGetSupportedStack`

```go
ctx := context.TODO()
id := appplatform.NewSupportedStackID("12345678-1234-9876-4563-123456789012", "example-resource-group", "springValue", "buildServiceValue", "supportedStackValue")

read, err := client.BuildServiceGetSupportedStack(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `AppPlatformClient.BuildServiceListBuildResults`

```go
ctx := context.TODO()
id := appplatform.NewBuildID("12345678-1234-9876-4563-123456789012", "example-resource-group", "springValue", "buildServiceValue", "buildValue")

// alternatively `client.BuildServiceListBuildResults(ctx, id)` can be used to do batched pagination
items, err := client.BuildServiceListBuildResultsComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `AppPlatformClient.BuildServiceListBuildServices`

```go
ctx := context.TODO()
id := appplatform.NewSpringID("12345678-1234-9876-4563-123456789012", "example-resource-group", "springValue")

// alternatively `client.BuildServiceListBuildServices(ctx, id)` can be used to do batched pagination
items, err := client.BuildServiceListBuildServicesComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `AppPlatformClient.BuildServiceListBuilds`

```go
ctx := context.TODO()
id := appplatform.NewBuildServiceID("12345678-1234-9876-4563-123456789012", "example-resource-group", "springValue", "buildServiceValue")

// alternatively `client.BuildServiceListBuilds(ctx, id)` can be used to do batched pagination
items, err := client.BuildServiceListBuildsComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `AppPlatformClient.BuildServiceListSupportedBuildpacks`

```go
ctx := context.TODO()
id := appplatform.NewBuildServiceID("12345678-1234-9876-4563-123456789012", "example-resource-group", "springValue", "buildServiceValue")

read, err := client.BuildServiceListSupportedBuildpacks(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `AppPlatformClient.BuildServiceListSupportedStacks`

```go
ctx := context.TODO()
id := appplatform.NewBuildServiceID("12345678-1234-9876-4563-123456789012", "example-resource-group", "springValue", "buildServiceValue")

read, err := client.BuildServiceListSupportedStacks(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `AppPlatformClient.BuildpackBindingCreateOrUpdate`

```go
ctx := context.TODO()
id := appplatform.NewBuildPackBindingID("12345678-1234-9876-4563-123456789012", "example-resource-group", "springValue", "buildServiceValue", "builderValue", "buildPackBindingValue")

payload := appplatform.BuildpackBindingResource{
	// ...
}


if err := client.BuildpackBindingCreateOrUpdateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `AppPlatformClient.BuildpackBindingDelete`

```go
ctx := context.TODO()
id := appplatform.NewBuildPackBindingID("12345678-1234-9876-4563-123456789012", "example-resource-group", "springValue", "buildServiceValue", "builderValue", "buildPackBindingValue")

if err := client.BuildpackBindingDeleteThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `AppPlatformClient.BuildpackBindingGet`

```go
ctx := context.TODO()
id := appplatform.NewBuildPackBindingID("12345678-1234-9876-4563-123456789012", "example-resource-group", "springValue", "buildServiceValue", "builderValue", "buildPackBindingValue")

read, err := client.BuildpackBindingGet(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `AppPlatformClient.BuildpackBindingList`

```go
ctx := context.TODO()
id := appplatform.NewBuilderID("12345678-1234-9876-4563-123456789012", "example-resource-group", "springValue", "buildServiceValue", "builderValue")

// alternatively `client.BuildpackBindingList(ctx, id)` can be used to do batched pagination
items, err := client.BuildpackBindingListComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `AppPlatformClient.BuildpackBindingListForCluster`

```go
ctx := context.TODO()
id := appplatform.NewSpringID("12345678-1234-9876-4563-123456789012", "example-resource-group", "springValue")

// alternatively `client.BuildpackBindingListForCluster(ctx, id)` can be used to do batched pagination
items, err := client.BuildpackBindingListForClusterComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `AppPlatformClient.CertificatesCreateOrUpdate`

```go
ctx := context.TODO()
id := appplatform.NewCertificateID("12345678-1234-9876-4563-123456789012", "example-resource-group", "springValue", "certificateValue")

payload := appplatform.CertificateResource{
	// ...
}


if err := client.CertificatesCreateOrUpdateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `AppPlatformClient.CertificatesDelete`

```go
ctx := context.TODO()
id := appplatform.NewCertificateID("12345678-1234-9876-4563-123456789012", "example-resource-group", "springValue", "certificateValue")

if err := client.CertificatesDeleteThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `AppPlatformClient.CertificatesGet`

```go
ctx := context.TODO()
id := appplatform.NewCertificateID("12345678-1234-9876-4563-123456789012", "example-resource-group", "springValue", "certificateValue")

read, err := client.CertificatesGet(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `AppPlatformClient.CertificatesList`

```go
ctx := context.TODO()
id := appplatform.NewSpringID("12345678-1234-9876-4563-123456789012", "example-resource-group", "springValue")

// alternatively `client.CertificatesList(ctx, id)` can be used to do batched pagination
items, err := client.CertificatesListComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `AppPlatformClient.ConfigServersGet`

```go
ctx := context.TODO()
id := appplatform.NewSpringID("12345678-1234-9876-4563-123456789012", "example-resource-group", "springValue")

read, err := client.ConfigServersGet(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `AppPlatformClient.ConfigServersUpdatePatch`

```go
ctx := context.TODO()
id := appplatform.NewSpringID("12345678-1234-9876-4563-123456789012", "example-resource-group", "springValue")

payload := appplatform.ConfigServerResource{
	// ...
}


if err := client.ConfigServersUpdatePatchThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `AppPlatformClient.ConfigServersUpdatePut`

```go
ctx := context.TODO()
id := appplatform.NewSpringID("12345678-1234-9876-4563-123456789012", "example-resource-group", "springValue")

payload := appplatform.ConfigServerResource{
	// ...
}


if err := client.ConfigServersUpdatePutThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `AppPlatformClient.ConfigServersValidate`

```go
ctx := context.TODO()
id := appplatform.NewSpringID("12345678-1234-9876-4563-123456789012", "example-resource-group", "springValue")

payload := appplatform.ConfigServerSettings{
	// ...
}


if err := client.ConfigServersValidateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `AppPlatformClient.ConfigurationServicesCreateOrUpdate`

```go
ctx := context.TODO()
id := appplatform.NewConfigurationServiceID("12345678-1234-9876-4563-123456789012", "example-resource-group", "springValue", "configurationServiceValue")

payload := appplatform.ConfigurationServiceResource{
	// ...
}


if err := client.ConfigurationServicesCreateOrUpdateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `AppPlatformClient.ConfigurationServicesDelete`

```go
ctx := context.TODO()
id := appplatform.NewConfigurationServiceID("12345678-1234-9876-4563-123456789012", "example-resource-group", "springValue", "configurationServiceValue")

if err := client.ConfigurationServicesDeleteThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `AppPlatformClient.ConfigurationServicesGet`

```go
ctx := context.TODO()
id := appplatform.NewConfigurationServiceID("12345678-1234-9876-4563-123456789012", "example-resource-group", "springValue", "configurationServiceValue")

read, err := client.ConfigurationServicesGet(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `AppPlatformClient.ConfigurationServicesList`

```go
ctx := context.TODO()
id := appplatform.NewSpringID("12345678-1234-9876-4563-123456789012", "example-resource-group", "springValue")

// alternatively `client.ConfigurationServicesList(ctx, id)` can be used to do batched pagination
items, err := client.ConfigurationServicesListComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `AppPlatformClient.ConfigurationServicesValidate`

```go
ctx := context.TODO()
id := appplatform.NewConfigurationServiceID("12345678-1234-9876-4563-123456789012", "example-resource-group", "springValue", "configurationServiceValue")

payload := appplatform.ConfigurationServiceSettings{
	// ...
}


if err := client.ConfigurationServicesValidateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `AppPlatformClient.ConfigurationServicesValidateResource`

```go
ctx := context.TODO()
id := appplatform.NewConfigurationServiceID("12345678-1234-9876-4563-123456789012", "example-resource-group", "springValue", "configurationServiceValue")

payload := appplatform.ConfigurationServiceResource{
	// ...
}


if err := client.ConfigurationServicesValidateResourceThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `AppPlatformClient.ContainerRegistriesCreateOrUpdate`

```go
ctx := context.TODO()
id := appplatform.NewContainerRegistryID("12345678-1234-9876-4563-123456789012", "example-resource-group", "springValue", "containerRegistryValue")

payload := appplatform.ContainerRegistryResource{
	// ...
}


if err := client.ContainerRegistriesCreateOrUpdateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `AppPlatformClient.ContainerRegistriesDelete`

```go
ctx := context.TODO()
id := appplatform.NewContainerRegistryID("12345678-1234-9876-4563-123456789012", "example-resource-group", "springValue", "containerRegistryValue")

if err := client.ContainerRegistriesDeleteThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `AppPlatformClient.ContainerRegistriesGet`

```go
ctx := context.TODO()
id := appplatform.NewContainerRegistryID("12345678-1234-9876-4563-123456789012", "example-resource-group", "springValue", "containerRegistryValue")

read, err := client.ContainerRegistriesGet(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `AppPlatformClient.ContainerRegistriesList`

```go
ctx := context.TODO()
id := appplatform.NewSpringID("12345678-1234-9876-4563-123456789012", "example-resource-group", "springValue")

// alternatively `client.ContainerRegistriesList(ctx, id)` can be used to do batched pagination
items, err := client.ContainerRegistriesListComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `AppPlatformClient.ContainerRegistriesValidate`

```go
ctx := context.TODO()
id := appplatform.NewContainerRegistryID("12345678-1234-9876-4563-123456789012", "example-resource-group", "springValue", "containerRegistryValue")

payload := appplatform.ContainerRegistryProperties{
	// ...
}


if err := client.ContainerRegistriesValidateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `AppPlatformClient.CustomDomainsCreateOrUpdate`

```go
ctx := context.TODO()
id := appplatform.NewDomainID("12345678-1234-9876-4563-123456789012", "example-resource-group", "springValue", "appValue", "domainValue")

payload := appplatform.CustomDomainResource{
	// ...
}


if err := client.CustomDomainsCreateOrUpdateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `AppPlatformClient.CustomDomainsDelete`

```go
ctx := context.TODO()
id := appplatform.NewDomainID("12345678-1234-9876-4563-123456789012", "example-resource-group", "springValue", "appValue", "domainValue")

if err := client.CustomDomainsDeleteThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `AppPlatformClient.CustomDomainsGet`

```go
ctx := context.TODO()
id := appplatform.NewDomainID("12345678-1234-9876-4563-123456789012", "example-resource-group", "springValue", "appValue", "domainValue")

read, err := client.CustomDomainsGet(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `AppPlatformClient.CustomDomainsList`

```go
ctx := context.TODO()
id := appplatform.NewAppID("12345678-1234-9876-4563-123456789012", "example-resource-group", "springValue", "appValue")

// alternatively `client.CustomDomainsList(ctx, id)` can be used to do batched pagination
items, err := client.CustomDomainsListComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `AppPlatformClient.CustomDomainsUpdate`

```go
ctx := context.TODO()
id := appplatform.NewDomainID("12345678-1234-9876-4563-123456789012", "example-resource-group", "springValue", "appValue", "domainValue")

payload := appplatform.CustomDomainResource{
	// ...
}


if err := client.CustomDomainsUpdateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `AppPlatformClient.CustomizedAcceleratorsCreateOrUpdate`

```go
ctx := context.TODO()
id := appplatform.NewCustomizedAcceleratorID("12345678-1234-9876-4563-123456789012", "example-resource-group", "springValue", "applicationAcceleratorValue", "customizedAcceleratorValue")

payload := appplatform.CustomizedAcceleratorResource{
	// ...
}


if err := client.CustomizedAcceleratorsCreateOrUpdateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `AppPlatformClient.CustomizedAcceleratorsDelete`

```go
ctx := context.TODO()
id := appplatform.NewCustomizedAcceleratorID("12345678-1234-9876-4563-123456789012", "example-resource-group", "springValue", "applicationAcceleratorValue", "customizedAcceleratorValue")

if err := client.CustomizedAcceleratorsDeleteThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `AppPlatformClient.CustomizedAcceleratorsGet`

```go
ctx := context.TODO()
id := appplatform.NewCustomizedAcceleratorID("12345678-1234-9876-4563-123456789012", "example-resource-group", "springValue", "applicationAcceleratorValue", "customizedAcceleratorValue")

read, err := client.CustomizedAcceleratorsGet(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `AppPlatformClient.CustomizedAcceleratorsList`

```go
ctx := context.TODO()
id := appplatform.NewApplicationAcceleratorID("12345678-1234-9876-4563-123456789012", "example-resource-group", "springValue", "applicationAcceleratorValue")

// alternatively `client.CustomizedAcceleratorsList(ctx, id)` can be used to do batched pagination
items, err := client.CustomizedAcceleratorsListComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `AppPlatformClient.CustomizedAcceleratorsValidate`

```go
ctx := context.TODO()
id := appplatform.NewCustomizedAcceleratorID("12345678-1234-9876-4563-123456789012", "example-resource-group", "springValue", "applicationAcceleratorValue", "customizedAcceleratorValue")

payload := appplatform.CustomizedAcceleratorProperties{
	// ...
}


read, err := client.CustomizedAcceleratorsValidate(ctx, id, payload)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `AppPlatformClient.DeploymentsCreateOrUpdate`

```go
ctx := context.TODO()
id := appplatform.NewDeploymentID("12345678-1234-9876-4563-123456789012", "example-resource-group", "springValue", "appValue", "deploymentValue")

payload := appplatform.DeploymentResource{
	// ...
}


if err := client.DeploymentsCreateOrUpdateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `AppPlatformClient.DeploymentsDelete`

```go
ctx := context.TODO()
id := appplatform.NewDeploymentID("12345678-1234-9876-4563-123456789012", "example-resource-group", "springValue", "appValue", "deploymentValue")

if err := client.DeploymentsDeleteThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `AppPlatformClient.DeploymentsDisableRemoteDebugging`

```go
ctx := context.TODO()
id := appplatform.NewDeploymentID("12345678-1234-9876-4563-123456789012", "example-resource-group", "springValue", "appValue", "deploymentValue")

if err := client.DeploymentsDisableRemoteDebuggingThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `AppPlatformClient.DeploymentsEnableRemoteDebugging`

```go
ctx := context.TODO()
id := appplatform.NewDeploymentID("12345678-1234-9876-4563-123456789012", "example-resource-group", "springValue", "appValue", "deploymentValue")

payload := appplatform.RemoteDebuggingPayload{
	// ...
}


if err := client.DeploymentsEnableRemoteDebuggingThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `AppPlatformClient.DeploymentsGenerateHeapDump`

```go
ctx := context.TODO()
id := appplatform.NewDeploymentID("12345678-1234-9876-4563-123456789012", "example-resource-group", "springValue", "appValue", "deploymentValue")

payload := appplatform.DiagnosticParameters{
	// ...
}


if err := client.DeploymentsGenerateHeapDumpThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `AppPlatformClient.DeploymentsGenerateThreadDump`

```go
ctx := context.TODO()
id := appplatform.NewDeploymentID("12345678-1234-9876-4563-123456789012", "example-resource-group", "springValue", "appValue", "deploymentValue")

payload := appplatform.DiagnosticParameters{
	// ...
}


if err := client.DeploymentsGenerateThreadDumpThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `AppPlatformClient.DeploymentsGet`

```go
ctx := context.TODO()
id := appplatform.NewDeploymentID("12345678-1234-9876-4563-123456789012", "example-resource-group", "springValue", "appValue", "deploymentValue")

read, err := client.DeploymentsGet(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `AppPlatformClient.DeploymentsGetLogFileUrl`

```go
ctx := context.TODO()
id := appplatform.NewDeploymentID("12345678-1234-9876-4563-123456789012", "example-resource-group", "springValue", "appValue", "deploymentValue")

read, err := client.DeploymentsGetLogFileUrl(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `AppPlatformClient.DeploymentsGetRemoteDebuggingConfig`

```go
ctx := context.TODO()
id := appplatform.NewDeploymentID("12345678-1234-9876-4563-123456789012", "example-resource-group", "springValue", "appValue", "deploymentValue")

read, err := client.DeploymentsGetRemoteDebuggingConfig(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `AppPlatformClient.DeploymentsList`

```go
ctx := context.TODO()
id := appplatform.NewAppID("12345678-1234-9876-4563-123456789012", "example-resource-group", "springValue", "appValue")

// alternatively `client.DeploymentsList(ctx, id, appplatform.DefaultDeploymentsListOperationOptions())` can be used to do batched pagination
items, err := client.DeploymentsListComplete(ctx, id, appplatform.DefaultDeploymentsListOperationOptions())
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `AppPlatformClient.DeploymentsListForCluster`

```go
ctx := context.TODO()
id := appplatform.NewSpringID("12345678-1234-9876-4563-123456789012", "example-resource-group", "springValue")

// alternatively `client.DeploymentsListForCluster(ctx, id, appplatform.DefaultDeploymentsListForClusterOperationOptions())` can be used to do batched pagination
items, err := client.DeploymentsListForClusterComplete(ctx, id, appplatform.DefaultDeploymentsListForClusterOperationOptions())
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `AppPlatformClient.DeploymentsRestart`

```go
ctx := context.TODO()
id := appplatform.NewDeploymentID("12345678-1234-9876-4563-123456789012", "example-resource-group", "springValue", "appValue", "deploymentValue")

if err := client.DeploymentsRestartThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `AppPlatformClient.DeploymentsStart`

```go
ctx := context.TODO()
id := appplatform.NewDeploymentID("12345678-1234-9876-4563-123456789012", "example-resource-group", "springValue", "appValue", "deploymentValue")

if err := client.DeploymentsStartThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `AppPlatformClient.DeploymentsStartJFR`

```go
ctx := context.TODO()
id := appplatform.NewDeploymentID("12345678-1234-9876-4563-123456789012", "example-resource-group", "springValue", "appValue", "deploymentValue")

payload := appplatform.DiagnosticParameters{
	// ...
}


if err := client.DeploymentsStartJFRThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `AppPlatformClient.DeploymentsStop`

```go
ctx := context.TODO()
id := appplatform.NewDeploymentID("12345678-1234-9876-4563-123456789012", "example-resource-group", "springValue", "appValue", "deploymentValue")

if err := client.DeploymentsStopThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `AppPlatformClient.DeploymentsUpdate`

```go
ctx := context.TODO()
id := appplatform.NewDeploymentID("12345678-1234-9876-4563-123456789012", "example-resource-group", "springValue", "appValue", "deploymentValue")

payload := appplatform.DeploymentResource{
	// ...
}


if err := client.DeploymentsUpdateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `AppPlatformClient.DevToolPortalsCreateOrUpdate`

```go
ctx := context.TODO()
id := appplatform.NewDevToolPortalID("12345678-1234-9876-4563-123456789012", "example-resource-group", "springValue", "devToolPortalValue")

payload := appplatform.DevToolPortalResource{
	// ...
}


if err := client.DevToolPortalsCreateOrUpdateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `AppPlatformClient.DevToolPortalsDelete`

```go
ctx := context.TODO()
id := appplatform.NewDevToolPortalID("12345678-1234-9876-4563-123456789012", "example-resource-group", "springValue", "devToolPortalValue")

if err := client.DevToolPortalsDeleteThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `AppPlatformClient.DevToolPortalsGet`

```go
ctx := context.TODO()
id := appplatform.NewDevToolPortalID("12345678-1234-9876-4563-123456789012", "example-resource-group", "springValue", "devToolPortalValue")

read, err := client.DevToolPortalsGet(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `AppPlatformClient.DevToolPortalsList`

```go
ctx := context.TODO()
id := appplatform.NewSpringID("12345678-1234-9876-4563-123456789012", "example-resource-group", "springValue")

// alternatively `client.DevToolPortalsList(ctx, id)` can be used to do batched pagination
items, err := client.DevToolPortalsListComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `AppPlatformClient.EurekaServersGet`

```go
ctx := context.TODO()
id := appplatform.NewSpringID("12345678-1234-9876-4563-123456789012", "example-resource-group", "springValue")

read, err := client.EurekaServersGet(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `AppPlatformClient.EurekaServersList`

```go
ctx := context.TODO()
id := appplatform.NewSpringID("12345678-1234-9876-4563-123456789012", "example-resource-group", "springValue")

read, err := client.EurekaServersList(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `AppPlatformClient.EurekaServersUpdatePatch`

```go
ctx := context.TODO()
id := appplatform.NewSpringID("12345678-1234-9876-4563-123456789012", "example-resource-group", "springValue")

payload := appplatform.EurekaServerResource{
	// ...
}


if err := client.EurekaServersUpdatePatchThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `AppPlatformClient.EurekaServersUpdatePut`

```go
ctx := context.TODO()
id := appplatform.NewSpringID("12345678-1234-9876-4563-123456789012", "example-resource-group", "springValue")

payload := appplatform.EurekaServerResource{
	// ...
}


if err := client.EurekaServersUpdatePutThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `AppPlatformClient.GatewayCustomDomainsCreateOrUpdate`

```go
ctx := context.TODO()
id := appplatform.NewGatewayDomainID("12345678-1234-9876-4563-123456789012", "example-resource-group", "springValue", "gatewayValue", "domainValue")

payload := appplatform.GatewayCustomDomainResource{
	// ...
}


if err := client.GatewayCustomDomainsCreateOrUpdateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `AppPlatformClient.GatewayCustomDomainsDelete`

```go
ctx := context.TODO()
id := appplatform.NewGatewayDomainID("12345678-1234-9876-4563-123456789012", "example-resource-group", "springValue", "gatewayValue", "domainValue")

if err := client.GatewayCustomDomainsDeleteThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `AppPlatformClient.GatewayCustomDomainsGet`

```go
ctx := context.TODO()
id := appplatform.NewGatewayDomainID("12345678-1234-9876-4563-123456789012", "example-resource-group", "springValue", "gatewayValue", "domainValue")

read, err := client.GatewayCustomDomainsGet(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `AppPlatformClient.GatewayCustomDomainsList`

```go
ctx := context.TODO()
id := appplatform.NewGatewayID("12345678-1234-9876-4563-123456789012", "example-resource-group", "springValue", "gatewayValue")

// alternatively `client.GatewayCustomDomainsList(ctx, id)` can be used to do batched pagination
items, err := client.GatewayCustomDomainsListComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `AppPlatformClient.GatewayRouteConfigsCreateOrUpdate`

```go
ctx := context.TODO()
id := appplatform.NewRouteConfigID("12345678-1234-9876-4563-123456789012", "example-resource-group", "springValue", "gatewayValue", "routeConfigValue")

payload := appplatform.GatewayRouteConfigResource{
	// ...
}


if err := client.GatewayRouteConfigsCreateOrUpdateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `AppPlatformClient.GatewayRouteConfigsDelete`

```go
ctx := context.TODO()
id := appplatform.NewRouteConfigID("12345678-1234-9876-4563-123456789012", "example-resource-group", "springValue", "gatewayValue", "routeConfigValue")

if err := client.GatewayRouteConfigsDeleteThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `AppPlatformClient.GatewayRouteConfigsGet`

```go
ctx := context.TODO()
id := appplatform.NewRouteConfigID("12345678-1234-9876-4563-123456789012", "example-resource-group", "springValue", "gatewayValue", "routeConfigValue")

read, err := client.GatewayRouteConfigsGet(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `AppPlatformClient.GatewayRouteConfigsList`

```go
ctx := context.TODO()
id := appplatform.NewGatewayID("12345678-1234-9876-4563-123456789012", "example-resource-group", "springValue", "gatewayValue")

// alternatively `client.GatewayRouteConfigsList(ctx, id)` can be used to do batched pagination
items, err := client.GatewayRouteConfigsListComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `AppPlatformClient.GatewaysCreateOrUpdate`

```go
ctx := context.TODO()
id := appplatform.NewGatewayID("12345678-1234-9876-4563-123456789012", "example-resource-group", "springValue", "gatewayValue")

payload := appplatform.GatewayResource{
	// ...
}


if err := client.GatewaysCreateOrUpdateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `AppPlatformClient.GatewaysDelete`

```go
ctx := context.TODO()
id := appplatform.NewGatewayID("12345678-1234-9876-4563-123456789012", "example-resource-group", "springValue", "gatewayValue")

if err := client.GatewaysDeleteThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `AppPlatformClient.GatewaysGet`

```go
ctx := context.TODO()
id := appplatform.NewGatewayID("12345678-1234-9876-4563-123456789012", "example-resource-group", "springValue", "gatewayValue")

read, err := client.GatewaysGet(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `AppPlatformClient.GatewaysList`

```go
ctx := context.TODO()
id := appplatform.NewSpringID("12345678-1234-9876-4563-123456789012", "example-resource-group", "springValue")

// alternatively `client.GatewaysList(ctx, id)` can be used to do batched pagination
items, err := client.GatewaysListComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `AppPlatformClient.GatewaysListEnvSecrets`

```go
ctx := context.TODO()
id := appplatform.NewGatewayID("12345678-1234-9876-4563-123456789012", "example-resource-group", "springValue", "gatewayValue")

read, err := client.GatewaysListEnvSecrets(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `AppPlatformClient.GatewaysRestart`

```go
ctx := context.TODO()
id := appplatform.NewGatewayID("12345678-1234-9876-4563-123456789012", "example-resource-group", "springValue", "gatewayValue")

if err := client.GatewaysRestartThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `AppPlatformClient.GatewaysUpdateCapacity`

```go
ctx := context.TODO()
id := appplatform.NewGatewayID("12345678-1234-9876-4563-123456789012", "example-resource-group", "springValue", "gatewayValue")

payload := appplatform.SkuObject{
	// ...
}


if err := client.GatewaysUpdateCapacityThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `AppPlatformClient.GatewaysValidateDomain`

```go
ctx := context.TODO()
id := appplatform.NewGatewayID("12345678-1234-9876-4563-123456789012", "example-resource-group", "springValue", "gatewayValue")

payload := appplatform.CustomDomainValidatePayload{
	// ...
}


read, err := client.GatewaysValidateDomain(ctx, id, payload)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `AppPlatformClient.MonitoringSettingsGet`

```go
ctx := context.TODO()
id := appplatform.NewSpringID("12345678-1234-9876-4563-123456789012", "example-resource-group", "springValue")

read, err := client.MonitoringSettingsGet(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `AppPlatformClient.MonitoringSettingsUpdatePatch`

```go
ctx := context.TODO()
id := appplatform.NewSpringID("12345678-1234-9876-4563-123456789012", "example-resource-group", "springValue")

payload := appplatform.MonitoringSettingResource{
	// ...
}


if err := client.MonitoringSettingsUpdatePatchThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `AppPlatformClient.MonitoringSettingsUpdatePut`

```go
ctx := context.TODO()
id := appplatform.NewSpringID("12345678-1234-9876-4563-123456789012", "example-resource-group", "springValue")

payload := appplatform.MonitoringSettingResource{
	// ...
}


if err := client.MonitoringSettingsUpdatePutThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `AppPlatformClient.PredefinedAcceleratorsDisable`

```go
ctx := context.TODO()
id := appplatform.NewPredefinedAcceleratorID("12345678-1234-9876-4563-123456789012", "example-resource-group", "springValue", "applicationAcceleratorValue", "predefinedAcceleratorValue")

if err := client.PredefinedAcceleratorsDisableThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `AppPlatformClient.PredefinedAcceleratorsEnable`

```go
ctx := context.TODO()
id := appplatform.NewPredefinedAcceleratorID("12345678-1234-9876-4563-123456789012", "example-resource-group", "springValue", "applicationAcceleratorValue", "predefinedAcceleratorValue")

if err := client.PredefinedAcceleratorsEnableThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `AppPlatformClient.PredefinedAcceleratorsGet`

```go
ctx := context.TODO()
id := appplatform.NewPredefinedAcceleratorID("12345678-1234-9876-4563-123456789012", "example-resource-group", "springValue", "applicationAcceleratorValue", "predefinedAcceleratorValue")

read, err := client.PredefinedAcceleratorsGet(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `AppPlatformClient.PredefinedAcceleratorsList`

```go
ctx := context.TODO()
id := appplatform.NewApplicationAcceleratorID("12345678-1234-9876-4563-123456789012", "example-resource-group", "springValue", "applicationAcceleratorValue")

// alternatively `client.PredefinedAcceleratorsList(ctx, id)` can be used to do batched pagination
items, err := client.PredefinedAcceleratorsListComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `AppPlatformClient.RuntimeVersionsListRuntimeVersions`

```go
ctx := context.TODO()


read, err := client.RuntimeVersionsListRuntimeVersions(ctx)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `AppPlatformClient.ServiceRegistriesCreateOrUpdate`

```go
ctx := context.TODO()
id := appplatform.NewServiceRegistryID("12345678-1234-9876-4563-123456789012", "example-resource-group", "springValue", "serviceRegistryValue")

if err := client.ServiceRegistriesCreateOrUpdateThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `AppPlatformClient.ServiceRegistriesDelete`

```go
ctx := context.TODO()
id := appplatform.NewServiceRegistryID("12345678-1234-9876-4563-123456789012", "example-resource-group", "springValue", "serviceRegistryValue")

if err := client.ServiceRegistriesDeleteThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `AppPlatformClient.ServiceRegistriesGet`

```go
ctx := context.TODO()
id := appplatform.NewServiceRegistryID("12345678-1234-9876-4563-123456789012", "example-resource-group", "springValue", "serviceRegistryValue")

read, err := client.ServiceRegistriesGet(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `AppPlatformClient.ServiceRegistriesList`

```go
ctx := context.TODO()
id := appplatform.NewSpringID("12345678-1234-9876-4563-123456789012", "example-resource-group", "springValue")

// alternatively `client.ServiceRegistriesList(ctx, id)` can be used to do batched pagination
items, err := client.ServiceRegistriesListComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `AppPlatformClient.ServicesCheckNameAvailability`

```go
ctx := context.TODO()
id := appplatform.NewLocationID("12345678-1234-9876-4563-123456789012", "locationValue")

payload := appplatform.NameAvailabilityParameters{
	// ...
}


read, err := client.ServicesCheckNameAvailability(ctx, id, payload)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `AppPlatformClient.ServicesCreateOrUpdate`

```go
ctx := context.TODO()
id := appplatform.NewSpringID("12345678-1234-9876-4563-123456789012", "example-resource-group", "springValue")

payload := appplatform.ServiceResource{
	// ...
}


if err := client.ServicesCreateOrUpdateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `AppPlatformClient.ServicesDelete`

```go
ctx := context.TODO()
id := appplatform.NewSpringID("12345678-1234-9876-4563-123456789012", "example-resource-group", "springValue")

if err := client.ServicesDeleteThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `AppPlatformClient.ServicesDisableApmGlobally`

```go
ctx := context.TODO()
id := appplatform.NewSpringID("12345678-1234-9876-4563-123456789012", "example-resource-group", "springValue")

payload := appplatform.ApmReference{
	// ...
}


if err := client.ServicesDisableApmGloballyThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `AppPlatformClient.ServicesDisableTestEndpoint`

```go
ctx := context.TODO()
id := appplatform.NewSpringID("12345678-1234-9876-4563-123456789012", "example-resource-group", "springValue")

read, err := client.ServicesDisableTestEndpoint(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `AppPlatformClient.ServicesEnableApmGlobally`

```go
ctx := context.TODO()
id := appplatform.NewSpringID("12345678-1234-9876-4563-123456789012", "example-resource-group", "springValue")

payload := appplatform.ApmReference{
	// ...
}


if err := client.ServicesEnableApmGloballyThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `AppPlatformClient.ServicesEnableTestEndpoint`

```go
ctx := context.TODO()
id := appplatform.NewSpringID("12345678-1234-9876-4563-123456789012", "example-resource-group", "springValue")

read, err := client.ServicesEnableTestEndpoint(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `AppPlatformClient.ServicesFlushVnetDnsSetting`

```go
ctx := context.TODO()
id := appplatform.NewSpringID("12345678-1234-9876-4563-123456789012", "example-resource-group", "springValue")

if err := client.ServicesFlushVnetDnsSettingThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `AppPlatformClient.ServicesGet`

```go
ctx := context.TODO()
id := appplatform.NewSpringID("12345678-1234-9876-4563-123456789012", "example-resource-group", "springValue")

read, err := client.ServicesGet(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `AppPlatformClient.ServicesList`

```go
ctx := context.TODO()
id := appplatform.NewResourceGroupID("12345678-1234-9876-4563-123456789012", "example-resource-group")

// alternatively `client.ServicesList(ctx, id)` can be used to do batched pagination
items, err := client.ServicesListComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `AppPlatformClient.ServicesListBySubscription`

```go
ctx := context.TODO()
id := appplatform.NewSubscriptionID("12345678-1234-9876-4563-123456789012")

// alternatively `client.ServicesListBySubscription(ctx, id)` can be used to do batched pagination
items, err := client.ServicesListBySubscriptionComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `AppPlatformClient.ServicesListGloballyEnabledApms`

```go
ctx := context.TODO()
id := appplatform.NewSpringID("12345678-1234-9876-4563-123456789012", "example-resource-group", "springValue")

read, err := client.ServicesListGloballyEnabledApms(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `AppPlatformClient.ServicesListSupportedApmTypes`

```go
ctx := context.TODO()
id := appplatform.NewSpringID("12345678-1234-9876-4563-123456789012", "example-resource-group", "springValue")

// alternatively `client.ServicesListSupportedApmTypes(ctx, id)` can be used to do batched pagination
items, err := client.ServicesListSupportedApmTypesComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `AppPlatformClient.ServicesListSupportedServerVersions`

```go
ctx := context.TODO()
id := appplatform.NewSpringID("12345678-1234-9876-4563-123456789012", "example-resource-group", "springValue")

// alternatively `client.ServicesListSupportedServerVersions(ctx, id)` can be used to do batched pagination
items, err := client.ServicesListSupportedServerVersionsComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `AppPlatformClient.ServicesListTestKeys`

```go
ctx := context.TODO()
id := appplatform.NewSpringID("12345678-1234-9876-4563-123456789012", "example-resource-group", "springValue")

read, err := client.ServicesListTestKeys(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `AppPlatformClient.ServicesRegenerateTestKey`

```go
ctx := context.TODO()
id := appplatform.NewSpringID("12345678-1234-9876-4563-123456789012", "example-resource-group", "springValue")

payload := appplatform.RegenerateTestKeyRequestPayload{
	// ...
}


read, err := client.ServicesRegenerateTestKey(ctx, id, payload)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `AppPlatformClient.ServicesStart`

```go
ctx := context.TODO()
id := appplatform.NewSpringID("12345678-1234-9876-4563-123456789012", "example-resource-group", "springValue")

if err := client.ServicesStartThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `AppPlatformClient.ServicesStop`

```go
ctx := context.TODO()
id := appplatform.NewSpringID("12345678-1234-9876-4563-123456789012", "example-resource-group", "springValue")

if err := client.ServicesStopThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `AppPlatformClient.ServicesUpdate`

```go
ctx := context.TODO()
id := appplatform.NewSpringID("12345678-1234-9876-4563-123456789012", "example-resource-group", "springValue")

payload := appplatform.ServiceResource{
	// ...
}


if err := client.ServicesUpdateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `AppPlatformClient.SkusList`

```go
ctx := context.TODO()
id := appplatform.NewSubscriptionID("12345678-1234-9876-4563-123456789012")

// alternatively `client.SkusList(ctx, id)` can be used to do batched pagination
items, err := client.SkusListComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `AppPlatformClient.StoragesCreateOrUpdate`

```go
ctx := context.TODO()
id := appplatform.NewStorageID("12345678-1234-9876-4563-123456789012", "example-resource-group", "springValue", "storageValue")

payload := appplatform.StorageResource{
	// ...
}


if err := client.StoragesCreateOrUpdateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `AppPlatformClient.StoragesDelete`

```go
ctx := context.TODO()
id := appplatform.NewStorageID("12345678-1234-9876-4563-123456789012", "example-resource-group", "springValue", "storageValue")

if err := client.StoragesDeleteThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `AppPlatformClient.StoragesGet`

```go
ctx := context.TODO()
id := appplatform.NewStorageID("12345678-1234-9876-4563-123456789012", "example-resource-group", "springValue", "storageValue")

read, err := client.StoragesGet(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `AppPlatformClient.StoragesList`

```go
ctx := context.TODO()
id := appplatform.NewSpringID("12345678-1234-9876-4563-123456789012", "example-resource-group", "springValue")

// alternatively `client.StoragesList(ctx, id)` can be used to do batched pagination
items, err := client.StoragesListComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```
