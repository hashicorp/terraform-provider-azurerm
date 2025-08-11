
## `github.com/hashicorp/go-azure-sdk/resource-manager/containerapps/2025-01-01/daprcomponents` Documentation

The `daprcomponents` SDK allows for interaction with Azure Resource Manager `containerapps` (API Version `2025-01-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/containerapps/2025-01-01/daprcomponents"
```


### Client Initialization

```go
client := daprcomponents.NewDaprComponentsClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `DaprComponentsClient.ConnectedEnvironmentsDaprComponentsCreateOrUpdate`

```go
ctx := context.TODO()
id := daprcomponents.NewConnectedEnvironmentDaprComponentID("12345678-1234-9876-4563-123456789012", "example-resource-group", "connectedEnvironmentName", "daprComponentName")

payload := daprcomponents.DaprComponent{
	// ...
}


read, err := client.ConnectedEnvironmentsDaprComponentsCreateOrUpdate(ctx, id, payload)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `DaprComponentsClient.ConnectedEnvironmentsDaprComponentsDelete`

```go
ctx := context.TODO()
id := daprcomponents.NewConnectedEnvironmentDaprComponentID("12345678-1234-9876-4563-123456789012", "example-resource-group", "connectedEnvironmentName", "daprComponentName")

read, err := client.ConnectedEnvironmentsDaprComponentsDelete(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `DaprComponentsClient.ConnectedEnvironmentsDaprComponentsGet`

```go
ctx := context.TODO()
id := daprcomponents.NewConnectedEnvironmentDaprComponentID("12345678-1234-9876-4563-123456789012", "example-resource-group", "connectedEnvironmentName", "daprComponentName")

read, err := client.ConnectedEnvironmentsDaprComponentsGet(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `DaprComponentsClient.ConnectedEnvironmentsDaprComponentsList`

```go
ctx := context.TODO()
id := daprcomponents.NewConnectedEnvironmentID("12345678-1234-9876-4563-123456789012", "example-resource-group", "connectedEnvironmentName")

// alternatively `client.ConnectedEnvironmentsDaprComponentsList(ctx, id)` can be used to do batched pagination
items, err := client.ConnectedEnvironmentsDaprComponentsListComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `DaprComponentsClient.ConnectedEnvironmentsDaprComponentsListSecrets`

```go
ctx := context.TODO()
id := daprcomponents.NewConnectedEnvironmentDaprComponentID("12345678-1234-9876-4563-123456789012", "example-resource-group", "connectedEnvironmentName", "daprComponentName")

read, err := client.ConnectedEnvironmentsDaprComponentsListSecrets(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `DaprComponentsClient.CreateOrUpdate`

```go
ctx := context.TODO()
id := daprcomponents.NewDaprComponentID("12345678-1234-9876-4563-123456789012", "example-resource-group", "managedEnvironmentName", "daprComponentName")

payload := daprcomponents.DaprComponent{
	// ...
}


read, err := client.CreateOrUpdate(ctx, id, payload)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `DaprComponentsClient.Delete`

```go
ctx := context.TODO()
id := daprcomponents.NewDaprComponentID("12345678-1234-9876-4563-123456789012", "example-resource-group", "managedEnvironmentName", "daprComponentName")

read, err := client.Delete(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `DaprComponentsClient.Get`

```go
ctx := context.TODO()
id := daprcomponents.NewDaprComponentID("12345678-1234-9876-4563-123456789012", "example-resource-group", "managedEnvironmentName", "daprComponentName")

read, err := client.Get(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `DaprComponentsClient.List`

```go
ctx := context.TODO()
id := daprcomponents.NewManagedEnvironmentID("12345678-1234-9876-4563-123456789012", "example-resource-group", "managedEnvironmentName")

// alternatively `client.List(ctx, id)` can be used to do batched pagination
items, err := client.ListComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `DaprComponentsClient.ListSecrets`

```go
ctx := context.TODO()
id := daprcomponents.NewDaprComponentID("12345678-1234-9876-4563-123456789012", "example-resource-group", "managedEnvironmentName", "daprComponentName")

read, err := client.ListSecrets(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```
