
## `github.com/hashicorp/go-azure-sdk/resource-manager/batch/2024-07-01/applications` Documentation

The `applications` SDK allows for interaction with Azure Resource Manager `batch` (API Version `2024-07-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/batch/2024-07-01/applications"
```


### Client Initialization

```go
client := applications.NewApplicationsClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `ApplicationsClient.ApplicationCreate`

```go
ctx := context.TODO()
id := applications.NewApplicationID("12345678-1234-9876-4563-123456789012", "example-resource-group", "batchAccountName", "applicationName")

payload := applications.Application{
	// ...
}


read, err := client.ApplicationCreate(ctx, id, payload)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `ApplicationsClient.ApplicationDelete`

```go
ctx := context.TODO()
id := applications.NewApplicationID("12345678-1234-9876-4563-123456789012", "example-resource-group", "batchAccountName", "applicationName")

read, err := client.ApplicationDelete(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `ApplicationsClient.ApplicationGet`

```go
ctx := context.TODO()
id := applications.NewApplicationID("12345678-1234-9876-4563-123456789012", "example-resource-group", "batchAccountName", "applicationName")

read, err := client.ApplicationGet(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `ApplicationsClient.ApplicationList`

```go
ctx := context.TODO()
id := applications.NewBatchAccountID("12345678-1234-9876-4563-123456789012", "example-resource-group", "batchAccountName")

// alternatively `client.ApplicationList(ctx, id, applications.DefaultApplicationListOperationOptions())` can be used to do batched pagination
items, err := client.ApplicationListComplete(ctx, id, applications.DefaultApplicationListOperationOptions())
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `ApplicationsClient.ApplicationUpdate`

```go
ctx := context.TODO()
id := applications.NewApplicationID("12345678-1234-9876-4563-123456789012", "example-resource-group", "batchAccountName", "applicationName")

payload := applications.Application{
	// ...
}


read, err := client.ApplicationUpdate(ctx, id, payload)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```
