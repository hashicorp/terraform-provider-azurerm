
## `github.com/hashicorp/go-azure-sdk/resource-manager/datafactory/2018-06-01/factories` Documentation

The `factories` SDK allows for interaction with the Azure Resource Manager Service `datafactory` (API Version `2018-06-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/datafactory/2018-06-01/factories"
```


### Client Initialization

```go
client := factories.NewFactoriesClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `FactoriesClient.ConfigureFactoryRepo`

```go
ctx := context.TODO()
id := factories.NewLocationID("12345678-1234-9876-4563-123456789012", "locationIdValue")

payload := factories.FactoryRepoUpdate{
	// ...
}


read, err := client.ConfigureFactoryRepo(ctx, id, payload)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `FactoriesClient.CreateOrUpdate`

```go
ctx := context.TODO()
id := factories.NewFactoryID("12345678-1234-9876-4563-123456789012", "example-resource-group", "factoryValue")

payload := factories.Factory{
	// ...
}


read, err := client.CreateOrUpdate(ctx, id, payload, factories.DefaultCreateOrUpdateOperationOptions())
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `FactoriesClient.Delete`

```go
ctx := context.TODO()
id := factories.NewFactoryID("12345678-1234-9876-4563-123456789012", "example-resource-group", "factoryValue")

read, err := client.Delete(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `FactoriesClient.Get`

```go
ctx := context.TODO()
id := factories.NewFactoryID("12345678-1234-9876-4563-123456789012", "example-resource-group", "factoryValue")

read, err := client.Get(ctx, id, factories.DefaultGetOperationOptions())
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `FactoriesClient.GetDataPlaneAccess`

```go
ctx := context.TODO()
id := factories.NewFactoryID("12345678-1234-9876-4563-123456789012", "example-resource-group", "factoryValue")

payload := factories.UserAccessPolicy{
	// ...
}


read, err := client.GetDataPlaneAccess(ctx, id, payload)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `FactoriesClient.GetGitHubAccessToken`

```go
ctx := context.TODO()
id := factories.NewFactoryID("12345678-1234-9876-4563-123456789012", "example-resource-group", "factoryValue")

payload := factories.GitHubAccessTokenRequest{
	// ...
}


read, err := client.GetGitHubAccessToken(ctx, id, payload)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `FactoriesClient.List`

```go
ctx := context.TODO()
id := factories.NewSubscriptionID("12345678-1234-9876-4563-123456789012")

// alternatively `client.List(ctx, id)` can be used to do batched pagination
items, err := client.ListComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `FactoriesClient.ListByResourceGroup`

```go
ctx := context.TODO()
id := factories.NewResourceGroupID("12345678-1234-9876-4563-123456789012", "example-resource-group")

// alternatively `client.ListByResourceGroup(ctx, id)` can be used to do batched pagination
items, err := client.ListByResourceGroupComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `FactoriesClient.Update`

```go
ctx := context.TODO()
id := factories.NewFactoryID("12345678-1234-9876-4563-123456789012", "example-resource-group", "factoryValue")

payload := factories.FactoryUpdateParameters{
	// ...
}


read, err := client.Update(ctx, id, payload)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```
