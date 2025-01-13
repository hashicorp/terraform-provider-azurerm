
## `github.com/hashicorp/go-azure-sdk/resource-manager/web/2023-01-01/staticsites` Documentation

The `staticsites` SDK allows for interaction with Azure Resource Manager `web` (API Version `2023-01-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
import "github.com/hashicorp/go-azure-sdk/resource-manager/web/2023-01-01/staticsites"
```


### Client Initialization

```go
client := staticsites.NewStaticSitesClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `StaticSitesClient.ApproveOrRejectPrivateEndpointConnection`

```go
ctx := context.TODO()
id := staticsites.NewStaticSitePrivateEndpointConnectionID("12345678-1234-9876-4563-123456789012", "example-resource-group", "staticSiteName", "privateEndpointConnectionName")

payload := staticsites.RemotePrivateEndpointConnectionARMResource{
	// ...
}


if err := client.ApproveOrRejectPrivateEndpointConnectionThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `StaticSitesClient.CreateOrUpdateBasicAuth`

```go
ctx := context.TODO()
id := staticsites.NewStaticSiteID("12345678-1234-9876-4563-123456789012", "example-resource-group", "staticSiteName")

payload := staticsites.StaticSiteBasicAuthPropertiesARMResource{
	// ...
}


read, err := client.CreateOrUpdateBasicAuth(ctx, id, payload)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `StaticSitesClient.CreateOrUpdateBuildDatabaseConnection`

```go
ctx := context.TODO()
id := staticsites.NewBuildDatabaseConnectionID("12345678-1234-9876-4563-123456789012", "example-resource-group", "staticSiteName", "buildName", "databaseConnectionName")

payload := staticsites.DatabaseConnection{
	// ...
}


read, err := client.CreateOrUpdateBuildDatabaseConnection(ctx, id, payload)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `StaticSitesClient.CreateOrUpdateDatabaseConnection`

```go
ctx := context.TODO()
id := staticsites.NewDatabaseConnectionID("12345678-1234-9876-4563-123456789012", "example-resource-group", "staticSiteName", "databaseConnectionName")

payload := staticsites.DatabaseConnection{
	// ...
}


read, err := client.CreateOrUpdateDatabaseConnection(ctx, id, payload)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `StaticSitesClient.CreateOrUpdateStaticSite`

```go
ctx := context.TODO()
id := staticsites.NewStaticSiteID("12345678-1234-9876-4563-123456789012", "example-resource-group", "staticSiteName")

payload := staticsites.StaticSiteARMResource{
	// ...
}


if err := client.CreateOrUpdateStaticSiteThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `StaticSitesClient.CreateOrUpdateStaticSiteAppSettings`

```go
ctx := context.TODO()
id := staticsites.NewStaticSiteID("12345678-1234-9876-4563-123456789012", "example-resource-group", "staticSiteName")

payload := staticsites.StringDictionary{
	// ...
}


read, err := client.CreateOrUpdateStaticSiteAppSettings(ctx, id, payload)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `StaticSitesClient.CreateOrUpdateStaticSiteBuildAppSettings`

```go
ctx := context.TODO()
id := staticsites.NewBuildID("12345678-1234-9876-4563-123456789012", "example-resource-group", "staticSiteName", "buildName")

payload := staticsites.StringDictionary{
	// ...
}


read, err := client.CreateOrUpdateStaticSiteBuildAppSettings(ctx, id, payload)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `StaticSitesClient.CreateOrUpdateStaticSiteBuildFunctionAppSettings`

```go
ctx := context.TODO()
id := staticsites.NewBuildID("12345678-1234-9876-4563-123456789012", "example-resource-group", "staticSiteName", "buildName")

payload := staticsites.StringDictionary{
	// ...
}


read, err := client.CreateOrUpdateStaticSiteBuildFunctionAppSettings(ctx, id, payload)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `StaticSitesClient.CreateOrUpdateStaticSiteCustomDomain`

```go
ctx := context.TODO()
id := staticsites.NewCustomDomainID("12345678-1234-9876-4563-123456789012", "example-resource-group", "staticSiteName", "customDomainName")

payload := staticsites.StaticSiteCustomDomainRequestPropertiesARMResource{
	// ...
}


if err := client.CreateOrUpdateStaticSiteCustomDomainThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `StaticSitesClient.CreateOrUpdateStaticSiteFunctionAppSettings`

```go
ctx := context.TODO()
id := staticsites.NewStaticSiteID("12345678-1234-9876-4563-123456789012", "example-resource-group", "staticSiteName")

payload := staticsites.StringDictionary{
	// ...
}


read, err := client.CreateOrUpdateStaticSiteFunctionAppSettings(ctx, id, payload)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `StaticSitesClient.CreateUserRolesInvitationLink`

```go
ctx := context.TODO()
id := staticsites.NewStaticSiteID("12345678-1234-9876-4563-123456789012", "example-resource-group", "staticSiteName")

payload := staticsites.StaticSiteUserInvitationRequestResource{
	// ...
}


read, err := client.CreateUserRolesInvitationLink(ctx, id, payload)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `StaticSitesClient.CreateZipDeploymentForStaticSite`

```go
ctx := context.TODO()
id := staticsites.NewStaticSiteID("12345678-1234-9876-4563-123456789012", "example-resource-group", "staticSiteName")

payload := staticsites.StaticSiteZipDeploymentARMResource{
	// ...
}


if err := client.CreateZipDeploymentForStaticSiteThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `StaticSitesClient.CreateZipDeploymentForStaticSiteBuild`

```go
ctx := context.TODO()
id := staticsites.NewBuildID("12345678-1234-9876-4563-123456789012", "example-resource-group", "staticSiteName", "buildName")

payload := staticsites.StaticSiteZipDeploymentARMResource{
	// ...
}


if err := client.CreateZipDeploymentForStaticSiteBuildThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `StaticSitesClient.DeleteBuildDatabaseConnection`

```go
ctx := context.TODO()
id := staticsites.NewBuildDatabaseConnectionID("12345678-1234-9876-4563-123456789012", "example-resource-group", "staticSiteName", "buildName", "databaseConnectionName")

read, err := client.DeleteBuildDatabaseConnection(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `StaticSitesClient.DeleteDatabaseConnection`

```go
ctx := context.TODO()
id := staticsites.NewDatabaseConnectionID("12345678-1234-9876-4563-123456789012", "example-resource-group", "staticSiteName", "databaseConnectionName")

read, err := client.DeleteDatabaseConnection(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `StaticSitesClient.DeletePrivateEndpointConnection`

```go
ctx := context.TODO()
id := staticsites.NewStaticSitePrivateEndpointConnectionID("12345678-1234-9876-4563-123456789012", "example-resource-group", "staticSiteName", "privateEndpointConnectionName")

if err := client.DeletePrivateEndpointConnectionThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `StaticSitesClient.DeleteStaticSite`

```go
ctx := context.TODO()
id := staticsites.NewStaticSiteID("12345678-1234-9876-4563-123456789012", "example-resource-group", "staticSiteName")

if err := client.DeleteStaticSiteThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `StaticSitesClient.DeleteStaticSiteBuild`

```go
ctx := context.TODO()
id := staticsites.NewBuildID("12345678-1234-9876-4563-123456789012", "example-resource-group", "staticSiteName", "buildName")

if err := client.DeleteStaticSiteBuildThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `StaticSitesClient.DeleteStaticSiteCustomDomain`

```go
ctx := context.TODO()
id := staticsites.NewCustomDomainID("12345678-1234-9876-4563-123456789012", "example-resource-group", "staticSiteName", "customDomainName")

if err := client.DeleteStaticSiteCustomDomainThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `StaticSitesClient.DeleteStaticSiteUser`

```go
ctx := context.TODO()
id := staticsites.NewUserID("12345678-1234-9876-4563-123456789012", "example-resource-group", "staticSiteName", "authProviderName", "userName")

read, err := client.DeleteStaticSiteUser(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `StaticSitesClient.DetachStaticSite`

```go
ctx := context.TODO()
id := staticsites.NewStaticSiteID("12345678-1234-9876-4563-123456789012", "example-resource-group", "staticSiteName")

if err := client.DetachStaticSiteThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `StaticSitesClient.DetachUserProvidedFunctionAppFromStaticSite`

```go
ctx := context.TODO()
id := staticsites.NewUserProvidedFunctionAppID("12345678-1234-9876-4563-123456789012", "example-resource-group", "staticSiteName", "userProvidedFunctionAppName")

read, err := client.DetachUserProvidedFunctionAppFromStaticSite(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `StaticSitesClient.DetachUserProvidedFunctionAppFromStaticSiteBuild`

```go
ctx := context.TODO()
id := staticsites.NewBuildUserProvidedFunctionAppID("12345678-1234-9876-4563-123456789012", "example-resource-group", "staticSiteName", "buildName", "userProvidedFunctionAppName")

read, err := client.DetachUserProvidedFunctionAppFromStaticSiteBuild(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `StaticSitesClient.GetBasicAuth`

```go
ctx := context.TODO()
id := staticsites.NewStaticSiteID("12345678-1234-9876-4563-123456789012", "example-resource-group", "staticSiteName")

read, err := client.GetBasicAuth(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `StaticSitesClient.GetBuildDatabaseConnection`

```go
ctx := context.TODO()
id := staticsites.NewBuildDatabaseConnectionID("12345678-1234-9876-4563-123456789012", "example-resource-group", "staticSiteName", "buildName", "databaseConnectionName")

read, err := client.GetBuildDatabaseConnection(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `StaticSitesClient.GetBuildDatabaseConnectionWithDetails`

```go
ctx := context.TODO()
id := staticsites.NewBuildDatabaseConnectionID("12345678-1234-9876-4563-123456789012", "example-resource-group", "staticSiteName", "buildName", "databaseConnectionName")

read, err := client.GetBuildDatabaseConnectionWithDetails(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `StaticSitesClient.GetBuildDatabaseConnections`

```go
ctx := context.TODO()
id := staticsites.NewBuildID("12345678-1234-9876-4563-123456789012", "example-resource-group", "staticSiteName", "buildName")

// alternatively `client.GetBuildDatabaseConnections(ctx, id)` can be used to do batched pagination
items, err := client.GetBuildDatabaseConnectionsComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `StaticSitesClient.GetBuildDatabaseConnectionsWithDetails`

```go
ctx := context.TODO()
id := staticsites.NewBuildID("12345678-1234-9876-4563-123456789012", "example-resource-group", "staticSiteName", "buildName")

// alternatively `client.GetBuildDatabaseConnectionsWithDetails(ctx, id)` can be used to do batched pagination
items, err := client.GetBuildDatabaseConnectionsWithDetailsComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `StaticSitesClient.GetDatabaseConnection`

```go
ctx := context.TODO()
id := staticsites.NewDatabaseConnectionID("12345678-1234-9876-4563-123456789012", "example-resource-group", "staticSiteName", "databaseConnectionName")

read, err := client.GetDatabaseConnection(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `StaticSitesClient.GetDatabaseConnectionWithDetails`

```go
ctx := context.TODO()
id := staticsites.NewDatabaseConnectionID("12345678-1234-9876-4563-123456789012", "example-resource-group", "staticSiteName", "databaseConnectionName")

read, err := client.GetDatabaseConnectionWithDetails(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `StaticSitesClient.GetDatabaseConnections`

```go
ctx := context.TODO()
id := staticsites.NewStaticSiteID("12345678-1234-9876-4563-123456789012", "example-resource-group", "staticSiteName")

// alternatively `client.GetDatabaseConnections(ctx, id)` can be used to do batched pagination
items, err := client.GetDatabaseConnectionsComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `StaticSitesClient.GetDatabaseConnectionsWithDetails`

```go
ctx := context.TODO()
id := staticsites.NewStaticSiteID("12345678-1234-9876-4563-123456789012", "example-resource-group", "staticSiteName")

// alternatively `client.GetDatabaseConnectionsWithDetails(ctx, id)` can be used to do batched pagination
items, err := client.GetDatabaseConnectionsWithDetailsComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `StaticSitesClient.GetLinkedBackend`

```go
ctx := context.TODO()
id := staticsites.NewLinkedBackendID("12345678-1234-9876-4563-123456789012", "example-resource-group", "staticSiteName", "linkedBackendName")

read, err := client.GetLinkedBackend(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `StaticSitesClient.GetLinkedBackendForBuild`

```go
ctx := context.TODO()
id := staticsites.NewBuildLinkedBackendID("12345678-1234-9876-4563-123456789012", "example-resource-group", "staticSiteName", "buildName", "linkedBackendName")

read, err := client.GetLinkedBackendForBuild(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `StaticSitesClient.GetLinkedBackends`

```go
ctx := context.TODO()
id := staticsites.NewStaticSiteID("12345678-1234-9876-4563-123456789012", "example-resource-group", "staticSiteName")

// alternatively `client.GetLinkedBackends(ctx, id)` can be used to do batched pagination
items, err := client.GetLinkedBackendsComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `StaticSitesClient.GetLinkedBackendsForBuild`

```go
ctx := context.TODO()
id := staticsites.NewBuildID("12345678-1234-9876-4563-123456789012", "example-resource-group", "staticSiteName", "buildName")

// alternatively `client.GetLinkedBackendsForBuild(ctx, id)` can be used to do batched pagination
items, err := client.GetLinkedBackendsForBuildComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `StaticSitesClient.GetPrivateEndpointConnection`

```go
ctx := context.TODO()
id := staticsites.NewStaticSitePrivateEndpointConnectionID("12345678-1234-9876-4563-123456789012", "example-resource-group", "staticSiteName", "privateEndpointConnectionName")

read, err := client.GetPrivateEndpointConnection(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `StaticSitesClient.GetPrivateEndpointConnectionList`

```go
ctx := context.TODO()
id := staticsites.NewStaticSiteID("12345678-1234-9876-4563-123456789012", "example-resource-group", "staticSiteName")

// alternatively `client.GetPrivateEndpointConnectionList(ctx, id)` can be used to do batched pagination
items, err := client.GetPrivateEndpointConnectionListComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `StaticSitesClient.GetPrivateLinkResources`

```go
ctx := context.TODO()
id := staticsites.NewStaticSiteID("12345678-1234-9876-4563-123456789012", "example-resource-group", "staticSiteName")

read, err := client.GetPrivateLinkResources(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `StaticSitesClient.GetStaticSite`

```go
ctx := context.TODO()
id := staticsites.NewStaticSiteID("12345678-1234-9876-4563-123456789012", "example-resource-group", "staticSiteName")

read, err := client.GetStaticSite(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `StaticSitesClient.GetStaticSiteBuild`

```go
ctx := context.TODO()
id := staticsites.NewBuildID("12345678-1234-9876-4563-123456789012", "example-resource-group", "staticSiteName", "buildName")

read, err := client.GetStaticSiteBuild(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `StaticSitesClient.GetStaticSiteBuilds`

```go
ctx := context.TODO()
id := staticsites.NewStaticSiteID("12345678-1234-9876-4563-123456789012", "example-resource-group", "staticSiteName")

// alternatively `client.GetStaticSiteBuilds(ctx, id)` can be used to do batched pagination
items, err := client.GetStaticSiteBuildsComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `StaticSitesClient.GetStaticSiteCustomDomain`

```go
ctx := context.TODO()
id := staticsites.NewCustomDomainID("12345678-1234-9876-4563-123456789012", "example-resource-group", "staticSiteName", "customDomainName")

read, err := client.GetStaticSiteCustomDomain(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `StaticSitesClient.GetStaticSitesByResourceGroup`

```go
ctx := context.TODO()
id := commonids.NewResourceGroupID("12345678-1234-9876-4563-123456789012", "example-resource-group")

// alternatively `client.GetStaticSitesByResourceGroup(ctx, id)` can be used to do batched pagination
items, err := client.GetStaticSitesByResourceGroupComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `StaticSitesClient.GetUserProvidedFunctionAppForStaticSite`

```go
ctx := context.TODO()
id := staticsites.NewUserProvidedFunctionAppID("12345678-1234-9876-4563-123456789012", "example-resource-group", "staticSiteName", "userProvidedFunctionAppName")

read, err := client.GetUserProvidedFunctionAppForStaticSite(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `StaticSitesClient.GetUserProvidedFunctionAppForStaticSiteBuild`

```go
ctx := context.TODO()
id := staticsites.NewBuildUserProvidedFunctionAppID("12345678-1234-9876-4563-123456789012", "example-resource-group", "staticSiteName", "buildName", "userProvidedFunctionAppName")

read, err := client.GetUserProvidedFunctionAppForStaticSiteBuild(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `StaticSitesClient.GetUserProvidedFunctionAppsForStaticSite`

```go
ctx := context.TODO()
id := staticsites.NewStaticSiteID("12345678-1234-9876-4563-123456789012", "example-resource-group", "staticSiteName")

// alternatively `client.GetUserProvidedFunctionAppsForStaticSite(ctx, id)` can be used to do batched pagination
items, err := client.GetUserProvidedFunctionAppsForStaticSiteComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `StaticSitesClient.GetUserProvidedFunctionAppsForStaticSiteBuild`

```go
ctx := context.TODO()
id := staticsites.NewBuildID("12345678-1234-9876-4563-123456789012", "example-resource-group", "staticSiteName", "buildName")

// alternatively `client.GetUserProvidedFunctionAppsForStaticSiteBuild(ctx, id)` can be used to do batched pagination
items, err := client.GetUserProvidedFunctionAppsForStaticSiteBuildComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `StaticSitesClient.LinkBackend`

```go
ctx := context.TODO()
id := staticsites.NewLinkedBackendID("12345678-1234-9876-4563-123456789012", "example-resource-group", "staticSiteName", "linkedBackendName")

payload := staticsites.StaticSiteLinkedBackendARMResource{
	// ...
}


if err := client.LinkBackendThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `StaticSitesClient.LinkBackendToBuild`

```go
ctx := context.TODO()
id := staticsites.NewBuildLinkedBackendID("12345678-1234-9876-4563-123456789012", "example-resource-group", "staticSiteName", "buildName", "linkedBackendName")

payload := staticsites.StaticSiteLinkedBackendARMResource{
	// ...
}


if err := client.LinkBackendToBuildThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `StaticSitesClient.List`

```go
ctx := context.TODO()
id := commonids.NewSubscriptionID("12345678-1234-9876-4563-123456789012")

// alternatively `client.List(ctx, id)` can be used to do batched pagination
items, err := client.ListComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `StaticSitesClient.ListBasicAuth`

```go
ctx := context.TODO()
id := staticsites.NewStaticSiteID("12345678-1234-9876-4563-123456789012", "example-resource-group", "staticSiteName")

// alternatively `client.ListBasicAuth(ctx, id)` can be used to do batched pagination
items, err := client.ListBasicAuthComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `StaticSitesClient.ListStaticSiteAppSettings`

```go
ctx := context.TODO()
id := staticsites.NewStaticSiteID("12345678-1234-9876-4563-123456789012", "example-resource-group", "staticSiteName")

read, err := client.ListStaticSiteAppSettings(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `StaticSitesClient.ListStaticSiteBuildAppSettings`

```go
ctx := context.TODO()
id := staticsites.NewBuildID("12345678-1234-9876-4563-123456789012", "example-resource-group", "staticSiteName", "buildName")

read, err := client.ListStaticSiteBuildAppSettings(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `StaticSitesClient.ListStaticSiteBuildFunctionAppSettings`

```go
ctx := context.TODO()
id := staticsites.NewBuildID("12345678-1234-9876-4563-123456789012", "example-resource-group", "staticSiteName", "buildName")

read, err := client.ListStaticSiteBuildFunctionAppSettings(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `StaticSitesClient.ListStaticSiteBuildFunctions`

```go
ctx := context.TODO()
id := staticsites.NewBuildID("12345678-1234-9876-4563-123456789012", "example-resource-group", "staticSiteName", "buildName")

// alternatively `client.ListStaticSiteBuildFunctions(ctx, id)` can be used to do batched pagination
items, err := client.ListStaticSiteBuildFunctionsComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `StaticSitesClient.ListStaticSiteConfiguredRoles`

```go
ctx := context.TODO()
id := staticsites.NewStaticSiteID("12345678-1234-9876-4563-123456789012", "example-resource-group", "staticSiteName")

read, err := client.ListStaticSiteConfiguredRoles(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `StaticSitesClient.ListStaticSiteCustomDomains`

```go
ctx := context.TODO()
id := staticsites.NewStaticSiteID("12345678-1234-9876-4563-123456789012", "example-resource-group", "staticSiteName")

// alternatively `client.ListStaticSiteCustomDomains(ctx, id)` can be used to do batched pagination
items, err := client.ListStaticSiteCustomDomainsComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `StaticSitesClient.ListStaticSiteFunctionAppSettings`

```go
ctx := context.TODO()
id := staticsites.NewStaticSiteID("12345678-1234-9876-4563-123456789012", "example-resource-group", "staticSiteName")

read, err := client.ListStaticSiteFunctionAppSettings(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `StaticSitesClient.ListStaticSiteFunctions`

```go
ctx := context.TODO()
id := staticsites.NewStaticSiteID("12345678-1234-9876-4563-123456789012", "example-resource-group", "staticSiteName")

// alternatively `client.ListStaticSiteFunctions(ctx, id)` can be used to do batched pagination
items, err := client.ListStaticSiteFunctionsComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `StaticSitesClient.ListStaticSiteSecrets`

```go
ctx := context.TODO()
id := staticsites.NewStaticSiteID("12345678-1234-9876-4563-123456789012", "example-resource-group", "staticSiteName")

read, err := client.ListStaticSiteSecrets(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `StaticSitesClient.ListStaticSiteUsers`

```go
ctx := context.TODO()
id := staticsites.NewAuthProviderID("12345678-1234-9876-4563-123456789012", "example-resource-group", "staticSiteName", "authProviderName")

// alternatively `client.ListStaticSiteUsers(ctx, id)` can be used to do batched pagination
items, err := client.ListStaticSiteUsersComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `StaticSitesClient.PreviewWorkflow`

```go
ctx := context.TODO()
id := staticsites.NewProviderLocationID("12345678-1234-9876-4563-123456789012", "locationName")

payload := staticsites.StaticSitesWorkflowPreviewRequest{
	// ...
}


read, err := client.PreviewWorkflow(ctx, id, payload)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `StaticSitesClient.RegisterUserProvidedFunctionAppWithStaticSite`

```go
ctx := context.TODO()
id := staticsites.NewUserProvidedFunctionAppID("12345678-1234-9876-4563-123456789012", "example-resource-group", "staticSiteName", "userProvidedFunctionAppName")

payload := staticsites.StaticSiteUserProvidedFunctionAppARMResource{
	// ...
}


if err := client.RegisterUserProvidedFunctionAppWithStaticSiteThenPoll(ctx, id, payload, staticsites.DefaultRegisterUserProvidedFunctionAppWithStaticSiteOperationOptions()); err != nil {
	// handle the error
}
```


### Example Usage: `StaticSitesClient.RegisterUserProvidedFunctionAppWithStaticSiteBuild`

```go
ctx := context.TODO()
id := staticsites.NewBuildUserProvidedFunctionAppID("12345678-1234-9876-4563-123456789012", "example-resource-group", "staticSiteName", "buildName", "userProvidedFunctionAppName")

payload := staticsites.StaticSiteUserProvidedFunctionAppARMResource{
	// ...
}


if err := client.RegisterUserProvidedFunctionAppWithStaticSiteBuildThenPoll(ctx, id, payload, staticsites.DefaultRegisterUserProvidedFunctionAppWithStaticSiteBuildOperationOptions()); err != nil {
	// handle the error
}
```


### Example Usage: `StaticSitesClient.ResetStaticSiteApiKey`

```go
ctx := context.TODO()
id := staticsites.NewStaticSiteID("12345678-1234-9876-4563-123456789012", "example-resource-group", "staticSiteName")

payload := staticsites.StaticSiteResetPropertiesARMResource{
	// ...
}


read, err := client.ResetStaticSiteApiKey(ctx, id, payload)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `StaticSitesClient.UnlinkBackend`

```go
ctx := context.TODO()
id := staticsites.NewLinkedBackendID("12345678-1234-9876-4563-123456789012", "example-resource-group", "staticSiteName", "linkedBackendName")

read, err := client.UnlinkBackend(ctx, id, staticsites.DefaultUnlinkBackendOperationOptions())
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `StaticSitesClient.UnlinkBackendFromBuild`

```go
ctx := context.TODO()
id := staticsites.NewBuildLinkedBackendID("12345678-1234-9876-4563-123456789012", "example-resource-group", "staticSiteName", "buildName", "linkedBackendName")

read, err := client.UnlinkBackendFromBuild(ctx, id, staticsites.DefaultUnlinkBackendFromBuildOperationOptions())
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `StaticSitesClient.UpdateBuildDatabaseConnection`

```go
ctx := context.TODO()
id := staticsites.NewBuildDatabaseConnectionID("12345678-1234-9876-4563-123456789012", "example-resource-group", "staticSiteName", "buildName", "databaseConnectionName")

payload := staticsites.DatabaseConnectionPatchRequest{
	// ...
}


read, err := client.UpdateBuildDatabaseConnection(ctx, id, payload)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `StaticSitesClient.UpdateDatabaseConnection`

```go
ctx := context.TODO()
id := staticsites.NewDatabaseConnectionID("12345678-1234-9876-4563-123456789012", "example-resource-group", "staticSiteName", "databaseConnectionName")

payload := staticsites.DatabaseConnectionPatchRequest{
	// ...
}


read, err := client.UpdateDatabaseConnection(ctx, id, payload)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `StaticSitesClient.UpdateStaticSite`

```go
ctx := context.TODO()
id := staticsites.NewStaticSiteID("12345678-1234-9876-4563-123456789012", "example-resource-group", "staticSiteName")

payload := staticsites.StaticSitePatchResource{
	// ...
}


read, err := client.UpdateStaticSite(ctx, id, payload)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `StaticSitesClient.UpdateStaticSiteUser`

```go
ctx := context.TODO()
id := staticsites.NewUserID("12345678-1234-9876-4563-123456789012", "example-resource-group", "staticSiteName", "authProviderName", "userName")

payload := staticsites.StaticSiteUserARMResource{
	// ...
}


read, err := client.UpdateStaticSiteUser(ctx, id, payload)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `StaticSitesClient.ValidateBackend`

```go
ctx := context.TODO()
id := staticsites.NewLinkedBackendID("12345678-1234-9876-4563-123456789012", "example-resource-group", "staticSiteName", "linkedBackendName")

payload := staticsites.StaticSiteLinkedBackendARMResource{
	// ...
}


if err := client.ValidateBackendThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `StaticSitesClient.ValidateBackendForBuild`

```go
ctx := context.TODO()
id := staticsites.NewBuildLinkedBackendID("12345678-1234-9876-4563-123456789012", "example-resource-group", "staticSiteName", "buildName", "linkedBackendName")

payload := staticsites.StaticSiteLinkedBackendARMResource{
	// ...
}


if err := client.ValidateBackendForBuildThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `StaticSitesClient.ValidateCustomDomainCanBeAddedToStaticSite`

```go
ctx := context.TODO()
id := staticsites.NewCustomDomainID("12345678-1234-9876-4563-123456789012", "example-resource-group", "staticSiteName", "customDomainName")

payload := staticsites.StaticSiteCustomDomainRequestPropertiesARMResource{
	// ...
}


if err := client.ValidateCustomDomainCanBeAddedToStaticSiteThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```
